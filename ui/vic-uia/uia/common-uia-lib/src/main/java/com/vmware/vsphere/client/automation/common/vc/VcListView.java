/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.vc;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.vcuilib.commoncode.ActionFunction;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.CommonUtil;

public class VcListView extends BaseView {

   private static final Logger logger = LoggerFactory.getLogger(ActionFunction.class);

   private static final String VC_GRID_NAME_COLUMN = CommonUtil
         .getLocalizedString("vc.column.name");
   private static final String ID_GRID = "list";
   private static final IDGroup ID_NEW_VAPP_FROM_LIBRARY_ICON = IDGroup
            .toIDGroup("vsphere.core.folder.relatedVMs/vappsForVCenter/vsphere.core.vApp.deployAction.global/button");

   /**
    * Executes an action from the context menu of the specified VC item.
    *
    * @param vcName        the name of the VC to be right-clicked on
    * @param menuItemID    id of the menu item to be clicked.
    * @param subMenuIDs    this parameter is ignored
    */
   public void executeContextMenuAction(String vcName, IDGroup menuItemID,
         IDGroup... subMenuIDs) {
      if (!rightClickVc(vcName)) {
         throw new RuntimeException("Right-click on the DC failed: " + vcName);
      }

      ActionNavigator.invokeMenuAction(menuItemID);
   }

   /**
    * Right-clicks on given VCs in VC list
    *
    * @param vcNames      the name of the VC to be right-clicked on
    * @return              true if the right-click is successful, false otherwise,
    *                      e.g. if the VC is not found in the VC list.
    */
   public boolean rightClickVc(String... vcNames) {
      GridControl.rightClickEntity(getGrid(), VC_GRID_NAME_COLUMN, vcNames);

      try {
         ActionNavigator.waitForContextMenu(BrowserUtil.flashSelenium, true);
      } catch (Exception e) {
         logger.error("Context menu was not opened correcly on VCs list");
         return false;
      }

      return true;
   }

   /**
    * Method that clicks on the New vApp From Library icon
    */
   public boolean clickNewVappFromLibraryIcon() {
       boolean result = false;

       if (UI.component.exists(ID_NEW_VAPP_FROM_LIBRARY_ICON)) {
           UI.component.click(ID_NEW_VAPP_FROM_LIBRARY_ICON);
           result = true;
       }
       return result;
   }

   /**
    * Finds and returns the advanced data grid on the VC list view.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    * is not the VC view.
    */
   private AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(ID_GRID));
   }
}
