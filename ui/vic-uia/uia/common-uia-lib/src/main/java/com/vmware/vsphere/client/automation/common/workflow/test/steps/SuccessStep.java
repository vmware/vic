/** Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.workflow.test.steps;

import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Step that does pass fatal and safely verification.
 */
public class SuccessStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      verifyFatal(true, "Pass fatal verification");
      verifySafely(true, "Pass safely verification");
   }

   @Override
   public void clean() throws Exception {
      verifyFatal(true, "Pass clean verification");
   }

}
