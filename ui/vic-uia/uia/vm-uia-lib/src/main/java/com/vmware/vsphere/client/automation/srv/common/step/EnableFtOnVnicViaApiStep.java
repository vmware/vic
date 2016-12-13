/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * Common workflow step that enables Fault Tolerance Logging via the API.
 */
public class EnableFtOnVnicViaApiStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private HostSpec _hostSpec;
   private boolean cleanUpFlag = false;

   @Override
   public void execute() throws Exception {
      if (HostBasicSrvApi.getInstance().checkFtLoggingState(false, _hostSpec)) {
         // Enable fault tolerance logging
         verifyFatal(HostBasicSrvApi.getInstance().enableFtLogging(_hostSpec),
               "Verify fault tolerance logging is enabled");

         // Mark for clean up
         cleanUpFlag = true;
      }
   }

   @Override
   public void clean() throws Exception {
      if (cleanUpFlag) {
         verifyFatal(HostBasicSrvApi.getInstance().disableFtLogging(_hostSpec),
               "Verify fault tolerance logging is disabled");
      }
   }
}