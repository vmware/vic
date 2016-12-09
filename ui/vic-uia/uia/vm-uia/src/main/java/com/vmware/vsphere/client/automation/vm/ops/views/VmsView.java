/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ops.views;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;

/**
 * Represents the list in Datacenter > Related Objects > Virtual Machines
 */
public class VmsView extends BaseView {

   private static final String GRID_ID = "vmsForDatacenter/list";

   /**
    * Represents a column in the list
    */
   public enum Column {
      NAME("grid.vms.column.name"), STATE("grid.vms.column.state");

      private String value;

      Column(String valueKey) {
         this.value = VmUtil.getLocalizedString(valueKey);
      }

      public String getValue() {
         return value;
      }
   }

   /**
    * Gets the value of a cell specified by column at row specified by
    * vmName
    *
    * @param vmName
    *           - the name of the VM
    * @param columnName
    *           - the column
    * @return the value of the cell specified by column at row specified by
    *         vmName
    */
   public String getCellValue(String vmName, Column column) {
      AdvancedDataGrid grid = GridControl.findGrid(GRID_ID);
      int rowIndex =  GridControl.getEntityIndex(grid, Column.NAME.getValue(), vmName);
      return GridControl.getEntityColumnValue(grid, rowIndex, column.getValue());
   }
}
