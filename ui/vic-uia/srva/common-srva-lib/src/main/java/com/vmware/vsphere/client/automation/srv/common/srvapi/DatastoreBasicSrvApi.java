/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import org.apache.commons.lang.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.sms.StorageManager;
import com.vmware.vim.binding.sms.storage.StorageContainer;
import com.vmware.vim.binding.sms.storage.StorageContainerResult;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.host.DatastoreSystem;
import com.vmware.vim.binding.vim.host.DatastoreSystem.VvolDatastoreSpec;
import com.vmware.vim.binding.vim.host.NasVolume.Specification;
import com.vmware.vim.binding.vim.host.ScsiDisk;
import com.vmware.vim.binding.vim.host.VmfsDatastoreCreateSpec;
import com.vmware.vim.binding.vim.host.VmfsDatastoreOption;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * Class that allows CRUD operations of a Datastore
 * in the VC Inventory based on a Datastore specification.
 */
public class DatastoreBasicSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(DatastoreBasicSrvApi.class);
   private static final Integer VMFS_VERSION_5 = 5;

   private static DatastoreBasicSrvApi instance = null;
   protected DatastoreBasicSrvApi() {}

   /**
    * Get instance of DatastoreSrvApi.
    *
    * @return  created instance
    */
   public static DatastoreBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized(DatastoreBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing DatastoreSrvApi.");
               instance = new DatastoreBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a datastore under a host.
    * @param datastoreSpec the spec of the datastore, parent should be a host
    * @return true if the creation was successful, false otherwise.
    * @throws Exception
    */
   public boolean createDatastore(DatastoreSpec datastoreSpec) throws Exception {
      validateDatastoreSpec(datastoreSpec);

      _logger.info(String.format("Creating datastore '%s'", datastoreSpec.name.get()));

      // get datastore system of the host
      DatastoreSystem datastoreSystem = getDatastoreSystem(datastoreSpec);

      if (datastoreSpec.type.isAssigned()
            && datastoreSpec.type.get() == DatastoreType.NFS) {
         Specification nfsSpec = getNfsSpec(datastoreSpec);
         return datastoreSystem.createNasDatastore(nfsSpec) != null;
      } else if (datastoreSpec.type.isAssigned() &&
            datastoreSpec.type.get() == DatastoreType.VMFS) {
         VmfsDatastoreCreateSpec vmfsSpec = getVmfsSpec(datastoreSystem, datastoreSpec);
         if (vmfsSpec == null) {
            throw new Exception("No available disks for vmfs creation!");
         }
         return datastoreSystem.createVmfsDatastore(vmfsSpec) != null;
      } else {
         throw new NotImplementedException(
               "Creation of datastores different from NFS is not supported for now!");
      }
   }

   /**
    * This method is deprecated. See
    * {@link #unmountDatastore(DatastoreSpec, HostSpec)}
    *
    * The unmount of a datastore itself is not defined. The unmount is only
    * meaningful for a host - datastore pair. Using the datastore specs parent
    * is dubious, error prone and lacks of ability to perform unmount for any
    * other host
    *
    * @param datastoreSpec
    *           the spec of the datastore, parent should be a host
    * @throws Exception
    */
   @Deprecated
   public void unmountDatastore(DatastoreSpec datastoreSpec) throws Exception {
      this.unmountDatastore(datastoreSpec,
            (HostSpec) datastoreSpec.parent.get());
   }

   /**
    * Unmounts datastore from a host
    *
    * @param datastoreSpec
    *           the spec for the datastore to be unmounted
    * @param hostSpec
    *           the spec for the hosts from which the datastore should be
    *           unmounted
    * @throws Exception
    */
   public void unmountDatastore(DatastoreSpec datastoreSpec, HostSpec hostSpec)
         throws Exception {
      _logger.info(String.format("Unmounting datastore '%s' from host '%s'",
            datastoreSpec.name.get(), hostSpec.name.get()));

      DatastoreSystem hostDatastoreSystem = HostBasicSrvApi.getInstance()
            .getDatastoreSystem(hostSpec);
      Datastore datastore = ManagedEntityUtil.getManagedObject(datastoreSpec);

      hostDatastoreSystem.removeDatastore(datastore._getRef());
   }

   /**
    * Delete datastore.
    *
    * @param datastoreSpec    spec of the datastore that will be destroyed
    * @return                 true if the deletion was successful, false otherwise
    * @throws Exception       if login to the VC service fails
    */
   public boolean deleteDatastore(DatastoreSpec datastoreSpec) throws Exception {
      validateDatastoreSpec(datastoreSpec);

      _logger.info(String.format("Deleting datastore '%s'", datastoreSpec.name.get()));

      Datastore datastore = ManagedEntityUtil.getManagedObject(datastoreSpec);
      ManagedObjectReference taskMoRef = datastore.destroy();

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, datastoreSpec.service.get())
            && ManagedEntityUtil.waitForEntityDeletion(
                  datastoreSpec,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified datastore exists and deletes it.
    *
    * @param datastoreSpec    spec of the datastore that will be deleted
    * @return                 true if the datastore doesn't exist or if the datastore is
    *                         deleted successfully, false otherwise.
    * @throws Exception       if login to the VC service fails, or the specified datastore
    *                         doesn't exist.
    */
   public boolean deleteDatastoreSafely(DatastoreSpec datastoreSpec)
         throws Exception {

      if (checkDatastoreExists(datastoreSpec)) {
         return deleteDatastore(datastoreSpec);
      }

      // Positive result if the folder doesn't exist
      return true;
   }

   /**
    * Checks whether the specified datastore exists.
    *
    * @param datastoreSpec    spec of the datastore that will be checked for existence
    * @return                 true is the datastore exists, false otherwise.
    * @throws Exception       if login to the VC service fails
    */
   public boolean checkDatastoreExists(DatastoreSpec datastoreSpec)
         throws Exception {
      validateDatastoreSpec(datastoreSpec);

      _logger.info(String.format(
            "Checking whether datastore '%s' exists",
            datastoreSpec.name.get()));

      try {
         ManagedEntityUtil.getManagedObject(datastoreSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }

      return true;
   }

   /**
    * Rename datastore.
    *
    * @param datastoreSpec spec of the datastore that will be destroyed
    * @return true if the deletion was successful, false otherwise
    * @throws Exception if login to the VC service fails
    */
   public boolean renameDatastore(DatastoreSpec datastoreToRenameSpec, String newName) throws Exception {
      validateDatastoreSpec(datastoreToRenameSpec);

      _logger.info(String.format("Renaming datastore '%s'", datastoreToRenameSpec.name.get()));

      Datastore datastore = ManagedEntityUtil.getManagedObject(datastoreToRenameSpec);
      ManagedObjectReference taskMoRef = datastore.rename(newName);

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, datastoreToRenameSpec);
   }

   //---------------------------------------------------------------------------
   // Private methods

   /**
    * This method is deprecated see
    * {@link HostBasicSrvApi#getDatastoreSystem(HostSpec)}
    *
    * The {@link DatastoreSystem} is property of a host and not of a datastore
    *
    * @param datastoreSpec
    * @return
    * @throws Exception
    */
   @Deprecated
   private DatastoreSystem getDatastoreSystem(DatastoreSpec datastoreSpec)
         throws Exception {
      VcService service = VcServiceUtil.getVcService(datastoreSpec);
      HostSystem host = ManagedEntityUtil.getManagedObject(datastoreSpec.parent
            .get());
      return service.getManagedObject(host.getConfigManager()
            .getDatastoreSystem());
   }

   // Method is default as it is used also in HostClientSrvApi
   void validateDatastoreSpec(DatastoreSpec spec)
         throws IllegalArgumentException {
      if (!spec.name.isAssigned()) {
         throw new IllegalArgumentException("Datastore name is not set.");
      }

      if (!spec.parent.isAssigned() || !(spec.parent.get() instanceof HostSpec)) {
         throw new IllegalArgumentException("Datastore parent has to be a host.");
      }

      // validate some more options for NFS datastores
      if (spec.type.isAssigned() && spec.type.get() == DatastoreType.NFS) {
         if (!spec.remoteHost.isAssigned()) {
            throw new IllegalArgumentException("Datastore remoteHost is not set.");
         }

         if (!spec.remotePath.isAssigned()) {
            throw new IllegalArgumentException("Datastore remotePath is not set.");
         }
      }
   }

   private Specification getNfsSpec(DatastoreSpec datastoreSpec) {
      Specification spec =
            new Specification(datastoreSpec.remoteHost.get(),
                  datastoreSpec.remotePath.get(), datastoreSpec.name.get(), "readWrite",
                  "nfs", null, null, null, null);
      return spec;
   }

   /**
    * Method that creates a VMFS Spec needed for creation of a VMFS datatsore on a host; the method
    * has default access as it is used also in HostClientBasicSrvApi.
    * @param dsSystem - the datastore system of the host on which to create the VMFS datastore
    * @param dsSpec - the Datastore spec - only name, type and parent host are needed
    * @return VmfsDatastoreCreateSpec, empty if no available disk
    * @throws Exception
    */
   // Method is also used in HostClientSrvApi
   VmfsDatastoreCreateSpec getVmfsSpec(DatastoreSystem dsSystem, DatastoreSpec dsSpec) throws Exception {
      VmfsDatastoreCreateSpec vmfsDsCreateSpec = null;
      // query for available disks
      ScsiDisk[] scsiDisks = dsSystem.queryAvailableDisksForVmfs(null);

      if (scsiDisks != null) {

         for (int i = 0; i < scsiDisks.length; i++) {
            // VMFS major version is 3 or 5, by default we are creating with latest
            // query disk for available creation options
            VmfsDatastoreOption[] vmfsOptions = dsSystem.queryVmfsDatastoreCreateOptions(
                  scsiDisks[i].getDevicePath(),
                  VMFS_VERSION_5);
            if (vmfsOptions.length > 0) {
               // if there are creation options, get creation spec
               vmfsDsCreateSpec = (VmfsDatastoreCreateSpec) vmfsOptions[0].getSpec();
               // set name from Datastore spec
               vmfsDsCreateSpec.getVmfs().setVolumeName(dsSpec.name.get());
               break;
            }
         }
      }

      return vmfsDsCreateSpec;
   }



   /**
    * Create VVOL datastore spec.
    *
    * @param datastoreSpec    holds all params of the datastore
    * @return                 created VVOL datastore spec
    * @throws Exception
    */
   private VvolDatastoreSpec getVvolSpec(DatastoreSpec datastoreSpec)
         throws Exception {
      validateVvolDatastoreSpec(datastoreSpec);

      StorageManager storageManager = StorageProviderBasicSrvApi.getInstance().getStorageManager(datastoreSpec);

      BlockingFuture<StorageContainerResult> storageContainerResultFuture =
            new BlockingFuture<StorageContainerResult>();
      storageManager.queryStorageContainer(null, storageContainerResultFuture);
      StorageContainerResult storageContainerResult = storageContainerResultFuture.get();

      StorageContainer[] storageContainer = storageContainerResult.getStorageContainer();
      if (storageContainer.length == 0) {
         _logger.error("There are no available storage containers.");
         throw new IllegalArgumentException("There are no available storage containers.");
      }

      //TODO: add filtering of storage containers based on the host of the datastore
      // for now only the first available container will be used
      VvolDatastoreSpec vvolSpec = new VvolDatastoreSpec(
            datastoreSpec.name.get(),
            storageContainer[0].getUuid()
            );

      return vvolSpec;
   }

   /**
    * Validates if the datastore spec is suitable for VVOL datastore creation.
    *
    * For now only VVOL datastores that are created on host available containers
    * are supported. This means that the datastore parent has to be a host.
    *
    * @param spec    datastore spec that will be validated
    */
   private void validateVvolDatastoreSpec(DatastoreSpec spec) {
      if (!spec.parent.isAssigned() || !(spec.parent.get() instanceof HostSpec)) {
         throw new IllegalArgumentException("VVOL Datastore parent has to be a host.");
      }
   }
}
