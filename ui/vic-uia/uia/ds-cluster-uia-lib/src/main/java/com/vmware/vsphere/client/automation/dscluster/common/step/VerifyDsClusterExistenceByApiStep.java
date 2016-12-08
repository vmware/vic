/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreClusterSrvApi;

public class VerifyDsClusterExistenceByApiStep extends BaseWorkflowStep {

   private DatastoreClusterSpec _dsClusterSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _dsClusterSpec = filteredWorkflowSpec.get(DatastoreClusterSpec.class);

      ensureNotNull(_dsClusterSpec, "DatastoreClusterSpec spec is missing.");
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(DatastoreClusterSrvApi.getInstance()
            .checkDsClusterExists(_dsClusterSpec), String.format(
            "Verifying that datastore cluster %s exists!",
            _dsClusterSpec.name.get()));
   }
}
