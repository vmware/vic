/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import org.apache.commons.lang.NotImplementedException;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Common operations on all related objects lists.
 * TODO iparaskov: Give more details for the purpose of the class and
 * how to determine the set of operations to cover in the class. Based on that will
 * design and implement the structure of the class.
 */
public class RelatedObjectsView extends BaseView {

   private final IDGroup listId;
   private final String toolBarButtonId;

   /**
    * Init the set of RO(Related Object view) identifies.
    * @param roPath   path to the specific RO page.
    * For example for the Datacenter > Related Objects > Distributed Switches view
    * the id to determine the gird is dvsForDatacenter/list.
    */
   public RelatedObjectsView(String roPath) {
      super();
      this.listId = IDGroup.toIDGroup(roPath + "/list");
      this.toolBarButtonId = roPath + "/allActions/button";
   }


   /**
    * Clicks on toolbar button.
    */
   public void clickToolbarButton() {
      throw new NotImplementedException("Implement me");
   }

   public void invokeActionForSelectedItem(IDGroup menuItemID) {
      ActionNavigator.invokeFromToolbarMenu(toolBarButtonId,
            menuItemID);
   }

   /**
    * Checks whether the specified item is listed in the related objects list.
    *
    * @param itemName
    *           the name of the related item to be searched for
    * @throws Exception
    */
   public boolean isFoundInGrid(String itemName) throws Exception {
      Integer rowIndex = getGrid().findItemByName(itemName);
      return (rowIndex != null) ? true : false;
   }

   /**
    * Selects the specified item in the related objects list.
    *
    * @param itemName
    *           the name of the related item to be selected
    * @throws Exception
    */
   public void selectItem(String itemName) throws Exception {
      Integer rowIndex = getGrid().findItemByName(itemName);
      if (rowIndex > -1) {
         getGrid().selectRows(rowIndex);
      }
   }

   /**
    * Finds and returns the advanced data grid on the distributed switches view.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the distributed switches view.
    * @throws Exception
    */
   private AdvancedDataGrid getGrid() throws Exception {
      return GridControl.findGrid(this.listId);
   }
}
