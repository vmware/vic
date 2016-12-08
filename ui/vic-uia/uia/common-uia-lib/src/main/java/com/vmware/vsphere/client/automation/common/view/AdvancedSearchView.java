/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;

/**
 * Implements the "advanced search" view.
 * Home -> New Search -> click the "advanced search" link-button.
 */
public class AdvancedSearchView extends BaseSearchView {

   private static final String ID_CB_TYPE = "advancedSearch/typeSelector";

   private static final String ID_SEARCH_ROW_PREFIX =
         "advancedSearch/advancedSearchRow_%s";
   private static final String ID_CB_PROPERTY = ID_SEARCH_ROW_PREFIX
         + "/propertySelector/className=TextInput";

   private static final String ID_CB_ENUM_OPERATOR = ID_SEARCH_ROW_PREFIX
         + "/operatorDropDown";
   private static final String ID_CB_ENUM_VALUE = ID_SEARCH_ROW_PREFIX
         + "/valueDropDown";

   private static final String ID_CB_STRING_OPERATOR = ID_SEARCH_ROW_PREFIX
         + "/nameOperator";
   protected static final String ID_TI_STRING_NAME = ID_SEARCH_ROW_PREFIX + "/nameInput";
   protected static final String ID_TI_STRING_VALUE = ID_SEARCH_ROW_PREFIX + "/valueInput";
   private static final String ID_CB_COMPLIANCE = ID_SEARCH_ROW_PREFIX
         + "/complianceDropDown";

   private static final String ID_LINK_ADD = "advancedSearch/addRowLink";

   private static final String ID_BTN_SEARCH = "advancedSearch/searchButton";

   /**
    * Selects an entity type to search for from the drop-down
    * @param type - the type to search for
    */
   public void selectSearchForType(String type) {
      UI.component.value.set(type, ID_CB_TYPE);
   }

   /**
    * Clicks the "Add new criteria..." link
    */
   public void clickAddNewCriteriaLink() {
      UI.component.click(ID_LINK_ADD);
   }

   /**
    * Selects a given property from the dropdown of the rowIndex specified
    * @param rowIndex - the index of the row containing the search criteria. This index is 0-based.
    * @param propertyName - the name of the property to select.
    */
   public void selectCriteriaProperty(int rowIndex, String propertyName) {
      UI.component.value.set(propertyName, String.format(ID_CB_PROPERTY, rowIndex));
      // lose focus to trigger change value event
      UI.component.setFocus(ID_CB_TYPE);
   }

   /**
    * Selects a given operator from the dropdown of the rowIndex specified when not comparing names but enums
    * @param rowIndex - the index of the row containing the search criteria. This index is 0-based.
    * @param operatorValue - the name of the operator to select.
    */
   public void selectOperator(int rowIndex, String operatorValue) {
      final String componentIdEnum = String.format(ID_CB_ENUM_OPERATOR, rowIndex);
      final String componentIdString = String.format(ID_CB_STRING_OPERATOR, rowIndex);

      Object evaluatorEnum = new Object() {
         @Override
         public boolean equals(Object other) {
            return UI.component.exists(componentIdEnum);
         };
      };

      Object evaluatorString = new Object() {
         @Override
         public boolean equals(Object other) {
            return UI.component.exists(componentIdString);
         }
      };

      if (UI.condition.isTrue(evaluatorEnum).await(
            SUITA.Environment.getPageLoadTimeout())) {
         UI.component.value.set(operatorValue, componentIdEnum);
      } else if (UI.condition.isTrue(evaluatorString).await(
            SUITA.Environment.getPageLoadTimeout())) {
         UI.component.value.set(operatorValue, componentIdString);
      } else {
         throw new RuntimeException("Component does not exist: " + componentIdEnum
               + " or " + componentIdString);
      }
   }

   /**
    * Sets a given value for the criteria specified by rowIndex
    * @param rowIndex - the index of the row containing the search criteria. This index is 0-based.
    * @param value - the value to set.
    */
   public void setValue(int rowIndex, String value) {
      if (UI.component.exists(String.format(ID_TI_STRING_NAME, rowIndex))) {
         UI.component.value.set(value, String.format(ID_TI_STRING_NAME, rowIndex));
      } else if (UI.component.exists(String.format(ID_TI_STRING_VALUE, rowIndex))) {
         UI.component.value.set(value, String.format(ID_TI_STRING_VALUE, rowIndex));
      } else if (UI.component.exists(String.format(ID_CB_ENUM_VALUE, rowIndex))) {
         UI.component.value.set(value, String.format(ID_CB_ENUM_VALUE, rowIndex));
      } else {
         throw new RuntimeException("Component does not exist: " + ID_TI_STRING_NAME
               + " or " + ID_TI_STRING_VALUE + " or " + ID_CB_ENUM_VALUE + " on row: "
               + rowIndex);
      }
   }

   /**
    * Selects a given compliance value from the dropdown of the rowIndex specified
    * @param rowIndex - the index of the row containing the search criteria. This index is 0-based.
    * @param complianceValue - the compliance value to select.
    */
   public void selectCriteriaCompliance(int rowIndex, String complianceValue) {
      if (UI.component.exists(String.format(ID_CB_COMPLIANCE, rowIndex))) {
         UI.component.value.set(
               complianceValue,
               String.format(ID_CB_COMPLIANCE, rowIndex));
      } else if (UI.component.exists(String.format(ID_CB_ENUM_VALUE, rowIndex))) {
         // This is a custom component that for some reason doesn't work with the getIndexForValue method used by SUITA
         // in UI.component.value.set(complianceValue, String.format(ID_CB_ENUM_VALUE, rowIndex));
         SparkDropDownList dropDown =
               new SparkDropDownList(String.format(ID_CB_ENUM_VALUE, rowIndex),
                     BrowserUtil.flashSelenium);
         boolean found = false;
         for (int i = 0; i < Integer.parseInt(dropDown.getNumElements()); i++) {
            dropDown.selectItem(String.valueOf(i));
            if (dropDown.getSelectedItemText().equals(complianceValue)) {
               found = true;
               break;
            }
         }
         if (!found) {
            throw new RuntimeException("Could not set compliance value: "
                  + complianceValue);
         }
      } else {
         throw new RuntimeException("Component does not exist: " + ID_CB_COMPLIANCE
               + " or " + ID_CB_ENUM_VALUE + " on row: " + rowIndex);
      }
   }

   /**
    * Clicks the "Search" button
    */
   public void clickSearchButton() {
      UI.component.click(ID_BTN_SEARCH);
   }
}
