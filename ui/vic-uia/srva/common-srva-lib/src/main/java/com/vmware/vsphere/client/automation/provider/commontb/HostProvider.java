/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.servicespec.HostServiceSpec;
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
import com.vmware.vsphere.client.automation.provider.connector.HostConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostClientBasicSrvApi;
import com.vmware.client.automation.delay.Sleep;

/**
 * The provided host by this class makes check that vmotion is enabled
 * on the host and enables it if not.
 * It should be used(requested) by the tests/providers when the underlying
 * infrastructure used for the deployed resource is not important for
 * the respective test/provider.
 */
public class HostProvider implements ProviderWorkflow {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.host.entity.default";

   // local datastore entity
   public static final String LOCAL_DS_ENTITY = "provider.host.entity.datastore";

   // Testbed publisher settings
   private static final String TESTBED_KEY_DATASTORE =  "testbed.local.datastore.name";

   // logger
   private static final Logger _logger = LoggerFactory.getLogger(HostProvider.class);

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      HostSpec hostSpec = new HostSpec();
      publisherSpec.links.add(hostSpec);
      hostSpec.service.set(new HostServiceSpec());
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, hostSpec);

      DatastoreSpec dsSpec = new DatastoreSpec();
      publisherSpec.links.add(dsSpec);
      hostSpec.service.set(new HostServiceSpec());
      publisherSpec.publishEntitySpec(LOCAL_DS_ENTITY, dsSpec);
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec,
         TestBedBridge testbedBridge) throws Exception {

      // Request host spec for the clustered host
      TestbedSpecConsumer hostProviderConsumer =
            testbedBridge.requestTestbed(BaseHostProvider.class, false);
      HostSpec requestedHostSpec =
            hostProviderConsumer.getPublishedEntitySpec(BaseHostProvider.DEFAULT_ENTITY);

      DatastoreSpec localDatastoreSpec = SpecFactory.getSpec(DatastoreSpec.class, requestedHostSpec);
      localDatastoreSpec.type.set(DatastoreType.VMFS);

      assemblerSpec.add(requestedHostSpec, localDatastoreSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.info("Start HostProvider assign published specs");
      HostSpec hostSpec = (HostSpec)publisherSpec.getPublishedEntitySpec(
            DEFAULT_ENTITY);

      String endpoint =
            SettingsUtil.getRequiredValue(testbedSettings, BaseHostProvider.TESTBED_KEY_ENDPOINT);
      String username =
            SettingsUtil.getRequiredValue(testbedSettings, BaseHostProvider.TESTBED_KEY_USERNAME);
      String password =
            SettingsUtil.getRequiredValue(testbedSettings, BaseHostProvider.TESTBED_KEY_PASSWORD);
      String servicePort =
            SettingsUtil.getRequiredValue(testbedSettings, BaseHostProvider.TESTBED_KEY_SERVICE_PORT);
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

      DatastoreSpec dsSpec = (DatastoreSpec)publisherSpec.getPublishedEntitySpec(LOCAL_DS_ENTITY);
      dsSpec.service.set(hostSeviceSpec);
      dsSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_DATASTORE));
      dsSpec.parent.set(hostSpec);

      _logger.info("Loaded publisherSpecs: ");
      _logger.info("Loaded Host spec: " + hostSpec.toString());
      _logger.info("Loaded Datastore spec: " + dsSpec.toString());
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.debug("No assembler spec to laod.");

   }

   @Override
   public void assignTestbedConnectors(Map<ServiceSpec,
         TestbedConnector> serviceConnectorsMap) throws Exception {
      // SimulatorConnectorsFactory.createAndSetConnectors(serviceConnectorsMap);
      // TODO: rkovachev - move to ConnectionFactory like the simulator sample
      _logger.info("Host spec provider assignTestbedConnectors");
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof HostServiceSpec) {
            try {
               serviceConnectorsMap.put(serviceSpec, new HostConnector(
                     (HostServiceSpec) serviceSpec));
            } catch (SsoException e) {
               e.printStackTrace();
            }
         }
      }
   }

   @Override
   public void composeProviderSteps(WorkflowStepsSequence<? extends WorkflowStepContext> flow)
         throws Exception {

      // save the settings for the base components of the TB
      flow.appendStep("Save Settings from base testbeds", new ProviderWorkflowStep() {

         private HostSpec _hostSpec;
         private DatastoreSpec _datastoreSpec;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filterAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) throws Exception {

            // collect data to be saved if in assemble mode
            if (isAssembling) {
               _hostSpec = filterAssemblerSpec.links.get(HostSpec.class);
               _datastoreSpec = filterAssemblerSpec.links.get(DatastoreSpec.class);
            } else {
               _hostSpec = filteredPublisherSpec.links.get(HostSpec.class);
               _datastoreSpec = filteredPublisherSpec.links.get(DatastoreSpec.class);
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

            testbedSettingsWriter.setSetting(BaseHostProvider.TESTBED_KEY_ENDPOINT, _hostSpec.name.get());
            testbedSettingsWriter.setSetting(BaseHostProvider.TESTBED_KEY_USERNAME, _hostSpec.userName.get());
            testbedSettingsWriter.setSetting(BaseHostProvider.TESTBED_KEY_PASSWORD, _hostSpec.password.get());
            testbedSettingsWriter.setSetting(BaseHostProvider.TESTBED_KEY_SERVICE_PORT, _hostSpec.port.get() + "");

            testbedSettingsWriter.setSetting(TESTBED_KEY_DATASTORE, _datastoreSpec.name.get());
         }

      });


      flow.appendStep("Verify vmotion is enabled on the host.", new ProviderWorkflowStep() {

         private HostSpec _hostSpec;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filterAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) throws Exception {

            // collect data to be saved if in assemble mode
            if (isAssembling) {
               _hostSpec = filterAssemblerSpec.get(HostSpec.class);
            } else {
               _hostSpec = filteredPublisherSpec.links.get(HostSpec.class);
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
            HostClientBasicSrvApi.getInstance().enableVmotion(_hostSpec);
         }

      });

      flow.appendStep("Re-create the local datastore.", new ProviderWorkflowStep() {

         private HostSpec _hostSpec;
         private DatastoreSpec _datastoreSpec;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filterAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) throws Exception {

            // collect data to be saved if in assemble mode
            if (isAssembling) {
               _hostSpec =  filterAssemblerSpec.get(HostSpec.class);
               _datastoreSpec = filterAssemblerSpec.get(DatastoreSpec.class);
            } else {
               _hostSpec = filteredPublisherSpec.links.get(HostSpec.class);
               _datastoreSpec = filteredPublisherSpec.links.get(DatastoreSpec.class);
            }
         }

         @Override
         public void disassemble() throws Exception {
            // Nothing to disassemble
         }

         @Override
         public boolean checkHealth() throws Exception {
            return HostClientBasicSrvApi.getInstance().findVmfsDatastoreByName(_datastoreSpec);
         }

         @Override
         public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
            // TODO: Remove this sleep/while or file a bug
            // We have to sleep here because if we try to delete the
            // datastore we see error that it has been used. This may be
            // because of PR1418588 but I can not tell for sure.
            boolean datastoreDeleted = HostClientBasicSrvApi
                  .getInstance().deleteAllVmfsDatastores(_hostSpec);
            int retry = 20;
            while (retry-- > 0 && !datastoreDeleted) {
               _logger
                     .info("Wait for 10 seconds before trying again to delete the datastore. Left retries: "
                           + retry);
               new Sleep(10000).consume();
               datastoreDeleted = HostClientBasicSrvApi.getInstance()
                     .deleteAllVmfsDatastores(_hostSpec);
            }

            if(!datastoreDeleted){
               _logger.error("Unable to deleted VMFS datastore");
               throw new Exception("Unable to deleted VMFS datastore");
            }
            HostClientBasicSrvApi.getInstance().createVmfsDatastore(_datastoreSpec);
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
      return BaseHostProvider.class;
   }


}
