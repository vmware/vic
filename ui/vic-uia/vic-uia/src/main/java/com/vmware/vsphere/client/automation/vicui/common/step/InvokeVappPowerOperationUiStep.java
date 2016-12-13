package com.vmware.vsphere.client.automation.vicui.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.view.YesNoDialog;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

public class InvokeVappPowerOperationUiStep extends CommonUIWorkflowStep {
	private VmPowerStateSpec _vmPowerStateSpec;
	private VappSpec _vAppSpec;
	private TaskSpec _taskSpec;
	private static final IDGroup AI_POWER_OFF_VAPP = IDGroup.toIDGroup("vsphere.core.vApp.powerOffAction");
	private static final IDGroup AI_POWER_ON_VAPP = IDGroup.toIDGroup("vsphere.core.vApp.powerOnAction");
	
	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
		_vmPowerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);
		_vAppSpec = filteredWorkflowSpec.get(VappSpec.class);
		_taskSpec = new TaskSpec();
		
		ensureNotNull(_vmPowerStateSpec, "_vmPowerStateSpec cannot be empty");
		ensureNotNull(_vAppSpec, "_vAppSpec cannot be empty");
		ensureAssigned(_vmPowerStateSpec.powerState, "VM is not assigned to the power spec.");
		if(_vmPowerStateSpec.powerState.get() == VmPowerState.POWER_ON) {
			_taskSpec.name.set("Start vApp");
		} else {
			_taskSpec.name.set("Stop vApp");
		}
		
		_taskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
		_taskSpec.target.set(_vAppSpec);
	}
	
	@Override
	public void execute() throws Exception {
		_logger.info("vApp to work in the step: " + _vAppSpec.name.get());
		
		switch(_vmPowerStateSpec.powerState.get()) {
			case POWER_ON:
				ActionNavigator.invokeFromActionsMenu(AI_POWER_ON_VAPP);
				break;
			case POWER_OFF:
				ActionNavigator.invokeFromActionsMenu(AI_POWER_OFF_VAPP);
				YesNoDialog.CONFIRMATION.clickYes();
				break;
		}
		
		// Wait for tasks to complete
//		new BaseView().waitForRecentTaskCompletion();
		boolean isTaskFound = new TasksUtil().waitForRecentTaskToMatchFilter(new RecentTaskFilter(_taskSpec));
		verifyFatal(isTaskFound, String.format("Verifying task %s for target %s has reached status %s", _taskSpec.name.get(), _taskSpec.target.get(), _taskSpec.status.get()));
		
	}

}
