package com.vmware.vsphere.client.automation.vicui.vchportlet.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.vm.lib.ops.model.VmOpsModel.VmPowerState;
import com.vmware.vsphere.client.automation.vm.lib.ops.spec.VmPowerStateSpec;

public class VerifyVchDockerEndpointIsValidStep extends CommonUIWorkflowStep {
	
	private VmPowerStateSpec _vmPowerStateSpec;
	private static final IDGroup VM_SUMMARY_VCHPORTLET_DOCKERAPIENDPOINT = IDGroup.toIDGroup("dockerApiEndpoint");
	private static final String DOCKER_API_ENDPOINT_PLACEHOLDER_VALUE = "-";
	
	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
		_vmPowerStateSpec = filteredWorkflowSpec.get(VmPowerStateSpec.class);
		
		ensureNotNull(_vmPowerStateSpec, "VmPowerStateSpec cannot be null");
		ensureAssigned(_vmPowerStateSpec.vm, "VM is not assigned to the power spec.");
		ensureAssigned(_vmPowerStateSpec.powerState, "Power state is not assigned to the power spec.");
	}

	@Override
	public void execute() throws Exception {
		if(_vmPowerStateSpec.powerState.get().equals(VmPowerState.POWER_ON)) {
			verifyFatal(!UI.component.property.get(Property.TEXT, VM_SUMMARY_VCHPORTLET_DOCKERAPIENDPOINT).equalsIgnoreCase(DOCKER_API_ENDPOINT_PLACEHOLDER_VALUE), "Verifying \"dockerApiEndpoint\" is not \"-\"");
			
		} else {
			verifyFatal(UI.component.property.get(Property.TEXT, VM_SUMMARY_VCHPORTLET_DOCKERAPIENDPOINT).equalsIgnoreCase(DOCKER_API_ENDPOINT_PLACEHOLDER_VALUE), "Verifying \"dockerApiEndpoint\" is \"-\"");
		}
	}
}

