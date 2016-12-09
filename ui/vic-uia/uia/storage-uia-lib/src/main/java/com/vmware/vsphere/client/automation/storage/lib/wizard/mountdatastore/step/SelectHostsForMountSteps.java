package com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.step;

import java.util.List;

import com.vmware.client.automation.assertions.ContainsAssertion;
import com.vmware.client.automation.assertions.FalseAssertion;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.specs.HostSelectionSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.SinglePageWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.view.SelectHostsPage;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;

/**
 * {@link NewDatastoreWizardStep} implementation for executing a selection on
 * Mount Datastore to Additional hosts
 */
public class SelectHostsForMountSteps extends SinglePageWizardStep {

   @UsesSpec
   private List<HostSelectionSpec> hostsToSelect;

   @Override
   protected void executeWizardOperation(
         SinglePageDialogNavigator wizardNavigator) throws Exception {
      SelectHostsPage page = new SelectHostsPage();

      for (HostSelectionSpec hostSpec : hostsToSelect) {
         if (!hostSpec.isEnabledExpected.isAssigned()
               || hostSpec.isEnabledExpected.get()) {
            page.selectHost(hostSpec.hostSpec.get().name.get());
         } else {
            String hostDisplayName = hostSpec.hostSpec.get().name.get()
                  + hostSpec.hostDisplayNameSuffix.get();

            List<String> hostDisplayNames = page.getHostDisplayNames();
            verifySafely(new ContainsAssertion(hostDisplayNames,
                  hostDisplayName, "Host displayed as " + hostDisplayName
                        + "is available"));
            verifySafely(new FalseAssertion(
                  page.isHostEnabled(hostDisplayName), String.format(
                        "%s is enabled", hostDisplayName)));
         }
      }

   }

}
