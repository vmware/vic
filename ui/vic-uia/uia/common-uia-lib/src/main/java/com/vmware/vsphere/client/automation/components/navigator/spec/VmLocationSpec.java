/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

/**
 * A location spec suitable for modelling a standard navigation in
 * the VM related tests.
 */
public class VmLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(
         String vmName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VM,
            vmName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(
         VmSpec vmSpec,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_VM,
            vmSpec,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(
         String vmName, String primaryTabNId, String secondaryTabNId) {
      this(vmName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(
         VmSpec vmSpec, String primaryTabNId, String secondaryTabNId) {
      this(vmSpec, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(String vmName, String primaryTabNId) {
      this(vmName, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(VmSpec vmSpec, String primaryTabNId) {
      this(vmSpec, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(String vmName) {
      this(vmName, null, null, null);
   }

   /**
    * Build a location path based on the provided VM navigation identifiers.
    */
   public VmLocationSpec(VmSpec vmSpec) {
      this(vmSpec, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the VM entity view.
    */
   public VmLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_VM);
   }
}
