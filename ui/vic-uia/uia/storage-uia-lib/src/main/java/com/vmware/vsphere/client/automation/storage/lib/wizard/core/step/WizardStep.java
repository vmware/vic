package com.vmware.vsphere.client.automation.storage.lib.wizard.core.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;

/**
 * Abstract BaseWorkflowStep implementation for wizard related steps
 */
public abstract class WizardStep extends EnhancedBaseWorkflowStep {

   @Override
   public final void execute() throws Exception {
      WizardNavigator wizardNavigator = new WizardNavigator();
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
         WizardNavigator wizardNavigator) throws Exception;

}
