/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Represents Select Name and device page in New Datastore wizard (VMFS flow)
 */
public class SelectNameAndLunPage extends WizardNavigator {

   private static final IDGroup ID_NAME_TF = IDGroup.toIDGroup("tiwoDialog/nameInput");
   private static final IDGroup ID_LUNS_LIST = IDGroup
         .toIDGroup("tiwoDialog/deviceList");

   /**
    * Types datastore name in the name text field.
    *
    * @param datastoreName
    */
   public void setDatastoreName(String datastoreName) {
      UI.component.value.set(datastoreName, ID_NAME_TF);
   }

   /**
    * @param rowIndex
    */
   public void selectDevice(int rowIndex) {
      AdvancedDataGrid deviceList = getGrid();
      if (deviceList.getRowsCount() > 0) {
         deviceList.selectRows(rowIndex);
      }
   }

   /**
    * Finds and returns the advanced data grid on Select Name and device view.
    *
    * @return data grid object, null if the current location
    *         is not the expected view.
    */
   private AdvancedDataGrid getGrid() {
      return GridControl.findGrid(ID_LUNS_LIST);
   }
}
