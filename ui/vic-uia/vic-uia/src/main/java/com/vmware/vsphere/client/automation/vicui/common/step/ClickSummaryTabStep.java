package com.vmware.vsphere.client.automation.vicui.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vicui.common.VicUIConstants;
import com.vmware.vsphere.client.automation.vicui.common.VicVcEnvSpec;

public class ClickSummaryTabStep extends CommonUIWorkflowStep {
	// This step is used to resolve the tab navigation issue on vCenter 6.0
	private boolean _isVC6_0;
	
	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
		VicVcEnvSpec vveSpec = filteredWorkflowSpec.get(VicVcEnvSpec.class);
		ensureNotNull(vveSpec, "vveSpec cannot be empty");
		_isVC6_0 = vveSpec.vcVersion.get().equals(VicUIConstants.VC_VERSION_6_0);
		
	}
	@Override
	public void execute() throws Exception {
		if(!_isVC6_0) {
			return;
		}

		LegacyPrimaryTabNav summaryNav = new LegacyPrimaryTabNav();
		summaryNav.selectPrimaryTab("Summary");
	}

}
