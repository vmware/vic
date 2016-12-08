package com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

/**
 * Localization for Datastore Spec related messages
 */
public interface DatastoreSpecMessages extends Messages {

   @DefaultMessage("NFS 3")
   String nfsVersion3();

   @DefaultMessage("NFS 4.1")
   String nfsVersion41();

   @DefaultMessage("Read-write")
   String datastoreAccessModeReadWrite();

   @DefaultMessage("Read-only")
   String datastoreAccessModeReadOnly();

   @DefaultMessage("Disabled")
   String datastoreAutheticationModeDisabled();

   @DefaultMessage("Enabled (krb5)")
   String datastoreAutheticationModeKrb();

   @DefaultMessage("Enabled (krb5i)")
   String datastoreAutheticationModeKrbi();

}
