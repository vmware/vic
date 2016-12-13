/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.suitaf.apl.sele;

import java.awt.event.KeyEvent;
import java.net.MalformedURLException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebDriverException;
import org.openqa.selenium.remote.RemoteWebDriver;
import org.openqa.selenium.remote.SessionId;
import org.openqa.selenium.remote.SessionNotFoundException;

import com.thoughtworks.selenium.FlashSelenium;
import com.thoughtworks.selenium.Selenium;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.selenium.MethodCallUtil;
import com.vmware.flexui.selenium.VMWareSelenium;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.Category;
import com.vmware.suitaf.apl.webdriver.WebDriverNode;
import com.vmware.suitaf.apl.webdriver.WebDriverUtils;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.suitaf.util.FailureCauseHolder;
import com.vmware.suitaf.util.Logger;
import com.vmware.suitaf.util.TypedParams;

/**
 * This class encapsulates the links to the main Selenium wrapper components.
 * These components are the base for most of the functions executed using the
 * Selenium library.<br>
 *
 * @author dkozhuharov
 */
public final class SeleLinks {
   public static final String PATH_TO_DESKTOP = null;
   public static final String PATH_TO_BROWSER_APP = null;
   public static final String PATH_TO_BROWSER_WIN = null;
   public static String PATH_TO_FLEX_APP = null;

   // =============================================================
   // Instance final fields and the constructor that initializes them
   // =============================================================
   private final SeleAPLImpl apl;
   private final String serverIP;
   private final int serverPort;
   private final String browser;
   private final String browserArgs;
   private final String baseURL;
   // Retry count for checking if the flash is loaded
   private final int FLASH_LOAD_CHECK_RETRY = 3;

   private String webDriverConnURLFormat = "http://%s:%s/wd/hub";
   private static WebDriverNode webDriverNode;

   public SeleLinks(SeleAPLImpl apl, String serverIP, int serverPort,
         String browser, String browserArgs, String baseURL) {
      this.apl = apl;
      this.serverIP = serverIP;
      this.serverPort = serverPort;
      this.browser = browser;
      this.browserArgs = browserArgs;
      this.baseURL = baseURL;
   }

   // =============================================================
   // The link holders to the main Selenium wrapper components.
   // =============================================================
   Selenium selenium = null;
   private FlashSelenium flashSelenium = null;
   WebDriver driver = null;

   // =============================================================
   // Standard time-outs.
   // =============================================================
   protected Integer PageLoadTimeout = null;
   protected Integer SmallFixTimeout = null;
   protected Integer BigFixTimeout = null;

   public void setPageLoadTimeout(String paramPageLoadTimeout) {
      PageLoadTimeout = Integer.valueOf(paramPageLoadTimeout);
      SmallFixTimeout = PageLoadTimeout * 4;
      BigFixTimeout = PageLoadTimeout * 8;
   }

   // =============================================================
   // Set of methods for repair/reset of the Selenium objects references
   // =============================================================

   private static final ArrayList<Selenium> CleanupList = new ArrayList<Selenium>();
   private static final ArrayList<WebDriver> driverCleanupList = new ArrayList<WebDriver>();

   static {
      Runtime.getRuntime().addShutdownHook(new Thread() {
         @Override
         public void run() {
            SUITA.Factory.UI_AUTOMATION_TOOL.logger
                  .warn("Selenium connection clean is invoked!");
            // Clean the Selenium RC sessions
            for (Selenium s : CleanupList) {
               s.stop();
            }
            // Clean the WebDriver sessions
            for (WebDriver driver : driverCleanupList) {
               SUITA.Factory.UI_AUTOMATION_TOOL.logger
                     .warn("Clean up driver instance to cleanup list:"
                           + driver.toString());
               driver.quit();
            }
         }
      });
   }

   /** Opens the browser instance with default parameters. */
   protected void openBrowserInstance() {
      selenium = null;
      flashSelenium = null;
      driver = null;

      int retries = 3;
      FailureCauseHolder fch = new FailureCauseHolder();
      do {
         Thread browserOpener = new Thread() {
            @Override
            public void run() {
               if (apl.useWebDriver) {
                  InitializeWebDriver driverInit;

                  String hubAddress = String.format(webDriverConnURLFormat,
                        serverIP, serverPort);
                  try {
                     driverInit = new InitializeWebDriver(
                           WebDriverType.getWebDriverType(browser), hubAddress, browserArgs);

                     driver = driverInit.getDriver();
                     SUITA.Factory.UI_AUTOMATION_TOOL.logger
                           .warn("Add web driver instance to cleanup list:"
                                 + driver.toString());
                     driverCleanupList.add(driver);
                  } catch (MalformedURLException e) {
                     apl.hostLogger.error(e.getMessage());
                  }
                  driver.get(hubAddress);

                  // Files in the node which is the selenium ip so we can easily
                  // triage the test run agaisnt selenium grid
                  SessionId sessionId =
                        ((RemoteWebDriver) driver).getSessionId();
                  webDriverNode = WebDriverUtils.getWebDriverNode(
                        serverIP,
                        serverPort,
                        sessionId);
                  SUITA.Factory.UI_AUTOMATION_TOOL.logger.info(
                        "Running on webDriverNode: " + webDriverNode);

                  // Maximizing the window luckily works for all browsers
                  driver.manage().window().maximize();
               } else {
                  selenium = new VMWareSelenium(serverIP, serverPort, browser,
                        baseURL);
                  CleanupList.add(selenium);

                  selenium.start();
                  apl.hostLogger.info("Selenium Browser Started");

                  // This is a workaround for the IE configuration dialog that
                  // pops-up occasionally
                  selenium.keyDownNative(KeyEvent.VK_ESCAPE + "");
                  selenium.keyUpNative(KeyEvent.VK_ESCAPE + "");

                  selenium.windowMaximize();
                  selenium.windowFocus();
                  selenium.setTimeout("" + apl.seleLinks.PageLoadTimeout);

                  MethodCallUtil.CUSTOM_TIMEOUT = apl.seleLinks.PageLoadTimeout;

                  apl.hostLogger.info("Selenium Links Initialized");
               }
            }
         };
         browserOpener.setDaemon(true);
         browserOpener.start();

         try {
            browserOpener.join(this.BigFixTimeout);
         } catch (InterruptedException e) {
         }

         // Check is the browser is open and available
         // NOTE: Workaround for the Selenium RC
         if (selenium != null) {
            try {
               selenium.getAllWindowIds();
               fch.setCause(null);
            } catch (Throwable e) {
               fch.setCause(e);
               apl.hostLogger.warn("Session check fail. Try to type 'N'.");
               // Try to close any JavaScript error dialogs
               selenium.keyPressNative("N".codePointAt(0) + "");
               // Try to close the failed Selenium session window
               selenium.keyDownNative(KeyEvent.VK_ALT + "");
               selenium.keyPressNative(KeyEvent.VK_F4 + "");
               selenium.keyUpNative(KeyEvent.VK_ALT + "");
               // Try to close the failed browser or its dialog
               selenium.keyDownNative(KeyEvent.VK_ALT + "");
               selenium.keyPressNative(KeyEvent.VK_F4 + "");
               selenium.keyUpNative(KeyEvent.VK_ALT + "");
            }
         }

         if (fch.isNull()) {
            break;
         } else {
            if (browserOpener.isAlive()) {
               browserOpener.interrupt();
            }
            closeBrowserInstance();
         }

         retries--;
      } while (retries > 0);

      fch.escalateCause();
   }

   /** Closes the browser instance. */
   protected void closeBrowserInstance() {
      try {
         if (selenium != null) {
            apl.hostLogger.info("Trying to close the browser.");
            selenium.close();
            selenium.stop();
         } else if (driver != null) {
            apl.hostLogger.info("Trying to close the Web Driver browser.");
            driver.quit();
         }
      } catch (SessionNotFoundException snfe) {
         Logger.error(snfe);
      } catch (WebDriverException wde) {
         Logger.error(wde);
      } catch (Throwable e) {
         apl.hostLogger.error("Failed browser close: " + e.getMessage());
      }

      driver = null;
      selenium = null;
      flashSelenium = null;
   }

   public void openUrl(String url) {
      if (selenium != null) {
         selenium.open(url);
      } else if (driver != null) {
         driver.get(url);
      } else {
         throw new RuntimeException("Selenium session not initialized!");
      }
      apl.hostLogger.info("Opened URL: " + url);

      // Release the old FlashSelenium API library instance and the UID cache
      flashSelenium = null;
   }

   /**
    * Returns the selenium object of the BrowserInstance.
    *
    * @return Selenium object of the BrowserInstance.
    */
   public Selenium selenium() {
      return selenium;
   }

   /**
    * Returns the WebDriver object of the BrowserInstance.
    *
    * @return WebDriver object of the BrowserInstance.
    */
   public WebDriver driver() {
      return driver;
   }

   /**
    * Returns the flashSelenium object of the BrowserInstance.
    *
    * @return FlashSelenium object of the BrowserInstance.
    */
   public FlashSelenium flashSelenium() {
      if (flashSelenium == null) {
         // Retrieve the name of the FlexApp on the opened web page
         // Will throw exception if no flex-application is loaded on page
         // String res = selenium.getEval(NativeCommand.FLEXAPP_ID.encode());
         // String id = NativeCommand.unpackArray(res)[0];
         // TODO: rkovachev - detect id dynamically
         String id = "container_app";

         // Initialize internal fields
         SeleLinks.PATH_TO_FLEX_APP = id;
         // Create a new instance of the FlashSelenium API library
         if (selenium != null) {
            flashSelenium = new FlashSelenium(selenium, id);
         } else if (driver != null) {
            // make sure that we are in the default content of the page
            driver.switchTo().defaultContent();
            flashSelenium = new FlexExtentionWrapper(driver, id);
         } else {
            throw new RuntimeException("Selenium session not initialized!");
         }

         int count = 0;
         int percent = 0;

         try {
            while ((percent = flashSelenium.PercentLoaded()) < 100
                  && count < FLASH_LOAD_CHECK_RETRY) {
               apl.hostLogger.info("Flash Loaded: " + percent + "%");
               CommonUtils.sleep(SUITA.Environment.getUIOperationTimeout());
               count++;
            }

            apl.hostLogger.info("Flash is loaded in " + percent + " %");
         } catch (org.openqa.selenium.WebDriverException wde) {
            apl.hostLogger.info("It is not Flex client");
            flashSelenium = null;
         }
      }
      return flashSelenium;
   }

   // ===================================================================
   // Component tree exploratory methods
   // ===================================================================
   /**
    * See {@link SeleLinks#findAny(boolean, Object...)} for parameterization
    * description.
    *
    * @param findParams
    *           - parameters of the "find" command
    * @return - the object that was retrieved or null
    */
   public <T> T find(Class<T> verifyClass, Object... findParams) {
      // Parse the "findAny" parameters and adjust its operational mode
      TypedParams tPars = new TypedParams(true, findParams);
      tPars.extractAll(Class.class);
      tPars.append(Arrays.asList((Object) verifyClass));

      return verifyClass.cast(find(tPars.toArray()));
   }

   /**
    * See {@link SeleLinks#findAny(boolean, Object...)} for parameterization
    * description.
    *
    * @param findParams
    *           - parameters of the "find" command
    * @return - the object that was retrieved or null
    */
   public DisplayObject find(Object... findParams) {
      List<DisplayObject> tmpRes = findAny(true, findParams);
      if (tmpRes != null && tmpRes.size() == 1)
         return tmpRes.get(0);
      else
         return null;
   }

   /**
    * See {@link SeleLinks#findAny(boolean, Object...)} for parameterization
    * description.
    *
    * @param findParams
    *           - parameters of the "findAll" command
    * @return - list of objects that were retrieved or null
    */
   public List<DisplayObject> findAll(Object... findParams) {
      return findAny(false, findParams, false);
   }

   public static WebDriverNode getWebDriverNode(){
      return webDriverNode;
   }

   /**
    * This method implements complex functionality for tracking and retrieving
    * of Silk TestObject instances. <BR>
    * To achieve maximum flexibility and applicability for diversity of
    * use-cases, the method uses the variable argument "findParams". It can
    * receive a list of objects of varying number and type.<br>
    * This allows for each input parameter to be provided 0, 1 or many values.
    * On the other hand - to allow deterministic parameter value recognition -
    * it is accepted as a contract that every parameters value will differ from
    * others by its class. Here follows a list of recognized classes and as what
    * parameter they are recognized:
    *
    * <li> {@link DirectID}; expected 1 times;<br>
    * Value of this type represents an ID descriptor to be used to locate
    * {@link DisplayObject}s. <br>
    * There must be exactly one such argument. <br>
    * If more than one instance is provided - the last one is used.
    *
    * <li> {@link Boolean}; expected 0 or 1 times;<br>
    * Value of this type represents the "throwNotFoundException" flag. <br>
    * Value of true forces throwing exception if no TestObject is found or if
    * more than one object is found, but just one is required. <br>
    * Value of false - causes null to be returned instead of exception <br>
    * If omitted - true is assumed. <br>
    * If more than one instance is provided - the last one is used.
    *
    * <li> {@link Class}; expected 0, 1 times;<br>
    * Values of this type represent a Silk class which will be used for type
    * checking of retrieved components. The components that does not apply are
    * removed from the resulting list. <br>
    * If omitted - the class TestObject is used: <br>
    * If more than one instance is provided - the last one is used.
    *
    * <br>
    * <br>
    *
    * @param findJustOne
    *           - if <b>true</b> forces retrieval of just one instance. If more
    *           or less instance were found - return null or throw Exception
    * @param findParams
    *           - the set of typed parameters
    * @return list of {@link DisplayObject} instances or <b>null</b> if none
    *         found
    */
   private List<DisplayObject> findAny(boolean findJustOne,
         Object... findParams) {
      // Parse the "findAny" parameters and adjust its operational mode
      TypedParams tPars = new TypedParams(true, findParams);

      DirectID id = tPars.extractLast(DirectID.class);
      boolean throwNotFoundException = tPars.extractLast(Boolean.class, true)
            .booleanValue();
      Class<?> verifyClass = tPars
            .extractLast(Class.class, DisplayObject.class);
      tPars.assertAllUsed();

      List<DisplayObject> result = null;
      FailureCauseHolder ward = new FailureCauseHolder();

      // Validate input parameters
      if (id != null) {
         result = findBase(id, verifyClass, ward);
      }

      // Check if more than one object was matched
      // and we are in single find mode
      if (findJustOne && result != null && result.size() > 1) {
         apl.hostLogger.warn("One object for '" + id + "' was requested"
               + " but " + result.size()
               + " were found! Picking just the first one.");

         // Log every single object with all its properties
         for (int i2 = 0; i2 < result.size(); i2++) {
            apl.seleHelper.logFlexObject(result.get(i2), 1,
                  new Class[] { DisplayObject.class }, Category.ID,
                  Category.INFO, Category.ACCESSIBILITY);
         }

         // Leave just the first one because only it is usable
         DisplayObject to = result.get(0);
         result.clear();
         result.add(to);
      }

      // If exception throwing mode is activated and reason for failure
      // has been found
      if (throwNotFoundException) {
         ward.escalateCause();
      }

      // Return the list of objects retrieved or null
      return result;
   }

   private List<DisplayObject> findBase(DirectID id, Class<?> verifyClass,
         FailureCauseHolder ward) {
      // Prepare the local resources for the method run
      List<DisplayObject> result = new ArrayList<DisplayObject>();
      try {
         if (id.hasIndex()) {
            result.add(ProxyFactory.getInstance(apl, id.getID()));
         } else {
            // Check if the ID matches 0, 1 or more components (max 200)
            for (int i = 0; i < 200; i++) {
               result.add(ProxyFactory.getInstance(apl, id.getID() + "[" + i
                     + "]"));
            }
         }
      } catch (Throwable e1) {
         // If no component was found - store the exception as failure cause
         if (result.size() == 0) {
            ward.setCause(e1);
         }
      }

      // Apply additional restrictions on the candidate list
      int i1 = 0;
      while (i1 < result.size()) {
         DisplayObject to = result.get(i1);
         // If was a component without UID - skip it
         if (to == null) {
            result.remove(i1);
         }
         // Apply additional restrictions on the candidate list
         else if (verifyClass != null && !verifyClass.isInstance(to)) {
            result.remove(i1);
         }
         // Apply property value restrictions
         else if (!id.matchProperties(to, Category.REGULAR)) {
            result.remove(i1);
         }
         // Move to next result entry if no restrictions apply
         else {
            i1++;
         }
      }

      // Check if no object was found - normalize the result
      if (result != null && result.size() == 0) {
         result = null;
      }

      // Generate a failure cause if no appropriate results were found
      if (ward.isNull() && result == null) {
         String className = ((verifyClass == null) ? "" : " of type: "
               + verifyClass.getSimpleName());
         ward.setCause(new RuntimeException("No object" + className
               + " matches locator: " + id));
      }

      return result;
   }

   /**
    * Get the existing count of the component
    * @param componentId the ID of the component
    * @return the existing count
    */
   public int getExistingCount(DirectID componentId) {
      List<DisplayObject> compList = findAll(componentId);
      return (compList == null) ? 0 : compList.size();
   }}
