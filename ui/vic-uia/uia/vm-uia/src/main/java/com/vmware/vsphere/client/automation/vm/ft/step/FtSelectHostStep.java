/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.vm.ft.view.FtSelectHostPage;

/**
 * Select a host on the Hosts page of the Fault Tolerance wizard
 */
public class FtSelectHostStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private HostSpec _hostSpec;

   @Override
   public void execute() throws Exception {
      FtSelectHostPage selectHostPage = new FtSelectHostPage();

      selectHostPage.waitForDialogToLoad();
      // Verify the host has been selected
      String hostToSelect = _hostSpec.name.get();
      verifyFatal(selectHostPage.selectHost(hostToSelect),
            "Verify the host has been selected.");

      // Wait for validation to complete
      selectHostPage.waitApplySavedDataProgressBar();

      // Click Next
      verifyFatal(selectHostPage.gotoNextPage(),
            "Verify navigation to the next page is successful.");
   }
}