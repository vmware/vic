/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.PortletControl;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.srv.common.spec.CdDvdDriveSpec.CdDvdDriveType;

/**
 * UI model of the Summary Tab view located VC > VM > Summary
 */
public class VmSummaryTabView extends BaseView {

   private final String ID_VM_HARDWARE_PORTLET = "vsphere.core.vm.hardwareSummaryView.chrome";
   private final String ID_HDD_VIEW = "viewDisk_Disk_?";
   private final String ID_VM_HARDWARE_HDD_PROP_GRID = "diskView";
   private final String ID_VM_HARDWARE_HDD_ARROW = "arrowImageDisk_?";
   private String ID_HDD_PROP_VIEW = ID_HDD_VIEW + "/"
         + ID_VM_HARDWARE_HDD_PROP_GRID;
   private final String ID_VM_HDD_VIEW_LABEL = ID_HDD_VIEW + "/"
         + "className=UITextField";
   // Capacity label without expand
   private final String ID_CAPACITY_LABEL = ID_HDD_VIEW + "/" + "capacityLabel";
   // Capacity and Location when hdd expanded
   private final String ID_VM_HDD_CAPACITY_LABEL = "diskCapacity/className=UITextField";
   private final String ID_VM_HDD_DATASTORE_LINK_LABEL = "datastoreLink/className=UITextField";
   // controls under CDRom property grid stack
   private final String ID_VM_CDROM_VIEW = "view_Cdrom_?";
   private final String ID_VM_CDROM_CONNECTION_LABEL = ID_VM_CDROM_VIEW
         + "/headerLabel";
   private final String ID_VM_CDROM_CONNECTION_BUTTON = ID_VM_CDROM_VIEW
         + "/menuButton";
   private final String ID_VM_CDROM_CONNECTION_MENU = "menu_view_Cdrom_?";
   private final String ID_VM_CDROM_CONNECTION_MENU_DISCONNECT_ITEM = ID_VM_CDROM_CONNECTION_MENU
         + "/automationName=Disconnect";
   private final String ID_VM_CDROM_CONNECTION_MENU_CONNECT_CL_ISO_ITEM = ID_VM_CDROM_CONNECTION_MENU
         + "/automationName=Connect to ISO image from a Content Library...";

   public void expandVmHardwarePortlet() {
      PortletControl hardwarePortlet = new PortletControl(
            ID_VM_HARDWARE_PORTLET);
      hardwarePortlet.expand();
   }

   /**
    * Get label for property including expand and collapse the view
    *
    * @param index
    *           - index of the device
    * @param labelId
    *           - ID of label
    * @return
    */
   public String getHddPropViewLabel(int index, String labelId) {
      String deviceIndex = String.valueOf(index + 1);
      String hddCapacityLabel, hddArrow = null;

      hddArrow = ID_VM_HARDWARE_HDD_ARROW.replace("?", deviceIndex);

      UI.component.click(hddArrow);
      isIdFound(ID_HDD_PROP_VIEW, deviceIndex);

      hddCapacityLabel = getFormattedIdWithIndex(labelId, deviceIndex);

      UI.component.click(hddArrow);

      return hddCapacityLabel;
   }

   /**
    * Formats IS containing "?" to String
    *
    * @param viewId
    *           - ID of property/view where "?" is encountered
    * @param index
    *           - index of the property/view
    * @return
    */
   private String getFormattedIdWithIndex(String viewId, String index) {
      return UI.component.value.get(IDGroup
            .toIDGroup(getFormattedIdStringWithIndex(viewId, index)));
   }

   /**
    * Formats ID containing "?" to String
    *
    * @param viewId
    *           - ID of property/view where "?" is encountered
    * @param index
    *           - index of the property/view
    * @return ID as String
    */
   private String getFormattedIdStringWithIndex(String viewId, String index) {
      return viewId.replace("?", index);
   }

   /**
    * Get hard disk datastore location label for according hdd index when hard
    * disk expanded Datastore location is visible only when vie wis expanded
    *
    * @param hddIndex
    *           - starts from 0
    * @return hard disk datastore location label
    */
   public String getHddDatastoreLabel(int hddIndex) {
      return getHddPropViewLabel(hddIndex, ID_VM_HDD_DATASTORE_LINK_LABEL);
   }

   /**
    * Get hard disk capacity label for according hdd index when hard disk
    * expanded
    *
    * @param hddIndex
    * @return hard disk label
    */
   public String getHddCapacityLabelExpanded(int hddIndex) {
      return getHddPropViewLabel(hddIndex, ID_VM_HDD_CAPACITY_LABEL);
   }

   /**
    * Get hard disk capacity label for according hdd index when hard disk not
    * expanded
    *
    * @param hddIndex
    * @return hard disk label
    */
   public String getHddCapacityLabel(int hddIndex) {
      String deviceIndex = String.valueOf(hddIndex + 1);
      isIdFound(ID_HDD_PROP_VIEW, deviceIndex);

      String hddCapacityLabel = getFormattedIdWithIndex(ID_CAPACITY_LABEL,
            deviceIndex);
      return hddCapacityLabel;
   }

   /**
    * Get hard disk label for according hdd index
    *
    * @param hddIndex
    * @return hard disk label
    */
   public String getHddLabel(int hddIndex) {
      String deviceIndex = String.valueOf(hddIndex + 1);

      isIdFound(ID_HDD_PROP_VIEW, deviceIndex);

      String hddLablel = getFormattedIdWithIndex(ID_VM_HDD_VIEW_LABEL,
            deviceIndex);

      return hddLablel;
   }

   /**
    * Get cdrom connection label for according cdrom index
    *
    * @param cdIndex
    * @return cdrom connection label
    */
   public String getCdRomConnectionLabel(int cdIndex) {
      String deviceIndex = String.valueOf(cdIndex);
      isIdFound(ID_VM_CDROM_VIEW, deviceIndex);

      String connectionLabel = getFormattedIdWithIndex(
            ID_VM_CDROM_CONNECTION_LABEL, deviceIndex);
      return connectionLabel;
   }

   /**
    * Return true if CD/DVD Drive exists
    *
    * @param index
    *           - index of property/view
    */
   public boolean doesCdRomDriveExist(String index) {
      waitForPageToRefresh();
      String driveId = getFormattedIdStringWithIndex(ID_VM_CDROM_VIEW, index);
      return UI.component.exists(driveId);
   }

   /**
    * Invokes cdrom's connection menu & selects "Disconnect" menu item
    *
    * @param driveType
    * @param driveIndex
    */
   public void selectDisconnectMenuItem(CdDvdDriveType driveType, int driveIndex) {
      String deviceIndex = String.valueOf(driveIndex);
      String menuButtonId = getFormattedIdStringWithIndex(
            ID_VM_CDROM_CONNECTION_BUTTON, deviceIndex);
      String menuId = getFormattedIdStringWithIndex(
            ID_VM_CDROM_CONNECTION_MENU_DISCONNECT_ITEM, deviceIndex);
      UI.component.click(menuButtonId);
      UI.component.click(menuId);
   }

   /**
    * Invokes cdrom's connection menu & selects
    * "Connect to ISO from a Content Library..." menu item
    *
    * @param driveType
    * @param driveIndex
    */
   public void selectConnectCLIsoMenuItem(CdDvdDriveType driveType,
         int driveIndex) {
      String deviceIndex = String.valueOf(driveIndex);
      String menuButtonId = getFormattedIdStringWithIndex(
            ID_VM_CDROM_CONNECTION_BUTTON, deviceIndex);
      String menuId = getFormattedIdStringWithIndex(
            ID_VM_CDROM_CONNECTION_MENU_CONNECT_CL_ISO_ITEM, deviceIndex);
      UI.component.click(menuButtonId);
      UI.component.click(menuId);
   }

   /**
    * Return true if ID us found
    *
    * @param viewId
    *           - ID of view found
    * @param index
    *           - index of property/view
    */
   void isIdFound(String viewId, String index) {
      waitForPageToRefresh();
      UI.condition.isFound(viewId.replace("?", index)).await(
            SUITA.Environment.getPageLoadTimeout());
   }

}
