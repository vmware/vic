package com.vmware.vsphere.client.automation.storage.lib.providers.spbm;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.suitaf.SUITA;
import com.vmware.vsphere.client.automation.srv.common.spec.StorageProviderSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.StorageProviderBasicSrvApi;

public class AttachVasaProviderStep implements ProviderWorkflowStep {

   private StorageProviderSpec _storageProvider;

   @Override
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filterAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) throws Exception {
      if (isAssembling) {
         _storageProvider = filterAssemblerSpec.links
               .get(StorageProviderSpec.class);
      } else {
         _storageProvider = filteredPublisherSpec.links
               .get(StorageProviderSpec.class);
      }

   }

   @Override
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      if (!StorageProviderBasicSrvApi.getInstance().isStorageProviderPresent(
            _storageProvider)) {

         boolean isProviderCreated = false;
         long timeout = SUITA.Environment.getBackendJobMid();
         long endTime = System.currentTimeMillis() + timeout;
         while (!isProviderCreated && System.currentTimeMillis() < endTime) {
            isProviderCreated = StorageProviderBasicSrvApi.getInstance()
                  .createStorageProvider(_storageProvider);
         }
         if (!isProviderCreated) {
            throw new Exception(String.format(
                  "Unable to add storage provider  '%s'",
                  _storageProvider.name.get()));
         }
      }

   }

   @Override
   public boolean checkHealth() throws Exception {
      // TODO: how should check health be realized for a storage provider :
      // StorageProviderBasicSrvApi.getInstance().isStorageProviderPresent(
      // _storageProvider);
      return true;
   }

   @Override
   public void disassemble() throws Exception {
      StorageProviderBasicSrvApi.getInstance().deleteStorageProvider(
            _storageProvider);

   }

}
