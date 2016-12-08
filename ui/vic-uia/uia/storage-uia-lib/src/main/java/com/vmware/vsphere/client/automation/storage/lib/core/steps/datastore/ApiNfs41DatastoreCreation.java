package com.vmware.vsphere.client.automation.storage.lib.core.steps.datastore;

import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Datastore.Info;
import com.vmware.vim.binding.vim.host.DatastoreSystem;
import com.vmware.vim.binding.vim.host.NasDatastoreInfo;
import com.vmware.vim.binding.vim.host.NasVolume;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.ApiOperationStep;

/**
 * ApiOperationStep implementation for creation of a NFS 4.1 Datastore
 */
public class ApiNfs41DatastoreCreation extends ApiOperationStep {

   @UsesSpec
   private Nfs41DatastoreSpec datastoreSpec;

   @UsesSpec
   private HostSpec host;

   private Datastore similarDatastore;

   @Override
   protected Similarity checkPresence() {
      for (Datastore datastore : HostBasicSrvApi.getInstance().getDatastores(
            host)) {
         Info dsInfo = datastore.getInfo();
         if (dsInfo instanceof NasDatastoreInfo) {
            NasDatastoreInfo nasDsInfo = (NasDatastoreInfo) dsInfo;

            if (isDatastoreBackingMatch(nasDsInfo.nas)) {
               if (nasDsInfo.getName().equals(datastoreSpec.name.get())) {
                  return Similarity.MATCH;
               } else {
                  this.similarDatastore = datastore;
                  return Similarity.SIMILAR;
               }
            }
         }
      }

      return Similarity.NO;
   }

   private boolean isDatastoreBackingMatch(NasVolume nas) {
      return nas.remoteHost.equals(datastoreSpec.remoteHost.get())
            && nas.remotePath.equals(datastoreSpec.remotePath.get())
            && nas.type.equals(datastoreSpec.nfsVersion.apiFileSystemEnum
                  .toString())
            && nas.securityType
                  .equals(datastoreSpec.authenticationMode.get().apiAuthentication
                        .toString());
   }

   @Override
   protected CleanOperation modify() throws Exception {
      final String initialName = this.similarDatastore.getName();

      this.similarDatastore.rename(this.datastoreSpec.name.get());

      return new CleanOperation() {

         @Override
         public void execute() throws Exception {
            similarDatastore.rename(initialName);
         }
      };
   }

   @Override
   protected CleanOperation perform() throws Exception {
      final NasVolume.Specification apiSpec = getApiCreationSpec();
      final DatastoreSystem datastoreSystem = HostBasicSrvApi.getInstance()
            .getDatastoreSystem(host);
      final ManagedObjectReference createdDatastore = datastoreSystem
            .createNasDatastore(apiSpec);
      if (createdDatastore == null) {
         throw new RuntimeException(
               "Datastore creation operation returned null datastore for spec "
                     + datastoreSpec);
      }

      return new CleanOperation() {
         @Override
         public void execute() throws Exception {
            datastoreSystem.removeDatastore(createdDatastore);
         }
      };
   }

   private NasVolume.Specification getApiCreationSpec() {
      NasVolume.Specification apiSpec = new NasVolume.Specification();
      apiSpec.setLocalPath(datastoreSpec.name.get());
      apiSpec.remoteHost = datastoreSpec.remoteHost.get();
      apiSpec
            .setRemoteHostNames(new String[] { datastoreSpec.remoteHost.get() });
      apiSpec.setRemotePath(datastoreSpec.remotePath.get());
      apiSpec
            .setAccessMode(datastoreSpec.accessMode.get().apiDatastoreAccessMode
                  .toString());
      apiSpec.setType(datastoreSpec.nfsVersion.apiFileSystemEnum.toString());
      apiSpec
            .setSecurityType(datastoreSpec.authenticationMode.get().apiAuthentication
                  .toString());
      if (datastoreSpec.remoteHostAdresses.isAssigned()) {
         apiSpec.setRemoteHostNames(datastoreSpec.remoteHostAdresses.get());
      }
      return apiSpec;
   }

}
