/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test.command;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.common.WorkflowPhaseState;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.client.automation.workflow.test.TestWorkflow;
import com.vmware.client.automation.workflow.test.TestWorkflowContext;
import com.vmware.client.automation.workflow.test.TestWorkflowController;

/**
 * Execute a test end to end by given test class name and scope
 */
@WorkflowCommandAnnotation(commandName = "test-execute")
public class ExecuteTestCommand extends BaseTestCommand {

   private String _testWorkflowClassString;
   private TestScope _testScope;
   private String _testReportFilePath;

   @Override
   public void prepare(String[] params) {
      // Param0 - canonical class name
      // Param1 - test scope
      // Param2..n - [out files for the test run] - Specify exactly which they are. Note: The registry is initialized by another command.

      if (params.length < 3) {
         throw new IllegalArgumentException("");
      }

      _testWorkflowClassString = params[0];
      if (Strings.isNullOrEmpty(_testWorkflowClassString)) {
         throw new IllegalArgumentException("");
      }

      if (Strings.isNullOrEmpty(params[1])) {
         throw new IllegalArgumentException("");
      }
      _testScope = TestScope.valueOf(params[1]);

      _testReportFilePath = params[2];
      if (Strings.isNullOrEmpty(_testReportFilePath)) {
         throw new IllegalArgumentException("");
      }
   }

   @SuppressWarnings("unchecked")
   @Override
   public void execute() throws Exception{
      // Note: Required test beds are provided separately.
      // Note: Consider how instructions should be passed to skip part of the test or what-so-ever.

      // As tasks:
      // 0. Validate the class exists - it's loading the registry class map
      // 1. Create TestWorkflowContext - it's stored in the registry
      // 2. Create a controller and initialize it with the context - it's stored here.
      // 3. Ask the controller to perform the preparation
      // 4. Ask the controller to perform the run
      // 5. Collect and output the results from the run

      // As API: TestWorkflowController and RegistryCommander
      // 0. registry.registerWorkflowInstance(testClassName) - load, validate,  return TestWorkflowContext;
      // 1. Create controller and initialize it with the context
      // 2. Run controller's execute method
      // 3. registry.saveResult - what about failed tests.
      // 4. Registry.unregisterWorkflowInstance() -> Unwire all maps. (finally)

      WorkflowRegistry registry = getRegistry();
      Class<TestWorkflow> testClass = null;
      TestWorkflowContext context = null;
      TestWorkflowController controller = null;
      try {
         testClass =
               (Class<TestWorkflow>) registry.getRegisteredWorkflowClass(
                     _testWorkflowClassString,
                     TestWorkflow.class);
         context = registry.registrerWorkflowContext(testClass);
         context.setTestScope(_testScope);
         controller = TestWorkflowController.create(context); // Context set
         controller.run(); // Test scope
      } catch (Exception e) {
         context.addExecutionErrors(e);
         context.setExecuteState(WorkflowPhaseState.FAILED);
      } finally {
         // Save the result from test execution
         try {
            controller.saveReport(_testReportFilePath);
         } catch (Exception e) {
            e.printStackTrace();
         }
         registry.unregistrerWorkflowContext(context);
      }
   }
}
