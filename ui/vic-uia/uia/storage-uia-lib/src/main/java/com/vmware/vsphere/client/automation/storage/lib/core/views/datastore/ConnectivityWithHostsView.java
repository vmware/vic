package com.vmware.vsphere.client.automation.storage.lib.core.views.datastore;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * BaseView implementation for NFS 41 datastore -> Manage -> Settings ->
 * Connectivity with hosts
 *
 */
public class ConnectivityWithHostsView extends BaseView {

   private static final String CONNECTIVITY_WITH_HOSTS_GRID_SELECTOR = "vsphere.core.datastore.manage.settings.nfsVVolDsConnectivityView/hostList";
   private static final ConnectivityWithHostsMessages messages = I18n
         .get(ConnectivityWithHostsMessages.class);

   /**
    * Get the value displayed in the Access Mode column for Connectivity with
    * Hosts grid
    *
    * @param hostDisplayName
    *           the host display name
    * @return
    */
   public String getHostAccessMode(String hostDisplayName) {
      int rowIndex = getHostRowIndex(hostDisplayName);

      return GridControl.findGrid(CONNECTIVITY_WITH_HOSTS_GRID_SELECTOR)
            .getCellValue(rowIndex,
                  messages.connectivityWithHostsGridHostAccessModeColumn());
   }

   /**
    * Get the value displayed in the Kerberos Authentication column for
    * Connectivity with Hosts grid
    *
    * @param hostDisplayName
    *           the host display name
    * @return
    */
   public String getHostKerberosAuthentication(String hostDisplayName) {
      int rowIndex = getHostRowIndex(hostDisplayName);

      return GridControl
            .findGrid(CONNECTIVITY_WITH_HOSTS_GRID_SELECTOR)
            .getCellValue(
                  rowIndex,
                  messages
                        .connectivityWithHostsGridHostKerberosAuthenticationColumn());
   }

   /**
    * Find the row index based on a host diplay name
    *
    * @param hostDisplayName
    *           the host display name
    * @return
    */
   private int getHostRowIndex(String hostDisplayName) {
      AdvancedDataGrid grid = GridControl
            .findGrid(CONNECTIVITY_WITH_HOSTS_GRID_SELECTOR);
      int result = GridControl.getEntityIndex(grid,
            messages.connectivityWithHostsGridHostDisplayNameColumn(),
            hostDisplayName);

      if (result < 0) {
         throw new RuntimeException(String.format(
               "Can not find row '%s' in column '%s' in grid %s",
               hostDisplayName,
               messages.connectivityWithHostsGridHostDisplayNameColumn(),
               CONNECTIVITY_WITH_HOSTS_GRID_SELECTOR));
      }

      return result;
   }
}
