/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.createvm.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.common.VmGlobalActions;
import com.vmware.vsphere.client.automation.vm.lib.createvm.view.SelectCreationTypePage;

/**
 * Launches New Virtual Machine Wizard from contexst menu of a managed entity
 */
public class LaunchNewVmStep extends CreateVmFlowStep {

   private static final String DIALOG_TITLE_FORMAT = CommonUtil
         .getLocalizedString("createVm.wizard.title.format");

   @Override
   public void execute() throws Exception {
      final SelectCreationTypePage selectPage = new SelectCreationTypePage();
      // Launch New VM wizard
      ActionNavigator.invokeFromActionsMenu(VmGlobalActions.AI_NEW_VM);
      selectPage.waitForDialogToLoad();

      // Verify that New VM wizard is opened
      verifyFatal(TestScope.BAT, selectPage.isOpen(), "New VM wizard is opened.");

      // Verify the title of the dialog
      verifySafely(TestScope.UI, new TestScopeVerification() {

         @Override
         public boolean verify() throws Exception {
            return selectPage.getTitle().contains(DIALOG_TITLE_FORMAT);
         }
      }, "Verifying the title of New VM wizard.");
   }

   /**
    * Close the wizard if something fails
    */
   @Override
   public void clean() throws Exception {
      SelectCreationTypePage selectPage = new SelectCreationTypePage();
      if (selectPage.isOpen()) {
         selectPage.cancel();
      }
   }
}
