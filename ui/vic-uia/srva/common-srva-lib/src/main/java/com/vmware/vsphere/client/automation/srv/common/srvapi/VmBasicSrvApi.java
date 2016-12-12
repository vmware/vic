/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.Calendar;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang.ArrayUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.pbm.profile.ProfileId;
import com.vmware.vim.binding.pbm.profile.ProfileManager;
import com.vmware.vim.binding.vim.ClusterComputeResource;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Description;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.ResourcePool;
import com.vmware.vim.binding.vim.Task;
import com.vmware.vim.binding.vim.VirtualMachine;
import com.vmware.vim.binding.vim.vm.ConfigSpec;
import com.vmware.vim.binding.vim.vm.DefinedProfileSpec;
import com.vmware.vim.binding.vim.vm.FileInfo;
import com.vmware.vim.binding.vim.vm.ProfileSpec;
import com.vmware.vim.binding.vim.vm.device.VirtualDevice.ConnectInfo;
import com.vmware.vim.binding.vim.vm.device.VirtualDeviceSpec;
import com.vmware.vim.binding.vim.vm.device.VirtualDeviceSpec.Operation;
import com.vmware.vim.binding.vim.vm.device.VirtualE1000;
import com.vmware.vim.binding.vim.vm.device.VirtualE1000e;
import com.vmware.vim.binding.vim.vm.device.VirtualEthernetCard;
import com.vmware.vim.binding.vim.vm.device.VirtualEthernetCard.NetworkBackingInfo;
import com.vmware.vim.binding.vim.vm.device.VirtualLsiLogicController;
import com.vmware.vim.binding.vim.vm.device.VirtualSCSIController.Sharing;
import com.vmware.vim.binding.vim.vm.device.VirtualVmxnet3;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NicSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NicSpec.AddressType;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * Provides utility methods for performing operations on VM via API.
 */
public class VmBasicSrvApi {

   private static final Logger _logger = LoggerFactory
         .getLogger(VmBasicSrvApi.class);

   private static final String DEFAULT_VM_GUEST_OS = "winVistaGuest";
   private static final int DEVICE_KEY_MAGIC_NUMBER = 105;

   private static VmBasicSrvApi instance = null;

   protected VmBasicSrvApi() {
   }

   /**
    * Get instance of VmSrvApi.
    *
    * @return created instance
    */
   public static VmBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (VmBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing VmSrvApi.");
               instance = new VmBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a VM with properties as specified in the <code>VmSpec</code>. This
    * method handles properly situations when VM resides on standalone host and
    * such inside a cluster. NOTE: Due to PR 1384397 the VM jhas no disk. NOTE:
    * this method does not support creating VM inside specified folder. NOTE:
    * the VM will be create on the first host datastore
    *
    * @param vmSpec
    *           VM spec that holds all VM properties
    * @return true if creation was successful and false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   // TODO: rreymer fix PR 1384397
   public boolean createVm(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Creating VM '%s'", vmSpec.name.get()));

      // NOTE: VM cannot be inside a folder only directly on host
      HostSpec hostSpec = (HostSpec) vmSpec.parent.get();
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

      // Compute resource for the vm
      ManagedObjectReference computeResourceMor = getVmComputeResourcePool(
            vmSpec, host);
      if (computeResourceMor == null) {
         new RuntimeException(String.format(
               "Compute resource for the VM %s is null!", vmSpec.name.get()));
      }

      // Getting the datastores available.
      ManagedObjectReference[] hostDatastoresMor = host.getDatastore();
      Map<String, Datastore> datastoresByNames = new HashMap<String, Datastore>(
            hostDatastoresMor.length);
      for (ManagedObjectReference datastoreMor : hostDatastoresMor) {
         Datastore datastore = VcServiceUtil.getVcService(vmSpec)
               .getManagedObject(datastoreMor);
         datastoresByNames.put(datastore.getName(), datastore);
      }

      FileInfo fileInfo = null;
      Datastore vmDatastore = null;

      // Check for preferred datastores.
      if (vmSpec.datastore.isAssigned()) {
         vmDatastore = datastoresByNames.get(vmSpec.datastore.get().name.get());
         if (vmDatastore == null) {
            throw new IllegalArgumentException(String.format(
                  "Failed to find datastore with name %s for VM with name %s",
                  vmSpec.datastore.get().name.get(), vmSpec.name.get()));
         }
      } else { // If not just use the first in the list.
         vmDatastore = (Datastore) datastoresByNames.values().toArray()[0];
      }

      fileInfo = new FileInfo(String.format("[%s]", vmDatastore.getName()),
            null, null, null, null);

      // check for storage policy
      ProfileSpec[] profiles = null;
      if (vmSpec.profile != null && vmSpec.profile.getAll().size() > 0) {
         profiles = new ProfileSpec[vmSpec.profile.getAll().size()];

         for (int i = 0; i < vmSpec.profile.getAll().size(); i++) {
            StoragePolicySpec policySpec = vmSpec.profile.getAll().get(i);
            final ProfileManager profileManager = StoragePolicyBasicSrvApi
                  .getInstance().getProfileManager(policySpec);
            if (!policySpec.name.isAssigned()) {
               throw new IllegalArgumentException(
                     "Invalid spec - policy has no name.");
            }
            final ProfileId profileId = StoragePolicyBasicSrvApi.getInstance()
                  .getProfileByName(policySpec.name.get(), profileManager);
            if (profileId == null) {
               throw new IllegalArgumentException("ProfileId is null.");
            }

            DefinedProfileSpec definedProfile = new DefinedProfileSpec();
            definedProfile.profileId = profileId.getUniqueId();
            profiles[i] = definedProfile;
         }
      }

      ConfigSpec vmConfigSpec = new ConfigSpec();
      vmConfigSpec.setName(vmSpec.name.get());
      vmConfigSpec.setFiles(fileInfo);

      // assign storage policy
      if (profiles != null) {
         vmConfigSpec.setVmProfile(profiles);
      }

      if (vmSpec.hardwareVersion.isAssigned()) {
         vmConfigSpec.setVersion(vmSpec.hardwareVersion.get());
      }

      if (vmSpec.guestId.isAssigned()) {
         vmConfigSpec.setGuestId(vmSpec.guestId.get());
      } else {
         // set default guest ID in order to be able to power on the VM
         vmConfigSpec.setGuestId(DEFAULT_VM_GUEST_OS);
      }

      if (vmSpec.memoryInMB.isAssigned()) {
         vmConfigSpec.setMemoryMB(vmSpec.memoryInMB.get());
      }

      if (vmSpec.numCPUs.isAssigned()) {
         vmConfigSpec.setNumCPUs(vmSpec.numCPUs.get());
      }

      DatacenterSpec datacenterSpec = null;
      ManagedEntitySpec hostParentSpec = hostSpec.parent.get();
      if (hostParentSpec instanceof DatacenterSpec) {
         datacenterSpec = (DatacenterSpec) hostParentSpec;
      } else if (hostParentSpec instanceof ClusterSpec) {
         datacenterSpec = (DatacenterSpec) hostParentSpec.parent.get();
      }

      Folder vmFolder = null;
      if (vmSpec.vmFolder.isAssigned()) {
         vmFolder = ManagedEntityUtil.getManagedObject(vmSpec.vmFolder.get());
      } else {
         vmFolder = FolderBasicSrvApi.getInstance().getVmFolder(datacenterSpec);
      }

      // Create the VM with a NIC if it's specified in the spec.
      if (vmSpec.nicList.isAssigned()) {
         vmConfigSpec.deviceChange = toVirtualDeviceSpec(vmSpec.nicList
               .getAll());
      }

      // Add SCSI controller
      vmConfigSpec.deviceChange = (VirtualDeviceSpec[]) ArrayUtils.addAll(
            vmConfigSpec.deviceChange, new VirtualDeviceSpec[1]);

      int deviceKey = getDeviceKey(vmConfigSpec.deviceChange.length - 1);
      VirtualDeviceSpec vdsController = buildSCSIControllerDeviceSpec(deviceKey);
      vmConfigSpec.deviceChange[vmConfigSpec.deviceChange.length - 1] = vdsController;

      // Start the task
      ManagedObjectReference createTaskMor = vmFolder.createVm(vmConfigSpec,
            computeResourceMor, host._getRef());

      boolean taskSuccess = VcServiceUtil.waitForTaskSuccess(createTaskMor,
            vmSpec);

      // Workaround for waitForTaskSuccess() returning true prematurely,
      // before the create VM task is completed - PR 1423973
      // TODO: rkovachev remove after PR 1423973 is fixed
      VcService service = VcServiceUtil.getVcService(vmSpec);
      waitForTaskTimeCompleted(service, createTaskMor);
      return taskSuccess;
   }

   /**
    * Deletes VM from inventory.
    *
    * @param vmSpec
    *           the spec of the VM that will be destroyed
    * @return true if deletion was successful, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean deleteVm(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Deleting VM '%s'", vmSpec.name.get()));

      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);
      ManagedObjectReference taskMoRef = virtualMachine.destroy();

      boolean taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vmSpec);

      // Workaround for waitForTaskSuccess() returning true prematurely,
      // before the create VM task is completed - PR 1423973
      // TODO: rkovachev remove after PR 1423973 is fixed
      VcService service = VcServiceUtil.getVcService(vmSpec);
      waitForTaskTimeCompleted(service, taskMoRef);
      return taskResult;
   }

   /**
    * Checks whether the specified VM exists.
    *
    * @param vmSpec
    *           the spec of the VM that will be checked for existence
    * @return true if VM is found, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean checkVmExists(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Checking whether VM '%s' exists",
            vmSpec.name.get()));

      try {
         getVirtualMachine(vmSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }

      return true;
   }

   /**
    * Checks whether the specified VM exists and deletes it.
    *
    * @param vmSpec
    *           the spec of the VM that will be destroyed
    * @return true if deletion is successful or VM is missing, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean deleteVmSafely(VmSpec vmSpec) throws Exception {
      if (checkVmExists(vmSpec)) {
         return deleteVm(vmSpec);
      }

      return true;
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Retrieves virtual machine from specified spec object.
    *
    * @param vmSpec
    *           spec of the VM that will be retrieved
    * @return retrieved VM object
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   private VirtualMachine getVirtualMachine(VmSpec vmSpec) throws Exception {
      return ManagedEntityUtil.getManagedObject(vmSpec,
            vmSpec.service.isAssigned() ? vmSpec.service.get() : null);
   }

   /**
    * Method that validates the VM spec: - name should be assigned - parent
    * should be assigned
    *
    * @param vmSpec
    *           VM spec that will be validated
    * @throws IllegalArgumentException
    *            if vm spec requirements are not met
    */
   private void validateVmSpec(VmSpec vmSpec) throws IllegalArgumentException {
      if (!vmSpec.name.isAssigned() || vmSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("VM name is not set.");
      }

      if (!vmSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("VM parent is not set.");
      }

      ManagedEntitySpec vmParentEntity = vmSpec.parent.get();
      if (!(vmParentEntity instanceof HostSpec)
            && !(vmParentEntity instanceof ClusterSpec)
            && !(vmParentEntity instanceof DatacenterSpec)
            && !(vmParentEntity instanceof FolderSpec)
            && !(vmParentEntity instanceof VappSpec)
            && !(vmParentEntity instanceof ResourcePoolSpec)) {
         throw new IllegalArgumentException(
               "VM parent association is not a host, cluster, datacenter, "
                     + "folder, vApp or resource pool.");
      }
   }

   /**
    * Populates properties from NicSpec to VirtualDeviceSpec. If spcefic
    * properties are not given, default ones will be used:
    * <ol>
    * <li>operation=add</li>
    * <li>wakeOnWan=true</li>
    * <li>status='untried'</li>
    * <li>startConnected=true</li>
    * <li>connected=true</li>
    * <li>allowGuestControl=true</li>
    * <li>addressType='Generated'</li>
    * <li>macAddress=""</li>
    * <li>deviceName='Vm network'</li>
    * </ol>
    *
    * @param nicList
    *           list of NicSpec objects found in a VmSpec.
    * @return Array of VirtualDeviceSpec objects - one for each NicSpec.
    */
   private VirtualDeviceSpec[] toVirtualDeviceSpec(List<NicSpec> nicList) {
      VirtualDeviceSpec[] aVds = new VirtualDeviceSpec[nicList.size()];
      for (int i = 0; i < aVds.length; i++) {
         NicSpec nicSpec = nicList.get(i);

         VirtualEthernetCard nic = createEthernetCard(nicSpec);
         VirtualDeviceSpec vds = new VirtualDeviceSpec();
         aVds[i] = vds;

         vds.operation = Operation.add;
         nic.wakeOnLanEnabled = true;
         ConnectInfo connInfo = new ConnectInfo();

         connInfo.status = "untried";
         connInfo.startConnected = nicSpec.startConnected.isAssigned() ? nicSpec.startConnected
               .get() : true;
         connInfo.connected = nicSpec.connected.isAssigned() ? nicSpec.connected
               .get() : true;
         connInfo.allowGuestControl = true;
         nic.connectable = connInfo;

         Description desc = new Description();
         desc.summary = "";
         desc.label = "New network - " + nicSpec.name;
         nic.deviceInfo = desc;

         AddressType nicAddressType = nicSpec.addressType.isAssigned() ? nicSpec.addressType
               .get() : AddressType.GENERATED;
         nic.addressType = nicAddressType.value();

         nic.macAddress = nicSpec.macAddress.isAssigned() ? nicSpec.macAddress
               .get() : "";

         NetworkBackingInfo bi = new NetworkBackingInfo();
         bi.deviceName = nicSpec.deviceName.isAssigned() ? nicSpec.deviceName
               .get() : "VM Network";
         nic.backing = bi;

         nic.key = getDeviceKey(i);

         vds.device = nic;
      }
      return aVds;
   }

   /**
    * Creates and returns new {@code VirtualEthernetCard} instance based on the
    * device type of the NIC
    *
    * @param nicSpec
    *           the spec for the network adapter to be created
    * @return the newly instantiated {@code VirtualEthernetCard}
    */
   private VirtualEthernetCard createEthernetCard(NicSpec nicSpec) {
      if (nicSpec.adapterType.isAssigned()) {
         switch (nicSpec.adapterType.get()) {
         case E1000:
            return new VirtualE1000();
         case E1000E:
            return new VirtualE1000e();
         case VMXNET3:
            return new VirtualVmxnet3();
         }
      }

      // return E1000e adapter by defaut
      return new VirtualE1000e();
   }

   /**
    * Builds an SCSI controller device spec to be added to a VM config. The
    * controller is configured with no sharing.
    *
    * @param controllerKey
    *           - the controller device hey
    * @return the spec we build
    */
   private VirtualDeviceSpec buildSCSIControllerDeviceSpec(int controllerKey) {
      VirtualLsiLogicController controller = new VirtualLsiLogicController();
      controller.setSharedBus(Sharing.noSharing);
      controller.setKey(controllerKey);

      VirtualDeviceSpec result = new VirtualDeviceSpec();
      result.operation = Operation.add;
      result.device = controller;
      return result;
   }

   /**
    * Gets a key based on the index of the device in the list of all devices
    *
    * @param deviceIndex
    *           - the index of the device in the list of all devices
    * @return a key based on the index of the device in the list of all devices
    */
   private int getDeviceKey(int deviceIndex) {
      return deviceIndex - DEVICE_KEY_MAGIC_NUMBER;
   }

   /**
    * Wait for the task to get time completed property populated for the
    * BackendDelay.LARGE timeout. TODO: rkovachev fix PR 1423973
    *
    * @param service
    * @param tasksMor
    */
   private void waitForTaskTimeCompleted(VcService service,
         ManagedObjectReference tasksMor) {
      // Workaround for waitForTaskSuccess() returning true prematurely,
      // before the create VM task is completed - PR 1423973
      // TODO: remove after PR 1423973 is fixed
      Task task = service.getManagedObject(tasksMor);
      Calendar timeCompleted = task.getInfo().getCompleteTime();
      int retry = 0;
      while ((timeCompleted == null)
            && (retry < BackendDelay.LARGE.getDuration())) {
         timeCompleted = task.getInfo().getCompleteTime();
         try {
            Thread.sleep(1000);
         } catch (InterruptedException e) {
            _logger.info(e.getMessage());
         }
         retry++;
      }
   }

   /**
    * Get the compute resource for a vm - either resource pool, cluster or
    * default compute resource
    *
    * @param vmSpec
    * @param host
    * @return mor of the compute resource
    * @throws Exception
    */
   private ManagedObjectReference getVmComputeResourcePool(VmSpec vmSpec,
         HostSystem host) throws Exception {

      ManagedObjectReference computeResourceMor = null;
      if (vmSpec.computeResource.isAssigned()) {
         if (vmSpec.computeResource.get() instanceof ResourcePoolSpec) {
            ResourcePool resPool = ManagedEntityUtil
                  .getManagedObject(vmSpec.computeResource.get());
            computeResourceMor = resPool._getRef();
         } else if (vmSpec.computeResource.get() instanceof ClusterSpec) {
            ClusterComputeResource clusterCompRes = ManagedEntityUtil
                  .getManagedObject(vmSpec.computeResource.get());
            computeResourceMor = clusterCompRes.getResourcePool();
         } else if (vmSpec.computeResource.get() instanceof HostSpec) {
            ComputeResource computeResource = VcServiceUtil
                  .getVcService(vmSpec).getManagedObject(host.getParent());
            computeResourceMor = computeResource.getResourcePool();
         } else {
            throw new RuntimeException(
                  "Unsupported type of compure resource. VM cannot be created under:"
                        + vmSpec.computeResource.get().name.get());
         }
      } else {
         ComputeResource computeResource = VcServiceUtil.getVcService(vmSpec)
               .getManagedObject(host.getParent());
         computeResourceMor = computeResource.getResourcePool();
      }

      return computeResourceMor;

   }

   /**
    * Reconfigure VM. NOTE: for now only VM name changes. If other params should
    * be introduced they have to be added at the marked place in the method.
    *
    * @param reconfigureVmSpec
    *           VM that will be reconfigured
    * @return true if reconfigure operation is successful, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean reconfigureVm(ReconfigureVmSpec reconfigureVmSpec)
         throws Exception {
      VmSpec targetVm = reconfigureVmSpec.targetVm.get();
      VmSpec newVmConfigs = reconfigureVmSpec.newVmConfigs.get();
      validateVmSpec(targetVm);

      _logger.info(String.format(
            "Reconfiguring VM '%s', it's new name is: '%s'",
            targetVm.name.get(), newVmConfigs.name.get()));

      VirtualMachine vm = getVirtualMachine(targetVm);

      ConfigSpec reconfigureSpec = new ConfigSpec();
      reconfigureSpec.setName(newVmConfigs.name.get());

      // TODO: add other VM params that could be reconfigured

      ManagedObjectReference taskMoRef = vm.reconfigure(reconfigureSpec);
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, targetVm);
   }
}
