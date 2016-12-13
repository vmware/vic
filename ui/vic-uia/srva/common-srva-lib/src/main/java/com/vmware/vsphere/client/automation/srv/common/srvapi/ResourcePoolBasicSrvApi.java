/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.ClusterComputeResource;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.ResourceAllocationInfo;
import com.vmware.vim.binding.vim.ResourceConfigSpec;
import com.vmware.vim.binding.vim.ResourcePool;
import com.vmware.vim.binding.vim.SharesInfo;
import com.vmware.vim.binding.vim.SharesInfo.Level;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

public class ResourcePoolBasicSrvApi {
   protected static final String DEFAULT_VM_GUEST_OS = "winVistaGuest";

   private static final Logger _logger =
         LoggerFactory.getLogger(ResourcePoolBasicSrvApi.class);

   private static ResourcePoolBasicSrvApi instance = null;
   protected ResourcePoolBasicSrvApi() {}

   /**
    * Get instance of ResourcePoolSrvApi.
    *
    * @return  created instance
    */
   public static ResourcePoolBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized(ResourcePoolBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing ResourcePoolSrvApi.");
               instance = new ResourcePoolBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates resource pool with properties as specified in the
    * <code>ResourcePoolSpec</code> parameter. Currently creation is possible
    * only on host and cluster.
    *
    * @param resPoolSpec
    *           <code>ResourcePoolSpec</code> containing the properties of the
    *           resource pool to be created.
    *
    * @return True if the creation was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean createResourcePool(ResourcePoolSpec resPoolSpec)
         throws Exception {
      validateResourcePoolSpec(resPoolSpec);

      ServiceSpec serviceSpec = resPoolSpec.service.get();
      VcService vcService = VcServiceUtil.getVcService(serviceSpec);
      ResourcePool parentRp = null;
      boolean resourcePoolCreated = false;


      if (resPoolSpec.parent.get() instanceof HostSpec) {
         HostSpec parentHostSpec = (HostSpec) resPoolSpec.parent.get();
         HostSystem vAppHost = ManagedEntityUtil.getManagedObject(
               parentHostSpec, serviceSpec);

         ComputeResource cr = (ComputeResource) vcService
               .getManagedObject(vAppHost.getParent());

         parentRp = (ResourcePool) vcService.getManagedObject(cr
               .getResourcePool());
      } else if (resPoolSpec.parent.get() instanceof ClusterSpec) {
         ClusterSpec parentClusterSpec = (ClusterSpec) resPoolSpec.parent.get();
         ClusterComputeResource vAppHost = ManagedEntityUtil.getManagedObject(
               parentClusterSpec, serviceSpec);

         parentRp = (ResourcePool) vcService.getManagedObject(vAppHost
               .getResourcePool());
      }

      // TODO add creation of resource pool on other resource pool or vApp

      if (parentRp != null) {
         resourcePoolCreated = parentRp.createResourcePool(
               resPoolSpec.name.get(), getDefaultResourceConfiguration()) != null;
      }

      return resourcePoolCreated;
   }

   /**
    * Deletes the specified resource pool from the inventory and waits for
    * timeout seconds to get the object really deleted. If timeout is over and
    * object still in inventory, it will return false
    *
    * @param resPoolSpec
    *           <code>ResourcePoolSpec</code> instance representing the resource
    *           pool to be deleted.
    *
    * @return True if the deletion was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteResourcePool(ResourcePoolSpec resPoolSpec)
         throws Exception {

      validateResourcePoolSpec(resPoolSpec);

      ServiceSpec serviceSpec = resPoolSpec.service.get();

      _logger.info(String.format(
            "Deleting resource pool '%s' from parent '%s'",
            resPoolSpec.name.get(), resPoolSpec.parent.get().name.get()));

      ResourcePool resPool = ManagedEntityUtil.getManagedObject(resPoolSpec,
            serviceSpec);

      ManagedObjectReference taskMoRef = resPool.destroy();

      return VcServiceUtil.waitForTaskSuccess(taskMoRef, serviceSpec)
            && ManagedEntityUtil.waitForEntityDeletion(resPoolSpec,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified resource pool exists and deletes it.
    *
    * @param resPoolSpec
    *           <code>ResourcePoolSpec</code> instance representing the resource
    *           pool to be deleted.
    *
    * @return True if the resource pool doesn't exist or if the resource pool is
    *         deleted successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteResourcePoolSafely(ResourcePoolSpec resPoolSpec)
         throws Exception {

      if (checkResourcePoolExists(resPoolSpec)) {
         return deleteResourcePool(resPoolSpec);
      }

      // Positive result if the resource pool doesn't exist
      return true;
   }

   /**
    * Updates the resource pool, currently only changes name
    * @param originalResPoolSpec - resource pool to update
    * @param newResPoolSpec - spec with new values
    * @return true if successful, false if not or if no such pool exists
    * @throws Exception
    */
   public boolean updateResourcePool(ResourcePoolSpec originalResPoolSpec, ResourcePoolSpec newResPoolSpec)
       throws Exception {

    if (checkResourcePoolExists(originalResPoolSpec)) {

        validateResourcePoolSpec(newResPoolSpec);

        ServiceSpec serviceSpec = originalResPoolSpec.service.get();

        _logger.info(String.format(
              "Updating resource pool '%s' from parent '%s'",
              originalResPoolSpec.name.get(), originalResPoolSpec.parent.get().name.get()));

        ResourcePool resPool = ManagedEntityUtil.getManagedObject(originalResPoolSpec,
              serviceSpec);
        // here actually only name is updated and the old resPool configuration is used
        // TODO if needed add parameters for configuration to ResourcePoolSpec and change
        // as necessary
        resPool.updateConfig(newResPoolSpec.name.get(), resPool.getConfig());

        return checkResourcePoolExists(newResPoolSpec) && !checkResourcePoolExists(originalResPoolSpec);
    }

    // resource pool doesn't exist
    return false;
 }

   /**
    * Checks whether the specified resource pool exists.
    *
    * @param resPoolSpec
    *           <code>ResourcePoolSpec</code> instance representing the resource
    *           pool to be queried.
    *
    * @return True is the resource pool exists, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean checkResourcePoolExists(ResourcePoolSpec resPoolSpec)
         throws Exception {

      validateResourcePoolSpec(resPoolSpec);

      ServiceSpec serviceSpec = resPoolSpec.service.get();

      _logger.info(String.format(
            "Checking whether resource pool '%s' exists in parent '%s'",
            resPoolSpec.name.get(), resPoolSpec.parent.get().name.get()));

      try {
         ManagedEntityUtil.getManagedObject(resPoolSpec, serviceSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }

      return true;
   }

   // ---------------------------------------------------------------------------
   // Private methods

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

      // set defaults for memory allocation
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
    * Method that validates the resource pool spec: <br>
    * - name should be assigned <br>
    * - parent host should be assigned<br>
    *
    * @param resourcePoolSpec
    *           Resource pool spec that will be validated
    * @throws IllegalArgumentException
    *            if resource pool spec requirements are not met
    */
   private void validateResourcePoolSpec(
         ResourcePoolSpec resourcePoolSpec) throws IllegalArgumentException {
      if (!resourcePoolSpec.name.isAssigned()
            || resourcePoolSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("Resource pool name is not set.");
      }

      if (!resourcePoolSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("Resource pool parent is not set.");
      }

      EntitySpec resourcePoolParentEntity = resourcePoolSpec.parent.get();
      if (!(resourcePoolParentEntity instanceof HostSpec)
            && !(resourcePoolParentEntity instanceof ClusterSpec)
            && !(resourcePoolParentEntity instanceof DatacenterSpec)
            && !(resourcePoolParentEntity instanceof FolderSpec)) {
         throw new IllegalArgumentException(
               "Resource pool parent association is not a host, cluster, datacenter or folder.");
      }
   }
}
