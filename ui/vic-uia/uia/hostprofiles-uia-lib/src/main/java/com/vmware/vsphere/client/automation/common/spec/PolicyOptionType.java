/*
 *  Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.vsphere.client.automation.common.spec;

/**
 * Enum representing the different policy option types in Edit Host Profile
 * wizard > step 2 Edit host Profile > policy tree > Advanced Configuration
 * Settings > Advanced Options
 */
public enum PolicyOptionType {
    FIXED_OPTION("FixedConfigOption"),
    DEFAULT_VALUE("SetDefaultConfigOption"),
    USER_EXPLICIT_POLICY_CHOICE("NoDefaultOption"),
    SPECIFIED_IN_HOST_CUSTOMIZATIONS("UserInputAdvancedConfigOption") ;

    private String optionTypeValue;

    public String getValue() {
        return optionTypeValue;
    }

    PolicyOptionType(String optionTypeValue) {
        this.optionTypeValue = optionTypeValue;
    }
}
