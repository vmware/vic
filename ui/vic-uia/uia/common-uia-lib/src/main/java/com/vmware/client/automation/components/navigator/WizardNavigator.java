/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator;

import java.util.List;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.AssertionFail;

/**
 * Utility class providing stateless navigation and reading of the current state of a wizard.
 */
public class WizardNavigator extends BaseMultiPageDialogNavigator {

   private static final IDGroup ID_GLOBAL_WIZARD_ERROR =
         IDGroup.toIDGroup("errorText");

   private static final IDGroup ID_PREV_BTN =
         IDGroup.toIDGroup("tiwoDialog$back");

   private static final IDGroup ID_NEXT_BTN =
         IDGroup.toIDGroup("tiwoDialog$next");

   private static final IDGroup ID_FINISH_BTN =
         IDGroup.toIDGroup("tiwoDialog$finish");

   private static final IDGroup ID_CANCEL_BTN =
         IDGroup.toIDGroup("tiwoDialog$cancel");

   /**
    * @return True if an open wizard page is found on the screen.
    */
   public boolean isOpen() {
      // We are checking for the Next button. This makes indirect conclusion
      // that this is a wizard page, but currently we don't know a better way.
      return UI.condition.isFound(ID_NEXT_BTN).
            await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Clicks the next button.
    */
   public void clickNext() {
      UI.component.click(ID_NEXT_BTN);
   }

   /**
    * Clicks the back button.
    */
   public void clickBack() {
      UI.component.click(ID_PREV_BTN);
   }

   /**
    * Click on the Next button of the current wizard page.
    *
    * @param validationErrors
    *    A string list in which validation error messages will be
    *    returned if any are displayed.
    *
    * @return
    *    True if the move to the next page was successful e.g.
    *    no validation errors were displayed.
    */
   public boolean gotoNextPage(List<String> validationErrors) {
      waitForLoadingProgressBar();
      _logger.info("Check availability of next button.");
      if (!waitForControlEnabled(ID_NEXT_BTN)) {
         _logger.error(
               "The wizard's Next button is not enanbled. Unsuccessful attempt to click on it.");

         return false;
      }

      _logger.info("Will try to move to next wizard step.");

      final Integer originalPage = getCurrentlyActivePage();
      clickNext();

      // wait around 10 secs to check if the next page has been reached
      boolean movedToNextPage = waitForNextPage(originalPage);

      _logger.info("Have been waiting to reach desired step: " + movedToNextPage);

      // Some wizards make heavy backend verifications - so wait
      // if loading bar is still present then the execution will be terminated
      waitForLoadingProgressBarFatal(SUITA.Environment.getBackendJobMid());

      _logger.info("Have been wating the progress bar to disappear.");

      if (!movedToNextPage) {

         _logger.info("The desired step has not been reached will wait some more time.");

         // check once again quickly if we have reached the next page
         movedToNextPage = waitForNextPage(originalPage);

         _logger.info("Finished waiting once again to reach the desired wizard step.");

         if (movedToNextPage) {

            _logger.info("Second try to move to the next page was successful.");

            return waitForFurtherNavigationAvailable();
         } else {

            _logger.info("Could not move to the next wizard page - collecting errors.");

            // next page has not been reached - collect errors
            List<String> errors = getMessagesFromValidationPanel();
            if (!errors.isEmpty()) {
               if (validationErrors != null) {
                  validationErrors.clear();
                  validationErrors.addAll(errors);
               }
            }

            return false;
         }
      }

      _logger.info("Moved to next wizard step after first try.");

      return waitForFurtherNavigationAvailable();
   }

   /**
    * Click on the Next button of the current wizard page.
    *
    * @return
    *    True if the move to the next page was successful e.g.
    *    no validation errors were displayed.
    */
   public boolean gotoNextPage() {
      return gotoNextPage(null);
   }

   /**
    * Click on the Back button of the current wizard page.
    *
    * @param validationErrors
    *    A string list in which validation error messages will be
    *    returned if any are displayed.
    *
    * @return
    *    True if the move to the previous page was successful e.g.
    *    no validation errors were displayed.
    */
   public boolean gotoPrevPage(List<String> validationErrors) {
      clickBack();

      List<String> errors = getMessagesFromValidationPanel();

      if (!errors.isEmpty()) {
         if (validationErrors != null) {
            validationErrors.clear();
            validationErrors.addAll(errors);
         }
         return false;
      }

      return true;
   }

   /**
    * Click on the Back button of the current wizard page.
    *
    * @return
    *    True if the move to the previous page was successful e.g.
    *    no validation errors were displayed.
    */
   public boolean gotoPrevPage() {
      return gotoPrevPage(null);
   }

   /**
    * Click on the Finish button of the current wizard page.
    *
    * @param validationErrors
    *    A string list in which validation error messages will be
    *    returned if any are displayed.
    *
    * @return
    *    True if wizard mutation operation was submitted successfully.
    */
   public boolean finishWizard(List<String> validationErrors) {
      if (!waitForControlEnabled(ID_FINISH_BTN)) {
         _logger.error(
               "The wizard's Finish button is not enanbled. Unsuccessful attempt to click on it.");

         return false;
      }

      UI.component.click(ID_FINISH_BTN);

      if (UI.condition.notFound(ID_FINISH_BTN)
            .await(SUITA.Environment.getUIOperationTimeout())) {
      // The Finish button has disappeared as indication that the wizard
      // is closed and the mutation operation is launched.
         return true;
      }

      // The Finish button is still visible, so the wizard is not close.
      // Checking for errors.
      List<String> errors = getMessagesFromValidationPanel();
      if (!errors.isEmpty()) {
         if (validationErrors != null) {
            validationErrors.clear();
            validationErrors.addAll(errors);
         }
      }

      return false;
   }

   /**
    * Click on the Finish button of the current wizard page.
    *
    * @return
    *    True if wizard mutation operation was submitted successfully.
    */
   public boolean finishWizard() {
      return finishWizard(null);
   }

   /**
    * @return true if the Next button is enabled.
    */
   public boolean isNextBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_NEXT_BTN);
   }

   /**
    * @return true if the Prev button is enabled.
    */
   public boolean isPrevBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_PREV_BTN);
   }

   /**
    * @return true if the Finish button is enabled.
    */
   public boolean isFinishBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_FINISH_BTN);
   }

   /**
    * Clicks the Cancel button
    */
   // Cannot use the superclass cancel() and isCancelBtnEnabled() methods
   // because the Cancel button in the wizard dialogs has a different ID
   @Override
   public void cancel() {
      UI.component.click(ID_CANCEL_BTN);
   }

   /**
    * Checks whether the Cancel button is enabled.
    *
    * @return true if, the Cancel button in enabled, false otherwise
    */
   @Override
   public boolean isCancelBtnEnabled() {
      return UI.component.property.getBoolean(Property.ENABLED, ID_CANCEL_BTN);
   }

   /**
    * Retrieves the global wizard error text.
    *
    * @return  the error message if present, empty string otherwise
    */
   public String getGlobalWizardErrorText() {
      if (UI.condition.isFound(ID_GLOBAL_WIZARD_ERROR)
            .await(SUITA.Environment.getUIOperationTimeout())) {
         return UI.component.property.get(Property.TEXT, ID_GLOBAL_WIZARD_ERROR);
      }

      return "";
   }

   //=============== Private methods ===============

    /**
     * Method which waits for next page to be loaded by comparing the current page index with the one passed as a parameter
     * @param startPageIndex - index of the original page before clicking in Next
     * @return - true if the current index is bigger than the startPageIndex
     */
    private boolean waitForNextPage(final int startPageIndex) {
        boolean movedToNextPage = UI.condition.isTrue(new Object() {
            @Override
            public boolean equals(Object other) {
                return getCurrentlyActivePage() > startPageIndex;
            }
        }).await(SUITA.Environment.getUIOperationTimeout());
        return movedToNextPage;
    }

    /**
     * Waits for further navigation to become available.
     * It check if any of the Back, Next or Finish buttons to become available.
     * Makes sure that loading progress bar disappears.
     *
     * @return                 if further navigation is available
     * @throws AssertionFail   if progress bar is still present after specific time
     */
    private boolean waitForFurtherNavigationAvailable() throws AssertionFail {

       _logger.info("About to check availability of navigation buttons.");

       if (!isNavigationAvailable()) {
          _logger.info("None of the Back, Next or Finish buttons is enabled!");
          _logger.info("As navigation buttons not available will "
                + "wait for progress bar to disappear");

          // check if loading progress bar is not present
          waitForLoadingProgressBarFatal(SUITA.Environment.getBackendJobSmall());

          _logger.info("About to check navigation button once again");

          // check once again if further navigation is available
          return isNavigationAvailable();
       } else {
          _logger.info("Navigation buttons are available - "
                + "succesfully moved to next wizard step");
          return true;
       }
    }

    /**
     * Check if any of the Back, Next or Finish buttons is enabled.
     */
    private boolean isNavigationAvailable() {
       return isNextBtnEnabled() || isPrevBtnEnabled() || isFinishBtnEnabled();
    }
}
