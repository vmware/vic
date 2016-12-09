/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.view;

import java.text.MessageFormat;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec.HddNodes;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.flexui.componentframework.controls.spark.SparkComboBox;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * Customize Hardware page
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.view.CustomizeHardwarePage}
 */
@Deprecated
public class CustomizeHardwarePage extends WizardNavigator {
   private static final String DISK_INDEX = "<DiskIndex>";
   private static final String ID_DEFAULT_DISK = "tiwoDialog/customizeHardwarePage/vmConfigControl/disk_"
         + DISK_INDEX + "/diskSize/capacity/textDisplay";

   private static final String ID_DEFAULT_HDD_NODE = "tiwoDialog/customizeHardwarePage/vmConfigControl/disk_"
         + DISK_INDEX + "controllerList";

   private static final String CUSTOMIZE_HARDWARE_TAB_ID_FORMAT = "vmConfigPages/label={0}";
   private static final String BOOT_OPTIONS_LABEL_ID = "BIOSPage/headerPropertyGrid/className=Label";
   private static final String BOOT_OPTIONS_SECTION_ARROW_ID = "arrowImageBIOSPage";
   private static final String UPGRADE_SECTION_ARROW_ID = "className=VmUpgradePage/arrowImagenull";
   private static final String FIRMWARE_COMBO_ID = "firmwareCombo";
   private static final String COMPATIBILITY_VERSION_COMBO_ID = "version";
   private static final String SCHEDULE_UPGRADE_CHECKBOX_ID = "scheduled";
   private static final String SECURITY_BOOT_CHECKBOX_ID = "secureBoot";

   private static final String ID_ADD_DEVICE_MENU = "menu";
   public static final String ID_ADD_NEW_DEVICE_BUTTON = "addDeviceButton";
   public static final String ID_ADD_DEVICE_COMBO = "addHardware/popUpButton";
   // Shared PCI Device
   private static final String ID_NEW_SHARED_PCI_LABEL = "automationName=Shared PCI Device";
   private static final String ID_NEW_SHARED_PCI_MENU_ITEM = ID_ADD_DEVICE_MENU
         + "/" + ID_NEW_SHARED_PCI_LABEL;
   private static final String ID_RESERVE_ALL_MEMORY_BTN = "_PciSharedPage_Button1";
   private static final String ID_VGPU_WARNING_LABEL = "_PciSharedPage_Text1";
   private static final String ID_VGPU_NOTE_LABEL = "_PciSharedPage_Text2";
   private static final String ID_VGPU_PROFILES_COMBOBOX = "vgpu";
   private static final String ID_MAX_DEVICES_MSG = "message";

   public static void openHwDevicesMenu() {
      // Click to open Hardware devices menu
      UI.component.click(ID_ADD_DEVICE_COMBO);

      // Wait until the menu is loaded
      UI.condition.isFound(ID_ADD_DEVICE_MENU)
            .await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Click Add button in VM Summary>Edit Settings>New Device
    */
   public static void clickAddDevice() {
      // Click "Add" button to confirm the choice
      UI.component.click(ID_ADD_NEW_DEVICE_BUTTON);
   }

   /**
    * Sets the size of the virtual disk
    *
    * @param hddIndex
    *           - virtual disk index (0 is the first default disk)
    * @param sizeInGb
    *           - size in GB
    * @throws Exception
    */
   public void setDiskSize(final int hddIndex, final int sizeInGb)
         throws Exception {
      final IDGroup disk = IDGroup.toIDGroup(
            ID_DEFAULT_DISK.replace(DISK_INDEX, Integer.toString(hddIndex)));
      waitForControlEnabled(IDGroup.toIDGroup(disk));
      UI.component.value.set(sizeInGb, disk);
   }

   /**
    * Sets default virtual disk size
    *
    * @param sizeInGb
    *           - size in GB
    * @throws Exception
    */
   public void setDiskSize(int sizeInGb) throws Exception {
      setDiskSize(0, sizeInGb);
   }

   /**
    * Set disk node for hard disk
    *
    * @param hddIndex
    * @param hddNode
    */
   public void setHddNode(String hddIndex, HddNodes hddNode) {
      final IDGroup disk = IDGroup
            .toIDGroup(ID_DEFAULT_HDD_NODE.replace(DISK_INDEX, hddIndex));
      waitForControlEnabled(IDGroup.toIDGroup(disk));
      UI.component.selectByIndex(hddNode.ordinal(), ID_DEFAULT_DISK);
   }

   /**
    * Select the Customize Hardware tab
    *
    * @param tabLabel
    *           Label of the tab
    */
   public void selectCustomizeHardwareTab(String tabLabel) {
      final IDGroup tab = IDGroup.toIDGroup(
            MessageFormat.format(CUSTOMIZE_HARDWARE_TAB_ID_FORMAT, tabLabel));
      waitForControlEnabled(tab);
      UI.component.click(tab);
   }

   /**
    * Expands the Boot options section
    */
   public void expandUpgradeSection() {
      final IDGroup arrow = IDGroup.toIDGroup(UPGRADE_SECTION_ARROW_ID);
      waitForControlEnabled(arrow);
      UI.component.click(arrow);
   }

   /**
    * Set the Scheduled upgrade
    *
    * @param toSchedule
    */
   public void setScheduleUpgrade(Boolean toSchedule) {
      final IDGroup checkBox = IDGroup.toIDGroup(SCHEDULE_UPGRADE_CHECKBOX_ID);
      waitForControlEnabled(checkBox);
      UI.component.value.set(toSchedule, SCHEDULE_UPGRADE_CHECKBOX_ID);
   }

   /**
    * Set the Compatibility version
    *
    * @param compatibilityVersion
    *           Label of Compatibility version
    */
   public void setCompatibilityVersion(String compatibilityVersion) {
      SparkComboBox combo = new SparkComboBox(COMPATIBILITY_VERSION_COMBO_ID,
            BrowserUtil.flashSelenium);
      combo.waitForElementVisibleOnPath(
            UiDelay.UI_OPERATION_TIMEOUT.getDuration());
      combo.selectItemByValue(compatibilityVersion);
   }

   /**
    * Expands the Boot options section
    */
   public void expandBootOptionsSection() {
      boolean isCollapsed = UI.component.property.getBoolean(Property.VISIBLE,
            BOOT_OPTIONS_LABEL_ID);
      if (isCollapsed) {
         final IDGroup arrow = IDGroup.toIDGroup(BOOT_OPTIONS_SECTION_ARROW_ID);
         waitForControlEnabled(arrow);
         UI.component.click(arrow);
      }
   }

   /**
    * Set the Firmware boot option
    *
    * @param firmwareName
    *           Label of firmware
    */
   public void setFirmwareBootOption(String firmwareName) {
      SparkComboBox combo = new SparkComboBox(FIRMWARE_COMBO_ID,
            BrowserUtil.flashSelenium);
      combo.waitForElementVisibleOnPath(
            UiDelay.UI_OPERATION_TIMEOUT.getDuration());
      combo.selectItemByValue(firmwareName);
   }

   /**
    * Set the Security boot option
    *
    * @param toEnableSecurityBoot
    */
   public void setSecurityBootOption(Boolean toEnableSecurityBoot) {
      final IDGroup checkBox = IDGroup.toIDGroup(SECURITY_BOOT_CHECKBOX_ID);
      waitForControlEnabled(checkBox);
      UI.component.value.set(toEnableSecurityBoot, SECURITY_BOOT_CHECKBOX_ID);
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
    * Get text of max devices message
    *
    * @return
    */
   public String getMaxDevicesMsg() {
      return UI.component.property.get(Property.TEXT,
            IDGroup.toIDGroup(ID_MAX_DEVICES_MSG));
   }
}
