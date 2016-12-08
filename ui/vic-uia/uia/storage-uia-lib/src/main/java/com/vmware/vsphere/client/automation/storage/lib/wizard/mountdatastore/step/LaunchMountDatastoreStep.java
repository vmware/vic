package com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.step;

import com.vmware.client.automation.components.menu.ContextMenuBuilder;
import com.vmware.client.automation.components.menu.MenuNode;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.MountMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * {@link BaseWorkflowStep} implementation for launching new datastore
 */
public class LaunchMountDatastoreStep extends BaseWorkflowStep {
   private final static MountMessages localizedMessages = I18n
         .get(MountMessages.class);

   @Override
   public void execute() throws Exception {
      ActionNavigator.openMoreActions();
      MenuNode rootNode = new ContextMenuBuilder().getRootMenuNode();
      MenuNode mountNewDatastore = rootNode
            .getChildByActionLabel(localizedMessages.mountDatastoreMenuOption());
      mountNewDatastore.expandTo();
      mountNewDatastore.leftMouseClick();

      final SinglePageDialogNavigator wizardNavigator = new SinglePageDialogNavigator();
      verifyFatal(wizardNavigator.isOpen(), "Mount datastore wizard is opened.");
   }

   @Override
   public void clean() throws Exception {
      WizardNavigator wizardNavigator = new WizardNavigator();
      if (wizardNavigator.isOpen()) {
         wizardNavigator.cancel();
      }
   }

}
