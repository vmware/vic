/** Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.workflow.test.steps;

import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Step with failed safely verification.
 */
public class SafelyFailStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      verifySafely(false, "Fail safely verification");
      verifyFatal(true, "Pass verification");
   }

   @Override
   public void clean() throws Exception {
      verifyFatal(true, "Pass clean verification");
   }

}
