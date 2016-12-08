package com.vmware.vsphere.client.automation.storage.lib.wizard.core.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * WizardStep implementation for canceling the wizard
 */
public class CancelWizardStep extends WizardStep {

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {
      wizardNavigator.cancel();
   }

}
