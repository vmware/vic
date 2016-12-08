/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test;

import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.StepPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;

/**
 * The class keep the runtime state of a test workflow step.
 * 
 **/
public class TestWorkflowStepContext extends WorkflowStepContext {

   private final TestWorkflowStep _testStep;

   private StepPhaseState _prepareState;
   private StepPhaseState _execState;
   private StepPhaseState _cleanState;

   public TestWorkflowStepContext(TestWorkflowStep testStep, TestScope testScope) {
      _testStep = testStep;

      _prepareState = StepPhaseState.BLOCKED;
      _execState = StepPhaseState.BLOCKED;
      _cleanState = StepPhaseState.BLOCKED;

      if(testScope != null) {
         _testStep.setStepTestScope(testScope);
      }
   }

   /**
    * Return the respective step test scope scope.
    * The step scope may be assigned during compose stage when adding step to
    * the test scenario.
    */
   public TestScope getTestScope() {
      return _testStep.getStepTestScope();
   }

   /**
    * Assign step test scope at runtime based on the test scope configuration
    * set for running the test.
    * @param runTestScope
    */
   public void setTestScope(TestScope runTestScope) {
      _testStep.setStepTestScope(runTestScope);
   }

   @Override
   public TestWorkflowStep getStep() {
      return _testStep;
   }


   // Runtime data

   public StepPhaseState getPrepareState() {
      return _prepareState;
   }

   public void setPrepareState(StepPhaseState state) {
      _prepareState = state;
   }

   public StepPhaseState getExecuteState() {
      return _execState;
   }

   public void setExecuteState(StepPhaseState state) {
      _execState = state;
   }

   public StepPhaseState getCleanState() {
      return _cleanState;
   }

   public void setCleanState(StepPhaseState state) {
      _cleanState = state;
   }

}
