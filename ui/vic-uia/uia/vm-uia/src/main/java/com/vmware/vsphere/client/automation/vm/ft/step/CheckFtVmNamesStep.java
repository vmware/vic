/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.util.FtUtil;
import com.vmware.vsphere.client.automation.vm.ft.view.ClusterVmView;

/**
 * Check the VM names are valid in the VMs > Virtual Machines tab
 */
public class CheckFtVmNamesStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private VmSpec _vmSpec;

   private final String primary = FtUtil.getLocalizedString("ft.vm.primary");
   private final String secondary = FtUtil
         .getLocalizedString("ft.vm.secondary");

   @Override
   public void execute() throws Exception {
      ClusterVmView vmView = new ClusterVmView();

      String _targetVm;
      _targetVm = _vmSpec.name.get();
      String primaryVmName = _targetVm + " " + primary;
      String secondaryVmName = _targetVm + " " + secondary;

      verifySafely(vmView.isVmPresentInGrid(primaryVmName),
            "Verify primary VM name");
      verifySafely(vmView.isVmPresentInGrid(secondaryVmName),
            "Verify secondary VM name");
   }
}