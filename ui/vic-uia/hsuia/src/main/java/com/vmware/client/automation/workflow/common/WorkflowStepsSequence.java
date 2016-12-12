/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.common;

import java.util.ArrayList;
import java.util.List;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStepContext;
import com.vmware.client.automation.workflow.test.TestWorkflowStep;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;

/**
 * The class is a holder of the Workflow steps. It provides the respective
 * methods to add and remove steps from a workflow.
 * @param <T>
 */
public class WorkflowStepsSequence<T extends WorkflowStepContext> {
   private final List<T> _workflowStepsContext = new ArrayList<T>();

   /**
    * Return all the defined workflow steps.
    *
    * @return
    *    A <code>List</code> of all workflow steps.
    */
   public List<T> getAllSteps() {
      List<T> workflowSteps = new ArrayList<T>();
      workflowSteps.addAll(_workflowStepsContext);
      return workflowSteps;
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * The method will do nothing if the step is already defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    *
    * @param title
    *    Optional step title.
    *
    * @param testScope
    *    Optional test scope. The step will be skipped if the test run doesn't
    *    have the required test scope.
    *
    * @param specTagsFilter
    *    Optional spec tag filter. Only specs containing at least one of the tags will
    *    be made available to the test step.
    */
   @SuppressWarnings("unchecked")
   public void appendStep(String title, WorkflowStep step, TestScope testScope, String[] specTagsFilter) {
      // TODO rkovachev: Verify that a step is not added twice by accident - e.g. sort out how to determine steps equality.

      T workflowStepContext;
      if (step instanceof ProviderWorkflowStep) {
         workflowStepContext =
               (T) new ProviderWorkflowStepContext((ProviderWorkflowStep) step);
      } else if (step instanceof TestWorkflowStep) {
         workflowStepContext = (T) new TestWorkflowStepContext((TestWorkflowStep) step, testScope);
      } else {
         throw new IllegalArgumentException();
      }

      if (!Strings.isNullOrEmpty(title)) {
         workflowStepContext.setTitle(title);
      }

      if (specTagsFilter != null) {
         workflowStepContext.setSpecTagsFilter(specTagsFilter);
      }

      _workflowStepsContext.add(workflowStepContext);
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * @param title
    *    Optional step title.
    * @param step
    *    A workflow step.
    *
    * @param testScope
    *    Optional test scope. The step will be skipped if the test run doesn't
    *    have the required test scope.
    */
   public void appendStep(String title, WorkflowStep step, TestScope testScope) {
      appendStep(title, step, testScope, null);
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * @param title
    *    Optional step title.
    * @param step
    *    A workflow step.
    * @param specTagsFilter
    *    Optional spec tag filter. Only specs containing at least one of the tags will
    *    be made available to the test step.
    */
   public void appendStep(String title, WorkflowStep step, String[] specTagsFilter) {
      appendStep(title, step, null, specTagsFilter);
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * @param title
    *    Optional step title.
    * @param step
    *    A workflow step.
    */
   public void appendStep(String title, WorkflowStep step) {
      appendStep(title, step, null, null);
   }

   /**
    * Remove a work flow step from the list. Do nothing if the step is
    * not defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    */
   public void removeStep(WorkflowStep step) {
      WorkflowStepContext stepContextForRemoval = null;

      for (T stepContext : _workflowStepsContext) {
         if (step == stepContext.getStep()) {
            stepContextForRemoval = stepContext;
            break;
         }
      }

      if (stepContextForRemoval != null) {
         _workflowStepsContext.remove(stepContextForRemoval);
      }
   }

   /**
    * Remove the specified number of steps from the end of the list.
    *
    * @param numberOfSteps number os steps to remove from the list.
    */
   public void removeLastSteps(int numberOfSteps) {
      if (numberOfSteps > 0 && numberOfSteps <= _workflowStepsContext.size()) {
         int newSize = _workflowStepsContext.size() - numberOfSteps;
         while (_workflowStepsContext.size() > newSize) {
            _workflowStepsContext.remove(_workflowStepsContext.size() - 1);
         }
      } else {
         throw new IllegalArgumentException(
               "numberOfSteps in outside the number of available steps.");
      }
   }
}
