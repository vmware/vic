/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util.testreporter.racetrack;

import static com.vmware.client.automation.util.testreporter.racetrack.LogWorker.setCurrentTestCaseId;

import org.apache.logging.log4j.core.appender.AbstractAppender;
import org.apache.logging.log4j.core.config.plugins.Plugin;
import org.apache.logging.log4j.core.config.plugins.PluginAttribute;
import org.apache.logging.log4j.core.config.plugins.PluginFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.util.testreporter.TestResultLogger;
import com.vmware.client.automation.util.testreporter.TestSetSpec;
import com.vmware.client.automation.workflow.BaseTest;

/**
 * Class for posting test output to Racetrack sites (racetrack and
 * racetrack-dev) using RacetrackWebServices More info here -
 * https://wiki.eng.vmware.com/RacetrackWebServices
 *
 * This class is registered as log4j2 appender plugin. In order to use be registered
 * add package scan to the log4j2.xml config file:
 * <Configuration status="WARN"
 * packages="com.vmware.client.automation.util.testreporter.racetrack">
 *
 * Create appender in log4j2.xml config file:
 * <RacetrackLogger name="RacetrackLogger" proprtiesPath="<fullpath to prop file>"/>
 *
 * Configure which logs you want to send to the racetrack logger.
 *
 */
@Plugin(name = "RacetrackLogger", category = "Core",
      elementType = "appender", printObject = true)
public class RacetrackLogger extends AbstractAppender implements TestResultLogger {

   private static final long serialVersionUID = 102424883394397658L;
   private static final Logger _logger =
         LoggerFactory.getLogger("RacetrackLogger");

   private static final String NULL_STRING = "null";
   private static final String EMPTY_STRING = "<empty string>";
   private static final String RACETRACK_RESULTS_PAGE_EXTENSION = "/result.php?id=";

   private RacetrackConnectionSpec _connectionSpec;
   private RacetrackWebservice _racetrackWebService;
   private boolean _isConnected;

   private String _testCaseId;
   private String _resultId = null;
   private Thread _logThread = null;
   private String _testRunnerOs = null;
   private String _language = null;

   /**
    * Initializes the racetrack logger - stores the connection spec.
    *
    * @param connectionSpec the preferred connection spec
    */
   public RacetrackLogger(String name, String propertiesFile) {
      super("RacetrackLogger", null, null);

      // if no config file is provided - do not connect to the racetrack server
      if (Strings.isNullOrEmpty(propertiesFile)) {
         setConnected(false);
         return;
      }

      RacetrackConnectionSpec racetrackConnectionSpec =
            new RacetrackConnectionSpec(propertiesFile);
      this._connectionSpec = racetrackConnectionSpec;
      this._racetrackWebService = new RacetrackWebservice(racetrackConnectionSpec);

      // Start connection to external test output logging system
      if (connect()) {

         // Create test set spec
         RacetrackTestSetSpec racetrackTestSetSpec =
               new RacetrackTestSetSpec(propertiesFile);
         // Start new test set
         startTestSet(racetrackTestSetSpec);
      } else {
         throw new RuntimeException("External reported object is null!");
      }
   }

   /**
    * Creates racetrack logger with specific settings in a properties file.
    *
    * @param name the name of the logger
    * @param propertiesPath pat to the properties file - absolute path
    * @return instantiated racetrack logger
    */
   @PluginFactory
   public static RacetrackLogger createAppender(@PluginAttribute("name") String name,
         @PluginAttribute("proprtiesPath") String propertiesPath) {
      return new RacetrackLogger(name, propertiesPath);
   }

   @Override
   public boolean connect() {
      // TODO: Find a way to validate the connection is live.
      setConnected(true);
      return _isConnected;
   }

   /**
    * Whether we post test output in a threaded way or not.
    *
    * @return true if we use threaded logging
    */
   public boolean isThreadedRacetrack() {
      return _connectionSpec.getThreadedLogging();
   }

   @Override
   public String getTestCaseId() {
      return _testCaseId;
   }

   /**
    * Gets the test runner machine OS
    *
    * @return String containing test runner OS
    */
   public String getTestRunnerOs() {
      return _testRunnerOs;
   }

   /**
    * Gets the display language for the OS where test runs
    *
    * @return language
    */
   public String getLanguage() {
      return _language;
   }

   @Override
   public boolean getConnected() {
      return _isConnected;
   }

   @Override
   public void logTestComment(String comment) {
      if (isThreadedRacetrack()) {
         LogWorker.queueEvent(new CommentEvent(this, comment));
      } else {
         doTestCaseComment(comment);
      }
   }

   @Override
   public void logTestVerification(String description, String actualValue,
         String expectedValue) {
      if (isThreadedRacetrack()) {
         LogWorker.queueEvent(new VerificationEvent(this, description, actualValue,
               expectedValue));
      } else {
         doTestCaseVerification(description, actualValue, expectedValue);
      }
   }

   @Override
   public String startTestSet(TestSetSpec testSetSpec) {
      if (getConnected()) {
         try {
            // Check if we already have existing resultId to use
            _resultId = testSetSpec.getExistingResultId();
            if (Strings.isNullOrEmpty(testSetSpec.getExistingResultId())) {
               _resultId =
                     _racetrackWebService.testSetBegin(
                           testSetSpec.getBuildNumber(),
                           testSetSpec.getTestOwner(),
                           testSetSpec.getProductName(),
                           testSetSpec.getTestSetDescription(),
                           testSetSpec.getBrowserOs(),
                           "",
                           testSetSpec.getBranch(),
                           "",
                           testSetSpec.getBuildType(),
                           testSetSpec.getTestType(),
                           testSetSpec.getLanguage());
            }
         } catch (Exception e) {
            _logger.error("Problem occured when starting test set: {}", e.getMessage());
         }

         if (isThreadedRacetrack()) {
            _logger.info("Starting new thread for racetrack logging");
            _logThread = new Thread(new LogWorker());
            _logThread.setDaemon(true);
            _logThread.start();
         }
      } else {
         _logger.info("Record to racetrack is not set or set to false. "
               + "Hence not recording to Racetrack.");
      }

      // Print the racetrack URL
      String racetrackUrl =
            this._connectionSpec.getTestLoggerURL() + RACETRACK_RESULTS_PAGE_EXTENSION
                  + _resultId;
      _logger.info("Racetrack URL: {}", racetrackUrl);

      setTestRunnerOs(testSetSpec.getBrowserOs());
      setLanguage(testSetSpec.getLanguage());

      return _resultId;
   }

   @Override
   public boolean endTestSet() {
      boolean endTestSet = false;

      // Stop the logging thread.
      if (isThreadedRacetrack()) {
         // Insert a special event to the queue which will trigger the thread
         // to finish.
         LogWorker.queueEvent(new TerminateEvent(this));
         try {
            // Wait for the log thread to finish.
            _logger.info("[Main Thread] Waiting for thread to terminate");
            _logThread.join();
         } catch (InterruptedException e) {
            _logger.error("Received InterruptedExpection while waiting "
                  + "for the log thread to finish: " + e.getMessage());
         }
      }

      try {
         if (getConnected()) {
            _racetrackWebService.testSetEnd(_resultId);
            endTestSet = true;
         }
      } catch (Exception e) {
         e.printStackTrace();
      }

      return endTestSet;
   }

   @Override
   public String startTestCase(String methodName, String packageName, String description) {
      String retVal = null;
      if (getConnected()) {
         if (isThreadedRacetrack()) {
            LogWorker.queueEvent(new CaseBeginEvent(this, methodName, packageName,
                  description));
         } else {
            retVal = doTestCaseBegin(methodName, packageName, description);
         }
      }

      return retVal;
   }

   @Override
   public void endTestCase(final String testCaseResult) {
      if (getConnected()) {
         if (isThreadedRacetrack()) {
            LogWorker.queueEvent(new CaseEndEvent(this, testCaseResult));
         } else {
            doTestCaseEnd(testCaseResult);
         }
      }
   }

   @Override
   public void postScreenshot(final String screenShot) {
      if (isThreadedRacetrack()) {
         LogWorker.queueEvent(new ScreenshotEvent(this, screenShot));
      } else {
         doTestCaseCaptureScreenShot(screenShot);
      }
   }

   @Override
   public void postScreenshot() {
      postScreenshot(null);
   }

   /**
    * Gets the Racetrack result ID
    * 
    * @return result ID
    */
   public String getResultId() {
      return _resultId;
   }

   // ======================================
   // Classes for threaded logging

   public class CaseBeginEvent extends LogEvent {

      // Event parameters.
      final String methodName;
      final String packageName;
      final String description;

      public CaseBeginEvent(RacetrackLogger logger, final String methodName,
            final String packageName, final String description) {
         super(logger);
         this.methodName = methodName;
         this.packageName = packageName;
         this.description = description;
      }

      @Override
      public void execute() {
         // Call the method implementation.
         _racetrackLogger.doTestCaseBegin(methodName, packageName, description);

         // Additionally, set the new test case ID.
         setCurrentTestCaseId(_racetrackLogger.getTestCaseId());
      }
   }

   /**
    * This event describes the end of the execution of a test case.
    */
   public class CaseEndEvent extends LogEvent {

      // Event parameters.
      final String testCaseResult;

      public CaseEndEvent(RacetrackLogger logger, final String testCaseResult) {
         super(logger);
         this.testCaseResult = testCaseResult;
      }

      @Override
      public void execute() {
         this._racetrackLogger.doTestCaseEnd(testCaseResult);
      }
   }

   /**
    * This log event describes a test case comment.
    */
   public class CommentEvent extends LogEvent {

      // Event parameters.
      private String comment;

      public CommentEvent(RacetrackLogger logger, final String comment) {
         super(logger);
         this.comment = comment;
      }

      @Override
      public void execute() {
         // Execute the implementation.
         this._racetrackLogger.doTestCaseComment(comment);
      }
   }

   /**
    * Describes a verification log event.
    */
   public class VerificationEvent extends LogEvent {
      private String description;
      private String actualValue;
      private String expectedValue;

      public VerificationEvent(RacetrackLogger logger, final String description,
            final String actualValue, final String expectedValue) {
         super(logger);
         this.description = description;
         this.actualValue = actualValue;
         this.expectedValue = expectedValue;
      }

      @Override
      public void execute() {
         this._racetrackLogger.doTestCaseVerification(
               description,
               actualValue,
               expectedValue);
      }
   }

   /**
    * This is a special event that signals the termination of the log thread.
    */
   public class TerminateEvent extends LogEvent {
      public TerminateEvent(RacetrackLogger logger) {
         super(logger);
      }

      @Override
      public void execute() {
         // Simply sets the interrupt flag for the current thread.
         Thread.currentThread().interrupt();
      }
   }

   /**
    * Describes a screenshot log event.
    */
   public class ScreenshotEvent extends LogEvent {
      // Event parameters.
      private final String screenShot;

      public ScreenshotEvent(RacetrackLogger logger, final String screenShot) {
         super(logger);
         this.screenShot = screenShot;
      }

      @Override
      public void execute() {
         this._racetrackLogger.doTestCaseCaptureScreenShot(screenShot);
      }
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Posts the passed screenshot file name to racetrack
    *
    * @param screenShot path and filename of the screenshot file
    */
   private void doTestCaseCaptureScreenShot(final String screenShot) {
      try {
         if (getConnected()) {
            this._racetrackWebService.testCaseScreenshot(_testCaseId, " ", screenShot);
         }
      } catch (Exception e) {
         _logger.error("Exception during screen cacpturing: {}", e);
      }
   }

   /**
    * TestCaseVerification action.
    */
   private void doTestCaseVerification(final String description, String actualValue,
         String expectedValue) {
      try {
         if (getConnected() && getTestCaseId() != null) {
            if ((actualValue == null) || (expectedValue == null)) {
               if (actualValue == null) {
                  actualValue = NULL_STRING;
               }
               if (expectedValue == null) {
                  expectedValue = NULL_STRING;
               }
            } else if (actualValue.trim().isEmpty() || expectedValue.trim().isEmpty()) {
               if (actualValue.trim().isEmpty()) {
                  actualValue = EMPTY_STRING;
               }
               if (expectedValue.trim().isEmpty()) {
                  expectedValue = EMPTY_STRING;
               }
            }
            boolean tcResult = actualValue.equals(expectedValue);
            this._racetrackWebService.testCaseVerification(
                  _testCaseId,
                  description,
                  actualValue,
                  expectedValue,
                  tcResult);
         }
      } catch (Exception e) {
         _logger.error("Exception during test case verification: {}", e);
      }
   }

   /**
    * TestCaseBegin action.
    */
   private String doTestCaseBegin(final String methodName, final String packageName,
         final String description) {
      try {
         if (getConnected()) {
            setTestCaseId(this._racetrackWebService.testCaseBegin(
                  _resultId,
                  methodName,
                  packageName,
                  description,
                  getTestRunnerOs(),
                  "",
                  getLanguage()));
            // _logger.info("testCaseId: {}", _testCaseId);
         }
      } catch (Exception ex) {
         _logger.info("Exception :: {}", ex.getMessage());
         _logger.info("Stack Trace == {}", ex.getStackTrace().toString());
         setConnected(false);
      }

      return getTestCaseId();
   }

   /**
    * TestCaseComment action.
    */
   private void doTestCaseComment(final String sMessage) {
      try {
         if (getConnected() && getTestCaseId() != null) {
            this._racetrackWebService.testCaseComment(_testCaseId, sMessage);
         }
      } catch (Exception e) {
         if (e.getMessage().contains("not be null or an empty string")) {
            _logger.info(e.getMessage());
         } else {
            _logger.error("Exception during test case verification: {}", e);
         }
      }
   }

   /**
    * TestCaseEnd action.
    */
   private void doTestCaseEnd(final String testCaseResult) {
      try {
         if (getConnected()) {
            this._racetrackWebService.testCaseEnd(_testCaseId, testCaseResult);
         }
      } catch (Exception e) {
         _logger.error("Exception during test case verification: {}", e);
      }
   }

   /**
    * Sets connected state for the Racetrack connection.
    *
    * @param true if connection is established
    */
   private void setConnected(boolean isConnected) {
      _isConnected = isConnected;
   }

   /**
    * Sets the test case id.
    *
    * @param testCaseId
    */
   private void setTestCaseId(String testCaseId) {
      _testCaseId = testCaseId;
   }

   /**
    * Sets the test runner OS string.
    *
    * @param testRunnerOs
    */
   private void setTestRunnerOs(String testRunnerOs) {
      _testRunnerOs = testRunnerOs;
   }

   /**
    * Sets the display language for the OS where test runs.
    *
    * @param language
    */
   private void setLanguage(String language) {
      _language = language;
   }

   // ======================================
   // Logging event handler

   @Override
   public void append(org.apache.logging.log4j.core.LogEvent logEvent) {

      // if not connected skip logging
      if (!getConnected()) {
         return;
      }

      String loggerName = logEvent.getLoggerName();

      Object[] parameters = null;
      switch (loggerName) {
      case BaseTest.EXT_REPORT_START_CASE:
         parameters = logEvent.getMessage().getParameters();
         if (parameters.length < 3) {
            throw new IllegalArgumentException(
                  "ExtRepStartCase has wrong number of parameters");
         }

         startTestCase(
               parameters[0].toString(),
               parameters[1].toString(),
               parameters[2].toString());
         break;
      case BaseTest.EXT_REPORT_END_SET:
         endTestSet();
         break;
      case BaseTest.EXT_REPORT_END_CASE:
         parameters = logEvent.getMessage().getParameters();
         if (parameters.length < 1) {
            throw new IllegalArgumentException(
                  "ExtRepEndCase has wrong number of parameters");
         }

         endTestCase(parameters[0].toString());
         break;
      case BaseTest.EXT_REPORT_CASE_COMMENT:
         parameters = logEvent.getMessage().getParameters();
         if (parameters.length < 1) {
            throw new IllegalArgumentException(
                  "ExtRepEndCase has wrong number of parameters");
         }

         logTestComment(parameters[0].toString());
         break;
      case BaseTest.EXT_REPORT_CASE_VERIFICATION:
         parameters = logEvent.getMessage().getParameters();
         if (parameters.length < 3) {
            throw new IllegalArgumentException(
                  "ExtRepCaseVerification has wrong number of parameters");
         }

         logTestVerification(
               parameters[0].toString(),
               parameters[1].toString(),
               parameters[2].toString());
         break;
      default:
         // check if the log is not a screenshot
         if (logEvent.getMessage().getFormattedMessage()
               .startsWith(BaseTest.EXT_REPORT_SCREENSHOT_PREFIX)) {
            parameters = logEvent.getMessage().getParameters();
            if (parameters.length < 1) {
               throw new IllegalArgumentException(
                     "Screenshot logging has wrong number of parameters");
            }

            postScreenshot(parameters[0].toString());
         }

         // ignore all logging from racetrack logger
         // ignore any logs before test case start
         if (RacetrackLogger.class.getName().equals(loggerName)
               || getTestCaseId() == null) {
            break;
         }

         // by default log as test case comment
         logTestComment(logEvent.getMessage().getFormattedMessage());

         break;
      }
   }
}
