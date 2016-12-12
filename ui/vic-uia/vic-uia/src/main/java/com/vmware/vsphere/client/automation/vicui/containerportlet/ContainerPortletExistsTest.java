package com.vmware.vsphere.client.automation.vicui.containerportlet;

import org.testng.annotations.Test;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.VmLocationSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.VmNavigationStep;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.vicui.common.VicEnvironmentProvider;
import com.vmware.vsphere.client.automation.vicui.common.VicUITestWorkflow;
import com.vmware.vsphere.client.automation.vicui.common.VicVcEnvSpec;
import com.vmware.vsphere.client.automation.vicui.common.step.ClickSummaryTabStep;
import com.vmware.vsphere.client.automation.vicui.containerportlet.step.VerifyContainerPortletAttributeStep;

/**
 * Test class for VCH VM portlet in the NGC client.
 * Executes the following test work-flow:
 *  1. Open a browser
 *  2. Login as admin user
 *  3. Navigate to the Container VM Summary tab
 *  4. Verify property "containerName" exists
 */ 

public class ContainerPortletExistsTest extends VicUITestWorkflow {	
	
	@Override
	public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
		TestbedSpecConsumer testBed = testbedBridge.requestTestbed(CommonTestBedProvider.class, true);
		TestbedSpecConsumer vveSpecConsumer = testbedBridge.requestTestbed(VicEnvironmentProvider.class, true);
		
		VicVcEnvSpec vveSpec = vveSpecConsumer.getPublishedEntitySpec(VicEnvironmentProvider.DEFAULT_ENTITY);
		
		// Spec for the VC
	    VcSpec requestedVcSpec = testBed.getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);
	    
	    // VmSpec for Container VM
	    VmSpec vmSpec = vveSpecConsumer.getPublishedEntitySpec(VicEnvironmentProvider.CONTAINER_VMSPEC_ENTITY);
	    
	    // Spec for the location to the VM
	    VmLocationSpec vmLocationSpec = new VmLocationSpec(vmSpec, NGCNavigator.NID_ENTITY_PRIMARY_TAB_SUMMARY);
	    
	    testSpec.add(requestedVcSpec, vveSpec, vmSpec, vmLocationSpec);
	    
	    super.initSpec(testSpec, testbedBridge);
	}
	
	@Override
	public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
		super.composeTestSteps(flow);
		
		flow.appendStep("Navigating to the Container VM", new VmNavigationStep());
		flow.appendStep("Clicking the Summary tab", new ClickSummaryTabStep());
	    flow.appendStep("Verifying a Container VM portlet property \"containerName\" exists", new VerifyContainerPortletAttributeStep());
	}
	
	@Override
	@Test(description = "Test if Container VM portlet exists")
	@TestID(id = "4")
	public void execute() throws Exception {
		super.execute();
	}
}
