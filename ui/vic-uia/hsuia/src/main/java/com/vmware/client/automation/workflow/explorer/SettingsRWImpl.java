/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.util.Properties;

/**
 * Provide ability to read and write settings.
 */
public class SettingsRWImpl extends Properties implements SettingsReader, SettingsWriter {

   /**
    * 
    */
   private static final long serialVersionUID = -7216557737421612184L;

   @Override
   public void setSetting(String key, String value) {
      this.setProperty(key, value);
   }

   @Override
   public String getSetting(String key) {
      return this.getProperty(key);
   }

}
