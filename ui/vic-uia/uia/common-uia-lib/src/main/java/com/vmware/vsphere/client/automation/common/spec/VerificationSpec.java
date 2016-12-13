/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * A spec for verifying an expected result.
 * It is used by steps that need to be parameterized in order to do the verification logs.
 * e.g. positive/negative verifications.
 */
public class VerificationSpec extends BaseSpec {
   /**
    * If a negative result is expected
    */
   public DataProperty<Boolean> isNegativeResultExpected;
}