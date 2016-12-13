/** Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.workflow.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.common.workflow.test.steps.FatalFailStep;
import com.vmware.vsphere.client.automation.common.workflow.test.steps.SuccessStep;

/**
 * Test that fails in the prerequisite stage.
 */
public class WorkflowPrerequisiteFailTest extends NGCTestWorkflow {

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
   }

   @Override
   @Test
   @TestID(id = "N/A")
   public void execute() throws Exception {
      super.execute();
   }

   @Override
   public void composePrereqSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      flow.appendStep("Fail prerequisite", new FatalFailStep());
   }

   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      flow.appendStep("Pass exe step 1", new SuccessStep());
   }

}
