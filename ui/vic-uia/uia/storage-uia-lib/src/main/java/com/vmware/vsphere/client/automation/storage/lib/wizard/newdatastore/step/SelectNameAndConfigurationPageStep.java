/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import java.util.HashMap;
import java.util.Map;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.NfsDatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.NfsDatastoreSpec.NfsVersion;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.SelectNameAndConfigurationNfsv4Page;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.SelectNameAndConfigurationPage;

/**
 * Slect name and configuration for nfs server
 */
public class SelectNameAndConfigurationPageStep extends NewDatastoreWizardStep {

   /**
    * Interface for Name and configuration selection
    */
   private static interface NameAndConfigurationSelctor {

      /**
       * Performs the name and configuration selection
       *
       * @param spec
       * @throws Exception
       */
      // TODO: Change DatastoreSpec to NfsSpec once nfs spec is supported by the
      // providers
      public void selectNameAndConfiguration(DatastoreSpec spec)
            throws Exception;
   }

   /**
    * NameAndConfigurationSelctor implementation for nfs v3
    */
   private static class Nfs3NameAndConfigurationSelctor implements
         NameAndConfigurationSelctor {

      @Override
      public void selectNameAndConfiguration(DatastoreSpec spec)
            throws Exception {
         SelectNameAndConfigurationPage selectPage = new SelectNameAndConfigurationPage();
         selectPage.waitForLoadingProgressBar();
         // Set datastore name
         selectPage.setDatastoreName(spec.name.get());
         // Set NFS folder
         selectPage.setNfsFolder(spec.remotePath.get());
         // Set NFS server
         selectPage.setNfsServer(spec.remoteHost.get());
      }

   }

   /**
    * NameAndConfigurationSelctor implementation for nfs v4.1
    */
   private static class Nfs41NameAndConfigurationSelctor implements
         NameAndConfigurationSelctor {

      @Override
      public void selectNameAndConfiguration(DatastoreSpec spec)
            throws Exception {
         SelectNameAndConfigurationNfsv4Page selectPage = new SelectNameAndConfigurationNfsv4Page();
         selectPage.waitForLoadingProgressBar();
         // Set datastore name
         selectPage.setDatastoreName(spec.name.get());
         // Set NFS folder
         selectPage.setNfsFolder(spec.remotePath.get());
         // Set NFS server
         selectPage.setNfsServer(spec.remoteHost.get());
      }

   }

   private static final Map<NfsVersion, NameAndConfigurationSelctor> nameAndConfigurationSelectors;

   static {
      nameAndConfigurationSelectors = new HashMap<>();
      nameAndConfigurationSelectors.put(NfsVersion.NFS3,
            new Nfs3NameAndConfigurationSelctor());
      nameAndConfigurationSelectors.put(NfsVersion.NFS41,
            new Nfs41NameAndConfigurationSelctor());
   }

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {

      // TODO: Change the type of datastoreSpec once the NFS3 storage
      // provider is ready
      NfsVersion nfsVersion = NfsVersion.NFS3;
      if (this.datastoreSpec instanceof NfsDatastoreSpec) {
         NfsDatastoreSpec nfsSpec = (NfsDatastoreSpec) this.datastoreSpec;
         nfsVersion = nfsSpec.nfsVersion;
      }

      NameAndConfigurationSelctor nfsVersionSelctor = nameAndConfigurationSelectors
            .get(nfsVersion);

      if (nfsVersionSelctor == null) {
         throw new RuntimeException("Unknow selector for nfs version: "
               + nfsVersion);
      }

      nfsVersionSelctor.selectNameAndConfiguration(this.datastoreSpec);

      // Click on Next and verify that next page is loaded
      boolean navigatedToNextWizardPage = wizardNavigator.gotoNextPage();
      verifyFatal(TestScope.BAT, navigatedToNextWizardPage,
            "Verifying clicking Next in the  page");
   }
}
