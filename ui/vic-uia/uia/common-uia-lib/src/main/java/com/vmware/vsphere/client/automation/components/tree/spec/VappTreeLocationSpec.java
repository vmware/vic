/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;

/**
 * A location spec suitable for modeling tree navigation in the vApp related
 * tests.
 */
public class VappTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public VappTreeLocationSpec(VappSpec vappSpec, TreeTabIDs entityViewNId,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), vappSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.VAPP;
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public VappTreeLocationSpec(VappSpec vappSpec, String primaryTabNId,
         String secondaryTabNId) {
      this(vappSpec, TreeTabIDs.VMS_AND_TEMPLATES, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public VappTreeLocationSpec(VappSpec vappSpec, TreeTabIDs entityTabNId,
         String primaryTabNId, String secondaryTabNId) {
      this(vappSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public VappTreeLocationSpec(VappSpec vappSpec, String primaryTabNId) {
      this(vappSpec, TreeTabIDs.VMS_AND_TEMPLATES, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public VappTreeLocationSpec(VappSpec vappSpec, TreeTabIDs entityTabNId,
         String primaryTabNId) {
      this(vappSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    */
   public VappTreeLocationSpec(VappSpec vappSpec) {
      this(vappSpec, TreeTabIDs.VMS_AND_TEMPLATES, null, null, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    *
    * @param vappSpec
    *           the vApp spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public VappTreeLocationSpec(VappSpec vappSpec, TreeTabIDs entityTabNId) {
      this(vappSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the vApp entity view.
    */
   public VappTreeLocationSpec() {
      super(TreeTabIDs.VMS_AND_TEMPLATES.getTreeTabID());
      _entityType = EntityTypes.VAPP;
   }
}