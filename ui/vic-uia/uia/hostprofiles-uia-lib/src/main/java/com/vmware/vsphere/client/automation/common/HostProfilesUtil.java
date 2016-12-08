/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common;

/**
 * Common utilities for Host Profiles tests.
 */
public class HostProfilesUtil {

    /**
     * Common method validating input is initialized
     * @param o - object to validate
     * @param msg - error message
     */
   public static void ensureNotNull(Object o, String msg) {
      if(o == null) {
         throw new RuntimeException(msg);
      }
   }

   /**
    * Removes leading and trailing characters from all the elements of the
    * array, if they exist.
    *
    * @param strings     the array of strings to trim
    * @param charsToTrim the characters to trim
    * @return the trimmed array
    */
   public static String[] trimLeadingTrailingChars(final String[] strings,
                                                   final String charsToTrim) {
      HostProfilesUtil.ensureNotNull(strings,
                                     "The array is null.");
      String[] result = new String[strings.length];
      for (int i = 0; i < strings.length; i++) {
         result[i] = trimLeadingTrailingChars(strings[i], charsToTrim);
      }
      return result;
   }

   /**
    * Removes leading and trailing characters, if they exist.
    *
    * @param string      the string to trim
    * @param charsToTrim the characters to trim
    * @return the resulting trimmed string
    */
   public static String trimLeadingTrailingChars(final String string,
                                                 final String charsToTrim) {
      HostProfilesUtil.ensureNotNull(string,
                                     "The string is null.");
      String result = string;

      if (string.startsWith(charsToTrim)) {
         result = result.substring(1, result.length());
      }

      if (string.endsWith(charsToTrim)) {
         result = result.substring(0, result.length() - 1);
      }

      return result;
   }
}
