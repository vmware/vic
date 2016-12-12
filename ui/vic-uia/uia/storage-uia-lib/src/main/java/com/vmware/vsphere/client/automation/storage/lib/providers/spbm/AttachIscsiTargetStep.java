package com.vmware.vsphere.client.automation.storage.lib.providers.spbm;

import java.util.List;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * Provider step that add iSCSI adapter to host/s
 */
public class AttachIscsiTargetStep implements ProviderWorkflowStep {

   private List<HostSpec> _hosts;

   @Override
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filterAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) throws Exception {
      if (isAssembling) {
         _hosts = filterAssemblerSpec.links.getAll(HostSpec.class);
      } else {
         _hosts = filteredPublisherSpec.links.getAll(HostSpec.class);
      }

   }

   @Override
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      for (HostSpec hostSpec : _hosts) {
         if (!HostBasicSrvApi.getInstance().createIscsiAdapter(hostSpec)) {
            throw new Exception(String.format(
                  "Unable to add iSCSI adapter to host '%s'",
                  hostSpec.name.get()));
         }
      }

   }

   @Override
   public boolean checkHealth() throws Exception {
      // TODO: how should check health be realized for a iscsi adapter
      return true;
   }

   @Override
   public void disassemble() throws Exception {
      for (HostSpec host : _hosts) {
         HostBasicSrvApi.getInstance().destroyIscsiAdapter(host);
      }
   }
}
