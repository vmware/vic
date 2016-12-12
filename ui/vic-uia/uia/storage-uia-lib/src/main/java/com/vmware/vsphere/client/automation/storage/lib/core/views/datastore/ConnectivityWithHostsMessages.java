package com.vmware.vsphere.client.automation.storage.lib.core.views.datastore;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

/**
 * Messages for datastore -> Manage -> Settings -> Connectivity with hosts
 */
public interface ConnectivityWithHostsMessages extends Messages {

   @DefaultMessage("Host")
   String connectivityWithHostsGridHostDisplayNameColumn();

   @DefaultMessage("Access Mode")
   String connectivityWithHostsGridHostAccessModeColumn();

   @DefaultMessage("Kerberos Authentication")
   String connectivityWithHostsGridHostKerberosAuthenticationColumn();
}
