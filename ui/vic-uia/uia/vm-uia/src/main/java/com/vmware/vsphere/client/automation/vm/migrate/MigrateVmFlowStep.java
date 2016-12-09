/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.migrate;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.vm.migrate.spec.MigrateVmSpec;

/**
 * A <code>MigrateVmFlowStep</code> is an abstract steps which can be
 * extended by any custom steps which will require valid <code>CreateVmSpec</code> to be present for their
 * execution
 *
 */
public abstract class MigrateVmFlowStep extends CommonUIWorkflowStep {

    protected MigrateVmSpec migrateVmSpec;
    protected VmSpec vmSpec;

    /**
     * {@inheritDoc}
     */
    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
       migrateVmSpec = filteredWorkflowSpec.get(MigrateVmSpec.class);
       ensureNotNull(migrateVmSpec, "No MigrateVmSpec object was linked to the spec.");

       vmSpec = filteredWorkflowSpec.get(VmSpec.class);
       ensureNotNull(vmSpec, "No VmSpec object was linked to the spec.");
    }
}
