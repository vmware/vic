/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.workflow.common.WorkflowContext;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowContext;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;
import com.vmware.client.automation.workflow.provider.PublisherSpec;

/**
 * Implementation of TestBedBridge. ass an addition it provides also methods
 * for used by the stages of the Prvider and Test Workflow controllers - claim
 * tesbted, release testbed and etc.
 */
public class ProviderBridgeImpl implements TestBedBridge {
   private final WorkflowContext _consumerContext;

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(ProviderBridgeImpl.class);

   // Need the caller context workflow
   public ProviderBridgeImpl(WorkflowContext consumerContext) {
      _consumerContext = consumerContext;
   }

   /**
    * Request a testbed specified by the provider class parameter.
    * The method returns the publisher spec of the requested provider intialized.
    * The publisher spec is not populated with actual data on that stage - it is
    * only registered that resources is needed by the Test/Provider.
    *
    * NOTE: The shared test bed support is not implemented!
    *
    * Internal Workflow
    * 1. A providerContext of the requested TB is instantiated.
    * 2. The providerContext instance is mapped to the _consumerContext.
    * 3. ProviderController is instantiated for the provider context and return getPublisherSpec()
    */
   @Override
   public TestbedSpecConsumer requestTestbed(
         Class<? extends ProviderWorkflow> testbedProviderClass, boolean isShared) {

      if (testbedProviderClass == null) {
         throw new IllegalArgumentException(
               "Required test bed provider class is not set");
      }

      // 0. Check the class through the registry
      if (!getRegistry().checkWorkflowExists(
            testbedProviderClass.getCanonicalName(),
            ProviderWorkflow.class)) {
         throw new IllegalArgumentException("Test bed provider class not found");
      }

      // 1. Instantiate it // This is moved in the registry
      ProviderWorkflowContext providerContext =
            getRegistry().registrerWorkflowContext(testbedProviderClass);

      // 2. Map provider with consumer context
      getRegistry().mapProviderToConsumer(_consumerContext, providerContext);

      // TODO rkovachev: Why instead of the 3. and 4.
      /**
       providerContext.getProviderWorkflow().initPublisherSpec(providerContext.getPublisherSpec());
       providerContext.getPublisherSpec();
       */
      // 3. Use the retrieved context to instantiate a controller
      ProviderWorkflowController controller =
            ProviderWorkflowController.create(providerContext, null);

      // 4. Call its getProvider method
      PublisherSpec publisherSpec = null;
      try {
         publisherSpec = controller.getPublisherSpec();
      } catch (Exception e) {
         throw new RuntimeException("Requesting provider spec failed.");
      }
      return publisherSpec;


      // Maps:
      // - test context 1 : n shared provider specs, hosted in provider contexts				->
      // - test context 1 : n dedicated provider specs, hosted in provider contexts				-> consumerContextProviderContextMap (Map<WorkflowContext, Set<ProviderWorkflowContext>>)
      // - provider context 1 : n dedicated provider specs, hosted in provider contexts	->
      // - provider spec n (e.g. provider context) : 1 test bed (context) 		-> providerContextTestbedContextMap (Map<ProviderWorkflowContext, TestbedContext>
      // - provider (class) 1 : n test bed (context) -> providerClassTestbedContextMap (Map<Class<ProviderWorkflow>, List<TestbedContext>>)

      // Query:
      // - Find my provider instances
      // - Find free test bed

      // Big question: Should we keep the map externally or internally to the context?
      // Big answer: Let's keep this externally. It's more flexible to quick strategy changes of workflow implementations and others.
   }

   /**
    * Release allocated testbeds for the test/provider _consumerContext.
    * No need to free connections here as they are kept in the provider controller.
    */
   public void releaseTestbeds() {
      // Connections are not in the registry - they are in the provider
      Set<ProviderWorkflowContext> providerContextSet =
            getRegistry().getRegisteredConsumerProviders(_consumerContext);
      if (providerContextSet.isEmpty()) {
         return;
      }

      // Free allocated test bed for each provider context
      for (ProviderWorkflowContext providerContext : providerContextSet) {
         try {
            // TODO rkovachev: Calling check health before clean???
            getRegistry().freeTestbed(providerContext);
         } catch (Exception e) {
            e.printStackTrace();
         }
      }
   }

   /**
    * Claims testbed for each requested resource. Only if for all the requests can be
    * allocated testbed contexts - assignTesbed is invoked for each ProviderWorkflow.
    * @throws ProviderWorkflowException
    * @throws ProviderControllerException
    */
   public void claimTestbeds() throws ProviderControllerException, ProviderWorkflowException {
      // Note: Only first level provider specs are assigned to test bests.
      // The assignment could work in all or non manner

      // Get all top level tests beds
      Set<ProviderWorkflowContext> providerContextSet =
            getRegistry().getRegisteredConsumerProviders(_consumerContext);
      if (providerContextSet.isEmpty()) {
         return;
      }

      // Discover free test bed for each provider context and map it
      for (ProviderWorkflowContext providerContext : providerContextSet) {
         // Claim
         getRegistry().claimTestbed(providerContext);
      }

      // Assign provided resources
      for (ProviderWorkflowContext providerContext : providerContextSet) {
         // Assign
         ProviderWorkflowController controller = ProviderWorkflowController
               .create(providerContext,
                     getRegistry().getProviderTestbed(providerContext)
                           .getTestbedSettingsFilePath());
         controller.assignTestBed();
      }

      // Note: Dedicated vs shared test bed support is not handles here.
   }

   /**
    * Establish connections for the assigned providers.
    * @throws ProviderWorkflowException
    * @throws ProviderControllerException
    */
   public void connectToTestbeds() throws ProviderControllerException, ProviderWorkflowException {
      // Get all top level tests beds
      Set<ProviderWorkflowContext> providerContextSet =
            getRegistry().getRegisteredConsumerProviders(_consumerContext);
      if (providerContextSet.isEmpty()) {
         return;
      }

      for (ProviderWorkflowContext providerContext : providerContextSet) {
         ProviderWorkflowController controller = ProviderWorkflowController
               .create(providerContext,
                     getRegistry().getProviderTestbed(providerContext)
                           .getTestbedSettingsFilePath());
         controller.connect();
      }
   }

   /**
    * TODO rkovachev - Document it
    */
   public void analyze() {
      Set<ProviderWorkflowContext> providerContextSet =
            getRegistry().getRegisteredConsumerProviders(_consumerContext);
      if (providerContextSet.isEmpty()) {
         return;
      }

      for (ProviderWorkflowContext providerContext : providerContextSet) {
         try {
            // Assign
            ProviderWorkflowController controller =
                  ProviderWorkflowController.create(providerContext, null);
            controller.analyze();
         } catch (Exception e) {
            break;
         }
      }
   }

   // Private methods

   private WorkflowRegistry getRegistry() {
      return WorkflowRegistry.getRegistry();
   }
}
