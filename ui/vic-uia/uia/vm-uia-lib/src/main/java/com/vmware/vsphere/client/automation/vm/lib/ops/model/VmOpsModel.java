/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.ops.model;

import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 *
 */
public class VmOpsModel {
   public enum VmPowerState {
      POWER_ON("vm.state.poweredOn"), POWER_OFF("vm.state.poweredOff");

      private String message;

      VmPowerState(String messageeKey) {
         this.message = CommonUtil.getLocalizedString(messageeKey);
      }

      public String getMessage() {
         return message;
      }
   }
}
