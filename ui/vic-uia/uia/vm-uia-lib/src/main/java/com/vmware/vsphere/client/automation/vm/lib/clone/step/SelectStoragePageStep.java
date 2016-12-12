/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.lib.clone.step;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.view.SelectStoragePage;

/**
 * The step select the and verify the storage to be used for the VM to clone.
 */
public class SelectStoragePageStep extends CloneVmFlowStep {
   private String _storageToSelect;
   private String _datastoreInDsClusterToSelect;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      super.prepare(filteredWorkflowSpec);
      _storageToSelect = cloneVmSpec.datastore.get().name.get();

      if (cloneVmSpec.targetDatastoreCluster.isAssigned()) {
         DatastoreClusterSpec dsClusterSpec = cloneVmSpec.targetDatastoreCluster
               .get();
         _storageToSelect = dsClusterSpec.name.get();

         if (!dsClusterSpec.sdrsEnabled.get()) {
            _datastoreInDsClusterToSelect = cloneVmSpec.datastore.get().name
                  .get();
         }
      }
   }

   @Override
   public void execute() throws Exception {
      SelectStoragePage storagePage = new SelectStoragePage();
      storagePage.waitForDialogToLoad();
      boolean isSelected = storagePage.selectStorage(_storageToSelect);
      verifyFatal(isSelected, "Storage is selected!");

      if (_datastoreInDsClusterToSelect != null) {
         boolean isDatastoreSelected = storagePage
               .selectDatastoreInDsCluster(_datastoreInDsClusterToSelect);
         verifyFatal(isDatastoreSelected, "Datastore is selected!");
      }
      // Wait for validation to complete
      new BaseView().waitForPageToRefresh();
      storagePage.waitApplySavedDataProgressBar();

      // Click on Next and verify that next page is loaded
      boolean isNextButtonClicked = storagePage.gotoNextPage();
      verifyFatal(isNextButtonClicked, "Verify the next button is clicked!");
   }
}
