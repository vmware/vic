package com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock;

import java.util.HashMap;
import java.util.Map;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProvider.MockProviderSpec;

/**
 * Mock implementation of {@link TestbedSpecConsumer} providing with a
 * predefined set of {@link MockProviderSpec} results
 */
public class MockTestbedSpecConsumer implements TestbedSpecConsumer {

   private final Map<String, MockProviderSpec> results = new HashMap<String, MockProviderSpec>();

   /**
    * Initializes new instance of {@link MockTestbedSpecConsumer}
    *
    * @param result
    *           the results for {@link #getPublishedEntitySpec(String)}
    */
   @SafeVarargs
   public MockTestbedSpecConsumer(MockProviderSpec... result) {
      for (MockProviderSpec mockedProviderSpec : result) {
         results.put(mockedProviderSpec.entityId, mockedProviderSpec);
      }

   }

   @SuppressWarnings("unchecked")
   @Override
   public <T extends EntitySpec> T getPublishedEntitySpec(String entitySpecId) {
      return (T) results.get(entitySpecId);
   }

}
