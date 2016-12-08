/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VAppSrvApi;

/**
 * Common workflow step for suspending vApps via API.
 * <p>
 *   Note that this step requires the vapp to be <b>powered on</b> first.
 *   After finishing this step, vapps will <b>remain suspended</b>
 *   so it's user's responsibility to power them back on if needed.
 * </p>
 */
public class SuspendVappByApiStep extends BaseWorkflowStep {
   private List<VappSpec> _vAppsToSuspend;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare() throws Exception {
      _vAppsToSuspend = getSpec().getAll(VappSpec.class);

      if (CollectionUtils.isEmpty(_vAppsToSuspend)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      for (VappSpec vapp : _vAppsToSuspend) {
         verifyFatal(
               TestScope.UI,
               VAppSrvApi.getInstance().suspendVapp(vapp),
               String.format(
                     "Suspending vApp with name '%s'", vapp.name.get())
         );
      }
   }
}
