package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * WizardNavigator for Name and configuration selection for nfs 4.1 datastores
 * in the create new datastore wizard
 */
public class SelectNameAndConfigurationNfsv4Page extends WizardNavigator {

   private static final String ID_NAME_TF = "nfsSettingsPage/nameInput";
   private static final Object ID_NFS_FOLDER = "nfsSettingsPage/nfsFolder";
   private static final Object ID_NFS_SERVER = "nfsSettingsPage/nfsServer";
   private static final Object ID_ADD_NFS_SERVER = "nfsSettingsPage/addServerButton";

   /**
    * Types datastore name in the name text field.
    *
    * @param datastoreName
    *           - datastore name to be set
    */
   public void setDatastoreName(String datastoreName) {
      UI.component.value.set(datastoreName, ID_NAME_TF);
   }

   /**
    * Types an NFS folder path
    *
    * @param datastoreSpec
    *           - datastore spec
    * @throws Exception
    */
   public void setNfsFolder(String folderPath) throws Exception {
      UI.component.value.set(folderPath, ID_NFS_FOLDER);
   }

   /**
    * Types an NFS server
    *
    * @param serverName
    *           - server name to be set for NFS datastore
    * @throws Exception
    */
   public void setNfsServer(String serverName) throws Exception {
      UI.component.value.set(serverName, ID_NFS_SERVER);
      UI.component.click(ID_ADD_NFS_SERVER);
   }
}
