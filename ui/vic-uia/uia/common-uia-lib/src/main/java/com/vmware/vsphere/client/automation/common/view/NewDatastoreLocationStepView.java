/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.control.ObjectSelectorControl3;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.flexui.componentframework.controls.mx.custom.InventoryTree;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Represents the first step of the New datacenter wizard.
 */
public class NewDatastoreLocationStepView extends WizardNavigator {

   private static final String ID_INV_TREE = "tiwoDialog/datastoreLocationPage/navTree";

   /**
    * Select the location where the datastore should be placed.
    *
    * @param spec          managed entity spec that represents the location object
    * @return              if the operation succeeded
    * @throws Exception
    */
   public boolean selectLocationByName(ManagedEntitySpec spec) throws Exception {

      // wait a short period of time as the tree seems not to be ready
      // right after the page is loaded
      Thread.sleep(3000);
      InventoryTree inventoryTree = new InventoryTree(ID_INV_TREE, BrowserUtil.flashSelenium);
      inventoryTree.expandNode("0");
      return ObjectSelectorControl3.selectBrowseViewItem(ID_INV_TREE, spec);
   }
}
