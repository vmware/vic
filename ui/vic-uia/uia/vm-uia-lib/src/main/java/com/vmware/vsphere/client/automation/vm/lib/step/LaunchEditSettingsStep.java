/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.common.VmGlobalActions;

/**
 * Step that invoked the "Edit Settings" dialog for a specified VM. <br>
 * The step expects that the scenario has already navigated to the VM which will
 * be edit
 */
public class LaunchEditSettingsStep extends CommonUIWorkflowStep {

   @Override
   public void execute() throws Exception {
      // Invoke edit action for the dvSwicth
      ActionNavigator.invokeFromActionsMenu(VmGlobalActions.AI_EDIT_SETTINGS);

      // Wait for the dialog to load
      new SinglePageDialogNavigator().waitForDialogToLoad();

      // Verify that Edit Settings dialog is opened
      verifyFatal(TestScope.BAT, new SinglePageDialogNavigator().isOpen(),
            "'Edit Settings' dialog is opened.");
   }
}