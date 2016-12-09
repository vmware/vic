/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.flexui.componentframework.controls.mx.custom.AnchoredDialog;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.AssertionFail;

/**
 * Provides methods common for NGC dialog components: single-page and multi-page
 * dialogs and wizards. The class doesn't represent any specific NGC dialog
 * navigator and is supposed to be sub-classed by specific dialog navigators.
 * That is why only access from the current package is allowed.
 */
//
public class BaseDialogNavigator {

   protected static final Logger _logger =
         LoggerFactory.getLogger(BaseDialogNavigator.class);

   // ---------------------------------------------------------------------------
   // Common dialog buttons IDs

   private static final IDGroup ID_CANCEL_BTN = IDGroup
         .toIDGroup("tiwoDialog$cancelButton");

   private static final IDGroup ID_MINIMIZE_TO_TIWO_BTN = IDGroup
         .toIDGroup("tiwoDialog$minimizeButton");

   private static final IDGroup ID_OPEN_HELP_BTN = IDGroup
         .toIDGroup("tiwoDialog$helpButton");

   // ---------------------------------------------------------------------------
   // Validation errors IDs

   private static final IDGroup ID_CLOSE_ERROR_BTN = IDGroup
         .toIDGroup("tiwoDialog/_closeButton");

   private static final IDGroup ID_MORE_ERROR_BTN = IDGroup
         .toIDGroup("tiwoDialog/_moreButton");

   private static final IDGroup ID_ERROR_LABEL = IDGroup
         .toIDGroup("tiwoDialog/_errorLabel");

   private static final IDGroup ID_GLOBAL_WIZARD_ERROR = IDGroup
         .toIDGroup("tiwoDialog/error_0");


   private static final String ERROR_LABEL_FORMAT = "_errorLabel[%s]";

   private static final String TIWO_DIALOG_ID = "tiwoDialog";

   private static final String TIWO_TASK_LIST_ID = "tiwoListView";

   private static final String TIWO_TASK_DESCRIPTION_FORMAT = "dataProvider.%d.description";

   private static final String TIWO_TASK_ID_FORMAT = "itemDescription[%d]";

   private static final String DIALOG_INITIALIZING_PROGRESS_BAR_ID = "tiwoDialog/vboxInfo/progressInfo";

   private static final String PAGE_LOADING_PROGRESS_BAR_ID = "tiwoDialog/progressBar";

   private static final String APPLY_SAVED_DATA_PROGRESS_BAR_ID =
         "animatedLoadingProgressBar";

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   // ---------------------------------------------------------------------------
   // Common methods for all NGC dialogs

   /**
    * Clicks the Cancel button
    */
   public void cancel() {
      UI.component.click(ID_CANCEL_BTN);
   }

   /**
    * Checks whether the Cancel button is enabled.
    *
    * @return true if, the Cancel button in enabled, false otherwise
    */
   public boolean isCancelBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_CANCEL_BTN);
   }

   /**
    * Gets the title of the dialog.
    *
    * @return String representing the title of the dialog
    */
   public String getTitle() {
      AnchoredDialog anchoredDialog = new AnchoredDialog(TIWO_DIALOG_ID,
            BrowserUtil.flashSelenium);
      return anchoredDialog.getTitle();
   }

   /**
    * Clicks the Help button
    */
   public void openHelp() {
      UI.component.click(ID_OPEN_HELP_BTN);
   }

   /**
    * Checks if the help button is enabled.
    */
   public boolean isHelpButtonEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_OPEN_HELP_BTN);
   }

   // ---------------------------------------------------------------------------
   // Methods that provide support for "Work in Progress" (TIWO) portlet

   /**
    * Minimize the dialog to the TIWO widget.
    */
   public void minimize() {
      UI.component.click(ID_MINIMIZE_TO_TIWO_BTN);
   }

   /**
    * Checks if the TIWO minimize button is enabled.
    */
   public boolean isMinimizeButonEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_MINIMIZE_TO_TIWO_BTN);
   }

   /**
    * Restores a dialog from the TIWO. Note that if several dialogs with the
    * same title are present in the TIWO portlet, the latest one is restored.
    *
    * @param dialogTitle
    *           The dialog's title
    *
    */
   public boolean restore(String dialogTitle) {

      com.vmware.flexui.componentframework.controls.mx.List tiwoTasks =
            new com.vmware.flexui.componentframework.controls.mx.List(
                  TIWO_TASK_LIST_ID,
                  BrowserUtil.flashSelenium
                  );

      if (!UI.condition.isFound(TIWO_TASK_LIST_ID).await(
            SUITA.Environment.getPageLoadTimeout())) {
         _logger.error("The TIWO tasks list is not available");
         return false;
      }

      int numberOfTasks = Integer.parseInt(tiwoTasks
            .getNumElements(SUITA.Environment.getPageLoadTimeout()));

      int taskIndex = -1;

      // The latest TIWO task is at index 0
      for (int i = 0; i < numberOfTasks; i++) {
         String taskDescription = tiwoTasks.getProperty(String.format(
               TIWO_TASK_DESCRIPTION_FORMAT, i));
         if (taskDescription.equals(dialogTitle)) {
            taskIndex = i;
            break;
         }
      }

      if (taskIndex == -1) {
         _logger.error(String.format("Unable to find %s in the TIWO portlet",
               dialogTitle));
         return false;
      }

      IDGroup taskId = IDGroup.toIDGroup(String.format(TIWO_TASK_ID_FORMAT,
            taskIndex));
      UI.component.click(taskId);

      return true;
   }

   // ---------------------------------------------------------------------------
   // Handling of current page's validation errors

   /**
    * Return the list of validation error messages displayed on the current
    * page.
    *
    * @return List of the validation error messages, empty ArrayList if no
    *         errors are found
    */
   public List<String> getMessagesFromValidationPanel() {
      List<String> errors = new ArrayList<String>();

      if (UI.condition.notFound(ID_CLOSE_ERROR_BTN).await(
            SUITA.Environment.getUIOperationTimeout()) ||
            (UI.component.property.get(Property.VISIBLE, ID_CLOSE_ERROR_BTN) == null ||
            !UI.component.property.getBoolean(Property.VISIBLE, ID_CLOSE_ERROR_BTN))) {

         return errors;
      }

      int errorCount = 1;
      if (UI.condition.isFound(ID_MORE_ERROR_BTN).await(
            SUITA.Environment.getUIOperationTimeout())) {
         UI.component.click(ID_MORE_ERROR_BTN);
         errorCount = UI.component.existingCount(ID_ERROR_LABEL);
      }

      for (int i = 0; i < errorCount; i++) {
         IDGroup ERROR_UI_LABEL = IDGroup.toIDGroup(String.format(
               ERROR_LABEL_FORMAT, i));
         String error = UI.component.property
               .get(Property.TEXT, ERROR_UI_LABEL);
         errors.add(error);
      }

      return errors;

   }

   /**
    * Extracts global wizard error message.
    * This error occurs in rare cases when the object on which wizard
    * was invoked does not longer exist.
    *
    * @return  the extracted error message.
    */
   public String getGlobalWizardErrorMessage() {
      if (UI.condition.isFound(ID_GLOBAL_WIZARD_ERROR).await(
            SUITA.Environment.getUIOperationTimeout())) {
         return UI.component.property.get(Property.TEXT, ID_GLOBAL_WIZARD_ERROR);
      }

      return null;
   }

   /**
    * Clicks the close-errors button
    */
   public void closeValidationPanel() {
      UI.component.click(ID_CLOSE_ERROR_BTN);
   }

   // ---------------------------------------------------------------------------------------
   // Handling of progress bars associated with dialog and page loading status

   /**
    * Waits until the progress bars indicating that the dialog is loading
    * disappear
    */
   public void waitForDialogToLoad() {
      waitForInitializingProgressBar();
      waitForLoadingProgressBar();
      new BaseView().waitForPageToRefresh();
   }

   /**
    * Waits for the "Loading..." progress bar of the current page to complete.
    */
   public void waitForLoadingProgressBar() {
      // Some wizard pages need longer timeout because of the backend calls made
      UI.condition.notFound(PAGE_LOADING_PROGRESS_BAR_ID).await(
            SUITA.Environment.getBackendJobMid());
   }

   /**
    * Waits for a specific time loading progress bar to disappear.
    * If it is still visible an exception is thrown - AssestionFatal.
    *
    * @param timeout    how long should be waited the loading bar to disappear
    */
   public void waitForLoadingProgressBarFatal(long timeout) {

      //wait a few seconds in order to be sure that the progress bar appears
      UI.condition.notFound(ID_CANCEL_BTN)
            .await(SUITA.Environment.getUIOperationTimeout() / 3);

      if (!UI.condition.notFound(PAGE_LOADING_PROGRESS_BAR_ID).await(timeout)) {
         throw new AssertionFail(
               String.format(
                     "Waited for %s but loading bar is still present", timeout / 1000)
            );
      }
   }

   /**
    * Waits for a specific time apply saved data progress bar to disappear.
    */
   public void waitApplySavedDataProgressBar() {
      if (!UI.condition.notFound(APPLY_SAVED_DATA_PROGRESS_BAR_ID)
            .await(SUITA.Environment.getBackendJobMid())) {
         throw new AssertionFail(
               String.format(
                     "Waited for apply saved data loading bar to dissapear but it is "
                     + "still visible after $d seconds",
                     SUITA.Environment.getBackendJobMid() / 1000
                 )
            );
      }
   }

   /**
    * Waits for the control to be enabled.
    */
   public boolean waitForControlEnabled(final IDGroup controlId) {
      Object evaluator = new Object() {
         @Override
         public boolean equals(Object other) {
            return UI.component.property.getBoolean(Property.ENABLED, controlId).equals(other);
         }
      };

      return UI.condition.isTrue(evaluator).await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Waits for the "Initializing..." progress bar of the current page to
    * complete.
    */
   protected void waitForInitializingProgressBar() {
      UI.condition.notFound(DIALOG_INITIALIZING_PROGRESS_BAR_ID).await(
            SUITA.Environment.getPageLoadTimeout());
   }
}
