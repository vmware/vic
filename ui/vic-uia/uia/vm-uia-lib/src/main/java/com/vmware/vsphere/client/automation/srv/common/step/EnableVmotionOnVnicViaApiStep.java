/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * Common workflow step that enables vMotion via the API.
 */
public class EnableVmotionOnVnicViaApiStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private HostSpec _hostSpec;
   private boolean cleanUpFlag = false;

   @Override
   public void execute() throws Exception {
      if (HostBasicSrvApi.getInstance().checkVmotionState(false, _hostSpec)) {
         // Enable vMotion
         verifyFatal(HostBasicSrvApi.getInstance().enableVmotion(_hostSpec),
               "Verify vMotion is enabled");

         // Mark for clean up
         cleanUpFlag = true;
      }
   }

   @Override
   public void clean() throws Exception {
      if (cleanUpFlag) {
         verifyFatal(HostBasicSrvApi.getInstance().disableVmotion(_hostSpec),
               "Verify vMotion is disabled");
      }
   }
}