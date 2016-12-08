/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.ContextHeader;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * View representing the cluster list view under the related objects of a vDC
 */
public class VdcClustersView extends BaseView{
   private static final String GRID_ID = "clustersForVdc/list";
   private static final String CLUSTER_GRID_NAME_COLUMN =
         CommonUtil.getLocalizedString("column.name");
   private static final IDGroup ID_GRID_CLUSTERS =
         IDGroup.toIDGroup("clustersForVdc/list");
   private static final String VDC_COLUMN =
         CommonUtil.getLocalizedString("clustersGrid.vdcColumn");

   /**
    * Returns the count clusters in the grid
    * @return the count clusters in the grid
    */
   public int getClustersCount() {
      return GridControl.getRowsCount(getGrid());
   }

   /**
    * Returns the index of the cluster in grid or -1 if not found
    * @param clusterName
    * @return the index of the cluster in grid or -1 if not found
    */
   public int getClusterIndex(String clusterName) {
      return GridControl.getEntityIndex(getGrid(), clusterName);
   }

   /**
    * Checks whether the specified cluster is listed in the data grid
    *
    * @param clusterName Name of the cluster to be searched for
    *
    * @return true if the cluster is found in the grid, false otherwise
    */
   public boolean isFoundInGrid(String clusterName) {
      return getClusterIndex(clusterName) != -1;
   }

   /**
    * Gets the value of a cell in the grid
    * @param rowIndex - the index of the row we are looking for
    * @param columnName - the name of the column we are looking for    *
    * @return the value of a cell in the grid
    */
   public String getCellValue(int rowIndex, String columnName) {
      return getGrid().getCellValue(rowIndex, columnName);
   }

   /**
    * Right-clicks on given cluster in clusters list
    *
    * @param clusterNames the names of the clusters to be right-clicked on
    * @return true if the right-click is successful, false otherwise, e.g. if
    *    the cluster is not found in the list.
    */
   public boolean rightClick(String ... clusterNames) {
      return GridControl.rightClickEntity(getGrid(), CLUSTER_GRID_NAME_COLUMN, clusterNames);
   }

   /**
    * Clicks on a cluster link in the grid
    * @param clusterName - the name of the cluster
    */
   public void clickOnClusterlink(String clusterName){
      AdvancedDataGrid grid = GridControl.findGrid(ID_GRID_CLUSTERS);
      int rowIndex = GridControl.getEntityIndex(grid, CommonUtil.getLocalizedString("column.name"), clusterName);
      grid.clickCell(rowIndex, CommonUtil.getLocalizedString("column.name"));
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

   /**
    * Unhides the given column
    */
   public void unhideVdcColumn(){
      if ( ! UI.condition.isFound( VDC_COLUMN ).await(
            SUITA.Environment.getUIOperationTimeout() / 5)) {

      AdvancedDataGrid grid = GridControl.findGrid(ID_GRID_CLUSTERS);
      ContextHeader header = grid.openShowHideCoumnHeader(CLUSTER_GRID_NAME_COLUMN);

      header.selectDeselectColums(
            new String[] { VDC_COLUMN }, null);
      header.clickOK();
      }

   }
}
