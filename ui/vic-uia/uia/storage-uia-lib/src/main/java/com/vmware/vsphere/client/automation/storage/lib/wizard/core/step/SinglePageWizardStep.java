package com.vmware.vsphere.client.automation.storage.lib.wizard.core.step;

import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;

public abstract class SinglePageWizardStep extends EnhancedBaseWorkflowStep {

   @Override
   public final void execute() throws Exception {
      SinglePageDialogNavigator wizardNavigator = new SinglePageDialogNavigator();
      wizardNavigator.waitForLoadingProgressBar();

      executeWizardOperation(wizardNavigator);
   }

   /**
    * Executes the step operation for the wizard
    *
    * @param wizardNavigator
    *           the WizardNavigator instance
    * @throws Exception
    */
   protected abstract void executeWizardOperation(
         SinglePageDialogNavigator wizardNavigator) throws Exception;

}
