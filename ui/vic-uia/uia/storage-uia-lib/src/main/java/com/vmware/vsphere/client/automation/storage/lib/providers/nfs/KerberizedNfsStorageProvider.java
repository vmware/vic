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
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.KerberosRealmSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec.Nfs41AuthenticationMode;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.NfsDatastoreSpec.DatastoreAccessMode;

public class KerberizedNfsStorageProvider extends BaseElementalProvider {

   public static final String KRB5_DATASTORE_ENTITY = "provider.nfs.kerberos.storage.entity.authenticationmode.krb5";
   public static final String KRB5_KERBEROS_REALM = "provider.nfs.kerberos.storage.entity.krbrealm";
   public static final String KRB5I_DATASTORE_ENTITY = "provider.nfs.kerberos.storage.entity.authenticationmode.krb5i";
   public static final String KRB5I_KERBEROS_REALM = "provider.nfs.kerberos.storage.entity.krb5irealm";

   // Testbed setting keys
   public static final String TESTBED_KEY_KRB5_DATASTORE_IP = "testbed.nfs.kerberos.storage.krb5.ip";
   public static final String TESTBED_KEY_KRB5_DATASTORE_FOLDER = "testbed.nfs.kerberos.storage.krb5.folder";
   public static final String TESTBED_KEY_KRB5_DATASTORE_NAME = "testbed.nfs.kerberos.storage.krb5.name";
   public static final String TESTBED_KEY_KRB5_DATASTORE_ACCESSMODE = "testbed.nfs.kerberos.storage.krb5.accessmode";
   public static final String TESTBED_KEY_KRB5I_DATASTORE_IP = "testbed.nfs.kerberos.storage.krb5i.ip";
   public static final String TESTBED_KEY_KRB5I_DATASTORE_FOLDER = "testbed.nfs.kerberos.storage.krb5i.folder";
   public static final String TESTBED_KEY_KRB5I_DATASTORE_NAME = "testbed.nfs.kerberos.storage.krb5i.name";
   public static final String TESTBED_KEY_KRB5I_DATASTORE_ACCESSMODE = "testbed.nfs.kerberos.storage.krb5i.accessmode";

   public static final String TESTBED_KEY_KERBEROS_DNS_SERVER = "testbed.nfs.kerberos.activedirectory.dns.ip";
   public static final String TESTBED_KEY_ACTIVE_DIRECTORY_SERVER = "testbed.nfs.kerberos.activedirectory.server";
   public static final String TESTBED_KEY_ACTIVE_DIRECTORY_DOMAIN = "testbed.nfs.kerberos.activedirectory.domain";
   public static final String TESTBED_KEY_ACTIVE_DIRECTORY_USERNAME = "testbed.nfs.kerberos.activedirectory.username";
   public static final String TESTBED_KEY_ACTIVE_DIRECTORY_PASSWORD = "testbed.nfs.kerberos.activedirectory.password";
   public static final String TESTBED_KEY_ACTIVE_DIRECTORY_NFS_USERNAME = "testbed.nfs.kerberos.activedirectory.nfs41.username";
   public static final String TESTBED_KEY_ACTIVE_DIRECTORY_NFS_PASSWORD = "testbed.nfs.kerberos.activedirectory.nfs41.password";

   private static final Logger _logger = LoggerFactory
         .getLogger(KerberizedNfsStorageProvider.class);

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      Nfs41DatastoreSpec datastoreSpec = new Nfs41DatastoreSpec();
      datastoreSpec.service.set(new NfsStorageServiceSpec());
      publisherSpec.links.add(datastoreSpec);
      KerberosRealmSpec krbRealm = new KerberosRealmSpec();
      publisherSpec.links.add(krbRealm);

      publisherSpec.publishEntitySpec(KRB5_DATASTORE_ENTITY, datastoreSpec);
      publisherSpec.publishEntitySpec(KRB5_KERBEROS_REALM, krbRealm);

      Nfs41DatastoreSpec datastore5iSpec = new Nfs41DatastoreSpec();
      datastore5iSpec.service.set(new NfsStorageServiceSpec());
      publisherSpec.links.add(datastore5iSpec);
      KerberosRealmSpec krb5iRealm = new KerberosRealmSpec();
      publisherSpec.links.add(krb5iRealm);

      publisherSpec.publishEntitySpec(KRB5I_DATASTORE_ENTITY, datastore5iSpec);
      publisherSpec.publishEntitySpec(KRB5I_KERBEROS_REALM, krb5iRealm);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {
      Nfs41DatastoreSpec datastoreSpec = publisherSpec
            .getPublishedEntitySpec(KerberizedNfsStorageProvider.KRB5_DATASTORE_ENTITY);
      NfsStorageServiceSpec serviceSpec = (NfsStorageServiceSpec) datastoreSpec.service
            .get();
      datastoreSpec.authenticationMode.set(Nfs41AuthenticationMode.KRB5);
      String ip = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5_DATASTORE_IP);
      datastoreSpec.remoteHost.set(ip);
      serviceSpec.endpoint.set(ip);

      String folder = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5_DATASTORE_FOLDER);
      datastoreSpec.remotePath.set(folder);
      String name = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5_DATASTORE_NAME);
      datastoreSpec.name.set(name);

      String accessMode = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5_DATASTORE_ACCESSMODE);
      datastoreSpec.accessMode.set(DatastoreAccessMode.valueOf(accessMode));

      KerberosRealmSpec kerberosRealmSpec = publisherSpec
            .getPublishedEntitySpec(KerberizedNfsStorageProvider.KRB5_KERBEROS_REALM);
      initializeKerberosRealmSpec(testbedSettings, kerberosRealmSpec);

      _logger
            .debug(String.format(
                  "Initialized KerberosNfsStorageProvider spec '%s' ",
                  datastoreSpec));

      Nfs41DatastoreSpec krb5iDatastoreSpec = publisherSpec
            .getPublishedEntitySpec(KerberizedNfsStorageProvider.KRB5I_DATASTORE_ENTITY);
      NfsStorageServiceSpec krb5iServiceSpec = (NfsStorageServiceSpec) krb5iDatastoreSpec.service
            .get();
      krb5iDatastoreSpec.authenticationMode.set(Nfs41AuthenticationMode.KRB5I);

      String krb5iIp = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5I_DATASTORE_IP);
      krb5iDatastoreSpec.remoteHost.set(krb5iIp);
      krb5iServiceSpec.endpoint.set(krb5iIp);

      String krb5iFolder = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5I_DATASTORE_FOLDER);
      krb5iDatastoreSpec.remotePath.set(krb5iFolder);

      String krb5iName = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5I_DATASTORE_NAME);
      krb5iDatastoreSpec.name.set(krb5iName);

      String krb5iAccessMode = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KRB5I_DATASTORE_ACCESSMODE);
      krb5iDatastoreSpec.accessMode.set(DatastoreAccessMode
            .valueOf(krb5iAccessMode));

      KerberosRealmSpec krb5iKerberosRealmSpec = publisherSpec
            .getPublishedEntitySpec(KerberizedNfsStorageProvider.KRB5I_KERBEROS_REALM);
      initializeKerberosRealmSpec(testbedSettings, krb5iKerberosRealmSpec);

      _logger.debug(String.format(
            "Initialized KerberosNfsStorageProvider spec '%s' ",
            krb5iDatastoreSpec));

   }

   private void initializeKerberosRealmSpec(SettingsReader testbedSettings,
         KerberosRealmSpec kerberosRealmSpec) {

      String dnsServer = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_KERBEROS_DNS_SERVER);
      kerberosRealmSpec.dnsServer.set(dnsServer);

      String activeDirectoryServer = SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_ACTIVE_DIRECTORY_SERVER);
      kerberosRealmSpec.activeDirectoryServer.set(activeDirectoryServer);

      String activeDirectoryDomain = SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_ACTIVE_DIRECTORY_DOMAIN);
      kerberosRealmSpec.activeDirectoryDomain.set(activeDirectoryDomain);

      String activeDirectoryUserName = SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_ACTIVE_DIRECTORY_USERNAME);
      kerberosRealmSpec.activeDirectoryUserName.set(activeDirectoryUserName);

      String activeDirectoryPassword = SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_ACTIVE_DIRECTORY_PASSWORD);
      kerberosRealmSpec.activeDirectoryPassword.set(activeDirectoryPassword);

      String storageUserName = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_ACTIVE_DIRECTORY_NFS_USERNAME);
      kerberosRealmSpec.storageUserName.set(storageUserName);

      String storagePassword = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_ACTIVE_DIRECTORY_NFS_PASSWORD);
      kerberosRealmSpec.storagePassword.set(storagePassword);
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
