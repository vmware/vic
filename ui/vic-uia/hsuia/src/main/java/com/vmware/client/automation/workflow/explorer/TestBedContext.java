/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.util.Properties;

/**
 * Keep the data defining a testbed context - path to
 * configuration file and settings and providing method to access
 * settings.
 */
public class TestBedContext {

   private final String _testbedSettingsFilePath;
   private final Properties _settings;

   public TestBedContext(String testbedSettingsFilePath, Properties settings) {
      _testbedSettingsFilePath = testbedSettingsFilePath;
      _settings = settings;
   }

   /**
    * @return path to the setting file defining testbed setting.
    */
   public String getTestbedSettingsFilePath() {
      return _testbedSettingsFilePath;
   }

   /**
    * Return setting by key value.
    * @param key settings key
    * @return setting value
    */
   public String getTestbedSetting(String key) {
      return _settings.getProperty(key);
   }

   /**
    * Return true if the properties from testbedSettings are
    * the same.
    * @param testbedSettings
    * @return
    */
   public boolean testbedEquals(Properties testbedSettings) {
      return this._settings.equals(testbedSettings);
   }
}
