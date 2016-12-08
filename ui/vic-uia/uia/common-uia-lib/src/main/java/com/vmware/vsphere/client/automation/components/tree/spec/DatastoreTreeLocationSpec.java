/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.vmware.vsphere.client.automation.components.tree.EntityTypes;
import com.vmware.vsphere.client.automation.components.tree.TreeTabIDs;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;

/**
 * A location spec suitable for modeling tree navigation in the Datastore
 * related tests.
 */
public class DatastoreTreeLocationSpec extends TreeLocationSpec {
   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    * @param tocTabNid
    *           the toc tab NID
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec,
         TreeTabIDs entityViewNId, String primaryTabNId,
         String secondaryTabNId, String tocTabNid) {

      super(entityViewNId.getTreeTabID(), datastoreSpec, primaryTabNId,
            secondaryTabNId, tocTabNid);
      _entityType = EntityTypes.DATASTORE;
   }

   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec,
         String primaryTabNId, String secondaryTabNId) {
      this(datastoreSpec, TreeTabIDs.STORAGE, primaryTabNId, secondaryTabNId,
            null);
   }

   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    * @param secondaryTabNId
    *           the secondary tab NID
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec,
         TreeTabIDs entityTabNId, String primaryTabNId, String secondaryTabNId) {
      this(datastoreSpec, entityTabNId, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec,
         String primaryTabNId) {
      this(datastoreSpec, TreeTabIDs.STORAGE, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    * @param entityViewNId
    *           the tree tab NID
    * @param primaryTabNId
    *           the primary tab NID
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec,
         TreeTabIDs entityTabNId, String primaryTabNId) {
      this(datastoreSpec, entityTabNId, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec) {
      this(datastoreSpec, TreeTabIDs.STORAGE, null, null, null);
   }

   /**
    * Build a location path based on the provided Datastore navigation
    * identifiers.
    *
    * @param datastoreSpec
    *           the datastore spec
    * @param entityViewNId
    *           the tree tab NID
    */
   public DatastoreTreeLocationSpec(DatastoreSpec datastoreSpec,
         TreeTabIDs entityTabNId) {
      this(datastoreSpec, entityTabNId, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the Datastore entity view.
    */
   public DatastoreTreeLocationSpec() {
      super(TreeTabIDs.STORAGE.getTreeTabID());
      _entityType = EntityTypes.DATASTORE;
   }
}