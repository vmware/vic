/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.vcuilib.commoncode.GlobalFunction;
import com.vmware.client.automation.vcuilib.commoncode.IDConstants;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.flexui.selenium.MethodCallUtil;

/**
 * Utility class that will hold common methods for working with Flex UI components.
 */
public class UIUtil {

   // Logger utility
   private static final Logger _logger = LoggerFactory.getLogger(UIUtil.class);

   /**
    * FlexSelenium framework returns all UIComponent properties in following format -
    * "propertyName=propertyValue" (i.e "visible=true"). Current util method will parse
    * the property text and will return only the property value ("true").
    *
    * @return property value
    */
   public static String getPropertyValue(String propertyPhrase) {

      // Sanity check
      if (Strings.isNullOrEmpty(propertyPhrase)) {
         return propertyPhrase;
      }

      String[] words = propertyPhrase.split("=");

      if (words.length != 2) {
         _logger.warn("Improper property phrase - '" + propertyPhrase + "'");
         return propertyPhrase;
      }

      // return the second part of the phrase
      return words[1];
   }

   /**
    * Checks whether a confirmation dialog is visible and if so, click Yes in it.
    */
   public static void confirmDialogIfVisible() {
      if (MethodCallUtil.getVisibleOnPath(BrowserUtil.flashSelenium,
            IDConstants.ID_CONFIRMATION_DIALOG)) {

         _logger.info("Confirmation dialog found!");
         GlobalFunction.handleConfirmationDialog(true);
      } else if (MethodCallUtil.getVisibleOnPath(BrowserUtil.flashSelenium,
            IDConstants.ID_YES_NO_DIALOG)) {

         _logger.info("Yes/No confirmation dialog found!");
         GlobalFunction.verifyAndHandleYesNoMsgBox(true, false, null);
      } else {
         _logger.error("No confirmation dialog is found!");
      }
   }
}
