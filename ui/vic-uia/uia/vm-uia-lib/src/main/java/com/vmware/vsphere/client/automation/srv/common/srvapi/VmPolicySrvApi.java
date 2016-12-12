/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.vim.binding.pbm.ServerObjectRef;
import com.vmware.vim.binding.pbm.profile.ProfileId;
import com.vmware.vim.binding.pbm.profile.ProfileManager;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;

/**
 * API commands for managing the policies associated to VMs
 */
public class VmPolicySrvApi {

   private static final Logger _logger = LoggerFactory.getLogger(VmPolicySrvApi.class);
   private static VmPolicySrvApi instance = null;
   protected VmPolicySrvApi() {}

   /**
    * Get instance of VmPolicySrvApi.
    *
    * @return  created instance
    */
   public static VmPolicySrvApi getInstance() {
      if (instance == null) {
         synchronized(VmPolicySrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing VmPolicySrvApi.");
               instance = new VmPolicySrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Assign a storage policy to an entity.
    * @param entity - the entity we are going to assign the policy to
    * @param policy - the policy to assign
    * @return
    * @throws Exception - if the API connection fails for some reason
    */
   public boolean addStoragePolicy(ManagedEntitySpec entity,
         StoragePolicySpec policy) throws Exception {
      ProfileManager profileManager = StoragePolicyBasicSrvApi.getInstance().getProfileManager(entity);

      // Get all the connecting parties
      ManagedObject entityMo = ManagedEntityUtil.getManagedObject(entity);
      ProfileId profileId =
            StoragePolicyBasicSrvApi.getInstance().getProfileByName(policy.name.get(), profileManager);

      ServerObjectRef srvObjRef =
            ManagedEntityUtil.managedObjectRefToServerObjectRef(entityMo._getRef());

      BlockingFuture<Void> taskResult = new BlockingFuture<Void>();
      profileManager.associate(srvObjRef, profileId, null, taskResult);
      return taskResult.get() != null;
   }

   /**
    * Remove a storage policy to an entity.
    * @param entity - the entity we are going to unassign the policy from
    * @param policy - the policy to unassign
    * @return
    * @throws Exception - if the API connection fails for some reason
    */
   public boolean removeStoragePolicy(ManagedEntitySpec entity,
         StoragePolicySpec policy) throws Exception {
      ProfileManager profileManager = StoragePolicyBasicSrvApi.getInstance().getProfileManager(entity);

      // Get all the connecting parties
      ManagedObject entityMo = ManagedEntityUtil.getManagedObject(entity);
      ProfileId profileId =
            StoragePolicyBasicSrvApi.getInstance().getProfileByName(policy.name.get(), profileManager);

      ServerObjectRef srvObjRef =
            ManagedEntityUtil.managedObjectRefToServerObjectRef(entityMo._getRef());

      BlockingFuture<Void> taskResult = new BlockingFuture<Void>();
      profileManager.dissociate(srvObjRef, profileId, taskResult);
      return taskResult.get() != null;
   }

}
