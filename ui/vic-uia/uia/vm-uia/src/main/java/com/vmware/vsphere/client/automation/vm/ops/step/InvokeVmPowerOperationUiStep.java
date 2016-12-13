/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ops.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.VmGlobalActions;
import com.vmware.vsphere.client.automation.common.view.YesNoDialog;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

/**
 * Invoke power on/off action from VM's action menu
 */
public class InvokeVmPowerOperationUiStep extends CommonUIWorkflowStep {

   private VmPowerStateSpec _powerStateSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _powerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);

      ensureNotNull(_powerStateSpec, "VmPowerStateSpec object is missing.");
      ensureAssigned(_powerStateSpec.vm, "VM is not assigned to the power spec.");
      ensureAssigned(
            _powerStateSpec.powerState,
            "Power state is not assigned to the power spec.");
   }

   @Override
   public void execute() throws Exception {
      _logger.info("VM to work in the step: " + _powerStateSpec.vm.get().name.get());
      boolean isPowerOnOp = true;
      switch (_powerStateSpec.powerState.get()) {
         case POWER_ON:
            ActionNavigator.invokeFromActionsMenu(VmGlobalActions.AI_POWER_ON_VM);
            break;
         case POWER_OFF:
            ActionNavigator.invokeFromActionsMenu(VmGlobalActions.AI_POWER_OFF_VM);
            // Confirm power off
            YesNoDialog.CONFIRMATION.clickYes();
            isPowerOnOp = false;
            break;
      }

      // Wait for tasks to complete
      new BaseView().waitForRecentTaskCompletion();

      verifyFatal(VmSrvApi.getInstance().waitForVmPowerState(_powerStateSpec.vm.get(), isPowerOnOp),
                  String.format("Verifying the VM %s reached the desired power state %s",
                  _powerStateSpec.vm.get().name.get(),
                  _powerStateSpec.powerState.get()));
   }

   @Override
   public void clean() throws Exception {
      _logger.info("VM to clean after the step: " + _powerStateSpec.vm.get().name.get());
      // Power off the VM if needed
      if (VmSrvApi.getInstance().isVmPoweredOn(_powerStateSpec.vm.get())) {
         VmSrvApi.getInstance().powerOffVm(_powerStateSpec.vm.get());
      }
   }
}
