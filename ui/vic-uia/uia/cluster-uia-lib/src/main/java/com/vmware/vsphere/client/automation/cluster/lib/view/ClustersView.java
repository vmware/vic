/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.ContextHeader;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * "vCenter > Clusters" is represented by current class.
 */
public class ClustersView extends BaseView {
   private static final String GRID_ID = "vsphere.core.viClusters.itemsView/list";
   private static final String CREATE_NEW_BUTTON_ID = "vsphere.core.cluster.createActionGlobal";

   private static final String CLUSTERS_GRID_NAME_COLUMN = CommonUtil
         .getLocalizedString("clustersGrid.clusterNameColumnName");

   /**
    * Finds and right-clicks on the specified cluster in the data grid. The method
    * expects that the current location is the clusters' view.
    *
    * @param clusterName Name of the cluster
    *
    * @return True if the right-click is successful, false otherwise, e.g. if
    *    the cluster is not found in the data grid.
    */
   public boolean rightClickCluster(String clusterName) {
      // Workaround, since sometimes the grid is not loaded
      // but the waiting methods don't work with the GRID_ID
      UI.condition.isFound(CREATE_NEW_BUTTON_ID).await(
            UiDelay.PAGE_LOAD_TIMEOUT.getDuration());
      return GridControl.rightClickEntity(
            getGrid(),
            CLUSTERS_GRID_NAME_COLUMN,
            clusterName);
   }

   /**
    * Checks whether the specified cluster is listed in the data grid
    *
    * @param clusterName Name of the cluster to be searched for
    *
    * @return true if the cluster is found in the grid, false otherwise
    */
   public boolean isFoundInGrid(String clusterName) {
      Integer rowIndex = getGrid().findItemByName(clusterName);
      boolean result = (rowIndex != null) ? true : false;

      return result;
   }

   /**
    * Gets the value of a cell specified by columnName at row specified by rowIndex
    * @param rowIndex - the index of the row
    * @param columnName - the name of the column
    * @return the value of the cell specified by columnName at row specified by rowIndex
    */
   public String getCellValue(int rowIndex, String columnName) {
      AdvancedDataGrid grid = GridControl.findGrid(GRID_ID);

      boolean isColumnVisible = grid.isColumnVisible(columnName);
      if (!isColumnVisible) {
         ContextHeader header = grid.openShowHideCoumnHeader(grid.getColumnNames()[0]);
         header.selectDeselectColums(new String[] { columnName }, null);
         header.clickOK();

         // The grid needs to be initialized again because it was refreshed by clicking "OK"
         grid = GridControl.findGrid(GRID_ID);
      }

      return GridControl.getEntityColumnValue(grid, rowIndex, columnName);
   }

   /**
    * Gets the index of a Cluster in the grid
    * @param name - the name of the Cluster
    * @return the index of a Cluster in the grid or -1 if not found
    */
   public int getClusterIndex(String name) {
      AdvancedDataGrid grid = GridControl.findGrid(IDGroup.toIDGroup(GRID_ID));
      return GridControl.getEntityIndex(grid, CLUSTERS_GRID_NAME_COLUMN, name);
   }

   /**
    * Finds and returns the advanced data grid on the clusters view.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    * is not the clusters view.
    */
   private AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(GRID_ID));
   }
}
