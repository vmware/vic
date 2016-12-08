/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.Set;

import org.apache.commons.collections4.CollectionUtils;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.workflow.common.Workflow;
import com.vmware.client.automation.workflow.common.WorkflowContext;
import com.vmware.client.automation.workflow.common.WorkflowContextFactory;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowContext;

/**
 * Singleton container for session configurations and for registered and
 * allocated testbed, workflow contexts and connections.
 */
public class WorkflowRegistry implements RegistryInitializationBridge {

   private static WorkflowRegistry _registry;

   private Set<Class<? extends Workflow>> _workflowClassSet;

   private Map<String, Class<? extends Workflow>> _workflowNameWorkflowClassMap;
   // Set of provider contexts for a consumer(Test/Provider)
   private Map<WorkflowContext, Set<ProviderWorkflowContext>> _consumerContextProviderContextMap;
   // Already allocated testbeds
   private Map<ProviderWorkflowContext, TestBedContext> _providerContextTestbedContextMap;
   // Map of registered testbeds
   private Map<Class<? extends ProviderWorkflow>, List<TestBedContext>> _providerClassTestbedContextMap;


   // Maps:
   // - test context 1 : n shared provider specs, hosted in provider contexts				->
   // - test context 1 : n dedicated provider specs, hosted in provider contexts				-> consumerContextProviderContextMap (Map<WorkflowContext, List<ProviderWorkflowContext>>)
   // - provider context 1 : n dedicated provider specs, hosted in provider contexts	->
   // - provider spec n (e.g. provider context) : 1 test bed (context) 		-> providerContextTestbedContextMap (Map<ProviderWorkflowContext, TestbedContext>
   // - provider (class) 1 : n test bed (context) -> providerClassTestbedContextMap (Map<Class<ProviderWorkflow>, List<TestbedContext>>)

   private SettingsRWImpl _sessionSettings;

   /**
    * Protected constructor to enforce singleton access.
    */
   protected WorkflowRegistry() {
      initRegistry();
   }

   /**
    * Get the workflows registry.
    *
    * Singleton. Single threaded.
    *
    * @return
    */
   public static WorkflowRegistry getRegistry() {
      if (_registry == null) {
         _registry = new WorkflowRegistry();
      }
      return _registry;
   }

   // Implementation of the WorkflowRegistry interface

   @Override
   public void setSessionSettings(Properties settings) {
      if (settings != null) {
         _sessionSettings.clear();
         _sessionSettings.putAll(settings);
      } else {
         throw new IllegalArgumentException("Uninitialized settings argument passed.");
      }
   }

   /**
    * Add a testbed to the registry.
    * NOTE: When the provider workflow is BaseElementalProvider it may be added as both the
    * provider specific class and as the provider base type.
    * For example PhysicalHostProvider will be added as both PhysicalHostProvider and as HostProvider.
    */
   @Override
   public void addTestBed(ProviderWorkflow providerWorkflow,
         String testbedFilePath, Properties testbedSettings) {

      //Testbeds are stored in a list mapped by the ProviderWorkflow class they represent.
      TestBedContext newTestBed = new TestBedContext(testbedFilePath, testbedSettings);

      Class<? extends ProviderWorkflow> providerClass = providerWorkflow.getClass();
      addTestbed(providerClass, newTestBed, testbedSettings);
      // When the provider is a BaseElementalProvider register it is as the base
      // type resource too.
      if(ProviderWorkflow.class.isAssignableFrom(providerClass)) {
         Class<? extends ProviderWorkflow> baseProviderType = providerWorkflow.getProviderBaseType();
         if(baseProviderType != null && !baseProviderType.equals(providerClass)) {
            addTestbed(baseProviderType, newTestBed, testbedSettings);
         }
      }
   }

   /**
    * Init and register a WorkFlowContext class in the Registry -
    * once registered the workflow link to it the dependent resources
    * (providers contexts).
    * @param workflowClass
    * @return WorkfloweContext
    */
   public <WC extends WorkflowContext> WC registrerWorkflowContext(
         Class<? extends Workflow> workflowClass) {
      WorkflowContextFactory factory = new WorkflowContextFactory();
      WC workflowContext = factory.createContext(workflowClass);

      _consumerContextProviderContextMap.put(
            workflowContext,
            new HashSet<ProviderWorkflowContext>());

      return workflowContext;
   }

   /**
    * Unregister the workflow context from the registry.
    */
   public void unregistrerWorkflowContext(WorkflowContext context) {
      if (_providerContextTestbedContextMap.containsKey(context)) {
         _providerContextTestbedContextMap.remove(context);
      }

      _consumerContextProviderContextMap.remove(context);

      // TODO: How to make sure the context is not referred in any other list - use some counter map?
      // Make sure related contexts are also freed up.
   }

   /**
    * Allocate the provider context to the consumer context that requested it.
    * After allocating it the provider context can not be allocated to other
    * consumer.
    *
    * @param consumerContext
    *       context of the workflow that requested the testbed.
    * @param providerContext
    *       provider contet allocated for the needs of the consumer.
    */
   public void mapProviderToConsumer(WorkflowContext consumerContext,
         ProviderWorkflowContext providerContext) {
      // Verify the consumer context is registered
      if (!_consumerContextProviderContextMap.containsKey(consumerContext)) {
         throw new IllegalArgumentException(
               "Consumer context unknown. It is not register in the registry.");
      }

      // Verify the provider context is registered
      if (!_consumerContextProviderContextMap.containsKey(providerContext)) {
         throw new IllegalArgumentException(
               "Provider context unknown. It is not register in the registry.");
      }

      Set<ProviderWorkflowContext> providerContextList =
            _consumerContextProviderContextMap.get(consumerContext);
      if (providerContextList.contains(providerContext)) {
         throw new IllegalArgumentException(
               "Specified provider context already registered with consumer.");
      }

      providerContextList.add(providerContext);
   }

   /**
    * Set of the provider contexts allocated for the consumer.
    * @param consumerContext context of the workflow.
    * @return
    */
   public Set<ProviderWorkflowContext> getRegisteredConsumerProviders(
         WorkflowContext consumerContext) {
      if (_consumerContextProviderContextMap.containsKey(consumerContext)) {
         return _consumerContextProviderContextMap.get(consumerContext);
      } else {
         return new HashSet<ProviderWorkflowContext>();
      }
   }

   /**
    * Browse for the first available(that is not already assigned) testbed
    * for the specified provider and assign it.
    * @param providerContext the needed provider workflow.
    * @return assigned testbed context.
    */
   public TestBedContext claimTestbed(ProviderWorkflowContext providerContext) {
      if (providerContext == null) {
         throw new IllegalArgumentException(
               "Required providerContext parameter is not set.");
      }

      // A tesbed is already assigned to that context
      if (_providerContextTestbedContextMap.containsKey(providerContext)) {
         throw new IllegalArgumentException(
               "A test bed is already associated with this provider context.");
      }

      // Load registered tesbeds for the needed workflow class.
      List<TestBedContext> testbedContextList =
            _providerClassTestbedContextMap.get(providerContext.getProviderWorkflow()
                  .getClass());
      if (testbedContextList == null || testbedContextList.size() == 0) {
         throw new RuntimeException(String.format(
               "No test beds are registered for the %s provider.",
               providerContext.getProviderWorkflow().getClass().toString()));
      }

      // Find available testbed - the first that is not already assigned.
      TestBedContext availableTestbedContext = null;
      for (TestBedContext testbedContext : testbedContextList) {
         if (!_providerContextTestbedContextMap.containsValue(testbedContext)) {
            availableTestbedContext = testbedContext;
            break;
         }
      }

      if (availableTestbedContext == null) {
         throw new RuntimeException(
               String.format(
                     "All of the '%s' testbeds registered for provider: \"%s\" are in use.",
                     testbedContextList.size(), providerContext
                           .getProviderWorkflow().getClass().toString()));
      }

      _providerContextTestbedContextMap.put(providerContext, availableTestbedContext);

      return availableTestbedContext;
   }

   /**
    * Free allocated resource(provider context) assigned.
    * @param providerContext  allocated resource to consumer that requested it.
    * TODO rkovachev: what about invoking checkhealth after releasing resource?
    */
   public void freeTestbed(ProviderWorkflowContext providerContext) {
      if (providerContext == null) {
         // nothing to do
         return;
      }

      // Unlink it from the consumer.
      if (_providerContextTestbedContextMap.containsKey(providerContext)) {
         _providerContextTestbedContextMap.remove(providerContext);
      }
   }

   /**
    * Return TestBedContext for the provided context.
    * @param providerContext
    * @return
    */
   public TestBedContext getProviderTestbed(ProviderWorkflowContext providerContext) {
      if (providerContext == null) {
         throw new IllegalArgumentException(
               "Required providerContext parameter is not set.");
      }

      if (!_providerContextTestbedContextMap.containsKey(providerContext)) {
         throw new IllegalArgumentException(
               "The specificed provider context does not have assigned test bed.");
      }

      return _providerContextTestbedContextMap.get(providerContext);
   }


   /**
    * Check if a corresponding class exists for a given name.
    *
    * The class name looked up in the map of <code>Workflow</code> classes. If found,
    * a check is made of the discovered class is has the specified super class
    * (workflowType).
    *
    * @param workflowClassName
    *      Name of the workflow class.
    * @param workflowType
    *      Super class for the given workflow class. If null, the check is omitted.
    * @return
    *      True - workflow class exists
    */
   public boolean checkWorkflowExists(String workflowClassName,
         Class<? extends Workflow> workflowType) {

      if (!_workflowNameWorkflowClassMap.containsKey(workflowClassName)) {
         return false;
      }

      if (workflowType != null) {
         Class<? extends Workflow> workflowClass =
               _workflowNameWorkflowClassMap.get(workflowClassName);

         // Check if the classes are the same type
         return (workflowType.isAssignableFrom(workflowClass));
      } else {
         return true;
      }
   }

   /**
    *
    * @param workflowClassName
    * @param workflowType
    * @return
    * @throws ClassNotFoundException
    */
   public Class<? extends Workflow> getRegisteredWorkflowClass(String workflowClassName,
         Class<? extends Workflow> workflowType) throws ClassNotFoundException {

      // Check workflowClassName is not empty
      if (Strings.isNullOrEmpty(workflowClassName)) {
         throw new IllegalArgumentException(
               "The required workflowClassName parameter is not set.");
      }

      // Check the require workflow class is available.
      if (!_workflowNameWorkflowClassMap.containsKey(workflowClassName)) {
         throw new ClassNotFoundException(String.format(
               "The specified workflow class \"%s\" cannot be found.",
               workflowClassName));
      }

      // Check it's from the desired type
      Class<? extends Workflow> workflowClass =
            _workflowNameWorkflowClassMap.get(workflowClassName);
      if (!workflowType.isAssignableFrom(workflowClass)) {
         throw new ClassNotFoundException(
               "The workflow class name is from unexpected workflow type.");
      }

      return workflowClass;
   }

   /**
    * Get the connector registered or the provided service spec.
    * @param serviceSpec
    * @return
    */
   public TestbedConnector getActiveTestbedConnection(ServiceSpec serviceSpec) {
      if (serviceSpec == null) {
         throw new IllegalArgumentException("Required serviceSpec parameter is not set.");
      }

      String connectionKey = serviceSpec.toString();
      Set<ProviderWorkflowContext> contexts = _providerContextTestbedContextMap.keySet();
      for (ProviderWorkflowContext context : contexts) {
         TestbedConnector connector = context.getKeyConncetionMap().get(connectionKey);
         if (connector != null) {
            return connector;
         }
      }

      return null;
   }

   /**
    * Provide session settings.
    * @return SettingReader object containing session configurations.
    */
   public SettingsReader getSessionSettingsReader() {
      return _sessionSettings;
   }

   /**
    * Add testbed context to the list of testbeds defined by providerClass.
    * @param providerClass
    * @param newTestBed
    * @param testbedSettings
    */
   private void addTestbed(Class<? extends ProviderWorkflow> providerClass,
         TestBedContext newTestBed, Properties testbedSettings ) {
      List<TestBedContext> providerTestbeds =
            _providerClassTestbedContextMap.get(providerClass);
      if (CollectionUtils.isEmpty(providerTestbeds)) {
         // Create new test bed list, if this is first registered one for the provider type.
         providerTestbeds = new ArrayList<TestBedContext>();
         _providerClassTestbedContextMap.put(providerClass, providerTestbeds);
      }

      // Validate that one provider in not loaded multiple times by comparing its properties for equal match among their values.
      for (TestBedContext testbedContext : providerTestbeds) {
         if (testbedContext.testbedEquals(testbedSettings)) {
            throw new IllegalArgumentException(
                  "The same test bed is already registered.");
         }
      }

      // Store the new test bed
      providerTestbeds.add(newTestBed);
   }

   @SuppressWarnings("unchecked")
   private void initRegistry() {
      _sessionSettings = new SettingsRWImpl();
      _workflowClassSet = new HashSet<Class<? extends Workflow>>();
      _workflowNameWorkflowClassMap = new HashMap<String, Class<? extends Workflow>>();

      _consumerContextProviderContextMap =
            new HashMap<WorkflowContext, Set<ProviderWorkflowContext>>();

      _providerContextTestbedContextMap =
            new HashMap<ProviderWorkflowContext, TestBedContext>();

      _providerClassTestbedContextMap =
            new HashMap<Class<? extends ProviderWorkflow>, List<TestBedContext>>();


      try {
         _workflowClassSet.addAll(WorkflowBrowser.findExecutableWorkflows());

         for (Class<? extends Workflow> workflowClass : _workflowClassSet) {

            _workflowNameWorkflowClassMap.put(
                  workflowClass.getCanonicalName(),
                  workflowClass);

            if (workflowClass.isAssignableFrom(ProviderWorkflow.class)) {
               _providerClassTestbedContextMap.put(
                     (Class<ProviderWorkflow>) workflowClass,
                     new ArrayList<TestBedContext>());
            }

         }

      } catch (Exception e) {
         // TODO rkovachev Throw the error properly
         throw new RuntimeException("Registry initilization failed!", e);
      }
   }

}
