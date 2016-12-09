/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;

/**
 * Host Provider based on vCloud deployment.
 */
public class OneCloudHostProvider extends BaseHostProvider {

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
      throw new RuntimeException("Implement me!");
   }
   
   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return BaseHostProvider.class;
   }
}
