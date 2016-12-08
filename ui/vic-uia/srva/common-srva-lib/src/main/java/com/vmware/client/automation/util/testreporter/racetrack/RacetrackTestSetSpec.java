/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.util.testreporter.racetrack;

import java.util.Properties;

import com.vmware.client.automation.util.testreporter.TestLoggerUtil;
import com.vmware.client.automation.util.testreporter.TestSetSpec;
import com.vmware.hsua.common.util.IOUtils;

/**
 * Test set specification for the Racetrack logging including application under
 * test info and environment info
 *
 */
public class RacetrackTestSetSpec implements TestSetSpec {

   private static final String DEFAULT_CONFIG_FILE = "racetrack.properties";

   private static final String BUILD_NUMBER_KEY = "racetrack.buildNumber";
   private static final String TEST_OWNER_KEY = "racetrack.username";
   private static final String PRODUCT_NAME_KEY = "racetrack.productName";
   // Product branch under test
   private static final String PRODUCT_BRANCH_KEY = "racetrack.branch";
   // Type of the test, i.e. regression, stress, load, etc.
   private static final String TEST_TYPE_KEY = "racetrack.testType";
   private static final String TEST_SET_DESCRIPTION_KEY = "racetrack.testSetDescription";
   // Type of the product build - release, beta, etc.
   private static final String BUILD_TYPE_KEY = "racetrack.buildType";
   private static final String BROWSER_OS_KEY = "racetrack.browserOs";
   private static final String OS_LOCALE_KEY = "racetrack.osLocale";
   private static final String RACETRACK_RESULTID_KEY = "racetrack.resultid";

   private final Properties config;

   public RacetrackTestSetSpec(String configurationFile) {
      if(configurationFile.isEmpty()){
			configurationFile = TestLoggerUtil.getCanonicalFileName(
					this.getClass(), DEFAULT_CONFIG_FILE);
      }

      config = IOUtils.readConfiguration(configurationFile);
   }

   @Override
   public String getBrowser() {
      // browser is not supported in racetrack reported
      // racetrack reporter only publishes results
      // it should not be coupled with selenium
      return null;
   }

   @Override
   public String getBrowserOs() {
      return config.getProperty(BROWSER_OS_KEY);
   }

   @Override
   public String getTestOwner() {
      return config.getProperty(TEST_OWNER_KEY);
   }

   @Override
   public String getProductName() {
      return config.getProperty(PRODUCT_NAME_KEY);
   }

   @Override
   public String getBuildNumber() {
      return config.getProperty(BUILD_NUMBER_KEY);
   }

   @Override
   public String getBuildType() {
      return config.getProperty(BUILD_TYPE_KEY);
   }

   @Override
   public String getTestSetDescription() {
      return config.getProperty(TEST_SET_DESCRIPTION_KEY);
   }

   @Override
   public String getBranch() {
      return config.getProperty(PRODUCT_BRANCH_KEY);
   }

   @Override
   public String getTestType() {
      return config.getProperty(TEST_TYPE_KEY);
   }

   @Override
   public String getLanguage() {
      return config.getProperty(OS_LOCALE_KEY);
   }

   @Override
   public String getExistingResultId() {
      return config.getProperty(RACETRACK_RESULTID_KEY);
   }
}
