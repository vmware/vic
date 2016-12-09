/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the datastore related tests.
 */
public class DatastoreLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided datastore navigation identifiers.
    */
   public DatastoreLocationSpec(
         String datastoreName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_DATASTORES,
            datastoreName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided datastore navigation
    * identifiers.
    */
   public DatastoreLocationSpec(DatastoreSpec datastoreSpec,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DATASTORES,
            datastoreSpec, primaryTabNId, secondaryTabNId, tocTabNid);
   }

   /**
    * Build a location path based on the provided datastore navigation identifiers.
    */
   public DatastoreLocationSpec(
         String datastoreName, String primaryTabNId, String secondaryTabNId) {
      this(datastoreName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided datastore navigation identifiers.
    */
   public DatastoreLocationSpec(String datastoreName, String primaryTabNId) {
      this(datastoreName, primaryTabNId, null, null);
   }


   /**
    * Build a location path based on the provided datastore navigation identifiers.
    */
   public DatastoreLocationSpec(String datastoreName) {
      this(datastoreName, null, null, null);
   }

   /**
    * Build a location path based on the provided datastore navigation
    * identifiers.
    */
   public DatastoreLocationSpec(DatastoreSpec datastoreSpec) {
      this(datastoreSpec, null, null, null);
   }

   /**
    * Build a location path based on the provided datastore navigation
    * identifiers.
    */
   public DatastoreLocationSpec(DatastoreSpec datastoreSpec,
         String primaryTabNId) {
      this(datastoreSpec, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided datastore navigation
    * identifiers.
    */
   public DatastoreLocationSpec(DatastoreSpec datastoreSpec,
         String primaryTabNId, String secondaryTabNId) {
      this(datastoreSpec, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location that will navigate the UI to the cluster entity view.
    */
   public DatastoreLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DATASTORES);
   }
}
