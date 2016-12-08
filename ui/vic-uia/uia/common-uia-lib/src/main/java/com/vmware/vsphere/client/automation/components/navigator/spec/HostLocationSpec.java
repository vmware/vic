/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the host related tests.
 */
public class HostLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(
         String hostName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_HOSTS,
            hostName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(
         HostSpec hostSpec,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_HOSTS,
            hostSpec,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(
         String hostName, String primaryTabNId, String secondaryTabNId) {
      this(hostName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(
         HostSpec hostSpec, String primaryTabNId, String secondaryTabNId) {
      this(hostSpec, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(String hostName, String primaryTabNId) {
      this(hostName, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(HostSpec hostSpec, String primaryTabNId) {
      this(hostSpec, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(String hostName) {
      this(hostName, null, null, null);
   }

   /**
    * Build a location path based on the provided host navigation identifiers.
    */
   public HostLocationSpec(HostSpec hostSpec) {
      this(hostSpec, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the cluster entity view.
    */
   public HostLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_HOSTS);
   }
}
