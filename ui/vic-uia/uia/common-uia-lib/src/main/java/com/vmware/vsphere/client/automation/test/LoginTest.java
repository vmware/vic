/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.common.spec.BrowserSpec;
import com.vmware.client.automation.common.step.ConnectBrowserStep;
import com.vmware.client.automation.common.step.LoginStep;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;

/**
 * Implementation of test for logging in NGC Client.
 * The purpose of this test is to check the connectivity to the selenium grid nodes
 * Test work-flow:
 * 1. Open a browser
 * 2. Logs in as admin user
 */
public class LoginTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      // Add the BrowserSpec
      BrowserSpec browserSpec = new BrowserSpec();
      browserSpec.url.set(testBed.getNGCURL());
      spec.links.add(browserSpec);

      // Add the UserSpec
      UserSpec userSpec = testBed.getAdminUser();
      spec.links.add(userSpec);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {
      // Connect Browser Step
      composition.appendStep(new ConnectBrowserStep());

      // Login Step
      composition.appendStep(new LoginStep());
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
