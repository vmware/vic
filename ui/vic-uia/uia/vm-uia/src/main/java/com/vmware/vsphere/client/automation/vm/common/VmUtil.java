/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.common;

import com.vmware.client.automation.util.ResourceUtil;

/**
 * Utilities for VM tests.
 */
public class VmUtil {

    private static final String RESOURCE_NAME = "VmTests";

    /**
     * Returns localized message
     *
     * @param key
     *            Resource key
     * @return localized message
     */
    public static String getLocalizedString(String key) {
        return ResourceUtil.getString(VmUtil.class.getClassLoader(), RESOURCE_NAME, key);
    }

}
