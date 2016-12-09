/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.vsphere.client.automation.vm.common.VmTemplateGlobalActions;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.common.step.CreateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.createvm.view.SelectCreationTypePage;

/**
 * Launches New Virtual Machine wizard by invoking "New VM from This Template"
 * context menu of a VM template.
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.step.LaunchNewVmFromThisTemplateStep}
 */
@Deprecated
public class LaunchNewVmFromThisTemplateStep extends CreateVmFlowStep {

   private static final String DIALOG_TITLE_FORMAT = VmUtil
         .getLocalizedString("deployFromTemplate.wizard.title.format");

   @Override
   public void execute() throws Exception {
      final SelectCreationTypePage selectPage = new SelectCreationTypePage();
      // Launch New VM wizard
      ActionNavigator
            .invokeFromActionsMenu(VmTemplateGlobalActions.AI_NEW_VM_FROM_THIS_TEMPLATE);
      selectPage.waitForDialogToLoad();

      // Verify that New VM wizard is opened
      verifyFatal(TestScope.BAT, selectPage.isOpen(),
            "New VM wizard is opened.");

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
