/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.lib.createvm.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec;

/**
 * A <code>CreateVmFlowStep</code> is an abstract steps which can be
 * extended by any custom steps which will require valid <code>CreateVmSpec</code> to be present for their
 * execution
 *
 */
public abstract class CreateVmFlowStep extends CommonUIWorkflowStep {

    protected CreateVmSpec createVmSpec;

    /**
     * {@inheritDoc}
     */
    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
       createVmSpec = filteredWorkflowSpec.get(CreateVmSpec.class);

       ensureNotNull(createVmSpec, "No CreateVmSpec object was linked to the spec.");
       ensureNotNull(createVmSpec.datastore, "No datastore object was linked to the CreateVmSpec.");
    }
}
