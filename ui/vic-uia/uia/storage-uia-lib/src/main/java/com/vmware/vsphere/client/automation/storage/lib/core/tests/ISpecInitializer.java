package com.vmware.vsphere.client.automation.storage.lib.core.tests;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;

/**
 * Interface for the spec initialization
 */
public interface ISpecInitializer {
   public void initSpec(final WorkflowSpec testSpec,
         final TestBedBridge testbedBridge);
}
