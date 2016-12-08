/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * A location spec suitable for modeling tree navigation in the Host related
 * tests.
 */
public class HostTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public HostTreeLocationSpec(HostSpec hostSpec, TreeTabIDs entityViewNId,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), hostSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.HOST;
   }

   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public HostTreeLocationSpec(HostSpec hostSpec, String primaryTabNId,
         String secondaryTabNId) {
      this(hostSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public HostTreeLocationSpec(HostSpec hostSpec, TreeTabIDs entityTabNId,
         String primaryTabNId, String secondaryTabNId) {
      this(hostSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public HostTreeLocationSpec(HostSpec hostSpec, String primaryTabNId) {
      this(hostSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public HostTreeLocationSpec(HostSpec hostSpec, TreeTabIDs entityTabNId,
         String primaryTabNId) {
      this(hostSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    */
   public HostTreeLocationSpec(HostSpec hostSpec) {
      this(hostSpec, TreeTabIDs.HOSTS_AND_CLUSTERS, null, null, null);
   }

   /**
    * Build a location path based on the provided Host navigation identifiers.
    *
    * @param hostSpec
    *           the host spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public HostTreeLocationSpec(HostSpec hostSpec, TreeTabIDs entityTabNId) {
      this(hostSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the Host entity view.
    */
   public HostTreeLocationSpec() {
      super(TreeTabIDs.HOSTS_AND_CLUSTERS.getTreeTabID());
      _entityType = EntityTypes.HOST;
   }
}