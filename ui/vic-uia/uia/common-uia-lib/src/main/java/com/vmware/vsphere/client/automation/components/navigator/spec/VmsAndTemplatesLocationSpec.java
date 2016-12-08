/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the VmsAndTemplates related tests.
 */
public class VmsAndTemplatesLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided VmsAndTemplates navigation identifiers.
    */
   public VmsAndTemplatesLocationSpec(
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(null,
            NGCNavigator.NID_VCENTER_VMSANDTEMPLATES,
            "",
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided VmsAndTemplates navigation identifiers.
    */
   public VmsAndTemplatesLocationSpec(String primaryTabNId, String secondaryTabNId) {
      this(primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided VmsAndTemplates navigation identifiers.
    */
   public VmsAndTemplatesLocationSpec(String primaryTabNId) {
      this(primaryTabNId, null, null);
   }


   /**
    * Build a location path based on the provided VmsAndTemplates navigation identifiers.
    */
   public VmsAndTemplatesLocationSpec() {
      this(null, null, null);
   }

}
