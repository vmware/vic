/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;

/**
 * Select creation type for the datastore
 *
 */
public class SelectDatastoreTypePage extends WizardNavigator {
   private static final IDGroup NFS_RADIO = IDGroup.toIDGroup("nfsRadio");
   private static final IDGroup VMFS_RADIO = IDGroup.toIDGroup("vmfsRadio");

   /**
    * Selects datastore type
    *
    * @param dsType - datastore type to be selected
    */
   public void selectCreationType(DatastoreType dsType) {
      IDGroup dsTypeRadioId = getComponentId(dsType);
      if (dsTypeRadioId != null) {
         UI.component.value.set(true, dsTypeRadioId);
      }
   }

   private IDGroup getComponentId(DatastoreType dsType) {
      switch (dsType) {
         case NFS:
            return NFS_RADIO;
         case VMFS:
            return VMFS_RADIO;
         default:
            throw new AssertionError(String.format(
                  "Invalid Datastore Type passed %s",
                  dsType.toString()));
      }
   }
}
