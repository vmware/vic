package com.vmware.vsphere.client.automation.vicui.vchportlet.step;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;

public class VerifyVchPortletAttributeStep extends CommonUIWorkflowStep {
	private static final IDGroup VM_SUMMARY_VCHPORTLET_DOCKERAPIENDPOINT = IDGroup.toIDGroup("dockerApiEndpoint");
	
	@Override
	public void execute() throws Exception {
		verifyFatal(UI.condition
				.isFound(VM_SUMMARY_VCHPORTLET_DOCKERAPIENDPOINT)
				.await(SUITA.Environment.getPageLoadTimeout()), "Checking if dockerApiEndpoint is visible");
	}
	
}
