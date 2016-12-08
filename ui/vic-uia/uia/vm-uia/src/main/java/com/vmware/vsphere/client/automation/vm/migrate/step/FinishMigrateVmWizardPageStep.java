package com.vmware.vsphere.client.automation.vm.migrate.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.vm.migrate.MigrateVmFlowStep;

public class FinishMigrateVmWizardPageStep extends MigrateVmFlowStep {

   @Override
   public void execute() throws Exception {
      new WizardNavigator().waitForLoadingProgressBar();
      boolean finishWizard = new WizardNavigator().finishWizard();
      verifyFatal(TestScope.BAT, finishWizard, "Verify wizard is closed");
      // Wait for recent task to complete
      new BaseView().waitForRecentTaskCompletion();
   }
}
