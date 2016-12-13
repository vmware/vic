/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.common.step.VerifyTaskByUiStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.common.step.ReconfigureClusterByApiStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.ClusterLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.ClusterNavigationStep;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.AdmissionControlFailoverPolicy;
import com.vmware.vsphere.client.automation.srv.common.spec.AdmissionControlSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpecBuilder;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec.FaultToleranceVmRoles;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateClusterStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.EnableFtOnVnicViaApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.EnableVmotionOnVnicViaApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.PowerOnVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.util.FtUtil;
import com.vmware.vsphere.client.automation.vm.ft.step.CheckFtPortletStep;
import com.vmware.vsphere.client.automation.vm.ft.step.CheckFtVmNamesStep;
import com.vmware.vsphere.client.automation.vm.ft.step.ClickYesOnWarningStep;
import com.vmware.vsphere.client.automation.vm.ft.step.FtFinishWizardStep;
import com.vmware.vsphere.client.automation.vm.ft.step.FtSelectDatastoreStep;
import com.vmware.vsphere.client.automation.vm.ft.step.FtSelectHostStep;
import com.vmware.vsphere.client.automation.vm.ft.step.LaunchTurnOnFaultToleranceWizardStep;
import com.vmware.vsphere.client.automation.vm.migrate.step.MountDatastoreByApiStep;

/**
 * Test class for turning on fault tolerance for a vm in the NGC client.<br>
 * Preconditions:<br>
 * <ol>
 * <li>Create a new cluster</li>
 * <li>Add two hosts to it</li>
 * <li>Mount the same NFS datastrore to both hosts</li>
 * <li>Turn on Fault Tolerance and vMotion for each host's vnic</li>
 * <li>Turn on HA for the cluster</li>
 * <li>Create a two CPU VM on the first host</li>
 * <li>Power on the VM</li>
 * </ol>
 * Steps:<br>
 * <ol>
 * <li>Open a browser</li>
 * <li>Login as admin user</li>
 * <li>Navigate to the VM Summary page</li>
 * <li>Turn on Fault Tolerance via More Actions</li>
 * <li>Click Yes in the pop up dialog</li>
 * <li>Select the NFS datastore</li>
 * <li>Select the second host</li>
 * <li>Complete the wizard by clicking Finish</li>
 * <li>Verify the tasks for primary and secondary VM complete successfully</li>
 * <li>Verify the Fault Tolerance status is Protected in the Fault Tolerance
 * portlet of the VM</li>
 * <li>Navigate to the Cluster > VMs > Virtual Machines</li>
 * <li>Verify there are two VMs and their names contain Primary and Secondary</li>
 * </ol>
 */
public class TurnOnFaultToleranceTest extends NGCTestWorkflow {
   protected static final String TAG_FIRST_HOST = "TAG_FIRST_HOST";
   protected static final String TAG_SECOND_HOST = "TAG_SECOND_HOST";
   protected static final String TAG_INITIAL_CLUSTER = "TAG_INITIAL_CLUSTER";
   protected static final String TAG_EDITED_CLUSTER = "TAG_EDITED_CLUSTER";
   protected static final String TAG_FIRST_HOST_DATASTORE = "TAG_FIRST_HOST_DATASTORE";
   protected static final String TAG_SECOND_HOST_DATASTORE = "TAG_SECOND_HOST_DATASTORE";
   protected static final String TAG_TURN_ON_FT = "TAG_TURN_ON_FT";
   protected static final String TAG_TURN_ON_SECONDARY = "TAG_TURN_ON_SECONDARY";
   protected static final String TAG_FAULT_TOLERANCE_SPEC = "TAG_FAULT_TOLERANCE_SPEC";

   @Override
   @Test(description = "Turn on Fault Tolerance for a VM")
   @TestID(id = "628109")
   public void execute() throws Exception {
      super.execute();
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {

      TestbedSpecConsumer testbed = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);
      TestbedSpecConsumer firstHostProvider = testbedBridge.requestTestbed(
            HostProvider.class, false);
      TestbedSpecConsumer secondHostProvider = testbedBridge.requestTestbed(
            HostProvider.class, false);

      // Spec for the VC
      VcSpec requestedVCSpec = testbed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Spec for the datacenter
      DatacenterSpec requestedDatacenterSpec = testbed
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      // Spec for the initial cluster
      ClusterSpec initialClusterSpec = SpecFactory.getSpec(ClusterSpec.class,
            requestedDatacenterSpec);
      initialClusterSpec.drsEnabled.set(false);
      initialClusterSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.OLD);
      initialClusterSpec.tag.set(TAG_INITIAL_CLUSTER);

      // Spec for cluster HA admission control
      AdmissionControlSpec admissionControlSpec = SpecFactory.getSpec(
            AdmissionControlSpec.class, initialClusterSpec);
      admissionControlSpec.failoverPolicy
            .set(AdmissionControlFailoverPolicy.DISABLED);

      // Spec for the edited cluster - HA is turned on
      ClusterSpec editedClusterSpec = SpecFactory.getSpec(ClusterSpec.class,
            requestedDatacenterSpec);
      editedClusterSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.NEW);
      editedClusterSpec.drsEnabled.set(false);
      editedClusterSpec.vsphereHA.set(true);
      editedClusterSpec.admissionControlSpec.set(admissionControlSpec);
      editedClusterSpec.tag.set(TAG_EDITED_CLUSTER);

      // Specs for the hosts
      HostSpec firstHostSpec = firstHostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      firstHostSpec.parent.set(initialClusterSpec);
      firstHostSpec.tag.set(TAG_FIRST_HOST);

      HostSpec secondHostSpec = secondHostProvider
            .getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      secondHostSpec.parent.set(initialClusterSpec);
      secondHostSpec.tag.set(TAG_SECOND_HOST);

      // Specs for the datastores
      DatastoreSpec firstDatastoreSpec = testbed
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_DS_ENTITY);
      firstDatastoreSpec.parent.set(firstHostSpec);
      firstDatastoreSpec.tag.set(TAG_FIRST_HOST_DATASTORE);

      DatastoreSpec secondDatastoreSpec = new DatastoreSpec();
      secondDatastoreSpec.copy(firstDatastoreSpec);
      secondDatastoreSpec.parent.set(secondHostSpec);
      secondDatastoreSpec.tag.set(TAG_SECOND_HOST_DATASTORE);

      // Spec for the VM
      VmSpec vmSpec = SpecFactory.getSpec(VmSpec.class, null, firstHostSpec);
      vmSpec.datastore.set(firstDatastoreSpec);
      vmSpec.numCPUs.set(2);
      vmSpec.ftRole.set(FaultToleranceVmRoles.PRIMARY);

      // Specs for the location to the VM > Summary page
      VmLocationSpec vmLocationSpec = new VmLocationSpec(vmSpec,
            NGCNavigator.NID_ENTITY_PRIMARY_TAB_SUMMARY);

      // Spec for the location to the Cluster -> VMs > Virtual Machines
      ClusterLocationSpec clusterLocationSpec = new ClusterLocationSpec(
            initialClusterSpec, NGCNavigator.NID_ENTITY_PRIMARY_TAB_VMS,
            NGCNavigator.NID_CLUSTER_VMS_II_VMS);

      // Spec for Fault Tolerance
      FaultToleranceSpec faultToleranceSpec = new FaultToleranceSpecBuilder()
            .setTags(TAG_FAULT_TOLERANCE_SPEC).build();

      // Task specs
      TaskSpec turnOnFtTaskSpec = new TaskSpec();
      turnOnFtTaskSpec.name
            .set(FtUtil.getLocalizedString("task.turnonft.name"));
      turnOnFtTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      turnOnFtTaskSpec.target.set(vmSpec);
      turnOnFtTaskSpec.tag.set(TAG_TURN_ON_FT);

      TaskSpec turnOnSecondaryVmTaskSpec = new TaskSpec();
      turnOnSecondaryVmTaskSpec.name.set(FtUtil
            .getLocalizedString("task.turnonsecondary.name"));
      turnOnSecondaryVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      turnOnSecondaryVmTaskSpec.target.set(vmSpec);
      turnOnSecondaryVmTaskSpec.tag.set(TAG_TURN_ON_SECONDARY);

      // Initialize test bed
      testSpec.add(requestedVCSpec, requestedDatacenterSpec,
            initialClusterSpec, editedClusterSpec, firstHostSpec,
            firstDatastoreSpec, secondDatastoreSpec, secondHostSpec, vmSpec,
            vmLocationSpec, clusterLocationSpec, faultToleranceSpec,
            turnOnFtTaskSpec, turnOnSecondaryVmTaskSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      // Create a cluster via API
      flow.appendStep("Create a cluster via API", new CreateClusterStep(),
            new String[] { TAG_INITIAL_CLUSTER });

      // Add hosts via API
      flow.appendStep("Add the first host to the inventory", new AddHostStep(),
            new String[] { TAG_FIRST_HOST });
      flow.appendStep("Add the second host to the inventory",
            new AddHostStep(), new String[] { TAG_SECOND_HOST });

      // Mount NFS to hosts via API
      flow.appendStep("Mount the NFS datastore to the first host via API",
            new MountDatastoreByApiStep(),
            new String[] { TAG_FIRST_HOST_DATASTORE });
      flow.appendStep("Mount the NFS datastore to the second host via API",
            new MountDatastoreByApiStep(),
            new String[] { TAG_SECOND_HOST_DATASTORE });

      // Turn on HA for the cluster
      flow.appendStep("Turn on HA for cluster",
            new ReconfigureClusterByApiStep());

      // Enable Fault Tolerance on host vnic via API
      flow.appendStep("Enable Fault Tolerance Logging on first host vnic",
            new EnableFtOnVnicViaApiStep(), new String[] { TAG_FIRST_HOST });
      flow.appendStep("Enable Fault Tolerance Logging on second host vnic",
            new EnableFtOnVnicViaApiStep(), new String[] { TAG_SECOND_HOST });

      // Enable vMotion on host vnic via API
      flow.appendStep("Enable vMotion on first host vnic",
            new EnableVmotionOnVnicViaApiStep(),
            new String[] { TAG_FIRST_HOST });
      flow.appendStep("Enable vMotion on second host vnic",
            new EnableVmotionOnVnicViaApiStep(),
            new String[] { TAG_SECOND_HOST });

      // Create VM via API
      flow.appendStep("Create VM with two CPUs on the first host via API",
            new CreateVmByApiStep());

      // Power on VM via API
      flow.appendStep("Power on VM via API", new PowerOnVmByApiStep());
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      // Step 1: Navigate to VM -> Summary page
      flow.appendStep("Navigate to VM Summary tab", new VmNavigationStep());
      // Step 2: Invoke Turn On Fault Tolerance
      flow.appendStep("Invoke Turn On Fault Tolerance via More Actions",
            new LaunchTurnOnFaultToleranceWizardStep());
      // Step 3: Handle warning
      flow.appendStep("Click Yes in the pop up warning",
            new ClickYesOnWarningStep());
      // Step 4: Select the NFS datastore
      flow.appendStep("Select the NFS datastore", new FtSelectDatastoreStep(),
            new String[] { TAG_SECOND_HOST_DATASTORE });
      // Step 5: Select the second host
      flow.appendStep("Select the second host", new FtSelectHostStep(),
            new String[] { TAG_SECOND_HOST });
      // Step 6: Finish the wizard
      flow.appendStep("Click Finish in the wizard", new FtFinishWizardStep());
      // Step 7: Verify the Turn On Fault Tolerance task for primary VM
      flow.appendStep("Verifying Turn On Fault Tolerance task via UI",
            new VerifyTaskByUiStep(), new String[] { TAG_TURN_ON_FT });
      // Step 8: Verify task for Start Fault Tolerance Secondary VM
      flow.appendStep(
            "Verifying Start Fault Tolerance Secondary VM task via UI",
            new VerifyTaskByUiStep(), new String[] { TAG_TURN_ON_SECONDARY });
      // Step 9: Check the VM Fault Tolerance portlet
      flow.appendStep("Check the VM Fault Tolerance portlet via UI",
            new CheckFtPortletStep(), new String[] { TAG_SECOND_HOST,
                  TAG_FAULT_TOLERANCE_SPEC });
      // Step 10: Navigate to Cluster > VMs > Virtual Machines
      flow.appendStep("Navigate to Cluster > VMs > Virtual Machines",
            new ClusterNavigationStep());
      // Step 11: Check the VM names
      flow.appendStep("Check the VM names contain PRIMARY and SECONDARY",
            new CheckFtVmNamesStep());
   }
}