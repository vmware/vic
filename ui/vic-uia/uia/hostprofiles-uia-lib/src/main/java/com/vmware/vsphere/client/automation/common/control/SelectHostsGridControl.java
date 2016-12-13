/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.control;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.flexui.componentframework.controls.mx.Image;
import com.vmware.flexui.componentframework.controls.mx.Label;
import com.vmware.flexui.componentframework.controls.mx.List;
import com.vmware.flexui.selenium.BrowserUtil;

public class SelectHostsGridControl {
   private static String ID_HOST_CLUSTER_AVAILABLE_LABEL = "%s/text=";
   private static String ID_EXPAND_COLLAPSE_IMAGE = "%s/Image_CollapseExpand_0_";

   private final String _gridId;

   public SelectHostsGridControl(String gridId) {
      _gridId = gridId;
   }

   /**
    * Method that clicks on an item in attach grid
    *
    * @return true if successful, false otherwise
    */
   public boolean selectItemInGrid(String name) {
      // SUITA doesn't use leftMouseClick method
      Label item =
            new Label(String.format(ID_HOST_CLUSTER_AVAILABLE_LABEL, _gridId) + name,
                  BrowserUtil.flashSelenium);
      if (item.isComponentExisting()) {
         item.leftMouseClick();
         return true;
      }
      return false;
   }

   /**
    * Method that clicks on all expand collapse images in the grid
    */
   public void expandAll() {
      int numRows = GridControl.getRowsCount(GridControl.findGrid(_gridId));

      for (int i = 0; i < numRows; i++) {
         // SUITA doesn't use leftmouseclick method
         Image img =
               new Image(String.format(ID_EXPAND_COLLAPSE_IMAGE, _gridId) + i,
                     BrowserUtil.flashSelenium);
         if (img.isVisibleOnPath()) {
            img.leftMouseClick();
         }
      }
   }

   /**
    * Method that performs multiple selection in the grid
    */
   public void multiSelectItemsInGrid(String... names) {
      List gridList = new List(_gridId, BrowserUtil.flashSelenium);

      for (String name : names) {
         gridList.selectAdditionalItemByIndex(String.valueOf(GridControl.getEntityIndex(
               GridControl.findGrid(_gridId),
               name)));
      }
   }
}
