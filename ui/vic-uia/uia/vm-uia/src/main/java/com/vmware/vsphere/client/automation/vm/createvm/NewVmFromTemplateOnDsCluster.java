/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm;

import org.apache.commons.lang.RandomStringUtils;
import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.common.step.VerifyTaskByUiStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.step.ClickNextWizardButtonStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.DatastoreClusterLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmTemplateInFolderLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.DatastoreClusterNavigationStep;
import com.vmware.vsphere.client.automation.components.navigator.step.VmTemplateInFolderNavigationStep;
import com.vmware.vsphere.client.automation.dscluster.common.spec.CreateDsClusterSpec;
import com.vmware.vsphere.client.automation.dscluster.common.step.CreateDsClusterByApiStep;
import com.vmware.vsphere.client.automation.dscluster.common.step.MoveDatastoresToDsClusterStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.provider.commontb.VmfsStorageProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SdrsBehavior;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateDatastoreStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmTemplateByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.VerifyVmExistenceByApiStep;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.AttachIscsiTargetStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec.VmCreationType;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.FinishCreateVmWizardPageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.LaunchNewVmFromThisTemplateStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectComputeResourcePageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectNameAndFolderPageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectStoragePageStep;
import com.vmware.vsphere.client.automation.vm.lib.createvm.step.VerifyVmExistsInDsClusterVmsListStep;

/**
 * Test class which deploys VM from VM Template on datastore cluster. Here are
 * the test steps:<br>
 * - Create datastore cluster with VMFS datastore.<br>
 * - Create VM template.<br>
 * - Launch the Deploy VM from VM template wizard.<br>
 * - In the Deploy VM from VM template wizard select the datastore cluster for
 * storage.<br>
 * - Complete the wizard.<br>
 * - Verify that the VM is successfully created on the datastore cluster.<br>
 */
public class NewVmFromTemplateOnDsCluster extends NGCTestWorkflow {
   private static final String TAG_STANDALONE_DATASTORE = "TAG_STANDALONE_DATASTORE";
   private static final String TAG_CLUSTERED_DATASTORE = "TAG_CLUSTERED_DATASTORE";
   private static final String TAG_SOURCE_VM = "TAG_SOURCE_VM";
   private static final String TAG_DESTINATION_VM = "TAG_DESTINATION_VM";

   @Override
   @Test(description = "Deploy VM from VM Template on datastore cluster", groups = { CAT, BAT })
   @TestID(id = "617000")
   public void execute() throws Exception {
      super.execute();
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer commonTestbed = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);

      TestbedSpecConsumer hostProvider = testbedBridge.requestTestbed(
            HostProvider.class, false);

      TestbedSpecConsumer vmfsStorageProvider = testbedBridge.requestTestbed(
            VmfsStorageProvider.class, false);

      // VC
      VcSpec vcSpec = commonTestbed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Datacenter
      DatacenterSpec datacenterSpec = commonTestbed
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      // Cluster
      ClusterSpec clusterSpec = commonTestbed
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_ENTITY);

      // Clustered host
      HostSpec hostSpec = hostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      hostSpec.parent.set(datacenterSpec);

      // Common datastore, where the VM template will reside
      DatastoreSpec saDatastoreSpec = hostProvider
            .getPublishedEntitySpec(HostProvider.LOCAL_DS_ENTITY);
      saDatastoreSpec.tag.set(TAG_STANDALONE_DATASTORE);

      // Datastore cluster
      DatastoreClusterSpec datastoreClusterSpec = SpecFactory.getSpec(
            DatastoreClusterSpec.class, datacenterSpec);
      datastoreClusterSpec.ioBalancingEnabled.set(true);
      datastoreClusterSpec.sdrsEnabled.set(false);
      datastoreClusterSpec.sdrsBehavior.set(SdrsBehavior.MANUAL);
      datastoreClusterSpec.tag.set(TAG_CLUSTERED_DATASTORE);

      // Setting the hostSpec.iscsiServerIp from the VmfsStorageProvider
      DatastoreSpec helperVmfsDatastoreSpec = vmfsStorageProvider
            .getPublishedEntitySpec(VmfsStorageProvider.DEFAULT_ENTITY);
      hostSpec.iscsiServerIp.set(helperVmfsDatastoreSpec.remoteHost);

      // Datastore in datastore cluster
      DatastoreSpec newDatastoreSpec = SpecFactory.getSpec(DatastoreSpec.class,
            hostSpec);
      newDatastoreSpec.type.set(DatastoreType.VMFS);
      newDatastoreSpec.tag.set(TAG_CLUSTERED_DATASTORE);

      // Spec for the 'Create datastore cluster' dialog
      CreateDsClusterSpec createDsClusterSpec = new CreateDsClusterSpec();
      createDsClusterSpec.datastoresParents.set(clusterSpec);
      createDsClusterSpec.datastores.set(newDatastoreSpec);

      // VM template
      VmSpec vmTemplateSpec = SpecFactory.getSpec(VmSpec.class, null, hostSpec);
      vmTemplateSpec.datastore.set(saDatastoreSpec);
      vmTemplateSpec.tag.set(TAG_SOURCE_VM);

      // Location specs
      VmTemplateInFolderLocationSpec vmTemplateLocation = new VmTemplateInFolderLocationSpec(
            vmTemplateSpec.name.get());
      DatastoreClusterLocationSpec datastoreClusterLocationSpec = new DatastoreClusterLocationSpec(
            datastoreClusterSpec.name.get(),
            NGCNavigator.NID_ENTITY_PRIMARY_TAB_VMS);

      // Spec for the VM that will be created
      CreateVmSpec createVmSpec = SpecFactory.getSpec(CreateVmSpec.class,
            hostSpec);
      createVmSpec.name.set("DeployedVm"
            + RandomStringUtils.randomAlphanumeric(10));
      createVmSpec.creationType.set(VmCreationType.CREATE_NEW_VM);
      createVmSpec.computeResource.set(hostSpec);
      createVmSpec.datastoreCluster.set(datastoreClusterSpec);
      createVmSpec.datastore.set(newDatastoreSpec);
      createVmSpec.tag.set(TAG_DESTINATION_VM);

      // Create VM task
      TaskSpec createVmTaskSpec = new TaskSpec();
      createVmTaskSpec.name.set(VmUtil.getLocalizedString("task.cloneVm.name"));
      createVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      createVmTaskSpec.target.set(vmTemplateSpec);

      testSpec.add(vcSpec, hostSpec, datacenterSpec, datastoreClusterSpec,
            saDatastoreSpec, newDatastoreSpec, createDsClusterSpec,
            vmTemplateLocation, datastoreClusterLocationSpec, vmTemplateSpec,
            createVmSpec, createVmTaskSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      flow.appendStep("Add standalone host to the inventory.",
            new AddHostStep());

      // Create VMFS datastore
      flow.appendStep("Attach iSCSI target", new AttachIscsiTargetStep());
      flow.appendStep("Create VMFS datastore", new CreateDatastoreStep(),
            new String[] { TAG_CLUSTERED_DATASTORE });

      // Create datastore cluster
      flow.appendStep("Create datastore cluster via API",
            new CreateDsClusterByApiStep());
      flow.appendStep(
            "Move the VMFS datastore to the datastore cluster via API",
            new MoveDatastoresToDsClusterStep(),
            new String[] { TAG_CLUSTERED_DATASTORE });

      // Create VM template
      flow.appendStep("Create VM template via API",
            new CreateVmTemplateByApiStep(), new String[] { TAG_SOURCE_VM });
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      // Launch the Deploy VM from Template wizard
      flow.appendStep("Navigate to the VM template",
            new VmTemplateInFolderNavigationStep());
      flow.appendStep("Launch the New VM wizard.",
            new LaunchNewVmFromThisTemplateStep());

      // Complete the Deploy VM from Template wizard
      flow.appendStep("Set name for the VM.", new SelectNameAndFolderPageStep());
      flow.appendStep("Select cluster compute resource.",
            new SelectComputeResourcePageStep());
      flow.appendStep("Select storage to be used for the vm.",
            new SelectStoragePageStep());
      flow.appendStep(
            "Pass through Select Clone options page with default selections.",
            new ClickNextWizardButtonStep());
      flow.appendStep("Finish the wizard.", new FinishCreateVmWizardPageStep());

      // Verify that the VM is deployed successfully
      flow.appendStep("Verify the clone VM task via UI",
            new VerifyTaskByUiStep());
      flow.appendStep("Verify that the deployed VM exists via API.",
            new VerifyVmExistenceByApiStep(),
            new String[] { TAG_DESTINATION_VM });
      flow.appendStep("Navigate to the Datastore Cluster",
            new DatastoreClusterNavigationStep());
      flow.appendStep("Verify the VM is deployed on the datastore cluster.",
            new VerifyVmExistsInDsClusterVmsListStep(),
            new String[] { TAG_DESTINATION_VM });
   }
}
