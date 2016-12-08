/*
 *  Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * A spec representing policy updates for a specific host in an answer file.
 */
public class AnswerFileUpdateSpec extends BaseSpec {
   /**
    * The host which should be updated
    */
   public DataProperty<HostSpec> host;

   /**
    * The policy or policies which should be updated for the host
    */
   public DataProperty<HostPolicyUpdateSpec> policyUpdateSpec;
}
