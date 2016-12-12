/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Specifies the length of time to sleep the thread for.
 */
public class SleepSpec extends BaseSpec {

   /**
    * The number of milliseconds to sleep the thread for.
    */
   public DataProperty<Long> sleepTimeInMillis;
}
