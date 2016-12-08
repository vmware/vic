/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;

/**
 * Select migration type for the VM
 */
public class SelectMigrationTypePage extends WizardNavigator {
   private static String PAGE_TITLE = VmUtil
         .getLocalizedString("migrateVm.wizard.titleFormat");

   /**
    * Validates that the view is present on the screen, before executing any
    * actions on it.
    */
   public void validate() {
      if (!verifyPageTitle()) {
         throw new IllegalStateException("Unexpected dialog title. Expected was: "
               + PAGE_TITLE);
      }
   }

   /**
    * Checks if the Select Creation Type page is open, by verifying the page title.
    *
    * @return True if the page title corresponds to Select Creation Type page, false
    *         otherwise
    */
   public boolean verifyPageTitle() {
      WizardNavigator wizardNavigator = new WizardNavigator();
      wizardNavigator.waitForDialogToLoad();
      return wizardNavigator.getCurrentlyActivePage().equals(1);
   }
}
