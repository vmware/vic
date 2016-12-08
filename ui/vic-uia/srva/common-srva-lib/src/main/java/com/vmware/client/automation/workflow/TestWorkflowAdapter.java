/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow;

import java.io.BufferedReader;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import com.vmware.client.automation.workflow.command.CommandController;
import com.vmware.client.automation.workflow.common.WorkflowPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.test.TestWorkflow;
import com.vmware.client.automation.workflow.test.TestWorkflowController;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;

/**
 *  Test work flow adapter for converting the existing workflow to the
 *  TestWorkflow and to keep the TestNG compatibility for the new
 *  test command tests.
 *  NOTE:
 *  Extends BaseTest        - TesNG integration for both workflows
 *  Implements WorkflowStep - Legacy workflow
 *  Implements TestWorkflow - Test command workflow
 */
public abstract class TestWorkflowAdapter extends BaseTest implements WorkflowStep,
      TestWorkflow {

   ///
   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      // Nothing to do

   }

   @Override
   public void composePrereqSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      // Nothing to do

   }

   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      // Nothing to do

   }

   ///

   /**
    * Method to invoke the test-execute command. The implementation of the method mimics the test-execute
    * command which takes care for executing tests.
    * @throws Exception
    */
   protected void invokeTestExecuteCommand() throws Exception {
      String testResultReportFile = "testResultReportFile";
      String testExecuteCommand = "test-execute %s %s %s";

      String command =
            String.format(
                  testExecuteCommand,
                  this.getClass().getCanonicalName(),
                  _testScope,
                  testResultReportFile);

      List<String> commandsList = new ArrayList<String>();
      commandsList.add(command);
      CommandController controller = new CommandController();
      // Init test execute command
      controller.initialize(commandsList);
      // Invoke prepare method
      controller.prepare();
      // Execute the command
      controller.execute();

      // Set TestNG test status based on the saved result
      Properties reportValue = new Properties();
      try (FileReader fileReader = new FileReader(testResultReportFile);
            BufferedReader bufferedReader = new BufferedReader(fileReader)) {
         reportValue.load(bufferedReader);

         assert reportValue.get(TestWorkflowController.TEST_RESULT_KEY).equals(
               WorkflowPhaseState.PASSED.name());

      } catch (FileNotFoundException e) {
         e.printStackTrace();
      } catch (IOException e) {
         e.printStackTrace();
      } catch (AssertionError ae) {
         StringBuilder testFailureReason = new StringBuilder();
         if(reportValue.get(TestWorkflowController.TEST_RESULT_KEY).equals(WorkflowPhaseState.SKIPPED.name())) {
            testFailureReason.append("SKIPPED TEST - ");
         } else {
            testFailureReason.append("FAILED TEST - ");
         }
         // Logerror message with the number of failed verification points!
         testFailureReason.append("Verification failure count is: " + (reportValue.size() - 1 ) +
               "! Check the log for them. Here is the first test failure: \n");

         // Retrieve the first error from the list, which is the fatal or the first non-fatal verification failure.
         String failureStacktrace = reportValue.getProperty(TestWorkflowController.TEST_EXECUTION_ERROR_KEY + "1");
         testFailureReason.append(failureStacktrace);
         throw new AssertionError(testFailureReason);
      }
   }
}
