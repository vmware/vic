/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec for Boot VM options
 */
public class VmBootOptionsSpec extends BaseSpec {
   /**
    * Property that shows VM's boot firmware
    */
   public DataProperty<String> firmware;

   /**
    * Property that shows VM's Security Boot option
    */
   public DataProperty<Boolean> securityBoot;
}
