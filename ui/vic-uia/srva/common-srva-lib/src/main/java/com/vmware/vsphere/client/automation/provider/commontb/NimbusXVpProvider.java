package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.Map;

import org.apache.commons.lang.NotImplementedException;
import org.apache.commons.lang.RandomStringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.XVpServiceSpec;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.BaseElementalProvider;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusServiceSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.XVpProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.connector.XVpConnector;
import com.vmware.vsphere.client.automation.provider.util.NimbusCommandsUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.StorageProviderSpec;

/**
 * Provider that deploys xVp on Nimbus using the nimbus-ovfdeploy command and
 * provides the ability to assign to VC already deployed xVp.
 */
public class NimbusXVpProvider extends BaseElementalProvider {
   public static final String DEFAULT_ENTITY = "provider.storage.vp.entity.default";

   // session info
   private static final String RESOURCE_KEY_VERSION = "resource.vp.version";
   private static final String RESOURCE_KEY_REPL_CONFIG = "resource.vp.replConfig";
   private static final String RESOURCE_KEY_URL = "resource.vp.url";

   // Assembler Spec
   private XVpProvisionerSpec vpProvisionerSpec;

   private static final String XVP_URL_TEMPLATE = "https://%s:8443/vasa/version.xml";

   // Testbed publisher settings
   protected static final String TESTBED_KEY_ENDPOINT = "testbed.storage.vp.endpoint";
   protected static final String TESTBED_KEY_URL = "testbed.storage.vp.url";
   protected static final String TESTBED_KEY_USERNAME = "testbed.storage.vp.username";
   protected static final String TESTBED_KEY_PASSWORD = "testbed.storage.vp.password";

   // Publisher Spec
   private StorageProviderSpec storageProviderSpec;

   // logger
   private static final Logger logger = LoggerFactory
         .getLogger(NimbusXVpProvider.class);

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec,
         TestBedBridge testbedBridge) throws Exception {
      vpProvisionerSpec = new XVpProvisionerSpec();
      assemblerSpec.links.add(vpProvisionerSpec);
   }

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      storageProviderSpec = new StorageProviderSpec();
      storageProviderSpec.service.set(new XVpServiceSpec());
      publisherSpec.links.add(storageProviderSpec);
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, storageProviderSpec);

   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {

      StorageProviderSpec storageProviderSpec = publisherSpec
            .getPublishedEntitySpec(NimbusXVpProvider.DEFAULT_ENTITY);

      String url = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_URL);
      storageProviderSpec.providerUrl.set(url);

      String username = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_USERNAME);
      storageProviderSpec.username.set(username);

      String password = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_PASSWORD);
      storageProviderSpec.password.set(password);

      XVpServiceSpec serviceSpec = (XVpServiceSpec) storageProviderSpec.service
            .get();
      String endpoint = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_ENDPOINT);
      serviceSpec.endpoint.set(endpoint);

      String name = "XVP " +  endpoint;
      storageProviderSpec.name.set(name);

      logger.debug(String.format("Initialized Storage Provider spec '%s' ",
            storageProviderSpec));
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      logger.debug("No assembler spec to laod.");
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap)
         throws Exception {
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof XVpServiceSpec) {
            serviceConnectorsMap.put(serviceSpec, new XVpConnector(
                  (XVpServiceSpec) serviceSpec));
         }
      }
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return NimbusXVpProvider.class;
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      vpProvisionerSpec = filteredAssemblerSpec.links
            .get(XVpProvisionerSpec.class);

      if (isAssembling) {
         // Only used in assemble command - for other stages use the
         // settings loaded from the settings file.

         vpProvisionerSpec.vmName.set("XVP_VM_"
               + RandomStringUtils.randomNumeric(5));
         vpProvisionerSpec.version.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader, RESOURCE_KEY_VERSION));
         vpProvisionerSpec.replConfig.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader, RESOURCE_KEY_REPL_CONFIG));
         vpProvisionerSpec.url.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader, RESOURCE_KEY_URL));
      }

      // Set Nimbus service
      NimbusServiceSpec nimbusSpec = new NimbusServiceSpec();
      nimbusSpec.endpoint.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader, NimbusConfigKeys.NIMBUS_IP_KEY));

      nimbusSpec.username.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader, NimbusConfigKeys.NIMBUS_USER_KEY));

      nimbusSpec.password.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader, NimbusConfigKeys.NIMBUS_PASSWORD_KEY));

      nimbusSpec.deployUser.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader, NimbusConfigKeys.NIMBUS_DEPLOY_USER_KEY));

      nimbusSpec.pod.set(SettingsUtil.getRequiredValue(sessionSettingsReader,
            NimbusConfigKeys.NIMBUS_POD_KEY));

      vpProvisionerSpec.service.set(nimbusSpec);

   }

   @Override
   public String determineResourceVersion() throws Exception {
      return String.format("xVp.v.%s", vpProvisionerSpec.version.get());
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new NotImplementedException(
            "Retrieve xvp as resource is not yet implemeted");

   }

   @Override
   public int providerWeight() {
      return 0;
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter)
         throws Exception {
      logger.info("Start deployment of new VP...");

      NimbusCommandsUtil.deployXvp(vpProvisionerSpec);

      // Publisher info
      testbedSettingsWriter.setSetting(TESTBED_KEY_ENDPOINT,
            vpProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(TESTBED_KEY_URL,
            String.format(XVP_URL_TEMPLATE, vpProvisionerSpec.ip.get()));
      testbedSettingsWriter.setSetting(TESTBED_KEY_USERNAME, "username");
      testbedSettingsWriter.setSetting(TESTBED_KEY_PASSWORD, "password");

      logger.info("Saving test bed connection data...");

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      return true;
   }

   @Override
   public void destroyTestbed() throws Exception {
      logger.info("Delete nimbus VM...");
      NimbusCommandsUtil.destroyVM(vpProvisionerSpec);
   }

}
