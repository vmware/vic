package com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

/**
 * Messages associated with the new datastore wizard
 */
public interface MountMessages extends Messages {

   @DefaultMessage("Mount Datastore to Additional Hosts...")
   String mountDatastoreMenuOption();

   @DefaultMessage("Host")
   String hostSelectionGridHostDisplayName();

}