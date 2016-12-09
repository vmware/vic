package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

/**
 * Messages associated with the new datastore wizard
 */
public interface NewDatastoreMessages extends Messages {

   @DefaultMessage("Enable Kerberos-based authentication")
   String enableDisableKrbAuthentication();

   @DefaultMessage("Use Kerberos for authentication only (krb5)")
   String useKrb5();

   @DefaultMessage("Use Kerberos for authentication and data integrity (krb5i)")
   String useKrb5i();

   @DefaultMessage("Host")
   String selectHostAccessibilityPageHostGridColumnHeader();

   @DefaultMessage("Host")
   String readyToCompleteHostGridHostNameColumn();

   @DefaultMessage("Name:")
   String readyToCompleteDatastoreNameLabel();

   @DefaultMessage("Type:")
   String readyToCompleteDatastoreTypeLabel();

   @DefaultMessage("NFS servers:")
   String readyToCompleteDatastoreNfsServersLabel();

   @DefaultMessage("Folder:")
   String readyToCompleteDatastoreNfsFolderLabel();

   @DefaultMessage("Access Mode:")
   String readyToCompleteDatastoreNfsAccessModeLabel();

   @DefaultMessage("Kerberos:")
   String readyToCompleteDatastoreNfsKerberosModeLabel();
}
