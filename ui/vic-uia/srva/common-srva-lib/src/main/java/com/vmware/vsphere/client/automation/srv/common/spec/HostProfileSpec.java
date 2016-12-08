/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.hsua.common.datamodel.BasePBox;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Host profile spec used to describe a host profile.
 */
public class HostProfileSpec extends ManagedEntitySpec {

   /**
    * The host from which this profile has been exported.
    */
   public DataProperty<ManagedEntitySpec> referenceHost;

   /**
    * Holds the hosts whose compliance status is to be verified
    */
   public DataProperty<ManagedEntitySpec> complianceStatusHosts;

    /**
     * Holds the description of the host profile
     */
   public DataProperty<String> description;

   @Override
   public boolean equals(Object o) {
      if (o == null || !HostProfileSpec.class.isInstance(o)) {
         return false;
      }

      boolean result = true;
      HostProfileSpec temp = (HostProfileSpec) o;

      result &= areEqualDataProperties(temp.name, this.name);
      result &= areEqualDataProperties(temp.description, this.description);
      return result;
   }


   // private methods
   private boolean areEqualDataProperties(BasePBox.DataProperty<?> prop1,
                                                BasePBox.DataProperty<?> prop2) {
      boolean result = true;
      if (prop1 == null || prop2 == null) {
         return false;
      } else if (prop1 == null && prop1 == null) {
         return true;
      }

      boolean areBothAssignedSame = prop1.isAssigned() == prop2.isAssigned();
      if (!areBothAssignedSame) {
         _logger.info("Props not equally assigned");
         result = false;
      } else {
         if (prop1.isAssigned() && !(prop1.get().equals(prop2.get()))) {
            _logger.info("Props do not have equal values");
            result = false;
         }
      }

      return result;
   }
}
