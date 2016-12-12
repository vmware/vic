package com.vmware.vsphere.client.automation.common.view;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.Condition;

public enum YesNoDialog {
   CONFIRMATION(
         "confirmationDialog",
         "confirmationDialog/YES",
         "confirmationDialog/NO"
      ),
   CONFIRMATION_OK(
         "confirmationDialog",
         "confirmationDialog/OK",
         "confirmationDialog/OK"
      ),
   WARNING_YN(
         "className=AlertForm",
         "className=AlertForm/YES",
         "className=AlertForm/NO"
      ),
   AUTHENTICITY(
         "className=Alert",
         "className=Alert/YES",
         "className=Alert/NO"
      ),
   ERROR(
         "className=Alert",
         "className=Alert/OK",
         "className=Alert/OK"
         );

   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;
   private static final String TEXT_ID = "className=UITextField";

   private String dialogId;
   private String yesButtonId;
   private String noButtonId;

   private YesNoDialog(String dialogId, String yesButtonId, String noButtonId) {
      this.dialogId = dialogId;
      this.yesButtonId = yesButtonId;
      this.noButtonId = noButtonId;
   }

   /**
    * Click the Yes button of the Yes/No dialog.
    */
   public void clickYes() {
      UI.component.click(yesButtonId);
   }

   /**
    * Click the No button in the Yes/No dialog.
    */
   public void clickNo() {
      UI.component.click(noButtonId);
   }

   /**
    * Gets Yes/No dialog text.
    *
    * @return  text listed in the dialog
    */
   public String getText(){
      return (UI.component.property.get(Property.TEXT, dialogId + "/" + TEXT_ID));
   }

   /**
    * Method that waits for the Yes/No dialog to be loaded on the screen.
    * @return true if the dialog is found on the screen.
    */
   public boolean waitForDialogToLoad() {
      return UI.condition.isFound(dialogId).await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Checks if the dialog is shown.
    */
   public boolean isVisible() {
      Condition condition = UI.condition.isFound(dialogId);
      condition.estimate();
      return !condition.isUnchecked() && condition.isTrue();
   }

   /**
    * Gets Yes/No dialog title.
    *
    * @return  the title of the dialog
    */
   public String getTitle() {
      return UI.component.property.get(Property.TITLE, dialogId);
   }
}
