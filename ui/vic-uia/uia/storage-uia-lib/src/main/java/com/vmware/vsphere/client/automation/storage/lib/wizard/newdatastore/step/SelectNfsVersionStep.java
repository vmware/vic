package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.NfsDatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.WizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.SelectNfsVersionPage;

/**
 * WizardStep implementation for selecting the Nfs version of the datastore
 */
public class SelectNfsVersionStep extends WizardStep {

   @UsesSpec
   private NfsDatastoreSpec datastoreSpec;

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {

      SelectNfsVersionPage nfsVersionPage = new SelectNfsVersionPage();
      nfsVersionPage.selectNfsVersion(this.datastoreSpec.nfsVersion);

      wizardNavigator.gotoNextPage();
   }

}
