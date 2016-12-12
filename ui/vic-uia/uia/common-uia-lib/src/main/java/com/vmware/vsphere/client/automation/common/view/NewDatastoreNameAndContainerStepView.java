/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * Represents the third step of the New datacenter wizard.
 */
public class NewDatastoreNameAndContainerStepView extends WizardNavigator {

   private static final String ID_DATASTORE_NAME = "tiwoDialog/nameInput";

   /**
    * Types datastore name.
    *
    * @param datastoreName     datastore name that will be typed into the text field
    */
   public void setDatastoreName(String datastoreName) {
      _logger.info(String.format("Entering datastore name '%s'", datastoreName));
      UI.component.value.set(datastoreName, ID_DATASTORE_NAME);
   }
}
