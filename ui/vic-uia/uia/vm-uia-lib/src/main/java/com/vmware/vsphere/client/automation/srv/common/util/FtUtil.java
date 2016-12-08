/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.util;

import com.vmware.client.automation.util.ResourceUtil;

/**
 * Utilities for Fault Tolerance tests.
 */
public class FtUtil {
   private static final String RESOURCE_NAME = "FtTests";

   /**
    * Returns localized message
    *
    * @param key
    *           Resource key
    * @return localized message
    */
   public static String getLocalizedString(String key) {
      return ResourceUtil.getString(FtUtil.class.getClassLoader(),
            RESOURCE_NAME, key);
   }
}