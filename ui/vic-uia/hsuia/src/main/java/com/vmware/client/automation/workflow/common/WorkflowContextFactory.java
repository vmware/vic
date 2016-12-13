/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.common;

import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowContext;
import com.vmware.client.automation.workflow.test.TestWorkflow;
import com.vmware.client.automation.workflow.test.TestWorkflowContext;

/**
 * Factory for WorkFlow - based on the WorkFlow type creates TestWorkflowContext
 * or ProviderWorkflowContext.
 */
public class WorkflowContextFactory {
   @SuppressWarnings("unchecked")
   public <WC extends WorkflowContext> WC createContext(
         Class<? extends Workflow> workflowClass) {
      if (TestWorkflow.class.isAssignableFrom(workflowClass)) {
         return (WC) new TestWorkflowContext((Class<TestWorkflow>)workflowClass);
      } else if (ProviderWorkflow.class.isAssignableFrom(workflowClass)) {
         return (WC) new ProviderWorkflowContext((Class<ProviderWorkflow>)workflowClass);
      } else {
         throw new RuntimeException(
               String.format(
                     "Cannot create workflow context for %s. The workflow type is unknown to the system and cannot be managed automatically.",
                     workflowClass.getCanonicalName()));
      }
   }
}
