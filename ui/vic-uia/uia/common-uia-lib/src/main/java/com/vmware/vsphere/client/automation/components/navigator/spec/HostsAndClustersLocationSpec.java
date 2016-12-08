/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the HostsAndClusters related tests.
 */
public class HostsAndClustersLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided HostsAndClusters navigation identifiers.
    */
   public HostsAndClustersLocationSpec(
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_HOSTSANDCLUSTERS,
            "",
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided HostsAndClusters navigation identifiers.
    */
   public HostsAndClustersLocationSpec(String primaryTabNId, String secondaryTabNId) {
      this(primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided HostsAndClusters navigation identifiers.
    */
   public HostsAndClustersLocationSpec(String primaryTabNId) {
      this(primaryTabNId, null, null);
   }


   /**
    * Build a location path based on the provided HostsAndClusters navigation identifiers.
    */
   public HostsAndClustersLocationSpec() {
      this(null, null, null);
   }

}
