package com.vmware.vsphere.client.automation.storage.lib.wizard.core.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.TestScope;

/**
 * Wizard step implementation for going to the next wizard page
 */
public class GoToNextWizardPageStep extends WizardStep {

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {

      boolean navigatedToNextWizardPage = wizardNavigator.gotoNextPage();

      verifyFatal(TestScope.BAT, navigatedToNextWizardPage,
            "Verifying clicking Next in the  page");
   }

}
