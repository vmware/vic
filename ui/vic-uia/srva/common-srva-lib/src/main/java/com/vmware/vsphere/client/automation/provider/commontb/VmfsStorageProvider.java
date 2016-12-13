/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.VmfsStorageServiceSpec;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.BaseElementalProvider;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.connector.VmfsStorageConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;

/**
 * Provider for VMFS shared iSCSI storage
 */
public class VmfsStorageProvider extends BaseElementalProvider {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.vmfs.storage.entity.default";

   private static final Logger _logger = LoggerFactory.getLogger(VmfsStorageProvider.class);

   // Testbed publisher settings
   protected static final String TESTBED_KEY_IP = "testbed.vmfs.storage.ip";

   @Override
   public void initPublisherSpec(PublisherSpec publsherSpec) throws Exception {
      DatastoreSpec datastoreSpec = new DatastoreSpec();
      publsherSpec.links.add(datastoreSpec);
      datastoreSpec.service.set(new VmfsStorageServiceSpec());
      publsherSpec.publishEntitySpec(DEFAULT_ENTITY, datastoreSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec, SettingsReader testbedSettings) throws Exception {
      _logger.info("Start VmfsStorageProvider assign published specs");
      DatastoreSpec datastoreSpec = publisherSpec.getPublishedEntitySpec(VmfsStorageProvider.DEFAULT_ENTITY);

      String ip = SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_IP);

      datastoreSpec.remoteHost.set(ip);
      datastoreSpec.type.set(DatastoreType.VMFS);

      _logger.info("Loaded publisherSpec: " + datastoreSpec.toString());
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec, SettingsReader testbedSettings) throws Exception {
      _logger.debug("No assembler spec to laod.");
   }

   @Override
   public void assignTestbedConnectors(Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) throws Exception {
      _logger.info("VmfsStorageProvider assignTestbedConnectors");
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof VmfsStorageServiceSpec) {
            serviceConnectorsMap.put(serviceSpec, new VmfsStorageConnector((VmfsStorageServiceSpec) serviceSpec));
         }
      }
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return this.getClass();
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpec, AssemblerSpec filteredAssemblerSpec,
      boolean isAssembling, SettingsReader sessionSettingsReader) {
      throw new RuntimeException("Use NimbusVmfsStorageProvider or implement me!");
   }

   @Override
   public String determineResourceVersion() throws Exception {
      throw new RuntimeException("Use NimbusVmfsStorageProvider or implement me!");
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Use NimbusVmfsStorageProvider or implement me!");
   }

   @Override
   public int providerWeight() {
      return 0;
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter) throws Exception {
      throw new RuntimeException("Use NimbusVmfsStorageProvider or implement me!");
   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      throw new RuntimeException("Use NimbusVmfsStorageProvider or implement me!");
   }

   @Override
   public void destroyTestbed() throws Exception {
      throw new RuntimeException("Use NimbusVmfsStorageProvider or implement me!");
   }

}
