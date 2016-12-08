package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.WizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.ConfigureKerberosAuthenticationPage;

/**
 * Step for interacting with New datastore wizard -> Configure Kerberos
 * authentication
 */
public class ConfigureKerberosAuthenticationStep extends WizardStep {

   @UsesSpec
   private Nfs41DatastoreSpec datastoreSpec;

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {
      ConfigureKerberosAuthenticationPage page = new ConfigureKerberosAuthenticationPage();
      page.setAuthenticationMode(this.datastoreSpec.authenticationMode.get());

      wizardNavigator.gotoNextPage();
   }
}
