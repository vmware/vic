/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.vm.ft.view.FtPortletView;

/**
 * Check the status and secondary VM location in the Fault Tolerance portlet on
 * the VM Summary page
 */
public class CheckFtPortletStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private FaultToleranceSpec _faultToleranceSpec;

   @UsesSpec
   private HostSpec _secondaryHostSpec;

   @Override
   public void execute() throws Exception {
      FtPortletView portletView = new FtPortletView();

      verifyFatal(portletView.waitForStatus(_faultToleranceSpec.status.get()),
            String.format("Status %s is reached", _faultToleranceSpec.status
                  .get().getValue()));

      verifyFatal(
            portletView.getSecondaryVmLocation().equals(_secondaryHostSpec.name.get()),
            String.format("Secondary VM Location is %s", _secondaryHostSpec.name.get()));
   }
}