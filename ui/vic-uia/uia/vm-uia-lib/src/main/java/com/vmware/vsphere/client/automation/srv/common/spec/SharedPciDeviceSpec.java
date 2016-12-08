/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec used for Shared PCI Device
 */
public class SharedPciDeviceSpec extends ManagedEntitySpec {
   /**
    * Operation on Shared PCI Device - edit or add
    */
   public enum SharedPCIDeviceActionType {
      EDIT("Edit"), ADD("Add");

      private String value;

      private SharedPCIDeviceActionType(String value) {
         this.value = value;
      }

      public String value() {
         return value;
      }
   }

   public DataProperty<SharedPCIDeviceActionType> vmDeviceAction;

   /**
    * Shared PCI Device device reserve all memory
    */
   public DataProperty<Boolean> SharedPCIDeviceReserveMemory;

   /**
    * VGPU Profiles for Shared PCI Device
    */
   public DataProperty<String> VgpuProfile;

}
