/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import com.vmware.client.automation.common.spec.EntitySpec;

/**
 * The interface provides method to published entity spec.
 *
 */
public interface TestbedSpecProvider {
   public void publishEntitySpec(String id, EntitySpec entitySpec);
}
