/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;

/**
 * A location spec suitable for modeling tree navigation in the Cluster related
 * tests.
 */
public class ClusterTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec,
         TreeTabIDs entityViewNId, String primaryTabNId,
         String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), clusterSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.CLUSTER;
   }

   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec,
         String primaryTabNId, String secondaryTabNId) {
      this(clusterSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec,
         TreeTabIDs entityTabNId, String primaryTabNId, String secondaryTabNId) {
      this(clusterSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec, String primaryTabNId) {
      this(clusterSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId, null,
            null);
   }

   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec,
         TreeTabIDs entityTabNId, String primaryTabNId) {
      this(clusterSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec) {
      this(clusterSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, null, null, null);
   }

   /**
    * Build a location path based on the provided cluster navigation
    * identifiers.
    *
    * @param clusterSpec
    *           the cluster spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public ClusterTreeLocationSpec(ClusterSpec clusterSpec,
         TreeTabIDs entityTabNId) {
      this(clusterSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the cluster entity view.
    */
   public ClusterTreeLocationSpec() {
      super(TreeTabIDs.HOSTS_AND_CLUSTERS.getTreeTabID());
      _entityType = EntityTypes.CLUSTER;
   }
}