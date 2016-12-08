/* Copyright $today.year VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Policy option spec used to represent the new values for the policy, as
 * well as the property to modify
 */
public class HostPolicyUpdateSpec extends EntitySpec {
    public DataProperty<String> name;

    public DataProperty<PolicyOptionType> originalPolicyType;
    public DataProperty<String> originalPropertyName;
    public DataProperty<String> originalPolicyValue;

    public DataProperty<PolicyOptionType> newPolicyType;
    public DataProperty<String> newPropertyName;
    public DataProperty<String> newPropertyValue;
}
