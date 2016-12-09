/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.SelectDatastoreTypePage;

public class SelectDatastoreTypePageStep extends NewDatastoreWizardStep {

   @Override
   public void executeWizardOperation(WizardNavigator wizardNavigator) {
      new SelectDatastoreTypePage()
            .selectCreationType(datastoreSpec.type.get());

      wizardNavigator.gotoNextPage();
   }
}
