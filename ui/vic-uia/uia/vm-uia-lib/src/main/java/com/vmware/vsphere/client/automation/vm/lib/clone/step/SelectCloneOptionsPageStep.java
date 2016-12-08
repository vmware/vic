/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.lib.clone.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.lib.clone.view.SelectCloneOptionsPage;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

/**
 * The step to select clone options
 */
public class SelectCloneOptionsPageStep extends CloneVmFlowStep {
   private VmPowerStateSpec _powerStateSpec;
   private Boolean _customizeGos;
   private Boolean _customizeHw;
   private Boolean _powerOnVm;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      super.prepare(filteredWorkflowSpec);

      if (cloneVmSpec.customizeGos.isAssigned()) {
         _customizeGos = cloneVmSpec.customizeGos.get();
      }
      if (cloneVmSpec.customizeHw.isAssigned()) {
         _customizeHw = cloneVmSpec.customizeHw.get();
      }

      if (cloneVmSpec.powerOnVm.isAssigned()) {
         _powerOnVm = cloneVmSpec.powerOnVm.get();

         _powerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);

         ensureNotNull(_powerStateSpec, "VmPowerStateSpec object is missing.");
         ensureAssigned(_powerStateSpec.vm,
               "VM is not assigned to the power spec.");
         ensureAssigned(_powerStateSpec.powerState,
               "Power state is not assigned to the power spec.");
      }
   }

   @Override
   public void execute() throws Exception {
      SelectCloneOptionsPage cloneOptionsPage = new SelectCloneOptionsPage();
      cloneOptionsPage.waitForDialogToLoad();

      if (_customizeGos != null) {
         cloneOptionsPage.setCustomizeGos(_customizeGos);
      }

      if (_customizeHw != null) {
         cloneOptionsPage.setCustomizeHw(_customizeHw);
      }

      if (_powerOnVm != null) {
         cloneOptionsPage.setPowerOnVm(_powerOnVm);
      }

      // Wait for validation to complete
      new BaseView().waitForPageToRefresh();
      cloneOptionsPage.waitApplySavedDataProgressBar();

      // Click on Next and verify that next page is loaded
      boolean isNextButtonClicked = cloneOptionsPage.gotoNextPage();
      verifyFatal(isNextButtonClicked, "Verify the next button is clicked!");
   }

}
