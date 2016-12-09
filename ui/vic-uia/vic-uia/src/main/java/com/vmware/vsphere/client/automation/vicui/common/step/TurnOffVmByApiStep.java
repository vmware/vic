package com.vmware.vsphere.client.automation.vicui.common.step;

import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

public class TurnOffVmByApiStep extends BaseWorkflowStep {
	
	private List<VmSpec> _vmsToPowerOff;

	@Override
	public void execute() throws Exception {
		for (VmSpec vm : _vmsToPowerOff) {
			
			if(VmSrvApi.getInstance().isVmPoweredOn(vm)) {
				if(!VmSrvApi.getInstance().powerOffVm(vm)) {
					throw new Exception(String.format("Unable to power off VM '%s'", vm.name.get()));
				}
			}
		}
	}
	
	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
		_vmsToPowerOff = filteredWorkflowSpec.getAll(VmSpec.class);
		
		if(CollectionUtils.isEmpty(_vmsToPowerOff)) {
			throw new IllegalArgumentException("The spec has no links to 'VmSpec' instances");
		}
	}

}
