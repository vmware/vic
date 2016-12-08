/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.BaseDialogNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.common.DatastoreGlobalActions;

/**
 * The class represents the datastores list view.
 * The list can be found in vCenter - > Datastores.
 */
public class DatastoreListView extends BaseView {

   private static final String ID_GRID = "vsphere.core.viDatastores.itemsView/list";
   private static final String DATASTORES_GRID_NAME_COLUMN =
         CommonUtil.getLocalizedString("datastoresList.grid.nameColumn");

   /**
    * Launches New datastore wizard.
    */
   public static void launchNewDatastoreWizard() {
      UI.component.click(DatastoreGlobalActions.AI_CREATE_DATASTORE);
      new BaseDialogNavigator().waitForDialogToLoad();
   }

   /**
    * Select given VDC in VDC list
    *
    * @param vdcName    the name of the vDC to be selected
    * @return           true if the selection was successful, false otherwise
    */
   public static boolean selectDatastoreByName(String vdcName) {
      return GridControl.selectEntity(getGrid(), DATASTORES_GRID_NAME_COLUMN, vdcName);
   }

   //---------------------------------------------------------------------------
   // Private methods

   /**
    * Finds and returns the advanced data grid on the datastore list view.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    * is not the datastore view.
    */
   private static AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(ID_GRID));
   }
}
