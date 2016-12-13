/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.srv.common.util.FtUtil;

/**
 * Select host page of the Fault Tolerance wizard
 */
public class FtSelectHostPage extends WizardNavigator {
   private static final IDGroup GRID_ID = IDGroup.toIDGroup("tiwoDialog/list");
   private static final String NAME_COLUMN = FtUtil
         .getLocalizedString("ft.wizard.selecthost.nameColumn");

   /**
    * Select a host in the hosts list by name in the Name column of the table.
    *
    * @param name
    *           - name of the host to select
    * @return true if item is successfully selected, false otherwise
    */
   public boolean selectHost(String hostName) {
      return GridControl.selectEntity(getHostsGrid(), NAME_COLUMN, hostName);
   }

   /**
    * Finds and returns the advanced data grid for the list of hosts in Select
    * host page of Fault Tolerance dialog.
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the hosts view.
    */
   private AdvancedDataGrid getHostsGrid() {
      return GridControl.findGrid(GRID_ID);
   }
}