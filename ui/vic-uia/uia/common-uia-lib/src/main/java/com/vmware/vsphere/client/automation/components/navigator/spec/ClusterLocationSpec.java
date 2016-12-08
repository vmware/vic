/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the cluster related tests.
 */
public class ClusterLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided cluster navigation identifiers.
    */
   public ClusterLocationSpec(
         String clusterName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_CLUSTERS,
            clusterName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided cluster navigation identifier
    * specified by the clusterSpec.
    * @param clusterSpec
    * @param primaryTabNId
    * @param secondaryTabNId
    * @param tocTabNid
    */
   public ClusterLocationSpec(
         ClusterSpec clusterSpec,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_CLUSTERS,
            clusterSpec,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }


   /**
    * Build a location path based on the provided cluster navigation identifiers.
    */
   public ClusterLocationSpec(
         String clusterName, String primaryTabNId, String secondaryTabNId) {
      this(clusterName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided cluster navigation identifier
    * specified by the clusterSpec.
    * @param clusterSpec
    * @param primaryTabNId
    * @param secondaryTabNId
    */
   public ClusterLocationSpec(
         ClusterSpec clusterSpec, String primaryTabNId, String secondaryTabNId) {
      this(clusterSpec, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided cluster navigation identifiers.
    */
   public ClusterLocationSpec(String clusterName, String primaryTabNId) {
      this(clusterName, primaryTabNId, null, null);
   }


   /**
    * Build a location path based on the provided cluster navigation identifiers.
    */
   public ClusterLocationSpec(String clusterName) {
      this(clusterName, null, null, null);
   }

   /**
    * Build a location path based on the provided cluster navigation identifiers.
    */
   public ClusterLocationSpec(ClusterSpec clusterSpec) {
      this(clusterSpec, null, null, null);
   }

   /**
    * Build a location path based on the provided cluster navigation identifiers.
    */
   public ClusterLocationSpec(ClusterSpec clusterSpec, String primaryTabNId) {
      this(clusterSpec, primaryTabNId, null, null);
   }

   /**
    * Build a location that will navigate the UI to the cluster entity view.
    */
   public ClusterLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_CLUSTERS);
   }

}
