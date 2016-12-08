/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.lib.clone.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.vm.lib.clone.spec.CloneVmSpec;

/**
 * A <code>CloneVmFlowStep</code> is an abstract steps which can be
 * extended by any custom steps which will require valid <code>CloneVmFlowStep</code>
 * to be present for their execution
 */
public abstract class CloneVmFlowStep extends CommonUIWorkflowStep {

    protected CloneVmSpec cloneVmSpec;
    protected VmSpec vmSpec;

    /**
     * {@inheritDoc}
     */
    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
       cloneVmSpec = filteredWorkflowSpec.get(CloneVmSpec.class);
       ensureNotNull(cloneVmSpec, "No CloneVmSpec object was linked to the spec.");

       vmSpec = filteredWorkflowSpec.get(VmSpec.class);
       ensureNotNull(vmSpec, "No VmSpec object was linked to the spec.");
    }
}
