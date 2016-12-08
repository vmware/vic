package com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.view;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.vsphere.client.automation.storage.lib.core.grid.CheckBoxGrid;
import com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.MountMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

public class SelectHostsPage {
   private final static MountMessages localizedMessages = I18n
         .get(MountMessages.class);

   private static final String HOST_GRID_SELECTOR = "tiwoDialog/hostList";
   private static final String CHECKBOX_PROPERTY = "data.name";

   /**
    * Select a host from the hosts grid
    *
    * @param hostName
    */
   public void selectHost(String hostName) {
      new CheckBoxGrid(HOST_GRID_SELECTOR, CHECKBOX_PROPERTY).select(hostName);
   }

   /**
    * Returns the host names from the grid in the page (with error/warning and state information).
    *
    * @return list of host names
    * @see {@link #getHostNames()}
    */
   public List<String> getHostDisplayNames() {
      List<String> hostDisplayNames = GridControl.getColumnContents(
            GridControl.findGrid(HOST_GRID_SELECTOR),
            localizedMessages.hostSelectionGridHostDisplayName());
      return hostDisplayNames;
   }

   /**
    * Returns the host names from the grid in the page (without error/warning and state information).
    *
    * @return list of host names
    * @see {@link #getHostDisplayNames()}
    */
   public List<String> getHostNames() {
      List<String> hostDisplayNames = getHostDisplayNames();

      List<String> result = new ArrayList<String>();
      for (String displayString : hostDisplayNames) {
         // Trim the host display name suffix. The hosts are displayed as:
         // hostname.domain (some additional data for the host e.g: host is
         // disabled, because..)
         result.add(displayString.replaceAll("[\\s]*\\(.*\\)", ""));
      }

      return result;
   }

   public boolean isHostEnabled(String hostDisplayeName) {
      return new CheckBoxGrid(HOST_GRID_SELECTOR, CHECKBOX_PROPERTY)
            .isEnabled(hostDisplayeName);
   }
}
