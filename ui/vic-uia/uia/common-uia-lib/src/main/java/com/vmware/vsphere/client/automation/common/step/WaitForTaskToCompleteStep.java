/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;

/**
 * Wait for task to complete
 */
public class WaitForTaskToCompleteStep extends CommonUIWorkflowStep{

   private TaskSpec _taskSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _taskSpec = filteredWorkflowSpec.get(TaskSpec.class);

      ensureNotNull(_taskSpec, "TaskSpec spec is missing.");
      ensureNotNull(_taskSpec.name.get(), "TaskSpec name is missing.");
      ensureNotNull(_taskSpec.target.get(), "TaskSpec target is missing.");
   }

   @Override
   public void execute() throws Exception {
      // TODO: There are plans for deprecation of the method TasksUtil.waitForTaskToComplete
      // Implementation should be changed when that happens
      // TODO: Find a way (maybe add property to TaskSpec) to take into account
      // the VC time after which the task should be started
      TasksUtil.waitForTaskToComplete(_taskSpec.name.get(),
            _taskSpec.target.get().name.get(),
            UiDelay.PAGE_LOAD_TIMEOUT.getDuration());
   }
}