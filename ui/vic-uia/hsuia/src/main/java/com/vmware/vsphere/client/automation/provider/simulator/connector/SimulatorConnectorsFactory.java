/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.simulator.connector;

import java.util.Map;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorServiceSpec;

/**
 * The factory is used for instantiating all kinds of simulated connectors. Such connectors
 * do not provide actual implementation for connecting to test beds. They are used for
 * testing purposes.
 * 
 * The API is designed to be integrated with the <code>ProviderWorkflow</code>
 * semantics.
 */
public class SimulatorConnectorsFactory {

	/**
	 * Call this method to instantiate and assign the required test bed connectors.
	 * 
	 * A separate connector instance is created for each service spec in the map.
	 * 
	 * No connector will be set if the factory does not know the service spec. However,
	 * the method will still finish successfully allowing the use of multiple factories
	 * on the same map.
	 * 
	 * @param serviceConnectorsMap
	 * 	Map between <code>ServiceSpecs<code> and <code>TestbedConnector</code>
	 * 	where the connectors are set.
	 */
	public static void createAndSetConnectors(
			Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) {
		
      if (serviceConnectorsMap == null) {
      	return;
      }

		for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
			// SimulatorServiceSpec -> SimulatorConnector
         if (serviceSpec instanceof SimulatorServiceSpec) {
         	serviceConnectorsMap.put(
         			serviceSpec, new SimulatorConnector((SimulatorServiceSpec)serviceSpec));
         } 
      }
	}
}
