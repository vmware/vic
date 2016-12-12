/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.SeleniumServiceSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.connector.SeleniumConnector;

/**
 * Selenium node provider is used to provide, register and disassemble Selenium
 * server used for execution of a UI tests.
 * NOTE: The current implementation does not provide assemble operation.
 * TODO: rkovachev implement assemble and disassemble commands for the provider.
 */
public class SeleniumNodeProvider implements ProviderWorkflow {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.selenium.node.default";

   // Selenium service info
   // Browser type
   private static final String SELENIUM_BROWSER = "testbed.seleniumBrowser";
   // (Optional) Browser command line arguments
   private static final String SELENIUM_BROWSER_ARGS = "testbed.seleniumBrowserArgs";
   // Selenium protocol - RC or WebDriver
   private static final String SELENIUM_PROTOCOL = "testbed.seleniumProtocol";
   // IP of the selenium Server
   private static final String SELENIUM_SERVER = "testbed.seleniumServer";
   private static final String SELENIUM_SERVER_PORT = "testbed.seleniumServerPort";

   // TODO rkovachev: For WebDriver these settings should not be used
   private static final String SELENIUM_SCREENSHOT_FOLDER = "testbed.screenShotFolder";
   private static final String SELENIUM_SCREENSHOT_WEB_SERVER = "testbed.screenShotWebServer";


   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      EntitySpec seleniumNode = new EntitySpec();
      publisherSpec.links.add(seleniumNode);
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, seleniumNode);
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      EntitySpec seleniumNode = publisherSpec.links.get(EntitySpec.class);

      SeleniumServiceSpec seleniumServiceSpec = new SeleniumServiceSpec();
      seleniumServiceSpec.seleniumServer.set(testbedSettings.getSetting(SELENIUM_SERVER));

      seleniumServiceSpec.seleniumBrowser.set(testbedSettings.getSetting(SELENIUM_BROWSER));
      if (!Strings.isNullOrEmpty(testbedSettings.getSetting(SELENIUM_BROWSER_ARGS))) {
         seleniumServiceSpec.seleniumBrowserArgs.set(testbedSettings.getSetting(SELENIUM_BROWSER_ARGS));
      }
      seleniumServiceSpec.seleniumServerPort.set(testbedSettings.getSetting(SELENIUM_SERVER_PORT));
      seleniumServiceSpec.seleniumProtocol.set(testbedSettings.getSetting(SELENIUM_PROTOCOL));

      seleniumServiceSpec.screenShotFolder.set(testbedSettings.getSetting(SELENIUM_SCREENSHOT_FOLDER));
      seleniumServiceSpec.screenShotWebServer.set(testbedSettings.getSetting(SELENIUM_SCREENSHOT_WEB_SERVER));

      seleniumNode.service.set(seleniumServiceSpec);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) throws Exception {
      // TODO Auto-generated method stub
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof SeleniumServiceSpec) {
            serviceConnectorsMap.put(serviceSpec, new SeleniumConnector(
                  (SeleniumServiceSpec) serviceSpec));

         }
      }
   }


   @Override
   public int providerWeight() {
      return 1;
   }

   @Override
   public void composeProviderSteps(
         WorkflowStepsSequence<? extends WorkflowStepContext> flow)
         throws Exception {
      // TODO Auto-generated method stub
      
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return this.getClass();
   }

}
