/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.net.URL;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Description;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.HttpNfcLease;
import com.vmware.vim.binding.vim.HttpNfcLease.DeviceUrl;
import com.vmware.vim.binding.vim.HttpNfcLease.Info;
import com.vmware.vim.binding.vim.HttpNfcLease.State;
import com.vmware.vim.binding.vim.ManagedEntity;
import com.vmware.vim.binding.vim.OvfManager;
import com.vmware.vim.binding.vim.OvfManager.CreateImportSpecParams;
import com.vmware.vim.binding.vim.OvfManager.CreateImportSpecResult;
import com.vmware.vim.binding.vim.OvfManager.FileItem;
import com.vmware.vim.binding.vim.OvfManager.ParseDescriptorParams;
import com.vmware.vim.binding.vim.OvfManager.ParseDescriptorResult;
import com.vmware.vim.binding.vim.ResourcePool;
import com.vmware.vim.binding.vim.Task;
import com.vmware.vim.binding.vim.VirtualMachine;
import com.vmware.vim.binding.vim.VirtualMachine.MovePriority;
import com.vmware.vim.binding.vim.VirtualMachine.PowerState;
import com.vmware.vim.binding.vim.vm.FaultToleranceConfigSpec;
import com.vmware.vim.binding.vim.vm.FaultToleranceMetaSpec;
import com.vmware.vim.binding.vim.vm.FaultToleranceVMConfigSpec;
import com.vmware.vim.binding.vim.vm.FaultToleranceVMConfigSpec.FaultToleranceDiskSpec;
import com.vmware.vim.binding.vim.vm.ConfigInfo;
import com.vmware.vim.binding.vim.vm.ConfigSpec;
import com.vmware.vim.binding.vim.vm.RelocateSpec;
import com.vmware.vim.binding.vim.vm.RuntimeInfo;
import com.vmware.vim.binding.vim.vm.device.VirtualDevice;
import com.vmware.vim.binding.vim.vm.VirtualHardware;
import com.vmware.vim.binding.vim.vm.device.VirtualDevice.ConnectInfo;
import com.vmware.vim.binding.vim.vm.device.VirtualDevice;
import com.vmware.vim.binding.vim.vm.device.VirtualDeviceSpec;
import com.vmware.vim.binding.vim.vm.device.VirtualDeviceSpec.FileOperation;
import com.vmware.vim.binding.vim.vm.device.VirtualDeviceSpec.Operation;
import com.vmware.vim.binding.vim.vm.device.VirtualDisk;
import com.vmware.vim.binding.vim.vm.device.VirtualDisk.FlatVer2BackingInfo;
import com.vmware.vim.binding.vim.vm.device.VirtualE1000e;
import com.vmware.vim.binding.vim.vm.device.VirtualEthernetCard.NetworkBackingInfo;
import com.vmware.vim.binding.vim.vm.device.VirtualLsiLogicController;
import com.vmware.vim.binding.vim.vm.device.VirtualSCSIController.Sharing;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NicSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VirtualDiskSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NicSpec.AddressType;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

/**
 * Provides utility methods for performing operations on VM via API.
 */
public class VmSrvApi extends VmBasicSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(VmSrvApi.class);
   private static final int DEVICE_KEY_MAGIC_NUMBER = 105;
   private static final long MB_IN_BYTES = 1024 * 1024;

   private static VmSrvApi instance = null;
   protected VmSrvApi() {}

   /**
    * Get instance of VmSrvApi.
    *
    * @return  created instance
    */
   public static VmSrvApi getInstance() {
      if (instance == null) {
         synchronized(VmSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing VmSrvApi.");
               instance = new VmSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Powers on virtual machine.
    *
    * @param vmSpec
    *           virtual machine specification that will be powered on
    * @return true if the operation was successful, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean powerOnVm(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Powering on VM '%s'", vmSpec.name.get()));

      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);
      ManagedObjectReference taskMoRef = virtualMachine.powerOn(null);

      boolean taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vmSpec);

      // Workaround for waitForTaskSuccess() returning true prematurely,
      // before the power on VM task is completed - PR 1423973
      // TODO: remove after PR 1423973 is fixed
      VcService service = VcServiceUtil.getVcService(vmSpec);
      Task powerOmVMTask = service.getManagedObject(taskMoRef);
      Calendar resultMor = powerOmVMTask.getInfo().getCompleteTime();
      int retry = 0;
      while ((resultMor == null) && (retry < BackendDelay.LARGE.getDuration())) {
         resultMor = powerOmVMTask.getInfo().getCompleteTime();
         Thread.sleep(1000);
         retry++;
      }

      return taskResult;
   }

   /**
    * Powers off virtual machine.
    *
    * @param vmSpec
    *           virtual machine specification that will be powered off
    * @return true if the operation was successful, false otherwise
    * @throws Exceptionif
    *            login to the VC service fails, or the specified VM settings are
    *            invalid
    */
   public boolean powerOffVm(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Powering off VM '%s'", vmSpec.name.get()));

      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);
      ManagedObjectReference taskMoRef = virtualMachine.powerOff();

      boolean taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vmSpec);
      // TODO: Fix it lgrigorova
      CommonUtils.sleep(10000);
      return taskResult;
   }

   /**
    * Suspends a VM.
    * @param vmSpec Spec representing the VM to be suspended.
    * @return true if VM was suspended successfully, false otherwise.
    */
   public boolean suspendVm(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Suspending VM '%s'", vmSpec.name.get()));

      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);
      ManagedObjectReference taskMoRef = virtualMachine.suspend();

      boolean taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vmSpec);
      return taskResult;
   }

   /**
    * Checks the power state of a VM and waits for it to become powered on or off. This
    * method is useful in order to verify the result of a power-on/off operation,
    * e.g. after calling powerOnVm().
    *
    * The method accounts for the fact that the state of a VM gets updated to
    * powered-on some time after the completion of the task.
    *
    * @param vmSpec
    *           the spec of the VM to check the state of
    * @param poweredOn true if should wait for VM to become powered on, false if powered off
    * @return <b>true</b> if the status of the VM becomes as expected in poweredOn,
    *         <b>false</b> otherwise
    * @throws Exception
    */
   public boolean waitForVmPowerState(VmSpec vmSpec, boolean poweredOn)
         throws Exception {
      long endTime = System.currentTimeMillis()
            + SUITA.Environment.getBackendJobSmall();
      while (isVmPoweredOn(vmSpec) != poweredOn) {
         if (System.currentTimeMillis() > endTime) {
            return false;
         }
         Thread.sleep(500);
      }
      return true;
   }


   /**
    * Check if the current power state of the VM is 'Powered On'.
    *
    * @param vmSpec
    *           the spec of the VM which status will be checked
    * @return true if the VM is powered on, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean isVmPoweredOn(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Checking if VM '%s' is powered on", vmSpec.name.get()));

      PowerState powerState = getVmPowerState(vmSpec);
      return powerState.equals(PowerState.poweredOn) ? true : false;
   }

   /**
    * Check if the current power state of the VM is 'Powered Off'.
    *
    * @param vmSpec
    *           the spec of the VM which status will be checked
    * @return true if the VM is powered off, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   public boolean isVmPoweredOff(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger
      .info(String.format("Checking if VM '%s' is powered off", vmSpec.name.get()));

      PowerState powerState = getVmPowerState(vmSpec);
      return powerState.equals(PowerState.poweredOff) ? true : false;
   }

   /**
    * Converts a VM to VM template.
    * @param vmSpec
    * @throws Exception
    */
   public void convertToTemplate(VmSpec vmSpec) throws Exception {
      validateVmSpec(vmSpec);

      _logger.info(String.format("Converting VM to template '%s'", vmSpec.name.get()));

      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);
      virtualMachine.markAsTemplate();
   }

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
   public VirtualMachine getVirtualMachine(VmSpec vmSpec) throws Exception {
      return ManagedEntityUtil.getManagedObject(
            vmSpec,
            vmSpec.service.isAssigned() ? vmSpec.service.get() : null
         );
   }

   /**
    * Migrates a VM to a different host and/or datastore
    * @param vmSpec     the VM to migrate
    * @param host       the new host. Set to null to keep it on the current host.
    * @param datastore  dest datastore. Set to null to keep it on the current datastore
    * @return           true if the migration succeeds
    * @throws Exception if there is an error when finding the host/datastore
    */
   public boolean migrateVm(VmSpec vmSpec, HostSpec host, DatastoreSpec datastore)
         throws Exception {
      validateVmSpec(vmSpec);

      RelocateSpec relocateSpec = new RelocateSpec();

      // Host can be null
      if (host != null) {
         ManagedObjectReference hostMoRef =
               ManagedEntityUtil.getManagedObject(host)._getRef();
         relocateSpec.host = hostMoRef;
      }

      // Datastore can be null
      if (datastore != null) {
         ManagedObjectReference datastoreMoRef =
               ManagedEntityUtil.getManagedObject(datastore)._getRef();
         relocateSpec.datastore = datastoreMoRef;
      }

      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);

      ManagedObjectReference taskMoRef =
            virtualMachine.relocate(relocateSpec, MovePriority.defaultPriority);

      return VcServiceUtil.waitForTaskSuccess(taskMoRef, vmSpec);
   }

   /**
    * Gets the datastore used by the VM.
    * The VM and the datastore should have the same parent which is a host.
    * @param vm - the VM to query for
    * @return the datastore used by the VM
    * @throws Exception - Exception if there is an error when finding the VM
    */
   public DatastoreSpec getVmDatastore(VmSpec vm) throws Exception {
      validateVmSpec(vm);
      ManagedEntitySpec vmParent = evaluateVmParent(vm);

      ServiceSpec serviceSpec = vm.service.get();

      VirtualMachine virtualMachine = getVirtualMachine(vm);
      ManagedObjectReference datastoreMoRef = virtualMachine.getDatastore()[0];

      ManagedEntity datastoreME = ManagedEntityUtil.getManagedObjectFromMoRef(datastoreMoRef, serviceSpec);

      DatastoreSpec datastore = SpecFactory.getSpec(DatastoreSpec.class, datastoreME.getName(), vmParent);

      return datastore;
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
   public void validateVmSpec(VmSpec vmSpec) throws IllegalArgumentException {
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
    * Imports VM from OVF.
    * The location of the OVF can be specified as URL.
    * Supported protocols are: HTTP, HTTPS, FTP.
    *
    * @param ovfUrl     the URL of the OVF that will be imported
    * @param vmSpec     the spec of the VM that will be imported, only name is used
    * @return           true if creation was successful and false otherwise
    * @throws Exception if login to the VC service fails, or the specified VM settings
    *                   are invalid
    */
   public boolean importVmFromOvf(final String ovfUrl, VmSpec vmSpec)
         throws Exception {
      validateVmSpec(vmSpec);
      validateOvfUrl(ovfUrl);

      _logger.info(String.format("Importing VM '%s'", vmSpec.name.get()));

      VcService vcService = VcServiceUtil.getVcService(vmSpec);

      ManagedObjectReference ovfManagerMor =
            vcService.getServiceInstanceContent().getOvfManager();
      OvfManager ovfManager = vcService.getManagedObject(ovfManagerMor);

      ParseDescriptorParams parseDescParams = new ParseDescriptorParams();
      parseDescParams.setDeploymentOption("");
      parseDescParams.setLocale("");

      String ovfFile = HttpConnector.getFileFromUrl(ovfUrl);
      ParseDescriptorResult parseDescriptor =
            ovfManager.parseDescriptor(ovfFile, parseDescParams);

      if (parseDescriptor.getError() != null && parseDescriptor.getError().length > 0) {
         throw new IllegalArgumentException("The passed OVF cannot be parsed!");
      }

      HostSpec hostSpec = (HostSpec) vmSpec.parent.get();
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

      ManagedObjectReference datastoreMor = null;
      ManagedObjectReference[] datastoreMorArr = host.getDatastore();

      // Check if there's a preferred datastore from the VmSpec.
      if (vmSpec.datastore.isAssigned() && datastoreMorArr.length > 1) {
         ManagedObjectReference datastoreReference =
               ManagedEntityUtil.getManagedObject(vmSpec.datastore.get())._getRef();
         String datastoreName = datastoreReference.getValue();
         for (ManagedObjectReference mor : datastoreMorArr) {
            if (mor.getValue().equals(datastoreName)) {
               datastoreMor = mor;
               break;
            }
         }
      }

      if (datastoreMor == null) {
         datastoreMor = host.getDatastore()[0];
      }

      CreateImportSpecParams createImportSpecParams = new CreateImportSpecParams();
      createImportSpecParams.setEntityName(vmSpec.name.get());
      createImportSpecParams.setLocale("");
      createImportSpecParams.setDeploymentOption("");
      createImportSpecParams.setHostSystem(host._getRef());

      ComputeResource computeResource = vcService.getManagedObject(host.getParent());
      ManagedObjectReference resourcePoolMor = computeResource.getResourcePool();
      final CreateImportSpecResult createImportSpec = ovfManager.createImportSpec(
            ovfFile,
            resourcePoolMor,
            datastoreMor,
            createImportSpecParams
            );

      DatacenterSpec datacenterSpec = null;
      ManagedEntitySpec hostParentSpec = hostSpec.parent.get();
      if (hostParentSpec instanceof DatacenterSpec) {
         datacenterSpec = (DatacenterSpec) hostParentSpec;
      } else if (hostParentSpec instanceof ClusterSpec) {
         datacenterSpec = (DatacenterSpec) hostParentSpec.parent.get();
      }

      ResourcePool resourcePool = vcService.getManagedObject(resourcePoolMor);
      ManagedObjectReference httpNfcLeaseMor = resourcePool.importVApp(
            createImportSpec.getImportSpec(),
            FolderBasicSrvApi.getInstance().getVmFolder(datacenterSpec)._getRef(),
            host._getRef()
            );

      Boolean uploaded = false;
      HttpNfcLease httpNfcLease = vcService.getManagedObject(httpNfcLeaseMor);
      if (httpNfcLease != null) {
         try {
            while (httpNfcLease.getState() == State.initializing) {
               Thread.sleep(2000);
            }

            if (httpNfcLease.getState() == State.ready) {
               final Info httpNfcLeaseInfo = httpNfcLease.getInfo();
               uploaded = uploadAppliance(createImportSpec, httpNfcLeaseInfo, ovfUrl);

               if (uploaded) {
                  _logger.info("Uploading of the entity is successful");
                  httpNfcLease.complete();
               } else {
                  _logger.error("Failed to upload the entity");
                  httpNfcLease.abort(null);
               }
            } else if (httpNfcLease.getState() == State.error) {
               Exception fault = httpNfcLease.getError();
               if (fault != null) {
                  _logger.error("Retrieved Fault from the HttpNfcLease - "
                        + fault.getClass());
                  throw new RuntimeException(fault.getMessage());
               }
            }
         } catch (Exception faultException) {
            State importState = httpNfcLease.getState();
            if (importState != State.done && importState != State.error) {
               httpNfcLease.abort(null);
               throw faultException;
            } else if (importState == State.error) {
               throw faultException;
            }
         }
      }

      return uploaded;
   }

   /**
    * Turn off Fault Tolerance
    *
    * @param vmSpec
    * @return true if successfully disabled, false otherwise
    * @throws Exception
    */
   public boolean turnOffFaultTolerance(VmSpec vmSpec) throws Exception {
      VirtualMachine primaryVM = getVirtualMachine(vmSpec);
      boolean taskResult = VcServiceUtil.waitForTaskSuccess(
            primaryVM.turnOffFaultTolerance(), vmSpec);
      return taskResult;
   }

   /**
    * Gets current VM's host name
    *
    * @param vm - spec for the VM
    * @return host name
    * @throws Exception
    */
   public String getVmHost(VmSpec vm) throws Exception {
      validateVmSpec(vm);

      evaluateVmParent(vm);

      ServiceSpec serviceSpec = vm.service.get();

      VirtualMachine virtualMachine = getVirtualMachine(vm);
      ManagedObjectReference hostMoRef = virtualMachine.getRuntime().getHost();

      ManagedEntity hostManagedEntity =
            ManagedEntityUtil.getManagedObjectFromMoRef(hostMoRef, serviceSpec);

      return hostManagedEntity.getName();
   }

   /**
    * Populates properties from NicSpec to VirtualDeviceSpec.
    * If specific properties are not given, default ones will be used:
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
    * @param nicList list of NicSpec objects found in a VmSpec.
    * @return Array of VirtualDeviceSpec objects - one for each NicSpec.
    */
   public VirtualDeviceSpec[] toVirtualDeviceSpec(List<NicSpec> nicList) {
      VirtualDeviceSpec[] aVds = new VirtualDeviceSpec[nicList.size()];
      for (int i = 0; i < aVds.length; i++) {
         VirtualE1000e nic = new VirtualE1000e();
         VirtualDeviceSpec vds = new VirtualDeviceSpec();
         aVds[i] = vds;

         NicSpec nicSpec = nicList.get(i);

         vds.operation = Operation.add;
         nic.wakeOnLanEnabled = true;
         ConnectInfo connInfo = new ConnectInfo();

         connInfo.status = "untried";
         connInfo.startConnected = nicSpec.startConnected.isAssigned() ? nicSpec.startConnected.get() : true;
         connInfo.connected = nicSpec.connected.isAssigned() ? nicSpec.connected.get() : true;
         connInfo.allowGuestControl = true;
         nic.connectable = connInfo;

         Description desc = new Description();
         desc.summary = "";
         desc.label = "New network - " + nicSpec.name;
         nic.deviceInfo = desc;

         AddressType nicAddressType = nicSpec.addressType.isAssigned() ? nicSpec.addressType.get() : AddressType.GENERATED;
         nic.addressType = nicAddressType.value();

         nic.macAddress = nicSpec.macAddress.isAssigned() ? nicSpec.macAddress.get() : "";

         NetworkBackingInfo bi = new NetworkBackingInfo();
         bi.deviceName = nicSpec.deviceName.isAssigned() ? nicSpec.deviceName.get() : "VM Network";
         nic.backing = bi;

         nic.key = getDeviceKey(i);

         vds.device = nic;
      }

      return aVds;
   }

   /**
    * Builds an SCSI controller device spec to be added to a VM config. The controller is configured with no sharing.
    * @param controllerKey - the controller device hey
    * @return the spec we build
    */
   public VirtualDeviceSpec buildSCSIControllerDeviceSpec(int controllerKey) {
      VirtualLsiLogicController controller = new VirtualLsiLogicController();
      controller.setSharedBus(Sharing.noSharing);
      controller.setKey(controllerKey);

      VirtualDeviceSpec result = new VirtualDeviceSpec();
      result.operation = Operation.add;
      result.device = controller;
      return result;
   }

   /**
    * Builds an HDD device spec to be added to a VM config. The HDD is thin provisioned.
    * @param datastore - the datastore where the HDD will be
    * @param sizeInMb - the HDD size in MB
    * @param controllerKey - the controller device key
    * @param hddKey - the HDD device hey
    * @return the spec we build
    */
   public VirtualDeviceSpec buildHddDeviceSpec(Datastore datastore, long sizeInMb, String fileName,
         int controllerKey, int hddKey) {
      VirtualDisk hdd = new VirtualDisk();

      FlatVer2BackingInfo hddBackingInfo = new FlatVer2BackingInfo();
      hddBackingInfo.setDiskMode("persistent");
      hddBackingInfo.setThinProvisioned(true);
      hddBackingInfo.setFileName(fileName);
      hdd.setBacking(hddBackingInfo);

      hdd.setCapacityInBytes(MB_IN_BYTES * sizeInMb);
      hdd.setUnitNumber(-1);
      hdd.setKey(hddKey);
      hdd.setControllerKey(controllerKey);

      VirtualDeviceSpec result = new VirtualDeviceSpec();
      result.operation = Operation.add;
      result.device = hdd;
      return result;
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Gets VM power state.
    *
    * @param vmSpec
    *           VM spec which power state will be retrieved
    * @return current power state of the VM
    * @throws Exception
    *            if login to the VC service fails, or the specified VM settings
    *            are invalid
    */
   private PowerState getVmPowerState(VmSpec vmSpec) throws Exception {
      VirtualMachine virtualMachine = getVirtualMachine(vmSpec);
      RuntimeInfo runtime = virtualMachine.getRuntime();
      return runtime.getPowerState();
   }

   /**
    * Validates the OVF URL that was passed.
    * For now only supported protocols are HTTP, HTTPS, FTP.
    *
    * @param ovfUrl     OVF URL to be validated
    */
   private void validateOvfUrl(String ovfUrl) {
      if (!ovfUrl.startsWith("http://") && !ovfUrl.startsWith("https://") &&
            !ovfUrl.startsWith("ftp://")) {
         throw new IllegalArgumentException(
               "The OVF url does not start with http, https or ftp"
               );
      }
   }

   /**
    * Uploads appliance based on the provided import spec and import lease.
    *
    * @param importSpecRes    import spec that is used to retrieve files from
    * @param httpNfcLeaseInfo import lease info that provides destination URLs
    * @param ovfPath          path to the source files
    * @return                 if the upload was successful
    * @throws Exception       if something wrong happens during upload
    */
   private boolean uploadAppliance(CreateImportSpecResult importSpecRes,
         Info httpNfcLeaseInfo, String ovfPath) throws Exception {
      validateOvfUrl(ovfPath);

      boolean uploaded = true;
      if (importSpecRes != null && httpNfcLeaseInfo != null && ovfPath != null) {

         FileItem[] fileItems = importSpecRes.getFileItem();
         DeviceUrl[] urls = httpNfcLeaseInfo.getDeviceUrl();
         String basePath = ovfPath.substring(0, ovfPath.lastIndexOf("/")) + "/";

         _logger.info("Upload started.....");

         if (fileItems != null && urls != null && fileItems.length == urls.length) {

            for (FileItem fileItem : fileItems) {
               boolean present = false;
               for (DeviceUrl deviceUrl : urls) {
                  if (fileItem.getDeviceId().equals(deviceUrl.getImportKey())) {
                     present = true;
                     uploaded &= HttpConnector.uploadToServer(
                           new URL(basePath + fileItem.getPath()),
                           new URL(deviceUrl.getUrl())
                           );

                     if (uploaded) {
                        _logger.info(
                              String.format("Uploaded file: %s", fileItem.getPath())
                              );
                     } else {
                        _logger.error(
                              String.format("File not uploaded: %s", fileItem.getPath())
                              );
                     }

                     break;
                  }
               }

               if (!present || !uploaded) {
                  _logger.error("File ID does not match any import key"
                        + "or upload was unsuccessful");
                  uploaded = false;
                  break;
               }
            }
         }
      } else {
         _logger.error("ImportSpec, HttpNfcLeaseInfo or OVF path are incorrect");
         uploaded = false;
      }

      return uploaded;
   }

   /**
    * Gets a key based on the index of the device in the list of all devices
    * @param deviceIndex - the index of the device in the list of all devices
    * @return a key based on the index of the device in the list of all devices
    */
   private int getDeviceKey(int deviceIndex) {
      return deviceIndex - DEVICE_KEY_MAGIC_NUMBER;
   }

   /**
    * Evaluates parent of the passed VmSpec in terms of validity. Returns parent
    * in general case and grandparent in case of a DRS cluster parent. It is
    * needed to conform with the current API limitations when querying for VM
    * host or datastore
    *
    * @param vm
    * @returns Original VmSpec's parent or the grandparent in case of a DRS
    *          cluster
    */
   private ManagedEntitySpec evaluateVmParent(VmSpec vm) {
      ManagedEntitySpec vmParent = vm.parent.get();
      if (!(vmParent instanceof HostSpec)) {
         // If user passed a DRS cluster as a parent
         if (vmParent instanceof ClusterSpec
               && ((ClusterSpec) vmParent).drsEnabled.get()) {
            // Take cluster's parent which
            // should be a datacenter
            vmParent = vmParent.parent.get();
         } else {
            throw new IllegalArgumentException(
                  "The VM's parent should be a host or a DRS cluster.");
         }
      }
      return vmParent;
   }

   /**
    * Removes a virtual disk from a vm without removing the disk and
    * disk-related files from the datastore.
    *
    * @param vmSpec
    *           vm whose disk is to be removed
    * @param diskSpec
    *           the disk to remove
    * @return true if the disk has been removed, false if an error has occurred
    *         while removing it
    * @throws Exception
    */
   public boolean removeDiskFromVm(VmSpec vmSpec, VirtualDiskSpec diskSpec)
         throws Exception {
      return deleteDisk(vmSpec, diskSpec, false);
   }

   /**
    * Removes a virtual disk from a vm and deletes its files from the datastore
    * where it is stored.
    *
    * @param vmSpec
    *           vm whose disk is to be removed
    * @param diskSpec
    *           the disk to delete
    * @return true if the disk has been deleted, false if an error has occurred
    *         while deleting it
    * @throws Exception
    */
   public boolean deleteDiskFromVm(VmSpec vmSpec, VirtualDiskSpec diskSpec)
         throws Exception {
      return deleteDisk(vmSpec, diskSpec, true);
   }

   private boolean deleteDisk(VmSpec vmSpec, VirtualDiskSpec diskSpec,
         boolean deleteDiskFiles) throws Exception {
      _logger
            .info(String.format(
                  "Deleting disk %s from vm %s %s",
                  diskSpec.name.get(),
                  vmSpec.name.get(),
                  deleteDiskFiles ? "while also deleting all disk-related files from datastore"
                        : "without deleting disk-related files from datastore"));
      final VirtualMachine vm = getVirtualMachine(vmSpec);
      final VirtualDisk existingDisk = getDisk(vm, diskSpec);
      if (existingDisk == null) {
         throw new IllegalStateException(String.format("%s not found.",
               diskSpec.name.get()));
      }
      VirtualDisk disk = new VirtualDisk();
      disk.setKey(existingDisk.getKey());
      VirtualDeviceSpec virtualDeviceSpec = new VirtualDeviceSpec();
      virtualDeviceSpec.setDevice(disk);
      virtualDeviceSpec.setOperation(Operation.remove);
      if (deleteDiskFiles) {
         virtualDeviceSpec.setFileOperation(FileOperation.destroy);
      }
      ConfigSpec configSpec = new ConfigSpec();
      configSpec.setDeviceChange(new VirtualDeviceSpec[] { virtualDeviceSpec });
      ManagedObjectReference taskRef = vm.reconfigure(configSpec);
      return VcServiceUtil.waitForTaskSuccess(taskRef, vmSpec.service.get());
   }

   /**
    * Checks if the given disk is listed in the virtual hardware devices of the
    * given vm.
    *
    * @param vmSpec
    *           the vm to check
    * @param diskSpec
    *           the disk to look for
    * @return true if the disk is attached to the given vm, false otherwise
    * @throws Exception
    */
   public boolean isDiskExisting(VmSpec vmSpec, VirtualDiskSpec diskSpec)
         throws Exception {
      final VirtualMachine vm = getVirtualMachine(vmSpec);
      return getDisk(vm, diskSpec) != null;
   }

   private VirtualDisk getDisk(VirtualMachine vm,
         VirtualDiskSpec diskSpec) {
      List<VirtualDisk> virtualDisks = getVirtualDevice(vm, VirtualDisk.class);
      for (VirtualDisk disk : virtualDisks) {
         if (disk.getBacking() instanceof FlatVer2BackingInfo) {
            FlatVer2BackingInfo backingInfo = (FlatVer2BackingInfo) disk
                  .getBacking();
            if (backingInfo.getFileName().equals(diskSpec.getAbsolutePath())) {
               return disk;
            }
         }
      }
      return null;
   }

   @SuppressWarnings("unchecked")
   private <E extends VirtualDevice> List<E> getVirtualDevice(
         VirtualMachine vm, Class<E> virtualDeviceType) {
      List<E> virtualDevices = new ArrayList<>();
      ConfigInfo configInfo = vm.getConfig();
      VirtualHardware virtuaHardware = configInfo.getHardware();
      VirtualDevice[] virtualDevicesArray = virtuaHardware.getDevice();
      for (VirtualDevice device : virtualDevicesArray) {
         if (device.getClass().equals(virtualDeviceType)) {
            virtualDevices.add((E) device);
         }
      }
      return virtualDevices;
   }
}
