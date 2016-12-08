/** Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.workflow.test.steps;

import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Step that fails fatal verification in the clean up.
 */
public class CleanVerificationFailStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      verifySafely(true, "Pass verification");
   }

   @Override
   public void clean() throws Exception {
      verifyFatal(false, "Fail clean verification");
   }

}
