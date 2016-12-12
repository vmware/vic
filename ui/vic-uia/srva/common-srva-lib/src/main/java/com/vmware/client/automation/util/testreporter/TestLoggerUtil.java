/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.util.testreporter;

import java.io.File;
import java.net.URL;

/**
 * Utility class for Test Logger common methods.
 */
public class TestLoggerUtil {

   /**
    * Finds file stored in class loader bootstrap directory.
    */
   public static String getCanonicalFileName(Class<?> resourceClass,
         String fileName) {
      try {
         URL url = resourceClass.getClassLoader()
               .getResource(fileName);
         File f = new File(url.toURI());

         return f.getCanonicalPath();
      } catch (Exception e) {
         throw new RuntimeException("Cannot locate file: " + fileName, e);
      }
   }
}
