/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.delay.Delay;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VAppSrvApi;

/**
 * Verifies whether vApp exists through API call
 * 
 */
public class VerifyVappExistenceByApiStep extends BaseWorkflowStep {

    protected VappSpec _vAppSpec;

    /**
     * {@inheritDoc}
     */
    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
        _vAppSpec = filteredWorkflowSpec.get(VappSpec.class);

        ensureNotNull(_vAppSpec, "The spec has no links to 'VappSpec' instances");
        ensureAssigned(_vAppSpec.name, "vApp name is not set");
        ensureAssigned(_vAppSpec.parent, "vApp parent is not set");
    }

    @Override
    public void execute() throws Exception {
        boolean vAppCreated = VAppSrvApi.getInstance().checkVAppExists(_vAppSpec);
        int retries = 10;

        // Wait for the vApp to be created
        while (!vAppCreated && retries < 0) {
            Delay.sleep.forSeconds(2).consume();
            vAppCreated = VAppSrvApi.getInstance().checkVAppExists(_vAppSpec);
            retries--;
        }

        verifyFatal(TestScope.BAT, vAppCreated, String.format("Verifying that vApp %s exists!", _vAppSpec.name.get()));
    }
}
