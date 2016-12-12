package com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock;

import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProvider.MockProviderSpec;

/**
 * Mock implementation of {@link TestBedBridge} providing all of the entities
 * from infinite amount of {@link MockProvider} instances
 */
public class MockProviderTestBedBridge implements TestBedBridge {

   private int providerRequestCounter = 0;

   public int getTotalRequestsForProvider() {
      return providerRequestCounter;
   }

   @Override
   public TestbedSpecConsumer requestTestbed(
         Class<? extends ProviderWorkflow> testbedProviderClass,
         boolean isShared) {

      MockTestbedSpecConsumer result = new MockTestbedSpecConsumer(
            new MockProviderSpec(MockProvider.ENTITY_1), new MockProviderSpec(
                  MockProvider.ENTITY_2), new MockProviderSpec(
                  MockProvider.ENTITY_3));
      providerRequestCounter++;
      return result;
   }

}
