package com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.step;

import com.vmware.client.automation.assertions.FalseAssertion;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.SinglePageWizardStep;

/**
 * {@link SinglePageWizardStep} implementation
 */
public class CancelMountDatastoreWizardStep extends SinglePageWizardStep {

   @Override
   protected void executeWizardOperation(
         SinglePageDialogNavigator wizardNavigator) throws Exception {
      wizardNavigator.waitForLoadingProgressBar();
      wizardNavigator.cancel();

      verifySafely(new FalseAssertion(wizardNavigator.isOpen(),
            "Wizard is closed"));
   }

}
