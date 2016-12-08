package com.vmware.vsphere.client.automation.vicui.common.step;

import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

public class TurnOnVmByApiStep extends BaseWorkflowStep {
	
	private List<VmSpec> _vmsToPowerOn;

	@Override
	public void execute() throws Exception {
		for (VmSpec vm : _vmsToPowerOn) {
			
			if(VmSrvApi.getInstance().isVmPoweredOff(vm)) {
				if(!VmSrvApi.getInstance().powerOnVm(vm)) {
					throw new Exception(String.format("Unable to power on VM '%s'", vm.name.get()));
				}
			}
		}
	}
	
	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
		_vmsToPowerOn = filteredWorkflowSpec.getAll(VmSpec.class);
		
		if(CollectionUtils.isEmpty(_vmsToPowerOn)) {
			throw new IllegalArgumentException("The spec has no links to 'VmSpec' instances");
		}
	}

}
