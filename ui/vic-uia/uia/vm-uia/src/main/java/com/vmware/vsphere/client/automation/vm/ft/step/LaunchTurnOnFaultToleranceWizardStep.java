/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.client.automation.components.menu.ContextMenuBuilder;
import com.vmware.client.automation.components.menu.MenuNode;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.util.FaultToleranceMessages;
import com.vmware.vsphere.client.automation.vm.ft.view.FtWarningPage;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * Launches Turn On Fault Tolerance Wizard from context menu of a VM
 */
public class LaunchTurnOnFaultToleranceWizardStep extends CommonUIWorkflowStep {
   private final FtWarningPage warningPage = new FtWarningPage();
   private final static FaultToleranceMessages localizedMessages = I18n
         .get(FaultToleranceMessages.class);

   @Override
   public void execute() throws Exception {

      // Invoke Turn On Fault Tolerance
      ActionNavigator.openMoreActions();
      MenuNode rootNode = new ContextMenuBuilder().getRootMenuNode();
      MenuNode turnOnFt = rootNode.getChildByActionLabel(localizedMessages
            .TurnOnFaultToleranceMenuOption());
      turnOnFt.expandTo();
      turnOnFt.leftMouseClick();

      // Verify the warning page is open
      verifySafely(warningPage.isOpen(), "Warning page is open.");

      // Verify the warning page title
      verifySafely(warningPage.verifyTitle(),
            "Verify the warning dialog title.");
   }

   /**
    * Close the dialog if something fails
    */
   @Override
   public void clean() throws Exception {
      if (warningPage.isOpen()) {
         warningPage.clickNo();
      }
   }
}