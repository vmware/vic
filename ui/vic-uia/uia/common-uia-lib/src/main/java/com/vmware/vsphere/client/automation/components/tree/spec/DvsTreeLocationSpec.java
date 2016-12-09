/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.DvsSpec;

/**
 * A location spec suitable for modeling tree navigation in the Dvs related
 * tests.
 */
public class DvsTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec, TreeTabIDs entityViewNId,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), dvsSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.DVS;
   }

   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec, String primaryTabNId,
         String secondaryTabNId) {
      this(dvsSpec, TreeTabIDs.NETWORKING, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec, TreeTabIDs entityTabNId,
         String primaryTabNId, String secondaryTabNId) {
      this(dvsSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec, String primaryTabNId) {
      this(dvsSpec, TreeTabIDs.NETWORKING, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec, TreeTabIDs entityTabNId,
         String primaryTabNId) {
      this(dvsSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec) {
      this(dvsSpec, TreeTabIDs.NETWORKING, null, null, null);
   }

   /**
    * Build a location path based on the provided Dvs navigation identifiers.
    *
    * @param dvsSpec
    *           the dvs spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public DvsTreeLocationSpec(DvsSpec dvsSpec, TreeTabIDs entityTabNId) {
      this(dvsSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the Dvs entity view.
    */
   public DvsTreeLocationSpec() {
      super(TreeTabIDs.NETWORKING.getTreeTabID());
      _entityType = EntityTypes.DVS;
   }
}