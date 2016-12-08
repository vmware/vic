package com.vmware.vsphere.client.automation.storage.lib.core.steps.hostconfiguration;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.lang3.ArrayUtils;

import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.host.ActiveDirectoryAuthentication;
import com.vmware.vim.binding.vim.host.ActiveDirectoryInfo;
import com.vmware.vim.binding.vim.host.AuthenticationManager;
import com.vmware.vim.binding.vim.host.AuthenticationStoreInfo;
import com.vmware.vim.binding.vim.host.ConfigManager;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.KerberosRealmSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.BasePrerequisiteStep;

/**
 * BasePrerequisiteStep implementation for Active Directory configuration based
 * on KerberosRealmSpec
 */
public class ActiveDirectoryConfigurationForKerberosRealmStep extends
      BasePrerequisiteStep {

   private static final String HOST_ACTIVE_DIRECTORY_AUTHENTICATION_TYPE = "HostActiveDirectoryAuthentication";

   /**
    * Spec for describing the Active Directory Domain name and credentials
    */
   @UsesSpec
   private KerberosRealmSpec kerberosReamlSpec;

   /**
    * Spec for the host to be used
    */
   @UsesSpec
   private List<HostSpec> hostSpecs;

   /**
    * Flag whether the host was already setup to use the same Active Directory
    * Domain
    */
   private final List<HostSpec> hostWithNoADPresetuped = new ArrayList<HostSpec>();

   @Override
   public void execute() throws Exception {
      for (HostSpec hostSpec : hostSpecs) {
         final HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
         ActiveDirectoryAuthentication activeDirectoryAuthenticationManager = getActiveDirectoryAuthenticationManager(
               host, hostSpec);

         ActiveDirectoryInfo initialActiveDirectoryConfiguration = getActiveDirectoryConfiguration(
               host, hostSpec);

         boolean isAdPresetuped = false;
         if (initialActiveDirectoryConfiguration.enabled) {
            if (initialActiveDirectoryConfiguration.joinedDomain
                  .equals(this.kerberosReamlSpec.activeDirectoryDomain.get())) {
               isAdPresetuped = true;
            } else {
               activeDirectoryAuthenticationManager.leaveCurrentDomain(true);
            }
         }

         if (!isAdPresetuped) {
            hostWithNoADPresetuped.add(hostSpec);
            ManagedObjectReference joinDomainTask = activeDirectoryAuthenticationManager
                  .joinDomain(kerberosReamlSpec.activeDirectoryDomain.get(),
                        kerberosReamlSpec.activeDirectoryUserName.get(),
                        kerberosReamlSpec.activeDirectoryPassword.get());
            if (!VcServiceUtil.waitForTaskSuccess(joinDomainTask, hostSpec)) {
               throw new RuntimeException(String.format(
                     "Task for joining active directory failed for spec: %s",
                     this.kerberosReamlSpec));
            }
         }
      }
   }

   @Override
   public void clean() throws Exception {
      for (HostSpec hostSpec : hostSpecs) {
         final HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
         getActiveDirectoryAuthenticationManager(host, hostSpec)
               .leaveCurrentDomain(true);
      }
   }

   /**
    * Get active directory authentication manager associated with the current
    * instance of the step
    *
    * @return
    * @throws Exception
    */
   private ActiveDirectoryAuthentication getActiveDirectoryAuthenticationManager(
         final HostSystem host, final HostSpec hostSpec) throws Exception {
      ConfigManager hostConfigManager = host.getConfigManager();

      AuthenticationManager hostAuthenticationManager = getService(
            hostConfigManager.getAuthenticationManager(),
            hostSpec.service.get());

      ActiveDirectoryAuthentication activeDirectoryAuthenticationManager = null;
      for (ManagedObjectReference supportedStore : hostAuthenticationManager
            .getSupportedStore()) {
         if (supportedStore.getType().equals(
               HOST_ACTIVE_DIRECTORY_AUTHENTICATION_TYPE)) {
            activeDirectoryAuthenticationManager = getService(supportedStore,
                  hostSpec.service.get());
         }
      }

      if (activeDirectoryAuthenticationManager == null) {
         throw new RuntimeException(
               String.format(
                     "Can not find ActiveDirectoryAuthentication manager for Host: %s",
                     hostSpec));
      }
      return activeDirectoryAuthenticationManager;
   }

   /**
    * Get active directory configuration for a host
    *
    * @param host
    * @return
    */
   private ActiveDirectoryInfo getActiveDirectoryConfiguration(
         final HostSystem host, final HostSpec hostSpec) {
      AuthenticationStoreInfo[] authConfig = host.getConfig().authenticationManagerInfo.authConfig;
      for (AuthenticationStoreInfo autheticationInfo : authConfig) {
         if (autheticationInfo instanceof ActiveDirectoryInfo) {
            return (ActiveDirectoryInfo) autheticationInfo;
         }
      }

      throw new RuntimeException(
            String.format(
                  "Can not find ActiveDirectoryInfo for host %s. Available AuthenticationStoreInfo: %s",
                  hostSpec, ArrayUtils.toString(authConfig)));

   }
}
