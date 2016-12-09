/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Common workflow step for suspending VMs via API.
 * <p>
 *   Note that this step requires the VM to be <b>powered on</b> first.
 *   After finishing this step, VMs will <b>remain suspended</b>
 *   so it's user's responsibility to power them back on if needed.
 * </p>
 */
public class SuspendVmByApiStep extends BaseWorkflowStep {

   private List<VmSpec> _vmsToSuspend;

   @Override
   public void prepare() throws Exception {
      _vmsToSuspend = getSpec().getAll(VmSpec.class);

      if (CollectionUtils.isEmpty(_vmsToSuspend)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }
   }

   @Override
   public void execute() throws Exception {
      for (VmSpec vm : _vmsToSuspend) {
         verifyFatal(
               TestScope.FULL,
               VmSrvApi.getInstance().suspendVm(vm),
               String.format(
                     "Suspending VM with name '%s'", vm.name.get())
         );
      }
   }
}
