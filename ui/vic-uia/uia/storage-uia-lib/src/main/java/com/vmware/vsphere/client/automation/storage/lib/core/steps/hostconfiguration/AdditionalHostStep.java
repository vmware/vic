package com.vmware.vsphere.client.automation.storage.lib.core.steps.hostconfiguration;

import java.util.ArrayList;
import java.util.List;

import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;

/**
 * BaseWorkflowStep for providing additional hosts in the inventory
 */
public class AdditionalHostStep extends EnhancedBaseWorkflowStep {

   @UsesSpec
   private List<HostSpec> additionalHost;

   private final List<HostSpec> attachedHosts = new ArrayList<>();

   @Override
   public void execute() throws Exception {
      for (HostSpec hostSpec : additionalHost) {
         if (!HostBasicSrvApi.getInstance().checkHostExists(hostSpec)) {
            if (!HostBasicSrvApi.getInstance().addHost(hostSpec, true)) {
               throw new Exception(String.format("Unable to add host '%s'",
                     hostSpec));
            }
            attachedHosts.add(hostSpec);
         }
      }
   }

   @Override
   public void clean() throws Exception {
      for (HostSpec hostSpec : attachedHosts) {
         HostBasicSrvApi.getInstance().deleteHostSafely(hostSpec);
      }
   }
}
