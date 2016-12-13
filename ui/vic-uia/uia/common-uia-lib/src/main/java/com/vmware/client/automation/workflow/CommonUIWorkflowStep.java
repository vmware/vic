/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.workflow;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.assertions.Assertion;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * UI steps should extend this class as it provides UI verification
 * screenshot/video etc.
 */
public abstract class CommonUIWorkflowStep extends BaseWorkflowStep {

   protected static final Logger _logger =
         LoggerFactory.getLogger(CommonUIWorkflowStep.class);

   public static final String FAILURE_PREFFIX = "FAIL";

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private static final String EXCEPTION_SCREENSHOT = "EXCEPTION_SCREENSHOT";

   @Override
   /**
    * The method is invoked by the test controller when exception is thrown during step execution.
    * The method creates screenshot of the browser to provide debug info.
    */
   public void logErrorInfo() {
      logScreenshot(EXCEPTION_SCREENSHOT);
   }

   @Override
   protected boolean performVerification(Assertion assertion, boolean verifySafely) {
      _logger.debug("Do UI logging - i.e. take screenshot when fails!");

      boolean result;
      try {
         result = super.performVerification(assertion, verifySafely);
      } catch (RuntimeException validationFailure) {
         result = false;
         throw validationFailure;
      }

      LoggerFactory.getLogger(BaseTest.EXT_REPORT_CASE_VERIFICATION).info("",
            assertion.getDescription(),
            assertion.getActual(),
            assertion.getExpected()
         );

      if(!result) {
         logScreenshot(assertion.getDescription());
      }
      return result;
   }

   /**
    * Helper Method that generates Failure Package prefixes to be used for
    * marking all log entries of one failure
    * @return a unique Failure Package Prefix
    */
   private static final String getFPID() {
      return FAILURE_PREFFIX + String.format("%014d", System.currentTimeMillis());
   }

   /**
    * Take screenshot.
    *
    * @param assertionDescription      assertion description
    */
   protected static void logScreenshot(String assertionDescription) {
      String fileName =
            SUITA.Factory.UI_AUTOMATION_TOOL.audit.snapshotAppScreen(
                  getFPID(),
                  assertionDescription);

      _logger.info(BaseTest.EXT_REPORT_SCREENSHOT_PREFIX, fileName);
   }
}
