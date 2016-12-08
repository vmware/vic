/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.view;

import com.vmware.client.automation.components.navigator.BaseDialogNavigator;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.srv.common.util.FtUtil;

/**
 * Warning pop up page when you try to turn on Fault Tolerance.
 */
public class FtWarningPage extends BaseDialogNavigator {
   private static final String YES_BUTTON = "yesButton";
   private static final String NO_BUTTON = "noButton";
   private static String DIALOG_TITLE = FtUtil
         .getLocalizedString("warning.title");
   private static final String DIALOG_TITLE_ID = "titleDisplay";

   /**
    * Method that verifies the warning dialog title.
    */
   public boolean verifyTitle() {
      String title = UI.component.property.get(Property.TEXT, DIALOG_TITLE_ID);
      return title.contains(DIALOG_TITLE);
   }

   /**
    * Method that clicks Yes the warning dialog title.
    */
   public void clickYes() {
      UI.component.click(YES_BUTTON);
   }

   /**
    * Method that clicks No the warning dialog title.
    */
   public void clickNo() {
      UI.component.click(NO_BUTTON);
   }

   /**
    * Checks whether the warning dialog is open.
    *
    * @return True if the warning dialog is open, false otherwise.
    */
   public boolean isOpen() {
      return UI.condition.isFound(YES_BUTTON).await(
            SUITA.Environment.getPageLoadTimeout());
   }
}