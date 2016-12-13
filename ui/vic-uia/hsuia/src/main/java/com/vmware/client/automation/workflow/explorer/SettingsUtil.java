/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import com.google.common.base.Strings;

/**
 * Util class for work with properties
 *
 */
public final class SettingsUtil {
   private SettingsUtil() {
   }

   public static String getRequiredValue(SettingsReader settingsReader, String key) {
      String value = settingsReader.getSetting(key);

      if (Strings.isNullOrEmpty(value)) {
         throw new IllegalArgumentException(
               String.format("Required %s key not set.",key));
      }

      return value;
   }

   public static boolean getBooleanValue(SettingsReader settingsReader, String key) {
      String value = settingsReader.getSetting(key);
      if (Strings.isNullOrEmpty(value)) {
         return false;
      }

      value = value.toLowerCase();
      return "true".equals(value);
   }
}
