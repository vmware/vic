/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.createvm.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Finish the create VM wizard step. It validated two things:
 * 1. The wizard is closed after clicking the finish button.
 * 2. That the recent task is completed.
 */
public class FinishCreateVmWizardPageStep extends CreateVmFlowStep {

   @Override
   public void execute() throws Exception {
      WizardNavigator wizardNavigator = new WizardNavigator();
      wizardNavigator.waitForLoadingProgressBar();
      boolean finishWizard = wizardNavigator.finishWizard();
      verifyFatal(finishWizard, "Verify wizard is closed");

      new BaseView().refreshPage();
      // Wait for recent task to complete
      boolean waitForTaskCompletion = new BaseView().waitForRecentTaskCompletion();
      verifyFatal(waitForTaskCompletion, "Verify the triggered task is completed");
   }

   /**
    * Delete the newly created vm and log if cleanup is not successful.
    */
   @Override
   public void clean() throws Exception {
      _logger.info("VM to clean after the step: " + createVmSpec.name.get());
      VmSrvApi.getInstance().deleteVmSafely(createVmSpec);
   }
}
