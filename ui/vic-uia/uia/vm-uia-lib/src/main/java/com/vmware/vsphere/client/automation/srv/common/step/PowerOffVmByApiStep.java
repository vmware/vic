/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Common workflow step for powering off VMs via API. <br>
 * To use this step in automation tests: <li>In the <code>initSpec()</code>
 * method of the <code>BaseTestWorkflow</code> test, create a
 * <code>VmSpec</code> instances and link them to the test spec. <li>Append a
 * <code>PowerOffVmByApiStep</code> instance to the test/prerequisite workflow
 * composition.
 */
public class PowerOffVmByApiStep extends BaseWorkflowStep {

   private List<VmSpec> _vmsToPowerOff;

   @Override
   public void execute() throws Exception {
      for (VmSpec vm : _vmsToPowerOff) {
         if (VmSrvApi.getInstance().powerOffVm(vm)) {
         } else {
            throw new Exception(String.format(
                  "Unable to power off VM with name '%s'", vm.name.get()));
         }
      }
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmsToPowerOff = filteredWorkflowSpec.getAll(VmSpec.class);

      if (CollectionUtils.isEmpty(_vmsToPowerOff)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }
   }
}
