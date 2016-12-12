package com.vmware.vsphere.client.automation.storage.lib.providers.spbm;

import java.util.AbstractMap.SimpleEntry;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.provider.commontb.NimbusXVpProvider;
import com.vmware.vsphere.client.automation.provider.connector.VcConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StorageProviderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * Provider workflow for Storage Policy Based Management test bed setup. It
 * provides Common Test Bed settings (VC with a connected clustered host. The
 * cluster is DRS enabled and the host has a mounted NFS data store.), xVp is
 * attached to the VC and iSCSI adapter is added to the clustered host
 *
 * -VC |---Admin user |---Data Center |--- Cluster |---Clustered Host |---iSCSI
 * Adapter added to the Host |---NFS Shared Storage for the Host |---xVP
 */
public class SPBMTestBedProvider implements ProviderWorkflow {

   // Publisher info
   public static final String VC_ENTITY = "provider.spbm.entity.vc";
   public static final String NGC_ADMIN_USER_ENTITY = "provider.spbm.entity.ngc.admin.user";
   public static final String DC_ENTITY = "provider.spbm.entity.dc";
   public static final String CLUSTER_ENTITY = "provider.spbm.entity.cluster";
   public static final String CLUSTER_HOST_ENTITY = "provider.spbm.entity.cluster.host";
   public static final String CLUSTER_HOST_DS_ENTITY = "provider.spbm.entity.cluster.host.ds";
   public static final String VP_ENTITY = "provider.spbm.entity.vp";
   public static final String STANDALONE_HOST_ENTITY = "provider.spbm.entity.standalone.host";
   public static final String CLUSTER_HOST_LOCAL_DS_ENTITY = "provider.spbm.entity.cluster.host.local.ds";

   // testbed key details
   // vc
   private static final String TESTBED_KEY_NAME = "testbed.name";
   private static final String TESTBED_KEY_ENDPOINT = "testbed.endpoint";
   private static final String TESTBED_KEY_VSC_URL = "testbed.vsc.url";
   private static final String TESTBED_KEY_USERNAME = "testbed.user";
   private static final String TESTBED_KEY_PASSWORD = "testbed.pass";
   // datacenter
   private static final String TESTBED_KEY_DC = "testbed.datacenter";
   // cluster
   private static final String TESTBED_KEY_CLUSTER = "testbed.cluster";
   // clustered host
   private static final String TESTBED_KEY_HOST = "testbed.clustered.host";
   private static final String TESTBED_KEY_HOST_USER = "testbed.clustered.host.user";
   private static final String TESTBED_KEY_HOST_PASSWORD = "testbed.clustered.host.password";
   // nfs datastore
   private static final String TESTBED_KEY_DATASTORE = "testbed.datastore";
   private static final String TESTBED_KEY_DATASTORE_IP = "testbed.datastore.ip";
   private static final String TESTBED_KEY_DATASTORE_FOLDER = "testbed.datastore.folder";
   private static final String TESTBED_KEY_DATASTORE_TYPE = "testbed.datastore.type";
   // local datastore
   private static final String TESTBED_KEY_HOST_LOCAL_DATASTORE_NAME = "testbed.host.datastore.name";
   // VP
   private static final String TESTBED_KEY_VP = "testbed.vp";
   private static final String TESTBED_KEY_VP_URL = "testbed.vp.url";
   private static final String TESTBED_KEY_VP_USERNAME = "testbed.vp.user";
   private static final String TESTBED_KEY_VP_PASSWORD = "testbed.vp.pass";

   // logger
   private static final Logger _logger = LoggerFactory
         .getLogger(SPBMTestBedProvider.class);

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      VcSpec vcSpec = new VcSpec();
      publisherSpec.links.add(vcSpec);
      publisherSpec.publishEntitySpec(VC_ENTITY, vcSpec);

      UserSpec ngcAdminUserSpec = new UserSpec();
      ngcAdminUserSpec.parent.set(vcSpec);
      publisherSpec.links.add(ngcAdminUserSpec);
      publisherSpec.publishEntitySpec(NGC_ADMIN_USER_ENTITY, ngcAdminUserSpec);

      DatacenterSpec datacenterSpec = new DatacenterSpec();
      datacenterSpec.parent.set(vcSpec);
      publisherSpec.links.add(datacenterSpec);
      publisherSpec.publishEntitySpec(DC_ENTITY, datacenterSpec);

      ClusterSpec clusterSpec = new ClusterSpec();
      clusterSpec.parent.set(datacenterSpec);
      publisherSpec.links.add(clusterSpec);
      publisherSpec.publishEntitySpec(CLUSTER_ENTITY, clusterSpec);

      HostSpec hostSpec = new HostSpec();
      hostSpec.parent.set(clusterSpec);
      publisherSpec.links.add(hostSpec);
      publisherSpec.publishEntitySpec(CLUSTER_HOST_ENTITY, hostSpec);

      DatastoreSpec datastoreSpec = new DatastoreSpec();
      datastoreSpec.parent.set(hostSpec);
      datastoreSpec.type.set(DatastoreType.NFS);
      publisherSpec.links.add(datastoreSpec);
      publisherSpec.publishEntitySpec(CLUSTER_HOST_DS_ENTITY, datastoreSpec);

      DatastoreSpec localDsSpec = new DatastoreSpec();
      localDsSpec.parent.set(hostSpec);
      localDsSpec.type.set(DatastoreType.VMFS);
      publisherSpec.links.add(localDsSpec);
      publisherSpec
            .publishEntitySpec(CLUSTER_HOST_LOCAL_DS_ENTITY, localDsSpec);

      StorageProviderSpec xVpSpec = new StorageProviderSpec();
      xVpSpec.parent.set(vcSpec);
      publisherSpec.links.add(xVpSpec);
      publisherSpec.publishEntitySpec(VP_ENTITY, xVpSpec);

   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec,
         TestBedBridge testbedBridge) throws Exception {

      TestbedSpecConsumer ctbProviderConsumer = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);

      // Request VC spec
      VcSpec requestedVcSpec = ctbProviderConsumer
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Request datacenter spec
      DatacenterSpec requestedDcSpec = ctbProviderConsumer
            .getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);
      requestedDcSpec.parent.set(requestedVcSpec);

      // Request cluster spec
      ClusterSpec clusterSpec = ctbProviderConsumer
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_ENTITY);
      clusterSpec.parent.set(requestedDcSpec);

      // Request host spec for the clustered host
      HostSpec requestedHostSpec = ctbProviderConsumer
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_ENTITY);
      requestedHostSpec.parent.set(clusterSpec);

      // Request VP spec
      TestbedSpecConsumer xVpProviderConsumer = testbedBridge.requestTestbed(
            NimbusXVpProvider.class, true);
      StorageProviderSpec requestedXVpSpec = xVpProviderConsumer
            .getPublishedEntitySpec(NimbusXVpProvider.DEFAULT_ENTITY);
      requestedXVpSpec.parent.set(requestedVcSpec);

      // Request NFS Datastore spec for clustered host
      DatastoreSpec requestedDatastoreSpec = ctbProviderConsumer
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_DS_ENTITY);
      requestedDatastoreSpec.parent.set(requestedHostSpec);
      requestedDatastoreSpec.type.set(DatastoreType.NFS);

      // Request local Datastore spec for clustered host
      DatastoreSpec requestedLocalDatastoreSpec = ctbProviderConsumer
            .getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_HOST_LOCAL_DS_ENTITY);
      requestedLocalDatastoreSpec.parent.set(requestedHostSpec);
      requestedLocalDatastoreSpec.type.set(DatastoreType.VMFS);

      assemblerSpec.links.add(requestedVcSpec, requestedDcSpec, clusterSpec,
            requestedHostSpec, requestedXVpSpec, requestedDatastoreSpec,
            requestedLocalDatastoreSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {
      _logger.info("SPBM Test Bed assignTestbed started");

      // Load published VC settings
      VcSpec vcSpec = publisherSpec.getPublishedEntitySpec(VC_ENTITY);

      vcSpec.vscUrl.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_VSC_URL));

      String name = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_NAME);
      String endpoint = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_ENDPOINT);
      String username = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_USERNAME);
      String password = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_PASSWORD);

      vcSpec.ssoLoginUsername.set(username);
      vcSpec.ssoLoginPassword.set(password);

      VcServiceSpec commonVcServiceSpec = new VcServiceSpec();
      commonVcServiceSpec.endpoint.set(endpoint);
      commonVcServiceSpec.username.set(username);
      commonVcServiceSpec.password.set(password);

      vcSpec.service.set(commonVcServiceSpec);
      vcSpec.name.set(name);

      // Load published VC user settings
      UserSpec ngcUserSpec = publisherSpec.links.get(UserSpec.class);
      ngcUserSpec.username.set(username);
      ngcUserSpec.password.set(password);
      ngcUserSpec.service.set(commonVcServiceSpec);

      // Load published DC settings
      DatacenterSpec datacenterSpec = publisherSpec.links
            .get(DatacenterSpec.class);
      datacenterSpec.service.set(commonVcServiceSpec);
      datacenterSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_DC));

      // Load published cluster settings
      ClusterSpec clusterSpec = publisherSpec.links.get(ClusterSpec.class);
      clusterSpec.service.set(commonVcServiceSpec);
      clusterSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_CLUSTER));

      // Load published clustered host settings
      HostSpec clusteredHostSpec = publisherSpec.links.get(HostSpec.class);
      clusteredHostSpec.service.set(commonVcServiceSpec);
      clusteredHostSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_HOST));
      clusteredHostSpec.userName.set(SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_HOST_USER));
      clusteredHostSpec.password.set(SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_HOST_PASSWORD));

      // Load published data stores settings
      List<DatastoreSpec> datastores = publisherSpec.links
            .getAll(DatastoreSpec.class);
      for (DatastoreSpec datastoreSpec : datastores) {
         if (datastoreSpec.type.get().equals(DatastoreType.VMFS)) {
            datastoreSpec.name.set(SettingsUtil.getRequiredValue(
                  testbedSettings, TESTBED_KEY_HOST_LOCAL_DATASTORE_NAME));
         } else if (datastoreSpec.type.get().equals(DatastoreType.NFS)) {
            datastoreSpec.name.set(SettingsUtil.getRequiredValue(
                  testbedSettings, TESTBED_KEY_DATASTORE));
            datastoreSpec.remoteHost.set(SettingsUtil.getRequiredValue(
                  testbedSettings, TESTBED_KEY_DATASTORE_IP));
            datastoreSpec.remotePath.set(SettingsUtil.getRequiredValue(
                  testbedSettings, TESTBED_KEY_DATASTORE_FOLDER));
         }
         datastoreSpec.service.set(commonVcServiceSpec);
      }

      // Load published xVp settings
      StorageProviderSpec xVpSpec = publisherSpec.links
            .get(StorageProviderSpec.class);
      xVpSpec.name.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_VP));
      xVpSpec.providerUrl.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_VP_URL));
      xVpSpec.username.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_VP_USERNAME));
      xVpSpec.password.set(SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_VP_PASSWORD));
      xVpSpec.service.set(commonVcServiceSpec);

      // set host iscsi server ip
      clusteredHostSpec.iscsiServerIp.set(StringUtils.substringBetween(
            xVpSpec.providerUrl.get(), "https://", ":8443/vasa/version.xml"));

      _logger.debug("Loaded CTB published specs list:");
      for (SimpleEntry<String, EntitySpec> simpleEntry : publisherSpec.entitySpecMap
            .getAll()) {
         _logger.debug("Entity key:" + simpleEntry.getKey() + " "
               + simpleEntry.getValue().service.get().endpoint);
      }

   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      // Nothing to do here
      _logger.debug("No assembler spec to load.");

   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) {
      _logger.debug("Virtual Center Provider assignTestbedConnectors");
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         if (serviceSpec instanceof VcServiceSpec) {
            try {
               serviceConnectorsMap.put(serviceSpec, new VcConnector(
                     (VcServiceSpec) serviceSpec));
            } catch (SsoException e) {
               _logger.error(e.getMessage());
               throw new RuntimeException("SSO authentication failed", e);
            }
         }
      }
   }

   @Override
   public void composeProviderSteps(
         WorkflowStepsSequence<? extends WorkflowStepContext> flow)
         throws Exception {
      flow.appendStep("Load settings from the base(requested testbeds",
            new ProviderWorkflowStep() {

               private VcSpec _vcSpec;
               private DatacenterSpec _dcSpec;
               private ClusterSpec _clusterSpec;
               private HostSpec _clusteredHostSpec;
               private DatastoreSpec _datastoreSpec;
               private DatastoreSpec _localVmfsDatastoreSpec;
               private StorageProviderSpec _xVpSpec;

               @Override
               public void prepare(PublisherSpec filteredPublisherSpec,
                     AssemblerSpec filterAssemblerSpec, boolean isAssembling,
                     SettingsReader sessionSettingsReader) throws Exception {

                  // collect data to be saved if in assemble mode
                  if (isAssembling) {
                     _vcSpec = filterAssemblerSpec.links.get(VcSpec.class);
                     _dcSpec = filterAssemblerSpec.links
                           .get(DatacenterSpec.class);
                     _clusterSpec = filterAssemblerSpec.links
                           .get(ClusterSpec.class);
                     _clusteredHostSpec = filterAssemblerSpec.links
                           .get(HostSpec.class);
                     List<DatastoreSpec> datastores = filterAssemblerSpec.links
                           .getAll(DatastoreSpec.class);
                     for (DatastoreSpec dsSpec : datastores) {
                        if (dsSpec.type.get().equals(DatastoreType.VMFS)) {
                           _localVmfsDatastoreSpec = dsSpec;
                        } else if (dsSpec.type.get().equals(DatastoreType.NFS)) {
                           _datastoreSpec = dsSpec;
                        }
                     }

                     _xVpSpec = filterAssemblerSpec.links
                           .get(StorageProviderSpec.class);

                  }
               }

               @Override
               public void disassemble() throws Exception {
                  // Nothing to disassemble
               }

               @Override
               public boolean checkHealth() throws Exception {
                  // Nothing to check
                  return true;
               }

               @Override
               public void assemble(SettingsWriter testbedSettingsWriter)
                     throws Exception {

                  String vcHostName = VcServiceUtil
                        .getVcHostname(_vcSpec.service.get());

                  testbedSettingsWriter
                        .setSetting(TESTBED_KEY_NAME, vcHostName);
                  testbedSettingsWriter.setSetting(TESTBED_KEY_ENDPOINT,
                        _vcSpec.service.get().endpoint.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_VSC_URL,
                        _vcSpec.vscUrl.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_USERNAME,
                        _vcSpec.ssoLoginUsername.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_PASSWORD,
                        _vcSpec.ssoLoginPassword.get());

                  testbedSettingsWriter.setSetting(TESTBED_KEY_DC,
                        _dcSpec.name.get());

                  testbedSettingsWriter.setSetting(TESTBED_KEY_CLUSTER,
                        _clusterSpec.name.get());

                  testbedSettingsWriter.setSetting(TESTBED_KEY_HOST,
                        _clusteredHostSpec.name.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_HOST_USER,
                        _clusteredHostSpec.userName.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_HOST_PASSWORD,
                        _clusteredHostSpec.password.get());

                  testbedSettingsWriter.setSetting(
                        TESTBED_KEY_HOST_LOCAL_DATASTORE_NAME,
                        _localVmfsDatastoreSpec.name.get());

                  testbedSettingsWriter.setSetting(TESTBED_KEY_DATASTORE,
                        _datastoreSpec.name.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_DATASTORE_IP,
                        _datastoreSpec.remoteHost.get());
                  testbedSettingsWriter.setSetting(
                        TESTBED_KEY_DATASTORE_FOLDER,
                        _datastoreSpec.remotePath.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_DATASTORE_TYPE,
                        _datastoreSpec.type.get().toString());

                  testbedSettingsWriter.setSetting(TESTBED_KEY_VP_URL,
                        _xVpSpec.providerUrl.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_VP_USERNAME,
                        _xVpSpec.username.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_VP_PASSWORD,
                        _xVpSpec.password.get());
                  testbedSettingsWriter.setSetting(TESTBED_KEY_VP,
                        _xVpSpec.name.get());
               }

            });

      // Set iscsiServerIp to host
      flow.appendStep("Add iscsiServerIp(xVp ip) to hosts in the SPBM Testbed",
            new ProviderWorkflowStep() {

               private HostSpec _host;
               private StorageProviderSpec _storageProvider;

               @Override
               public void prepare(PublisherSpec filteredPublisherSpec,
                     AssemblerSpec filterAssemblerSpec, boolean isAssembling,
                     SettingsReader sessionSettingsReader) throws Exception {
                  if (isAssembling) {
                     _host = filterAssemblerSpec.links.get(HostSpec.class);
                     _storageProvider = filterAssemblerSpec.links
                           .get(StorageProviderSpec.class);
                  } else {
                     _host = filteredPublisherSpec.links.get(HostSpec.class);
                     _storageProvider = filteredPublisherSpec.links
                           .get(StorageProviderSpec.class);
                  }

               }

               @Override
               public void assemble(SettingsWriter testbedSettingsWriter)
                     throws Exception {
                  String iscsiServerIp = StringUtils.substringBetween(
                        _storageProvider.providerUrl.get(), "https://",
                        ":8443/vasa/version.xml");
                  _host.iscsiServerIp.set(iscsiServerIp);
               }

               @Override
               public boolean checkHealth() throws Exception {
                  return true;
               }

               @Override
               public void disassemble() throws Exception {
                  // do nothing
               }
            });

      // Attach ISCSI in the Testbed
      flow.appendStep("Attach ISCSI to hosts in the SPBM Testbed",
            new AttachIscsiTargetStep());

      // Attach VP in the Testbed
      flow.appendStep("Attach VASA provider in the SPBM Testbed",
            new AttachVasaProviderStep());

   }

   @Override
   public int providerWeight() {
      return 0;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return this.getClass();
   }

}
