/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmBasicSrvApi;

/**
 * Common workflow step for creating VMs via the vc API.
 * <br>To use this step in automation tests:
 * <li>In the <code>initSpec()</code> method of the
 * <code>BaseTestWorkflow</code> test, create a <code>VmSpec</code> instances
 * and link them to the test spec.
 * <li>Append a <code>CreateVmByApiStep</code> instance to the test/prerequisite
 *  workflow composition.
 */
public class CreateVmByApiStep extends BaseWorkflowStep {

   protected List<VmSpec> _vmsToCreate;
   protected List<VmSpec> _vmsToDelete;

   @Override
   public void prepare() throws Exception {
      _vmsToCreate = getSpec().links.getAll(VmSpec.class);

      if (_vmsToCreate == null || _vmsToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances"
               );
      }

      _vmsToDelete = new ArrayList<VmSpec>();
   }

   @Override
   public void execute() throws Exception {
      for (VmSpec vm : _vmsToCreate) {
         if (VmBasicSrvApi.getInstance().createVm(vm)) {
            _vmsToDelete.add(vm);
         } else {
            throw new Exception(
                  String.format("Unable to create VM with name '%s'", vm.name.get())
                  );
         }
      }
   }

   @Override
   public void clean() throws Exception {
      for (VmSpec vmSpec : _vmsToDelete) {
         VmBasicSrvApi.getInstance().deleteVmSafely(vmSpec);
      }
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmsToCreate = filteredWorkflowSpec.getAll(VmSpec.class);

      if (_vmsToCreate == null || _vmsToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances"
            );
      }

      _vmsToDelete = new ArrayList<VmSpec>();
   }
}
