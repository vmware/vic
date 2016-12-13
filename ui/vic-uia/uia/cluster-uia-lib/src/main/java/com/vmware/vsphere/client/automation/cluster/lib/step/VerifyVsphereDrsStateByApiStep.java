/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.cluster.lib.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;

/**
 * Performs verification of cluster DRS state via API call.
 */
public class VerifyVsphereDrsStateByApiStep extends BaseWorkflowStep {

    protected ClusterSpec _clusterSpec;
    protected Boolean _expectedDrsState;

    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
        _clusterSpec = filteredWorkflowSpec.get(ClusterSpec.class);
        ensureNotNull(_clusterSpec, "ClusterSpec spec is missing.");

        _expectedDrsState = _clusterSpec.drsEnabled.get();
        ensureNotNull(_expectedDrsState, "drsEnabled property for ClusterSpec is not found.");
    }

   @Override
   public void execute() throws Exception {
      verifyFatal(new EqualsAssertion(ClusterBasicSrvApi.getInstance()
            .isDrsEnabled(_clusterSpec), _expectedDrsState, String.format(
            "Verifying DRS state for cluster %s", _clusterSpec.name.get())));
   }
}