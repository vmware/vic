/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ops;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.step.LogoutStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * Test class for login and logout of NGC Executes the following test
 * work-flow:
 * 1. Open a browser
 * 2. Login as admin user
 * 3. Logout
 * 4. Verify the Logout is successful
 */
public class LogoutTest extends NGCTestWorkflow {

   /**
    * {@inheritDoc}
    */
   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBed = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);

      // Spec for the VC
      VcSpec requestedVcSpec = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);
      testSpec.add(requestedVcSpec);
      super.initSpec(testSpec, testbedBridge);
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);
      flow.appendStep("Verify the Logout is successful", new LogoutStep());
   }

   /**
    * {@inheritDoc}
    */
   @Override
   @Test(description = "Verifies if Logout from the NGC navigates to the login page.", groups = { BAT, CAT })
   @TestID(id = "620409")
   public void execute() throws Exception {
      super.execute();
   }
}