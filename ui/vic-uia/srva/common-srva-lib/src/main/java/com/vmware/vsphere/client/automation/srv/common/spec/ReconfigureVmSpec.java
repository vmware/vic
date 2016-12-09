/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Container spec class that holds VM configs for reconfigure VM operation.
 */
public class ReconfigureVmSpec extends BaseSpec {

   /**
    * Holds spec of the VM that will be reconfigured.
    * This VM should be present in the inventory.
    */
   public DataProperty<VmSpec> targetVm;

   /**
    * Represents all new VM ploperties that will be applied on
    * rarget VM.
    */
   public DataProperty<VmSpec> newVmConfigs;
}
