/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostStateSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * A class that is used for verifying host will enter certain state
 */
public class WaitHostToEnterStateStep extends BaseWorkflowStep {

    private HostSpec host;
    private HostStateSpec hostState;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      host = filteredWorkflowSpec.get(HostSpec.class);
      hostState = filteredWorkflowSpec.get(HostStateSpec.class);

      TestSpecValidator.ensureNotNull(host, "Please, provide a HostSpec!");
      TestSpecValidator.ensureNotNull(hostState, "Please, provide a HostStateSpec!");
   }

    @Override
    public void execute() throws Exception {
        boolean isHostInExpectedState = HostBasicSrvApi.getInstance().waitForHostToEnterState(host, hostState);
        verifySafely(
            isHostInExpectedState,
            "Host entered in the  expected " + hostState.state.get().toString() + " state.");
    }
}
