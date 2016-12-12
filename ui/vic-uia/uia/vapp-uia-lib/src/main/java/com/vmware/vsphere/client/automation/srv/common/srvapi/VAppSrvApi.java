/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.vim.binding.vim.ClusterComputeResource;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.ManagedEntity;
import com.vmware.vim.binding.vim.Network;
import com.vmware.vim.binding.vim.ResourceAllocationInfo;
import com.vmware.vim.binding.vim.ResourceConfigSpec;
import com.vmware.vim.binding.vim.ResourcePool;
import com.vmware.vim.binding.vim.SharesInfo;
import com.vmware.vim.binding.vim.SharesInfo.Level;
import com.vmware.vim.binding.vim.VirtualApp;
import com.vmware.vim.binding.vim.VirtualApp.Summary;
import com.vmware.vim.binding.vim.VirtualApp.VAppState;
import com.vmware.vim.binding.vim.vApp.VAppConfigSpec;
import com.vmware.vim.binding.vim.vm.ConfigSpec;
import com.vmware.vim.binding.vim.vm.FileInfo;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

public class VAppSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(VAppSrvApi.class);

   protected static final String DEFAULT_VM_GUEST_OS = "winVistaGuest";

   private static VAppSrvApi instance = null;
   protected VAppSrvApi() {}

   /**
    * Get instance of VAppSrvApi.
    *
    * @return  created instance
    */
   public static VAppSrvApi getInstance() {
      if (instance == null) {
         synchronized(VAppSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing VAppSrvApi.");
               instance = new VAppSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a Vapp with properties as specified in the <code>VappSpec
    * </code> parameter. Currently creation is possible only on host. If the
    * vapp spec contains vm specs the vms will also be created in the vapp.
    * Otherwise the vapp will remain empty. The VMs would be added on the
    * datastore of the first host encountered.
    *
    * @param vappSpec
    *           <code>VappSpec</code> containing the properties of the Vapp to
    *           be created.
    *
    * @return True if the creation was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean createVApp(VappSpec vappSpec) throws Exception {
      validateVAppSpec(vappSpec);

      // TODO refactor the method to be shorter

      _logger.info(String.format("Creating vApp '%s'", vappSpec.name.get()));

      if (vappSpec.parent.get() instanceof HostSpec) {
         HostSpec parentHostSpec = (HostSpec) vappSpec.parent.get();
         HostSystem vAppHost = ManagedEntityUtil.getManagedObject(parentHostSpec);
         Folder vmFolder = FolderBasicSrvApi.getInstance().getVmFolder(getDatacenterSpec(parentHostSpec));

         VcService vcService = VcServiceUtil.getVcService(vappSpec);

         ComputeResource cr =
               (ComputeResource) vcService.getManagedObject(vAppHost.getParent());

         ResourcePool parentRp =
               (ResourcePool) vcService.getManagedObject(cr.getResourcePool());

         boolean vappCreated =
               parentRp.createVApp(
               vappSpec.name.get(),
               getDefaultResourceConfiguration(),
               new VAppConfigSpec(),
               vmFolder._getRef()) != null;

         if (vappSpec.vmList.getAll().isEmpty() || !vappCreated) {
            return vappCreated;
         }

         /* Create the VMs inside the vApp. */

         VirtualApp virtualApp = getVirtualApp(vappSpec);

         // Getting the datastores available.
         ManagedObjectReference[] hostDatastoresMor = vAppHost.getDatastore();
         Map<String, Datastore> datastoresByNames =
               new HashMap<String, Datastore>(hostDatastoresMor.length);
         for (ManagedObjectReference datastoreMor : hostDatastoresMor) {
            Datastore datastore =
                  vcService.getManagedObject(datastoreMor);

            datastoresByNames.put(datastore.getName(), datastore);
         }

         FileInfo fileInfo = null;
         Datastore vmDatastore = null;
         for (VmSpec vmSpec : vappSpec.vmList.getAll()) {
            ConfigSpec vmConfigSpec = new ConfigSpec();
            vmConfigSpec.setName(vmSpec.name.get());

            // Check for preferred datastores.
            if (vmSpec.datastore.isAssigned()) {
               vmDatastore = datastoresByNames.get(vmSpec.datastore.get().name.get());
               if (vmDatastore == null) {
                  throw new IllegalArgumentException(
                        String.format(
                              "Failed to find datastore with name %s for VM with name %s",
                              vmSpec.datastore.get().name.get(), vmSpec.name.get()));
               }
            } else { // If not just re-use the prevous one or use the first in the list.
               if (vmDatastore == null) {
                  vmDatastore = (Datastore) datastoresByNames.values().toArray()[0];
               }
            }

            fileInfo = new FileInfo(
                     String.format("[%s]", vmDatastore.getName()),
                     null, null, null, null
            );

            vmConfigSpec.setFiles(fileInfo);

            // Checking for preferred guest.
            if (vmSpec.guestId.isAssigned()) {
               vmConfigSpec.setGuestId(vmSpec.guestId.get());
            } else {
               // set default guest ID in order to be able to power on the VM
               vmConfigSpec.setGuestId(DEFAULT_VM_GUEST_OS);
            }

            //Checking for NIC adapters in VmSpec.
            if (vmSpec.nicList.isAssigned()) {
               vmConfigSpec.deviceChange =
                     VmSrvApi.getInstance().toVirtualDeviceSpec(vmSpec.nicList.getAll());
            }

            ManagedObjectReference taskCreateVm =
                  virtualApp.createVm(vmConfigSpec, vAppHost._getRef());

            if (!VcServiceUtil.waitForTaskSuccess(taskCreateVm, vappSpec)) {
               _logger.error("Failed to create VM with name: '" + vmSpec.name.get()
                     + "'");
               return false;
            }
         }

         return true;
      }

      if (vappSpec.parent.get() instanceof ClusterSpec) {

         ClusterSpec parentClusterSpec = (ClusterSpec) vappSpec.parent.get();
         // Check if DRS is enalbed
         if (!parentClusterSpec.drsEnabled.get()) {
            _logger.error("DRS on cluster has to be enalbed");
            return false;
         }
         ClusterComputeResource vAppCluster = ManagedEntityUtil
               .getManagedObject(parentClusterSpec);
         Folder vmFolder = FolderBasicSrvApi.getInstance()
               .getVmFolder((DatacenterSpec) parentClusterSpec.parent.get());

         VcService vcService = VcServiceUtil.getVcService(vappSpec);

         ResourcePool parentRp = (ResourcePool) vcService
               .getManagedObject(vAppCluster.getResourcePool());

         boolean vappCreated = parentRp.createVApp(vappSpec.name.get(),
               getDefaultResourceConfiguration(), new VAppConfigSpec(),
               vmFolder._getRef()) != null;

         return vappCreated;

         // TODO add creation of vms under the vapp and verification if the
         // cluster is under folder

      }
      return false;
   }

   /**
    * Deletes the specified vApp from the inventory and waits for timeout
    * seconds to get the object really deleted. If timeout is over and object
    * still in inventory, it will return false
    *
    * @param vappSpec
    *           <code>VappSpec</code> instance representing the vApp to be
    *           deleted.
    *
    * @return True if the deletion was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteVApp(VappSpec vappSpec) throws Exception {

      validateVAppSpec(vappSpec);

      _logger.info(String.format(
            "Deleting vApp '%s' from parent '%s'",
            vappSpec.name.get(),
            vappSpec.parent.get().name.get()));

      VirtualApp vApp = ManagedEntityUtil.getManagedObject(vappSpec);

      // TODO check vApp status
      ManagedObjectReference taskMoRef = vApp.destroy();

      return VcServiceUtil.waitForTaskSuccess(taskMoRef, vappSpec)
            && ManagedEntityUtil.waitForEntityDeletion(
                  vappSpec,
                  (int) SUITA.Environment.getBackendJobSmall() / 1000);
   }

   /**
    * Checks whether the specified vApp exists and deletes it.
    *
    * @param vappSpec
    *           <code>VappSpec</code> instance representing the vApp to be
    *           deleted.
    *
    * @return True if the vApp doesn't exist or if the vApp is deleted
    *         successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteVAppSafely(VappSpec vappSpec) throws Exception {

      if (checkVAppExists(vappSpec)) {
         return deleteVApp(vappSpec);
      }

      // Positive result if the vApp doesn't exist
      return true;
   }

   /**
    * Renames the vApp
    * @param originalVAppSpec - vApp to update
    * @param newVAppSpec - spec with new values
    * @return true if successful, false if not or if no such pool exists
    * @throws Exception
    */
   public boolean renameVapp(VappSpec originalVAppSpec, VappSpec newVAppSpec)
       throws Exception {

    if (checkVAppExists(originalVAppSpec)) {

        validateVAppSpec(newVAppSpec);

        ServiceSpec serviceSpec = originalVAppSpec.service.get();

        _logger.info(String.format(
              "Updating vApp '%s' from parent '%s'",
              originalVAppSpec.name.get(), originalVAppSpec.parent.get().name.get()));

        VirtualApp vapp = ManagedEntityUtil.getManagedObject(originalVAppSpec,
              serviceSpec);
        vapp.rename(newVAppSpec.name.get());

        return checkVAppExists(newVAppSpec) && !checkVAppExists(originalVAppSpec);
    }

    // vApp doesn't exist
    return false;
 }

   /**
    * Checks whether the specified Vapp exists in the specified resource pool.
    *
    * @param vappSpec
    *           <code>VappSpec</code> instance representing the vApp to be
    *           queried.
    *
    * @return True is the vApp exists, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean checkVAppExists(VappSpec vappSpec) throws Exception {

      validateVAppSpec(vappSpec);

      _logger.info(String.format(
            "Checking whether Vapp '%s' exists in parent '%s'",
            vappSpec.name.get(),
            vappSpec.parent.get().name.get()));

      try {
         ManagedEntityUtil.getManagedObject(vappSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }

      return true;
   }

   /**
    * Powers on vApp.
    *
    * @param vappSpec
    *           vApp specification that will be powered on
    * @return              true if the operation was successful, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified vApp
    *                      settings are invalid
    */
   public boolean powerOnVapp(VappSpec vappSpec) throws Exception {
      validateVappSpec(vappSpec);

      _logger.info(String.format("Powering on vApp '%s'", vappSpec.name.get()));

      VirtualApp virtualApp = getVirtualApp(vappSpec);
      ManagedObjectReference taskMoRef = virtualApp.powerOn();

      boolean taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vappSpec);
      // TODO: fix it lgrigorova
      CommonUtils.sleep(10000);
      return taskResult;
   }

   /**
    * Suspends a vApp.
    *
    * @param vappSpec
    *           vApp specification that will be suspended
    * @return              true if the operation was successful, false otherwise
    * @throws Exception
    *            if login to the VC service fails, or the specified vApp
    *                      settings are invalid
    */
   public boolean suspendVapp(VappSpec vappSpec) throws Exception {
      validateVappSpec(vappSpec);

      _logger.info(String.format("Suspending vApp '%s'", vappSpec.name.get()));

      VirtualApp virtualApp = getVirtualApp(vappSpec);
      ManagedObjectReference taskMoRef = virtualApp.suspend();

      boolean taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vappSpec);
      // TODO: fix it lgrigorova
      CommonUtils.sleep(10000);
      return taskResult;
   }

   /**
    * Powers off vApp.
    *
    * @param vappSpec
    *           vApp specification that will be powered off
    * @param safely
    *           if <code>true</code> the method will do some pre-checks before
    *                      trying to power off the vapp. For example whether the vapp
    *                      exists and if it's powered on. If <code>false</code> the method
    *                      directly tryies to power off the vapp which might raise some
    *           exceptions in case the vapp doesn't exist or is already powered
    *           off.
    * @return              true if the operation was successful, false otherwise
    * @throws Exceptionif
    *            login to the VC service fails, or the specified vApp settings
    *            are invalid
    */
   public boolean powerOffVapp(VappSpec vappSpec, boolean safely)
         throws Exception {
      validateVappSpec(vappSpec);

      _logger.info(String.format("Powering off vApp '%s'", vappSpec.name.get()));

      VirtualApp virtualApp = null;
      try {
         virtualApp = getVirtualApp(vappSpec);
      } catch (ObjectNotFoundException onfe) {
         if (safely) {
            return false;
         }
         throw onfe;
      }

      boolean started = true;

      if (safely) {
         VirtualApp.Summary vappSummary = (Summary) virtualApp.getSummary();
         VAppState powerState = vappSummary.getVAppState();
         started = (powerState == VAppState.started || powerState == VAppState.starting);
      }

      boolean taskResult = false;

      if (started) {
         ManagedObjectReference taskMoRef = virtualApp.powerOff(true); // Force
                                                                       // power
                                                                       // off.
         taskResult = VcServiceUtil.waitForTaskSuccess(taskMoRef, vappSpec);
         // TODO: fix it lgrigorova
         CommonUtils.sleep(10000);
      }

      return taskResult;
   }

   /**
    * Retrieves virtual machine from specified spec object.
    *
    * @param vappSpec
    *           spec of the vApp that will be retrieved
    * @return              retrieved vApp object
    * @throws Exception
    *            if login to the VC service fails, or the specified vApp
    *                      settings are invalid
    */
   public VirtualApp getVirtualApp(VappSpec vappSpec) throws Exception {
      return ManagedEntityUtil.getManagedObject(vappSpec);
   }


   /**
    * Retrieves the datastore of the specified vApp.
    *
    * @param vappSpec
    *           - the vApp to query for
    * @return the datastore of the vApp
    * @throws Exception
    *            - Exception if there is an error when finding the vApp
    */
   public DatastoreSpec getVappDatastore(VappSpec vappSpec)
         throws Exception {
      validateVappSpec(vappSpec);

      ServiceSpec serviceSpec = vappSpec.service.get();

      VirtualApp vapp = getVirtualApp(vappSpec);
      ManagedObjectReference datastoreMoRef = vapp.getDatastore()[0];

      ManagedEntity datastoreME = ManagedEntityUtil.getManagedObjectFromMoRef(
            datastoreMoRef, serviceSpec);

      DatastoreSpec datastore = SpecFactory.getSpec(DatastoreSpec.class,
            datastoreME.getName(), vappSpec.parent.get());

      return datastore;
   }

   /**
    * Retrieves list of vApp network names.
    *
    * @param vappSpec
    *           the spec for the vApp that will be searched
    * @throws RuntimeException
    *            if login to the VC server fails or host is not present
    */
   public List<String> getNetworkNames(VappSpec vappSpec) {
      List<String> networkNames = new ArrayList<String>();
      VirtualApp vapp = null;
      VcService service = null;
      ServiceSpec serviceSpec = vappSpec.service.get();

      try {
         service = VcServiceUtil.getVcService(serviceSpec);
         vapp = ManagedEntityUtil.getManagedObject(vappSpec, serviceSpec);
      } catch (Exception e) {
         String errorMessage =
               String.format("Cannot retrieve vApp %s", vappSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }

      for (ManagedObjectReference networkMor : vapp.getNetwork()) {
         Network network = service.getManagedObject(networkMor);
         networkNames.add(network.getName());
      }

      return networkNames;
   }

   ///////////////////////////////////////////////////////////////////////////////////////
   //                                PRIVATE METHODS                                    //
   ///////////////////////////////////////////////////////////////////////////////////////

   /**
    * Method that validates that both name and parent are assigned
    *
    * @param vappSpec
    *           - the VappSpec to be verified
    */
   private void validateVAppSpec(VappSpec vappSpec) {
      if (!vappSpec.name.isAssigned()) {
         throw new IllegalArgumentException("Vapp name is not set");
      }

      if (!vappSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("Vapp parent is not set");
      }
   }

   /**
    * Method that creates a default resource allocation spec with same values
    * for memory and cpu: sets expandable reservation, unlimited cpu and memory,
    * no overhead limit, 0 reservation and normal level of shares
    *
    * @return ResourceConfigSpec - a default config spec
    */
   private ResourceConfigSpec getDefaultResourceConfiguration() {

      ResourceConfigSpec resCfgSpec = new ResourceConfigSpec();

      // set defaults for CPU allocation
      ResourceAllocationInfo cpuAlloc = new ResourceAllocationInfo();
      cpuAlloc.setExpandableReservation(Boolean.TRUE);
      cpuAlloc.setLimit(new Long(-1));
      cpuAlloc.setOverheadLimit(null);
      cpuAlloc.setReservation(new Long(0));
      SharesInfo cpuSharesInfo = new SharesInfo();
      cpuSharesInfo.setLevel(Level.normal);
      cpuSharesInfo.setShares(0);
      cpuAlloc.setShares(cpuSharesInfo);

      // set defaults for Mem allocation
      ResourceAllocationInfo memAlloc = new ResourceAllocationInfo();
      memAlloc.setExpandableReservation(true);
      memAlloc.setLimit(new Long(-1));
      memAlloc.setOverheadLimit(null);
      memAlloc.setReservation(new Long(0));
      SharesInfo memSharesInfo = new SharesInfo();
      memSharesInfo.setLevel(Level.normal);
      memSharesInfo.setShares(0);
      memAlloc.setShares(memSharesInfo);

      resCfgSpec.setCpuAllocation(cpuAlloc);
      resCfgSpec.setMemoryAllocation(memAlloc);

      return resCfgSpec;
   }

   /**
    * Method to get the Datacenter Spec of the parent datacenter of a host; The
    * method won't return correctly if the cluster of the host is in a folder
    *
    * @param hostSpec
    *           - the spec of the host whose parent datacenter is looked for
    * @return DatacenterSpec - the spec of the parent datacenter
    */
   private DatacenterSpec getDatacenterSpec(HostSpec hostSpec) {
      // For standalone host return the spec of its direct parent who is
      // datacenter
      if (hostSpec.parent.get() instanceof DatacenterSpec) {
         return (DatacenterSpec) hostSpec.parent.get();
      }

      // TODO not working the cluster is in folder
      // For clustered host return the spec of the parent of the host's parent-
      // i.e. host's parent is a cluster whose parent is a datacenter
      ClusterSpec parentClusterSpec = (ClusterSpec) hostSpec.parent.get();
      return (DatacenterSpec) parentClusterSpec.parent.get();
   }

   /**
    * Method that validates the vApp spec: - name should be assigned - parent
    * host should be assigned
    *
    * @param vAppSpec
    *           vApp spec that will be validated
    * @throws IllegalArgumentException
    *            if vApp spec requirements are not met
    */
   private void validateVappSpec(VappSpec vAppSpec)
         throws IllegalArgumentException {
      if (!vAppSpec.name.isAssigned() || vAppSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("vApp name is not set.");
      }

      if (!vAppSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("vApp host is not set.");
      }

      EntitySpec vAppParentEntity = vAppSpec.parent.get();
      if (!(vAppParentEntity instanceof HostSpec)
            && !(vAppParentEntity instanceof ClusterSpec)
            && !(vAppParentEntity instanceof DatacenterSpec)
            && !(vAppParentEntity instanceof FolderSpec)
            && !(vAppParentEntity instanceof VappSpec)
            && !(vAppParentEntity instanceof ResourcePoolSpec)) {
         throw new IllegalArgumentException(
               "vApp parent is not a host, cluster, datacenter, folder, vApp or resource pool.");
      }
   }
}
