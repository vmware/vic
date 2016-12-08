/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.ContextHeader;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * "vCenter > Hosts" is represented by current class.
 */
public class HostsView extends BaseView {

   private static final String GRID_NAME_COLUMN = CommonUtil
         .getLocalizedString("column.name");
   private static final IDGroup AI_EDIT_HOST_CUSTOMIZATIONS = IDGroup
         .toIDGroup("vsphere.core.host.editHostCustomizationsAction");
   private static final IDGroup AI_REMEDIATE = IDGroup
         .toIDGroup("vsphere.core.hostprofile.remediateHostAction");

   private final String gridId;

   /**
    * Instantiates a hosts view that operates on a grid containing hosts.
    *
    * @param gridId - Id of the hosts grid.
    */
   public HostsView(String gridId) {
      this.gridId = gridId;
   }

   /**
    * Invoke Edit Host Profile wizard from context menu of a host
    *
    * @param hostName - the name of the host to click
    */
   public void invokeEditHostCustomizationsContextMenu(String hostName) {
      rightClick(hostName);
      ActionNavigator.invokeMenuAction(AI_EDIT_HOST_CUSTOMIZATIONS);
   }

   /**
    * Invoke Remediate wizard from actions menu
    *
    * @param entityNames - entities in grid to select in order to invoke action on them
    */
   public void invokeRemediateActionsMenu(String... entityNames) {
      rightClick(entityNames);
      ActionNavigator.invokeMenuAction(AI_REMEDIATE);
   }

   /**
    * Finds and right-clicks on the specified item in the data grid. The method expects
    * that the current location is the items' view. If the items are more than 1, it
    * selects them all and then right clicks on the first one with keeping selection
    *
    * @param name Name of the item (s) to right click
    * @return True if the right-click is successful, false otherwise, e.g. if the item is
    *         not found in the data grid.
    */
   public boolean rightClick(String... name) {
      return GridControl.rightClickEntity(getGrid(), GRID_NAME_COLUMN, name);
   }

   /**
    * Checks whether the specified item is listed in the data grid
    *
    * @param name Name of the item to be searched for
    * @return true if the item is found in the grid, false otherwise
    */
   public boolean isFoundInGrid(String name) {
      return GridControl.getEntityIndex(getGrid(), name) >= 0;
   }

   /**
    * Gets the value of a cell specified by columnName at row specified by rowIndex
    *
    * @param rowIndex - the index of the row
    * @param columnName - the name of the column
    * @return the value of the cell specified by columnName at row specified by rowIndex
    */
   public String getCellValue(int rowIndex, String columnName) {
      AdvancedDataGrid grid = GridControl.findGrid(gridId);

      boolean isColumnVisible = GridControl.isColumnVisible(grid, columnName);
      if (!isColumnVisible) {
         ContextHeader header = grid.openShowHideCoumnHeader(grid.getColumnNames()[0]);
         header.selectDeselectColums(new String[] { columnName }, null);
         header.clickOK();

         // The grid needs to be initialized again because it was refreshed by clicking
         // "OK"
         grid = GridControl.findGrid(gridId);
      }

      return GridControl.getEntityColumnValue(grid, rowIndex, columnName);
   }

   /**
    * Gets the index of a item in the grid
    *
    * @param name - the name of the item
    * @return the index of a item in the grid or -1 if not found
    */
   public int getIndex(String name) {
      AdvancedDataGrid grid = GridControl.findGrid(IDGroup.toIDGroup(gridId));
      return GridControl.getEntityIndex(grid, GRID_NAME_COLUMN, name);
   }

   /**
    * Method that clicks Yes on the Remediation confirmation dialog, which popups when
    * more then one host is selected for remediation.
    */
   public void confirmRemediationPopup() {
      YesNoDialog expectedConfirmationDialog = YesNoDialog.CONFIRMATION;
      // if more then one host selected
      if (expectedConfirmationDialog.isVisible()) {
         expectedConfirmationDialog.clickYes();
      }
   }

   /**
    * Finds and returns the advanced data grid on the items view.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location is not
    *         the items view.
    */
   private AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(gridId));
   }
}
