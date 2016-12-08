/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * Refreshing the UI.
 */
public class SelectAllUsersTasksInRecentTasksStep extends BaseWorkflowStep {

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private static final String USER_FILTER_BTN = "portletButtons/userFilterButton";
   private static final String USER_FILTER_MENU = "userFilterMenu";
   private static final String USER_FILTER_MENU_ITEM_ALLUSERS = USER_FILTER_MENU +"/automationName=All Users\\\' Tasks[1]";

   @Override
   public void execute() throws Exception {
      openHwDevicesMenu();
      selectAllUsersTasksMenuItem();
   }

   public static void openHwDevicesMenu() {
      // Click to open Hardware devices menu
      UI.component.click(USER_FILTER_BTN);

      // Wait until the menu is loaded
      UI.condition.isFound("userFilterMenu").await(
            SUITA.Environment.getPageLoadTimeout());
   }

   public static void selectAllUsersTasksMenuItem() {
      UI.component.click(USER_FILTER_MENU_ITEM_ALLUSERS);
   }

}
