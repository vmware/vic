/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * A class that is used for entering maintenance mode for all the hosts,
 * specified in a test
 */
public class EnterMaintenanceModeStep extends BaseWorkflowStep {

	private List<HostSpec> hostsToEnterMM;
	private List<HostSpec> hostsToExitMM;

	@Override
	public void prepare() throws Exception {

		// Get all Hosts' specs
		hostsToEnterMM = getSpec().links.getAll(HostSpec.class);

		TestSpecValidator
				.ensureNotEmpty(hostsToEnterMM,
						"Please, provide HostSpecs for the hosts that need to enter maintenance mode!");

		hostsToExitMM = new ArrayList<HostSpec>();
	}

	@Override
	public void execute() throws Exception {
		for (HostSpec host : hostsToEnterMM) {
			if (!HostBasicSrvApi.getInstance().enterMaintenanceMode(host)) {
				throw new RuntimeException(String.format("Unable to enter maintenance mode for host '%s'",
						host.name.get()));
			}
			hostsToExitMM.add(host);
		}
	}

	@Override
	public void clean() throws Exception {
		for (HostSpec host : hostsToExitMM) {
			HostBasicSrvApi.getInstance().exitMaintenanceMode(host);
		}
	}

	// TestWorkflowStep methods
	@Override
	public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

		// Get all Hosts' specs
		hostsToEnterMM = filteredWorkflowSpec.getAll(HostSpec.class);

		TestSpecValidator
		.ensureNotEmpty(hostsToEnterMM,
				"Please, provide HostSpecs for the hosts that need to enter maintenance mode!");

		hostsToExitMM = new ArrayList<HostSpec>();
	}
}
