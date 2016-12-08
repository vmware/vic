package com.vmware.vsphere.client.automation.storage.lib.core.steps.datastore;

import java.util.List;

import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.DatastoreSpecMessages;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.ApiOperationStep;
import com.vmware.vsphere.client.automation.storage.lib.core.views.datastore.ConnectivityWithHostsView;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * ApiOperationStep implementation for performing verifications over NFS 41
 * datastore -> Manage -> Settings -> Connectivity with hosts
 */
public class VerifyNfs41DatastoreConnectivityWithHostsStep extends
      ApiOperationStep {

   @UsesSpec
   private Nfs41DatastoreSpec datastoreSpec;
   @UsesSpec
   private List<HostSpec> hostSpecs;
   private final ConnectivityWithHostsView view = new ConnectivityWithHostsView();

   private final static DatastoreSpecMessages messages = I18n
         .get(DatastoreSpecMessages.class);

   @Override
   protected CleanOperation perform() {
      for (HostSpec hostSpec : hostSpecs) {
         final String hostName = hostSpec.name.get();

         final String actualDatastoreAccessMode = view
               .getHostAccessMode(hostName);
         verifyFatal(new EqualsAssertion(actualDatastoreAccessMode,
               messages.datastoreAccessModeReadWrite(), "Access mode for host "
                     + hostName));

         final String actualDatastoreAuthenticationMode = view
               .getHostKerberosAuthentication(hostName);
         verifyFatal(new EqualsAssertion(actualDatastoreAuthenticationMode,
               datastoreSpec.authenticationMode.get().localizedDisplayName,
               "Authentication mode for host " + hostName));
      }

      // No cleanup of this verification step
      return null;
   }
}
