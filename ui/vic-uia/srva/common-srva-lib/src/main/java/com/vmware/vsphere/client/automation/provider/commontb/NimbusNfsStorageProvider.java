/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.io.IOException;
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
import com.vmware.vsphere.client.automation.provider.commontb.spec.NfsStorageProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusServiceSpec;
import com.vmware.vsphere.client.automation.provider.util.NimbusCommandsUtil;

/**
 * Provider workflow for common test bed setup:
 *
 * shared nfs storage deployed on a Nimbus Vm with configured NFS server
 *
 */
public class NimbusNfsStorageProvider extends NfsStorageProvider {

   // Session info

   // Product deploy settings - should be in session settings. It is part of
   // the configuration run.
   private static final String NFS_STORAGE_ROOT_USER = "resource.nfs.storage.root.user";
   private static final String NFS_STORAGE_ROOT_PASSWORD =
         "resource.nfs.storage.root.password";
   private static final String NFS_STORAGE_FOLDER = "resource.nfs.storage.folder";

   // Assemble info

   // Testbed assembler settings
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME =
         "testbed.assembler.nimbus.vm.name";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER =
         "testbed.assembler.nimbus.vm.owner";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD =
         "testbed.assembler.nimbus.vm.pod";

   private static final String TESTBED_ASSEMBLER_KEY_SSH_IP = "testbed.ssh.ip";
   private static final String TESTBED_ASSEMBLER_KEY_SSH_USER = "testbed.ssh.user";
   private static final String TESTBED_ASSEMBLER_KEY_SSH_PASSWORD = "testbed.ssh.pass";

   // Assembler Spec
   private NfsStorageProvisionerSpec _nfsStorageProvisionerSpec;

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(NimbusNfsStorageProvider.class);

   private static final int NFS_SERVER_TELNET_PORT = 2049;
   private static final String NIMBUS_NFS_STORAGE_PREFIX = "NFS_STORAGE_";

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
      _nfsStorageProvisionerSpec = new NfsStorageProvisionerSpec();
      assemblerSpec.links.add(_nfsStorageProvisionerSpec);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {

      // _nfsStorageProvisionerSpec is added in the init stage but add
      // the line for consistency
      _nfsStorageProvisionerSpec =
            assemblerSpec.links.get(NfsStorageProvisionerSpec.class);

      // NOTE: Nimbus service is set in the prepare stage - check
      // preapreForOperations

      _nfsStorageProvisionerSpec.ip.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_SSH_IP));

      _nfsStorageProvisionerSpec.user.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_SSH_USER));

      _nfsStorageProvisionerSpec.password.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_SSH_PASSWORD));
      _nfsStorageProvisionerSpec.vmName.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME));

      _nfsStorageProvisionerSpec.vmOwner.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER));

      _nfsStorageProvisionerSpec.pod.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD));

      _nfsStorageProvisionerSpec.name.set(SettingsUtil.getRequiredValue(
            testbedSettings, TESTBED_KEY_NAME));

      _logger.info("Loaded assemblerSpec: " + _nfsStorageProvisionerSpec);
   }


   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpecSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      _nfsStorageProvisionerSpec =
            filteredAssemblerSpec.links.get(NfsStorageProvisionerSpec.class);

      if (isAssembling) {
         final String name =
               NIMBUS_NFS_STORAGE_PREFIX + RandomStringUtils.randomNumeric(5);
         // Only used in assemble command - for other stages use the
         // settings loaded from the settings file.
         _nfsStorageProvisionerSpec.vmName.set(name);
         _nfsStorageProvisionerSpec.name.set(name);

         // Load nfs storage admin user and password as it is not provided by
         // Nimbus result file.
         _nfsStorageProvisionerSpec.user.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               NFS_STORAGE_ROOT_USER));
         _nfsStorageProvisionerSpec.password.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               NFS_STORAGE_ROOT_PASSWORD));
         _nfsStorageProvisionerSpec.folder.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               NFS_STORAGE_FOLDER));
      }

      // Set Nimbus service
      NimbusServiceSpec nimbusSpec = new NimbusServiceSpec();
      nimbusSpec.endpoint.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader,
            NimbusConfigKeys.NIMBUS_IP_KEY));

      nimbusSpec.username.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader,
            NimbusConfigKeys.NIMBUS_USER_KEY));

      nimbusSpec.password.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader,
            NimbusConfigKeys.NIMBUS_PASSWORD_KEY));

      nimbusSpec.deployUser.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader,
            NimbusConfigKeys.NIMBUS_DEPLOY_USER_KEY));

      nimbusSpec.pod.set(SettingsUtil.getRequiredValue(
            sessionSettingsReader,
            NimbusConfigKeys.NIMBUS_POD_KEY));

      _nfsStorageProvisionerSpec.service.set(nimbusSpec);
   }

   // nimbus template version
   @Override
   public String determineResourceVersion() throws Exception {
      return "1.0";
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Retrieve nfs as resource is not implemented");
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter) throws Exception {

      _logger.info("Start deployment of new nfs storage...");

      NimbusCommandsUtil.deployNfsStorage(_nfsStorageProvisionerSpec);

      // Assembler info
      // Nimbus provider specific info
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME,
            _nfsStorageProvisionerSpec.vmName.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD,
            _nfsStorageProvisionerSpec.pod.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER,
            ((NimbusServiceSpec) _nfsStorageProvisionerSpec.
                  service.get()).username.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_SSH_IP,
            _nfsStorageProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_SSH_USER,
            _nfsStorageProvisionerSpec.user.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_SSH_PASSWORD,
            _nfsStorageProvisionerSpec.password.get());

      // Publisher info
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_IP,
            _nfsStorageProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_FOLDER,
            _nfsStorageProvisionerSpec.folder.get());
      testbedSettingsWriter.setSetting(TESTBED_KEY_NAME, _nfsStorageProvisionerSpec.name.get());

      _logger.info("Saving test bed connection data...");

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {

      boolean isHealthy = true;

      Socket socket = null;
      try {
         socket =
               new Socket(_nfsStorageProvisionerSpec.ip.get(), NFS_SERVER_TELNET_PORT);
      } catch (IOException e) {
         isHealthy = false;
      } finally {
         if (socket != null) {
            socket.close();
         }
      }

      return isHealthy;

      /*
       * This here is the correct way to prove that the nfs server is up and
       * running Open connection to NFS server VM SshConnectionManager
       * sshConnectionManager = new SshConnectionManager(
       * _nfsStorageProvisionerSpec.ip.get(),
       * _nfsStorageProvisionerSpec.user.get(),
       * _nfsStorageProvisionerSpec.password.get()); try { SshConnection
       * sshConn = sshConnectionManager.getDefaultConnection(); if (sshConn !=
       * null) { // checking all needed processes are up and running String //
       * rpcinfo -p shows all nfs server processes result =
       * NimbusCommandsUtil.execSshCommand(sshConn, "rpcinfo -p"); //
       * portmapper, nfs and mountd should be up and running isHealthy =
       * result.contains("portmapper") && result.contains("nfs") &&
       * result.contains("mountd"); } return isHealthy; } finally {
       * sshConnectionManager.close(); }
       */
   }

   @Override
   public void destroyTestbed() throws Exception {
      _logger.info("Delete nimbus VM...");

      NimbusCommandsUtil.destroyVM(_nfsStorageProvisionerSpec);
      NimbusCommandsUtil.deleteResultFile(_nfsStorageProvisionerSpec);
   }

   @Override
   public int providerWeight() {
      return 1;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      // It should be base nfs storage
      return NfsStorageProvider.class;
   }
}
