/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Select Storage page of the Deploy OVF wizard
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.view.SelectStoragePage}
 */
@Deprecated
public class SelectStoragePage extends WizardNavigator {
   private static final String ID_GRID = "storageList";
   private static final String ID_DATASTORES_GRID = "datastoreList";

   /**
    * Select an item in the storage list, by its name in the Name column of the
    * table.
    *
    * @param name
    *           - name of the datastore to select
    * @return true if item is successfully selected, false otherwise
    */
   public boolean selectStorage(String storageName) throws Exception {
      return GridControl.selectEntities(getGrid(), storageName);
   }

   /**
    * Select a datastore item by its name in the Name column of the table.
    *
    * @param name
    *           - name of the datastore to select
    * @return true if item is successfully selected, false otherwise
    */
   public boolean selectDatastoreInDsCluster(String datastoreName)
         throws Exception {
      UI.condition.isFound(ID_DATASTORES_GRID).await(
            SUITA.Environment.getUIOperationTimeout());
      return GridControl.selectEntities(getDatastoresGrid(), datastoreName);
   }

   /**
    * Finds and returns the advanced data grid for the storage list.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the storages view.
    */
   private AdvancedDataGrid getGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(ID_GRID));
   }

   /**
    * Finds and returns the advanced data grid for the list of datastores in
    * datastore cluster.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the datastores view.
    */
   private AdvancedDataGrid getDatastoresGrid() {
      return GridControl.findGrid(IDGroup.toIDGroup(ID_DATASTORES_GRID));
   }
}