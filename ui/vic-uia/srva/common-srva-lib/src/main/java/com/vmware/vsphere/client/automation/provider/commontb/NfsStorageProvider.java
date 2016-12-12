/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.NfsStorageServiceSpec;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.BaseElementalProvider;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.connector.NfsStorageConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;

/**
 *
 */
public class NfsStorageProvider extends BaseElementalProvider {

   // Publisher info
   public static final String DEFAULT_ENTITY = "provider.nfs.storage.entity.default";

   private static final Logger _logger =
         LoggerFactory.getLogger(NfsStorageProvider.class);

   // Testbed publisher settings
   protected static final String TESTBED_KEY_FOLDER = "testbed.storage.folder";
   protected static final String TESTBED_KEY_IP = "testbed.storage.ip";
   protected static final String TESTBED_KEY_NAME = "testbed.storage.name";

   @Override
   public void initPublisherSpec(PublisherSpec publsherSpec) throws Exception {
      DatastoreSpec datastoreSpec = new DatastoreSpec();
      publsherSpec.links.add(datastoreSpec);
      datastoreSpec.service.set(new NfsStorageServiceSpec());
      publsherSpec.publishEntitySpec(DEFAULT_ENTITY, datastoreSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      DatastoreSpec datastoreSpec = publisherSpec.links.get(DatastoreSpec.class);

      String ip = SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_IP);
      String name = SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_NAME);
      String folder = SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_FOLDER);

      datastoreSpec.name.set(name);
      datastoreSpec.remoteHost.set(ip);
      datastoreSpec.remotePath.set(folder);
      datastoreSpec.type.set(DatastoreType.NFS);

      _logger.info("Loaded publisherSpec: " + datastoreSpec.toString());
   }



   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) throws Exception {
      _logger.info("NfsStorageProvider assignTestbedConnectors");
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof NfsStorageServiceSpec) {
            serviceConnectorsMap.put(serviceSpec, new NfsStorageConnector(
                  (NfsStorageServiceSpec) serviceSpec));
         }
      }
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {
      // TODO Auto-generated method stub

   }

   @Override
   public String determineResourceVersion() throws Exception {
      // TODO Auto-generated method stub
      return null;
   }

   @Override
   public void retrieveResource() throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public int providerWeight() {
      // TODO Auto-generated method stub
      return 0;
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter)
         throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      // TODO Auto-generated method stub
      return false;
   }

   @Override
   public void destroyTestbed() throws Exception {
      // TODO Auto-generated method stub

   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      // TODO Auto-generated method stub
      return null;
   }

}
