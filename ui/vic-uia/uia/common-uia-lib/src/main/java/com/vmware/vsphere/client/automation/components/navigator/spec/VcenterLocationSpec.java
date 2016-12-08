/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * A <code>LocationSpec</code> suitable for modelling a standard navigation in
 * the vCenter server related tests.
 */
public class VcenterLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided vCenter navigation identifiers.
    */
   public VcenterLocationSpec(
         String vcenterName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VCS,
            vcenterName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided vCenter navigation identifiers.
    */
   public VcenterLocationSpec(
         VcSpec vcSpec,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VCS,
            vcSpec,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided vCenter server navigation identifiers.
    */
   public VcenterLocationSpec(
         String vcenterName, String primaryTabNId, String secondaryTabNId) {
      this(vcenterName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided vCenter server navigation identifiers.
    */
   public VcenterLocationSpec(String vcenterName, String primaryTabNId) {
      this(vcenterName, primaryTabNId, null, null);
   }


   /**
    * Build a location path based on the provided vCenter server navigation identifiers.
    */
   public VcenterLocationSpec(String vcenterName) {
      this(vcenterName, null, null, null);
   }

   /**
    * Build a location path based on the provided vCenter navigation identifiers.
    */
   public VcenterLocationSpec(VcSpec vcSpec) {
      this(vcSpec, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the vCenter servers entity view.
    */
   public VcenterLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VCS);
   }
}
