/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.HostServiceSpec;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.BaseElementalProvider;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.connector.HostConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * Base host provider class that provides the ability to assign and publish
 * already deployed Host/ESX.
 * NOTE: The base provider is no able to assemble hosts. For assembling
 * use one of the providers that extends it and is able to deploy the host
 * on the specified infrastructure - i.e. NimbusHostProvider.
 */
public class BaseHostProvider extends BaseElementalProvider {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.host.entity.default";

   // Testbed publisher settings
   protected static final String TESTBED_KEY_ENDPOINT = "testbed.endpoint";
   protected static final String TESTBED_KEY_USERNAME = "testbed.user";
   protected static final String TESTBED_KEY_PASSWORD = "testbed.pass";
   protected static final String TESTBED_KEY_SERVICE_PORT = "testbed.service.port";

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(BaseHostProvider.class);

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      HostSpec hostSpec = new HostSpec();
      publisherSpec.links.add(hostSpec);
      hostSpec.service.set(new HostServiceSpec());
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, hostSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      _logger.info("Start HostProvider assign published specs");
      HostSpec hostSpec = (HostSpec)publisherSpec.getPublishedEntitySpec(
            DEFAULT_ENTITY);

      String endpoint =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_ENDPOINT);
      String username =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_USERNAME);
      String password =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_PASSWORD);

      String servicePort =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_SERVICE_PORT);
      hostSpec.port.set(new Integer(servicePort));

      HostServiceSpec hostSeviceSpec = new HostServiceSpec();
      hostSeviceSpec.endpoint.set(endpoint);
      hostSeviceSpec.username.set(username);
      hostSeviceSpec.password.set(password);
      hostSeviceSpec.isHostClient.set(Boolean.TRUE);

      hostSpec.service.set(hostSeviceSpec);
      hostSpec.userName.set(username);
      hostSpec.password.set(password);
      hostSpec.name.set(endpoint);

      _logger.info("Loaded publisherSpec: " + hostSpec.toString());
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) throws Exception {
      // SimulatorConnectorsFactory.createAndSetConnectors(serviceConnectorsMap);
      // TODO: rkovachev - move to ConnectionFactory like the simulator sample
      _logger.info("Host provider assignTestbedConnectors");
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof HostServiceSpec) {
            // TODO: rkovachev provide HostConnector constructor parameter - see
            // the HostConnector class.
            serviceConnectorsMap.put(serviceSpec, new HostConnector((HostServiceSpec) serviceSpec));
         }
      }
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Retrieve host as resource is not yet implemeted");
   }

   @Override
   public int providerWeight() {
      return 1;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return this.getClass();
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.warn("Nothing to do here. The HostProvider can not assemble just publish!");

   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {
      throw new RuntimeException("Use NimbusHostProvider or implement me!");

   }

   @Override
   public String determineResourceVersion() throws Exception {
      throw new RuntimeException("Use NimbusHostProvider or implement me!");
   }


   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter)
         throws Exception {
      throw new RuntimeException("Use NimbusHostProvider or implement me!");

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      throw new RuntimeException("Use NimbusHostProvider or implement me!");
   }

   @Override
   public void destroyTestbed() throws Exception {
      throw new RuntimeException("Use NimbusHostProvider or implement me!");
   }
}
