/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.migrate.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Verifies the Vm's host through API
 */
public class VerifyVmHostByApiStep extends BaseWorkflowStep {

   private VmSpec _vmSpec;
   private HostSpec _targetHost;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmSpec = filteredWorkflowSpec.get(VmSpec.class);
      ensureNotNull(_vmSpec, "The spec has no links to 'MigrateVmSpec' instances");


      _targetHost = filteredWorkflowSpec.get(HostSpec.class);
      ensureNotNull(_targetHost, "The spec has no links to 'HostSpec' instances");
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      verifyFatal(
            TestScope.BAT,
            VmSrvApi.getInstance().getVmHost(_vmSpec).equals(_targetHost.name.get()),
            String.format(
                  "Verify that the VM %s has parent host %s, actual: %s ",
                  _vmSpec.name.get(),
                  _targetHost.name.get(),
                  VmSrvApi.getInstance().getVmHost(_vmSpec)));
   }
}
