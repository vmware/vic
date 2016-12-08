/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.SsoUtil;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.ClusterComputeResource;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.Network;
import com.vmware.vim.binding.vim.cluster.ConfigSpecEx;
import com.vmware.vim.binding.vim.cluster.DasConfigInfo;
import com.vmware.vim.binding.vim.cluster.DrsConfigInfo;
import com.vmware.vim.binding.vim.cluster.DrsConfigInfo.DrsBehavior;
import com.vmware.vim.binding.vim.cluster.FailoverHostAdmissionControlPolicy;
import com.vmware.vim.binding.vim.cluster.FailoverLevelAdmissionControlPolicy;
import com.vmware.vim.binding.vim.cluster.FailoverResourcesAdmissionControlPolicy;
import com.vmware.vim.binding.vim.cluster.SlotPolicy;
import com.vmware.vim.binding.vim.vsan.cluster.ConfigInfo;
import com.vmware.vim.binding.vim.vsan.cluster.ConfigInfo.HostDefaultInfo;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.query.client.exception.NotImplementedException;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * Provides utility methods for performing operations on Virtual Center clusters
 * via the Vim API.
 */
public class ClusterBasicSrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(ClusterBasicSrvApi.class);

   private static ClusterBasicSrvApi instance = null;

   protected ClusterBasicSrvApi() {
   }

   /**
    * Get instance of ClusterSrvApi.
    *
    * @return created instance
    */
   public static ClusterBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (ClusterBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing ClusterSrvApi.");
               instance = new ClusterBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a cluster with properties as specified in the <code>ClusterSpec
    * </code> parameter.
    *
    * @param clusterSpec
    *           <code>ClusterSpec</code> containing the properties of the
    *           cluster to be created.
    *
    * @return True if the creation was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified cluster
    *            settings are invalid.
    */
   public boolean createCluster(ClusterSpec clusterSpec)
         throws Exception {
      validateClusterSpec(clusterSpec);

      _logger.info(String.format("Creating cluster '%s' on datacenter '%s'",
            clusterSpec.name.get(), clusterSpec.parent.get().name.get()));

      // Get the host folder of the datacenter
      DatacenterSpec datacenterSpec = (DatacenterSpec) clusterSpec.parent.get();
      Folder hostFolder = FolderBasicSrvApi.getInstance().getHostFolder(datacenterSpec);

      // Build the necessary parameters for cluster creation
      ConfigSpecEx configSpec = buildClusterConfigSpecEx(clusterSpec);

      return hostFolder.createClusterEx(clusterSpec.name.get(), configSpec) != null;
   }

   /**
    * Deletes the specified cluster from the inventory and waits for timeout
    * seconds to get the object really deleted. If timeout is over and object
    * still in inventory, it will return false
    *
    * @param clusterSpec
    *           <code>ClusterSpec</code> instance representing the cluster to be
    *           deleted.
    *
    * @return True if the deletion was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified cluster
    *            settings are invalid.
    */
   public boolean deleteCluster(ClusterSpec clusterSpec)
         throws Exception {

      validateClusterSpec(clusterSpec);

      _logger.info(String.format("Deleting cluster '%s' from datacenter '%s'",
            clusterSpec.name.get(), clusterSpec.parent.get().name.get()));

      ClusterComputeResource cluster = ManagedEntityUtil
            .getManagedObject(clusterSpec);

      ManagedObjectReference taskMoRef = cluster.destroy();

      return VcServiceUtil.waitForTaskSuccess(taskMoRef,
            clusterSpec.service.get())
            && ManagedEntityUtil.waitForEntityDeletion(clusterSpec,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Method to reconfigure the cluster, currently it does so only for DRS
    *
    * @param clusterToReconfigure
    *           - spec fo teh cluster that should be reconfigured
    * @param clusterNewSettings
    *           - the spec containing the new settings to be applied
    * @return true if it was successful, false otherwise
    * @throws Exception
    *            - in case VC connection fails
    */
   public boolean reconfigureCluster(ClusterSpec clusterToReconfigure,
         ClusterSpec clusterNewSettings)

         throws Exception {

      validateClusterSpec(clusterToReconfigure);

      _logger.info(String.format(
            "Reconfiguring cluster '%s' from datacenter '%s'",
            clusterToReconfigure.name.get(),
            clusterToReconfigure.parent.get().name.get()));

      ClusterComputeResource cluster = ManagedEntityUtil
            .getManagedObject(clusterToReconfigure);
      // modify = true means that the configuration will be applied
      // incrementally, i.e. it applies only partial changes;
      // if modify is set to false then the unset values in the reconfigure spec
      // will be set to defaults or unset on cluster
      // TODO it currently reconfigures cluster only for DRS - on/off and
      // behavior,
      // in order to be more general buildClusterConfigSpec should be reworked
      ManagedObjectReference taskMoRef = cluster.reconfigureEx(
            buildClusterConfigSpecEx(clusterNewSettings), true);

      return VcServiceUtil.waitForTaskSuccess(taskMoRef,
            clusterToReconfigure.service.get());
   }

   /**
    * Method that allows the renaming of an existing cluster
    *
    * @param clusterToRename
    *           - specification of the cluster to be renamed
    * @param clusterNewName
    *           - cluster specification that contains the new name
    * @return - true if the task was successful, otherwise false
    * @throws Exception
    *            - if the connection to VC fails or if the new name is not
    *            assigned or is empty
    */
   public boolean renameCluster(ClusterSpec clusterToRename,
         ClusterSpec clusterNewName) throws Exception {
      validateClusterSpec(clusterToRename);

      _logger
            .info(String.format(
                  "Reconfiguring cluster '%s' from datacenter '%s'",
                  clusterToRename.name.get(),
                  clusterToRename.parent.get().name.get()));

      if (Strings.isNullOrEmpty(clusterNewName.name.get())) {
         throw new IllegalArgumentException(
               "Cannot rename cluster with empty name");
      }

      ClusterComputeResource cluster = ManagedEntityUtil
            .getManagedObject(clusterToRename);
      ManagedObjectReference taskMoRef = cluster.rename(clusterNewName.name
            .get());

      return VcServiceUtil.waitForTaskSuccess(taskMoRef,
            clusterToRename.service.get());

   }

   /**
    * Checks whether the specified cluster exists and deletes it.
    *
    * @param clusterSpec
    *           <code>ClusterSpec</code> instance representing the cluster to be
    *           deleted.
    *
    * @return True if the cluster doesn't exist or if the cluster is deleted
    *         successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified datacenter
    *            doesn't exist.
    */
   public boolean deleteClusterSafely(ClusterSpec clusterSpec)
         throws Exception {

      if (checkClusterExists(clusterSpec)) {
         return deleteCluster(clusterSpec);
      }

      // Positive result if the cluster doesn't exist
      return true;
   }

   /**
    * Checks whether the specified cluster exists in the specified datacenter.
    *
    * @param clusterSpec
    *           <code>ClusterSpec</code> instance representing the cluster to be
    *           queried.
    *
    * @return True is the cluster exists, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified datacenter
    *            doesn't exist.
    */
   public boolean checkClusterExists(ClusterSpec clusterSpec)
         throws Exception {

      validateClusterSpec(clusterSpec);

      _logger.info(String.format(
            "Checking whether cluster '%s' exists in datacenter '%s'",
            clusterSpec.name.get(), clusterSpec.parent.get().name.get()));

      try {
         ManagedEntityUtil.getManagedObject(clusterSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }

      return true;
   }

   /**
    * Retrieves list of cluster network names.
    *
    * @param clusterSpec
    *           the spec for the cluster that will be searched
    */
   public List<String> getNetworkNames(ClusterSpec clusterSpec) {
      List<String> networkNames = new ArrayList<String>();
      ClusterComputeResource cluster = null;
      VcService service = null;

      try {
         service = SsoUtil.getVcConnector(clusterSpec).getConnection()
               .getVcService();
         cluster = ManagedEntityUtil.getManagedObject(clusterSpec,
               clusterSpec.service.get());
      } catch (Exception e) {
         String errorMessage = String.format("Cannot retrieve cluster %s",
               clusterSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }

      for (ManagedObjectReference networkMor : cluster.getNetwork()) {
         Network network = service.getManagedObject(networkMor);
         networkNames.add(network.getName());
      }

      return networkNames;
   }

   /**
    * Retrieves list of cluster datastore names.
    *
    * @param clusterSpec
    *           the spec of the cluster that will be searched
    * @throws RuntimeException
    *            if login to the VC server fails or cluster is not present
    */
   public static List<String> getDatastoreNames(ClusterSpec clusterSpec) {
      ClusterComputeResource cluster = null;
      VcService service = null;
      try {
         service = VcServiceUtil.getVcService(clusterSpec);
         cluster = ManagedEntityUtil.getManagedObject(clusterSpec);
      } catch (Exception e) {
         String errorMessage = String.format("Cannot retrieve cluster %s", clusterSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }

      List<String> datastoreNames = new ArrayList<String>();
      for (ManagedObjectReference datastoreRef : cluster.getDatastore()) {
         Datastore datastore = service.getManagedObject(datastoreRef);
         datastoreNames.add(datastore.getName());
      }

      return datastoreNames;
   }

   /**
    * Checks if DRS is enabled on the specified cluster.
    *
    * @param clusterSpec
    *           <code>ClusterSpec</code> instance representing the cluster to
    *           check
    * @return true if DRS is enabled on the cluster, false otherwise
    * @throws Exception
    */
   public boolean isDrsEnabled(ClusterSpec clusterSpec) throws Exception {
      ClusterComputeResource cluster = null;
      try {
         cluster = ManagedEntityUtil.getManagedObject(clusterSpec);
      } catch (Exception e) {
         String errorMessage = String.format("Cannot retrieve cluster %s", clusterSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }

      return cluster.getConfiguration().getDrsConfig().getEnabled();
   }

   /*
    * Creates a ConfigSpecEx instance necessary for creation of a DRS-enabled
    * cluster.
    */
   private ConfigSpecEx buildClusterConfigSpecEx(ClusterSpec clusterSpec)
         throws NotImplementedException {
      ConfigSpecEx configSpec = new ConfigSpecEx();

      // Setup DRS
      DrsConfigInfo drsConfig = new DrsConfigInfo();
      if (clusterSpec.drsEnabled.isAssigned()) {
         drsConfig.setEnabled(clusterSpec.drsEnabled.get());
      } else {
         drsConfig.setEnabled(true);
      }

      if (clusterSpec.drsBehavior.isAssigned()) {
         switch (clusterSpec.drsBehavior.get()) {
         case MANUAL:
            drsConfig.setDefaultVmBehavior(DrsBehavior.manual);
            break;
         case PARTIALLY_AUTOMATED:
            drsConfig.setDefaultVmBehavior(DrsBehavior.partiallyAutomated);
            break;
         case FULLY_AUTOMATED:
         default:
            drsConfig.setDefaultVmBehavior(DrsBehavior.fullyAutomated);
            break;
         }
      }

      configSpec.setDrsConfig(drsConfig);

      // Configure High Availability
      if (clusterSpec.vsphereHA.isAssigned()) {
         DasConfigInfo dasConfig = new DasConfigInfo();
         dasConfig.setEnabled(clusterSpec.vsphereHA.get());
         dasConfig.setVmMonitoring("enabled");
         dasConfig.setFailoverLevel(1);

         if (clusterSpec.admissionControlSpec.isAssigned()) {
            switch (clusterSpec.admissionControlSpec.get().failoverPolicy.get()) {
            case CLUSTER_RESOURCE_PERCENTAGE:
               // TODO Finish implementation - use
               // FailoverResourcesAdmissionControlPolicy
               throw new NotImplementedException(
                     "Cluster resource percentage option not implemented!");
            case DEDICATED_FAILOVER_HOSTS:
               // TODO Finish implementation - use
               // FailoverHostAdmissionControlPolicy
               throw new NotImplementedException(
                     "Dedicated failover hosts option not implemented!");
            case SLOT_POLICY:
               // TODO Finish implementation - use
               // FailoverLevelAdmissionControlPolicy
               throw new NotImplementedException(
                     "Slot policy option not implemented!");
            case DISABLED:
            default:
               dasConfig.setAdmissionControlEnabled(false);
               break;
            }
         } else {
            FailoverLevelAdmissionControlPolicy admissionPolicy = new FailoverLevelAdmissionControlPolicy();
            admissionPolicy.setFailoverLevel(1);
            dasConfig.setAdmissionControlPolicy(admissionPolicy);
            dasConfig.setAdmissionControlEnabled(true);
         }

         dasConfig.setHBDatastoreCandidatePolicy("allFeasibleDsWithUserPreference");
         configSpec.setDasConfig(dasConfig);
      }

      // Configure SPBM (storage policies)
      if (clusterSpec.isSpbmEnabled.isAssigned()) {
         configSpec.setSpbmEnabled(clusterSpec.isSpbmEnabled.get());
      }

      // Setup vSAN
      if (clusterSpec.vsanEnabled.isAssigned()) {
         ConfigInfo vsanConfig = new ConfigInfo();

         vsanConfig.setEnabled(clusterSpec.vsanEnabled.get());

         HostDefaultInfo hostDefaultConfig = new HostDefaultInfo();
         if (clusterSpec.vsanAutoClaimStorage.isAssigned()) {
            hostDefaultConfig.setAutoClaimStorage(clusterSpec.vsanAutoClaimStorage.get());
         } else {
            hostDefaultConfig.setAutoClaimStorage(true);
         }
         vsanConfig.setDefaultConfig(hostDefaultConfig);

         configSpec.setVsanConfig(vsanConfig);
      }

      return configSpec;
   }

   private void validateClusterSpec(ClusterSpec clusterSpec) {
      if (!clusterSpec.name.isAssigned()) {
         throw new IllegalArgumentException("Cluster name is not set");
      }

      if (!clusterSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("Cluster parent is not set");
      }
   }
}
