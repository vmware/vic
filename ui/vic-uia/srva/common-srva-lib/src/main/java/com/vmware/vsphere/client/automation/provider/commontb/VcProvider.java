/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.connector.VcConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * The class provides VC which is available in the system no matter if it is
 * physical or Nimbus deployed.
 */
public class VcProvider implements ProviderWorkflow {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.vc.entity.default";

   // Testbed publisher settings

   // logger
   private static final Logger _logger = LoggerFactory.getLogger(VcProvider.class);

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      VcSpec vcSpec = new VcSpec();
      // Set empty service spec
      vcSpec.service.set(new VcServiceSpec());
      publisherSpec.links.add(vcSpec);
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, vcSpec);
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec,
         TestBedBridge testbedBridge) throws Exception {

      // Request vc spec
      TestbedSpecConsumer vcProviderConsumer =
            testbedBridge.requestTestbed(BaseVcProvider.class, false);
      VcSpec requestedVcSpec =
            vcProviderConsumer.getPublishedEntitySpec(BaseVcProvider.DEFAULT_ENTITY);

      assemblerSpec.add(requestedVcSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      VcSpec vcSpec = publisherSpec.links.get(VcSpec.class);

      vcSpec.vscUrl.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            BaseVcProvider.TESTBED_KEY_VSC_URL));

      String name =
            SettingsUtil.getRequiredValue(testbedSettings, BaseVcProvider.TESTBED_KEY_NAME);
      String endpoint =
            SettingsUtil.getRequiredValue(testbedSettings, BaseVcProvider.TESTBED_KEY_ENDPOINT);
      String username =
            SettingsUtil.getRequiredValue(testbedSettings, BaseVcProvider.TESTBED_KEY_USERNAME);
      String password =
            SettingsUtil.getRequiredValue(testbedSettings, BaseVcProvider.TESTBED_KEY_PASSWORD);

      vcSpec.ssoLoginUsername.set(username);
      vcSpec.ssoLoginPassword.set(password);

      vcSpec.service.get().endpoint.set(endpoint);
      vcSpec.service.get().username.set(username);
      vcSpec.service.get().password.set(password);

      vcSpec.name.set(name);

      _logger.info("Loaded publisherSpec: " + vcSpec.toString());
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.info("Nothing to do here - the resource is requested!");
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) {
      // SimulatorConnectorsFactory.createAndSetConnectors(serviceConnectorsMap);
      // TODO: rkovachev - move ot ConnectionFactory like the simulator sample
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
   public void composeProviderSteps(
         WorkflowStepsSequence<? extends WorkflowStepContext> flow)
         throws Exception {
      // save the settings for the base components of the TB
      flow.appendStep("Save Settings from base testbeds", new ProviderWorkflowStep() {

         private VcSpec _vcSpec;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filterAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) throws Exception {

            // collect data to be saved if in assemble mode
            if (isAssembling) {
               _vcSpec = filterAssemblerSpec.links.get(VcSpec.class);
            } else {
               _vcSpec = filteredPublisherSpec.links.get(VcSpec.class);
            }
         }

         @Override
         public void disassemble() throws Exception {
            // Nothing to disassemble
         }

         @Override
         public boolean checkHealth() throws Exception {
            // Nothing to check
            return true;
         }

         @Override
         public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
            String vcHostName = VcServiceUtil.getVcHostname(_vcSpec.service.get());

            testbedSettingsWriter.setSetting(BaseVcProvider.TESTBED_KEY_NAME, vcHostName);
            testbedSettingsWriter.setSetting(BaseVcProvider.TESTBED_KEY_ENDPOINT, _vcSpec.service.get().endpoint.get());
            testbedSettingsWriter.setSetting(BaseVcProvider.TESTBED_KEY_USERNAME, _vcSpec.ssoLoginUsername.get());
            testbedSettingsWriter.setSetting(BaseVcProvider.TESTBED_KEY_PASSWORD, _vcSpec.ssoLoginPassword.get());
            testbedSettingsWriter.setSetting(BaseVcProvider.TESTBED_KEY_VSC_URL, _vcSpec.vscUrl.get());
         }

      });

   }

   @Override
   public int providerWeight() {
      // TODO Auto-generated method stub
      return 0;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return BaseVcProvider.class;
   }

}
