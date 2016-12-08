/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator;

import com.vmware.client.automation.vcuilib.commoncode.GlobalFunction;
import com.vmware.flexui.componentframework.controls.mx.custom.AnchoredDialog;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.apl.sele.SeleAPLImpl;

/**
 * Provides utility methods for state-less navigation in a pop-up dialog.
 * An example of a pop-up dialog is the one opened on top of a TIWO wizard for
 * additional customization of the objects managed by the wizard. To view a
 * concrete example of a pop-up window, navigate to the "Add virtual machines"
 * page of the "Create cloud vApp" wizard and click the  "Add from Library" button.
 *
 */
public class PopupDialogNavigator {

   //---------------------------------------------------------------------------
   // Private Constants

   private static final IDGroup ID_OK_BTN =
         IDGroup.toIDGroup("dialogPopup$okButton");

   private static final IDGroup ID_CANCEL_BTN =
         IDGroup.toIDGroup("dialogPopup$cancelButton");

   private static final String POPUP_DIALOG_ID = "dialogPopup";

   private static final String LOADING_PROGRESS_BAR_ID = "progressBar";

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;


   //---------------------------------------------------------------------------
   // Class Methods

   /**
    * Checks whether a pop-up dialog is open,
    *
    * @return True if a pop-up dialog is open, false otherwise.
    */
   public boolean isOpen() {
      return UI.condition.isFound(ID_OK_BTN).
            await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Gets the title of an open pop-up dialog.
    *
    * @return The title of the dialog.
    */
   public String getTitle() {
      AnchoredDialog dialog =
            new AnchoredDialog(POPUP_DIALOG_ID,
                  ((SeleAPLImpl)SUITA.Factory.apl()).getFlashSelenium());

      return dialog.getTitle();
   }

   /**
    * Clicks the OK button of an open pop-up dialog.
    *
    * @return True if no validation errors occur, false otherwise
    */
   public boolean clickOk() {
      isOpen();
      UI.component.click(ID_OK_BTN);
      // Wait for dialog to dissapear
      if (UI.condition.notFound(ID_OK_BTN).
            await(SUITA.Environment.getBackendJobMid())) {
       // The OK button is not visible after waiting for the operation to complete
       return true;
      }
      return false;
   }

   /**
    * Clicks the Cancel button of an open pop-up dialog.
    */
   public boolean clickCancel() {
      isOpen();
      UI.component.click(ID_CANCEL_BTN);
      // Wait for dialog to dissapear
      if (UI.condition.notFound(ID_CANCEL_BTN).
            await(SUITA.Environment.getBackendJobMid())) {
       // The Cancel button is not visible after waiting for the operation to complete
       return true;
      }
      return false;   }

   /**
    * Checks whether the OK button is enabled
    *
    * @return True if the OK button is enabled, false otherwise.
    */
   public boolean isOkBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_OK_BTN);
   }

   /**
    * Checks whether the Cancel button is enabled
    *
    * @return True if the Cancel button is enabled, false otherwise.
    */
   public boolean isCancelBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_CANCEL_BTN);
   }

   /**
    * Waits for the "Loading..." progress bar to complete.
    */
   public void waitForLoadingProgressBar() {
      GlobalFunction.waitForProgressBar(
            LOADING_PROGRESS_BAR_ID,
            SUITA.Environment.getBackendJobMid(),
            ((SeleAPLImpl)SUITA.Factory.apl()).getFlashSelenium());
   }
}
