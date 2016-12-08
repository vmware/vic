/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;

/**
 * A location spec suitable for modeling tree navigation in the Datacenter
 * related tests.
 */
public class DatacenterTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec,
         TreeTabIDs entityViewNId, String primaryTabNId,
         String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), datacenterSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.DATACENTER;
   }

   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec,
         String primaryTabNId, String secondaryTabNId) {
      this(datacenterSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec,
         TreeTabIDs entityTabNId, String primaryTabNId, String secondaryTabNId) {
      this(datacenterSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec,
         String primaryTabNId) {
      this(datacenterSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId, null,
            null);
   }

   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec,
         TreeTabIDs entityTabNId, String primaryTabNId) {
      this(datacenterSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec) {
      this(datacenterSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, null, null, null);
   }

   /**
    * Build a location path based on the provided Datacenter navigation
    * identifiers.
    *
    * @param datacenterSpec
    *           the datacenter spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public DatacenterTreeLocationSpec(DatacenterSpec datacenterSpec,
         TreeTabIDs entityTabNId) {
      this(datacenterSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the Datacenter entity view.
    */
   public DatacenterTreeLocationSpec() {
      super(TreeTabIDs.HOSTS_AND_CLUSTERS.getTreeTabID());
      _entityType = EntityTypes.DATACENTER;
   }
}