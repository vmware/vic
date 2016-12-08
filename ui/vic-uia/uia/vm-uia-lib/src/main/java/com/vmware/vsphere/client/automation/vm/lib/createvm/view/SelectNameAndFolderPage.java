/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.createvm.view;

import com.vmware.client.automation.components.control.ObjectSelectorControl3;
import com.vmware.client.automation.components.control.ObjectSelectorControl3.Tab;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Represents the Select Name and folder page of the Deploy OVF wizard
 */
public class SelectNameAndFolderPage extends WizardNavigator {
   private static final IDGroup ID_NAME_TF = IDGroup
         .toIDGroup("tiwoDialog/nameTextInput");
   private static final String ID_INV_TREE = "tiwoDialog/navTree";
   private static final String ID_TAB_BAR = "className=TabBar";

   /**
    * Types VM name name in the name text field.
    *
    * @param vmName - name of the new VM
    */
   public void setVmName(String vmName) {
      UI.condition.isTrue(new Object() {
         @Override
         public boolean equals(Object other) {
            return UI.component.exists(ID_NAME_TF)
                  && UI.component.property.getBoolean(Property.ENABLED, ID_NAME_TF);
         }
      }).await(SUITA.Environment.getPageLoadTimeout());
      UI.component.value.set(vmName, ID_NAME_TF);
   }

   /**
    * Select a VM Folder item by its name in the inventory tree
    *
    * @param spec
    *            - ManagedEntitySpec of the vApps target folder to select
    * @return
    * @return true if item is successfully selected, false otherwise
    * @throws Exception
    */
   public boolean selectParentFolder(ManagedEntitySpec locationSpec) throws Exception {
      UI.condition.isFound(ID_TAB_BAR).await(SUITA.Environment.getBackendJobMid());
      ObjectSelectorControl3.selectView(ID_TAB_BAR, Tab.BROWSE);
      return ObjectSelectorControl3.selectBrowseViewItem(ID_INV_TREE, locationSpec);
   }
}
