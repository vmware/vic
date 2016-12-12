/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatacenterBasicSrvApi;

/**
 * Performs verification of datacenter existence via API call.
 */
public class VerifyDatacenterExistenceByApiStep extends BaseWorkflowStep {

   private DatacenterSpec _datacenterSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _datacenterSpec = filteredWorkflowSpec.get(DatacenterSpec.class);

      ensureNotNull(_datacenterSpec, "DatacenterSpec is missing.");
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.BAT,
            DatacenterBasicSrvApi.getInstance().checkDatacenterExists(_datacenterSpec),
            String.format("Verifying by API that datacenter %s exists",
                  _datacenterSpec.name.get()));
   }
}
