/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;

/**
 * A location spec suitable for modelling a standard navigation in
 * the vApp related tests.
 */
public class VappLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    */
   public VappLocationSpec(
         String vappName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VAPP,
            vappName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    */
   public VappLocationSpec(
         VappSpec vappSpec,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VAPP,
            vappSpec,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    */
   public VappLocationSpec(
         String vappName, String primaryTabNId, String secondaryTabNId) {
      this(vappName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    */
   public VappLocationSpec(String vappName, String primaryTabNId) {
      this(vappName, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    */
   public VappLocationSpec(String vappName) {
      this(vappName, null, null, null);
   }

   /**
    * Build a location path based on the provided vApp navigation identifiers.
    */
   public VappLocationSpec(VappSpec vappSpec) {
      this(vappSpec, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the vApp entity view.
    */
   public VappLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_VAPP);
   }
}
