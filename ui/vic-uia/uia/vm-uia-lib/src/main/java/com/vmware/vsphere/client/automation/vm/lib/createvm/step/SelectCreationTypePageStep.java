/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.createvm.step;

import com.vmware.client.automation.workflow.View;
import com.vmware.vsphere.client.automation.vm.lib.createvm.view.SelectCreationTypePage;

@View(SelectCreationTypePage.class)
public class SelectCreationTypePageStep extends CreateVmFlowStep {

   @Override
   public void execute() throws Exception {
      SelectCreationTypePage selectPage = new SelectCreationTypePage();
      selectPage.waitForLoadingProgressBar();
      selectPage.selectCreationType(createVmSpec.creationType.get());
      selectPage.gotoNextPage();
   }
}
