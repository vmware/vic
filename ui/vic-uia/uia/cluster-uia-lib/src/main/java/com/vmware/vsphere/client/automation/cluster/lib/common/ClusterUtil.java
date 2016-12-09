/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.common;

import com.vmware.client.automation.util.ResourceUtil;

/**
 * Utilities for cluster tests.
 */
public class ClusterUtil {

   //---------------------------------------------------------------------------
   // Localization Utilities

   private static final String RESOURCE_NAME = "ClusterTests";

   /**
    * Returns localized message
    *
    * @param key
    *       Resource key
    * @return
    *       localized message
    */
   public static String getLocalizedString(String key) {
      return ResourceUtil.getString(ClusterUtil.class.getClassLoader(), RESOURCE_NAME, key);
   }
}
