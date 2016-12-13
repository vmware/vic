/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.vsphere.client.automation.storage.lib.core.DatastoreGlobalActions;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.SelectDatastoreTypePage;

/**
 * Launches New Datastore wizard
 */
public class LaunchNewDatastoreStep extends CommonUIWorkflowStep {

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      final WizardNavigator wizardNavigator = new WizardNavigator();

      // Launch create new datastore wizard
      ActionNavigator
            .invokeFromActionsMenu(DatastoreGlobalActions.AI_NEW_DATASTORE);
      wizardNavigator.waitForDialogToLoad();

      // Verify that create new datastore wizard is opened
      verifyFatal(wizardNavigator.isOpen(),
            "Create new datastore wizard is opened.");
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void clean() throws Exception {
      SelectDatastoreTypePage selectPage = new SelectDatastoreTypePage();
      if (selectPage.isOpen()) {
         selectPage.cancel();
      }
   }
}
