/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.servicespec;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec class defining a Selenium service settings configuration.
 */
public class SeleniumServiceSpec extends ServiceSpec {

   public DataProperty<String> seleniumServer;

   public DataProperty<String> seleniumServerPort;

   public DataProperty<String> seleniumProtocol;

   public DataProperty<String> seleniumBrowser;

   public DataProperty<String> seleniumBrowserArgs;

   public DataProperty<String> seleniumOS;

   // RC Legacy settings
   // TODO rkovachev: WebDriver should not use these settings.
   public DataProperty<String> screenShotFolder;

   public DataProperty<String> screenShotWebServer;
}
