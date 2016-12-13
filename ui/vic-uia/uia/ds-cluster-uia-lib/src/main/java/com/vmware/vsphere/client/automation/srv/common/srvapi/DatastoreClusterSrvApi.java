/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.client.automation.vcuilib.commoncode.ActionFunction;
import com.vmware.suitaf.SUITA;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.StoragePod;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * Class that encapsulates basic API operations related to a Datastore cluster
 * in the VC Inventory based on a DatastoreClusterSpec.
 *
 */
public class DatastoreClusterSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(ActionFunction.class);

   private static DatastoreClusterSrvApi instance = null;
   protected DatastoreClusterSrvApi() {}

   /**
    * Get instance of DatastoreClusterSrvApi.
    *
    * @return  created instance
    */
   public static DatastoreClusterSrvApi getInstance() {
      if (instance == null) {
         synchronized(DatastoreClusterSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing DatastoreClusterSrvApi.");
               instance = new DatastoreClusterSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates datastore cluster
    *
    * @param dsClusterSpec
    *           - spec for the ds cluster
    * @return - true if creation is successful, fals otherwise
    * @throws Exception
    */
   public boolean createDatastoreCluster(DatastoreClusterSpec dsClusterSpec)
         throws Exception {
      validateDsClusterSpec(dsClusterSpec);

      _logger.info(String.format("Creating datastore cluster '%s' on datacenter '%s'",
            dsClusterSpec.name.get(), dsClusterSpec.parent.get().name.get()));

      // Get the storage folder of the datacenter
      DatacenterSpec datacenterSpec = (DatacenterSpec) dsClusterSpec.parent.get();
      Folder storageFolder = FolderBasicSrvApi.getInstance().getStorageFolder(datacenterSpec);

      if (storageFolder.createStoragePod(dsClusterSpec.name.get()) == null) {
         throw new Exception(String.format(
               "Unable to create datastore cluster '%s'",
               dsClusterSpec.name.get()));
      }

      return true;
   }

   /**
    * Check for Datastore cluster existence
    *
    * @param dsClusterSpec
    *           - spec for the ds cluster
    * @return true if the ds cluster exists
    * @throws Exception
    */
   public boolean checkDsClusterExists(DatastoreClusterSpec dsClusterSpec)
         throws Exception {

      _logger.info(String.format(
            "Checking for datastore cluster existence '%s' exists in datacenter '%s'",
            dsClusterSpec.name.get(), dsClusterSpec.parent.get()));

      try {
         ManagedEntityUtil.getManagedObject(dsClusterSpec, dsClusterSpec.service.get());
      } catch (ObjectNotFoundException onfe) {
         return false;
      }
      return true;
   }

   /**
    * Moves specified datastores to datastore cluster.
    *
    * @param dsClusterSpec
    *           - spec for the datastore cluster
    * @param datastores
    *           - specs for the datastores to be moved to datastore cluster
    * @return true if the ds cluster exists
    * @throws Exception
    */
   public boolean moveDatastoresToDsCluster(List<DatastoreSpec> datastores,
         DatastoreClusterSpec dsClusterSpec) throws Exception {

      _logger.info(String.format("Moving datastores to datastore cluster '%s'",
            dsClusterSpec.name.get()));

      int dsClusterMembersCount = getDsClusterMembers(dsClusterSpec).size();

      List<ManagedObjectReference> datastoreMors = new LinkedList<ManagedObjectReference>();
      for (DatastoreSpec datastore : datastores) {
         datastoreMors
               .add(ManagedEntityUtil.getManagedObject(datastore, datastore.service.get())._getRef());
      }
      StoragePod dsCluster = ManagedEntityUtil.getManagedObject(dsClusterSpec, dsClusterSpec.service.get());
      dsCluster.moveInto(datastoreMors.toArray(new ManagedObjectReference[datastoreMors
            .size()]));

      // Verify the datastores were moved successfully
      if ((dsClusterMembersCount + datastores.size()) == getDsClusterMembers(
            dsClusterSpec).size()) {
         return true;
      } else {
         return false;
      }
   }

   /**
    * Removes datastore cluster
    *
    * @param dsClusterSpec
    *           - spec for the ds cluster
    * @return true if the ds cluster is successfully removed
    * @throws Exception
    */
   public boolean deleteDsCluster(DatastoreClusterSpec dsClusterSpec)
         throws Exception {
      _logger.info(String.format("Deleting datastore cluster '%s' from datacenter '%s'",
            dsClusterSpec.name.get(), dsClusterSpec.parent.get()));
      StoragePod dsCluster = ManagedEntityUtil.getManagedObject(dsClusterSpec,
            dsClusterSpec.service.isAssigned() ? dsClusterSpec.service.get() : null);
      ManagedObjectReference taskMoRef = dsCluster.destroy();
      return VcServiceUtil.waitForTaskSuccess(taskMoRef, dsClusterSpec)
            && ManagedEntityUtil.waitForEntityDeletion(dsClusterSpec,
                  (int) SUITA.Environment.getBackendJobSmall() / 1000);
   }

   /**
    * Checks whether SDRS is enabled for datastore cluster
    *
    * @param dsClusterSpec
    *           - spec for the ds cluster
    * @return true if SDRS is enabled
    * @throws Exception
    */
   public boolean isSdrsEnabled(DatastoreClusterSpec dsClusterSpec)
         throws Exception {
      _logger
            .info(String
                  .format(
                        "Checking whether SDRS is enabled on datastore cluster '%s' from datacenter '%s'",
                        dsClusterSpec.name.get(), dsClusterSpec.parent.get()));
      StoragePod dsCluster = ManagedEntityUtil.getManagedObject(dsClusterSpec, dsClusterSpec.service.get());
      boolean sdrsEnabled = dsCluster.getPodStorageDrsEntry().getStorageDrsConfig().getPodConfig().isEnabled();
      return sdrsEnabled;
   }

   /**
    * Gets datastore members of the datastore cluster
    *
    * @param dsClusterSpec
    *           - spec for the ds cluster
    * @return List of datastore names under the datastore cluster
    * @throws Exception
    */
   public List<String> getDsClusterMembers(DatastoreClusterSpec dsClusterSpec)
         throws Exception {
      VcService service = null;
      ServiceSpec serviceSpec = dsClusterSpec.service.get();
      service = VcServiceUtil.getVcService(serviceSpec);
      _logger.info(String.format(
            "Getting datastores for ds cluster '%s' from datacenter '%s'",
            dsClusterSpec.name.get(), dsClusterSpec.parent.get()));
      StoragePod dsCluster = ManagedEntityUtil.getManagedObject(dsClusterSpec, dsClusterSpec.service.get());
      ManagedObjectReference[] members = dsCluster.getChildEntity();
      List<String> datastoreNames = new ArrayList<String>();
      for (ManagedObjectReference ds : members) {
         Datastore datastore = service.getManagedObject(ds);
         datastoreNames.add(datastore.getName());
      }

      return datastoreNames;
   }

   private void validateDsClusterSpec(DatastoreClusterSpec dsClusterSpec) {
      if (dsClusterSpec == null) {
         throw new IllegalArgumentException("DatastoreClusterSpec is null");
      }
      if (!dsClusterSpec.name.isAssigned()) {
         throw new IllegalArgumentException("Datastore cluster name is not set");
      }

      if (!dsClusterSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("Datastore cluster parent is not set");
      }
   }

   /**
    * Checks whether the specified datastore cluster exists and deletes it.
    *
    * @param dsClusterSpec
    *           <code>DatastoreClusterSpec</code> instance representing the datastore
    *           cluster to be deleted.
    *
    * @return true if the datastore cluster doesn't exist OR if the datastore cluster
    *            is deleted successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteDsClusterSafely(DatastoreClusterSpec dsClusterSpec)
         throws Exception {

      if (checkDsClusterExists(dsClusterSpec)) {
         return deleteDsCluster(dsClusterSpec);
      }

      // Positive result if the datastore cluster doesn't exist
      return true;
   }
}
