package com.vmware.client.automation.util;


/**
 * Server Localization Utility Class.
 */
public class SrvLocalizationUtil {

   private static final String RESOURCE_NAME = "SrvTests";

   /**
    * Returns localized message.
    *
    * @param key  resource key
    * @return     localized message
    */
   public static String getLocalizedString(String key) {
      return ResourceUtil.getString(SrvLocalizationUtil.class.getClassLoader(), RESOURCE_NAME, key);
   }

   /**
    * Returns localized message
    *
    * @param key     Resource key
    * @param params  Substituted parameters for localized message
    * @return        localized message
    */
   public static String getLocalizedString(String key, String... params) {
      return ResourceUtil.getString(SrvLocalizationUtil.class.getClassLoader(), RESOURCE_NAME, key, params);
   }
}
