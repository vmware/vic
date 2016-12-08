/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.lib.clone.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.view.SelectComputeResourcePage;

/*
 * Workflow step that handles the selection of a resource on 'Select Compute Resource'
 * page of the Clone VM wizard.
 */
public class SelectComputeResourcePageStep extends CloneVmFlowStep {
   private ManagedEntitySpec _resourceToSelect;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      super.prepare(filteredWorkflowSpec);
      ensureNotNull(cloneVmSpec.targetComputeResource,
            "No compute resource is linked to the ClonedVmSpec.");
      _resourceToSelect = cloneVmSpec.targetComputeResource.get();
   }

   @Override
   public void execute() throws Exception {
      SelectComputeResourcePage selectComputeResourcePage = new SelectComputeResourcePage();
      selectComputeResourcePage.waitForDialogToLoad();

      // Select target resource
      selectComputeResourcePage.selectResource(_resourceToSelect);

      // Wait for validation to complete
      new BaseView().waitForPageToRefresh();
      selectComputeResourcePage.waitApplySavedDataProgressBar();

      // Click on Next and verify that next page is loaded
      boolean isNextButtonClicked = selectComputeResourcePage.gotoNextPage();
      verifyFatal(isNextButtonClicked, "Verify the next button is clicked!");
   }
}
