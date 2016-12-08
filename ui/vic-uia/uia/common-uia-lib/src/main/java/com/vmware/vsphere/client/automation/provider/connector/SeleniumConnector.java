/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.connector;

import org.apache.commons.lang.NotImplementedException;

import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.SeleniumServiceSpec;
import com.vmware.client.automation.util.SeleniumConnection;

/**
 * Selenium client connector
 */
public class SeleniumConnector implements TestbedConnector {

   private static final String WEBDRIVER_PROTOCOL_NAME = "webdriver";

   // TODO: rename the TestTarget class to SeleniumConnection
   private final SeleniumConnection _seleniumConnection;

   public SeleniumConnector(SeleniumServiceSpec seleniumServiceSpec) {
      _seleniumConnection =
            new SeleniumConnection(seleniumServiceSpec.seleniumServer.get(),
                  seleniumServiceSpec.seleniumBrowser.get(),
                  seleniumServiceSpec.seleniumBrowserArgs.isAssigned() ?
                  seleniumServiceSpec.seleniumBrowserArgs.get() : "",
                  seleniumServiceSpec.screenShotFolder.get(),
                  seleniumServiceSpec.screenShotWebServer.get());
      _seleniumConnection.setUseWebDriver(
            WEBDRIVER_PROTOCOL_NAME.equalsIgnoreCase(seleniumServiceSpec.seleniumProtocol.get()) + "");
      _seleniumConnection.setSeleniumServerPort(seleniumServiceSpec.seleniumServerPort
            .get());
   }

   @Override
   public boolean isAlive() {
      throw new NotImplementedException();
   }

   @SuppressWarnings("unchecked")
   @Override
   public SeleniumConnection getConnection() {
      return _seleniumConnection;
   }

   @Override
   public void connect() {
      // Nothing to do - no need to keep active connection to the selenium.
      // The connection is established at starting the browser.
   }

   /**
    * Establish connection and start the web browser configured by the
    * SeleniumConnection object.
    */
   public void startBrowser() {
      _seleniumConnection.startUp();
   }

   @Override
   public void disconnect() {
      _seleniumConnection.cleanUp();
   }

}
