/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ops.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;
import com.vmware.vsphere.client.automation.vm.ops.views.VmsView;

/**
 * Verify VM power state via UI's Virtual Machines view
 */
public class VerifyVmPowerStateOnVmsViewStep extends CommonUIWorkflowStep {

   private VmPowerStateSpec _vmPowerStateSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmPowerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);

      ensureNotNull(_vmPowerStateSpec, "VmPowerStateSpec object is missing.");
      ensureAssigned(_vmPowerStateSpec.vm,
            "VM is not assigned to the power spec.");
      ensureAssigned(_vmPowerStateSpec.powerState,
            "Power state is not assigned to the power spec.");
   }

   @Override
   public void execute() throws Exception {
      VmsView vmsView = new VmsView();

      String actualPowerState = vmsView.getCellValue(
            _vmPowerStateSpec.vm.get().name.get(), VmsView.Column.STATE);
      String expectedPowerState = _vmPowerStateSpec.powerState.get()
            .getMessage();

      verifyFatal(TestScope.BAT, actualPowerState.equals(VmUtil
            .getLocalizedString(expectedPowerState)), String.format(
            "Verifying power state of %s by UI",
            _vmPowerStateSpec.vm.get().name.get()));
   }
}
