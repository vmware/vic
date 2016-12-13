/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.util.AbstractMap.SimpleEntry;
import java.util.List;
import java.util.Map;

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
import com.vmware.vsphere.client.automation.provider.commontb.steps.AttachHostProviderStep;
import com.vmware.vsphere.client.automation.provider.commontb.steps.AttachStorageProviderStep;
import com.vmware.vsphere.client.automation.provider.commontb.steps.CreateClusterProviderStep;
import com.vmware.vsphere.client.automation.provider.commontb.steps.CreateDcProviderStep;
import com.vmware.vsphere.client.automation.provider.connector.VcConnector;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * Provider workflow for common test bed setup. It provides a VC with a
 * connected clustered host.
 * The cluster is DRS enabled and the host has a mounted NFS data store.
 *
 * -VC
 *   |---Admin user spec
 *   |---Data Center
 *      |--- Cluster
 *         |---Clustered Host
 *            |---NFS Shared Storage for the Host
 */
public class CommonTestBedProvider implements ProviderWorkflow {

   // Publisher info
   public static final String VC_ENTITY = "provider.ctb.entity.vc";
   public static final String NGC_ADMIN_USER_ENTITY = "provider.ctb.entity.ngc.admin.user";
   public static final String DC_ENTITY = "provider.ctb.entity.dc";
   public static final String CLUSTER_ENTITY = "provider.ctb.entity.cluster";
   public static final String CLUSTER_HOST_ENTITY = "provider.ctb.entity.cluster.host";
   public static final String CLUSTER_HOST_DS_ENTITY = "provider.ctb.entity.cluster.host.ds";
   public static final String CLUSTER_HOST_LOCAL_DS_ENTITY = "provider.ctb.entity.cluster.host.local.ds";

   // Testbed publisher settings

   // VC details
   private static final String TESTBED_KEY_NAME = "testbed.name";
   private static final String TESTBED_KEY_ENDPOINT = "testbed.endpoint";
   private static final String TESTBED_KEY_VSC_URL = "testbed.vsc.url";
   private static final String TESTBED_KEY_USERNAME = "testbed.user";
   private static final String TESTBED_KEY_PASSWORD = "testbed.pass";

   private static final String TESTBED_KEY_DC = "testbed.datacenter";
   private static final String TESTBED_KEY_CLUSTER = "testbed.cluster";
   private static final String TESTBED_KEY_HOST = "testbed.host";
   private static final String TESTBED_KEY_HOST_USER = "testbed.host.user";
   private static final String TESTBED_KEY_HOST_PASSWORD = "testbed.host.password";
   private static final String TESTBED_KEY_HOST_SERVICE_PORT =
         "testbed.host.service.port";
   private static final String TESTBED_KEY_HOST_LOCAL_DATASTORE_NAME = "testbed.host.datastore.name";
   private static final String TESTBED_KEY_DATASTORE = "testbed.datastore";
   private static final String TESTBED_KEY_DATASTORE_IP = "testbed.datastore.ip";
   private static final String TESTBED_KEY_DATASTORE_FOLDER = "testbed.datastore.folder";
   private static final String TESTBED_KEY_DATASTORE_TYPE = "testbed.datastore.type";

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(CommonTestBedProvider.class);

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

      DatastoreSpec localDsSpec = new DatastoreSpec();
      localDsSpec.parent.set(hostSpec);
      localDsSpec.type.set(DatastoreType.VMFS);
      publisherSpec.links.add(localDsSpec);
      publisherSpec.publishEntitySpec(CLUSTER_HOST_LOCAL_DS_ENTITY, localDsSpec);

      DatastoreSpec datastoreSpec = new DatastoreSpec();
      datastoreSpec.parent.set(hostSpec);
      datastoreSpec.type.set(DatastoreType.NFS);
      publisherSpec.links.add(datastoreSpec);
      publisherSpec.publishEntitySpec(CLUSTER_HOST_DS_ENTITY, datastoreSpec);
   }

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {

      // TODO: rkovachev for the Test workflow long keys - auto updates

      // Request VC spec
      TestbedSpecConsumer vcProviderConsumer =
            testbedBridge.requestTestbed(VcProvider.class, true);
      VcSpec requestedVcSpec =
            vcProviderConsumer.getPublishedEntitySpec(VcProvider.DEFAULT_ENTITY);

      // Init datacenter spec
      DatacenterSpec dcSpec = SpecFactory.getSpec(DatacenterSpec.class, requestedVcSpec);

      // Init cluster spec
      ClusterSpec clusterSpec = SpecFactory.getSpec(ClusterSpec.class, dcSpec);

      // Request host spec for the clustered host
      TestbedSpecConsumer hostProviderConsumer =
            testbedBridge.requestTestbed(HostProvider.class, false);
      HostSpec requestedHostSpec =
            hostProviderConsumer.getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      requestedHostSpec.parent.set(clusterSpec);

      // Request local Datastore spec for clustered host
      DatastoreSpec requestedLocalDatastoreSpec =
            hostProviderConsumer.getPublishedEntitySpec(HostProvider.LOCAL_DS_ENTITY);
      requestedLocalDatastoreSpec.parent.set(requestedHostSpec);
      requestedLocalDatastoreSpec.type.set(DatastoreType.VMFS);

      // Request Datastore spec for clustered host
      TestbedSpecConsumer datastoreProviderConsumer =
            testbedBridge.requestTestbed(NfsStorageProvider.class, true);
      DatastoreSpec requestedDatastoreSpec =
            datastoreProviderConsumer.getPublishedEntitySpec(NfsStorageProvider.DEFAULT_ENTITY);
      requestedDatastoreSpec.parent.set(requestedHostSpec);
      requestedDatastoreSpec.type.set(DatastoreType.NFS);

      // Link the specs to the assembler spec
      assemblerSpec.links.add(requestedVcSpec, requestedHostSpec, dcSpec, clusterSpec,
            requestedLocalDatastoreSpec, requestedDatastoreSpec);
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) {
      // SimulatorConnectorsFactory.createAndSetConnectors(serviceConnectorsMap);
      // TODO: rkovachev - move to ConnectionFactory like the simulator sample
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
   public void composeProviderSteps(WorkflowStepsSequence<? extends WorkflowStepContext> flow) throws Exception {
      // save the settings for the base components of the TB - common ESX and VC
      // in the current case.
      flow.appendStep("Load settings from the base(requested0 testbeds", new ProviderWorkflowStep() {

         private VcSpec _vcSpec;
         private DatacenterSpec _dcSpec;
         private ClusterSpec _clusterSpec;
         private HostSpec _hostSpec;
         private DatastoreSpec _localVmfsDatastoreSpec;
         private DatastoreSpec _datastoreSpec;

         // private HostSpec _hostSpec;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filterAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) throws Exception {

            // collect data to be saved if in assemble mode
            if (isAssembling) {
               List<DatastoreSpec> datastores = filterAssemblerSpec.links.getAll(DatastoreSpec.class);
               for (DatastoreSpec dsSpec : datastores) {
                  if (dsSpec.type.get().equals(DatastoreType.VMFS)) {
                     _localVmfsDatastoreSpec = dsSpec;
                  } else if (dsSpec.type.get().equals(DatastoreType.NFS)) {
                     _datastoreSpec = dsSpec;
                  }
               }
               _vcSpec = filterAssemblerSpec.links.get(VcSpec.class);
               _dcSpec = filterAssemblerSpec.links.get(DatacenterSpec.class);
               _clusterSpec = filterAssemblerSpec.links.get(ClusterSpec.class);
               // _hostSpec = _vcSpec.links.get(HostSpec.class);
               _hostSpec = filterAssemblerSpec.links.get(HostSpec.class);
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
         public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
            String vcHostName = VcServiceUtil.getVcHostname(_vcSpec.service.get());

            testbedSettingsWriter.setSetting(TESTBED_KEY_NAME, vcHostName);
            testbedSettingsWriter.setSetting(TESTBED_KEY_ENDPOINT, _vcSpec.service.get().endpoint.get());
            testbedSettingsWriter.setSetting(TESTBED_KEY_VSC_URL, _vcSpec.vscUrl.get());
            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_USERNAME,
                  _vcSpec.ssoLoginUsername.get());
            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_PASSWORD,
                  _vcSpec.ssoLoginPassword.get());

            testbedSettingsWriter.setSetting(TESTBED_KEY_DC, _dcSpec.name.get());

            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_CLUSTER,
                  _clusterSpec.name.get());

            testbedSettingsWriter.setSetting(TESTBED_KEY_HOST, _hostSpec.name.get());
            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_HOST_USER,
                  _hostSpec.userName.get());
            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_HOST_PASSWORD,
                  _hostSpec.password.get());
            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_HOST_SERVICE_PORT,
                  _hostSpec.port.get().toString());
            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_HOST_LOCAL_DATASTORE_NAME,
                  _localVmfsDatastoreSpec.name.get());

            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_DATASTORE,
                  _datastoreSpec.name.get());

            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_DATASTORE_IP,
                  _datastoreSpec.remoteHost.get());

            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_DATASTORE_FOLDER,
                  _datastoreSpec.remotePath.get());

            testbedSettingsWriter.setSetting(
                  TESTBED_KEY_DATASTORE_TYPE,
                  _datastoreSpec.type.get().toString());

         }

      });

      // Create Data Center in the Common Testbed
      flow.appendStep("Create Data Center in Common Testbed", new CreateDcProviderStep());

      // Create Cluster in the Common Testbed
      flow.appendStep(
            "Create Cluster in the Common Testbed",
            new CreateClusterProviderStep());

      // Attach host to cluster in the Common Testbed
      flow.appendStep(
            "Attach host to cluster in the Common Testbed",
            new AttachHostProviderStep());

      // Attach all available datastores
      flow.appendStep(
            "Attach datastores to hosts in the Common Testbed",
            new AttachStorageProviderStep());

   }

   @Override
   public void assignTestbedSettings(PublisherSpec providerSpec,
         SettingsReader testbedSettings) throws Exception {

      _logger.info("Common Test Bed assignTestbed started");

      // Load published VC settings
      VcSpec vcSpec = providerSpec.getPublishedEntitySpec(VC_ENTITY);

      vcSpec.vscUrl.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_VSC_URL));

      String name =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_NAME);
      String endpoint =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_ENDPOINT);
      String username =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_USERNAME);
      String password =
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_PASSWORD);

      vcSpec.ssoLoginUsername.set(username);
      vcSpec.ssoLoginPassword.set(password);

      VcServiceSpec commonVcServiceSpec = new VcServiceSpec();
      commonVcServiceSpec.endpoint.set(endpoint);
      commonVcServiceSpec.username.set(username);
      commonVcServiceSpec.password.set(password);

      vcSpec.service.set(commonVcServiceSpec);

      vcSpec.name.set(name);

      // Load published VC user settings
      UserSpec ngcUserSpec = providerSpec.links.get(UserSpec.class);
      ngcUserSpec.username.set(username);
      ngcUserSpec.password.set(password);
      ngcUserSpec.service.set(commonVcServiceSpec);

      // Load published DC settings
      DatacenterSpec datacenterSpec = providerSpec.links.get(DatacenterSpec.class);
      datacenterSpec.service.set(commonVcServiceSpec);
      datacenterSpec.name.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_DC));

      // Load published cluster settings
      ClusterSpec clusterSpec = providerSpec.links.get(ClusterSpec.class);
      clusterSpec.service.set(commonVcServiceSpec);
      clusterSpec.name.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_CLUSTER));

      // Load published clustered host settings
      HostSpec hostSpec = providerSpec.links.get(HostSpec.class);
      hostSpec.service.set(commonVcServiceSpec);
      hostSpec.name.set(
            SettingsUtil.getRequiredValue(testbedSettings, TESTBED_KEY_HOST));
      hostSpec.userName.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_HOST_USER));
      hostSpec.password.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_HOST_PASSWORD));

      // Load published shared data store settings
      List<DatastoreSpec> datastores = providerSpec.links.getAll(DatastoreSpec.class);
      for (DatastoreSpec datastoreSpec : datastores) {
         if (datastoreSpec.type.get().equals(DatastoreType.VMFS)) {
            datastoreSpec.name.set(SettingsUtil.getRequiredValue(
                  testbedSettings,
                  TESTBED_KEY_HOST_LOCAL_DATASTORE_NAME));
         } else if (datastoreSpec.type.get().equals(DatastoreType.NFS)) {
            datastoreSpec.name.set(SettingsUtil.getRequiredValue(
                  testbedSettings,
                  TESTBED_KEY_DATASTORE));
            datastoreSpec.remoteHost.set(SettingsUtil.getRequiredValue(
                  testbedSettings,
                  TESTBED_KEY_DATASTORE_IP));
            datastoreSpec.remotePath.set(SettingsUtil.getRequiredValue(
                  testbedSettings,
                  TESTBED_KEY_DATASTORE_FOLDER));
         }
         datastoreSpec.service.set(commonVcServiceSpec);
      }

      _logger.debug("Loaded CTB published specs list:");
      for (SimpleEntry<String, EntitySpec> simpleEntry : providerSpec.entitySpecMap.getAll()) {
         _logger.debug("Entity key:" + simpleEntry.getKey() + " " + simpleEntry.getValue().service.get().endpoint);
      }
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
      // Nothing to do here
      _logger.debug("No CTB assembler spec to laod.");
   }

   @Override
   public int providerWeight() {
      // TODO Auto-generated method stub
      return 0;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return this.getClass();
   }
}
