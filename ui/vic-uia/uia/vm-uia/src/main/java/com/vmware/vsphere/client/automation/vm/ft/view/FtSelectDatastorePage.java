/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Select Datastores page of the Fault Tolerance wizard
 */
public class FtSelectDatastorePage extends WizardNavigator {
   private static final String ID_DATASTORES_GRID = "storageList";

   /**
    * Select a datastore in the storage list by name in the Name column of the
    * table.
    *
    * @param name
    *           - name of the datastore to select
    * @return true if item is successfully selected, false otherwise
    */
   public boolean selectDatastore(String datastoreName) {
      return GridControl.selectEntities(getDatastoresGrid(), datastoreName);
   }

   /**
    * Finds and returns the advanced data grid for the list of datastores in
    * Select datastore page of Fault Tolerance dialog.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the datastores view.
    */
   private AdvancedDataGrid getDatastoresGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(ID_DATASTORES_GRID));
   }
}