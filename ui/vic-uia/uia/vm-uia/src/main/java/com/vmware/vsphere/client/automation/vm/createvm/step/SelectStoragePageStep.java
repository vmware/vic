/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.vm.common.step.CreateVmFlowStep;
import com.vmware.vsphere.client.automation.vm.createvm.view.SelectStoragePage;

/**
 * The step select the and verify the storage to be used for the VM to create.
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.step.SelectStoragePageStep}
 */
@Deprecated
public class SelectStoragePageStep extends CreateVmFlowStep {
   private String _storageToSelect;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      super.prepare(filteredWorkflowSpec);
      _storageToSelect = createVmSpec.datastore.get().name.get();

      if (createVmSpec.datastoreCluster.isAssigned()) {
         DatastoreClusterSpec dsClusterSpec = createVmSpec.datastoreCluster
               .get();
         if (dsClusterSpec.sdrsEnabled.get()) {
            _storageToSelect = dsClusterSpec.name.get();
         }
      }
   }

   @Override
   public void execute() throws Exception {
      SelectStoragePage storagePage = new SelectStoragePage();
      storagePage.waitForDialogToLoad();
      boolean isSelected = storagePage.selectStorage(_storageToSelect);
      verifyFatal(isSelected, "Storage is selected!");

      // Wait for validation to complete
      new BaseView().waitForPageToRefresh();
      storagePage.waitApplySavedDataProgressBar();

      // Click on Next and verify that next page is loaded
      boolean isNextButtonClicked = storagePage.gotoNextPage();
      verifyFatal(isNextButtonClicked, "Verify the next button is clicked!");
   }
}
