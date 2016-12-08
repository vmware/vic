/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Represents the rename page dialog.
 */
public class RenameEntityPage extends SinglePageDialogNavigator {
   private static final IDGroup ID_NAME_TF = IDGroup.toIDGroup("inputLabel");

   /**
    * Gets the name value if the name text input.
    * @param name  name to set
    * @return
    */
   public String getName() {
      return UI.component.value.get(ID_NAME_TF);
   }

   /**
    * Set the name parameter value to the name text input.
    * @param name  name to set
    */
   public void setName(String name) {
      UI.component.value.set(name, ID_NAME_TF);
   }
}
