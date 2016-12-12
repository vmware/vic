/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang.ArrayUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.DefaultAdvancedDataGrid;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.custom.ItemRenderer;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.apl.sele.SeleAPLImpl;

/**
 * The only task for now is to be the entry point for accessing the VCUIQE_LIB
 * AdvancedDataGrid object.
 */
public class GridControl {

   private static final Logger _logger = LoggerFactory.getLogger(GridControl.class);
   private static final IDGroup REFRESH_BTN = IDGroup.toIDGroup("refreshButton");

   /**
    * Constants holding IDs of 'Actions' toolbar button for various grids.
    */
   public enum GridActionsButton {
      DATACENTER_TOP_LEVEL("TODO: Write correct id"), DATACENTER_CLUSTERS("TODO: Write correct id"), DATACENTER_HOSTS(
            "TODO: Write correct id"), DATACENTER_VMS("vmsForDatacenter/allActions/button"), DATACENTER_VMT_IN_FOLDERS(
                  "TODO: Write correct id"), DATACENTER_VAPPS("TODO: Write correct id"), DATACENTER_DATASTORES(
                        "TODO: Write correct id"),

                        CLUSTER_VMS("vmsForCluster/allActions/button"),

                        VCENTER_VMS("vmsForVCenter/allActions/button"),

                        CLUSTERED_HOST_VMS("vmsForClusteredHost/allActions/button"),

                        DATASTORE_VMS("vmsForDatastore/allActions/button"),

                        CL_TEMPLATES("ovfTemplateForLibrary/allActions/button"),

                        CL_OTHERTYPES("otherTypeForLibrary/allActions/button"),

                        ;

      private final String allActionsButtonId;

      GridActionsButton(String allActionsButtonId) {
         this.allActionsButtonId = allActionsButtonId;
      }

      public String getAllActionsButtonId() {
         return allActionsButtonId;
      }
   }

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private static final int IDX_COL_NAME = 1;

   private static final String GRID_CHECK_VAL = "true";

   /**
    * Finds the <code>AdvancedDataGrid</code> specified by the gridId.
    *
    * @param gridId
    *            Id of the grid
    *
    * @return <code>AdvancedDataGrid</code> instance
    */
   public static AdvancedDataGrid findGrid(IDGroup gridId) {
      return getGrid(gridId, 0, null);
   }

   /**
    * Finds the <code>AdvancedDataGrid</code> specified by the gridId.
    *
    * @param gridId
    *            Id of the grid
    *
    * @return <code>AdvancedDataGrid</code> instance
    */
   public static AdvancedDataGrid findGrid(String gridId) {
      return getGrid(IDGroup.toIDGroup(gridId), 0, null);
   }

   /**
    * Finds the <code>AdvancedDataGrid</code> specified by the gridId. Use this
    * method if the grid contains specific column elements, such as radio
    * buttons.
    *
    * @param gridId
    *            Id of the grid
    * @param numberOfColumns
    *            Total count of the columns
    *
    * @return <code>AdvancedDataGrid</code> instance
    */
   public static AdvancedDataGrid findGrid(IDGroup gridId, int numberOfColumns) {
      return getGrid(gridId, numberOfColumns, null);
   }

   /**
    * Finds the <code>AdvancedDataGrid</code> specified by the gridId. Use this
    * method if the grid contains specific Item Renderers such as TextInputItemRenderer
    *
    * @param gridId
    *            Id of the grid
    * @param numberOfColumns
    *            Total count of the columns
    *
    * @return <code>AdvancedDataGrid</code> instance
    */
   public static AdvancedDataGrid findGrid(IDGroup gridId, Map<Integer, String> colSpec) {
      return getGrid(gridId, 0, colSpec);
   }

   /**
    * Gets the String value in a given column of the selected element in the
    * grid for a datagrid with radiobutton selection.
    *
    * @param dataGrid
    *            - the grid
    * @param columnName
    *            - the name of the column in which to look for the selected
    *            entity
    * @return - returns a the column value of the selected entity in the grid
    */
   public static String getSelectedEntityColumnValue(IDGroup gridId, String columnName) {
      String selectedRow = new UIComponent(gridId.toString(), BrowserUtil.flashSelenium).getProperty("selectedIndex");
      return findGrid(gridId).getCellValue(Integer.getInteger(selectedRow), columnName);
   }

   /**
    * Gets the String value of the specified cell.
    *
    * @param grid - the grid in which to search
    * @param columnName - the column in which the cell is located
    * @param rowIndex - the row on which the cell is located
    * @return - the contents of the cell
    */
   public static String getCellContents(AdvancedDataGrid grid, String columnName, int rowIndex) {
      return grid.getCellValue(rowIndex, columnName);
   }

   /**
    * Gets the String value in a given column of a given element in the grid.
    * Retrieves the value of a cell in the grid.
    *
    * @param dataGrid
    *            - the grid
    * @param index
    *            - row index
    * @param valueColumn
    *            - the column from which to get the value
    * @return - returns the column value of the specified entity in the grid
    */
   public static String getEntityColumnValue(AdvancedDataGrid grid, int index, String valueColumn) {
      return grid.getCellValue(index, valueColumn);
   }

   /**
    * Gets the index of the entity in the Grid
    *
    * @param dataGrid
    *            - the grid
    * @param columnName
    *            - the column name we are going to search in
    * @param entityName
    *            - the entity name we are going to search for
    * @return - the index of the item or -1 if not found
    */
   public static int getEntityIndex(AdvancedDataGrid dataGrid, String columnName, String entityName) {
      if (dataGrid == null || Strings.isNullOrEmpty(columnName)) {
         throw new IllegalArgumentException("dataGrid or columnName not set");
      }

      if (Strings.isNullOrEmpty(entityName)) {
         throw new IllegalArgumentException("entityName not set");
      }

      String[] entities = dataGrid.getColumnContents(columnName);

      int entityIndex = ArrayUtils.indexOf(entities, entityName);

      return entityIndex;
   }

   /**
    * Gets a list of indexes that have the expected value in the given column
    *
    * @param dataGrid
    *            - the grid
    * @param columnName
    *            - the column name we are going to search in
    * @param entityName
    *            - the entity name we are going to search for
    * @return - the index of the item or -1 if not found
    */
   public static List<Integer> getEntityIndexes(AdvancedDataGrid dataGrid, String columnName, String entityName) {
      if (dataGrid == null || Strings.isNullOrEmpty(columnName)) {
         throw new IllegalArgumentException("dataGrid or columnName not set");
      }

      if (Strings.isNullOrEmpty(entityName)) {
         throw new IllegalArgumentException("entityName not set");
      }

      String[] entities = dataGrid.getColumnContents(columnName);
      List<Integer> result = new ArrayList<Integer>();
      for (int i = 0; i < entities.length; i++) {
         if (entities[i].equals(entityName)) {
            result.add(i);
         }
      }
      return result;
   }

   /**
    * Gets the index of the entity in the Grid. Looks the first grid column for a match.
    *
    * @param dataGrid - the grid
    * @param entityName - the entity name we are going to search for
    * @return - the index of the item or -1 if not found
    */
   public static int getEntityIndex(AdvancedDataGrid dataGrid, String entityName) {
      if (dataGrid == null) {
         throw new IllegalArgumentException("dataGrid not set");
      }

      if (Strings.isNullOrEmpty(entityName)) {
         throw new IllegalArgumentException("entityName not set");
      }

      String[] columns = dataGrid.getColumnNames();
      String[] entities = dataGrid.getColumnContents(columns[0]);

      int entityIndex = ArrayUtils.indexOf(entities, entityName);

      return entityIndex;
   }

   /**
    * Method to get the column contents of a datagrid
    *
    * @param dataGrid - datagrid that is under test
    * @param columnName - column name whose contents to get
    * @return List of Stings with the data in the column
    */
   public static List<String> getColumnContents(AdvancedDataGrid dataGrid, String columnName) {
      String[] columnContents = dataGrid.getColumnContents(columnName);
      return columnContents == null ? new ArrayList<String>() : Arrays.asList(columnContents);
   }

   /**
    * Gets the number of rows in the grid
    *
    * @param dataGrid
    *            - the grid
    * @return - the number of rows
    */
   public static int getNumberOfEntities(AdvancedDataGrid dataGrid) {
      return dataGrid.getRowsCount();
   }

   /**
    * Gets the total number of rows in the grid.
    *
    * @param dataGrid
    *            - the grid
    * @return - total number of rows in the grid
    */
   public static int getRowsCount(AdvancedDataGrid dataGrid) {
      return dataGrid.getRowsCount();
   }

   /**
    * Sets a value in a cell
    *
    * @param dataGrid - the grid
    * @param index - row index where to set the value
    * @param column - column name
    * @param value - value to set
    * @return - true if the value is set, false otherwise
    */
   public static boolean setCellValue(AdvancedDataGrid dataGrid, Integer index, String column, String value) {
      dataGrid.setCellValue(index, column, value);
      return getEntityIndex(dataGrid, column, value) > -1;
   }

   /**
    * Sets a value in a cell
    *
    * @param dataGrid - the grid
    * @param index - row index where to set the value
    * @param column - column name
    * @param value - value to set
    */
   public static void setCellValueNoResult(AdvancedDataGrid dataGrid, Integer index, String column, String value) {
      // there are cases in which getEntityIndex in above method doesn't execute correctly
      dataGrid.setCellValue(index, column, value);
   }

   /**
    * /**
    * Selects the specified entity in the grid. The method uses <code>
    * AdvancedDataGrid.getColumnContents</code> internally and is recommended to
    * be used when the grid contains specific columns such as radio buttons.
    *
    * @param dataGrid
    *            <code>AdvancedDataGrid</code> instance
    * @param columnName
    *            Name of the column where the entity will be searched for
    * @param entityName
    *            Name of the entity to be selected
    *
    * @return True if the selection of all specified entity is successful, false
    *         otherwise, e.g. invalid column name or entity name is passed
    */
   public static boolean selectEntity(AdvancedDataGrid dataGrid, String columnName, String entityName) {
      int entityIndex = getEntityIndex(dataGrid, columnName, entityName);
      if (entityIndex < 0) {
         _logger.error("Unable to find the specified entity " + entityName + "in the grid");
         return false;
      }

      dataGrid.selectRows(entityIndex);
      return true;
   }

   /**
    * Select the specified entities in the specified data grid.
    *
    * @param dataGrid
    *            <code>AdvancedDataGrid</code> instance
    * @param entityNames
    *            One or more names of entities to be selected
    *
    * @return True if the selection of all specified entities is successful,
    *         false otherwise, e.g. invalid entity name is passed
    */
   public static boolean selectEntities(AdvancedDataGrid dataGrid, String... entityNames) {

      if (dataGrid == null) {
         throw new IllegalArgumentException("dataGrid not set");
      }

      if (ArrayUtils.isEmpty(entityNames)) {
         throw new IllegalArgumentException("entityNames not set");
      }

      _logger.info("Search for the row indexes of the entities");
      Integer rowIndexes[] = new Integer[entityNames.length];

      for (int i = 0; i < entityNames.length; i++) {
         String entityName = entityNames[i];

         Integer entityIndex = dataGrid.findItemByName(entityName);
         if (entityIndex == null) {
            _logger.error(String.format("Unable to find entity: %s", entityName));
            return false;
         }
         rowIndexes[i] = entityIndex;
      }

      _logger.info("Select the entities by their row indexes");
      dataGrid.selectRows(rowIndexes[0], Arrays.copyOfRange(rowIndexes, 1, entityNames.length));

      return true;
   }

   /**
    * Selects all the items in the grid
    * @param dataGrid - <code>AdvancedDataGrid</code> instance
    */
   public static void selectAll(AdvancedDataGrid dataGrid) {
      if (dataGrid.getRowsCount() == 0) {
         return;
      }

      List<Integer> moreRows = new ArrayList<Integer>();
      for (int i = 1; i < dataGrid.getRowsCount(); i++) {
         moreRows.add(i);
      }
      dataGrid.selectRows(0, moreRows.toArray(new Integer[] {}));
   }

   /**
    * Select the specified entities in the specified data grid.
    *
    * @param id
    *            Identifier of the <code>AdvancedDataGrid</code> instance
    *
    * @param entityNames
    *            One or more names of entities to be selected
    *
    * @return True if the selection of all specified entities is successful,
    *         false otherwise, e.g. invalid entity name is passed
    */
   // TODO: This method and checkGridCheckBox() have similar functionality
   // and should be merged to a single method.
   public static boolean checkEntities(String id, String... entityNames) {

      if (Strings.isNullOrEmpty(id)) {
         throw new IllegalArgumentException("Required id argument is not set.");
      }

      if (ArrayUtils.isEmpty(entityNames)) {
         throw new IllegalArgumentException("Required entityNames argument is empty.");
      }

      waitForGridToLoad(IDGroup.toIDGroup(id));

      // Get a reference to the advanced data grid.
      com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid grid = new com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid(
            id, BrowserUtil.flashSelenium);

      // Get all entity names specified in the Name column.
      String[] names = grid.getColumnContents(String.valueOf(IDX_COL_NAME));

      boolean result = true;

      // Find the grid indices of the entity names specified for selection and
      // make the selection.
      for (int i = 0; i < entityNames.length; ++i) {
         int idx = ArrayUtils.indexOf(names, entityNames[i]);

         if (idx < 0) {
            result = false;
         }

         grid.checkUncheckAdvDatagridCheckBox(
               String.valueOf(idx),
               GRID_CHECK_VAL,
               SUITA.Environment.getUIOperationTimeout());
      }

      return result;
   }

   /**
    * Method that uses "doFlexClickDataGridCell" flex call in order to click an
    * AdvancedDatagrid cell
    *
    * @param gridId - String id of the grid, whose cell is to be clicked
    * @param rowIndex - String index of the row
    * @param columnIndex - String index of the column
    */
   public static void clickDataGridCell(String gridId, String rowIndex, String columnIndex) {
      com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid dataGrid = new com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid(
            gridId, BrowserUtil.flashSelenium);

      dataGrid.clickDataGridCell(rowIndex, columnIndex, (int) SUITA.Environment.getUIOperationTimeout());
   }

   /**
    * Method that uses the AdvancedDataGrid component to relay the click call.
    *
    * @param grid - the grid, whose cell is to be clicked
    * @param rowIndex - String index of the row
    * @param columnName - String name of the column
    */
   public static void clickCell(AdvancedDataGrid grid, int rowIndex, String columnName) {
      grid.clickCell(rowIndex, columnName);
   }

   /**
    * Finds and right-clicks on the specified entity in the data grid.
    *
    * @param dataGrid
    *            <code>AdvancedDataGrid</code> instance
    * @param columnName
    *            Name of the column where the entity will be searched for
    * @param entityName
    *            Name of the entity to be right-clicked on
    *
    * @return True if the right-click is successful, false otherwise, e.g. if
    *         the entity is not found in the data grid.
    */
   public static boolean rightClickEntity(AdvancedDataGrid dataGrid, String columnName, String entityName) {
      return rightClickEntity(dataGrid, columnName, new String[] { entityName });
   }

   /**
    * Finds and right-clicks on the specified entities in the data grid.
    *
    * @param dataGrid
    *            the datagrid where entities will be selected
    * @param columnName
    *            the name of the column where the entity will be searched for
    * @param entityNames
    *            the names of the entities to be right-clocked on
    *
    * @return true if the right-click is successful, false otherwise, e.g. if
    *         the entity is not found in the data grid.
    */
   public static boolean rightClickEntity(AdvancedDataGrid dataGrid, String columnName, String... entityNames) {

      if (dataGrid == null || Strings.isNullOrEmpty(columnName)) {
         throw new IllegalArgumentException("dataGrid or columnName not set");
      }

      if (entityNames == null || entityNames.length == 0) {
         throw new IllegalArgumentException("entityNames not set");
      }

      _logger.info(String.format("Invoke right-click on %s", (Object[]) entityNames));

      _logger.info("Selecting entities in the grid");
      LinkedList<Integer> rowIndexes = new LinkedList<Integer>();

      for (String entityName : entityNames) {
         Integer rowIndex = getEntityIndex(dataGrid, columnName, entityName);
         if (rowIndex == null) {
            _logger.error(String.format("Unable to find entity: %s in the grid", entityName));
            return false;
         } else {
            rowIndexes.add(rowIndex);
         }
      }

      Integer firstRowIndex = rowIndexes.pop();
      if (rowIndexes.contains(firstRowIndex)) {
         throw new IllegalArgumentException("The primary entity was added twice to the list: " + entityNames[0]);
      }

      dataGrid.selectRows(firstRowIndex, rowIndexes.toArray(new Integer[0]));

      _logger.info("Select and right-click on the entities");
      dataGrid.rightClickCell(firstRowIndex, columnName);

      return true;
   }

   /**
    * Finds an entity in a data grid and checks a check-box found on the same
    * row as the entity.
    *
    * @param gridId
    *            Id of the data grid
    * @param columnIndex
    *            Index of the column where the entity will be searched for.
    * @param entityName
    *            Name of the entity whose check-box will be checked.
    *
    * @return True if the check-box is checked successfully, false otherwise,
    *         e.g. invalid column index or entity name is passed.
    */
   public static boolean checkGridCheckBox(IDGroup gridId, int columnIndex, String entityName) {

      waitForGridToLoad(gridId);

      com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid grid = new com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid(
            gridId.getValue(Property.DIRECT_ID), ((SeleAPLImpl) SUITA.Factory.apl()).getFlashSelenium());

      String[] entities = grid.getColumnContents(String.valueOf(columnIndex));

      int entityIndex = ArrayUtils.indexOf(entities, entityName);

      if (entityIndex < 0) {
         _logger.error("Unable to find the entity");
         return false;
      }

      grid.checkUncheckAdvDatagridCheckBox(
            String.valueOf(entityIndex),
            "true",
            SUITA.Environment.getUIOperationTimeout());

      return true;
   }

   /**
    * Returns true if the column is visible, otherwise false
    *
    * @param String column name
    * @return boolean
    */
   public static boolean isColumnVisible(AdvancedDataGrid grid, String columnName) {
      return grid.isColumnVisible(columnName);
   }

   /**
    * Waits until the grid is loaded. The method first waits for a global refresh to finish.
    * It then waits for the grid's progress bar to disappear. It then waits for the
    * contents to appear. And at the end checks if a page refresh has started.
    */
   public static void waitForGridToLoad(IDGroup gridId) {
      _logger.info("Start waiting for the grid to load");

      //TODO: change this once all views are refactored
      UI.condition.isFound(REFRESH_BTN).await(SUITA.Environment.getBackendJobMid());
      waitForGridProgressBarToDisappear(gridId);
      waitForGridContentToRender(gridId);
      //TODO: change this once all views are refactored
      UI.condition.isFound(REFRESH_BTN).await(SUITA.Environment.getBackendJobMid());

      _logger.info("Finished waiting for the grid to load");
   }

   /**
    * Waiting for the showProgress
    *
    * @param gridId
    */
   private static void waitForGridProgressBarToDisappear(IDGroup gridId) {
      while (isGridProgressBarVisible(gridId)) {
         waitOneMillisecond();
      }
   }

   /**
    * Returns the status of the showProgress grid property
    *
    * @param gridId - the grid from which to get the showProgress status from
    * @return boolean
    */
   private static Boolean isGridProgressBarVisible(IDGroup gridId) {
      return UI.component.property.getBoolean(Property.SHOW_PROGRESS, gridId);
   }

   private static void waitOneMillisecond() {
      try {
         Thread.sleep(1);
      } catch (InterruptedException e) {
         e.printStackTrace();
      }
   }

   /**
    * In some cases it takes time for the grid contents to render after
    * the grid has loaded. This method waits for the contents of the grid
    * to appear.
    *
    * @param gridId - the grid to get contents holder enabled status from
    */
   private static void waitForGridContentToRender(IDGroup gridId) {
      while (!isGridContentHolderEnabled(gridId)) {
         waitOneMillisecond();
      }
   }

   /**
    * Returns true if the grid content holder object has been loaded.
    * This is the object, which directly contains the contents. If it is loaded
    * the we can be sure the contents are visible.
    *
    * @param gridId - the grid to get contents holder enabled status from
    * @return boolean
    */
   private static Boolean isGridContentHolderEnabled(IDGroup gridId) {
      final String contentHolderId = gridId.getValue(Property.DIRECT_ID)
            + "/className=AdvancedListBaseContentHolder[0]";
      return UI.component.property.getBoolean(Property.ENABLED, IDGroup.toIDGroup(contentHolderId));
   }

   /**
    * Instantiate and return the <code>AdvancedDataGrid</code> specified by the
    * gridId. When the grid contains specific column elements like radio buttons
    * specify their column count by providing numberOfColumns. NOTE: If the data
    * grid is not found on the screen due to loading issues the method will make
    * screenshot and will return null.
    *
    * @param gridId
    *            Id of the grid
    * @param numberOfColumns
    *            Total count of the columns
    *
    * @return <code>AdvancedDataGrid</code> instance if AdvancedDataGrid control
    *         is found otherwise return null.
    */
   private static AdvancedDataGrid getGrid(final IDGroup gridId, int numberOfColumns, Map<Integer, String> colSpec) {
      waitForGridToLoad(gridId);

      if (numberOfColumns > 0) {
         colSpec = new HashMap<Integer, String>();
         for (int i = 0; i < numberOfColumns; i++) {
            colSpec.put(i, ItemRenderer.class.getCanonicalName());
         }
      }

      final List<AdvancedDataGrid> result = new ArrayList<AdvancedDataGrid>();
      final Map<Integer, String> colSpecParam;
      if (colSpec == null) {
         colSpecParam = null;
      } else {
         colSpecParam = new HashMap<Integer, String>(colSpec);
      }

      try {
         // We try to get the grid until the timeout expires
         Object condition = new Object() {
            @Override
            public boolean equals(Object other) {
               AdvancedDataGrid dataGrid = DefaultAdvancedDataGrid.getInstance(
                     gridId.getValue(Property.DIRECT_ID),
                     true,
                     colSpecParam,
                     ((SeleAPLImpl) SUITA.Factory.apl()).getFlashSelenium());

               if (dataGrid != null) {
                  result.add(dataGrid);
               }
               return dataGrid != null;
            }
         };
         UI.condition.isTrue(condition).await(SUITA.Environment.getPageLoadTimeout());
      } catch (AssertionError err) {
         // Add error log and make screenshot when the grid is not found
         _logger.error("AdvancedDatagrid with ID: " + gridId.toString() + " is not found on the screem!");
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "Grid not found!");
         err.printStackTrace();
      }

      if (!result.isEmpty()) {
         return result.get(0);
      } else {
         return null;
      }
   }

   /**
    * Returns a list of the selected items indexes
    *
    * @param dataGrid - id of the data grid
    * @return - list of indexes for the selected items
    */
   public static List<Integer> getSelectedItemIndexes(AdvancedDataGrid dataGrid) {
      if (dataGrid == null) {
         throw new IllegalArgumentException("dataGrid is not set");
      }

      List<Integer> selectedRows = dataGrid.getSelectedRows();

      return selectedRows;
   }
}
