/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.dvs.step;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vim.binding.vim.ConfigSpecOperation;
import com.vmware.vim.binding.vim.DistributedVirtualSwitch;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.dvs.HostMember;
import com.vmware.vim.binding.vim.dvs.HostMember.PnicSpec;
import com.vmware.vim.binding.vim.host.HostProxySwitch;
import com.vmware.vim.binding.vim.host.NetworkInfo;
import com.vmware.vim.binding.vim.host.NetworkSystem;
import com.vmware.vim.binding.vim.host.PhysicalNic;
import com.vmware.vim.binding.vim.host.VirtualSwitch;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.spec.DvsSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;

/**
 * Adds hosts to existing distributed switch
 */
public class AddHostsToDistributedSwitchApiStep extends BaseWorkflowStep {

   private static final Logger _logger =
         LoggerFactory.getLogger(AddHostsToDistributedSwitchApiStep.class);

   private List<HostSpec> _hostSpecs;
   private DvsSpec _dvsSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _hostSpecs = filteredWorkflowSpec.getAll(HostSpec.class);
      if (_hostSpecs.isEmpty()) {
         throw new IllegalArgumentException("HostSpec objects are missing.");
      }

      _dvsSpec = filteredWorkflowSpec.get(DvsSpec.class);
      if (this._dvsSpec == null) {
         throw new IllegalArgumentException("DvsSpec object is missing.");
      }
   }
   
   @Override
   public void prepare() throws Exception {
      _hostSpecs = getSpec().getAll(HostSpec.class);
      if (_hostSpecs.isEmpty()) {
         throw new IllegalArgumentException("HostSpec objects are missing.");
      }

      _dvsSpec = getSpec().get(DvsSpec.class);
      if (this._dvsSpec == null) {
         throw new IllegalArgumentException("DvsSpec object is missing.");
      }
   }

   @Override
   public void execute() throws Exception {
      for (HostSpec host : this._hostSpecs) {
         // Build Member host spec
         HostMember.PnicSpec pnicSpec = new HostMember.PnicSpec();
         // first free pnic on the host
         String pnic = this.getFreePhysicalNicDevice(host).get(0);
         pnicSpec.setPnicDevice(pnic);

         HostMember.ConfigSpec hostConfigSpec = new HostMember.ConfigSpec();
         hostConfigSpec.setHost(this.getHostSystem(host)._getRef());
         hostConfigSpec.setBacking(new HostMember.PnicBacking(
               new PnicSpec[] { pnicSpec }));
         hostConfigSpec.setOperation(ConfigSpecOperation.add.name());

         // Build switch config spec
         DistributedVirtualSwitch.ConfigSpec dvsConfigSpec =
               new DistributedVirtualSwitch.ConfigSpec();
         dvsConfigSpec.setConfigVersion(this.getDvs(this._dvsSpec).getConfig()
               .getConfigVersion());
         dvsConfigSpec.setHost(new HostMember.ConfigSpec[] { hostConfigSpec });

         // Reconfigure dvs
         _logger
               .info(String
                     .format(
                           "Adding host '%s' with its free pnic adapter '%s' to distributed switch '%s'",
                           host.name.get(),
                           pnic,
                           this._dvsSpec.name.get()));
         ManagedObjectReference taskMoRef = this.getDvs(this._dvsSpec)
               .reconfigure(dvsConfigSpec);

         // Waits for task to complete
         verifyFatal(
               VcServiceUtil.waitForTaskSuccess(taskMoRef, host),
               String.format(
                     "Added host '%s' with its free pnic adapter '%s' to distributed switch '%s'",
                     host.name.get(), pnic, this._dvsSpec.name.get()));
      }
   }

   @Override
   public void clean() throws Exception {
   }

   private HostSystem getHostSystem(HostSpec hostSpec) throws Exception {
      return ManagedEntityUtil.getManagedObject(hostSpec, hostSpec.service.get());
   }

   private NetworkSystem getNetworkSystem(HostSpec hostSpec) throws Exception {
      HostSystem host = this.getHostSystem(hostSpec);
      return ManagedEntityUtil.getManagedObjectFromMoRef(host.getConfigManager()
            .getNetworkSystem(), hostSpec.service.get());
   }

   private NetworkInfo getNetworkInfo(HostSpec hostSpec) throws Exception {
      return this.getNetworkSystem(hostSpec).getNetworkInfo();
   }

   private VirtualSwitch[] getVirtualSwitch(HostSpec hostSpec) throws Exception {
      VirtualSwitch[] virtualSwitchesOnHost = this.getNetworkInfo(hostSpec).getVswitch();
      return virtualSwitchesOnHost == null ? new VirtualSwitch[] {}
            : virtualSwitchesOnHost;
   }

   private HostProxySwitch[] getHostProxySwitch(HostSpec hostSpec) throws Exception {
      HostProxySwitch[] proxySwitchesOnHost =
            this.getNetworkInfo(hostSpec).getProxySwitch();
      return proxySwitchesOnHost == null ? new HostProxySwitch[] {}
            : proxySwitchesOnHost;
   }

   private PhysicalNic[] getPhysicalNic(HostSpec hostSpec) throws Exception {
      return this.getNetworkInfo(hostSpec).getPnic();
   }

   private PhysicalNic getPhysicalNic(HostSpec hostSpec, String pnicKey)
         throws Exception {
      for (PhysicalNic pnic : this.getPhysicalNic(hostSpec)) {
         if (pnic.getKey().equals(pnicKey)) {
            return pnic;
         }
      }

      return null;
   }

   private List<String> getPhysicalNicDevice(HostSpec hostSpec) throws Exception {
      List<String> pnicKeys = new ArrayList<String>();
      for (PhysicalNic pnic : this.getPhysicalNic(hostSpec)) {
         pnicKeys.add(pnic.getDevice());
      }
      return pnicKeys;
   }

   private List<String> getFreePhysicalNicDevice(HostSpec hostSpec) throws Exception {
      List<String> allPnics = new ArrayList<String>();
      allPnics.addAll(this.getPhysicalNicDevice(hostSpec));

      List<String> usedPnics = new ArrayList<String>();
      for (VirtualSwitch virtualSwitch : this.getVirtualSwitch(hostSpec)) {
         for (String pnicKey : virtualSwitch.getPnic()) {
            usedPnics.add(this.getPhysicalNic(hostSpec, pnicKey).getDevice());
         }
      }

      for (HostProxySwitch proxySwitch : this.getHostProxySwitch(hostSpec)) {
         for (String pnicKey : proxySwitch.getPnic()) {
            usedPnics.add(this.getPhysicalNic(hostSpec, pnicKey).getDevice());
         }
      }

      List<String> freePnics = new ArrayList<String>();
      for (String pnicDevice : allPnics) {
         if (!usedPnics.contains(pnicDevice)) {
            freePnics.add(pnicDevice);
         }
      }

      return freePnics;
   }

   private DistributedVirtualSwitch getDvs(DvsSpec dvsSpec) throws Exception {
      return ManagedEntityUtil.getManagedObject(dvsSpec, dvsSpec.service.get());
   }
}
