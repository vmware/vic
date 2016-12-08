/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator;

import java.util.List;

import org.apache.commons.lang.NotImplementedException;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * Provides methods for navigation through a multi-page NGC dialog. Note that
 * the class should be sub-classed since it doesn't represent a specific dialog
 * navigator.
 */
public class BaseMultiPageDialogNavigator extends BaseDialogNavigator {

   private static final String SELECTED_TOC_STATE = "selected";
   private static final String SELECTED_AND_COMPLETED_TOC_STATE = "selectedAndComplete";
   private static final IDGroup ID_PAGE_TITLE =
         IDGroup.toIDGroup("pageHeaderTitleElement");

   private static final IDGroup ID_PAGE_HEADER_DESCRIPTION =
         IDGroup.toIDGroup("pageHeaderDescriptionElement");

   private static final IDGroup ID_TOC_NAVIGATOR =
         IDGroup.toIDGroup("wizardPageNavigator");
   private static final IDGroup ID_TOC_DATAGROUP =
         IDGroup.toIDGroup("wizardPageNavigator/dataGroup");
   private static final String ID_STEP_STATUS_FORMAT =
         "className=WizardTableOfContentsSkinInnerClass0[%s]";
   private static final String ID_STEP_LABEL_FORMAT =
         ID_STEP_STATUS_FORMAT + "/stepGroup/className=Label[1]";

   //---------------------------------------------------------------------------
   // Methods that provide support for  navigating through the TOC items

   /**
    * Gets number of pages in the dialog.
    *
    * @return  number of pages
    */
   public Integer getNumberOfPages() {
      String childrenNumber =
            UI.component.property.get(Property.CHILDREN_NUMBER, ID_TOC_DATAGROUP);
      return Integer.valueOf(childrenNumber);
   }

   /**
    * Checks if the page is selected.
    *
    * @param pageId  the ID of the page
    * @return        true if the page is selected, false otherwise
    */
   public boolean isPageSelected(String pageId) {
      Integer elements = getNumberOfPages();

      for (int i = 0; i < elements; i++) {
         String currentPageId = UI.component.property.get(
               Property.ID, String.format(ID_STEP_LABEL_FORMAT, i)
            );
         if (currentPageId.equals(pageId)) {
            String currentState = UI.component.property.get(
                  Property.CURRENT_STATE, String.format(ID_STEP_STATUS_FORMAT, i)
               );
            if (currentState.equals(SELECTED_AND_COMPLETED_TOC_STATE) ||
                  currentState.equals(SELECTED_TOC_STATE)) {
               return true;
            }

            break;
         }
      }

      return false;
   }

   /**
    * Gets the index of the currently active page.
    * Indexes start from 0.
    *
    * @return     the retrieved index
    */
   public Integer getCurrentlyActivePage() {
      return Integer.valueOf(
            UI.component.property.get(Property.VALUE_INDEX, ID_TOC_NAVIGATOR)
         );
   }

   /**
    * Gets the pages available in the left pane of the dialog.
    *
    * @return List of the available page titles
    */
   public List<String> getPageTitles() {
      // TODO inikolova: Implement
      throw new NotImplementedException();
   }

   /**
    * Navigates to the specified page.
    *
    * @param pageId ID of the page
    * @param validationErrors A string list in which validation error messages
    *    will be returned if any are displayed.
    *
    * @return True if the navigation to the specified page was successful,
    *    false otherwise
    */
   public boolean goToPage(
         final String pageId, List<String> validationErrors) {
      _logger.info("Clicking on page ID: " + pageId);
      UI.component.click(IDGroup.toIDGroup(pageId));

      // wait around 10 secs to check if the next page has been reached
      boolean pageSelected = UI.condition.isTrue(
            new Object() {
               @Override
               public boolean equals(Object other) {
                  return isPageSelected(pageId);
               }
            }
         ).await(SUITA.Environment.getUIOperationTimeout());

      // Some wizards make heavy backend verifications - so wait
      // if loading bar is still present then the execution will be terminated
      waitForLoadingProgressBarFatal(SUITA.Environment.getBackendJobSmall());

      if (!pageSelected) {
         // check once again very quickly if we have reached the desired page
         pageSelected = UI.condition.isTrue(
               new Object() {
                  @Override
                  public boolean equals(Object other) {
                     return isPageSelected(pageId);
                  }
               }
            ).await(SUITA.Environment.getUIOperationTimeout() / 3);

         if (pageSelected) {
            return true;
         } else {
            // the desired page has not been reached - collect errors
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

      return true;
   }

   /**
    * Navigates to the specified page.
    *
    * @param pageId ID of the page
    *
    * @return True if the navigation to the specified page was successful,
    *    false otherwise
    */
   public boolean goToPage(String pageId) {
      return goToPage(pageId, null);
   }

   /**
    * Gets the title of a page in the left navigation pane..
    *
    * @param pageId ID of a page
    *
    * @return The name of the page as shown in the left pane of the dialog.
    */
   public String getPageTitle(String pageId) {
      return UI.component.property.get(Property.TEXT, pageId);
   }

   /**
    * Gets the header title of the current page.
    *
    * @return The name of the current page as shown in the page header.
    */
   public String getPageTitle() {
      return UI.component.property.get(Property.TEXT, ID_PAGE_TITLE);
   }

   /**
    * Gets the header description of the current page.
    *
    * @return     the description of the current page
    */
   public String getPageHeaderDescription() {
      return UI.component.property.get(Property.TEXT, ID_PAGE_HEADER_DESCRIPTION);
   }

   /**
    * Set the focus to the page header. Use it to trigger the change action
    * for the numeric stepper controls. Their change is triggered once the
    * focus is moved away from them.
    */
   public void setFocusToPageHeader() {
      UI.component.setFocus(ID_PAGE_TITLE);
   }
}
