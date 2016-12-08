/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * Provides utility methods for state-less navigation in a object selector pop-up dialog.
 * An example of such a dialog is the one opened on top of a TIWO wizard for
 * selecting objects to add to a list managed by the wizard. To view a
 * concrete example of a pop-up window, navigate to the "Create new policy"
 * wizard and click the  "Select Resources" button.
 */
public class ObjectSelectorDialog {

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   protected static final String DIALOG_PREFIX = "dialogPopup";

   protected static final String ID_GRID_SELECTED_ITEMS = DIALOG_PREFIX + "/selectedItemsList";
   protected static final String ID_BTN_REMOVE_ITEMS = DIALOG_PREFIX + "/removeItems";

   protected static final String ID_OK_BTN = DIALOG_PREFIX + "$okButton";

   private static final String ID_CANCEL_BTN = DIALOG_PREFIX + "$cancelButton";

   private static final String ID_TAB = "selectorViewTabs";

   /**
    * Navigates between the Filter/Selected Objects tabs
    * @param tabName - the name of the tab that will be switched to
    */
   public void goToTab(String tabName) {
      UI.component.value.set(tabName, ID_TAB);
   }

   /**
    * Clicks the OK button of an open pop-up dialog.
    *
    * @return True if no validation errors occur, false otherwise
    */
   public boolean clickOk() {
      UI.component.click(ID_OK_BTN);
      // TODO: Add error handling functionality

      return true;
   }

   /**
    * Clicks the Cancel button of an open pop-up dialog.
    */
   public void clickCancel() {
      UI.component.click(ID_CANCEL_BTN);
   }

   /**
    * Check if the dialog is opened.
    *
    * @return     true if the dialog is opened, false otherwise
    */
   public boolean isOpen() {
      return UI.condition.isFound(ID_OK_BTN).
            await(SUITA.Environment.getPageLoadTimeout());
   }

   public void selectAndRemoveAllItems() {
      AdvancedDataGrid grid = GridControl.findGrid(ID_GRID_SELECTED_ITEMS);

      GridControl.selectAll(grid);
      UI.component.click(ID_BTN_REMOVE_ITEMS);
   }
}
