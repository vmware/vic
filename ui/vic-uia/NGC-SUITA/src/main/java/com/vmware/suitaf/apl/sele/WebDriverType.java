/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.suitaf.apl.sele;

/**
 * This enum contains the different types of WebDrivers
 * Firefox, Chrome, Internet Explorer and their remote versions.
 */
public enum WebDriverType {
   FF, // Firefox run on the local system
   RFF, // Firefox run remotely using remote WebDriver
   GH, // Google Chrome run on the local system
   RGH, // Google Chrome run remotely using remote WebDriver
   IE, // Internet Explorer run on the local system
   RIE, // Internet Explorer run remotely using remote WebDriver
   SF, // Safari run on the local system
   RSF; // Safari run remotely using remote WebDriver

   private static String RC_IE = "*iexplore";
   private static String RC_FF = "*firefox";
   private static String RC_SF = "*safari";
   private static String RC_GCH = "*googlechrome";

   /**
    * Converts Selenium RC browser types string to WebDriverType enum.
    * @param seleniumRCBrowserType
    * @return respective WebDriverType corresponding to the given property.
    * TODO: Add converter for the local running browser.
    */
   public static WebDriverType getWebDriverType(String seleniumRCBrowserType) {

      if(seleniumRCBrowserType.equalsIgnoreCase(RC_IE)) {
         return WebDriverType.RIE;
      } else if (seleniumRCBrowserType.equalsIgnoreCase(RC_FF)) {
         return WebDriverType.RFF;
      } else if (seleniumRCBrowserType.equalsIgnoreCase(RC_SF)) {
         return WebDriverType.RSF;
      } else if (seleniumRCBrowserType.equalsIgnoreCase(RC_GCH)) {
         return WebDriverType.RGH;
      } else {
         throw new RuntimeException(
               "The specified browser type is not supported by the Selenium Web Drvier");
      }
   }
}
