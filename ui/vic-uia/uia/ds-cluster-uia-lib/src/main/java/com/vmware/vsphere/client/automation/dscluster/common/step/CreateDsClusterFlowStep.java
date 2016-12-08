/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.dscluster.common.spec.CreateDsClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;

/**
 * A <code>CreateDsClusterFlowStep</code> is an abstract steps which can be
 * extended by any custom steps which will require valid
 * <code>DatastoreClusterSpec</code> to be present for their execution
 */
public abstract class CreateDsClusterFlowStep extends CommonUIWorkflowStep {
   protected DatastoreClusterSpec dsClusterSpec;
   protected CreateDsClusterSpec createDsClusterSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      dsClusterSpec = filteredWorkflowSpec.get(DatastoreClusterSpec.class);
      ensureNotNull(dsClusterSpec,
            "No DatastoreClusterSpec object was linked to the spec.");

      createDsClusterSpec = filteredWorkflowSpec.get(CreateDsClusterSpec.class);
      ensureNotNull(createDsClusterSpec,
            "No CreateDsClusterSpec object was linked to the spec.");
   }
}
