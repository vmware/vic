/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;

/**
 * A <code>LocationSpec</code> suitable for modelling a standard navigation in
 * the datastore cluster related tests.
 */
public class DatastoreClusterLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(String dsClusterName,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DATASTORE_CLUSTERS,
            dsClusterName, primaryTabNId, secondaryTabNId, tocTabNid);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(DatastoreClusterSpec dsCluster,
         String primaryTabNId, String secondaryTabNId, String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DATASTORE_CLUSTERS,
            dsCluster, primaryTabNId, secondaryTabNId, tocTabNid);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(String dsClusterName,
         String primaryTabNId, String secondaryTabNId) {
      this(dsClusterName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(DatastoreClusterSpec dsCluster,
         String primaryTabNId, String secondaryTabNId) {
      this(dsCluster, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(String dsClusterName,
         String primaryTabNId) {
      this(dsClusterName, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(DatastoreClusterSpec dsCluster,
         String primaryTabNId) {
      this(dsCluster, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(String dsClusterName) {
      this(dsClusterName, null, null, null);
   }

   /**
    * Build a location path based on the provided datastore cluster navigation
    * identifiers.
    */
   public DatastoreClusterLocationSpec(DatastoreClusterSpec dsCluster) {
      this(dsCluster, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the datastore cluster entity view.
    */
   public DatastoreClusterLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DATASTORE_CLUSTERS);
   }
}
