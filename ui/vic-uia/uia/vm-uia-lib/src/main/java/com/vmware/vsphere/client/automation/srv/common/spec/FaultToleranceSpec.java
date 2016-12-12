/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.srv.common.util.FtUtil;

/**
 * Spec for Fault Tolerance
 *
 */
public class FaultToleranceSpec extends BaseSpec {

   /**
    * Enumeration for Fault Tolerance states
    */
   public enum FaultToleranceState {
      STARTING("portlet.state.starting"), SUSPENDED("portlet.state.suspended"), NEED_SECONDARY_VM(
            "portlet.state.needsecondaryvm"), VM_NOT_RUNNING(
            "portlet.state.vmnotrunning");

      private final String localizedDisplayName;

      private FaultToleranceState(String key) {
         localizedDisplayName = FtUtil.getLocalizedString(key);
      }

      public String getValue() {
         return localizedDisplayName;
      }
   };

   /**
    * Enumeration for Fault Tolerance statuses
    */
   public enum FaultToleranceStatus {
      PROTECTED("portlet.status.protected"), NOT_PROTECTED(
            "portlet.status.notprotected");

      private final String localizedDisplayName;

      private FaultToleranceStatus(String key) {
         localizedDisplayName = FtUtil.getLocalizedString(key);
      }

      public String getValue() {
         return localizedDisplayName;
      }
   };


   /**
    * Expected status
    */
   public DataProperty<FaultToleranceStatus> status;

   /**
    * Expected state
    */
   public DataProperty<FaultToleranceState> state;
}