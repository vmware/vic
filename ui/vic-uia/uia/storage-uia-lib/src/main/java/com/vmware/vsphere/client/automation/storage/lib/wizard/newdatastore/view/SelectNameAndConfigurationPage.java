/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Represents the create new Select name and Configuration page in New Datastore wizard
 */
public class SelectNameAndConfigurationPage extends WizardNavigator {
   private static final IDGroup ID_NAME_TF = IDGroup.toIDGroup("tiwoDialog/nameInput");
   private static final IDGroup ID_NFS_FOLDER = IDGroup
         .toIDGroup("tiwoDialog/nfsFolder");
   private static final IDGroup ID_NFS_SERVER = IDGroup
         .toIDGroup("tiwoDialog/nfsServer");

   /**
    * Types datastore name in the name text field.
    *
    * @param datastoreName - datastore name to be set
    */
   public void setDatastoreName(String datastoreName) {
      UI.component.value.set(datastoreName, ID_NAME_TF);
   }

   /**
    * Types an NFS folder path
    *
    * @param datastoreSpec - datastore spec
    * @throws Exception
    */
   public void setNfsFolder(String folderPath) throws Exception {
      UI.component.value.set(folderPath, ID_NFS_FOLDER);
   }

   /**
    * Types an NFS server
    *
    * @param serverName - server name to be set for NFS datastore
    * @throws Exception
    */
   public void setNfsServer(String serverName) throws Exception {
      UI.component.value.set(serverName, ID_NFS_SERVER);
   }
}
