/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.NfsDatastoreSpec.NfsVersion;

/**
 * Select NFS version page part of New datastore wizard
 */
public class SelectNfsVersionPage extends WizardNavigator {

   private static final String NFS_VERSION_PAGE = "nfsVersionPage";
   private static final String NFS_VERSION_SELCTOR = NFS_VERSION_PAGE + "/label=%s";

   /**
    * Select the nfs version on the page
    *
    * @param nfsVersion
    *           the desired nfs version
    */
   public void selectNfsVersion(NfsVersion nfsVersion) {
      UI.component.click(String.format(NFS_VERSION_SELCTOR,
            nfsVersion.localizedNfsVersionString));
   }

   /**
    * Checks whether the NFS version page is available
    *
    * @return true if the NFS version page is available
    */
   public boolean isVersionPageAvailable() {
      return UI.component.exists(NFS_VERSION_PAGE);
   }
}
