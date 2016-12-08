/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.BaseElementalProvider;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.connector.VcConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * Base VC provider class that provides the ability to assign and publish
 * already deployed VC.
 * NOTE: The base provider is no able to assemble resources. For assembling
 * use one of the providers that extends it and is able to deploy the vC
 * on the specified infrastructure - i.e. NimbusVcProvider.
 */
public class BaseVcProvider extends BaseElementalProvider {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.vc.entity.default";

   // logger
   private static final Logger _logger = LoggerFactory.getLogger(BaseVcProvider.class);

   // Testbed publisher settings
   protected static final String TESTBED_KEY_NAME = "testbed.name";
   protected static final String TESTBED_KEY_ENDPOINT = "testbed.endpoint";
   protected static final String TESTBED_KEY_VSC_URL = "testbed.vsc.url";
   protected static final String TESTBED_KEY_USERNAME = "testbed.user";
   protected static final String TESTBED_KEY_PASSWORD = "testbed.pass";

   protected static final String TESTBED_ASSEMBLER_KEY_VC_IP = "testbed.vc.ip";
   protected static final String TESTBED_ASSEMBLER_KEY_VC_USER = "testbed.vc.user";
   protected static final String TESTBED_ASSEMBLER_KEY_VC_PASSWORD = "testbed.vc.pass";


   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      VcSpec vcSpec = new VcSpec();
      // Set empty service spec
      vcSpec.service.set(new VcServiceSpec());
      publisherSpec.links.add(vcSpec);
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, vcSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      VcSpec vcSpec = publisherSpec.links.get(VcSpec.class);

      vcSpec.vscUrl.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_VSC_URL));

      String name =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_NAME);
      String endpoint =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_ENDPOINT);
      String username =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_USERNAME);
      String password =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_PASSWORD);

      vcSpec.ssoLoginUsername.set(username);
      vcSpec.ssoLoginPassword.set(password);

      vcSpec.service.get().endpoint.set(endpoint);
      vcSpec.service.get().username.set(username);
      vcSpec.service.get().password.set(password);

      vcSpec.name.set(name);

      _logger.info("Loaded publisherSpec: " + vcSpec.toString());
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) {
      // SimulatorConnectorsFactory.createAndSetConnectors(serviceConnectorsMap);
      // TODO: rkovachev - move to ConnectionFactory like the simulator sample
      _logger.info("Virtual Center Provider assignTestbedConnectors");
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof VcServiceSpec) {
            try {
               serviceConnectorsMap.put(serviceSpec, new VcConnector(
                     (VcServiceSpec) serviceSpec));
            } catch (SsoException e) {
               // TODO Auto-generated catch block
               e.printStackTrace();
            }
         }
      }
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Retrieve VC as resource is not yet implemeted");
   }

   @Override
   public int providerWeight() {
      return 4;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return this.getClass();
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.warn("Nothing to do here. The BaseVcProvider can not assemble just publish!");
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {
      throw new RuntimeException("Nothing to do here. The BaseVcProvider can not assemble just publish!");
   }

   @Override
   public String determineResourceVersion() throws Exception {
      throw new RuntimeException("Nothing to do here. The BaseVcProvider can not assemble just publish!");
   }


   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter)
         throws Exception {
      throw new RuntimeException("Nothing to do here. The BaseVcProvider can not assemble just publish!");
   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      throw new RuntimeException("Nothing to do here. The BaseVcProvider can not assemble just publish!");
   }

   @Override
   public void destroyTestbed() throws Exception {
      throw new RuntimeException("Nothing to do here. The BaseVcProvider can not assemble just publish!");
   }

}
