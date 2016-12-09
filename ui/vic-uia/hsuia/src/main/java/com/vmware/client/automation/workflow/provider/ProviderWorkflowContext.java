/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import java.util.HashMap;
import java.util.Map;

import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.workflow.common.WorkflowContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsRWImpl;

/**
 * The context class keeps the configuration, runtime state and statistics
 * data about a <code>ProviderWorkflow</code>.
 * 
 * In a nut shell it keeps all the state of a provider workflow.
 */
public class ProviderWorkflowContext implements WorkflowContext {

   private final Class<ProviderWorkflow> _providerWorkflowClass;

   private ProviderWorkflow _providerWorkflow;

   private PublisherSpec _publisherSpec;

   private AssemblerSpec _assemblerSpec;

   private WorkflowStepsSequence<ProviderWorkflowStepContext> _flow;

   private final Map<String, TestbedConnector> _keyConnectionMap;

   private final SettingsRWImpl _testbedSettings;

   private final String _resourceType;

   /**
    * Default constructor.
    * Initialize all members, so the getters are accessible from this point on.
    * 
    * @param providerWorkflow
    *       Reference to the provider workflow
    */
   public ProviderWorkflowContext(
         Class<ProviderWorkflow> providerWorkflowClass) {
      _providerWorkflowClass = providerWorkflowClass;
      _testbedSettings = new SettingsRWImpl();
      _keyConnectionMap = new HashMap<String, TestbedConnector>();
      _resourceType = "unknown";
   }

   /**
    * Call this method to create an instance of the workflow.
    * 
    * If an instance was already present it will be deleted. Related
    * data, which could be restored by running the worflow's spec
    * initialization and steps' composition will be also deleted.
    * 
    * Note that this method does not delete the all the data stored
    * in the context.
    * 
    * @throws InstantiationException
    * @throws IllegalAccessException
    *       In case of error, the old instances will be deleted.
    * 
    * 
    */
   public void createWorkflowInstance() throws InstantiationException, IllegalAccessException {
      _providerWorkflow = null;
      _publisherSpec = null;
      _assemblerSpec = null;
      _flow = null;

      // TODO rkovachev: it is quite ambiguous instantiating provider
      // workflow by reflection here - see the requestTestbed method of
      // the ProviderBridgeImpl class.
      _providerWorkflow = _providerWorkflowClass.newInstance();
      _publisherSpec = new PublisherSpec();
      _assemblerSpec = new AssemblerSpec();
      _flow = new WorkflowStepsSequence<ProviderWorkflowStepContext>();
   }

   /**
    * Return the provider's workflow instance.
    */
   public ProviderWorkflow getProviderWorkflow() {
      return _providerWorkflow;
   }

   /**
    * Return the provider's publisher spec.
    */
   public PublisherSpec getPublisherSpec() {
      return _publisherSpec;
   }

   /**
    * Return the provider's assembler spec.
    */
   public AssemblerSpec getAssemblerSpec() {
      return _assemblerSpec;
   }

   /**
    * Return the provider's steps holder.
    */
   public WorkflowStepsSequence<ProviderWorkflowStepContext> getFlow() {
      return _flow;
   }

   /**
    * Returns the test bed setting's holder.
    */
   public SettingsRWImpl getTestbedSettings() {
      return _testbedSettings;
   }

   /**
    * Return the map of registered connectors for the provider.
    * @return
    */
   public Map<String, TestbedConnector> getKeyConncetionMap() {
      return _keyConnectionMap;
   }

   public String getResourceType() {
      return _resourceType;
   }
}
