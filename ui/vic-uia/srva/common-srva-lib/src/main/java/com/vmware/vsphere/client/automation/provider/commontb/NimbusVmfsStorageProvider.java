/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;

import org.apache.commons.lang.RandomStringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusServiceSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.VmfsStorageProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.util.NimbusCommandsUtil;

/**
 * Provider for shared VMFS storage deployed on a Nimbus Vm with configured VMFS iSCSI server
 */
public class NimbusVmfsStorageProvider extends VmfsStorageProvider {

   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME = "testbed.assembler.nimbus.vm.name";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER = "testbed.assembler.nimbus.vm.owner";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD = "testbed.assembler.nimbus.vm.pod";

   private static final String TESTBED_ASSEMBLER_KEY_SSH_IP = "testbed.ssh.ip";
   private static final String TESTBED_ASSEMBLER_KEY_SSH_IPV4 = "testbed.ssh.ipv4";
   private static final String TESTBED_ASSEMBLER_KEY_SSH_IPV6 = "testbed.ssh.ipv6";

   private VmfsStorageProvisionerSpec _vmfsStorageProvisionerSpec;

   private static final Logger _logger = LoggerFactory.getLogger(NimbusVmfsStorageProvider.class);

   private static final int VMFS_SERVER_TCP_PORT = 3260;
   private static final int TCP_CONNECTION_TIMEOUT = 10000;
   private static final String NIMBUS_VMFS_STORAGE_PREFIX = "VMFS_STORAGE_";

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge) throws Exception {
      _vmfsStorageProvisionerSpec = new VmfsStorageProvisionerSpec();
      assemblerSpec.links.add(_vmfsStorageProvisionerSpec);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec, SettingsReader testbedSettings) throws Exception {

      // _vmfsStorageProvisionerSpec is added in the init stage but add the line for consistency
      _vmfsStorageProvisionerSpec = assemblerSpec.links.get(VmfsStorageProvisionerSpec.class);

      // NOTE: Nimbus service is set in the prepare stage - check preapreForOperations

      _vmfsStorageProvisionerSpec.vmName.set(SettingsUtil.getRequiredValue(
         testbedSettings,
         TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME));
      _vmfsStorageProvisionerSpec.ip.set(SettingsUtil.getRequiredValue(testbedSettings, TESTBED_ASSEMBLER_KEY_SSH_IP));
      _vmfsStorageProvisionerSpec.ipv4.set(SettingsUtil.getRequiredValue(
         testbedSettings,
         TESTBED_ASSEMBLER_KEY_SSH_IPV4));
      _vmfsStorageProvisionerSpec.ipv6.set(SettingsUtil.getRequiredValue(
         testbedSettings,
         TESTBED_ASSEMBLER_KEY_SSH_IPV6));
      _vmfsStorageProvisionerSpec.vmOwner.set(SettingsUtil.getRequiredValue(
         testbedSettings,
         TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER));
      _vmfsStorageProvisionerSpec.pod.set(SettingsUtil.getRequiredValue(
         testbedSettings,
         TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD));

      _logger.info("Loaded assemblerSpec: " + _vmfsStorageProvisionerSpec);
   }


   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpecSpec, AssemblerSpec filteredAssemblerSpec,
      boolean isAssembling, SettingsReader sessionSettingsReader) {

      _vmfsStorageProvisionerSpec = filteredAssemblerSpec.links.get(VmfsStorageProvisionerSpec.class);

      if (isAssembling) {
         final String name = NIMBUS_VMFS_STORAGE_PREFIX + RandomStringUtils.randomNumeric(5);
         // Only used in assemble command - for other stages use the
         // settings loaded from the settings file.
         _vmfsStorageProvisionerSpec.vmName.set(name);
      }

      // Set Nimbus service
      NimbusServiceSpec nimbusSpec = new NimbusServiceSpec();
      nimbusSpec.endpoint.set(SettingsUtil.getRequiredValue(sessionSettingsReader, NimbusConfigKeys.NIMBUS_IP_KEY));

      nimbusSpec.username.set(SettingsUtil.getRequiredValue(sessionSettingsReader, NimbusConfigKeys.NIMBUS_USER_KEY));

      nimbusSpec.password.set(SettingsUtil
         .getRequiredValue(sessionSettingsReader, NimbusConfigKeys.NIMBUS_PASSWORD_KEY));

      nimbusSpec.deployUser.set(SettingsUtil.getRequiredValue(
         sessionSettingsReader,
         NimbusConfigKeys.NIMBUS_DEPLOY_USER_KEY));

      nimbusSpec.pod.set(SettingsUtil.getRequiredValue(sessionSettingsReader, NimbusConfigKeys.NIMBUS_POD_KEY));

      _vmfsStorageProvisionerSpec.service.set(nimbusSpec);
   }

   // nimbus template version
   @Override
   public String determineResourceVersion() throws Exception {
      return "1.0";
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Retrieve vmfs as resource is not implemented");
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter) throws Exception {

      _logger.info("Start deployment of new vmfs storage...");

      NimbusCommandsUtil.deployVmfsStorage(_vmfsStorageProvisionerSpec);

      // Assembler info
      // Nimbus provider specific info
      testbedSettingsWriter.setSetting(TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME, _vmfsStorageProvisionerSpec.vmName.get());
      testbedSettingsWriter.setSetting(TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD, _vmfsStorageProvisionerSpec.pod.get());
      testbedSettingsWriter.setSetting(
         TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER,
         _vmfsStorageProvisionerSpec.service.get().username.get());
      testbedSettingsWriter.setSetting(TESTBED_ASSEMBLER_KEY_SSH_IP, _vmfsStorageProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(TESTBED_ASSEMBLER_KEY_SSH_IPV4, _vmfsStorageProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(TESTBED_ASSEMBLER_KEY_SSH_IPV6, _vmfsStorageProvisionerSpec.ip.get());

      // Publisher info
      testbedSettingsWriter.setSetting(TESTBED_KEY_IP, _vmfsStorageProvisionerSpec.ip.get());

      _logger.info("Saving test bed connection data...");

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      boolean isAlive = true;

      // TODO: Find a common location to this code. Probably most of the health check will be with similar
      // implementation
      Socket socket = new Socket();
      try {
         socket.connect(
            new InetSocketAddress(_vmfsStorageProvisionerSpec.ip.get(), VMFS_SERVER_TCP_PORT),
            TCP_CONNECTION_TIMEOUT);
      } catch (IOException e) {
         _logger.debug("Unable to connect to TCP port " + VMFS_SERVER_TCP_PORT + "; VMFS server might be dead", e);
         isAlive = false;
      } finally {
         try {
            socket.close();
         } catch (IOException e) {
            _logger.debug("Unable to close VMFS Server TCP connection", e);
         }
      }

      return isAlive;
   }

   @Override
   public void destroyTestbed() throws Exception {
      _logger.info("Delete nimbus VMFS VM...");

      NimbusCommandsUtil.destroyVM(_vmfsStorageProvisionerSpec);
      NimbusCommandsUtil.deleteResultFile(_vmfsStorageProvisionerSpec);
   }

   @Override
   public int providerWeight() {
      return 1;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return VmfsStorageProvider.class;
   }
}
