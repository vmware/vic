/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.vm.ft.view.FtSelectDatastorePage;

/**
 * Select a datastore on the Datastores page of the Fault Tolerance wizard
 */
public class FtSelectDatastoreStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private DatastoreSpec _datastoreSpec;

   @Override
   public void execute() throws Exception {
      FtSelectDatastorePage selectDatastorePage = new FtSelectDatastorePage();
      selectDatastorePage.waitForDialogToLoad();

      // Verify the datastore has been selected
      String datastoreToSelect = _datastoreSpec.name.get();
      verifyFatal(selectDatastorePage.selectDatastore(datastoreToSelect),
            "Verify the datastore has been selected.");

      // Wait for validation to complete
      selectDatastorePage.waitApplySavedDataProgressBar();

      // Click Next
      verifyFatal(selectDatastorePage.gotoNextPage(),
            "Verify navigation to the next page is successful.");
   }
}