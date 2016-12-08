/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * This class is a location spec which represents navigation to HOME > Policies and
 * Profile > <host profile>.
 */
public class HostProfileLocationSpec extends NGCLocationSpec {

   /**
    * Navigation to HOME > Policies and Profiles.
    */
   public HostProfileLocationSpec() {
      this(null);
   }

   /**
    * Navigation to HOME > Policies and Profile > <host profile>.
    */
   public HostProfileLocationSpec(String hostPofileName) {
      super(
            NGCNavigator.NID_HOME_RULES_AND_PROFILES,
            NGCNavigator.NID_HOST_PROFILES,
            hostPofileName);
   }

   /**
    * Build a location path based on the provided host profile navigation
    * identifiers.
    */
   public HostProfileLocationSpec(String hostPofileName, String primaryTabNId,
         String secondaryTabNId) {
      super(NGCNavigator.NID_HOME_RULES_AND_PROFILES,
            NGCNavigator.NID_HOST_PROFILES, hostPofileName, primaryTabNId,
            secondaryTabNId);
   }

}

