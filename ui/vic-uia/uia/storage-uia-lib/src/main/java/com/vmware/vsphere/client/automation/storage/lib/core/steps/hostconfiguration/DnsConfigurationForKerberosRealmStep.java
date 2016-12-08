package com.vmware.vsphere.client.automation.storage.lib.core.steps.hostconfiguration;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.host.ConfigChange;
import com.vmware.vim.binding.vim.host.ConfigManager;
import com.vmware.vim.binding.vim.host.DnsConfig;
import com.vmware.vim.binding.vim.host.NetStackInstance;
import com.vmware.vim.binding.vim.host.NetworkConfig;
import com.vmware.vim.binding.vim.host.NetworkConfig.NetStackSpec;
import com.vmware.vim.binding.vim.host.NetworkSystem;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.KerberosRealmSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.BasePrerequisiteStep;

/**
 * BasePrerequisiteStep implementation for performing the DNS configuration
 * based on KerberosRealmSpec
 */
public class DnsConfigurationForKerberosRealmStep extends BasePrerequisiteStep {

   /**
    * KerberosRealmSpec to extract the DNS configuration from
    */
   @UsesSpec
   private KerberosRealmSpec kerberosReamlSpec;

   /**
    * Spec describing the host to be configured
    */
   @UsesSpec
   private List<HostSpec> hostSpecs;

   /**
    * Saved initial state of the DNS configuration
    */
   private final Map<HostSpec, DnsConfig> initialDnsConfiguration = new HashMap<HostSpec, DnsConfig>();

   @Override
   public void execute() throws Exception {
      for (HostSpec hostSpec : hostSpecs) {
         final HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

         // Save current DNS configuration for the default TCP/IP stack
         initialDnsConfiguration
               .put(hostSpec,
                     getTheDefaultTcpIpStack(host.getConfig().network.netStackInstance).dnsConfig);

         // Prepare a spec for modification of the hosts default TCP/IP stack
         final NetworkConfig networkConfigSpec = getSpecForTheDefaultHostTcpIpStack(host);
         applyKerberosSpecConfiguration(networkConfigSpec.netStackSpec[0].netStackInstance.dnsConfig);

         // TODO: This should be in API class
         final NetworkSystem hostNetworkSystem = getHostNetworkSystem(host,
               hostSpec);
         hostNetworkSystem.updateNetworkConfig(networkConfigSpec,
               ConfigChange.Mode.modify.name());

      }
   }

   @Override
   public void clean() throws Exception {
      for (HostSpec hostSpec : hostSpecs) {
         final HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

         NetworkConfig initialNetworkConfigurationSpec = getSpecForTheDefaultHostTcpIpStack(host);
         initialNetworkConfigurationSpec.netStackSpec[0].netStackInstance.dnsConfig = this.initialDnsConfiguration
               .get(hostSpec);

         final NetworkSystem hostNetworkSystem = getHostNetworkSystem(host,
               hostSpec);
         hostNetworkSystem.updateNetworkConfig(initialNetworkConfigurationSpec,
               ConfigChange.Mode.modify.name());
      }
   }

   /**
    * Get a NetworkConfig instance populated for the modification of the hosts
    * Default TCP/IP stack
    *
    * @param host
    * @return
    */
   private NetworkConfig getSpecForTheDefaultHostTcpIpStack(
         final HostSystem host) {
      final NetStackSpec tcpIpStackSpec = new NetStackSpec();

      tcpIpStackSpec.netStackInstance = getTheDefaultTcpIpStack(host
            .getConfig().network.netStackInstance);
      tcpIpStackSpec.operation = ConfigChange.Operation.edit.name();

      final NetworkConfig result = new NetworkConfig();
      result.netStackSpec = new NetStackSpec[1];
      result.netStackSpec[0] = tcpIpStackSpec;
      return result;
   }

   /**
    * Filters the provided NetStackInstances and finds the NetStack intance for
    * the default TCP/IP Stack
    *
    * @param netStackInstance
    * @return
    */
   private NetStackInstance getTheDefaultTcpIpStack(
         NetStackInstance[] netStackInstance) {
      for (NetStackInstance currentNetStack : netStackInstance) {
         if (currentNetStack.key.equals("defaultTcpipStack")) {
            return currentNetStack;
         }
      }

      throw new RuntimeException("Can not find defaultTcpipStack in "
            + Arrays.asList(netStackInstance));
   }

   /**
    * Get instance of NetworkSystem for a given hsot
    *
    * @param host
    * @return
    * @throws Exception
    */
   private NetworkSystem getHostNetworkSystem(HostSystem host, HostSpec hostSpec)
         throws Exception {
      ConfigManager hostConfigManager = host.getConfigManager();
      ManagedObjectReference hostNetworkSystemMor = hostConfigManager
            .getNetworkSystem();

      NetworkSystem hostNetworkSystem = getService(hostNetworkSystemMor,
            hostSpec.service.get());
      return hostNetworkSystem;
   }

   /**
    * Apply the DNS configuration for the current instance of Step
    *
    * @param newDnsConfig
    *           the DnsConfig instance to apply the changes to
    * @return
    */
   private DnsConfig applyKerberosSpecConfiguration(DnsConfig newDnsConfig) {
      newDnsConfig.address[0] = kerberosReamlSpec.dnsServer.get();
      newDnsConfig.searchDomain = new String[1];
      newDnsConfig.searchDomain[0] = kerberosReamlSpec.activeDirectoryDomain
            .get();
      newDnsConfig.dhcp = false;

      return newDnsConfig;
   }

}
