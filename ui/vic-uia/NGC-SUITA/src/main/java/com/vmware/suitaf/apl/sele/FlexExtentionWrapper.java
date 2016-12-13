/**
 * Copyright 2014 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele;

import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.WebDriver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.thoughtworks.selenium.FlashSelenium;
import com.thoughtworks.selenium.Selenium;

/**
 * The class goal is to provide a Selenium Flash extension for Selenium Web Driver.
 * The main point to extend the FlashSelenium which is the the official Flash extension
 * for Selenium RC is to continue use SUITA Selenium implementation without changes.
 * As the FlashSelenium provides the Flash extension to the Selenium RC it will be reused.
 * For Web Driver session will be provided implementation of javascript invocation mechanism
 * in the class.
 */
public class FlexExtentionWrapper extends FlashSelenium {

   private WebDriver webDriver = null;
   private String flashObjectId = "";
   protected static final Logger _logger = LoggerFactory.getLogger(FlexExtentionWrapper.class);
   private static final String JS_FUNC_FORMAT = "return document.%1$s.%2$s";

   /**
    * Create Selenium Flash extension for Selenium RC session.
    * @param selenium      - Selenium RC session object.
    * @param flashObjectId - Flash object id.
    */
   public FlexExtentionWrapper(Selenium selenium, String flashObjectId) {
      super(selenium, flashObjectId);
   }

   /**
    * Create Selenium Flash extension for Selenium Web Driver session.
    * @param selenium      - Selenium Web Driver session object.
    * @param flashObjectId - Flash object id.
    */
   public FlexExtentionWrapper(final WebDriver webDriver, final String flashObjectId) {
      super(null, "");
      this.webDriver = webDriver;
      this.flashObjectId = flashObjectId;
   }


   /**
    * Invoke the respective javascript function specified by the functionName parameter.
    * @param functionName  - name of the javacript function to be invoked.
    * @param args          - list of arguments to be provided to the invoked function.
    * @return the return value of the invoked javascript function.
    */
   @Override
   public String call(String functionName, String ... args) {
      if(webDriver != null) {
         return callFlashObject(functionName, args);
      } else {
         return super.call(functionName, args);
      }
   }

   /**
    * Use the Web Driver javascript executor to invoke the functionName fucntion.
    * @param functionName  - name of the javacript function to be invoked.
    * @param args          - list of arguments to be provided to the invoked function.
    * @return the return value of the invoked javascript function.
    */
   private String callFlashObject(final String functionName,
         final String... args) {
      String jsFunction = makeJsFunction(functionName, args);
      String jsStatement = String.format(JS_FUNC_FORMAT, flashObjectId, jsFunction);

      _logger.debug("Executing JavaScript: " + jsFunction);

      Object result =
            ((JavascriptExecutor) webDriver).executeScript(jsStatement, new Object[0]);
      String executionResult = result != null ? result.toString() : null;

      _logger.debug("Got JavaScript result: " + executionResult);

      return executionResult;
   }

   /**
    * Constructs the javascript invocation string.
    * @param functionName  - name of the javacript function to be generated.
    * @param args          - arguments for the javascript function to be invoked.
    * @return javascript to be invoked.
    */
   private String makeJsFunction(final String functionName, final String... args) {
      final StringBuffer functionArgs = new StringBuffer();

      if (args.length > 0) {
         for (int i = 0; i < args.length; i++) {
            if (i > 0) {
               functionArgs.append(",");
            }
            functionArgs.append(String.format("'%1$s'", args[i]));
         }
      }
      return String.format("%1$s(%2$s);", functionName, functionArgs);
   }
}
