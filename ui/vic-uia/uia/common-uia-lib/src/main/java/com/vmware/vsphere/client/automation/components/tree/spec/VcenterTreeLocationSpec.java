/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * A location spec suitable for modeling tree navigation in the vCenter related
 * tests.
 */
public class VcenterTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec, TreeTabIDs entityViewNId,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), vcSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.FOLDER;
   }

   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec, String primaryTabNId,
         String secondaryTabNId) {
      this(vcSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec, TreeTabIDs entityTabNId,
         String primaryTabNId, String secondaryTabNId) {
      this(vcSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec, String primaryTabNId) {
      this(vcSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec, TreeTabIDs entityTabNId,
         String primaryTabNId) {
      this(vcSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec) {
      this(vcSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, null, null, null);
   }

   /**
    * Build a location path based on the provided vCenter navigation
    * identifiers.
    *
    * @param vcSpec
    *           the VC spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public VcenterTreeLocationSpec(VcSpec vcSpec, TreeTabIDs entityTabNId) {
      this(vcSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the vCenter entity view.
    */
   public VcenterTreeLocationSpec() {
      super(TreeTabIDs.HOSTS_AND_CLUSTERS.getTreeTabID());
      _entityType = EntityTypes.FOLDER;
   }
}