/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for Network Resource Allocation settings of Virtual Machine NIC.
 */
public class NetworkResourceAllocationSpec extends ManagedEntitySpec {

   /**
    * An enumeration of the predefined shares levels in resource allocation settings.
    */
   public enum SharesLevel {
      LOW("Low"), NORMAL("Normal"), HIGH("High"), CUSTOM("Custom");

      private String value;

      private SharesLevel(String value) {
         this.value = value;
      }

      public String value() {
         return value;
      }
   }

   public enum DataRateUnit { MBPS, GBPS }

   /**
    * The shares to configured for this network adapter or traffic type.
    */
   public DataProperty<SharesLevel> sharesLevel;

   /**
    * The amount of shares for this adapter or traffic type. The {@code sharesValue} is
    * mandatory only when {@code sharesLevel} is set to {@code SharesLevel.CUSTOM}
    */
   public DataProperty<Integer> sharesValue;

   /**
    * The amount of bandwidth to set as reservation. By default it is in the same unit
    * as the one appearing in the UI (typically Mbit/s) unless another unit is
    * configured in {@code reservationUnit}.
    */
   public DataProperty<Integer> reservationValue;

   /**
    * The data rate unit for the {@code reservationValue}
    */
   public DataProperty<DataRateUnit> reservationUnit;

   /**
    * Maximum allowed network bandwidth usage. By default it is in the same unit
    * as the one appearing in the UI (typically Mbit/s) unless another unit is
    * configured in {@code limitUnit}. A value of -1 is treated as unlimited.
    */
   public DataProperty<Integer> limitValue;

   /**
    * The data rate unit for the {@code limitValue}
    */
   public DataProperty<DataRateUnit> limitUnit;

}
