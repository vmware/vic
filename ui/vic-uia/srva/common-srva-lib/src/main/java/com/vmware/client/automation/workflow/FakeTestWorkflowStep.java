/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow;

import com.vmware.client.automation.workflow.common.WorkflowSpec;

/**
 * The purpose of the class is to provide empty implementation so the legacy steps
 * are TestWorkflowStep.
 */
public class FakeTestWorkflowStep extends CommonWorkflowStep {

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void execute() throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void clean() throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void logErrorInfo() {
      // TODO Auto-generated method stub
   }

}
