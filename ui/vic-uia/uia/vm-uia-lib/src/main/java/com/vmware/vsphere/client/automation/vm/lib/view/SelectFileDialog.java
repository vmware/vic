/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.BaseDialogNavigator;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.srv.common.messages.VmMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * A dialog of this kind can be opened from VM -> Edit Settings -> Add Existing Hard Disk
 */
public class SelectFileDialog extends BaseDialogNavigator {
   private static final VmMessages messages = I18n.get(VmMessages.class);
   private static final String DIALOG_OK_BUTTON_ID = "buttonOk";
   private static final String FILE_LIST_ID = "fileList";
   private static final int HUNDRED_MS = 100;
   private static final int SECOND = 1000;

   /**
    * Selects the disk identified by the given datastore, path and name.
    *
    * @param datastoreName
    *           the name of the containing datastore
    * @param pathToDiskFolder
    *           location of the disk within the datastore
    * @param diskName
    *           the name of the disk
    * @return true if selection was successful, false otherwise
    */
   public static boolean selectExistingDisk(String datastoreName,
         String pathToDiskFolder, String diskName) {
      UI.component.click("dbTreeView/text=" + datastoreName);
      String[] folders = pathToDiskFolder.contains("/") ? pathToDiskFolder
            .split("\\/") : new String[] { pathToDiskFolder };
      for (String folder : folders) {
         selectInGrid(FILE_LIST_ID, messages.getContentsColumnHeader(), folder);
      }
      selectInGrid(FILE_LIST_ID, messages.getContentsColumnHeader(), diskName);
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
