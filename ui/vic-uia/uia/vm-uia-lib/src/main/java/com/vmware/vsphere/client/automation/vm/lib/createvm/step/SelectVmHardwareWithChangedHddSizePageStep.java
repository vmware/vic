/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.createvm.step;

import com.vmware.vsphere.client.automation.vm.lib.createvm.view.CustomizeHardwarePage;
import com.vmware.vsphere.client.automation.vm.lib.messages.VmHardwareMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

public class SelectVmHardwareWithChangedHddSizePageStep extends CreateVmFlowStep {

   @Override
   public void execute() throws Exception {
      CustomizeHardwarePage customizePage = new CustomizeHardwarePage();

      customizePage.waitForLoadingProgressBar();
      customizePage.selectCustomizeHardwareTab(I18n.get(VmHardwareMessages.class).vmVirtualHardwareTab());
      customizePage.waitForLoadingProgressBar();
      // Set virtual disk size
      customizePage.setDiskSize(1);
      customizePage.waitForLoadingProgressBar();
      // Click on Next and verify that next page is loaded
      customizePage.gotoNextPage();
   }
}
