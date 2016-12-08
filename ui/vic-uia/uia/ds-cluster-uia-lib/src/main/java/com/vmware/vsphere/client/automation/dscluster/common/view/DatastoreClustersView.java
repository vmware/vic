/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.ContextHeader;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.dscluster.common.DsClusterUtil;

/**
 * "vCenter > Datastore Clusters" is represented by current class.
 */
public class DatastoreClustersView extends BaseView {
   private static final String DS_CLUSTERS_GRID_ID = "vsphere.core.viDsCluster.itemsView/list";
   private static final String CREATE_NEW_DS_CL_BUTTON_ID = "vsphere.core.dscluster.createActionGlobal";

   private static final String DS_CLUSTERS_GRID_NAME_COLUMN =
          DsClusterUtil.getLocalizedString("grid.dsclusters.column.name");

   /**
    * Finds and right-clicks on the specified datastore cluster in the data grid.
    * The method expects that the current location is the datastore clusters' view.
    *
    * @param dsClusterName Name of the datastore cluster
    *
    * @return True if the right-click is successful, false otherwise, e.g. if
    *    the datastore cluster is not found in the data grid.
    */
   public boolean rightClickDsCluster(String dsClusterName) {
      // Workaround, since sometimes the grid is not loaded
      // but the waiting methods don't work with the GRID_ID
      UI.condition.isFound(CREATE_NEW_DS_CL_BUTTON_ID).await(
            UiDelay.PAGE_LOAD_TIMEOUT.getDuration());
      return GridControl.rightClickEntity(
            getGrid(),
            DS_CLUSTERS_GRID_NAME_COLUMN,
            dsClusterName);
   }

   /**
    * Checks whether the specified datastore cluster is listed in the data grid
    *
    * @param dsClusterName Name of the datastore cluster to be searched for
    *
    * @return true if the datastore cluster is found in the grid, false otherwise
    */
   public boolean isFoundInGrid(String dsClusterName) {
      Integer rowIndex = getGrid().findItemByName(dsClusterName);
      boolean result = rowIndex != null;

      return result;
   }

   /**
    * Gets the value of a cell specified by columnName at row specified by rowIndex
    * @param rowIndex - the index of the row
    * @param columnName - the name of the column
    * @return the value of the cell specified by columnName at row specified by rowIndex
    */
   public String getCellValue(int rowIndex, String columnName) {
      AdvancedDataGrid grid = getGrid();

      boolean isColumnVisible = grid.isColumnVisible(columnName);
      if (!isColumnVisible) {
         ContextHeader header = grid.openShowHideCoumnHeader(grid.getColumnNames()[0]);
         header.selectDeselectColums(new String[] { columnName }, null);
         header.clickOK();

         // The grid needs to be initialized again because it was refreshed by clicking "OK"
         grid = GridControl.findGrid(DS_CLUSTERS_GRID_ID);
      }

      return GridControl.getEntityColumnValue(grid, rowIndex, columnName);
   }

   /**
    * Gets the index of a Datastore Cluster in the grid
    * @param dsClusterName - the name of the Datastore Cluster
    * @return the index of a Datastore Cluster in the grid or -1 if not found
    */
   public int getDsClusterIndex(String dsClusterName) {
      AdvancedDataGrid grid = getGrid();
      return GridControl.getEntityIndex(grid, DS_CLUSTERS_GRID_NAME_COLUMN, dsClusterName);
   }

   /**
    * Finds and returns the advanced data grid on the datastore clusters view.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    * is not the datastore clusters view.
    */
   private AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(DS_CLUSTERS_GRID_ID));
   }
}
