package com.vmware.vsphere.client.automation.vicui.common;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.connector.VcConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

public class VicEnvironmentProvider implements ProviderWorkflow {
	
	public static final String DEFAULT_ENTITY = "provider.vic.environment";
	public static final String VCH_VMSPEC_ENTITY = "provider.vic.vch_vm";
	public static final String CONTAINER_VMSPEC_ENTITY = "provider.vic.container_vm";
	
	private static final String ENDPOINT = "testbed.endpoint";
	private static final String USERNAME = "testbed.user";
	private static final String PASSWORD = "testbed.pass";
	public static final String VC_VER = "testbed.vc_version";
	public static final String VCH_VM_NAME = "testbed.vch_vm_name";
	public static final String CONTAINER_VM_NAME = "testbed.container_vm_name";

	private static final Logger _logger = LoggerFactory.getLogger(VicEnvironmentProvider.class);
	
	public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
		VicVcEnvSpec vicVcEnvSpec = new VicVcEnvSpec();
		publisherSpec.links.add(vicVcEnvSpec);
		publisherSpec.publishEntitySpec(DEFAULT_ENTITY, vicVcEnvSpec);
		
		VmSpec vchVmSpec = new VmSpec();
		publisherSpec.links.add(vchVmSpec);
		publisherSpec.publishEntitySpec(VCH_VMSPEC_ENTITY, vchVmSpec);
		
		VmSpec containerVmSpec = new VmSpec();
		publisherSpec.links.add(containerVmSpec);
		publisherSpec.publishEntitySpec(CONTAINER_VMSPEC_ENTITY, containerVmSpec);
		
	}
	
	public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge) throws Exception {
		TestbedSpecConsumer ctbProviderConsumer = testbedBridge.requestTestbed(CommonTestBedProvider.class, true);
		VcSpec requestedVcSpec = ctbProviderConsumer.getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);
		
		assemblerSpec.links.add(requestedVcSpec);
	}
	
	public void assignTestbedConnectors(Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) {
	      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
	    	  if (serviceSpec instanceof VcServiceSpec) {
	    		  try {
	    			  serviceConnectorsMap.put(serviceSpec, new VcConnector((VcServiceSpec) serviceSpec));
	    		  } catch (SsoException e) {
	    			  _logger.error(e.getMessage());
	    			  throw new RuntimeException("SSO authentication failed", e);
	    		  }
	    	  }
	      }
	}

	public void assignTestbedSettings(PublisherSpec publisherSpec, SettingsReader settingsReader) throws Exception {
		VcServiceSpec commonVcServiceSpec = new VcServiceSpec();
		commonVcServiceSpec.endpoint.set(SettingsUtil.getRequiredValue(settingsReader, ENDPOINT));
		commonVcServiceSpec.username.set(SettingsUtil.getRequiredValue(settingsReader, USERNAME));
		commonVcServiceSpec.password.set(SettingsUtil.getRequiredValue(settingsReader, PASSWORD));
		
		VicVcEnvSpec vveSpec = publisherSpec.getPublishedEntitySpec(DEFAULT_ENTITY);
		vveSpec.vcVersion.set(settingsReader.getSetting(VC_VER));
		vveSpec.service.set(commonVcServiceSpec);

		String vchVmName = SettingsUtil.getRequiredValue(settingsReader, VCH_VM_NAME);
		String containerVmName = SettingsUtil.getRequiredValue(settingsReader, CONTAINER_VM_NAME);
		VmSpec vchVmSpec = publisherSpec.getPublishedEntitySpec(VCH_VMSPEC_ENTITY);
		vchVmSpec.name.set(vchVmName);
		vchVmSpec.service.set(commonVcServiceSpec);
		
		VmSpec containerVmSpec = publisherSpec.getPublishedEntitySpec(CONTAINER_VMSPEC_ENTITY);
		containerVmSpec.name.set(containerVmName);
		containerVmSpec.service.set(commonVcServiceSpec);
		
	}

	public void assignTestbedSettings(AssemblerSpec assemblerSpec, SettingsReader settingsReader) throws Exception {
		
	}

	public void composeProviderSteps(WorkflowStepsSequence<? extends WorkflowStepContext> flow) throws Exception {

		
	}

	public Class<? extends ProviderWorkflow> getProviderBaseType() {
		return this.getClass();
	}

	public int providerWeight() {
		return 0;
	}

}
