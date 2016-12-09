/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test;

import java.io.BufferedWriter;
import java.io.ByteArrayOutputStream;
import java.io.FileWriter;
import java.io.IOException;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Properties;

import org.apache.commons.lang3.exception.ExceptionUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.StepPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowController;
import com.vmware.client.automation.workflow.common.WorkflowPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.ProviderBridgeImpl;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.hsua.common.datamodel.AbstractProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Controller for a test workflow.
 *
 */
public class TestWorkflowController extends WorkflowController {

   // The key used to store the test execution result
   public static final String TEST_RESULT_KEY = "test.result";
   // Keys used to store the errors in the test execution result file
   public static final String TEST_EXECUTION_ERROR_KEY = "test.error.";
   public static final String TEST_CLEAN_UP_ERROR_KEY = "clean.error.";

   private TestWorkflowContext _testWorkflowContext;

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(TestWorkflowController.class);

   private TestBedBridge _testbedBridge;

   @Override
   public WorkflowRegistry getRegistry() {
      return WorkflowRegistry.getRegistry();
   }

   /**
    * Factory to create test work flow controller for specifies test workflow
    * context.
    * @param workflowContext
    * @return
    */
   public static TestWorkflowController create(TestWorkflowContext workflowContext) {
      TestWorkflowController controller = new TestWorkflowController();
      controller._testWorkflowContext = workflowContext;
      controller._testbedBridge =
            new ProviderBridgeImpl(controller._testWorkflowContext);
      return controller;
   }

   /**
    * Save the report from test execution.
    * @param filePath path to the file in which to write the result.
    * @throws TestControllerException
    */
   public void saveReport(String filePath) throws TestControllerException {
      Properties resultValues = new Properties();
      try {
         resultValues.put(TEST_RESULT_KEY, _testWorkflowContext.getExecuteState()
               .toString());
      } catch (Exception e) {
         e.printStackTrace();
      }

      // Save errors thrown during prerequisite stage
      if(_testWorkflowContext.getPrerequisiteErrors().size() > 0) {
         _logger.error("Failures count in the prerequsited stage:" + _testWorkflowContext.getPrerequisiteErrors());
         addErrors(TEST_EXECUTION_ERROR_KEY, resultValues,
               _testWorkflowContext.getPrerequisiteErrors());
      }

      // Save errors thrown during test scenario execution stage
      if(_testWorkflowContext.getExecutionErrors().size() > 0) {
         _logger.error("Failures count in the execute stage: " + _testWorkflowContext.getExecutionErrors().size());
         addErrors(TEST_EXECUTION_ERROR_KEY, resultValues,
               _testWorkflowContext.getExecutionErrors());
      }

      // Save errors thrown during test clean up stage
      if(_testWorkflowContext.getCleanUpErrors().size() > 0) {
         _logger.error("Failures count in the clean up stage:" + _testWorkflowContext.getCleanUpErrors().size());
         addErrors(TEST_CLEAN_UP_ERROR_KEY, resultValues,
               _testWorkflowContext.getCleanUpErrors());
      }

      // Write the error and state of the test execution in the output file
      try (FileWriter fileWriter = new FileWriter(filePath);
            BufferedWriter bufferedWriter = new BufferedWriter(fileWriter)) {
         resultValues.store(bufferedWriter, null);
      } catch (IOException e) {
         e.printStackTrace();
         throw new TestControllerException("Failed to save the report!", e);
      }
   }

   /**
    * Run sequentially all steps marked as runnable in the workflow. The method
    * invokes the prepare and execute stages of the steps.
    */
   public void run() throws Exception {
      guardSingleRequest();
      try {
         // Instantiate the workflow for assembling
         runWorkflowIntantiate();

         // 0. Init spec - done
         // 1. Compose steps - done
         runWorkflowStepsCompose();

         // 2. Claim test bed instances as participating in this test bed / workflow / run 
         // (call call initSpecs(), give connections settings(), open connections/)
         ((ProviderBridgeImpl) _testbedBridge).claimTestbeds();

         // 2.1. Automatic assign of the respective service specs based on the hierarchy 
         // defined by the parent property.
         runWorkflowAssignConnectors();

         // 3. Execute prepare() for all steps - done
         runStepsSequencePrepare();

         // 4. Establish the connections as listed in the corresponding
         // ServiceSpec classes. If needed the connections might be stored in the registry.

         ((ProviderBridgeImpl) _testbedBridge).connectToTestbeds();

         // 5.1. Execute test prerequisite steps
         // 5.2. Execute test scenario steps
         executeTestSteps();

         // 5.3. Iterate through the steps and check for aggregated non-fatal
         // failures.
         calculateNonFatalValidationFails();
      } catch (Exception e) {
         // The executeTestSteps() method can not throw exception.
         // Something failed during testbed allocation or retrieving data from it - test should not be executed.
         // Set test status and add execution error.
         _testWorkflowContext.addExecutionErrors(new TestWorkflowException("Failure during resource allocation!", e));
         _testWorkflowContext.setExecuteState(WorkflowPhaseState.SKIPPED);
      } finally {
         // 6. Clean up
         cleanUp();
         // Close the connections - this is quite a mess at this point - should
         // we moved inside the sequencers.
         ((ProviderBridgeImpl) _testbedBridge).releaseTestbeds();
      }
   }

   // Private methods
   /**
    * Execute clean stage from the test workflow.
    * It executes the clean methods of the already executed steps in reverse
    * order.
    * @throws TestControllerException
    */
   private void cleanUp() throws TestControllerException {
      try {
         _logger.info("Start Test Clean up stage");
         _logger.info("=========================");
         runWorkflowClean(_testWorkflowContext.getTestFlow());
         runWorkflowClean(_testWorkflowContext.getPrepFlow());
      } catch (TestWorkflowException e) {
         // Probably damaged environment
         _testWorkflowContext.setCleanState(WorkflowPhaseState.FAILED);
         _testWorkflowContext.addCleanUpErrors(e);
         _logger.error(e.toString());
      } finally {
         _logger.info("End Test Clean up stage");
         _logger.info("=========================");
      }
   }

   /**
    * Execute test prerequisite and test scenario steps.
    * @throws TestControllerException
    */
   private void executeTestSteps() throws TestControllerException {
      // Execute test prerequisite steps
      try {
         _logger.info("Start Prerequisite Steps");
         _logger.info("=========================");
         runStepsSequenceExecute(_testWorkflowContext.getPrepFlow().getAllSteps(),
               _testWorkflowContext.getTestScope());

         _testWorkflowContext.setPrerequisiteState(WorkflowPhaseState.PASSED);

         _logger.info("Complete Prerequisite Steps");
         _logger.info("=========================");

         // Execute test scenario steps
         try {
            _logger.info("Start Scenario Steps");
            _logger.info("=========================");
            runStepsSequenceExecute(_testWorkflowContext.getTestFlow().getAllSteps(),
                  _testWorkflowContext.getTestScope());

            // Set the test global state to PASSED as no fatal validation has failed.
            _testWorkflowContext.setExecuteState(WorkflowPhaseState.PASSED);

            _logger.info("Complete Scenario Steps");
            _logger.info("=========================");

         } catch (TestWorkflowException twe) {
            // Verify fatal has failed
            _logger.error("Error during test execution!");
            _logger.error("=========================");
            _testWorkflowContext.setExecuteState(WorkflowPhaseState.FAILED);
            _testWorkflowContext.addExecutionErrors(twe);
            _logger.error(stacktraceToString(twe.getCause()));
         }

      } catch (TestWorkflowException twe) {
         _logger.error("Error during test prerequiste steps!");
         _logger.error("=========================");
         _testWorkflowContext.setPrerequisiteState(WorkflowPhaseState.FAILED);
         _testWorkflowContext.addPrerequisiteErrors(twe);
         _testWorkflowContext.setExecuteState(WorkflowPhaseState.SKIPPED);
         _logger.error(twe.getMessage());
      }
   }

   /**
    * Add the errors from the list to the provided property object using the
    * specified errorPrefix as a key and index for uniqueness.
    * @param errorPrefix
    * @param resultValues
    * @param errorList
    */
   private void addErrors(String errorPrefix, Properties resultValues,
         List<Throwable> errorList) {
      int counter = 1;
      for (Throwable throwable : errorList) {
         String errorStack = ExceptionUtils.getStackTrace(throwable);
         _logger.error("Failure number: " + counter);
         _logger.error(errorStack);
         resultValues.put(errorPrefix + counter, errorStack);
         counter++;
      }
   }

   /**
    *
    * @param numberOfStepsToRemove
    * @throws TestWorkflowException
    */
   private void runWorkflowClean(
         WorkflowStepsSequence<TestWorkflowStepContext> stepsContext)
               throws TestControllerException, TestWorkflowException {

      int numberOfNotStartedSteps = getNotExecutedStepsCount(stepsContext.getAllSteps());

      if (numberOfNotStartedSteps > 0) {
         stepsContext.removeLastSteps(numberOfNotStartedSteps);
      }
      runStepsSequenceClean(stepsContext);
   }

   private void runStepPhaseClean(TestWorkflowStepContext stepContext)
         throws TestWorkflowException {

      _logger.info("Start Step Clean up stage for: " + stepContext.getTitle() + " - "
            + stepContext.getStep().getClass());
      _logger.info("=========================");

      if (stepContext.getCleanState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      try {
         stepContext.setCleanState(StepPhaseState.IN_PROGRESS);
         stepContext.getStep().clean();
      } catch (Throwable t) {
         stepContext.setCleanState(StepPhaseState.FAILED);
         _logger.error("Failed test clean up stage on step: " + stepContext.getStep().getClass(), t);
         throw new TestWorkflowException(
               "Failed test clean up stage on step: " + stepContext.getStep().getClass(), t);
      } finally {
         _logger.info("End Step Clean up stage for: " + stepContext.getTitle() + " "
               + stepContext.getStep().getClass());
         _logger.info("=========================");
      }
      stepContext.setCleanState(StepPhaseState.DONE);
   }

   private void enableStepPhaseClean(TestWorkflowStepContext stepContext)
         throws TestControllerException {
      if (stepContext.getCleanState() == StepPhaseState.BLOCKED) {
         stepContext.setCleanState(StepPhaseState.READY_TO_START);
      } else {
         throw new TestControllerException(
               "The test attempted invalid enabling of step phase.");
      }
   }

   /**
    * Iterate the test scenario steps and check if for non-fatal validation fails.
    * If such fails are detected change the test status to fail and print the errors.
    */
   private void calculateNonFatalValidationFails() {
      boolean hasNonFatalFails = false;
      for (TestWorkflowStepContext stepContex : _testWorkflowContext.getTestFlow()
            .getAllSteps()) {
         if (stepContex.getExecuteState().equals(StepPhaseState.FAILED)) {
            if (!hasNonFatalFails) {
               _logger.debug("Print NON-FATAL validation fails");
            }
            hasNonFatalFails = true;
            List<RuntimeException> failedValidations =
                  stepContex.getStep().getFailedValidations();
            for (RuntimeException runtimeException : failedValidations) {
               _logger.debug("Step Failed Validation : "
                     + stepContex.getStep().getClass());
               _testWorkflowContext.addExecutionErrors(runtimeException);
            }
            _testWorkflowContext.setExecuteState(WorkflowPhaseState.FAILED);
         }
      }
   }

   private void runStepsSequenceClean(
         WorkflowStepsSequence<TestWorkflowStepContext> cleanStepsContext)
               throws TestWorkflowException, TestControllerException {
      List<TestWorkflowStepContext> reversedSteps =
            new ArrayList<TestWorkflowStepContext>(cleanStepsContext.getAllSteps());

      Collections.reverse(reversedSteps);

      for (TestWorkflowStepContext stepContext : reversedSteps) {
         // Note: Important - keep on cleaning to the first one.
         enableStepPhaseClean(stepContext);
         try {
            runStepPhaseClean(stepContext);
         } catch (TestWorkflowException e) {
            _testWorkflowContext.addCleanUpErrors(e);
            _testWorkflowContext.setCleanState(WorkflowPhaseState.FAILED);
         }
      }
   }


   /**
    * Run the prepare() methods of all tests in the workflow.
    *
    * Each method is given isolated copy of the workflow spec.
    * @throws Exception
    */
   private void runStepsSequencePrepare() throws Exception {
      List<TestWorkflowStepContext> prerequsiteStepsContext =
            _testWorkflowContext.getPrepFlow().getAllSteps();

      for (TestWorkflowStepContext stepContext : prerequsiteStepsContext) {
         enableStepPhasePrepare(stepContext);
         runStepPhasePrepare(stepContext); // TODO: Handle
         // exception
      }

      List<TestWorkflowStepContext> testScenarioStepsContext = _testWorkflowContext
            .getTestFlow().getAllSteps();

      for (TestWorkflowStepContext stepContext : testScenarioStepsContext) {
         enableStepPhasePrepare(stepContext);
         runStepPhasePrepare(stepContext); // TODO: Handle
         // exception
      }
   }

   private void enableStepPhasePrepare(TestWorkflowStepContext stepContext)
         throws TestControllerException {
      if (stepContext.getPrepareState() == StepPhaseState.BLOCKED) {
         stepContext.setPrepareState(StepPhaseState.READY_TO_START);
      } else {
         throw new TestControllerException(
               "The controller attempted invalid enabling of step phase.");
      }
   }

   private void runStepPhasePrepare(TestWorkflowStepContext stepContext) throws Exception {

      if (stepContext.getPrepareState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      WorkflowSpec filteredTestSpec = null;
      try {
         filteredTestSpec =
               (WorkflowSpec) getDeepClonedSpec(
                     _testWorkflowContext.getTestSpec(),
                     stepContext.getSpecTagsFilter());
         _logger.debug("Preapare STEP: " + stepContext.getTitle());
      } catch (Exception e) {
         throw new TestControllerException(
               "runStepPhasePrepare getDeepClonedSpec failure", e);
      }

      try {
         stepContext.setPrepareState(StepPhaseState.IN_PROGRESS);
         stepContext.getStep().prepare(filteredTestSpec);
      } catch (Exception e) {
         stepContext.setPrepareState(StepPhaseState.FAILED);
         throw e;
      }

      stepContext.setPrepareState(StepPhaseState.DONE);
   }

   // TODO rkovachev: move the logic for the ManagedEntitySpec to another class
   // and inject it in the class - as the other team might not need such entity
   private void runWorkflowAssignConnectors() throws TestWorkflowException {
      WorkflowSpec testSpec = _testWorkflowContext.getTestSpec();
      List<BaseSpec> specList = testSpec.links.getAll(BaseSpec.class);

      List<BaseSpec> processedSpecs = new ArrayList<BaseSpec>();
      for (BaseSpec spec : specList) {
         deepAssignServiceSpec(spec, processedSpecs);
      }
   }

   private void assignParentService(ManagedEntitySpec spec) {

      if (!spec.parent.isAssigned() || spec.parent.get() == null) {
         // Container spec- vc, host and etc.
         if (spec.service.get() == null) {
            throw new RuntimeException("Spec " + spec
                  + " is container and has no service assigned!");
         }
      } else {
         assignParentService(spec.parent.get());
         spec.service.set(spec.parent.get().service.get());
      }
   }

   /**
    * Assign service spec to all properties of this spec and linked ManagedEntitySpecs.
    * This method keeps list of all processed specs in order not to
    * fell in infinite loop.
    *
    * @param sourceSpec       spec that will be traversed
    * @param processedSpecs   list of already processed specs
    */
   private void deepAssignServiceSpec(BaseSpec sourceSpec,
         List<BaseSpec> processedSpecs) {

      // Skip already processed specs. This is necessary because spec links
      // provide automatically references to their container specs.
      if (processedSpecs.contains(sourceSpec)) {
         return;
      }

      processedSpecs.add(sourceSpec);
      _logger.debug("About to assign service spec to: " + sourceSpec);

      // process current spec
      if (sourceSpec instanceof ManagedEntitySpec) {
         assignParentService((ManagedEntitySpec) sourceSpec);
      }

      // check all spec properties
      // all of the EntityManagedSpec's should be processed
      for (Field property : sourceSpec.getPropertyFields()) {
         if (!"parent".equals(property.getName())) {
            Object fieldObject = null;
            try {
               AbstractProperty<?> abstrProperty = (AbstractProperty<?>) property.get(sourceSpec);
               if (abstrProperty.isAssigned()) {
                  fieldObject = abstrProperty.get();
               }
            } catch (Exception e) {
               _logger.error(
                     String.format(
                           "Could not acces a spec prperty %s: %s",
                           property.getName(),
                           sourceSpec
                        )
                  );
               throw new RuntimeException("Could not access a spec property");
            }

            // if the field is of type BaseSpec process it
            if (fieldObject instanceof BaseSpec) {
               _logger.debug(
                     String.format(
                           "About to assign service spec to the %s property of %s",
                           property.getName(),
                           sourceSpec
                        )
                  );
               deepAssignServiceSpec((BaseSpec) fieldObject, processedSpecs);
            }
         }
      }

      // check all linked specs
      List<ManagedEntitySpec> linkedSourceSpecs = sourceSpec.links.getAll(ManagedEntitySpec.class);
      for (BaseSpec linkedSourceSpec : linkedSourceSpecs) {
         _logger.debug(
               String.format(
                     "About to traverse %s links and to assign service specs",
                     sourceSpec
                  )
            );
         deepAssignServiceSpec(linkedSourceSpec, processedSpecs);
      }
   }

   /**
    * Calls the workflow to build the list of its steps.
    */
   private void runWorkflowStepsCompose() throws TestWorkflowException {
      TestWorkflow workflow = _testWorkflowContext.getTestWorkflow();

      try {
         workflow.composePrereqSteps(_testWorkflowContext.getPrepFlow());
         // TODO: Add validation that the workflow has at least one step
      } catch (Exception e) {
         throw new TestWorkflowException("Error composing prerequisite test steps.", e);
      }

      try {
         workflow.composeTestSteps(_testWorkflowContext.getTestFlow());
         // TODO: Add validation that the workflow has at least one step
      } catch (Exception e) {
         throw new TestWorkflowException("Error composing test scenario steps.", e);
      }
   }

   /**
    * Instantiates the workflow and related data and stores it in context.
    *
    * @throws TestControllerException
    *            The exception is triggered if a system error occurs.
    */
   private void runWorkflowIntantiate() throws TestControllerException,
   TestWorkflowException {

      try {
         _testWorkflowContext.createWorkflowInstance();
      } catch (InstantiationException | IllegalAccessException e1) {
         throw new TestControllerException(
               "The controller cannot instantiate the test workflow", e1);
      }

      //      ProviderWorkflow workflow = _context.getProviderWorkflow();
      TestWorkflow workflow = _testWorkflowContext.getTestWorkflow();

      try {
         workflow.initSpec(_testWorkflowContext.getTestSpec(), _testbedBridge);
      } catch (Exception e) {
         throw new TestWorkflowException("Error initializing provider specs.", e);
      }
   }


   /**
    * Returns the number of steps on which the assembling was not performed.
    *
    * @return Number of steps on which the assembling was not started.
    */
   private int getNotExecutedStepsCount(List<TestWorkflowStepContext> stepsContext) {
      int count = 0;
      for (TestWorkflowStepContext stepContext : stepsContext) {
         if (stepContext.getExecuteState() == StepPhaseState.BLOCKED) {
            // TODO: Review the step phase states.
            count++;
         }
      }
      return count;
   }


   private void runStepPhaseExecute(TestWorkflowStepContext stepContext)
         throws TestWorkflowException {

      _logger.info("Start Exec Step: " + stepContext.getTitle() + " - "
            + stepContext.getStep().getClass());
      _logger.info("=========================");
      if (stepContext.getExecuteState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      try {
         stepContext.setExecuteState(StepPhaseState.IN_PROGRESS);
         stepContext.getStep().execute();
      } catch (Throwable t) {
         stepContext.setExecuteState(StepPhaseState.FAILED);
         _logger.error("======= FAILURE IN STEP: "
               + stepContext.getStep().getClass());
         // Invoke step log collector
         try {
            stepContext.getStep().logErrorInfo();
         } catch(Throwable te) {
            _logger.error("Failed to collect step error log:" + te.getMessage());
         }
         throw new TestWorkflowException(t.getMessage(), t);
      }

      if (stepContext.getStep().getFailedValidations().size() > 0) {
         stepContext.setExecuteState(StepPhaseState.FAILED);
      } else {
         stepContext.setExecuteState(StepPhaseState.DONE);
      }
      _logger.info("End Exec Step: " + stepContext.getTitle());
      _logger.info("=========================");
   }

   private void skipStepPhaseExecute(TestWorkflowStepContext stepContext) {
      _logger.info("Skip Step : " + stepContext.getTitle() +
            " due to test scope " + stepContext.getTestScope().toString() + " "
            + stepContext.getStep().getClass());
      _logger.info("=========================");

      stepContext.setExecuteState(StepPhaseState.DONE);
   }


   private void runStepsSequenceExecute(List<TestWorkflowStepContext> stepsContex, TestScope testScope)
         throws TestWorkflowException, TestControllerException {

      for (TestWorkflowStepContext stepContext : stepsContex) {
         // Check if test step requires compatible scope with the
         // one the test in run on in order to determine if the
         // step should be executed.
         if (stepContext.getTestScope().getScopeNumber() <= testScope.getScopeNumber()) {
            // Set the test scope to the step.
            // stepContext.setTestScope(testScope);

            enableStepPhaseExecute(stepContext);
            runStepPhaseExecute(stepContext);
            if (stepContext.getExecuteState().equals(StepPhaseState.FAILED)) {
               _testWorkflowContext.setExecuteState(WorkflowPhaseState.FAILED);
            }

         } else {
            // Skip the step if it requires higher scope than the test
            // is run on.
            skipStepPhaseExecute(stepContext);
         }


      }
   }

   private void enableStepPhaseExecute(TestWorkflowStepContext stepContext)
         throws TestControllerException {
      if (stepContext.getExecuteState() == StepPhaseState.BLOCKED) {
         stepContext.setExecuteState(StepPhaseState.READY_TO_START);
      } else {
         throw new TestControllerException(
               "The controller attempted invalid enabling of step phase.");
      }
   }

   private String stacktraceToString(Throwable e) {
      ByteArrayOutputStream baos = new ByteArrayOutputStream();
      PrintStream ps = new PrintStream(baos);
      e.printStackTrace(ps);
      ps.close();
      return baos.toString();
   }

}
