package com.vmware.vsphere.client.automation.storage.lib.core.steps.hostconfiguration;

import java.util.List;

import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.host.ConfigManager;
import com.vmware.vim.binding.vim.host.StorageSystem;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.KerberosRealmSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.BasePrerequisiteStep;

/**
 * BasePrerequisiteStep implementation for setting the nfs user
 */
public class SetHostKerberosNfsUserStep extends BasePrerequisiteStep {

   /**
    * KerberosRealmSpec to extract the NFS user from
    */
   @UsesSpec
   private KerberosRealmSpec kerberosReamlSpec;

   /**
    * Spec for the host to be used
    */
   @UsesSpec
   private List<HostSpec> hostSpecs;

   @UsesSpec
   private DatacenterSpec dc;

   @Override
   public void execute() throws Exception {
      for (HostSpec hostSpec : hostSpecs) {
         StorageSystem hostStorageSystem = getHostStorageSystem(hostSpec);
         hostStorageSystem.setNFSUser(kerberosReamlSpec.storageUserName.get(),
               kerberosReamlSpec.storagePassword.get());
      }
   }

   @Override
   public void clean() throws Exception {
      for (HostSpec hostSpec : hostSpecs) {
         getHostStorageSystem(hostSpec).clearNFSUser();
      }
   }

   /**
    * Get the storage system associated with the current instance host
    *
    * @return
    * @throws Exception
    */
   private StorageSystem getHostStorageSystem(HostSpec hostSpec)
         throws Exception {
      final HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

      ConfigManager hostConfigManager = host.getConfigManager();
      StorageSystem hostStorageSystem = getService(
            hostConfigManager.getStorageSystem(), hostSpec.service.get());
      return hostStorageSystem;
   }

}
