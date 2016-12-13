/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.host.manage.graphics.view;

import java.util.Arrays;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.ContextHeader;
import com.vmware.flexui.componentframework.controls.mx.custom.ViClientPermanentTabBar;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.host.common.ManageHostUtil;

/**
 * UI model of the Host > Manage > Settings > Graphics > Graphics Devices view
 */
public class GraphicsDevicesView extends BaseView {

   private static final String ID_HOST_GRAPHICS_TAB_NAVIGATOR = "vsphere.core.host.manage.settings.graphicsTabs/tabNavigator";
   private static Integer graphicDeviceRowIndex = 0;
   private static String configTypeColumnName = ManageHostUtil
         .getLocalizedString(
               "host.edit.graphicdevices.settings.datagrig.column.config.type");
   /**
    * Graphics Devices > Finds and returns the advanced data grid for VMs
    * associated for a graphic device
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the graphic devices view.
    */
   private static final IDGroup ID_GRAPHICS_DEVICES_VMS_DATAGRID = IDGroup
         .toIDGroup("vsphere.core.host.manage.settings.graphics.vmList/list");

   private static final IDGroup ID_GRAPHICS_DEVICES_VMS_ASSOCIATED_EMPTY_LABEL = IDGroup
         .toIDGroup("emptyListIndicator");

   private static final IDGroup ID_GRAPHICS_DEVICES_MULTIPLE_VMS_SELECTED_LABEL = IDGroup
         .toIDGroup("multipleItemsSelectedWarning");

   private static final IDGroup ID_GRAPHICS_DEVICES_NO_VMS_SELECTED_LABEL = IDGroup
         .toIDGroup("noItemsSelectedWarning");

   private static final IDGroup ID_GRAPHICS_DEVICES_DATAGRID = IDGroup
         .toIDGroup("vsphere.core.host.manage.settings.graphics.vgaList/list");

   private static final IDGroup ID_GRAPHICS_DEVICES_VMS_ASSOCIATED_LABEL = IDGroup
         .toIDGroup(
               "vsphere.core.host.manage.settings.graphics.vmList/titleLabel");

   private static final IDGroup ID_GRAPHICS_DEVICES_TITLE_LABEL = IDGroup
         .toIDGroup(
               "vsphere.core.host.manage.settings.graphics.vgaList/titleLabel");
   /**
    * Graphics Devices > Gets the graphics devices Edit button
    */
   private static final IDGroup ID_GRAPHICS_DEVICES_EDIT_BTN = IDGroup
         .toIDGroup("dataGridToolbar/button");

   public static void selectGraphicsDevicesTab() {
      ViClientPermanentTabBar hostSettingsTabBar = new ViClientPermanentTabBar(
            ID_HOST_GRAPHICS_TAB_NAVIGATOR, BrowserUtil.flashSelenium);

      hostSettingsTabBar
            .waitForElementEnable(SUITA.Environment.getUIOperationTimeout());
      hostSettingsTabBar.selectTabAt("1");
   }

   /**
    * Click Edit Graphics Devices button
    */
   public static void clickEditGraphicsDevicesButton() {
      UI.component.click(ID_GRAPHICS_DEVICES_EDIT_BTN);
   }

   /**
    * Graphics Devices > Gets the tooltip of Edit graphics device button
    */
   public static String getGraphicDeviceEditBtnTooltip() {
      return UI.component.property.get(Property.TOOLTIP,
            ID_GRAPHICS_DEVICES_EDIT_BTN);
   }

   /**
    * Graphics Devices > Gets the graphics devices title label
    */
   public static String getGraphicsDevicesLabel() {
      return UI.component.property.get(Property.TEXT,
            ID_GRAPHICS_DEVICES_TITLE_LABEL);
   }

   /**
    * Graphics Devices > Gets the VMs associated label of the selected device
    */
   public static String getGraphicsDevicesVmsAssociatedLabel() {
      return UI.component.property.get(Property.TEXT,
            ID_GRAPHICS_DEVICES_VMS_ASSOCIATED_LABEL);
   }

   /**
    * Select the first graphic device in graphic devices datagrid
    *
    * @param graphicDevice
    */
   public static void selectGraphicDevice() {
      getDeviceGraphicsGrid().selectRows(graphicDeviceRowIndex);
   }

   /**
    * Select the first graphic device in graphic devices datagrid
    *
    * @param graphicDevice
    */
   public static void selectMultipleGraphicDevices(
         int multipleGraphicDeviceRowIndex) {
      getDeviceGraphicsGrid().selectRows(0, multipleGraphicDeviceRowIndex);
   }

   /**
    * Graphics Devices > Returns true if VM associated with graphic device found
    *
    * @param vmName
    * @return
    */
   public boolean isVmFoundInGrid(String vmName) {
      Integer rowIndex = getVmsGrid().findItemByName(vmName);
      return (rowIndex != null) ? true : false;
   }

   /**
    * Graphics Devices > Returns content of row index 0 for column Configured
    * Type
    *
    * @return
    */
   public static String getGraphicsDeviceConfigType() {
      return getDeviceGraphicsGrid().getCellValue(graphicDeviceRowIndex,
            configTypeColumnName);
   }

   /**
    * Graphics Devices > Returns the values in column Configured Type
    *
    * @return
    */
   public static String[] getConfigTypeColumnValues() {
      return getDeviceGraphicsGrid().getColumnContents(configTypeColumnName);
   }

   /**
    * Graphics Devices > Returns row number for Graphic devices datagrid
    *
    * @return
    */
   public static int getConfigTypeRowNumber() {
      return getDeviceGraphicsGrid().getRowsCount();
   }

   /**
    * Graphics Devices > Returns column names for Graphics devices datagrid
    *
    * @return
    */
   public static String[] getGraphicsDevicesColumns() {
      return getDeviceGraphicsGrid().getColumnNames();
   }

   /**
    * Graphics Devices > Returns column names for vms associated with graphic
    * device datagrid
    *
    * @return
    */
   public static String[] getVmsListColumns() {
      return getVmsGrid().getColumnNames();
   }

   /**
    * Get empty list label when no graphic device selected
    *
    * @return
    */
   public static String getEmptyVmsListText() {
      return UI.component.property.get(Property.TEXT,
            ID_GRAPHICS_DEVICES_VMS_ASSOCIATED_EMPTY_LABEL);
   }

   /**
    * Get text in vms datagrid when no devices selected
    *
    * @return
    */
   public static String getNoDevicesSelectedText() {
      return UI.component.property.get(Property.TEXT,
            ID_GRAPHICS_DEVICES_NO_VMS_SELECTED_LABEL);
   }

   /**
    * Get text in vms datagrid when multiple devices selected
    *
    * @return
    */
   public static String getMultipleDevicesSelectedText() {
      return UI.component.property.get(Property.TEXT,
            ID_GRAPHICS_DEVICES_MULTIPLE_VMS_SELECTED_LABEL);
   }

   /**
    * Common method that sort 2 arrays and then compares them
    *
    * @param arr1
    * @param arr2
    * @return
    */
   public static boolean arraysEquals(String[] arr1, String[] arr2) {
      Arrays.sort(arr1);
      Arrays.sort(arr2);
      return Arrays.equals(arr1, arr2);
   }

   /**
    * Common method for select and deselect of datagrid column
    *
    * @param columnName
    *           - name of column that to be selected/deselected
    * @param select
    *           - if true column is selected
    */
   public static void selectColumn(String columnName, boolean select) {
      boolean isColumnVisible = getDeviceGraphicsGrid()
            .isColumnVisible(columnName);
      ContextHeader header = getDeviceGraphicsGrid().openShowHideCoumnHeader(
            getDeviceGraphicsGrid().getColumnNames()[0]);
      if (!isColumnVisible && select) {
         header.selectDeselectColums(new String[] { columnName }, null);
      } else if (isColumnVisible && !select) {
         header.selectDeselectColums(null, new String[] { columnName });
      }
      header.clickOK();
   }

   /**
    * True if column is visible
    *
    * @param columnName
    *           - name of the column
    * @return
    */
   public static boolean isColumnVisible(String columnName) {
      return getDeviceGraphicsGrid().isColumnVisible(columnName);
   }

   /**
    * Graphics Devices > Finds and returns the advanced data grid for graphic
    * devices
    *
    * @return <code>AdvancedDataGrid</code> object, null if the current location
    *         is not the graphics devices view.
    */
   private static AdvancedDataGrid getDeviceGraphicsGrid() {
      return GridControl
            .findGrid(IDGroup.toIDGroup(ID_GRAPHICS_DEVICES_DATAGRID));
   }

   /**
    * Graphics Devices > Finds and returns the advanced data grid for vms
    * associated to some graphic device
    * 
    * @return<code>AdvancedDataGrid</code> object, null if the current location
    *                                      is not the graphics devices view.
    */
   private static AdvancedDataGrid getVmsGrid() {
      return GridControl
            .findGrid(IDGroup.toIDGroup(ID_GRAPHICS_DEVICES_VMS_DATAGRID));
   }

}
