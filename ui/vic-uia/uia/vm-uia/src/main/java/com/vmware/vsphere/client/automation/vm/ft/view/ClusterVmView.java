/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Represents the VMs tab and Virtual Machines subtab of a cluster.
 */
public class ClusterVmView extends BaseView {

   private final IDGroup ID_VMS_GRID = IDGroup.toIDGroup("vmsForCluster/list");

   /**
    * A method to get the advanced data grid.
    *
    * @return The advanced data grid for VMs
    */
   private AdvancedDataGrid getVmGrid() {
      return GridControl.findGrid(ID_VMS_GRID);
   }

   /**
    * Checks if the VM is present in the grid
    *
    * @param String
    *           vmName - the name of the VM
    * @return true if the VM is found, false otherwise
    */
   public boolean isVmPresentInGrid(String vmName) {
      Integer rowIndex = getVmGrid().findItemByName(vmName);
      return rowIndex != null ? true : false;
   }
}