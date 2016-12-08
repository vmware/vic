/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.simulator.connector;

import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorServiceSpec;

/**
 * Test bed connector implementation used for simulation purposes.
 * It does not hold actual connection.
 */
public final class SimulatorConnector implements TestbedConnector {

	public SimulatorConnector(SimulatorServiceSpec serviceSpec) {
	}
	
	@Override
	public void connect() {
	}

	@Override
	public boolean isAlive() {
		return true;
	}

	@Override
	public void disconnect() {
	}

	@Override
	public <T> T getConnection() {
		return null;
	}
}
