/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Provides common methods for the search view
 */
public class BaseSearchView extends BaseView {
   private static String ID_LINK_ADVANCED_SEARCH = "advancedSearchLink";
   private static String ID_LINK_SIMPLE_SEARCH = "simpleSearchLink";
   private static String ID_TABS_RESULTS = "resultTabs";
   private static String ID_GRID_RESULTS = "list";


   /**
    * Switches the view to "advanced search"
    */
   public void switchToAdvancedSearch() {
      UI.component.click(ID_LINK_ADVANCED_SEARCH);
   }

   /**
    * Switches the view to "simple search"
    */
   public void switchToSimpleSearch() {
      UI.component.click(ID_LINK_SIMPLE_SEARCH);
   }

   /**
    * Switches to a results tab
    * @param tabName - the name of the results tab to switch to
    */
   public void switchToTab(String tabName) {
      UI.component.value.set(tabName, ID_TABS_RESULTS);
   }

   /**
    * The count of results in the tab specified
    * @param tabName - the name of the results tab for the grid
    * @return count of results
    */
   public int getResultsCount(String tabName) {
      switchToTab(tabName);
      return GridControl.getRowsCount(getGrid(tabName));
   }

   /**
    * Gets the index of a grid item
    * @param tabName - the name of the results tab for the grid
    * @param name - the name of the item
    */
   public int getItemIndex(String tabName, String name) {
      switchToTab(tabName);
      return GridControl.getEntityIndex(getGrid(tabName), name);
   }

   // Private methods
   private AdvancedDataGrid getGrid(String tabName) {
      IDGroup gridId =
            IDGroup.toIDGroup("automationName=" + tabName + "/" + ID_GRID_RESULTS);
      return GridControl.findGrid(gridId);
   }
}
