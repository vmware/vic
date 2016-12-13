/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test;

import com.vmware.client.automation.workflow.common.Workflow;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;

/**
 * Implement this interface to define a test case scenario.
 */
public interface TestWorkflow extends Workflow {

   /**
    * Implement this method to initialize the workflow test spec.
    * 
    * @param testSpec
    * 	Reference to the container of all test specs.
    * 
    * @param testbedGateway
    * 	Access point for requesting and consuming test bed resources.
    */
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge);

   /**
    * Implement this method to build the list of the workflow steps building the
    * prerequisites. They will be executed in the order they are provided.
    * 
    * @param flow
    * 	Reference to the steps container.
    */
   public void composePrereqSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow);

   /**
    * Implement this method to build the list of the workflow steps related to
    * the actual test. They will be executed in the order they are provided.
    * 
    * @param flow
    *		Reference to the steps container.
    */
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow);
}
