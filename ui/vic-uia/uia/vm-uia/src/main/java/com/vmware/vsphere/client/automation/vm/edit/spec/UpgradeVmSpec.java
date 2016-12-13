package com.vmware.vsphere.client.automation.vm.edit.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec for upgrade VM operation
 */
public class UpgradeVmSpec extends BaseSpec {
   /**
    * Compatibility version. E.g 'ESXi 6.5 and later'
    */
   public DataProperty<String> compatibilityVersion;

   /**
    * VM hardware version. E.g 'VM version 13'
    */
   public DataProperty<String> vmHardwareVersion;

   /**
    * Indicates if the new update should be scheduled
    */
   public DataProperty<Boolean> scheduleUpdate;
}
