/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.Datacenter;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.fault.SSLVerifyFault;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * Class that allows creation, deletion and check for existence of a Datacenter
 * in the VC Inventory based on a Datacenter specification
 *
 */
public class DatacenterBasicSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(DatacenterBasicSrvApi.class);

   private static DatacenterBasicSrvApi instance = null;

   protected DatacenterBasicSrvApi() {
   }

   /**
    * Get instance of DatacenterSrvApi.
    *
    * @return created instance
    */
   public static DatacenterBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (DatacenterBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing DatacenterSrvApi.");
               instance = new DatacenterBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a datacenter with name as its parameter in the root folder
    *
    * @param name
    *           - string with the name of the datacenter to be created
    *
    * @return True if the creation was successful, false otherwise.
    * @throws Exception
    */
   public boolean createDatacenter(DatacenterSpec datacenterSpec)
         throws Exception {
      validateDatacenterSpec(datacenterSpec);

      _logger.info(String.format("Creating datacenter '%s'", datacenterSpec.name.get()));

      // Get the host folder of the datacenter
      Folder parentFolder = null;

      if (!datacenterSpec.parent.isAssigned()
            || !datacenterSpec.parent.get().name.isAssigned()
            || datacenterSpec.parent.get() instanceof VcSpec) {
         parentFolder = FolderBasicSrvApi.getInstance().getRootFolder(datacenterSpec.service.get());
      } else {
         parentFolder = ManagedEntityUtil.getManagedObject(datacenterSpec.parent.get());
      }

      return parentFolder.createDatacenter(datacenterSpec.name.get()) != null;
   }

   /**
    * Renames a datacenter
    * 
    * @param datacenterSpec - the datacenter to rename
    * @param newName - the new name
    */
   public boolean renameDatacenter(DatacenterSpec datacenterSpec, String newName) throws Exception {
      validateDatacenterSpec(datacenterSpec);

      _logger.info(String.format("Renaming datacenter '%s", datacenterSpec.name.get()));

      // Get the datacenter
      Datacenter datacenter = ManagedEntityUtil.getManagedObject(datacenterSpec, datacenterSpec.service.get());

      ManagedObjectReference taskMoRef = datacenter.rename(newName);

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, datacenterSpec.service.get());
   }

   /**
    * Deletes the specified datacenter from the inventory and wait the specified
    * timeout for the entity to get deleted, otherwise return false if timeout
    * is over or the entity cannot be deleted
    *
    * @param datacenterSpec
    *           - specification of the Datacenter to be deleted
    * @return True if the deletion was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteDatacenter(DatacenterSpec datacenterSpec)
         throws Exception {
      validateDatacenterSpec(datacenterSpec);

      _logger.info(String.format("Deleting datacenter '%s", datacenterSpec.name.get()));

      // Get the datacenter
      Datacenter datacenter =
            ManagedEntityUtil.getManagedObject(
                  datacenterSpec,
                  datacenterSpec.service.get());

      ManagedObjectReference taskMoRef = datacenter.destroy();

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, datacenterSpec.service.get())
            && ManagedEntityUtil.waitForEntityDeletion(
                  datacenterSpec,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified datacenter exists and deletes it.
    *
    * @param datacenterSpec
    *           <code>DatacenterSpec</code> instance representing the datacenter
    *           to be deleted.
    *
    * @return True if the datacenter doesn't exist or if the datacenter is
    *         deleted successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified datacenter
    *            doesn't exist.
    */
   public boolean deleteDatacenterSafely(DatacenterSpec datacenterSpec)
         throws Exception {

      if (checkDatacenterExists(datacenterSpec)) {
         return deleteDatacenter(datacenterSpec);
      }

      // Positive result if the datacenter doesn't exist
      return true;
   }

   /**
    * Checks whether the specified datacenter exists.
    *
    * @param datacenterSpec
    *           - specification of the datacenter
    *
    * @return True is the datacenter exists, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean checkDatacenterExists(DatacenterSpec datacenterSpec)
         throws Exception {
      validateDatacenterSpec(datacenterSpec);

      _logger.info(String.format(
            "Checking whether datacenter '%s' exists",
            datacenterSpec.name.get()));

      try {
         ManagedEntityUtil.getManagedObject(datacenterSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }

      return true;
   }

   /**
    * Method to validate the connection between the host and vc
    *
    * @param hostSpec
    *           - specification of the host the connection to which to validate
    * @param datacenter
    *           - specification of the datacenter to which the host will be
    *           added
    * @return if the connection is validated, the thumbprint returned will be
    *         null and the host addition can proceed, if the connection is not
    *         valid a thumbprint will be returned that has to be set in the host
    *         specification and only after that the host can be added
    * @throws Exception
    */
   public String validateSslThumbprint(HostSpec hostSpec, DatacenterSpec datacenterSpec)
         throws Exception {
      String thumbprint = null;
      Datacenter datacenter = ManagedEntityUtil.getManagedObject(datacenterSpec);
      try {
         datacenter.queryConnectionInfo(
               hostSpec.name.get(),
               hostSpec.port.get(),
               hostSpec.userName.get(),
               hostSpec.password.get(),
               null);
      } catch (SSLVerifyFault sslFault) {
         thumbprint = sslFault.getThumbprint();
      }
      return thumbprint;
   }

   private void validateDatacenterSpec(DatacenterSpec datacenterSpec)
         throws IllegalArgumentException {
      if (!datacenterSpec.name.isAssigned() || datacenterSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("Datacenter name is not set");
      }
   }
}
