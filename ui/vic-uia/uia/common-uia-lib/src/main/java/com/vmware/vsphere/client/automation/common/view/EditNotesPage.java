/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.suitaf.apl.IDGroup;

/**
 * View that specifies the VDC Notes Portlet
 */
public class EditNotesPage extends SinglePageDialogNavigator {
   private static final IDGroup ID_NOTES_PORTLET_EDIT_TB = IDGroup.toIDGroup("textDisplay");

   /**
    * Set description in Notes dialog.
    *
    * @param descripton - the VDC description to be set
    */
   public void setNotesDescription(String description) {
      UI.component.value.set(description, ID_NOTES_PORTLET_EDIT_TB);
   }
}
