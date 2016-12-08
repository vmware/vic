/** Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.workflow.test.steps;

import com.vmware.client.automation.workflow.BaseWorkflowStep;

/**
 * Step with failed fatal verification.
 */
public class FatalFailStep extends BaseWorkflowStep {

   @Override
   public void execute() throws Exception {
      verifyFatal(false, "Fatal fail verification");
      verifySafely(true, "Pass verification");
   }

   @Override
   public void clean() throws Exception {
      verifyFatal(true, "Pass clean verification");
   }

}
