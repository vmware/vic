/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.simulator;

import java.util.List;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Sample provider
 */
public abstract class BaseCompositeProvider implements ProviderWorkflow {

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
   }

   @Override
   public void assignTestbedSettings(PublisherSpec providerSpec,
         SettingsReader testbedSettings) throws Exception {

      runWorkflowAssignServices(providerSpec);

   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
   }

   // Private methods

   private void runWorkflowAssignServices(BaseSpec containerSpec)
         throws ProviderWorkflowException {
      List<ManagedEntitySpec> specList =
            containerSpec.links.getAll(ManagedEntitySpec.class);

      // Elemental assembler workflow - no ManagedEntitySpec objects
      for (ManagedEntitySpec spec : specList) {
         assignParentService(spec);
      }
   }

   private void assignParentService(ManagedEntitySpec spec) {
      if (!spec.parent.isAssigned() || spec.parent.get() == null) {
         // Container spec- vc, host and etc.
         if (spec.service.get() == null) {
            throw new RuntimeException("Spec " + spec
                  + " is container and has no service assigned!");
         }
      } else {
         assignParentService(spec.parent.get());
         spec.service.set(spec.parent.get().service.get());
      }
   }

}
