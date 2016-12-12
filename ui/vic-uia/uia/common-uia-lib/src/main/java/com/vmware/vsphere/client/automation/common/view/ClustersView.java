/*
 *  Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * A view class representing vCenter > Clusters
 */
public class ClustersView extends BaseView {
   private static final String GRID_NAME_COLUMN = CommonUtil
         .getLocalizedString("column.name");

   private String gridId;

   /**
    * Instantiates a clusters view that operates on a grid containing clusters.
    *
    * @param gridId - Id of the clusters grid.
    */
   public ClustersView(String gridId) {
      this.gridId = gridId;
   }

   /**
    * Finds and right-clicks on the specified item in the data grid. The method expects
    * that the current location is the items' view. If the items are more than 1, it
    * selects them all and then right clicks on the first one with keeping selection
    *
    * @param name Name of the item (s) to right click
    * @return True if the right-click is successful, false otherwise, e.g. if the item is
    * not found in the data grid.
    */
   public boolean rightClick(String... name) {
      return GridControl.rightClickEntity(getGrid(), GRID_NAME_COLUMN, name);
   }

   protected AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(gridId));
   }
}
