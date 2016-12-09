/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.clone.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.vsphere.client.automation.common.VmGlobalActions;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.createvm.view.SelectCreationTypePage;
import com.vmware.vsphere.client.automation.vm.createvm.view.SelectNameAndFolderPage;
import com.vmware.vsphere.client.automation.vm.clone.CloneVmFlowStep;

/**
 * Launches Clone Virtual Machine Wizard from context menu of a VM
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.clone.step.LaunchCloneVmStep}
 */
@Deprecated
public class LaunchCloneVmStep extends CloneVmFlowStep {

   private static final String DIALOG_TITLE_FORMAT = VmUtil
         .getLocalizedString("cloneExistingVm.wizard.titleFormat");

   @Override
   public void execute() throws Exception {
      final SelectNameAndFolderPage selectPage = new SelectNameAndFolderPage();

      // Launch Clone VM wizard
      ActionNavigator.invokeFromActionsMenu(VmGlobalActions.AI_CLONE_VM_TO_VM);
      selectPage.waitForDialogToLoad();

      // Verify that Clone VM wizard is opened
      verifyFatal(selectPage.isOpen(), "Clone VM wizard is opened.");


      // Verify the title of the dialog
      verifySafely(TestScope.UI, new TestScopeVerification() {

         @Override
         public boolean verify() throws Exception {
            return selectPage.getTitle().contains(DIALOG_TITLE_FORMAT);
         }
      }, "Verifying the title of Clone Existing VM wizard.");
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