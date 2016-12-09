package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import java.util.List;

import com.vmware.client.automation.assertions.ContainsAssertion;
import com.vmware.client.automation.assertions.FalseAssertion;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.specs.HostSelectionSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.HostAccessibilityPage;

/**
 * Step for interacting with New datastore wizard -> Host accessibility page.
 * <br />
 * With {@link HostSelectionSpec#isEnabledExpected} set to true, the step selects
 * all the hosts from a provided list of {@link HostSelectionSpec}s.
 * <br />
 * With {@link HostSelectionSpec#isEnabledExpected} set to false, the step verifies
 * if the host is available and disabled.
 * <br />
 * NOTE: <i>Next</i> is clicked in both cases.
 */
public class SelectHostAccessibilityStep extends NewDatastoreWizardStep {

   @UsesSpec
   private List<HostSelectionSpec> hostsToSelect;

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {
      HostAccessibilityPage page = new HostAccessibilityPage();

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
      wizardNavigator.gotoNextPage();
   }
}