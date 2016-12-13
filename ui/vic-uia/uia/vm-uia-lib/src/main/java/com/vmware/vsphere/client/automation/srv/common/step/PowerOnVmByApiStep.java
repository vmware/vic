/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Common workflow step for powering on VMs via API.
 * <br>To use this step in automation tests:
 * <li>In the <code>initSpec()</code> method of the
 * <code>BaseTestWorkflow</code> test, create a <code>VmSpec</code> instances
 * and link them to the test spec.
 * <li>Append a <code>PowerOnVmByApiStep</code> instance to the test/prerequisite
 *  workflow composition.
 */
public class PowerOnVmByApiStep extends BaseWorkflowStep {

   private List<VmSpec> _vmsToPowerOn;
   private List<VmSpec> _vmsToPowerOff;

   @Override
   public void prepare() throws Exception {
      _vmsToPowerOn = getSpec().links.getAll(VmSpec.class);

      if (CollectionUtils.isEmpty(_vmsToPowerOn)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }

      _vmsToPowerOff = new ArrayList<VmSpec>();
   }

   @Override
   public void execute() throws Exception {
      for (VmSpec vm : _vmsToPowerOn) {
         if (VmSrvApi.getInstance().powerOnVm(vm)) {
            _vmsToPowerOff.add(vm);
         } else {
            throw new Exception(String.format(
                  "Unable to power on VM with name '%s'",
                  vm.name.get()));
         }
      }
   }

   @Override
   public void clean() throws Exception {
      for (VmSpec vm : _vmsToPowerOff) {
         VmSrvApi.getInstance().powerOffVm(vm);
      }
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmsToPowerOn = filteredWorkflowSpec.getAll(VmSpec.class);

      if (CollectionUtils.isEmpty(_vmsToPowerOn)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }

      _vmsToPowerOff = new ArrayList<VmSpec>();
   }
}
