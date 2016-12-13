/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator;

import java.util.List;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * Provides utility methods for state-less navigation in a multi-page NGC dialog.
 *
 */
public class MultiPageDialogNavigator extends BaseMultiPageDialogNavigator {

   private static final IDGroup ID_OK_BTN = IDGroup.toIDGroup("tiwoDialog$okButton");

   //---------------------------------------------------------------------------
   // Common methods for Ok-Cancel dialogs

   /**
    * Checks whether a multi-page dialog is open.
    *
    * @return True if a single-page dialog is open, false otherwise.
    */
   public boolean isOpen() {
      return UI.condition.isFound(ID_OK_BTN).
            await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Clicks the OK button and checks for validation errors on the page.
    *
    * @return True if no validation errors occur after clicking the OK button,
    *    false otherwise
    */
   public boolean clickOk() {
      return clickOk(null);
   }

   /**
    * Clicks the OK button and checks for validation errors on the page.
    *
    * @param validationErrors A string list in which validation error messages
    *    will be returned if any are displayed.
    *
    * @return True if no validation errors occur after clicking the OK button,
    *    false otherwise
    */
   public boolean clickOk(List<String> validationErrors) {
      UI.component.click(ID_OK_BTN);

      if (UI.condition.notFound(ID_OK_BTN).
              await(SUITA.Environment.getBackendJobSmall())) {
         // The OK button is not visible after waiting for the operation to complete
         return true;
      }

      List<String> pageErrors = getMessagesFromValidationPanel();

      if ((!pageErrors.isEmpty()) && (validationErrors != null)) {
         validationErrors.clear();
         validationErrors.addAll(pageErrors);
      }

      return false;
   }

   /**
    * Checks whether the OK button is enabled
    *
    * @return True if the OK button is enabled, false otherwise.
    */
   public boolean isOkBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_OK_BTN);
   }
}
