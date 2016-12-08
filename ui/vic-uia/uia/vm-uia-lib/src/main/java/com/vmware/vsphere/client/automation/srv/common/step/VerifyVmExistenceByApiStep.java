/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Verifies whether VM exists through API call
 *
 */
public class VerifyVmExistenceByApiStep extends BaseWorkflowStep {

   protected VmSpec _vmSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmSpec = filteredWorkflowSpec.get(VmSpec.class);

      ensureNotNull(_vmSpec, "The spec has no links to 'VmSpec' instances");
      ensureAssigned(_vmSpec.name, "VM name is not set");
      ensureAssigned(_vmSpec.parent, "VM parent is not set");
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.BAT, VmSrvApi.getInstance().checkVmExists(_vmSpec),
            String.format("Verifying that VM %s exists!", _vmSpec.name.get()));
   }
}
