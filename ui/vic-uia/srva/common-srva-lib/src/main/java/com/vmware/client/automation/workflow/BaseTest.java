/**
 * Copyright 2012 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.workflow;

import java.io.File;
import java.io.FilenameFilter;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.ITestContext;
import org.testng.ITestNGMethod;
import org.testng.ITestResult;
import org.testng.annotations.AfterClass;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.AfterSuite;
import org.testng.annotations.BeforeMethod;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.Optional;
import org.testng.annotations.Parameters;

import com.vmware.client.automation.servicespec.SeleniumServiceSpec;
import com.vmware.client.automation.testbed.LocalTestBed;
import com.vmware.client.automation.testbed.TestBed;
import com.vmware.client.automation.util.SsoUtil;
import com.vmware.client.automation.workflow.command.CommandController;
import com.vmware.client.automation.workflow.command.CommandException;

/**
 * This class will be used as base class for all vCD UI tests. All classes that
 * contain UI tests should inherit it.
 */
public class BaseTest {

   // External reporting unique logger names
   // Each of these logger names stands for a specific logging operation
   public static final String EXT_REPORT_START_CASE = "ExtRepStartCase";
   public static final String EXT_REPORT_END_SET = "ExtRepEndSet";
   public static final String EXT_REPORT_END_CASE = "ExtRepEndCase";
   public static final String EXT_REPORT_CASE_VERIFICATION = "ExtRepCaseVerification";
   public static final String EXT_REPORT_CASE_COMMENT = "ExtRepCaseComment";
   public static final String EXT_REPORT_SCREENSHOT_PREFIX = "SCREENSHOT";

   protected static final Logger _logger = LoggerFactory.getLogger(BaseTest.class);

   /* Testng group names */
   protected static final String SVS = "svs";
   protected static final String BAT = "bat";
   protected static final String CAT = "cat";
   protected static final String HPTC = "HPTC";
   protected static final String LPTC = "LPTC";
   protected static final String P0 = "P0";
   protected static final String P1 = "P1";
   protected static final String P2 = "P2";
   protected static final String P3 = "P3";
   protected static final String CODB = "CODB";
   protected static final String CONTENT_LIBRARY = "CL";
   protected static final String PERMISSIONS = "PERMISSIONS";
   protected static final String PBM = "PBM";
   protected static final String VDC = "VDC";
   protected static final String NEW_PROV = "NEW_PROVISIONING";
   protected static final String MIGRATE = "MIGRATE";
   protected static final String CLONE = "CLONE";
   protected static final String EDIT_POLICIES = "EDIT_POLICIES";
   protected static final String DEPRECATED = "deprecated";
   protected static final String DELETED = "DELETED";
   protected static final String INTEGRATION_TEST = "INTEGRATION_TEST";

   // If true the test will be executed with the provider workflow
   protected static boolean isProviderRun = false;

   protected static TestBed testBed;
   protected SeleniumServiceSpec _seleniumServiceSpec;

   protected boolean _isCleanupEnabled = true;

   // Used by both work flows
   protected TestScope _testScope = TestScope.FULL;

   /**
    * This annotation will be used for assigning test ID to the test case. It
    * will map the test method to the test definition in the HPQC. It will be
    * used for uploading of test results. Each test method must have TestID
    * annotation set.
    */
   @Retention(RetentionPolicy.RUNTIME)
   public static @interface TestID {
      String[] id();
   }

   /**
    * This method is invoked before the test suite. It sets the test
    * scope, testbed type used for test execution. First it loads the common
    * resources properties provided by the default configuration file.
    *
    * All the parameters are optional and if are not provided the default value
    * will be used.
    *
    * @param isCleanupEnabled
    *           If set to false - the <code>clean</code> methods won't be called
    *           after the test finishes
    * @param testScope
    *           The scope which tests will be executed. Supported values are:
    *           BAT, SAT, UI, FULL.
    * @param testBedType
    *           The type of test bed configuration. Supported values: local,
    *           nimbus.
    *           <ul>
    *           <li>If 'local' value is chosen then testBed configuration will
    *           be read from 'localTestBedResourcesFile' property</li>
    *           <li>If 'nimbus' value is chosen then testBed configuration will
    *           be downloaded from AutomationNimbusService published on
    *           'nimbusTestBedServer' property</li>
    *           </ul>
    * @param localTestBedResourceFile
    *           path to the configuration file used to describe the common
    *           resources for each test. If not path is set the
    *           testbed.properties file in the resources folder of the project
    *           will be used.
    * @param providersResourceFolder
    *
    */
   @BeforeSuite(alwaysRun = true)
   @Parameters({ "isCleanupEnabled", "testScope",
      "localTestBedResourceFile", "providersResourceFolder" })
   public final void setTestBedConfiguration(
         @Optional("true") boolean isCleanupEnabled,
         @Optional("FULL") String testScope,
         @Optional("") String localTestBedResourceFile,
         @Optional("/Users/kjosh/Desktop/vic/ui/vic-uia/resources") String providersResourceFolder) {
//         @Optional("C:\\providerslist\\main") String providersResourceFolder) {

      _logger.info("========= Load TestBed Configuration =========");

      this._isCleanupEnabled = isCleanupEnabled;

      // Set test scope
      this._testScope = TestScope.valueOf(testScope);

      _logger.warn("Requested test scope: {}", testScope);
      _logger.info("Test scope that will be used: {}", _testScope);
      _logger.warn("localTestBedResourceFile: {}", localTestBedResourceFile);
      _logger.warn("providersResourceFolder: {}", providersResourceFolder);

      // Local testbed properties configuration
      if (providersResourceFolder.isEmpty()) {
         _isCleanupEnabled = isCleanupEnabled;
         loadLocalTestBed(localTestBedResourceFile);
      } else {
         isProviderRun = true;
         // Load all providers from the folder
         loadAndRegisterTestbeds(providersResourceFolder);
      }
   }

   @BeforeMethod(alwaysRun = true)
   public final void beforeMethodBase(Object[] parameter, Method m, ITestContext testContext) {
      // Run prepare methods for the non-provider workflow
      if(!isProviderRun) {
         beforeClassBase();
      }

      _logger.info("========= Before Method of Test Case: "
            + getClass().getSimpleName() + " =========");

      this.beforeMethod(parameter, m);
      String testCaseName = this.getClass().getSimpleName();
      _logger.info("========= Test Case: " + testCaseName + " =========");

      // Get Test Description
      String testDescription = getTestDescription(m, testContext);

      // send event for case start
      LoggerFactory.getLogger(EXT_REPORT_START_CASE).info("",
            testCaseName,
            this.getClass().getPackage().getName(),
            testDescription
         );
      _logger.info("Starting test...");
   }

   @AfterSuite(enabled = true, alwaysRun = true)
   public final void afterSuiteMethod() {
      _logger.info("========= After Suite =========");

      // send event for ending test set
      LoggerFactory.getLogger(EXT_REPORT_END_SET).info("");
   }

   public final void beforeClassBase() {
      _logger.info("========= Before Class: " + getClass().getName() + " =========");
      _logger.info("Test runner version: "
            + getClass().getPackage().getImplementationVersion());

      beforeClass();
   }

   /**
    * Override this method if you need to execute some code before starting any
    * test method from a class.
    */
   public void beforeClass() {
      // Do some prerequisites for the tests in this class here.
   }

   @AfterClass(alwaysRun = true)
   public final void afterClassBase() {
      _logger.info("========= After Class " + getClass().getName() + " =========");
      afterClass();

      try {
         if (!isProviderRun) {
            testBed.cleanUp();
         }
      } catch (Exception e) {
         _logger.error(String.format(
               "Cleaning testbed has failed. Exception: %s\n%s",
               e.getMessage(),
               e.getStackTrace().toString()));
      }
   }

   /**
    * Override this method if you need to execute some code after running all
    * test method from a class.
    */
   public void afterClass() {
      // Do some clean up for the tests in this class here.
   }

   /**
    * Override this method with initializing code to run at before each test
    * run. This method is called in
    * {@link UITestBase#beforeMethodBase(Object[], Method)}.
    *
    * @param parameter
    * @param m
    */
   public void beforeMethod(Object[] parameter, Method m) {
      // Do some prerequisites for the test method here.
   }

   @AfterMethod(alwaysRun = true)
   public final void afterMethodBase(Object[] parameter, Method m, ITestResult tr) {
      _logger.info("========= After Method of Test Case: "
            + getClass().getSimpleName() + " =========");
      try {
         this.afterMethod(parameter, m);
      } catch (Exception e) {
         e.printStackTrace();
         _logger.error("Something bad happened: " + e.getMessage());
      }

      _logger.info("Test completed");
      LoggerFactory.getLogger(EXT_REPORT_END_CASE).info("",
            Integer.toString(tr.getStatus()));
   }

   /**
    * Override this method with initializing code to run at after each test run.
    * This method is called in
    * {@link UITestBase#beforeMethodBase(Object[], Method)}.
    *
    * @param parameter
    * @param m
    */
   public void afterMethod(Object[] parameter, Method m) {
      // Do some prerequisites for the test method here.
   }

   /**
    * @return the testBed
    */
   public TestBed getTestBed() {
      return testBed;
   }

   /**
    * Testbed method factory responsible to create testbed based on input
    * parameters.
    */
   private TestBed createTestBed(String localTestBedConfig) {
      return new LocalTestBed(localTestBedConfig);
   }

   /**
    * Load and register testbeds listed in the providersResourceFolder.
    * @param testbedListFolder to iterate for providers configuration
    */
   private void loadAndRegisterTestbeds(String testbedListFolder) {
      final String fileExtensionFilter = ".properties";
      List<String> testBedList = new ArrayList<String>();
      testbedListFolder = testbedListFolder + File.separator;

      File folder = new File(testbedListFolder);
      if(folder.isDirectory()) {
         FilenameFilter fileFilter = new FilenameFilter() {
            @Override
            public boolean accept(File dir, String name) {
               return name.toLowerCase().endsWith(fileExtensionFilter);
            }
         };

         for (File file : folder.listFiles(fileFilter)) {
            testBedList.add(file.getAbsolutePath());
         }
      }
      registerTestBed(testBedList);
   }

   /**
    * Register testbed.
    * @param testBedList list of testbeds to register.
    */
   private void registerTestBed(List<String> testBedList) {
      _logger.info("Register testbeds.");
      CommandController _controller = new CommandController();
      _logger.info("Controller: " + _controller.toString());
      if(testBedList.isEmpty()) {
         return;
      }

      List<String> registerTB = new ArrayList<String>();
      for (String commandsFilePath : testBedList) {
         registerTB.add("registry-add-testbed " + commandsFilePath);
      }

      try {
         _controller.initialize(registerTB);
         _controller.prepare();
         _controller.execute();
      } catch (Exception e) {
         // There is an issue with the provided configuration - the test is skipped/
         throw new RuntimeException("Register testbed issue!", e);
      }
   }

   /**
    * Load the testbed from the legacy properties file.
    * @param localTestBedResourceFile
    */
   private void loadLocalTestBed(String localTestBedResourceFile) {
      _logger.info("Load testbed from provided file: " + localTestBedResourceFile);
      // TODO: move to controller
      testBed = createTestBed(localTestBedResourceFile);

      // Set common login credentials for the vCDe services
      // TODO: OBSOLETE will be moved once the workflow 2
      SsoUtil.setLoginCredentials(
            testBed.getVc(),
            testBed.getAdminUser().username.get(),
            testBed.getAdminUser().password.get());
   }

   /**
    * Gets the test description for a @Test method
    **
    * @param m
    *           - test method which is executed
    * @param testContext
    *           - TestNG object holding the current test suite information
    * @return The test description set on the respective test method
    */
   private String getTestDescription(Method m, ITestContext testContext) {
      String packageName = m.getDeclaringClass().getPackage().getName();
      String methodName = m.getName();
      String testDescription = methodName;
      String testId;
      // Get test description
      for (ITestNGMethod ngMethod : testContext.getAllTestMethods()) {
         String methodClass = ngMethod.getRealClass().getSimpleName();
         if (ngMethod.getRealClass().getPackage().getName().equals(packageName)
               && methodClass.equals(m.getDeclaringClass().getSimpleName())
               && ngMethod.getMethodName().equals(m.getName())) {
            testId = ngMethod.getConstructorOrMethod().getMethod()
                  .getAnnotation(TestID.class).id()[0];
            testDescription = String.format("Test ID: %s, %s", testId,
                  ngMethod.getDescription());
            break;
         }
      }
      return testDescription;
   }
}
