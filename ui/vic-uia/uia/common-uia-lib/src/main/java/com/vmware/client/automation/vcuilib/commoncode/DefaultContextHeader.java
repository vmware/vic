package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CONTEXT_HEADER_ALL_COLUMN_NAMES_PROPERTY;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CONTEXT_HEADER_LIST;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_DATAGRID_CONTEXT_HEADER_SCROLLER;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_DATAGRID_CONTEXT_HEADER_SCROLLER_DOWN_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_DATAGRID_CONTEXT_HEADER_SCROLL_POSITION;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_OK;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_SELECT_ALL_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_SHOW_HIDE_COLUMNS_CLOSE_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.Timeout;
import com.vmware.flexui.componentframework.InteractiveObject;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.componentframework.controls.mx.CheckBox;

/**
 * Context header implementation.
 *
 * NOTE: this class is a copy of the one from VCUI-QE-LIB
 */
public class DefaultContextHeader implements ContextHeader {

   private static int PAGE_SIZE = 12;
   private final Button okButton;
   private final Button closeButton;
   private final Button selectAllLink;
   private String contexHeadedID = null;
   private FlashSelenium flashSelenium = null;

   private DefaultContextHeader(String id, FlashSelenium flashSelenium) {
      this.flashSelenium = flashSelenium;
      okButton = new Button(ID_OK, flashSelenium);
      closeButton = new Button(ID_SHOW_HIDE_COLUMNS_CLOSE_BUTTON, flashSelenium);
      selectAllLink = new Button(ID_SELECT_ALL_BUTTON, flashSelenium);
      contexHeadedID = id;

   }

   /**
    * Create ContextHeader instance
    *
    * @param String id - ContextHeader id
    * @param String flashSelenium
    *
    * @return ContextHeader
    */
   public static ContextHeader getInstance(String id, FlashSelenium flashSelenium) {
      return new DefaultContextHeader(id, flashSelenium);
   }

   /**
    * Click OK button of context header
    *
    */
   @Override
   public void clickOK() {

      okButton.waitForElementEnable(DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE);
      okButton.click();
   }

   /**
    * Click close button of context header
    *
    */
   @Override
   public void clickClose() {

      closeButton.waitForElementEnable(DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE);
      closeButton.click();
   }

   /**
    * Click link in context header
    *
    */
   @Override
   public void clickLink() {

      selectAllLink.click();
   }

   /**
    * Get name of the link in context header
    *
    */
   @Override
   public String getLinkName() {
      return selectAllLink.getAutomationName();
   }

   /**
    * Return true if OK button is enables, otherwise false
    *
    * @return true
    */
   @Override
   public boolean isOKButtonEnabled() {
      return okButton.getEnabled();
   }

   /**
    * Return column check boxes state in context header
    *
    * @param String[] columns - columns which check box state to be return
    * @return Map<String, Boolean> - key is column name, value is true if the check box is checked, false if the check
    *         box is not checked and <code>null</code> if this column deosn't exist
    */
   @Override
   public Map<String, Boolean> getColumnsSelectedState(String[] columns) {
      CheckBox check = null;
      Map<String, Boolean> result = new HashMap<String, Boolean>();
      String[] allColumns = getAllColumns();
      for (String columnName : columns) {
         check = getColumnCheckBox(columnName, allColumns);
         result.put(
               columnName,
               check != null ? Boolean.valueOf(check.getProperty("selected", false))
                     : null);
      }

      return result;
   }


   /**
    * Check and uncheck column check boxes in context header
    *
    * @param String[] selectColumns - columns which check boxes to be checked
    * @param String[] deselectColumns - columns which check boxes to be unchecked
    * @return Map<String, Boolean> - key is column
    */
   @Override
   public void selectDeselectColums(String[] selectColumns, String[] deselectColumns) {
      selectColumns(selectColumns, true);
      selectColumns(deselectColumns, false);
   }

   /**
    * Set provided columns with desired selection state
    *
    * @param columns
    * @param doSelect
    */
   private void selectColumns(String[] columns, boolean doSelect) {
      Map<String, Boolean> states = null;
      if (columns != null) {
         String[] allColumns = getAllColumns();
         states = getColumnsSelectedState(columns);
         for (String columnName : columns) {
            if (states.get(columnName) != doSelect) {
               clickOnColumn(columnName, allColumns);
            }
         }
      }
   }

   /**
    * Click on checkbox for a gived column in the header context menu
    *
    * @param columnName
    * @param allColumns
    */
   private void clickOnColumn(String columnName, String[] allColumns) {
      CheckBox check = getColumnCheckBox(columnName, allColumns);
      check.leftMouseClick();
      Timeout.FIVE_HUNDRED_MILLIS.consume();
   }

   @Override
   public String[] getAllColumns() {
      final String COLUMNS_DELIMITER = ";;";
      InteractiveObject list =
            new InteractiveObject(contexHeadedID + ID_CONTEXT_HEADER_LIST, flashSelenium);
      String rawColumnNames =
            list.getProperties(
                  ID_CONTEXT_HEADER_ALL_COLUMN_NAMES_PROPERTY,
                  COLUMNS_DELIMITER);
      return rawColumnNames.split(COLUMNS_DELIMITER);
   }

   private String escapeColumnName(String columnName) {
      return columnName.replaceAll("/", "\\\\\\\\\\\\/");
   }

   /**
    * Scroll to the check box and get it for specific column
    *
    * @param columnName
    * @param allColumns
    * @return
    */
   private CheckBox getColumnCheckBox(String columnName, String[] allColumns) {
      List<String> columns = Arrays.asList(allColumns);
      int columnIndex = columns.indexOf(columnName);
      if (columnIndex == -1) {
         return null;
      }

      if (allColumns.length > PAGE_SIZE) {
         int checkPage = columnIndex / PAGE_SIZE;
         UIComponent scroller =
               new UIComponent(ID_DATAGRID_CONTEXT_HEADER_SCROLLER, flashSelenium);
         scroller.setProperty(
               ID_DATAGRID_CONTEXT_HEADER_SCROLL_POSITION,
               Integer.toString(checkPage * PAGE_SIZE - 1));
         UIComponent downButton =
               new UIComponent(ID_DATAGRID_CONTEXT_HEADER_SCROLLER_DOWN_BUTTON,
                     flashSelenium);
         downButton.leftMouseClick();
         Timeout.FIVE_HUNDRED_MILLIS.consume();
      }

      return new CheckBox(contexHeadedID + "/label=" + escapeColumnName(columnName),
            flashSelenium);
   }
}
