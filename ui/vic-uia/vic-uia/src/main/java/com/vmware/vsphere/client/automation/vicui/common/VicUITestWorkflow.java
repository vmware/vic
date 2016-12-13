package com.vmware.vsphere.client.automation.vicui.common;

import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

public class VicUITestWorkflow extends NGCTestWorkflow {
	/*
	 * This custom NGCTestWorkflow class disables admin user creation and forces the framework to use the administrator@vsphere.local account
	 * because vCenter 6.0 has an issue with SSO in the NGC plugin test project
	 */

	@Override
	protected UserSpec generateUserSpec(VcSpec vcSpec) {
		UserSpec userSpec = new UserSpec();
		userSpec.username.set(System.getProperty("vcAdminUsername"));
		userSpec.password.set(System.getProperty("vcAdminPassword"));
		userSpec.parent.set(vcSpec);
		userSpec.tag.set(NGCTestWorkflow.TEST_USER_SPEC_TAG);
		
		return userSpec;
	}
	
	@Override
	public void composePrereqSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
		
	}
}
