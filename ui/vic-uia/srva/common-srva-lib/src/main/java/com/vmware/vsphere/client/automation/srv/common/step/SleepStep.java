/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SleepSpec;

/**
 * Sleeps the thread for a specified amount of time.
 */
public class SleepStep extends BaseWorkflowStep {

   private SleepSpec spec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) {
      spec = filteredWorkflowSpec.get(SleepSpec.class);

      ensureNotNull(spec, "A required SleepSpec is not provided.");
   }

   @Override
   public void execute() throws Exception {
      Thread.sleep(spec.sleepTimeInMillis.get());
   }

}
