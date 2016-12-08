/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common;

import com.vmware.client.automation.util.ResourceUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.SdrsBehavior;

/**
 * Utilities for datastore cluster tests.
 */
public class DsClusterUtil {

    private static final String RESOURCE_NAME = "DsClusterTests";

    /**
     * Returns localized message
     *
     * @param key
     *            Resource key
     * @return localized message
     */
    public static String getLocalizedString(String key) {
        return ResourceUtil.getString(DsClusterUtil.class.getClassLoader(), RESOURCE_NAME, key);
    }

    /**
     * Gets expected UI Sdrs label based on SdrsBehavior enum value passed
     *
     * @param sdrsBehavior
     * @return
     */
    public static String getSdrsLabel(SdrsBehavior sdrsBehavior) {
        String sdrsLabel;
        switch (sdrsBehavior) {
        case MANUAL:
            sdrsLabel = DsClusterUtil.getLocalizedString("datastoreCluster.summary.sdrsAutomation.manual");
            break;
        case FULLY_AUTOMATED:
            sdrsLabel = DsClusterUtil.getLocalizedString("datastoreCluster.summary.sdrsAutomation.fullyAutomated");
            break;
        default:
            throw new IllegalArgumentException("Invalid SdrsBehavior type passed!");
        }
        return sdrsLabel;
    }
}
