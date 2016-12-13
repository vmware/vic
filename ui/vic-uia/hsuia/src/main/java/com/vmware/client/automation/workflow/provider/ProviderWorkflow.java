/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import java.util.Map;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.workflow.common.Workflow;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;

/**
 * Implement this interface to construct a provider workflow. Provider workflows
 * are used to deploy, build, configure, validate and destroy test beds from
 * scratch.
 *
 * One provider workflow may include entities prepared by other provider
 * workflows.
 */
public interface ProviderWorkflow extends Workflow {

   /**
    * Invoked first - define the spec to be populated. Add mapping to it in the
    * registry.
    * @param publisherSpec
    * @throws Exception
    */
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception;

   /**
    * Define the spec to be used in assemble, disassemble and checkhealth stages.
    * Here the provider may request testbed through the TestBedBridge.
    * @param assemblerSpec
    * @param testbedBridge
    * @throws Exception
    */
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception;

   /**
    * Populate data in the publisher spec of the provider.
    *
    * @param providerSpec
    * @param testbedSettings
    * @param contextSettings
    */
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception;

   /**
    * Assign settings for the assembler spec of the provider.
    * @param assemblerSpec
    * @param testbedSettings
    * @throws Exception
    */
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception;

   /**
    * Set proper connectors for the specified service specs.
    * @param serviceConnectorsMap
    * @throws Exception
    */
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) throws Exception;

   /**
    * Compose steps to be executed during the provider workflow stages.
    * @param flow
    * @throws Exception
    */
   public void composeProviderSteps(WorkflowStepsSequence<? extends WorkflowStepContext> flow) throws Exception;

   /**
    * Define the weight of the testbed. The higher the weight means that more
    * resources is need to provide the published testbed.
    * @return
    */
   public int providerWeight();
   
   /**
    * The base elemental provider takes care for the resource deployment -
    * i.e. deployment in Nimbus, deployment in Cloud or other system
    * (it could be even a physical resource).
    * Bases on that there are NimbusHostProvider, CloudHostProvider and etc..
    * By providing this method it would possible to register a NimbusHostProvider
    * resource and HostProvider resource.
    * @return
    */
   public Class<? extends ProviderWorkflow> getProviderBaseType();
}
