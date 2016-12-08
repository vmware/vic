/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.simulator;

import java.util.Map;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.simulator.connector.SimulatorConnectorsFactory;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorClusterSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorDatacenterSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorHostSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorServiceSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorVcSpec;
import com.vmware.vsphere.client.automation.provider.simulator.step.AddHostToClusterSimStep;
import com.vmware.vsphere.client.automation.provider.simulator.step.AddHostToDatacenterSimStep;
import com.vmware.vsphere.client.automation.provider.simulator.step.CreateVirtualInventorySimStep;

/**
 *
 * -Virtual Center (d1-group) | -Datacenter (basic-datacenter) | -Host 1
 * (dynamically generated) | -Cluster (basic-cluster) | -Host 2 (dynamically
 * generated)
 *
 */
public class BasicInventorySimProvider extends BaseCompositeProvider {

   public static final String ENTITY_VC = "provider.simulator.basicinventory.entity.vc";
   public static final String ENTITY_DATACENTER = "provider.simulator.basicinventory.entity.datacenter";
   public static final String ENTITY_HOST_DATACENTER = "provider.simulator.basicinventory.entity.sahost";
   public static final String ENTITY_CLUSTER = "provider.simulator.basicinventory.entity.cluster";
   public static final String ENTITY_HOST_CLUSTERED = "provider.simulator.basicinventory.entity.clusterhost";

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {

      SimulatorVcSpec publishedVcSpec = new SimulatorVcSpec();
      publisherSpec.publishEntitySpec(ENTITY_VC, publishedVcSpec);
      publisherSpec.links.add(publishedVcSpec);

      SimulatorDatacenterSpec publishedDatacenterSpec = new SimulatorDatacenterSpec();
      publishedDatacenterSpec.parent.set(publishedVcSpec);
      publisherSpec.add(publishedDatacenterSpec);
      publisherSpec.publishEntitySpec(ENTITY_DATACENTER,
            publishedDatacenterSpec);

      SimulatorClusterSpec publishedClusterSpec = new SimulatorClusterSpec();
      publishedClusterSpec.parent.set(publishedDatacenterSpec);
      publisherSpec.add(publishedClusterSpec);
      publisherSpec.publishEntitySpec(ENTITY_CLUSTER, publishedClusterSpec);

      SimulatorHostSpec publishedClusterHostSpec = new SimulatorHostSpec();
      publishedClusterHostSpec.parent.set(publishedClusterSpec);
      publisherSpec.add(publishedClusterHostSpec);
      publisherSpec.publishEntitySpec(ENTITY_HOST_CLUSTERED,
            publishedClusterHostSpec);

      publisherSpec.links.getAll(EntitySpec.class);
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec,
         TestBedBridge testbedBridge) throws Exception {

      // ===========================================================================
      // VC
      SimulatorVcSpec assemblerVcSpec = testbedBridge.requestTestbed(
            VcSimProvider.class, true).getPublishedEntitySpec(
                  VcSimProvider.DEFAULT_ENTITY);
      assemblerSpec.links.add(assemblerVcSpec);

      // Datacenter TODO: Replace with spec builder after specs relationships
      // model is finally sorted.
      SimulatorDatacenterSpec assemblerDcSpec = new SimulatorDatacenterSpec();
      assemblerDcSpec.name.set("generate name");
      assemblerDcSpec.parent.set(assemblerVcSpec);
      assemblerSpec.links.add(assemblerDcSpec);

      // Cluster TODO: Replace with spec builder after specs relationships model
      // is finally sorted.
      SimulatorClusterSpec assemblerClusterSpec = new SimulatorClusterSpec();
      assemblerClusterSpec.name.set("cluster-1");
      assemblerClusterSpec.parent.set(assemblerDcSpec);
      assemblerSpec.links.add(assemblerClusterSpec);

      // Cluster host
      SimulatorHostSpec assemblerClusterHostSpec =
            testbedBridge.requestTestbed(HostSimProvider.class, false).getPublishedEntitySpec(HostSimProvider.DEFAULT_ENTITY);

      assemblerClusterHostSpec.parent.set(assemblerClusterSpec);
      assemblerSpec.links.add(assemblerClusterHostSpec);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      // Get settings
      String endpoint = SettingsUtil.getRequiredValue(testbedSettings,
            CreateVirtualInventorySimStep.TESTBED_KEY_SERVICE_ENDPOINT);
      String username = SettingsUtil.getRequiredValue(testbedSettings,
            CreateVirtualInventorySimStep.TESTBED_KEY_SERVICE_USERNAME);
      String password = SettingsUtil.getRequiredValue(testbedSettings,
            CreateVirtualInventorySimStep.TESTBED_KEY_SERVICE_PASSWORD);

      // Create service spec
      SimulatorServiceSpec serviceSpec = new SimulatorServiceSpec();
      serviceSpec.endpoint.set(endpoint);
      serviceSpec.username.set(username);
      serviceSpec.password.set(password);

      SimulatorVcSpec publishedVcSpec = publisherSpec
            .getPublishedEntitySpec(ENTITY_VC);
      SimulatorDatacenterSpec publishedDatacenterSpec = publisherSpec
            .getPublishedEntitySpec(ENTITY_DATACENTER);
      SimulatorClusterSpec publishedClusterSpec = publisherSpec
            .getPublishedEntitySpec(ENTITY_CLUSTER);
      SimulatorHostSpec publishedClusterHostSpec = publisherSpec
            .getPublishedEntitySpec(ENTITY_HOST_CLUSTERED);

      // Assign the service spec to all published entities.

      publishedVcSpec.service.set(serviceSpec);
      publishedDatacenterSpec.service.set(serviceSpec);
      publishedClusterSpec.service.set(serviceSpec);
      publishedClusterHostSpec.service.set(serviceSpec);

      // Where to load the specs?:
      // O1 - Do it in the steps; - This will probably lead to extracting and
      // assigning the test bed settings inside the step;
      // We'll have to start using non-cloned specs, which is difficult to
      // process. In order to solve it, we'll need the spec querying API,
      // so we can get actual mutable references instead of clones.
      // O2 - Do it here in this method; Taking this one until the obstacles for
      // O1 are removed.
      // TODO: Consider and implement O1

      // Logically it's where we do the assembling where we should perform the
      // spec initialization.

      // Set the spec names
      publishedVcSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            CreateVirtualInventorySimStep.TESTBED_KEY_VC_NAME));

      publishedDatacenterSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            CreateVirtualInventorySimStep.TESTBED_KEY_DC_NAME));

      publishedClusterSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            CreateVirtualInventorySimStep.TESTBED_KEY_CL_NAME));

      publishedClusterHostSpec.name.set(SettingsUtil.getRequiredValue(
            testbedSettings, AddHostToClusterSimStep.TESTBED_KEY_CLHOST_NAME));
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap)
               throws Exception {
      SimulatorConnectorsFactory.createAndSetConnectors(serviceConnectorsMap);
   }


   @Override
   public int providerWeight() {
      // TODO Auto-generated method stub
      return 0;
   }

   @Override
   public void composeProviderSteps(
         WorkflowStepsSequence<? extends WorkflowStepContext> flow)
         throws Exception {
      flow.appendStep("Create inventory", new CreateVirtualInventorySimStep());

      flow.appendStep("Add host in datacenter", new AddHostToDatacenterSimStep());

      flow.appendStep("Add host in cluster", new AddHostToClusterSimStep());

   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      // TODO Auto-generated method stub
      return null;
   }

}
