/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * Represents the second step of the New datacenter wizard.
 */
public class NewDatastoreTypeStepView extends WizardNavigator {

   private static final String ID_VVOL_RADIO_BUTTON = "vvolRadio";

   /**
    * Select VVOL type.
    */
   public void selectVvolType() {
      UI.component.value.set("true", ID_VVOL_RADIO_BUTTON);
   }
}
