/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;

/**
 * Describes expected host state and timeout for entering it
 */
public class HostStateSpec extends BaseSpec {

    /**
     * Property that describes the expected state of the host
     */
    public DataProperty<HostUtil.HostStates> state;

    /**
     * Property that describes the timeout to wait for entering the state
     */
    public DataProperty<Integer> numOfRetries;

}
