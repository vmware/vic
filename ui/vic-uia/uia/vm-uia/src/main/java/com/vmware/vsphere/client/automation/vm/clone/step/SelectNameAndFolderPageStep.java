/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.clone.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.vsphere.client.automation.vm.clone.CloneVmFlowStep;
import com.vmware.vsphere.client.automation.vm.createvm.view.SelectNameAndFolderPage;

/**
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.clone.step.SelectNameAndFolderPageStep}
 */
@Deprecated
public class SelectNameAndFolderPageStep extends CloneVmFlowStep {

   @Override
   public void execute() throws Exception {
      SelectNameAndFolderPage selectPage = new SelectNameAndFolderPage();

      selectPage.waitForDialogToLoad();
      // Set VM name
      selectPage.setVmName(cloneVmSpec.name.get());
      selectPage.waitApplySavedDataProgressBar();
      new BaseView().waitForPageToRefresh();
      // Click on Next and verify that next page is loaded
      boolean isNextButtonClicked = selectPage.gotoNextPage();
      verifyFatal(isNextButtonClicked, "Verify the next button is clicked!");
   }
}
