/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import com.vmware.client.automation.common.spec.EntitySpec;

/**
 * The interface provides access to the test bed specs.
 */
public interface TestbedSpecConsumer {

	/**
	 * Returns an entity spec for a given entity ID in the test bed.
	 * @param entitySpecId
	 *    Entity spec Id
	 * @return
	 *    Entity Spec
	 */
	public <T extends EntitySpec> T getPublishedEntitySpec(String entitySpecId);
}
