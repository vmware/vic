/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Finish the Fault Tolerance wizard
 */
public class FtFinishWizardStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private VmSpec _vmSpec;

   @Override
   public void execute() throws Exception {
      WizardNavigator wizardNavigator = new WizardNavigator();
      boolean finishWizard = wizardNavigator.finishWizard();
      verifyFatal(finishWizard, "Verify wizard is closed");
   }

   /**
    * Turn off fault tolerance
    */
   @Override
   public void clean() throws Exception {
      verifyFatal(VmSrvApi.getInstance().turnOffFaultTolerance(_vmSpec),
            "Verify fault tolerance has been turned off");
   }
}