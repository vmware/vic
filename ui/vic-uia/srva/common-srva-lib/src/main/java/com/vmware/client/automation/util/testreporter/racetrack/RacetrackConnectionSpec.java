/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.util.testreporter.racetrack;

import java.util.Properties;

import com.google.common.base.Strings;
import com.vmware.client.automation.util.testreporter.TestLoggerConnectSpec;
import com.vmware.client.automation.util.testreporter.TestLoggerUtil;
import com.vmware.hsua.common.util.IOUtils;

/**
 * Connection specification for the Racetrack connection.
 */
public class RacetrackConnectionSpec implements TestLoggerConnectSpec {

   private static final String RACETRACK_URL_KEY = "racetrack.url";
   private static final String THREADED_LOGGING_KEY = "racetrack.useThreadedlogging";
   private static final String DEFAULT_CONFIG_FILE = "racetrack.properties";

   private final Properties config;

   /**
    * Initializes racetrack connection spec and
    * loads racetrack test run configuration data like browser OS, locale, etc.
    *
    * @param configurationFile   the name of the property file holding the data
    */
   public RacetrackConnectionSpec(String configurationFile) {
      if (configurationFile.isEmpty()) {
			configurationFile = TestLoggerUtil.getCanonicalFileName(
					this.getClass(), DEFAULT_CONFIG_FILE);
      }

      config = IOUtils.readConfiguration(configurationFile);
   }

   @Override
   public String getTestLoggerURL() {
      return config.getProperty(RACETRACK_URL_KEY);
   }

   @Override
   public boolean getThreadedLogging() {
      return Boolean.parseBoolean(config.getProperty(THREADED_LOGGING_KEY));
   }
}
