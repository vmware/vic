package com.vmware.vsphere.client.automation.vicui.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VAppSrvApi;

public class EnsureVappIsOnStep extends BaseWorkflowStep {
	private VappSpec _vAppSpec;
	private VmSpec _vmSpec;

	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
		_vAppSpec = filteredWorkflowSpec.get(VappSpec.class);
		ensureNotNull(_vAppSpec, "vAppSpec cannot be null");
		
		_vmSpec = filteredWorkflowSpec.get(VmSpec.class);
		ensureNotNull(_vmSpec, "vmSpec cannot be null");
		_vAppSpec.name.set(_vmSpec.name.get());
		
	}
	
	@Override
	public void execute() throws Exception {
		// do not throw an exception unlike in PowerOnVappByApiStep.java. rather do nothing
		VAppSrvApi.getInstance().powerOnVapp(_vAppSpec);
	}
}
