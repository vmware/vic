/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the Cloud Resource Pool related tests.
 */
public class CrpLocationSpec extends NGCLocationSpec {
   /**
    * Build a location path based on the provided CRP navigation identifiers.
    */
   public CrpLocationSpec(
         String crpName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_RESOURCES,
            NGCNavigator.NID_RESOURCES_CRPS,
            crpName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided CRP navigation identifiers.
    */
   public CrpLocationSpec(
         String crpName, String primaryTabNId, String secondaryTabNId) {
      this(crpName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided CRP navigation identifiers.
    */
   public CrpLocationSpec(String crpName, String primaryTabNId) {
      this(crpName, primaryTabNId, null, null);
   }


   /**
    * Build a location path based on the provided CRP navigation identifiers.
    */
   public CrpLocationSpec(String crpName) {
      this(crpName, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the CRP entity view.
    */
   public CrpLocationSpec() {
      this(null, null, null, null);
   }
}
