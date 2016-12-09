package com.vmware.vsphere.client.automation.vicui.plugininstalled;

import org.testng.annotations.Test;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.vicui.common.VicUITestWorkflow;
import com.vmware.vsphere.client.automation.vicui.plugininstalled.spec.AdminNavigationSpec;
import com.vmware.vsphere.client.automation.vicui.plugininstalled.step.AdminNavigationStep;
import com.vmware.vsphere.client.automation.vicui.plugininstalled.step.FindVicUIStep;

/**
 * Test class for VCH VM portlet in the NGC client.
 * Executes the following test work-flow:
 *  1. Open a browser
 *  2. Login as admin user
 *  3. Navigate to Administration -> Client Plug-Ins 
 *  4. Verify if item "VicUI" exists
 */ 

public class PluginInstalledTest extends VicUITestWorkflow {
	@Override
	public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
		TestbedSpecConsumer testBed = testbedBridge.requestTestbed(CommonTestBedProvider.class, true);
		VcSpec requestedVcSpec = testBed.getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);
		AdminNavigationSpec adminNavigationSpec = new AdminNavigationSpec();
		
		testSpec.add(requestedVcSpec, adminNavigationSpec);
		
		super.initSpec(testSpec, testbedBridge);
	}
	
	@Override
	public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
		super.composeTestSteps(flow);
		
		flow.appendStep("Navigating to the administration -> client plugins menu", new AdminNavigationStep());
		flow.appendStep("Verifying if VicUI is installed", new FindVicUIStep());
	}
	
	@Override
	@Test(description = "Test if VIC UI plugin is installed correctly")
	@TestID(id = "0")
	public void execute() throws Exception {
		super.execute();
	}
}
