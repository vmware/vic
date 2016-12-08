/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.step.ClickOkSinglePageDialogStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ReconfigureVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.VerifyVmExistenceByApiStep;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.edit.step.AddDiskStep;
import com.vmware.vsphere.client.automation.vm.lib.step.LaunchEditSettingsStep;
import com.vmware.vsphere.client.automation.vm.edit.step.VerifyAddHddStep;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec.HddActionType;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Test class for edit VM in the NGC client. Executes the following test
 * work-flow: 1. Open a browser 2. Login as Administrator user 3. Navigate to
 * the cluster host 4. Create new VM via API 5. Verify via the API that the VM
 * has been created 6. Edit VM - add new HDD 7. Verify via UI that the VM has
 * been hdd added
 */
public class EditVmAddHddTest extends NGCTestWorkflow {

   private static final String TAG_INIT_CONFIG = "oldConfig";
   private static final String TAG_NEW_CONFIG = "newConfig";

   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);

      flow.appendStep("Create new test VM through the API",
            new CreateVmByApiStep(), new String[] { TAG_INIT_CONFIG });

      flow.appendStep("Verified that VM exists through API.",
            new VerifyVmExistenceByApiStep(), new String[] { TAG_INIT_CONFIG });
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBed = testbedBridge
            .requestTestbed(CommonTestBedProvider.class, true);

      // Spec for the VC
      VcSpec requestedVcSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Spec for the host
      HostSpec requestedHostSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_ENTITY);

      // Spec for the datastore
      DatastoreSpec requestedDastartoreSpec = testBed.getPublishedEntitySpec(
            CommonTestBedProvider.CLUSTER_HOST_DS_ENTITY);

      // Spec for the vm that will be edit (init configuration)
      CustomizeHwVmSpec vmSpec = SpecFactory.getSpec(CustomizeHwVmSpec.class, requestedHostSpec);
      vmSpec.datastore.set(requestedDastartoreSpec);
      vmSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.OLD);
      vmSpec.tag.set(TAG_INIT_CONFIG);

      // Spec for the HDD of the vm that to be added
      HddSpec hddSpec = SpecFactory.getSpec(HddSpec.class, vmSpec);
      hddSpec.hddCapacity.set("1.00");
      hddSpec.hddCapacityType.set("GB");
      hddSpec.name
            .set(VmUtil.getLocalizedString("vm.summary.hardware.harddisk"));
      hddSpec.vmDeviceAction.set(HddActionType.ADD);

      // Spec for the vm that will be edit (target configuration)
      CustomizeHwVmSpec editVmSpec = SpecFactory.getSpec(CustomizeHwVmSpec.class, requestedHostSpec);
      editVmSpec.name.set(vmSpec.name.get());
      editVmSpec.reconfigurableConfigVersion
            .set(ManagedEntitySpec.ReconfigurableConfigSpecVersion.NEW);
      editVmSpec.tag.set(TAG_NEW_CONFIG);
      editVmSpec.hddList.set(hddSpec);

      // Spec for the reconfigurable VMs
      ReconfigureVmSpec reconfigVmSpec = new ReconfigureVmSpec();
      reconfigVmSpec.newVmConfigs.set(vmSpec);
      reconfigVmSpec.targetVm.set(editVmSpec);

      VmLocationSpec vmLocationSpec = new VmLocationSpec(vmSpec,
            NGCNavigator.NID_VM_SUMMARY);
      vmLocationSpec.tag.set(TAG_INIT_CONFIG);

      // Spec for the edit VM task
      TaskSpec editVmTaskSpec = new TaskSpec();
      editVmTaskSpec.name
            .set(VmUtil.getLocalizedString("task.reconfigureVm.name"));
      editVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      editVmTaskSpec.target.set(vmSpec);

      VmLocationSpec vmEditSummaryLocationSpec = new VmLocationSpec(editVmSpec,
            NGCNavigator.NID_VM_SUMMARY);
      vmEditSummaryLocationSpec.tag.set(TAG_NEW_CONFIG);

      // Specs only used in the steps directly
      testSpec.add(requestedVcSpec, vmSpec, hddSpec, editVmSpec, vmLocationSpec,
            editVmTaskSpec, vmEditSummaryLocationSpec, reconfigVmSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);

      flow.appendStep("Navigate to VM in object navigator",
            new VmNavigationStep(), new String[] { TAG_INIT_CONFIG });

      flow.appendStep("Launch Edit Settings of the VM",
            new LaunchEditSettingsStep());

      flow.appendStep("Add new hdd for the VM", new AddDiskStep());

      flow.appendStep("Click OK on the 'Edit Vm Settings' dialog",
            new ClickOkSinglePageDialogStep());

      flow.appendStep("Navigate to VM > Summary tab", new VmNavigationStep(),
            new String[] { TAG_NEW_CONFIG });

      flow.appendStep("Verify the new hdd in Hardware portlet",
            new VerifyAddHddStep(), new String[] { TAG_NEW_CONFIG });

   }

   // Test added under HPQC>vSphere2016>VC_UI>p:VCUI CAT
   @Override
   @Test(description = "Edit VM > Add hard disk", groups = { BAT, CAT })
   @TestID(id = "616703")
   public void execute() throws Exception {
      super.execute();
   }
}
