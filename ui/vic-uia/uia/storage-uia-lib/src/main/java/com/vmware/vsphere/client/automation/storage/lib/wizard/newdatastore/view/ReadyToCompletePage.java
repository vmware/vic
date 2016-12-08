/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.NewDatastoreMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * Ready to complete page of the New datastore wizard (NFS flow)
 */
public class ReadyToCompletePage extends WizardNavigator {

   protected final static NewDatastoreMessages localizedMessages = I18n
         .get(NewDatastoreMessages.class);

   private static final String ACCESSIBLE_HOSTS_GRID_SELCTOR = "tiwoDialog/className=DatastoreHostAccessList";

   public String getDatastoreName() {
      return getPropertyGridValue(localizedMessages
            .readyToCompleteDatastoreNameLabel());
   }

   public String getDatastoreType() {
      return getPropertyGridValue(localizedMessages
            .readyToCompleteDatastoreTypeLabel());
   }

   public String[] getNfsServers() {
      return getPropertyGridMultiValue(
            localizedMessages.readyToCompleteDatastoreNfsServersLabel()).split(
            ",");
   }

   public String getNfsFolder() {
      return getPropertyGridValue(localizedMessages
            .readyToCompleteDatastoreNfsFolderLabel());
   }

   public String getNfsAccessMode() {
      return getPropertyGridValue(localizedMessages
            .readyToCompleteDatastoreNfsAccessModeLabel());
   }

   public String getNfsKerberosMode() {
      return getPropertyGridValue(localizedMessages
            .readyToCompleteDatastoreNfsKerberosModeLabel());
   }

   /**
    * Gets the value of the value property for the PropertyViewPropValue ui
    * component
    *
    * @param localizedProperty
    * @return
    */
   private String getPropertyGridValue(String localizedProperty) {
      // TODO this implementation does not extract the actual label value
      // displayed. It rather extracts the value of a PropertyViewPropValue UI
      // component. The HSUIA does not allow to take properties other than the
      // enumeration. Once it is supported add proper value extraction
      return UI.component.property.get(Property.VALUE,
            String.format("tiwoDialog/title=%s", localizedProperty));

   }

   /**
    * WORKAROUND! Get the value of the values property
    * PropertyViewPropMultivalue iu component
    *
    * @param localizedProperty
    * @return
    */
   private String getPropertyGridMultiValue(String localizedProperty) {
      // TODO this method is WORKARAOUND for the poor designed UI components
      // value extractions. Remove once the workaround is no longer needed
      return new UIComponent(String.format("tiwoDialog/title=%s",
            localizedProperty), BrowserUtil.flashSelenium)
            .getProperty("values");
   }

   /**
    * Get the accessible hosts displayed
    *
    * @return
    */
   public String[] getAccessibleHosts() {
      return GridControl.findGrid(ACCESSIBLE_HOSTS_GRID_SELCTOR)
            .getColumnContents(
                  localizedMessages.readyToCompleteHostGridHostNameColumn());
   }
}
