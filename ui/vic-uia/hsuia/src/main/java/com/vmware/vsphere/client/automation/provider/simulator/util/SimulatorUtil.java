/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.simulator.util;

import java.util.Random;
import java.util.UUID;

/**
 * Common helper methods used in the simulation providers.
 */
public final class SimulatorUtil {
	private SimulatorUtil() {
	}
	

	/**
	 * Use this method to simulate randomly generated failure.
	 * 
	 * @param allowRandomFailure
	 * 		Set to true if failure is allowed.
	 * 
	 * @throws Exception
	 * 		The exception is thrown to simulate failure.
	 */
	public static void trySucceed(boolean allowRandomFailure) throws Exception {
		boolean result = true;
		
		if (allowRandomFailure) {
			Random random = new Random();
			result = random.nextBoolean();
		}
		
		if (!result) {
			throw new Exception(
					"Operation failed for a random reason. It's not your fault. It's just a bad luck.");
		}
	}
	
	/**
	 * Returns a random string by using a fixed part and unique identifier.
	 */
	public static String getRandomString() {
		return "Simulator-" + UUID.randomUUID().toString();
	}
}
