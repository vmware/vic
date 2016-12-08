/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.view;

import java.awt.event.KeyEvent;

import org.openqa.selenium.By;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.TimeoutException;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.support.ui.ExpectedCondition;
import org.openqa.selenium.support.ui.ExpectedConditions;
import org.openqa.selenium.support.ui.WebDriverWait;

import com.thoughtworks.selenium.Selenium;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.apl.sele.SeleAPLImpl;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;

public class LoginView extends BaseView {

   private static IDGroup ID_USERNAME_TF =
         IDGroup.toIDGroup("username");

   private static IDGroup ID_PASSWORD_TF =
         IDGroup.toIDGroup("password");

   private static IDGroup ID_LOGIN_BTN =
         IDGroup.toIDGroup("submit");

   private static IDGroup ID_WINDOWS_SESION_CHECK_BTN =
         IDGroup.toIDGroup("sspiCheckbox");

   private static IDGroup ID_IFRAME =
         IDGroup.toIDGroup("websso");

   /**
    * Logs in a user to the NGC.
    *
    * @param user <code>UserSpec</code> instance containing the
    *    credentials for login
    *
    * @throws Exception If the login page is not available or login fails
    */
   public void login(UserSpec user) throws Exception {
      if (isLoginPageOpen()) {
         doLogin(user);
      } else if (isMainPageOpen()) {
         if (!getLoggedInUsername().equalsIgnoreCase(user.username.get())) {
            _logger.warn("Another user is logged in. Trying to login with the correct user");
            logout();
            doLogin(user);
         }
      } else {
         UI.audit.aplFailure(new Exception("Unable to detect login page!"));
      }

      // Push environment variables to the lower frameworks
      initBrowserUtil();
   }

   /**
    * Checks whether the NGC login page is open. The result is based on the
    * visibility of the user name text input.
    *
    * @return True if the login page is open, false otherwise
    */
   public boolean isLoginPageOpen() {
      if(((SeleAPLImpl) SUITA.Factory.apl()).getSelenium() != null) {
         return ((SeleAPLImpl)SUITA.Factory.apl()).getSelenium().isElementPresent(ID_USERNAME_TF.getValue(Property.DIRECT_ID));
      } else if(((SeleAPLImpl) SUITA.Factory.apl()).getWebDriver() != null) {
         try {
            WebDriverWait waiting = new WebDriverWait(
                  ((SeleAPLImpl)SUITA.Factory.apl()).getWebDriver(), 10);

            try {
               waiting.until(ExpectedConditions.frameToBeAvailableAndSwitchToIt(
                     ID_IFRAME.getValue(Property.DIRECT_ID)));
            } catch (TimeoutException tie) {
               _logger.warn("The WebSSO is old!");
               _logger.warn(tie.getMessage());
            }

            waiting.until(ExpectedConditions.visibilityOfElementLocated (
                  By.id(ID_USERNAME_TF.getValue(Property.DIRECT_ID))));

            return true;
         } catch (RuntimeException e) {
            _logger.warn("The login page is not found!");
            return false;
         }
      } else {
         throw new RuntimeException("Selenium session not initialized!");
      }
   }


   public void clickElementJs(WebDriver driver, WebElement e) {
      getExecutor(driver).executeScript("arguments[0].click();", e);
   }

   // ---------------------------------------------------------------------------
   // Private methods

   private JavascriptExecutor getExecutor(WebDriver driver) {
      return ((JavascriptExecutor) driver);
   }

   /**
    * Performs login of a user to the NGC. The method returns when the main
    * page of the NGC is loaded.
    *
    * @param user <code>UserSpec</code> instance containing the
    *    credentials for login
    *
    * @throws Exception If login fails.
    */
   private void doLogin(UserSpec user) throws Exception {
      _logger.info("Enter user credentials in the login form!");
      // Handle WebDriver Session
      if(((SeleAPLImpl)SUITA.Factory.apl()).getWebDriver() != null) {
         WebDriver driver = ((SeleAPLImpl)SUITA.Factory.apl()).getWebDriver();
         WebDriverWait waiting = new WebDriverWait(driver, SUITA.Environment.getUIOperationTimeout() / 10);
         ExpectedCondition<WebElement> condition = ExpectedConditions.presenceOfElementLocated(
               By.id(ID_USERNAME_TF.getValue(Property.DIRECT_ID)));
         waiting.until(condition);

         WebElement usernameField = driver.findElement(By.id(ID_USERNAME_TF.getValue(Property.DIRECT_ID)));
         usernameField.sendKeys(user.username.get());

         WebElement passwordField = driver.findElement(By.id(ID_PASSWORD_TF.getValue(Property.DIRECT_ID)));
         passwordField.sendKeys(user.password.get());
         WebElement submitBtn = driver.findElement(By.id(ID_LOGIN_BTN.getValue(Property.DIRECT_ID)));

         submitBtn.click();

         ExpectedCondition<Boolean> loginVanish = ExpectedConditions.invisibilityOfElementLocated(
               By.id(ID_USERNAME_TF.getValue(Property.DIRECT_ID)));
         try {
            waiting = new WebDriverWait(driver, 5);
            waiting.until(loginVanish);
         } catch(TimeoutException timeoutException) {
            _logger.warn("The login form not submitted!");
            _logger.warn("Try to click on the login button using javascript!");
            ((JavascriptExecutor) driver).executeScript("arguments[0].click();", submitBtn);
            // submitBtn.sendKeys(Keys.ENTER);
         }

      } else if(((SeleAPLImpl)SUITA.Factory.apl()).getSelenium() != null) {
         Selenium selenium = ((SeleAPLImpl)SUITA.Factory.apl()).getSelenium();

         // Workaround for PR 1315781 to simulate native button press to enable the login button view.
         selenium.focus(ID_USERNAME_TF.getValue(Property.DIRECT_ID));
         selenium.keyPressNative(KeyEvent.VK_R + "");

         // Check if the workaround with native key press buggon pass.
         // If for soem reason the browser is not on focus it will fail.
         // For such cases try second workaround - by using the
         // "Use Windows session authentication" checkbox.
         // NOTE: isEditable in selenium RC is equal to isEnabled.
         if(!selenium.isEditable(ID_LOGIN_BTN.getValue(Property.DIRECT_ID))) {
            // Check and un-check on the "Use Windows session authentication"
            // checkbox.
            selenium.click(ID_WINDOWS_SESION_CHECK_BTN.getValue(Property.DIRECT_ID));
            // Timeout as the call is asynchronous and there is
            // no guarantee that the check box is selected.
            Thread.sleep(1000);
            selenium.click(ID_WINDOWS_SESION_CHECK_BTN.getValue(Property.DIRECT_ID));
            // Timeout as the call is asynchronous and there is
            // no guarantee that the check box is selected.
            Thread.sleep(1000);
         }
         selenium.type(ID_USERNAME_TF.getValue(Property.DIRECT_ID), user.username.get());
         selenium.type(ID_PASSWORD_TF.getValue(Property.DIRECT_ID), user.password.get());
         selenium.click(ID_LOGIN_BTN.getValue(Property.DIRECT_ID));
         //Workaround for PR 1174271 - Press Enter Key
         selenium.keyPressNative(KeyEvent.VK_ENTER + "");
      }

      if (!isMainPageOpen()) {
         UI.audit.aplFailure(new Exception("Main page not loaded after entering login credentials"));
      }
   }


   /**
    * Initalize <code>BrowserUtil</code>
    *  (a class that is mainly used by <code>ObjNav</code>).
    * Current method push some environment variables (Selenium reference, screenshot
    * directory, etc.) to VCUI-QE-LIB project (lower level library).
    *
    * TODO mdzhokanov: BrowserUtil class initialization should take place in more
    * generic class. This should be updated once the general architecture of our
    * test scenarios is clearly defined and established.
    *
    * NOTE mdzhokanov: I have already tried to put the BrowserUtil init logic in
    * <code>TestTarget.startUp()</code>.
    * However, the initialization fails in that case. The first location where BrowserUtil
    * init logic successfully passes is right after the login.
    *
    * NOTE: more appropriate place will be <code>SubToolBrowser.open()<code>
    */
   private void initBrowserUtil(){
      BrowserUtil.flashSelenium = ((SeleAPLImpl)SUITA.Factory.apl()).getFlashSelenium();
      BrowserUtil.selenium = ((SeleAPLImpl)SUITA.Factory.apl()).getSelenium();
   }
}
