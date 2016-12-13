/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.clone.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Select clone options for the VM in Clone VM wizard
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.clone.view.SelectCloneOptionsPage}
 */
@Deprecated
public class SelectCloneOptionsPage extends WizardNavigator {
   private static final IDGroup CUSTOMIZE_GOS =
         IDGroup.toIDGroup("selectCloneOptionsPage/customizeGOSCheckBox");
   private static final IDGroup CUSTOMIZE_HW =
         IDGroup.toIDGroup("selectCloneOptionsPage/customizeHardwareCheckBox");
   private static final IDGroup POWERON_VM =
         IDGroup.toIDGroup("selectCloneOptionsPage/powerOnCheckBox");


   public void setCustomizeHw(boolean value) {
      UI.component.value.set(value, CUSTOMIZE_HW);
  }

   public void setCustomizeGos(boolean value) {
      UI.component.value.set(value, CUSTOMIZE_GOS);
  }

   public void setPowerOnVm(boolean value) {
      UI.component.value.set(value, POWERON_VM);
  }

}
