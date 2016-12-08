/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.vsphere.client.automation.common.VmGlobalActions;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.migrate.MigrateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.migrate.view.SelectMigrationTypePage;

/**
 * Launches Migrate Virtual Machine Wizard from context menu of a VM
 */
public class LaunchMigrateVmStep extends MigrateVmFlowStep {

   private static final String DIALOG_TITLE_FORMAT = VmUtil
         .getLocalizedString("migrateVm.wizard.titleFormat");

   @Override
   public void execute() throws Exception {
      final SelectMigrationTypePage selectPage = new SelectMigrationTypePage();

      // Launch Migrate VM wizard
      ActionNavigator.invokeFromActionsMenu(VmGlobalActions.AI_MIGRATE_VM);
      selectPage.waitForDialogToLoad();

      // Verify that Migrate VM wizard is opened
      verifyFatal(TestScope.BAT, selectPage.isOpen(), "Migrate VM wizard is opened.");

      // Verify the title of the dialog

      verifySafely(TestScope.UI, new TestScopeVerification() {

         @Override
         public boolean verify() throws Exception {
            return selectPage.getTitle().equals(
                  String.format(DIALOG_TITLE_FORMAT, vmSpec.name.get()));
         }
      },
            String.format(
                  "Verifying the title of Migrate VM dialog, expected %s, actual %s",
                  selectPage.getTitle(),
                  String.format(DIALOG_TITLE_FORMAT, vmSpec.name.get())));
   }

   /**
    * Close the wizard if something fails
    */
   @Override
   public void clean() throws Exception {
      SelectMigrationTypePage selectPage = new SelectMigrationTypePage();
      if (selectPage.isOpen()) {
         selectPage.cancel();
      }
   }

}
