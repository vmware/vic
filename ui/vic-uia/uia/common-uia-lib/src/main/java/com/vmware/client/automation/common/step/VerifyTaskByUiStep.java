/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;

/**
 * Verifies task by UI.
 *
 * Waits for task in Recent Task bar to reach expected state or provided by
 * TaskSpec maxTimeInQueue time to expire. If maxTimeInQueue is not set,
 * SUITA.Environment .getBackendJobMid() is assumed for default timeout period.
 */
public class VerifyTaskByUiStep extends CommonUIWorkflowStep {

   protected TaskSpec _taskSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _taskSpec = filteredWorkflowSpec.get(TaskSpec.class);

      ensureNotNull(_taskSpec, "Task spec is missing.");
      ensureAssigned(_taskSpec.name, "Name is not assigned.");
      ensureAssigned(_taskSpec.target, "Target is not assigned.");
      ensureAssigned(_taskSpec.status, "Status is not assigned.");
   }

   @Override
   public void execute() throws Exception {
      boolean isTaskFound = new TasksUtil()
            .waitForRecentTaskToMatchFilter(new RecentTaskFilter(_taskSpec));

      verifyFatal(TestScope.BAT, isTaskFound, String.format(
            "Verifying task '%s' for target '%s' has reached status '%s'",
            _taskSpec.name.get(), _taskSpec.target.get(),
            _taskSpec.status.get()));
   }
}
