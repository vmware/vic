/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.step;

import com.vmware.client.automation.workflow.View;
import com.vmware.vsphere.client.automation.vm.common.step.CreateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.createvm.view.SelectCreationTypePage;

/**
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectCreationTypePageStep}}
 */
@View(SelectCreationTypePage.class)
@Deprecated
public class SelectCreationTypePageStep extends CreateVmFlowStep {

   @Override
   public void execute() throws Exception {
      SelectCreationTypePage selectPage = new SelectCreationTypePage();
      selectPage.waitForLoadingProgressBar();
      selectPage.selectCreationType(createVmSpec.creationType.get());
      selectPage.gotoNextPage();
   }
}
