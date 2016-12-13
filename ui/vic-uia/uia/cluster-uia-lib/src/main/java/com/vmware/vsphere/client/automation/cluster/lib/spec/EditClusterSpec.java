/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec.DrsAutoLevel;

/**
 * Container for test-specific data for the edit-cluster test.
 */
public class EditClusterSpec extends BaseSpec {

    /**
     * Specifies if the test should turn on or off vSphere DRS.
     */
    public DataProperty<Boolean> drsEnabled;

    /**
     * Specifies the automation level of DRS to be set by the test.
     * This setting depends on {@link EditClusterSpec#drsEnabled}}
     * which has to be set to true.
     */
    public DataProperty<String> drsAutomationLevel;

    /**
     * Specifies if individual VM DRS automation level is enabled or disabled.
     * This setting depends on {@link EditClusterSpec#drsEnabled}}
     * which has to be set to true.
     */
    public DataProperty<Boolean> vmDrsAutomationLevelEnabled;

    /**
     * Specifies the enum for the DRS automation level
     */
    public DataProperty<DrsAutoLevel> autoLevel;

    /**
     * Specifies if advanced Enforce Even Distribution option is enabled or
     * disabled.
     */
    public DataProperty<Boolean> advancedEnforceEvenDistributionEnabled;

    /**
     * Specifies if advanced Consumed Memory option is enabled or disabled.
     */
    public DataProperty<Boolean> advancedConsumedMemoryEnabled;

    /**
     * Specifies if advanced CPU over-commitment option is enabled or disabled.
     */
    public DataProperty<Boolean> advancedCPUOverCommitmentEnabled;

    /**
     * Specifies advanced CPU over-commitment option's value
     * This setting depends on {@link EditClusterSpec#advancedCPUOverCommitmentEnabled}}
     * which has to be set to true.
     */
    public DataProperty<String> advancedCPUOverCommitmentValue;
}
