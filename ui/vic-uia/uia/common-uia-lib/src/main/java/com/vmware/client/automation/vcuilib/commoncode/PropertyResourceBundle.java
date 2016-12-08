/**
 * Copyright 2011 VMWare, Inc. All rights reserved. -- VMWare Confidential
 */
package com.vmware.client.automation.vcuilib.commoncode;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

import org.apache.commons.io.IOUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * This class loads the property file with VCUI QE test constants and use
 * them for initialization of the constants defined in TestConstantsKey file
 *
 * NOTE: this class is a copy of the one from VCUI-QE-LIB
 */
public class PropertyResourceBundle {

   private static final Logger logger =
         LoggerFactory.getLogger(PropertyResourceBundle.class);

   // public access required for other teams to use
   public static final String VCUI_TEST_CONSTANTS_FILE = "vcuiqe.test.constants.file";

   private static String propertiesFilePath = "TestConstants.properties";
   // "C:\\install\\VC_SETUP\\TestConstants.properties";

   private static Properties testConstantProperty = null;

   /**
    * Method which will be used to get the constants from the properties file.
    * This is the replacement implementation of current TestConstants file.
    * Going forward, TestConstats.java will only have constants required by
    * framework and Structures (Enums). All other key-value constants should be
    * added to properties file. This implementation is easy to maintain than
    * single java file and also easy to implement internationalization using
    * i18N. Throwing assertion error as consumers of this method need not handle
    * assertion errors and framework handles any assertions errors (UITestUtil
    * catches these and does a throwsafely to only fail that particular test).
    * The location of test constants property file can be specified by java system property
    * "vcuiqe.test.constants.file", if such property is not specified the the test constants
    * will be loaded from default default test constants file =
    * "vcuiqe\main\VCUIQA-FLEX-UI\resources\TestConstants.properties properties"
    *
    *
    * @param constant ID
    * @return String value corresponding to the input ID
    * @exception AssertionError
    */
   public static final String getValue(String key) {
      if (testConstantProperty == null) {
         logger.info("Load properties file");
         testConstantProperty = new Properties();

         try {
            String externalConstantsFilePath = System.getProperty(VCUI_TEST_CONSTANTS_FILE);
            InputStream is;
            if (externalConstantsFilePath != null) {
               File externalFile = new File(externalConstantsFilePath);
               if (!externalFile.isFile()) {
                  logger.error("External constants file not found in the path specified: " + externalConstantsFilePath);
                  throw new AssertionError("External constants file not found in the path specified "
                        + externalConstantsFilePath);
               }
               propertiesFilePath = externalConstantsFilePath;
               is = new FileInputStream(externalFile);
            } else {
               is = PropertyResourceBundle.class.getClassLoader().getResourceAsStream(propertiesFilePath);
            }

            try {
               logger.info("Test constants property file is: " + propertiesFilePath);
               testConstantProperty.load(is);
            } finally {
               IOUtils.closeQuietly(is);
            }
         } catch (FileNotFoundException e) {
            logger.error("Constants file not found in the path specified: "
                  + propertiesFilePath);
            throw new AssertionError("Constants file not found in the path specified "
                  + propertiesFilePath);
         } catch (IOException e) {
            logger.error("Constants file not loaded, path= " + propertiesFilePath);
            throw new AssertionError("Constants file not loaded, path= "
                  + propertiesFilePath);
         }
      }
      String returnString;
      if (testConstantProperty.containsKey(key)) {
         returnString = testConstantProperty.getProperty(key);
      } else {
         logger.error("Property not found " + key);
         throw new AssertionError(
               "testConstantProperty file is missing entry for property " + key);
      }
      return returnString;
   }
}
