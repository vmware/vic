/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.view;

import com.vmware.client.automation.components.control.ObjectSelectorControl3;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Class that represents the 'Select Compute Resource' page of the New VM
 * wizard.
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.view.SelectComputeResourcePage}
 */
@Deprecated
public class SelectComputeResourcePage extends WizardNavigator {

   // Component IDs
   private static final String ID_TIWO_DIALOG = "tiwoDialog/";
   private static final String ID_INV_TREE = ID_TIWO_DIALOG
         + "selectResourcePage/navTree";

   /**
    * Select resource by its name in the Inventory tree.
    *
    * @param spec
    *           - ManagedEntitySpec of the compute resource to select (host,
    *           cluster, resource pool, vApp)
    * @return true if item is successfully selected, false otherwise
    * @throws Exception
    *            if ActionScript error occurs
    */
   public boolean selectResource(ManagedEntitySpec spec) throws Exception {
      return ObjectSelectorControl3.selectBrowseViewItem(ID_INV_TREE, spec);
   }
}
