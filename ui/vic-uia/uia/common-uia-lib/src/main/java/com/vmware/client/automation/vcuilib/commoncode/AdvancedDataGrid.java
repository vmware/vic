package com.vmware.client.automation.vcuilib.commoncode;

import java.util.List;
import java.util.Map;

import com.vmware.client.automation.vcuilib.commoncode.TestConstants.GRID_MENU_ITEM;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.IDataGridProvider;

/**
 * The interface for data grids.
 *
 * NOTE: this interface is a copy of the one from VCUI-QE-LIB
 */
public interface AdvancedDataGrid {

   /**
    * Returns true if the column is visible, otherwise false
    *
    * @param String column name
    * @return boolean
    */
   public boolean isColumnVisible(String columnName);

   /**
    * Returns true if the column is visible on the screen, otherwise false
    *
    * @param String column name
    * @return boolean
    */
   boolean isColumnVisibleOnScreen(String columnName);

   /**
    * Returns true if the first column is locked, otherwise false
    *
    * @param String column name
    * @return boolean
    */
   public boolean isFirstColumnLocked();

   /**
    * Lock/Unlock the first column
    *
    * @param boolean - if true lock the column, if false unlock it
    * @return boolean
    */
   public boolean setFirstColumnLockState(boolean lock);

   /**
    * Right click on column header
    *
    * @param String column name
    */
   public void rightClickColumnHeader(String columnName);

   /**
    * Sets Filter component visibility.
    * If showFilter is true, it shows the Filter component.
    * If showFilter is false, it hides the Filter component.
    *
    * @param filterId
    * @param dataGridToolbarId
    * @param showFilter
    */
   public void setFilterVisibility(String filterId, String dataGridToolbarId,
         boolean showFilter);

   /**
    * Set value of the cell
    *
    * @param Integer rowIndex - row index
    * @param String columnName - column name
    * @param String value - cell value
    */
   public void setCellValue(Integer rowIndex, String columnName, String value);

   /**
    * Invoke column header menu
    *
    * @param String columnName - column name
    * @param GRID_MENU_ITEM menuItem - column headed menu item
    */
   public void invokeColumnHeaderMenu(String columnName, GRID_MENU_ITEM menuItem);

   /**
    * Open Show/Hide context header
    *
    * @param String columnName - column name
    * @param ContextHeader
    */
   public ContextHeader openShowHideCoumnHeader(String columnName);

   /**
    * Open filter context header
    *
    * @param String columnName - column name
    * @param ContextHeader
    */
   public ContextHeader openFilterContextHeader();

   /**
    * Filter grid data
    *
    * @param String value - filter value
    * @param String[] excludeColumns - columns which to be excluded from the filter, can be null
    */
   public int filter(String value, String[] excludeColumns);

   /**
    * Get contents of a row
    *
    * @param Integer - row name
    * @return Map<String,String> - row content the key is column the value is
    *         its content
    */
   public Map<String, String> getRowContent(Integer rowIndex);

   /**
    * Get rows count
    *
    * @return int - row count
    */
   public int getRowsCount();

   /**
    * Select grid item (item in most grid is represented by row) by key
    *
    * @param String itemKey - item key
    */
   public void selectItemByKey(String itemKey);

   /**
    * Sort grid
    *
    * @param String - the name of the column by which to be sort the grid
    * @param String - additional columns names (multisort)
    * @return if the sort operation is executed, otherwise false
    */
   public boolean sort(String columnName, String... columnNames);

   /**
    * Click a cell based on name of the column and row index
    *
    * @param Integer rowIndex - row index
    * @param String columnName - column name
    */
   public void clickCell(Integer rowIndex, String columnName);

   /**
    * Double click a cell based on name of the column and row index
    *
    * @param Integer rowIndex - row index
    * @param String columnName - column name
    */
   public void doubleClickCell(Integer rowIndex, String columnName);

   /**
    * Right click a cell based on name of the column and row index
    *
    * @param Integer rowIndex - row index
    * @param String columnName - column name
    */
   public void rightClickCell(Integer rowIndex, String columnName);


   /**
    * Get contents of a column
    *
    * @param String column name
    * @return String[] contents of the column
    */
   public String[] getColumnContents(String columnName);

   /**
    * Get column names
    *
    * @return String[] column names
    */
   public String[] getColumnNames();

   /**
    * Get data grid provider
    */
   public IDataGridProvider getDataGridProvider();


   /**
    * Get value of the cell
    *
    * @param Integer rowIndex - row index
    * @param String columnName - column name
    * @return - cell value as string
    */
   public String getCellValue(Integer rowIndex, String columnName);

   /**
    * Return is the sorting in the column is ascending
    *
    * @param columnName
    * @return
    */
   public Boolean isColumnSortedAscending(String columnName);

   /**
    * Get sorting index for the column
    *
    * @param columnName
    * @return
    */
   public Integer getColumnSortingIndex(String columnName);

   /**
    * Refresh the data in the grid
    */
   public void refresh();

   /**
    * Refresh the data in the grid with new key column
    */
   public void refresh(String keyColumnName);

   /**
    * Scroll to left
    */
   public void scrollOneColumnLeft();

   /**
    * Scroll to right
    */
   public void scrollOneColumnRight();

   /**
    * Returns true if the horizontal scroll is vidible, otherwise false
    *
    * @return
    */
   boolean isHorizontalScrollVisible();

   /**
    * Get column width
    *
    * @param String column name
    * @return int
    */
   public int getColumnWidth(String columnName);

   /**
    * Set column width
    *
    * @param String column name
    * @param column width
    */
   public void setColumnWidth(String columnName, int size);

   /**
    * Get the UI component inside a data grid cell.
    *
    * @param rowIndex - row index
    * @param columnName - column name
    * @param componentClass - the class of the UI component to be instantiated
    * @param propertyName - a unique property to filter the cell's children by
    * @param propertyValue - the value of the unique property
    * @return the UI component of the given type.
    */
   public <T extends DisplayObject> T getCellDisplayObject(Integer rowIndex,
         String columnName, Class<T> componentClass, String propertyName,
         String propertyValue);

   /**
    * Move column with specified name to specified column position
    *
    * @param String column name
    * @param int - new column position in the grid
    */
   public void moveColumn(String columnName, int pos);

   /**
    * Returns the location (rowIndex) of the specified item in hierarchical data
    * grid
    *
    * @param itemKey
    * @param itemParentKey
    * @return the row index if found, otherwise <code>null</code>
    */
   public Integer findItemByName(String itemKey, String itemParentKey);

   /**
    * Returns the location (rowIndex) of the specified item in non-hierarchical
    * data grid
    *
    * @param itemKey
    * @return the row index if found, otherwise <code>null</code>
    */
   public Integer findItemByName(String itemKey);

   /**
    * Select rows in the grid
    *
    * @param Integer rowIndex - row index
    * @param Integer ... rows - additional rows which to be selected
    */
   public void selectRows(Integer rowIndex, Integer... rows);

   /**
    * Get the indexes of selected rows
    *
    * @return List<Integer> - indexes of selected rows
    */
   public List<Integer> getSelectedRows();

   /**
    * Expands all the nodes of the grid tree control
    */
   public void expandAll();
}
