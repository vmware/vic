/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;

/**
 * A location spec suitable for modeling a standard navigation in resource pool
 * related tests.
 */
public class ResourcePoolLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided resource pool navigation
    * identifiers.
    */
   public ResourcePoolLocationSpec(String resPoolName, String primaryTabNId,
         String secondaryTabNId, String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_RPS,
            resPoolName, primaryTabNId, secondaryTabNId, tocTabNid);
   }

   /**
    * Build a location path based on the provided resource pool navigation
    * identifiers.
    */
   public ResourcePoolLocationSpec(ResourcePoolSpec resPoolSpec, String primaryTabNId,
         String secondaryTabNId, String tocTabNid) {

      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_RPS,
            resPoolSpec, primaryTabNId, secondaryTabNId, tocTabNid);
   }

   /**
    * Build a location path based on the provided resource pool navigation
    * identifiers.
    */
   public ResourcePoolLocationSpec(String resPoolName, String primaryTabNId,
         String secondaryTabNId) {
      this(resPoolName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided resource pool navigation
    * identifiers.
    */
   public ResourcePoolLocationSpec(String resPoolName, String primaryTabNId) {
      this(resPoolName, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided resource pool navigation
    * identifiers.
    */
   public ResourcePoolLocationSpec(String resPoolName) {
      this(resPoolName, null, null, null);
   }

   /**
    * Build a location path based on the provided resource pool navigation
    * identifiers.
    */
   public ResourcePoolLocationSpec(ResourcePoolSpec resPoolSpec) {
      this(resPoolSpec, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the resource pool entity
    * view.
    */
   public ResourcePoolLocationSpec() {
      super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_RPS);
   }
}
