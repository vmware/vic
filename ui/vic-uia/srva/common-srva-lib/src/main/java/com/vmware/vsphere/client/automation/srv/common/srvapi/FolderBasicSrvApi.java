/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.Datacenter;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.vmomi.core.Future;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderType;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

public class FolderBasicSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(FolderBasicSrvApi.class);

   private static FolderBasicSrvApi instance = null;
   protected FolderBasicSrvApi() {}

   /**
    * Get instance of FolderSrvApi.
    *
    * @return  created instance
    */
   public static FolderBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized(FolderBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing FolderSrvApi.");
               instance = new FolderBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a folder with FolderSpec as its parameter
    *
    * @param folderSpec
    *           - specification of the folder to be created
    *
    * @return True if the creation was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean createFolder(FolderSpec folderSpec) throws Exception {
      validateFolderSpec(folderSpec);

      DatacenterSpec dcSpec = getParentDcSpec(folderSpec);
      // Only if it is a Datacenter folder, it will no parent, otherwise its
      // parent should be a datacenter
      if (!folderSpec.type.get().equals(FolderType.DATACENTER)
            && !folderSpec.parent.isAssigned() && dcSpec == null) {
         throw new IllegalArgumentException(
               "Folder spec is not correct, change type or assign parent.");
      }

      _logger.info(String.format("Creating folder '%s'", folderSpec.name.get()));

      // Get the root folder
      Folder parentFolder = null;

      if (folderSpec.parent.isAssigned()
            && folderSpec.parent.get() instanceof FolderSpec) {
         parentFolder = ManagedEntityUtil.getManagedObject(folderSpec.parent.get());
      } else {
         switch (folderSpec.type.get()) {
            case DATACENTER:
               parentFolder = getRootFolder(folderSpec.service.get());
               break;
            case HOST:
               parentFolder = getHostFolder(dcSpec);
               break;
            case NETWORK:
               parentFolder = getNetworkFolder(dcSpec);
               break;
            case STORAGE:
               parentFolder = getStorageFolder(dcSpec);
               break;
            case VM:
               parentFolder = getVmFolder(dcSpec);
               break;
            default:
               throw new IllegalArgumentException("No such Folder type.");
         }
      }
      return parentFolder.createFolder(folderSpec.name.get()) != null;
   }

   /**
    * Deletes the specified folder from the inventory and waits the specified
    * timeout to let the entity get deleted, returns false if the timeout is
    * over or the entity cannot get deleted
    *
    * @param folder
    *           - specification of the Folder to be deleted
    * @return True if the deletion was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteFolder(FolderSpec folder) throws Exception {
      validateFolderSpec(folder);

      _logger.info(String.format("Deleting folder '%s", folder.name.get()));

      // Get the folder
      Folder folderObj = ManagedEntityUtil.getManagedObject(folder);

      ManagedObjectReference taskMoRef = null;
      taskMoRef = folderObj.destroy();

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, folder.service.get())
            && ManagedEntityUtil.waitForEntityDeletion(
                  folder,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified folder exists.
    *
    * @param folder
    *           - specification of the folder
    *
    * @return True if the folder exists, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean checkFolderExists(FolderSpec folder) throws Exception {
      validateFolderSpec(folder);

      _logger.info(String.format(
            "Checking whether folder '%s' exists",
            folder.name.get()));

      try {
         ManagedEntityUtil.getManagedObject(folder);
      } catch (ObjectNotFoundException e) {
         return false;
      }
      return true;
   }

   /**
    * Method that renames a folder
    * @param folder - the spec of teh folder to be renamed
    * @param name - the new name
    * @return true if successful, false otherwise
    * @throws Exception
    */
   public boolean renameFolder(FolderSpec folder, String name) throws Exception {
      validateFolderSpec(folder);

      _logger.info(String.format("Renaming folder '%s", folder.name.get()));

      // Get the folder
      Folder folderObj = ManagedEntityUtil.getManagedObject(folder);

      Future<ManagedObjectReference> blockingFuture =
            new BlockingFuture<ManagedObjectReference>();
      folderObj.rename(name, blockingFuture);

      // success or failure of the task
      ManagedObjectReference taskMoRef = blockingFuture.get();
      if (taskMoRef == null) {
         return false;
      }
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, folder)
            && ManagedEntityUtil.waitForEntityDeletion(
                  folder,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified folder exists and deletes it.
    *
    * @param folderSpec
    *           <code>FolderSpec</code> instance representing the fodler to be
    *           deleted.
    *
    * @return True if the folder doesn't exist or if the folder is deleted
    *         successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified folder
    *            doesn't exist.
    */
   public boolean deleteFolderSafely(FolderSpec folderSpec) throws Exception {

      if (checkFolderExists(folderSpec)) {
         return deleteFolder(folderSpec);
      }

      // Positive result if the folder doesn't exist
      return true;
   }

   /**
    * Retrieves the root folder of the VC inventory.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   Folder getRootFolder(ServiceSpec serviceSpec) throws Exception {
      VcService service = VcServiceUtil.getVcService(serviceSpec);

      // Get the root folder of the VC inventory
      ManagedObjectReference rootFolderMoRef =
            service.getServiceInstanceContent().getRootFolder();
      return service.getManagedObject(rootFolderMoRef);
   }

   /**
    * Retrieves the vm folder of specified datacenter
    *
    * @param datacenterSpec   specification of the datacenter
    * @throws Exception       if login to the VC service fails
    */
   Folder getVmFolder(DatacenterSpec datacenterSpec) throws Exception {
      Datacenter datacenterObj = ManagedEntityUtil.getManagedObject(datacenterSpec);

      ManagedObjectReference vmFolderMoRef = datacenterObj.getVmFolder();
      return (Folder) VcServiceUtil.getVcService(datacenterSpec).getManagedObject(
            vmFolderMoRef);
   }

   /**
    * Retrieves the host folder of specified datacenter
    *
    * @param datacenterSpec
    *           - specification of the datacenter
    * @throws Exception
    *            If login to the VC service fails
    */
   Folder getHostFolder(DatacenterSpec datacenterSpec) throws Exception {
      Datacenter datacenterObj = ManagedEntityUtil.getManagedObject(datacenterSpec);

      ManagedObjectReference hostFolderMoRef = datacenterObj.getHostFolder();
      return (Folder) VcServiceUtil.getVcService(datacenterSpec)
            .getManagedObject(hostFolderMoRef);
   }

   /**
    * Retrieves the network folder of specified datacenter
    *
    * @param datacenterSpec
    *           - specification of the datacenter
    * @throws Exception
    *            If login to the VC service fails
    */
   Folder getNetworkFolder(DatacenterSpec datacenterSpec) throws Exception {
      Datacenter datacenterObj = ManagedEntityUtil.getManagedObject(datacenterSpec);

      ManagedObjectReference networkFolderMoRef = datacenterObj.getNetworkFolder();
      return VcServiceUtil.getVcService(datacenterSpec).getManagedObject(
            networkFolderMoRef);
   }

   /**
    * Retrieves the storage folder of specified datacenter
    *
    * @param datacenterSpec
    *           - specification of the datacenter
    * @throws Exception
    *            If login to the VC service fails
    */
   Folder getStorageFolder(DatacenterSpec datacenterSpec) throws Exception {
      Datacenter datacenterObj = ManagedEntityUtil.getManagedObject(datacenterSpec);

      ManagedObjectReference storageFolderMoRef = datacenterObj.getDatastoreFolder();
      return VcServiceUtil.getVcService(datacenterSpec).getManagedObject(
            storageFolderMoRef);
   }

   // Get the spec of the parent Datacenter of the folder who might not be
   // direct parent of the folder if there is such
   private DatacenterSpec getParentDcSpec(FolderSpec folderSpec) {
      ManagedEntitySpec tempSpec = folderSpec;
      while (tempSpec.parent.isAssigned()) {
         if (tempSpec.parent.get() instanceof DatacenterSpec) {
            return (DatacenterSpec) tempSpec.parent.get();
         }
      }

      return null;
   }

   // Method that validates the folder spec has a non-empty name assigned
   private void validateFolderSpec(FolderSpec folderSpec) {
      if (!folderSpec.name.isAssigned() || folderSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("Folder name is not set.");
      }
   }
}
