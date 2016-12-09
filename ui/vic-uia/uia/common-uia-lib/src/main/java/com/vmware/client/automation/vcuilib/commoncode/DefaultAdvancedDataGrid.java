package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.AUTOMATIONNAME;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.CLASS_NAME;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_ADVANCED_GRID_MENU;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_BUTTON_DROP_DOWN_ARROW;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CONTEXT_HEADER;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_FILTER_CONTEXT_HEADER;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_GRID_HORIZONTAL_SCROLL;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_GRID_HORIZONTAL_SCROLL_LEFT_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_GRID_HORIZONTAL_SCROLL_RIGHT_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_TOOLBAR_FILTERCONTROL_TEXT_INPUT;
import static com.vmware.client.automation.vcuilib.commoncode.TestBaseUI.verifyTrueSafely;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.SELECT_COLUMNS;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.Timeout.ONE_SECOND;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.Timeout.TWO_SECONDS;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.FILTER;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.STRING_SELECT_ALL;

import java.util.List;
import java.util.Map;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.GRID_MENU_ITEM;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.IAdvancedDataGrid;
import com.vmware.flexui.componentframework.IDataGridProvider;
import com.vmware.flexui.componentframework.InteractiveObject;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.AdvancedDataGridEx;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.componentframework.controls.mx.TextInput;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.flexui.selenium.MethodCallUtil;

/**
 * Datagrid implementation.
 *
 * NOTE: this class is a copy of the one from VCUI-QE-LIB
 */
public class DefaultAdvancedDataGrid implements AdvancedDataGrid {


   private IAdvancedDataGrid advancedDataGrid = null;
   private FlashSelenium flashSelenium = null;
   private String uniqueId = null;

   private DefaultAdvancedDataGrid(String Id, Boolean rebuildInstanceCache,
         Map<Integer, String> columnsSpec, FlashSelenium flashSelenium) {
      this.uniqueId = Id;
      this.flashSelenium = flashSelenium;
      MethodCallUtil.waitForElementVisible(flashSelenium, Id, true);
      this.advancedDataGrid =
            new AdvancedDataGridEx(uniqueId, flashSelenium, columnsSpec,
                  rebuildInstanceCache);
   }

   public static AdvancedDataGrid getInstance(String Id, FlashSelenium flashSelenium) {
      return getInstance(Id, true, flashSelenium);
   }

   public static AdvancedDataGrid getInstance(String Id, Boolean rebuildInstanceCache,
         FlashSelenium flashSelenium) {
      return new DefaultAdvancedDataGrid(Id, rebuildInstanceCache, null, flashSelenium);
   }

   public static AdvancedDataGrid getInstance(String Id, Boolean rebuildInstanceCache,
         Map<Integer, String> columnsSpec, FlashSelenium flashSelenium) {
      return new DefaultAdvancedDataGrid(Id, rebuildInstanceCache, columnsSpec,
            flashSelenium);
   }

   @Override
   public void invokeColumnHeaderMenu(String columnName, GRID_MENU_ITEM menuItem) {
      rightClickColumnHeader(columnName);
      invokeMenuItem(menuItem.getIndex());
   }

   @Override
   public ContextHeader openShowHideCoumnHeader(String columnName) {

      invokeColumnHeaderMenu(columnName, GRID_MENU_ITEM.SHOW_HIDE_COLUMNS);
      return DefaultContextHeader.getInstance(ID_CONTEXT_HEADER, flashSelenium);
   }

   @Override
   public ContextHeader openFilterContextHeader() {

      Button filterDropDown = new Button(ID_BUTTON_DROP_DOWN_ARROW, flashSelenium);
      filterDropDown.click(ONE_SECOND.getDuration());
      invokeMenuItem(SELECT_COLUMNS);
      return DefaultContextHeader.getInstance(ID_FILTER_CONTEXT_HEADER, flashSelenium);
   }

   @Override
   public int filter(String value, String[] excludeColumns) {

      if (excludeColumns != null) {
         ContextHeader contextHeader = openFilterContextHeader();
         if (contextHeader.getLinkName().equals(STRING_SELECT_ALL)) {
            // Select all
            contextHeader.clickLink();
         }
         if (excludeColumns.length > 0) {
            contextHeader.selectDeselectColums(null, excludeColumns);
         }
         // close context header
         contextHeader.clickOK();
      }

      TextInput filter =
            new TextInput(ID_TOOLBAR_FILTERCONTROL_TEXT_INPUT, flashSelenium);
      filter.type(value, ONE_SECOND.getDuration());
      filter.setFocus(TWO_SECONDS.getDuration());
      advancedDataGrid.rebuildInternalCache();
      return this.advancedDataGrid.getRowsCount();
   }

   /**
    * Sets Filter component visibility.
    * If showFilter is true, it shows the Filter component.
    * If showFilter is false, it hides the Filter component.
    *
    * @param filterId
    * @param dataGridToolbarId
    * @param showFilter
    */
   @Override
   public void setFilterVisibility(String filterId, String dataGridToolbarId,
         boolean showFilter) {
      boolean filterVisibility =
            MethodCallUtil.getVisibleOnPath(flashSelenium, filterId);


      // If actual Filter visibility is the same as the expected
      // we don't need to change the component state
      if (filterVisibility == showFilter) {
         verifyTrueSafely(true, "Filter's visibility is in already expected state");
      } else {
         // Right click on the datagrid toolbar and select menu option Filter
         UIComponent dataGridToolbar = new UIComponent(dataGridToolbarId, flashSelenium);
         dataGridToolbar.rightMouseClick();
         invokeActionFromContextMenuOnGrid(FILTER);
      }
   }

   /**
    * Invokes an action from advanced menu. The function assume that the menu is already shown.
    * Action is specified by its name. Actions could be 'Filter' and 'Hide Toolbar'.
    *
    * @param actionName
    */
   public void invokeActionFromContextMenuOnGrid(String actionName) {
      InteractiveObject menuOption =
            new InteractiveObject(CLASS_NAME + ID_ADVANCED_GRID_MENU + "/"
                  + AUTOMATIONNAME + actionName, flashSelenium);
      menuOption.leftMouseClick();
   }

   private void invokeMenuItem(String menuItem) {

      try {
         ActionFunction.invokeActionFromContextMenuForDataGrid(
               menuItem,
               BrowserUtil.flashSelenium,
               ONE_SECOND.getDuration());
      } catch (Exception e) {
         e.printStackTrace();
         TestBaseUI.verifyTrueSafely(false, "Can't invoke menu item: " + menuItem);
      }
   }

   @Override
   public boolean isColumnVisible(String columnName) {
      return advancedDataGrid.isColumnVisible(columnName);
   }

   @Override
   public boolean isColumnVisibleOnScreen(String columnName) {
      return advancedDataGrid.isColumnVisibleOnScreen(columnName);
   }

   @Override
   public boolean isFirstColumnLocked() {
      return advancedDataGrid.isFirstColumnLocked();
   }

   @Override
   public boolean setFirstColumnLockState(boolean lock) {
      return advancedDataGrid.setFirstColumnLockState(lock);
   }

   @Override
   public void rightClickColumnHeader(String columnName) {
      advancedDataGrid.rightClickColumnHeader(columnName);
   }

   @Override
   public void selectItemByKey(String itemKey) {
      advancedDataGrid.selectItemByKey(itemKey);
   }


   @Override
   public void setCellValue(Integer rowIndex, String columnName, String value) {
      advancedDataGrid.setCellValue(rowIndex, columnName, value);
   }

   @Override
   public boolean sort(String columnName, String... columnNames) {
      return advancedDataGrid.sort(columnName, columnNames);
   }

   @Override
   public void rightClickCell(Integer rowIndex, String columnName) {
      advancedDataGrid.rightClickCell(rowIndex, columnName);
   }

   /**
    * Method to get the AdvancedDataGrid column names
    *
    * @return String[] column name array
    */
   @Override
   public String[] getColumnNames() {
      return advancedDataGrid.getColumnNames();
   }

   @Override
   public IDataGridProvider getDataGridProvider() {
      return advancedDataGrid.getDataGridProvider();
   }

   @Override
   public void clickCell(Integer rowIndex, String columnName) {
      advancedDataGrid.clickCell(rowIndex, columnName);

   }

   @Override
   public void doubleClickCell(Integer rowIndex, String columnName) {
      advancedDataGrid.doubleClickCell(rowIndex, columnName);

   }

   @Override
   public String getCellValue(Integer rowIndex, String columnName) {
      return advancedDataGrid.getCellValue(rowIndex, columnName);
   }

   @Override
   public String[] getColumnContents(String columnName) {
      return advancedDataGrid.getColumnContents(columnName);
   }

   @Override
   public Map<String, String> getRowContent(Integer rowIndex) {
      return advancedDataGrid.getRowContent(rowIndex);
   }

   @Override
   public int getRowsCount() {
      return advancedDataGrid.getRowsCount();
   }

   @Override
   public Boolean isColumnSortedAscending(String columnName) {
      return advancedDataGrid.isColumnSortedAscending(columnName);
   }

   @Override
   public Integer getColumnSortingIndex(String columnName) {
      return advancedDataGrid.getColumnSortingIndex(columnName);
   }

   @Override
   public void refresh() {
      advancedDataGrid.rebuildInternalCache();
   }

   @Override
   public void refresh(String keyColumnName) {
      advancedDataGrid.buildInternalCache(null, null, keyColumnName);
   }

   @Override
   public void scrollOneColumnLeft() {
      if (isHorizontalScrollVisible()) {
         UIComponent leftArrow =
               new UIComponent(uniqueId + "/" + ID_GRID_HORIZONTAL_SCROLL_LEFT_BUTTON,
                     flashSelenium);
         leftArrow.leftMouseClick();
      }
   }

   @Override
   public void scrollOneColumnRight() {
      if (isHorizontalScrollVisible()) {
         UIComponent rightArrow =
               new UIComponent(uniqueId + "/" + ID_GRID_HORIZONTAL_SCROLL_RIGHT_BUTTON,
                     flashSelenium);
         rightArrow.leftMouseClick();
      }
   }

   @Override
   public boolean isHorizontalScrollVisible() {
      UIComponent scroll =
            new UIComponent(uniqueId + "/" + ID_GRID_HORIZONTAL_SCROLL, flashSelenium);
      return scroll.getVisible();
   }

   @Override
   public int getColumnWidth(String columnName) {
      return advancedDataGrid.getColumnWidth(columnName);
   }

   @Override
   public void setColumnWidth(String columnName, int size) {
      advancedDataGrid.setColumnWidth(columnName, size);
   }

   @Override
   public <T extends DisplayObject> T getCellDisplayObject(Integer rowIndex,
         String columnName, Class<T> componentClass, String propertyName,
         String propertyValue) {
      return advancedDataGrid.getCellDisplayObject(
            rowIndex,
            columnName,
            componentClass,
            propertyName,
            propertyValue);
   }


   @Override
   public void moveColumn(String columnName, int pos) {
      advancedDataGrid.moveColumn(columnName, pos);
   }

   @Override
   public Integer findItemByName(String itemKey, String itemParentKey) {
      return advancedDataGrid.findItemByName(itemKey, itemParentKey);
   }

   @Override
   public Integer findItemByName(String itemKey) {
      return advancedDataGrid.findItemByName(itemKey);
   }

   @Override
   public void selectRows(Integer rowIndex, Integer... rows) {
      advancedDataGrid.selectRows(rowIndex, rows);
   }

   /**
    * Get the indexes of selected rows
    *
    * @return List<Integer> - indexes of selected rows
    */
   @Override
   public List<Integer> getSelectedRows() {
      return advancedDataGrid.getSelectedRows();
   }

   /**
    * Expands all the nodes of the grid tree control
    */
   @Override
   public void expandAll() {
      advancedDataGrid.expandAll();
      refresh();
   }
}
