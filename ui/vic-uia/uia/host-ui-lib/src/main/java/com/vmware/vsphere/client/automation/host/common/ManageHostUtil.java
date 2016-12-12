/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.host.common;

import com.vmware.client.automation.util.ResourceUtil;


/**
 * Utilities for Manage host tests.
 */
public class ManageHostUtil {
   //---------------------------------------------------------------------------
   // Localization Utilities

   private static final String RESOURCE_NAME = "ManageHostTests";

   /**
    * Returns localized message
    *
    * @param key
    *       Resource key
    * @return
    *       localized message
    */
   public static String getLocalizedString(String key) {
      return ResourceUtil.getString(ManageHostUtil.class.getClassLoader(), RESOURCE_NAME, key);
   }
}
