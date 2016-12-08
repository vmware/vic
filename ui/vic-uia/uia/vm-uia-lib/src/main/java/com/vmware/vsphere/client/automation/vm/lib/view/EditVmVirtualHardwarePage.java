/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.control.StackBlockControl;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.flexui.componentframework.controls.mx.ComboBox;
import com.vmware.flexui.componentframework.controls.mx.TextInput;
import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.srv.common.messages.VmMessages;
import com.vmware.vsphere.client.automation.srv.common.spec.CdDvdDriveSpec.CdDvdDriveType;
import com.vmware.vsphere.client.automation.srv.common.spec.VirtualDiskSpec;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * VM > Edit Settings dialog
 */
public class EditVmVirtualHardwarePage extends SinglePageDialogNavigator {

   private static final VmMessages messages = I18n.get(VmMessages.class);
   // Add Device controls IDs
   public static final String ID_ADD_NEW_DEVICE_BUTTON = "addDeviceButton";
   public static final String ID_EDIT_ADD_DEVICE_COMBO = "addHardware/popUpButton";
   // Add Device's Popup menu item labels
   private static final String ID_EDIT_ADD_DEVICE_MENU = "menu";
   private static final String ID_NEW_HARD_DISK_LABEL = "automationName=New Hard Disk";
   private static final String ID_NEW_HDD_MENU_ITEM = ID_EDIT_ADD_DEVICE_MENU
         + "/" + ID_NEW_HARD_DISK_LABEL;
   private static final String ID_NEW_EXISTING_HDD_MENU_ITEM = ID_EDIT_ADD_DEVICE_MENU
         + "/automationName=Existing Hard Disk";
   private static final String ID_NEW_CDDVDDRIVE_DEVICE_LABEL = "automationName=CD\\\\/DVD Drive";
   private static final IDGroup ID_NEW_CDDVDDRIVE_MENU_ITEM = IDGroup
         .toIDGroup(ID_EDIT_ADD_DEVICE_MENU + "/"
               + ID_NEW_CDDVDDRIVE_DEVICE_LABEL);
   // HDD device controls IDs
   private static final String DISK_INDEX = "<DiskIndex>";
   private static final String ID_DISK = "vmConfigControl/disk_" + DISK_INDEX
         + "/diskSize/capacity/textDisplay";
   // CD/DVD Drive device controls IDs
   private static final String ID_BROWSE_BUTTON = "hardwareStack/cdrom_%s/browse";
   private static final String ID_CONNECT_AT_STARTUP_CHECKBOX = "hardwareStack/cdrom_%s/startConnected";
   private static final String ID_CDDVDDRIVE_TYPE_COMBOBOX = "tiwoDialog/hardwareStack/cdrom_%s/type";
   private static final String ID_CDDVDDRIVE_SECTION_TITLE = "hardwareStack/cdrom_?/titleLabel.cdrom_?";
   private static final String ID_CDDVDDRIVE_REMOVE_BUTTON = "hardwareStack/cdrom_?/remove";
   private static final String ID_CDDVDDRIVE_WARNING_ICON = "hardwareStack/cdrom_?/className=Image[2]";
   private static final String ID_CDDVDDRIVE_MEDIAFILE_TEXTBOX = "hardwareStack/cdrom_?/bodySection/file";
   // Shared PCI Device
   private static final String ID_NEW_SHARED_PCI_LABEL = "automationName=Shared PCI Device";
   private static final String ID_NEW_SHARED_PCI_MENU_ITEM = ID_EDIT_ADD_DEVICE_MENU
         + "/" + ID_NEW_SHARED_PCI_LABEL;
   private static final String ID_RESERVE_ALL_MEMORY_BTN = "_PciSharedPage_Button1";
   private static final String ID_VGPU_WARNING_LABEL = "_PciSharedPage_Text1";
   private static final String ID_VGPU_NOTE_LABEL = "_PciSharedPage_Text2";
   private static final String ID_VGPU_PROFILES_COMBOBOX = "vgpu";
   private static final String ID_TITLE_LABEL = "titleLabel";
   // Only 1 Shared PCI Device can be added
   private static final String ID_PCI_DEVICE_LABEL = "PCI device 0";
   private static final String ID_SHARED_PCI_TITLE_LABEL = ID_TITLE_LABEL + "."
         + ID_PCI_DEVICE_LABEL;

   private static final String ID_SHARED_PCI_TYPE_DROPDOWN = ID_PCI_DEVICE_LABEL
         + "/" + "type";
   private static final String ID_SHARED_PCI_REMOVE = ID_PCI_DEVICE_LABEL + "/"
         + "remove";

   private static final String ID_SHARED_PCI_TXT_REMOVE = "_PciSharedPage_Label1";

   private static final String ID_SHARED_PCI_EXPAND = ID_PCI_DEVICE_LABEL + "/"
         + "arrowImagenull";

   private static final String STORAGE_LOCATOR_ID = "storageList";
   private static final String DIALOG_OK_BUTTON_ID = "buttonOk";
   private static final int HUNDRED_MS = 100;
   private static final int SECOND = 1000;

   public static enum StackBlockEnum {
      HARD_DISK("hardwareStack/disk_%d/");
      private String id;

      StackBlockEnum(String id) {
         this.id = id;
      }

      public String id() {
         return this.id;
      }

      public String getBlockTitleId(int index) {
         return String.format(this.id(), index) + "titleLabel.disk_" + index;
      }

      public String getBlockPropertyGridId(int index) {
         return String.format(this.id(), index) + "bodySection";
      }
   }

   public static enum ComboEnum {
      PROVISIONING("provisioningList"), SHARING("sharing"), LOCATION("location");
      private String id;

      ComboEnum(String id) {
         this.id = id;
      }

      public String id() {
         return this.id;
      }
   }

   // Add Devices control ops //
   /**
    * Launch Add Devices menu in VM Summary>Edit Settings
    */
   public static void openHwDevicesMenu() {
      // Click to open Hardware devices menu
      UI.component.click(ID_EDIT_ADD_DEVICE_COMBO);

      // Wait until the menu is loaded
      UI.condition.isFound(ID_EDIT_ADD_DEVICE_MENU).await(
            SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Click Add button in VM Summary>Edit Settings>New Device
    */
   public static void clickAddDevice() {
      // Click "Add" button to confirm the choice
      UI.component.click(ID_ADD_NEW_DEVICE_BUTTON);
   }

   /**
    * Clicks the "New Hard Disk" item in "New Device" menu
    */
   public static void selectHddAddDevice() {
      UI.component.click(ID_NEW_HDD_MENU_ITEM);
   }

   /**
    * Clicks on the 'Existing Hard Disk' in 'New Device' menu.
    */
   public static void selectHddAddExistingDevice() {
      UI.component.click(ID_NEW_EXISTING_HDD_MENU_ITEM);
   }

   /**
    * Clicks the "CD/DVD Drive" item in "New Device" menu
    */
   public static void selectCdDvdDriveMenuItem() {
      UI.component.click(ID_NEW_CDDVDDRIVE_MENU_ITEM);
   }

   // HDD device control ops //
   /**
    * Sets the size of the virtual disk
    *
    * @param hddIndex
    *           - virtual disk index (0 is the first default disk)
    * @param sizeInGb
    *           - size in GB - default one
    * @throws Exception
    */
   public void setDiskSize(final int hddIndex, final int sizeInGb)
         throws Exception {
      final IDGroup disk = IDGroup.toIDGroup(ID_DISK.replace(DISK_INDEX,
            Integer.toString(hddIndex)));
      waitForControlEnabled(IDGroup.toIDGroup(disk));
      UI.component.value.set(sizeInGb, disk);
   }

   // CD/DVD Drive device controls ops //
   /**
    * Click browse button for the given drive index
    *
    * @param driveIndex
    */
   public static void clickBrowse(String driveIndex) {
      UI.component.click(generateIdwithIndex(ID_BROWSE_BUTTON, driveIndex));
   }

   /**
    * Click Remove (x) button for the given drive index
    *
    * @param driveIndex
    */
   public static void clickRemoveCdRom(String driveIndex) {
      String id = ID_CDDVDDRIVE_REMOVE_BUTTON.replace("?", driveIndex);
      UI.component.click(id);
   }

   /**
    * Method that expands or collapses the cdRom stack block
    *
    * @param expand
    *           - if true, expands it, if false - collapses it
    */
   public static void expandCollapseCdDvdDriveStackblock(boolean expand,
         String driveIndex) {
      String id = ID_CDDVDDRIVE_SECTION_TITLE.replace("?", driveIndex);
      if (StackBlockControl.isExpanded(id) != expand) {
         UI.component.click(id);
      }
   }

   /**
    * Check Connect at Power On checkbox for the given drive index
    *
    * @param driveIndex
    * @param value
    *           if true, checks box
    */
   public static void setConnectPowerOn(String driveIndex, boolean value) {
      UI.component.value.set(value,
            generateIdwithIndex(ID_CONNECT_AT_STARTUP_CHECKBOX, driveIndex));
   }

   /**
    * Selects drive type
    *
    * @param driveType
    * @param driveIndex
    */
   public static void selectCdDvdDriveType(CdDvdDriveType driveType,
         String driveIndex) {
      String controlID = generateIdwithIndex(ID_CDDVDDRIVE_TYPE_COMBOBOX,
            driveIndex);
      UI.component.value.set(driveType.value(), controlID);
   }

   /**
    * Selects drive type by dropdown index
    *
    * @param driveTypeIndex
    *           index of item in dropdown
    * @param driveIndex
    */
   public static void selectCdDvdDriveTypeByIndex(int driveTypeIndex,
         String driveIndex) {
      String controlID = generateIdwithIndex(ID_CDDVDDRIVE_TYPE_COMBOBOX,
            driveIndex);
      UI.component.selectByIndex(driveTypeIndex, controlID);
   }

   /**
    * Get error icon's visibility
    *
    * @param driveIndex
    *           of cdRom drive
    */
   public static boolean isErrorIconVisible(String driveIndex) {
      String id = ID_CDDVDDRIVE_WARNING_ICON.replace("?", driveIndex);
      return UI.condition.isFound(IDGroup.toIDGroup(id)).await(
            SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Get Media's file path
    *
    * @param driveIndex
    *           of cdRom drive
    */
   public String getMediaPath(String driveIndex) {
      String id = ID_CDDVDDRIVE_MEDIAFILE_TEXTBOX.replace("?", driveIndex);
      TextInput textInput = new TextInput(id, BrowserUtil.flashSelenium);
      return textInput.getProperty("text", false);
   }

   // Utility methods
   private static String generateIdwithIndex(String id, String driveIndex) {
      return String.format(id, driveIndex);
   }

   /**
    * Select New Shared PCI Device to be added
    */
   public static void selectSharedPciAddDevice() {
      UI.component.click(ID_NEW_SHARED_PCI_MENU_ITEM);
   }

   /**
    * Clicks on Reserve all memory button
    */
   public void clickReserveAllMemoryBtn() {
      UI.component.click(ID_RESERVE_ALL_MEMORY_BTN);
   }

   /**
    * Get label of warning message that appears when Shared PCI Device added
    *
    * @return
    */
   public String getVgpuWarning() {
      return UI.component.property.get(Property.TEXT,
            IDGroup.toIDGroup(ID_VGPU_WARNING_LABEL));
   }

   /**
    * Get label of note message that appears when Shared PCI Device added
    *
    * @return
    */
   public String getVgpuNoteLabel() {
      return UI.component.property.get(Property.TEXT,
            IDGroup.toIDGroup(ID_VGPU_NOTE_LABEL));
   }

   /*
    * @return true if vgpu warning visible
    */
   public boolean isVgpuWarningVisible() {
      return UI.condition.isFound(ID_VGPU_WARNING_LABEL).estimate();
   }

   /**
    * @return true if Reserve All Memory button visible
    */
   public boolean isReserveAllMemoryBtnVisible() {
      return UI.condition.isFound(ID_RESERVE_ALL_MEMORY_BTN).estimate();
   }

   /**
    * Select vGPU profile
    *
    * @param valueToSelect
    *           - vGPU profile to be selected
    */
   public void selectVgpuProfile(String valueToSelect) {
      ComboBox dropDown = new ComboBox(ID_VGPU_PROFILES_COMBOBOX,
            BrowserUtil.flashSelenium);
      dropDown.selectByLabel(valueToSelect);
   }

   /**
    * Get selected vGPU Profile
    *
    * @return selected vGPU Profile
    */
   public String getSelectedVgpuProfile() {
      ComboBox dropDown = new ComboBox(ID_VGPU_PROFILES_COMBOBOX,
            BrowserUtil.flashSelenium);
      return dropDown.getSelectedLabel();
   }

   /**
    * Return true if ID is found
    *
    * @param viewId
    *           - ID of view found
    * @param index
    *           - index of property/view
    */
   public void isIdFound(String viewId, String index) {
      UI.condition.isFound(viewId.replace("?", index)).await(
            SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Selected Shared PCI Device type
    *
    * @return selected Shared PCI Device type
    */
   public String getSharedPciDeviceType() {
      SparkDropDownList vgpuDeviceType = new SparkDropDownList(
            ID_SHARED_PCI_TYPE_DROPDOWN, BrowserUtil.flashSelenium);
      return vgpuDeviceType.getSelectedItemText();
   }

   /**
    * Return "Device will be removed" label when Remove clicked
    *
    * @return title label for Shared PCI Device
    */
   public String getSharedPciTitleLabel() {
      return UI.component.property.get(Property.TEXT,
            IDGroup.toIDGroup(ID_SHARED_PCI_TITLE_LABEL));
   }

   /**
    * Click Remove button to remove Shared Pci Device
    */
   public void clickRemoveSharedPciDevice() {
      // Click "Remove" button to remove Shared Pci Device
      UI.component.click(ID_SHARED_PCI_REMOVE);
   }

   /**
    * Return "Device will be removed" label when Remove clicked
    *
    * @return "Device will be removed" label when Remove clicked
    */
   public String getVgpuDeviceRemoveLabel() {
      return UI.component.property.get(Property.TEXT,
            IDGroup.toIDGroup(ID_SHARED_PCI_TXT_REMOVE));
   }

   /**
    * Return true if vGPU profiles combo enabled
    *
    * @return true if Remove button enabled
    */
   public boolean getVgpuDeviceRemoveBtnEnabled() {
      return Boolean.valueOf((UI.component.property.get(Property.ENABLED,
            IDGroup.toIDGroup(ID_SHARED_PCI_REMOVE))));
   }

   /**
    * Return true if vGPU profiles combo enabled
    *
    * @return true if vGPU profiles combo enabled
    */
   public boolean getVgpuProfileEnabled() {
      ComboBox dropDown = new ComboBox(ID_VGPU_PROFILES_COMBOBOX,
            BrowserUtil.flashSelenium);
      return dropDown.getEnabled();
   }

   /**
    * Expand Shared Pci Device stack block
    */
   public void expandSharedPciDeviceStackblock() {
      if (!UI.component.exists(ID_VGPU_PROFILES_COMBOBOX)) {
         UI.component.click(ID_SHARED_PCI_EXPAND);
      }
   }

   /**
    * Return true if Shared Pci Device is present
    *
    * @return true if Shared Pci device added
    */
   public boolean isSharedPciDevicePresent() {
      return UI.component.exists(ID_SHARED_PCI_TITLE_LABEL);
   }

   /**
    * Toggles the state of the first stack block of the given type.
    *
    * @param stackBlock
    *           the type of the stack block
    */
   public static void toggleStackblock(StackBlockEnum stackBlock) {
      toggleStackblock(stackBlock, 0);
   }

   /**
    * Toggles the state of a stack block of the given type.
    *
    * @param stackBlock
    *           the type of the stack block
    * @param index
    *           zero-based index of the stack block of the given type
    */
   public static void toggleStackblock(StackBlockEnum stackBlock, int index) {
      UI.component.click(stackBlock.getBlockTitleId(index));
   }

   /**
    * Selects the given value in the given combo.
    *
    * @param combo
    *           the combo to operate on
    * @param value
    *           the value to set
    */
   public static void selectComboValue(ComboEnum combo, String value) {
      UI.component.value.set(value, combo.id());
   }

   /**
    * Returns the path to the file where the disk is stored.
    *
    * @param index
    *          the 0-based index of the Hard Disk block whose
    *          "Disk File" property value to look for.
    * @see {@link VirtualDiskSpec#getAbsolutePath()}
    * @return value
    *          value of the "Disk File" property
    */
   public static String getDiskFile(int index) {
      return UI.component.property.get(Property.TEXT,
            StackBlockEnum.HARD_DISK.getBlockPropertyGridId(index) + "/file");
   }

   /**
    * Selects the given datastore for the 'Location' field.
    *
    * @param dsName
    *           name of the datastore to select
    * @return true if the datastore selection dialog was closed
    */
   public static boolean selectLocationDatastore(String dsName) {
      EditVmVirtualHardwarePage.selectComboValue(ComboEnum.LOCATION,
            "Browse...");
      selectInGrid(STORAGE_LOCATOR_ID, messages.getNameColumnHeader(), dsName);
      return closeDialog(DIALOG_OK_BUTTON_ID);
   }

   private static boolean closeDialog(String closeBtnId) {
      boolean isCloseButtonEnabled = UI.condition.isTrue(
            UI.component.property
                  .getBoolean(Property.ENABLED, closeBtnId)).await(HUNDRED_MS);
      if (isCloseButtonEnabled) {
         UI.component.click(closeBtnId);
         return UI.condition.notFound(closeBtnId).await(5 * SECOND);
      }
      return false;
   }

   private static void selectInGrid(String gridId, String columnName,
         String columnValueToSelect) {
      GridControl.waitForGridToLoad(IDGroup.toIDGroup(gridId));
      GridControl.selectEntity(GridControl.findGrid(gridId),
            columnName, columnValueToSelect);
   }
}
