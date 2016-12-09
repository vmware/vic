package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import java.util.List;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.storage.lib.core.grid.CheckBoxGrid;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.NewDatastoreMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * WizardNavigator implementation for new datastore wizard -> Host accessibility
 * page
 */
public class HostAccessibilityPage extends WizardNavigator {

   protected final static NewDatastoreMessages localizedMessages = I18n
         .get(NewDatastoreMessages.class);

   private static final String HOST_GRID_SELECTOR = "nfsSelectHostsPage/hostList";
   private static final String CHECKBOX_PROPERTY = "data.name";

   /**
    * Select a host from the hosts grid
    *
    * @param hostName
    */
   public void selectHost(String hostName) {
      new CheckBoxGrid(HOST_GRID_SELECTOR, CHECKBOX_PROPERTY).select(hostName);
   }

   public List<String> getHostDisplayNames() {
      return GridControl
            .getColumnContents(GridControl.findGrid(HOST_GRID_SELECTOR),
                  localizedMessages
                        .selectHostAccessibilityPageHostGridColumnHeader());
   }

   public boolean isHostEnabled(String hostDisplayeName) {
      return new CheckBoxGrid(HOST_GRID_SELECTOR, CHECKBOX_PROPERTY)
            .isEnabled(hostDisplayeName);
   }

}
