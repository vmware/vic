/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;

/**
 * A location spec suitable for modeling tree navigation in the Datastore
 * cluster related tests.
 */
public class DatastoreClusterTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec, TreeTabIDs entityViewNId,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), datastoreClusterSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.DS_CLUSTER;
   }

   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec, String primaryTabNId,
         String secondaryTabNId) {
      this(datastoreClusterSpec, TreeTabIDs.STORAGE, primaryTabNId,
            secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec, TreeTabIDs entityTabNId,
         String primaryTabNId, String secondaryTabNId) {
      this(datastoreClusterSpec, entityTabNId, primaryTabNId, secondaryTabNId,
            null);
   }

   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec, String primaryTabNId) {
      this(datastoreClusterSpec, TreeTabIDs.STORAGE, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec, TreeTabIDs entityTabNId,
         String primaryTabNId) {
      this(datastoreClusterSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec) {
      this(datastoreClusterSpec, TreeTabIDs.STORAGE, null, null, null);
   }

   /**
    * Build a location path based on the provided Datastore cluster navigation
    * identifiers.
    *
    * @param datastoreClusterSpec
    *           the datastore cluster spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public DatastoreClusterTreeLocationSpec(
         DatastoreClusterSpec datastoreClusterSpec, TreeTabIDs entityTabNId) {
      this(datastoreClusterSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the Datastore cluster entity
    * view.
    */
   public DatastoreClusterTreeLocationSpec() {
      super(TreeTabIDs.STORAGE.getTreeTabID());
      _entityType = EntityTypes.DS_CLUSTER;
   }
}
