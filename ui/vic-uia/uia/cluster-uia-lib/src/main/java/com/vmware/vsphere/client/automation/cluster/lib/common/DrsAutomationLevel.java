/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.common;

/**
 * Represents all available cluster DRS automation level options.
 */
public enum DrsAutomationLevel {
    MANUAL(
            ClusterUtil.getLocalizedString(
                    "editCluster.dialog.drs.automationLevel.manual")),
    PARTIALLY_AUTOMATED(
            ClusterUtil.getLocalizedString(
                    "editCluster.dialog.drs.automationLevel.partiallyAutomated")),
    FULLY_AUTOMATED(
            ClusterUtil.getLocalizedString(
                    "editCluster.dialog.drs.automationLevel.fullyAutomated"));

    private String _value;

    private DrsAutomationLevel(String value) {
        this._value = value;
    }

    /**
     * Returns localised DRS automation level value.
     *
     * @return
     */
    public String getValue() {
        return _value;
    }
}
