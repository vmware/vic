/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;

/**
 * A location spec suitable for modeling tree navigation in the Resource pool
 * related tests.
 */
public class ResourcePoolTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec,
         TreeTabIDs entityViewNId, String primaryTabNId,
         String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), resourcePoolSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.RESOURCE_POOL;
   }

   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec,
         String primaryTabNId, String secondaryTabNId) {
      this(resourcePoolSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec,
         TreeTabIDs entityTabNId, String primaryTabNId, String secondaryTabNId) {
      this(resourcePoolSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec,
         String primaryTabNId) {
      this(resourcePoolSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId,
            null, null);
   }

   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec,
         TreeTabIDs entityTabNId, String primaryTabNId) {
      this(resourcePoolSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec) {
      this(resourcePoolSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, null, null, null);
   }

   /**
    * Build a location path based on the provided Resource pool navigation
    * identifiers.
    *
    * @param resourcePoolSpec
    *           the resource pool spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public ResourcePoolTreeLocationSpec(ResourcePoolSpec resourcePoolSpec,
         TreeTabIDs entityTabNId) {
      this(resourcePoolSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the Resource pool entity
    * view.
    */
   public ResourcePoolTreeLocationSpec() {
      super(TreeTabIDs.HOSTS_AND_CLUSTERS.getTreeTabID());
      _entityType = EntityTypes.RESOURCE_POOL;
   }
}