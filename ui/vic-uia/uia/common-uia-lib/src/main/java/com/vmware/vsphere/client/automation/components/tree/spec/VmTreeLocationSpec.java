/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

/**
 * A location spec suitable for modeling tree navigation in the VM related
 * tests.
 */
public class VmTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    * @param vmSpec
    *           the VM spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public VmTreeLocationSpec(VmSpec vmSpec, TreeTabIDs entityViewNId,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), vmSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.VM;
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    * @param vmSpec
    *           the VM spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public VmTreeLocationSpec(VmSpec vmSpec, String primaryTabNId,
         String secondaryTabNId) {
      this(vmSpec, TreeTabIDs.VMS_AND_TEMPLATES, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    *@param vmSpec
    *           the VM spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public VmTreeLocationSpec(VmSpec vmSpec, TreeTabIDs entityTabNId,
         String primaryTabNId, String secondaryTabNId) {
      this(vmSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    * @param vmSpec
    *           the VM spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public VmTreeLocationSpec(VmSpec vmSpec, String primaryTabNId) {
      this(vmSpec, TreeTabIDs.VMS_AND_TEMPLATES, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    * @param vmSpec
    *           the VM spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public VmTreeLocationSpec(VmSpec vmSpec, TreeTabIDs entityTabNId,
         String primaryTabNId) {
      this(vmSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    * @param vmSpec
    *           the VM spec
    */
   public VmTreeLocationSpec(VmSpec vmSpec) {
      this(vmSpec, TreeTabIDs.VMS_AND_TEMPLATES, null, null, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    *
    * @param vmSpec
    *           the VM spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public VmTreeLocationSpec(VmSpec vmSpec, TreeTabIDs entityTabNId) {
      this(vmSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the VM entity view.
    */
   public VmTreeLocationSpec() {
      super(TreeTabIDs.VMS_AND_TEMPLATES.getTreeTabID());
      _entityType = EntityTypes.VM;
   }
}