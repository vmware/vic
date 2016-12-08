/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele;

import java.util.Arrays;

import org.openqa.selenium.NoSuchElementException;
import org.openqa.selenium.TimeoutException;
import org.openqa.selenium.UnhandledAlertException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.thoughtworks.selenium.SeleniumException;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.SpecialStateHandler;
import com.vmware.suitaf.apl.SpecialStates;

public class SeleSpecialStateHandlers {

   protected static final Logger _logger = LoggerFactory.getLogger(SeleSpecialStateHandlers.class);

    public static final SpecialStateHandler getHandler(SpecialStates state) {
        // TODO: Implement the special state handling
        // ===================================================================
        switch (state) {

        // Handle the IE certificate error page.
        case WIN_IE_CERT_ERROR_OVERRIDE_LINK:
            return new SpecialStateHandler() {
                String certOverrideLinkID = "overridelink";
                String winTitleMatch = "certificate error";
                String webDriverCertOverrideIEJS = "javascript:document.getElementById('overridelink').click()";
                String IE_BROWSER = "internet explorer";
                String TIMEOUT_EXCEPTION_MSG = "Timed out after";

                @Override
                public void passOptionalData(Object... optionals) {
                }

                @Override
                public void stateHandle() {
                    SeleAPLImpl apl = ((SeleAPLImpl) SUITA.Factory.aplx());

                    // Code for the Selenium RC
                    if(apl.seleLinks.selenium != null) {
                       apl.seleLinks.selenium.click(certOverrideLinkID);
                       try {
                          // Timeout of 30 seconds for waiting to load the page
                          apl.seleLinks.selenium.waitForPageToLoad("30000");
                       } catch (SeleniumException se) {
                          // Check if timeout exception
                          if(se.getMessage().contains(TIMEOUT_EXCEPTION_MSG)) {
                             // check the certificate error is still present after page load timeout
                             if(apl.seleLinks.selenium.getTitle().toLowerCase().
                                   contains(winTitleMatch)) {
                                throw se;
                             }
                             return;
                          }
                          throw se;
                       }
                       return;
                    } else if (apl.seleLinks.driver != null &&
                          apl.seleLinks.driver.toString().contains(IE_BROWSER)) {
                       try {
                          apl.seleLinks.driver.navigate().to(webDriverCertOverrideIEJS);
                       } catch (TimeoutException te) {
                          // check the certificate error is still present after page load timeout
                          if(apl.seleLinks.driver.getTitle().contains(winTitleMatch)) {
                             throw te;
                          }
                       }
                    }
                }

                @Override
                public boolean stateRecognize() {
                   SeleAPLImpl apl = ((SeleAPLImpl) SUITA.Factory.aplx());

                    // NOTE: it is a quick and dirty implementation for the web
                    // driver to be used to run all the BAT tests.
                    // TODO: rkovachev provide proper implementation for recognizing
                    // certificate error on IE.
                   if(apl.seleLinks.driver != null) {
                      try {
                         return apl.seleLinks.driver.getTitle().toLowerCase().contains(winTitleMatch);
                      } catch (NoSuchElementException e) {
                         return true;
                      } catch (UnhandledAlertException uae) {
                         // When the CIP is installed the getTitle method throws UnhandledAlertException as the the IE WebdDiver is configured.
                         _logger.warn("UnhandledAlertException is thrown while trying to get the page titel:" + uae.getMessage());
                         return false;
                      }
                   }

                   String[] allLinks = null;
                   String[] winTitles = null;
                   try {
                      allLinks = apl.seleLinks.selenium.getAllLinks();
                      winTitles = apl.seleLinks.selenium.getAllWindowTitles();
                   } catch (Throwable e) {  }
                   // If previous calls failed or returned null - no detection
                   if (allLinks == null || winTitles == null) {
                      return false;
                   }
                   // First necessary condition for recognition
                   if (!Arrays.asList(allLinks).contains(certOverrideLinkID)) {
                      return false;
                   }
                   // Second necessary condition for recognition
                   for (String title : winTitles) {
                      if (title.toLowerCase().contains(winTitleMatch)) {
                         return true;
                      }
                   }
                   return false;
                }
            };
        // ===================================================================
        // ===================================================================
        default:
            return state.blankHandler;
        }
    }
}

