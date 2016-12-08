/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.step;

import com.vmware.vsphere.client.automation.vm.createvm.view.SelectNameAndFolderPage;
import com.vmware.vsphere.client.automation.vm.migrate.MigrateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.migrate.view.SelectComputeResourcePage;

public class SelectComputeResourcePageStep extends MigrateVmFlowStep {

   @Override
   public void execute() throws Exception {

      SelectComputeResourcePage selectComputeResourcePage =
            new SelectComputeResourcePage();
      SelectNameAndFolderPage selectNameAndFolderPage = new SelectNameAndFolderPage();

      selectComputeResourcePage.waitForDialogToLoad();
      // Select target resource
      selectComputeResourcePage.selectEntity(migrateVmSpec);
      selectNameAndFolderPage.waitApplySavedDataProgressBar();
      // Click on Next and verify that next page is loaded
      selectNameAndFolderPage.gotoNextPage();
   }
}
