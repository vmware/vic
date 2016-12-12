package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.assertions.Assertion;
import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.ReadyToCompletePage;

/**
 * NewDatastoreWizardStep implementation for verifications of a DatastoreSpec on
 * the Ready to complete page
 */
public class ReadyToCompleteVerificationStep extends NewDatastoreWizardStep {

   /**
    * Interface for providing with assertions based on a {@link DatastoreSpec}
    */
   private static interface IDatastorePropertyAssetionProvider {

      /**
       * Get assertion based on a datastore spec
       *
       * @param datastoreSpec
       *           the datastore to be asserted
       * @return
       */
      List<Assertion> getAssertions(DatastoreSpec datastoreSpec);
   }

   /**
    * {@link IDatastorePropertyAssetionProvider} implementation for NFS 4.1
    * datastore
    */
   private static final class Nfs41AssertionProvider implements
         IDatastorePropertyAssetionProvider {

      @Override
      public List<Assertion> getAssertions(DatastoreSpec datastoreSpec) {
         if (!(datastoreSpec instanceof Nfs41DatastoreSpec)) {
            throw new RuntimeException(
                  "Can not perform verification for NFS 4.1 Datastore with spec of "
                        + datastoreSpec.getClass().getCanonicalName());
         }

         Nfs41DatastoreSpec nfsDatastore = (Nfs41DatastoreSpec) datastoreSpec;
         List<Assertion> result = new ArrayList<Assertion>();

         ReadyToCompletePage page = new ReadyToCompletePage();

         String actualDatastoreName = page.getDatastoreName();
         result.add(new EqualsAssertion(actualDatastoreName, nfsDatastore.name
               .get(), "Ready to complete datastore name"));

         String actualDatastoreType = page.getDatastoreType();
         result.add(new EqualsAssertion(actualDatastoreType,
               nfsDatastore.nfsVersion.localizedNfsVersionString,
               "Ready to complete datastore type"));

         String actualNfsServers = page.getNfsServers()[0];
         result.add(new EqualsAssertion(actualNfsServers,
               nfsDatastore.remoteHost.get(), "Ready to NFS remote path"));

         String actualNfsFolder = page.getNfsFolder();
         result.add(new EqualsAssertion(actualNfsFolder,
               nfsDatastore.remotePath.get(), "Ready to complete Nfs folder"));

         String actualNfsAccessMode = page.getNfsAccessMode();
         result.add(new EqualsAssertion(actualNfsAccessMode,
               nfsDatastore.accessMode.get().localizedDisplayName,
               "Ready to complete Nfs access mode"));

         String actualKerberosMode = page.getNfsKerberosMode();
         result.add(new EqualsAssertion(actualKerberosMode,
               nfsDatastore.authenticationMode.get().localizedDisplayName,
               "Ready to complete Nfs kerberos mode"));

         return result;
      }

   }

   /**
    * Enumeration of all supported verification mode for
    * {@link ReadyToCompleteVerificationStep}
    */
   public static enum ReadyToCompleteVerificationMode {
      /**
       * Verifies the NFS 4.1 Datastore
       */
      NFS41(new Nfs41AssertionProvider());

      /**
       * Assertion provide for current enumeration value
       */
      final IDatastorePropertyAssetionProvider assertionProvider;

      private ReadyToCompleteVerificationMode(
            IDatastorePropertyAssetionProvider assertionProvider) {
         this.assertionProvider = assertionProvider;
      }
   }

   private final ReadyToCompleteVerificationMode verificationMode;

   public ReadyToCompleteVerificationStep(
         ReadyToCompleteVerificationMode verificationMode) {
      this.verificationMode = verificationMode;
   }

   @UsesSpec
   private DatastoreSpec datastoreSpec;

   @UsesSpec
   private List<HostSpec> hostSpecs;

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {
      for (Assertion assertion : this.verificationMode.assertionProvider
            .getAssertions(datastoreSpec)) {
         verifySafely(assertion);
      }
   }

}
