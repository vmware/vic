/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowContext;
import com.vmware.client.automation.workflow.common.WorkflowPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;

/**
 * The context class keeps the configuration, runtime state and statistics
 * data about a <code>TestWorkflow</code>.
 * 
 * In a nut shell it keeps all the state of a test workflow.
 */
public class TestWorkflowContext implements WorkflowContext {

   // Configuration
   private final Class<TestWorkflow> _testWorkflowClass;

   // test scenario workflow
   private TestWorkflow _testWorkflow;

   // Test spec describing the Inventory Structure used by the test
   private WorkflowSpec _testSpec;

   private WorkflowStepsSequence<TestWorkflowStepContext> _prepFlow;

   private WorkflowStepsSequence<TestWorkflowStepContext> _testFlow;

   // Move to configuration
   private TestScope _testScope = TestScope.FULL;

   // Test
   private WorkflowPhaseState _prerequisiteState;
   private List<Throwable> _prerequisiteErrors;

   private WorkflowPhaseState _executeState;
   private List<Throwable> _executeErrors;

   private WorkflowPhaseState _cleanState;
   private List<Throwable> _cleanUpErrors;

   /**
    * Default constructor.
    * Initialize all members, so the getters are accessible from this point on.
    * 
    * @param providerWorkflow
    *       Reference to the provider workflow
    */
   public TestWorkflowContext(Class<TestWorkflow> testWorkflowClass) {
      _testWorkflowClass = testWorkflowClass;
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
   public void createWorkflowInstance() throws InstantiationException,
   IllegalAccessException {
      _testWorkflow = null;
      _testSpec = null;
      _prepFlow = null;
      _testFlow = null;

      _testWorkflow = _testWorkflowClass.newInstance();
      _testSpec = new WorkflowSpec();
      _prepFlow = new WorkflowStepsSequence<TestWorkflowStepContext>();
      _testFlow = new WorkflowStepsSequence<TestWorkflowStepContext>();

      _executeErrors = new ArrayList<Throwable>();
      _prerequisiteErrors = new ArrayList<Throwable>();
      _cleanUpErrors = new ArrayList<Throwable>();
   }


   /**
    * Returns a reference to the test workflow class.
    */
   public TestWorkflow getTestWorkflow() {
      return _testWorkflow;
   }

   /**
    * Returns reference to a root spec.
    */
   public WorkflowSpec getTestSpec() {
      return _testSpec;
   }

   /**
    * Returns a reference to the prerequisite steps container.
    */
   public WorkflowStepsSequence<TestWorkflowStepContext> getPrepFlow() {
      return _prepFlow;
   }

   /**
    * Returns a reference to the test steps container.
    */
   public WorkflowStepsSequence<TestWorkflowStepContext> getTestFlow() {
      return _testFlow;
   }

   /**
    * Returns the scope on which the test runs.
    */
   public TestScope getTestScope() {
      return _testScope;
   }

   /**
    * Set the scope on which the tet will run. By default it's full;
    * 
    * @param testScope
    *      One a <code>TestScope</code> enum value;
    */
   public void setTestScope(TestScope testScope) {
      _testScope = testScope;
   }

   /**
    * Set test scenario execution state - PASS, FAIL or SKIP.
    * @param state
    */
   public void setExecuteState(WorkflowPhaseState state) {
      this._executeState = state;
   }

   /**
    * Get test scenario execution state - PASS, FAIL or SKIP.
    * @return
    */
   public WorkflowPhaseState getExecuteState() {
      return this._executeState;
   }

   /**
    * Set test prerequisite execution state - PASS, FAIL.
    * @param state
    */
   public void setPrerequisiteState(WorkflowPhaseState state) {
      this._prerequisiteState = state;
   }

   /**
    * Get test scenario prerequisite state - PASS, FAIL.
    * @return
    */
   public WorkflowPhaseState getPrerequisiteState() {
      return this._prerequisiteState;
   }

   /**
    * Set test scenario clean up state - PASS, FAIL.
    * @param state
    */
   public void setCleanState(WorkflowPhaseState state) {
      this._cleanState = state;
   }

   /**
    * Get test scenario clean up state - PASS, FAIL.
    * @return
    */
   public WorkflowPhaseState getCleanState() {
      return this._cleanState;
   }


   // Runtime
   // TODO: Use this spot to place anything the controller will need to
   // keep when executing the test.
   public void addExecutionErrors(Throwable error) {
      _executeErrors.add(error.getCause() == null ? error : error.getCause());
   }

   public List<Throwable> getExecutionErrors() {
      return this._executeErrors;
   }

   public void addPrerequisiteErrors(Throwable error) {
      _prerequisiteErrors.add(error.getCause() == null ? error : error.getCause());
   }

   public List<Throwable> getPrerequisiteErrors() {
      return this._prerequisiteErrors;
   }

   public void addCleanUpErrors(Throwable error) {
      _cleanUpErrors.add(error.getCause() == null ? error : error.getCause());
   }

   public List<Throwable> getCleanUpErrors() {
      return this._cleanUpErrors;
   }


   // Report
   // TODO: Report data for the test run. Need to be designed first.
}
