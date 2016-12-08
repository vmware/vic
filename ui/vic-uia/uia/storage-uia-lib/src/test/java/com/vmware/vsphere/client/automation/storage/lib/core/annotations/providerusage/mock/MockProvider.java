package com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock;

import java.util.Map;

import org.apache.commons.lang.NotImplementedException;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.common.annotations.MockedSpec;

/**
 * Class "implementing" the {@link ProviderWorkflow} Please note: This class
 * does not provide mock implementation of ProviderWorkflow. It is need only to
 * be used as "key" when using the {@link MockProviderTestBedBridge}
 *
 * See {@link MockProviderTestBedBridge}
 */
public class MockProvider implements ProviderWorkflow {

   /**
    * {@link MockedSpec} implementation providing entity id
    */
   public static class MockProviderSpec extends MockedSpec {

      /**
       * The entity id with which the provider is publishing the spec
       */
      public final String entityId;

      /**
       * Initializes new instance of {@link MockProviderSpec}
       *
       * @param entityId
       *           the entity id with which the entity is published by the
       *           provider
       */
      public MockProviderSpec(String entityId) {
         this.entityId = entityId;
      }

   }

   public static final String ENTITY_1 = "mock.entity.1";
   public static final String ENTITY_2 = "mock.entity.2";
   public static final String ENTITY_3 = "mock.entity.3";

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec,
         TestBedBridge testbedBridge) throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap)
         throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void composeProviderSteps(
         WorkflowStepsSequence<? extends WorkflowStepContext> flow)
         throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public int providerWeight() {
      throw new NotImplementedException();
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      throw new NotImplementedException();
   }

}
