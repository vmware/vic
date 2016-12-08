/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test;

import java.util.List;

import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStep;

/**
 * The interface defines the API of test workflow step.
 */
public interface TestWorkflowStep extends WorkflowStep {

   /**
    * Use this method to complete the test spec initialization or other
    * data related preparation. It is called before the <code>execute</code>
    * method of the same step.
    *
    * The implementation of this method is optional.
    *
    * Note: Do not use it to perform mutation operations.
    */
   public void prepare(WorkflowSpec filteredWorkflowSpec /*TestBedConnectionsBridge ; TestBedProvidersBridge*/)
         throws Exception;

   /**
    * Use this method to perform the actual execution and, if required,
    * verification of the test step. It is called after the
    * <code>prepare</code> method of the same step.
    *
    * The implementation of this method is required.
    */
   public void execute() throws Exception;


   /**
    * Use this method to revert the changes committed to the backend in
    * the <code>execute</code> method of the same step.
    *
    * The implementation of this method is required only of the
    * <code>execute</code> method is design to commit a change to the
    * backend.
    *
    * The method is called in the clean-up phase backwards on each workflow
    * step.
    */
   public void clean() throws Exception;

   /**
    * Return list of the failed non-fatal validations.
    *
    * @return list of RuntimeException objects.
    */
   public List<RuntimeException> getFailedValidations();

   /**
    * Set the step run scope defined by respective test scenario.
    * @param stepScope
    */
   public void setStepTestScope(TestScope stepScope);

   /**
    * Get the step scope defined during test scenario.
    * @return
    */
   public TestScope getStepTestScope();

   /**
    * The method is invoked by the test controller if exception is thrown during step execution.
    * Its goal is to provide more detailed information about the state of the test/setup during the test step execution.
    * For example the CommonUIWorkflowStep take screenshot if exception is thrown.
    */
   public void logErrorInfo();

}
