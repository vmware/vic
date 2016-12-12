package com.vmware.vsphere.client.automation.storage.lib.providers.nfs;

import java.util.Map;

import org.apache.commons.lang.NotImplementedException;
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
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec.Nfs41AuthenticationMode;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.NfsDatastoreSpec.DatastoreAccessMode;

/**
 * BaseElementalProvider for NFS 4.1 Storage
 */
public class Nfs41StorageProvider extends BaseElementalProvider {

   public static final String DEFAULT_ENTITY = "provider.nfs41.storage.entity.default";

   protected static final String TESTBED_KEY_IPS = "testbed.storage.ips";
   protected static final String TESTBED_KEY_FOLDER = "testbed.storage.folder";
   protected static final String TESTBED_KEY_NAME = "testbed.storage.name";
   protected static final String TESTBED_KEY_ACCESS_MODE = "testbed.storage.accessmode";

   private static final Logger _logger = LoggerFactory
         .getLogger(Nfs41StorageProvider.class);

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      Nfs41DatastoreSpec datastoreSpec = new Nfs41DatastoreSpec();
      datastoreSpec.service.set(new NfsStorageServiceSpec());
      publisherSpec.links.add(datastoreSpec);
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, datastoreSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.debug("No assembler spec to laod.");
      Nfs41DatastoreSpec datastoreSpec = publisherSpec
            .getPublishedEntitySpec(Nfs41StorageProvider.DEFAULT_ENTITY);
      datastoreSpec.authenticationMode.set(Nfs41AuthenticationMode.DISABLED);
      NfsStorageServiceSpec serviceSpec = (NfsStorageServiceSpec) datastoreSpec.service
            .get();

      String[] ips = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_IPS).split(",");
      datastoreSpec.remoteHost.set(ips[0]);
      datastoreSpec.remoteHostAdresses.set(ips);
      serviceSpec.endpoint.set(ips[0]);

      String folder = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_FOLDER);
      datastoreSpec.remotePath.set(folder);
      String name = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_NAME);
      datastoreSpec.name.set(name);
      String accessMode = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_ACCESS_MODE);
      datastoreSpec.accessMode.set(DatastoreAccessMode.valueOf(accessMode));

      _logger.debug(String.format("Initialized Nfs41DatastoreSpec spec '%s' ",
            datastoreSpec));
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.debug("No assembler spec to laod.");
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap)
         throws Exception {
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof NfsStorageServiceSpec) {
            // TODO rename the KerberizedNfsConnector to NFSConnector prior to
            // submit for review
            serviceConnectorsMap.put(serviceSpec, new KerberizedNfsConnector(
                  serviceSpec));
         }
      }
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return null;
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {
      throw new NotImplementedException();
   }

   @Override
   public String determineResourceVersion() throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public int providerWeight() {
      throw new NotImplementedException();
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter)
         throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      throw new NotImplementedException();
   }

   @Override
   public void destroyTestbed() throws Exception {
      throw new NotImplementedException();
   }

}
