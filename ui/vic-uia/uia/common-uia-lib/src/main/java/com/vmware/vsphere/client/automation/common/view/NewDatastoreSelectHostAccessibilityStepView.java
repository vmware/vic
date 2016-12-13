/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * Represents the forth step of the New datacenter wizard.
 */
public class NewDatastoreSelectHostAccessibilityStepView extends WizardNavigator {

   private static final String ID_GRID = "tiwoDialog/vvolSelectHostsPage/hostList";

   /**
    * Select given host in host list.
    *
    * @param hostName   the name of the host to be selected
    * @return           true if the selection was successful, false otherwise
    */
   public boolean selectHostByName(String hostName) {
      return GridControl.checkEntities(ID_GRID, hostName);
   }
}
