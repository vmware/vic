/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util;

import java.util.HashMap;
import java.util.Map;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.sele.InitParams;
import com.vmware.suitaf.apl.sele.SeleAPLImpl;

/**
 * The class represents Selenium client connection.
 * It needs Selenium server and browser configuration.
 */
public class SeleniumConnection {

   private String seleniumServer = "127.0.0.1";
   private String seleniumServerPort = "4444";
   private String seleniumBrowser = "*iexplore";
   private String seleniumBrowserArgs = "";
   private String screenShotFolder = "C:\\imgs";
   private String screenShotWebServer = "";
   private String useWebDriver = "false";

   public SeleniumConnection() {
      // Nothing to do
   }

   /**
    * Set to use the WebDriver Selenium.
    * @param useWebDriver if not true use the RC Selenium.
    */
   public void setUseWebDriver(String useWebDriver) {
      this.useWebDriver = useWebDriver;
   }

   /**
    * Set Selenium server port. If not set use the default value - 4444.
    * @param seleniumServerPort  port of the selenium server.
    */
   public void setSeleniumServerPort(String seleniumServerPort) {
      if(seleniumServerPort.isEmpty()) {
         return;
      }
      this.seleniumServerPort = seleniumServerPort;
   }

   /**
    * If the parameters are empty use the default values. That will be the case
    * in development state.
    *
    * @param seleniumServer
    *           ip of the machine where RC server is running.
    * @param seleniumBrowser
    *           browser that will be used in the test.
    * @param screenShotFolder
    *           folder where the screenshot will be stored.
    */
   public SeleniumConnection(String seleniumServer, String seleniumBrowser, String seleniumBrowserArgs,
      String screenShotFolder, String screenShotWebServer) {
      if (!seleniumServer.isEmpty()) {
         this.seleniumServer = seleniumServer;
      }

      if (!screenShotFolder.isEmpty()) {
         this.screenShotFolder = screenShotFolder;
      }

      if (!seleniumBrowser.isEmpty()) {
         this.seleniumBrowser = seleniumBrowser;
      }

      if (!seleniumBrowserArgs.isEmpty()) {
         this.seleniumBrowserArgs = seleniumBrowserArgs;
      }

      this.screenShotWebServer = screenShotWebServer;
   }

   /**
    * This utility method configures and APL instance and prepares it to work
    * with particular test-target host. Then it initiates a test session on the
    * Selenium-RC server and requests opening of the browser.
    */
   // TODO: Externalize setting to configuration file.
   public void startUp() {
      // Set-up the environment needed by the SUITA framework
      SUITA.Environment.set(false, screenShotFolder, (long) 30000, screenShotWebServer);

      // Retrieves the APL implementation class name
      String aplClassName = SeleAPLImpl.class.getCanonicalName();

      // Retrieve the configuration parameters for the APL implementation
      Map<String, String> aplPar = new HashMap<String, String>();
      // The IP of the Selenium-RC server
      aplPar.put(InitParams.PN_SERVER_IP, seleniumServer);
      // The Port of the Selenium-RC server
      aplPar.put(InitParams.PN_SERVER_PORT, seleniumServerPort);
      // The browser selection key-word
      aplPar.put(InitParams.PN_BROWSER, seleniumBrowser);
      aplPar.put(InitParams.PN_BROWSER_ARGS, seleniumBrowserArgs);
      // The initial url for the opening browser (keep neutral URL here)
      aplPar.put(InitParams.PN_BASE_URL, "http://127.0.0.1/");
      // The max time for loading of an URL in the browser
      // This constant is a base for calculation of all used timeouts, so
      // increasing it - will scale all timeouts
      aplPar.put(InitParams.PN_PAGE_LOAD_TIMEOUT, "30001");

      aplPar.put(InitParams.PN_USE_WEB_DRIVER, useWebDriver);

      // Configure the instance of the APL interface
      SUITA.Factory.aplSetup(aplClassName, aplPar);
      // Initiate connecting with the Selenium-RC server
      // Request opening of the browser by the Selenium-RC server
      SUITA.Factory.aplReset();
   }

   /**
    * This utility method closes the test session on the test-target host and
    * releases the APL instance.
    */
   public void cleanUp() {
      SUITA.Factory.aplClose();
   }

}
