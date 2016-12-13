/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;

/**
 * Performs verification of cluster existence via API call.
 */
public class VerifyClusterExistanceByApiStep extends BaseWorkflowStep {

    private ClusterSpec _clusterSpec;

    /**
     * {@inheritDoc}
     */
    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
        _clusterSpec = filteredWorkflowSpec.get(ClusterSpec.class);

        ensureNotNull(_clusterSpec, "ClusterSpec spec is missing.");
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void execute() throws Exception {
        verifyFatal(ClusterBasicSrvApi.getInstance().checkClusterExists(_clusterSpec),
                String.format("Verifying that cluster %s exists!", _clusterSpec.name.get()));
    }
}
