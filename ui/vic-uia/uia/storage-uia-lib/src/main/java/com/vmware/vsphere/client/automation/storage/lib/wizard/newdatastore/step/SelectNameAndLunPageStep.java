/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.SelectNameAndLunPage;

public class SelectNameAndLunPageStep extends NewDatastoreWizardStep {

   @Override
   public void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {
      SelectNameAndLunPage selectNameAndLunPage = new SelectNameAndLunPage();
      wizardNavigator.waitForLoadingProgressBar();
      // Set datastore name
      selectNameAndLunPage.setDatastoreName(datastoreSpec.name.get());
      // Set NFS folder
      selectNameAndLunPage.selectDevice(0);
      // Click on Next and verify that next page is loaded
      boolean navigatedToNextWizardPage = wizardNavigator.gotoNextPage();
      verifyFatal(TestScope.BAT, navigatedToNextWizardPage,
            "Verifying clicking Next in the  page");
   }
}
