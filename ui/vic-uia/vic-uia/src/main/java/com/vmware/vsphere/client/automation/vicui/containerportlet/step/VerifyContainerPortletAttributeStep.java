package com.vmware.vsphere.client.automation.vicui.containerportlet.step;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.SUITA;

public class VerifyContainerPortletAttributeStep extends CommonUIWorkflowStep {
	private static final IDGroup VM_SUMMARY_CONTAINERPORTLET_NAME = IDGroup.toIDGroup("containerName");
	
	@Override
	public void execute() throws Exception {
                UI.condition.isFound(VM_SUMMARY_CONTAINERPORTLET_NAME).await(SUITA.Environment.getPageLoadTimeout());
		verifyFatal(UI.component.exists(VM_SUMMARY_CONTAINERPORTLET_NAME), "Checking if containerName is available");
	}

}
