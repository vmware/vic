/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common;

import com.vmware.client.automation.util.ResourceUtil;

/**
 * Common localization utility class.
 */
public class CommonUtil {

   private static final String RESOURCE_NAME = "CommonTests";

   /**
    * Returns localized message.
    *
    * @param key     resource key
    * @return        localized message
    */
   public static String getLocalizedString(String key) {
      return ResourceUtil.getString(CommonUtil.class.getClassLoader(), RESOURCE_NAME, key);
   }
}
