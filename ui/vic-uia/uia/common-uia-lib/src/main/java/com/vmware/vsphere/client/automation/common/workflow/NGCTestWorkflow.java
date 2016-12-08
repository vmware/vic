/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.workflow;

import java.util.List;

import org.apache.commons.lang.RandomStringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.annotations.BeforeSuite;
import org.testng.annotations.Optional;
import org.testng.annotations.Parameters;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.URLSpec;
import com.vmware.client.automation.common.step.ConnectBrowserStep;
import com.vmware.client.automation.common.step.LoadURLStep;
import com.vmware.client.automation.common.step.LoginStep;
import com.vmware.client.automation.servicespec.SeleniumServiceSpec;
import com.vmware.client.automation.util.SeleniumUtil;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.vsphere.client.automation.common.step.LoadVcLoginPageStep;
import com.vmware.vsphere.client.automation.provider.commontb.SeleniumNodeProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.step.CreateAdminUserByApiStep;

/**
 * Abstract test workflow implementation common for basic NGC UI tests.
 *
 * UI Base test
 * 1. Verify that the user for the test exist if not creates one.
 * 2. Connect to Selenium server.
 * 3. Load url of the vSphere client to test - provided VcSpec. The step search for vcSpec tagged with NGC_URL_TAG.
 * If no such spec is defined get the first spec found in the link of specs.
 * 4. Login using the created user in step 1.
 */
public abstract class NGCTestWorkflow extends BaseTestWorkflow {

   // Tag to mark the user used for executing the UI steps ofthe test.
   protected static final String TEST_USER_SPEC_TAG = "TEST_USER_SPEC_TAG";

   // Tag used for marking the VcSpec used for opening vSphere client.
   // It can be used to specify specific vC in case of multi -vc setup.
   protected static final String NGC_URL_TAG = "NGC_URL_TAG";

   // username template
   private static String USERNAME_TEMMPLATE = "%s@vsphere.local";

   // Default password - there is a restriction that the password should be strong.
   // Keep for now hard coded value.
   private static String DEFAULT_PASSWORD = "Admin!23";

   protected static final Logger _logger = LoggerFactory.getLogger(NGCTestWorkflow.class);

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   /**
    * This method will be invoked before any test suite. First it loads the
    * common resources properties provided by the default configuration file.
    * Than it creates connection to the Selenium server by using the parameters
    * provided in the TestNG configuration file.
    *
    * All the parameters are optional and if are not provided the default value
    * will be used.
    *
    * @param seleniumServer
    *           IP of the Selenium RC server. If no IP is provided the localhost
    *           will be used.
    * @param seleniumBrowser
    *           the browser type to be used by the Selenium RC server. If no
    *           browser type is specified the IE will be used.
    * @param seleniumBrowserArgs
    *           the browser command-line arguments (comma-separated) to be used by the Selenium RC server.
    * @param screenShotFolder
    *           path to the folder where the screenshots from the test execution
    *           will be kept. If no folder is specified the default value in
    *           TestTarget class will be used.
    * @param screenShotWebServer
    *           downloaded from AutomationNimbusService published on
    *           'nimbusTestBedServer' web server address used for storing the
    *           screen shots from the property</li> tests
    * @param seleniumProtocol
    *           whether to use Selenium Protocl - Web Driver or RC <li>
    *           The possible values are webdriver and rc. If no valid value is
    *           specified RC is protocol is used.
    * @param seleniumServerPort
    *           from 'localTestBedResourcesFile' property</li> port where
    *           Selenium runs <li>If 'nimbus' value is chosen then testBed
    *           configuration will be
    */
   @BeforeSuite(alwaysRun = true, dependsOnMethods = {"setTestBedConfiguration"})
   @Parameters({ "seleniumServer", "seleniumBrowser", "seleniumBrowserArgs", "screenShotFolder",
      "screenShotWebServer", "seleniumProtocol", "seleniumServerPort" })
   public final void setSeleniumParameters(
         @Optional("127.0.0.1") String seleniumServer,
         @Optional("*firefox") String seleniumBrowser,
         @Optional("") String seleniumBrowserArgs,
         @Optional("") String screenShotFolder,
         @Optional("") String screenShotWebServer,
         @Optional("webdriver") String seleniumProtocol,
         @Optional("4444") String seleniumServerPort) {

      // In provider workflow mode the Selenium settings are loaded through the registered providers
      if(!isProviderRun) {
         _logger.info("Use WebDriver: " + seleniumProtocol);
         _logger.warn("Selenium Server: " + seleniumServer);

         // Init default selenium service
         initDefaultSeleniumService(seleniumProtocol, seleniumBrowser, seleniumBrowserArgs, seleniumServer,
               seleniumServerPort, screenShotFolder, screenShotWebServer);
         SeleniumUtil.setSeleniumSpec(_seleniumServiceSpec);
      }
   }

   @Override
   public void initSpec() {
      BaseSpec spec = getSpec();

      // Add the BrowserSpec
      URLSpec URLSpec = new URLSpec();
      URLSpec.url.set(testBed.getNGCURL());
      spec.links.add(URLSpec);

      // Add the UserSpec
      UserSpec userSpec = testBed.getAdminUser();
      spec.links.add(userSpec);

      // TODO: fixtures have to set that value for us
      spec.links.add(SeleniumUtil.getSeleniumSpec());

   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Nothing to do

   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {
      // Connect to Selenium server
      composition.appendStep(new ConnectBrowserStep());

      // Connect Browser Step
      composition.appendStep(new LoadURLStep());

      // Login Step
      composition.appendStep(new LoginStep());
   }

   // TestWorkflow methods

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer seleniumProviderConsumer =
            testbedBridge.requestTestbed(SeleniumNodeProvider.class, false);

      EntitySpec seleniumSpec =
            seleniumProviderConsumer.getPublishedEntitySpec(
                  SeleniumNodeProvider.DEFAULT_ENTITY);
      testSpec.add(seleniumSpec);

      // This vcSpec will be used to open NGC URL
      VcSpec vcSpecToLogin = getVcSpecToLogin(testSpec.getAll(VcSpec.class));

      // Generate and tag admin user with unique username
      UserSpec ngcTestUser = generateUserSpec(vcSpecToLogin);
      testSpec.add(ngcTestUser);
   }


   @Override
   // TODO: rkovachev what about moving it to upper class?
   public void execute() throws Exception {
      if (isProviderRun) {
         _logger.info("Running provider testbed step workflow.");
         super.invokeTestExecuteCommand();
      } else {
         _logger.info("Running local testbed step workflow.");
         super.execute();
      }
   }

   @Override
   public void composePrereqSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      // Validate and created the user needed to UI test execution
//      flow.appendStep(
//            "Create admin user to run the test!",
//            new CreateAdminUserByApiStep(),
//            new String[] { TEST_USER_SPEC_TAG });
   }

   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {

      // Connect to Selenium browser
      flow.appendStep("Connect to Selenium Browser", new ConnectBrowserStep());

      // Load the vSphere client in the browser
      flow.appendStep(
            "Open vSphere URL",
            new LoadVcLoginPageStep(),
            new String[] { NGC_URL_TAG });

      // Login in the vSphere client
      flow.appendStep(
            "Login in vSphere Client",
            new LoginStep(),
            new String[] { TEST_USER_SPEC_TAG });
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Generate UserSpec with unique(pseudo unique - it assumes that random string creating is unique for our case) to
    * be used for UI steps of the scenario.
    *
    * @param vcSpec parent of the USerSpec
    * @return UserSpec
    */
   protected UserSpec generateUserSpec(VcSpec vcSpec) {
      UserSpec userSpec = new UserSpec();

      String randomUsername = RandomStringUtils.randomAlphabetic(8);
//      userSpec.username.set(String.format(USERNAME_TEMMPLATE, randomUsername));
//      userSpec.password.set(DEFAULT_PASSWORD);
      userSpec.username.set(String.format(USERNAME_TEMMPLATE, "administrator"));
      userSpec.password.set(DEFAULT_PASSWORD);
      userSpec.parent.set(vcSpec);
      userSpec.tag.set(TEST_USER_SPEC_TAG);

      return userSpec;
   }

   private void initDefaultSeleniumService(String seleniumProtocol, String seleniumBrowser, String seleniumBrowserArgs,
         String seleniumServer, String seleniumServerPort, String screenShotFolder,
         String screenShotWebServer) {

      // TODO: move to controller
      _seleniumServiceSpec = new SeleniumServiceSpec();
      _seleniumServiceSpec.seleniumBrowser.set(seleniumBrowser);
      _seleniumServiceSpec.seleniumBrowserArgs.set(seleniumBrowserArgs);
      _seleniumServiceSpec.seleniumOS.set("Windows");
      _seleniumServiceSpec.seleniumServer.set(seleniumServer);
      _seleniumServiceSpec.seleniumServerPort.set(seleniumServerPort);
      // TODO rkovachev: fix the workaround
      if(seleniumProtocol.equalsIgnoreCase("true")) {
        seleniumProtocol = "webdriver";
      }
      _seleniumServiceSpec.seleniumProtocol.set(seleniumProtocol);
      _seleniumServiceSpec.screenShotFolder.set(screenShotFolder);
      _seleniumServiceSpec.screenShotWebServer.set(screenShotWebServer);
   }

   /**
    * Iterates the list of VcSpec to find the one tagged for opening the vSphere URL.
    * If no VcSpec is tagged - get the first one.
    * @param vcList list of VcSpec linked in the TestSpec
    * @return VcSpec
    */
   private VcSpec getVcSpecToLogin(List<VcSpec> vcList) {
      VcSpec vcToLogin = null;
      boolean isLoginVcTagged = false;

      if (vcList.size() == 0) {
         throw new IllegalArgumentException(
               "VcSpec spec is mandatory for each vSphere client test.");
      }

      for (VcSpec vcSpec : vcList) {
         if (vcSpec.tag.isAssigned() && vcSpec.tag.getAll().contains(NGC_URL_TAG)) {
            isLoginVcTagged = true;
            vcToLogin = vcSpec;
            break;
         }
      }

      if (!isLoginVcTagged) {
         vcToLogin = vcList.get(0);
         List<String> tagList = vcToLogin.tag.getAll();
         tagList.add(NGC_URL_TAG);
         vcToLogin.tag.set(tagList);
      }
      return vcToLogin;
   }
}
