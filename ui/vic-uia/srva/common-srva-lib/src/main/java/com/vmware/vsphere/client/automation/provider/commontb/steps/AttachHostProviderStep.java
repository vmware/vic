/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.steps;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * Provider work flow step to attach a host.
 */
public class AttachHostProviderStep implements ProviderWorkflowStep {

   private HostSpec _hostToAttachSpec;

   @Override
   /**
    * Init the host spec with details of the host to be attached.
    */
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssmblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      if (isAssembling) {
         _hostToAttachSpec = filteredAssmblerSpec.links.get(HostSpec.class);
      } else {
         _hostToAttachSpec = filteredPublisherSpec.links.get(HostSpec.class);
      }

   }

   @Override
   /**
    * Check if the host specified by the step spec data is attached. If not attach it.
    */
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      if (!HostBasicSrvApi.getInstance().checkHostExists(_hostToAttachSpec)) {
         if (!HostBasicSrvApi.getInstance().addHost(_hostToAttachSpec, true)) {
            throw new Exception(String.format(
                  "Unable to add host '%s'",
                  _hostToAttachSpec.name.get()));
         }
      }
   }

   @Override
   /**
    * Return true if the host specified is connected.
    */
   public boolean checkHealth() throws Exception {
      // TODO: rkovachev Check if need other validation points.
      return HostBasicSrvApi.getInstance().isConnected(_hostToAttachSpec);
   }

   @Override
   /**
    * Remove the host from the inventory.
    */
   public void disassemble() throws Exception {
      HostBasicSrvApi.getInstance().deleteHostSafely(_hostToAttachSpec);

   }

}
