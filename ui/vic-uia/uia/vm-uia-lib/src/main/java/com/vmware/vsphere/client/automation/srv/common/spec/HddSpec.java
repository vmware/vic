/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec used for Hard Disk
 */
public class HddSpec extends ManagedEntitySpec {
   /**
    * Property that shows HDD capacity/size
    */
   public DataProperty<String> hddCapacity;

   /**
    * Property that shows HDD capacity type MB/GB/TB
    */
   public DataProperty<String> hddCapacityType;

   /**
    * Operation on the device - edit or add
    */
   public enum HddActionType {
      EDIT("Edit"), ADD("Add");

      private String value;

      private HddActionType(String value) {
         this.value = value;
      }

      public String value() {
         return value;
      }
   }

   public DataProperty<HddActionType> vmDeviceAction;

   /**
    * Contains HDD node types
    */
   public static enum HddNodes {
      SCSI_CONTROLER_0, IDE0, IDE1
   }

   /**
    * Hdd device node types
    */
   public DataProperty<HddNodes> hddNodes;
}
