/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.io.IOException;
import java.io.StringWriter;
import java.net.Socket;

import org.apache.commons.lang.RandomStringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.ssh.SshConnectionManager;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.PhysicalNfsStorageProvisionerSpec;

/**
 * Provider workflow for common test bed setup:
 *
 * shared nfs storage based on a physical NFS server described by
 * user credentials and parent folder in which to create shared folders
 * session.settings should contain
 * resource.physical.nfs.storage.ip = "ip"
 * resource.physical.nfs.storage.root.user = "user"
 * resource.physical.nfs.storage.root.password = "password"
 * resource.physical.nfs.storage.parent.folder = "parent folder"
 */
public class PhysicalNfsStorageProvider extends NfsStorageProvider {

   // Commands
   private static final String PHYSICAL_NFS_STORAGE_IP =
         "resource.physical.nfs.storage.ip";
   private static final String PHYSICAL_NFS_STORAGE_USER =
         "resource.physical.nfs.storage.root.user";
   private static final String PHYSICAL_NFS_STORAGE_PASSWORD =
         "resource.physical.nfs.storage.root.password";
   // this is the parent folder in which the folder representing the shared storage is created
   private static final String PHYSICAL_NFS_STORAGE_PARENT_FOLDER =
         "resource.physical.nfs.storage.parent.folder";
   private static String DEPLOY_PHYSICAL_NFS_STORAGE_CMD = "mkdir --mode=777 %s";
   private static String DESTROY_PHYSICAL_NFS_STORAGE_CMD = "rm -r -f %s";

   // Assembler Spec
   private PhysicalNfsStorageProvisionerSpec _physicalNfsStorageProvisionerSpec;

   // Testbed Assembler Settings
   private static final String TESTBED_ASSEMBLER_PHYSICAL_STORAGE_USER =
         "testbed.assembler.physical.nfs.storage.root.user";
   private static final String TESTBED_ASSEMBLER_PHYSICAL_STORAGE_PASSWORD =
         "testbed.assembler.physical.nfs.storage.root.password";

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(PhysicalNfsStorageProvider.class);

   private static final int NFS_SERVER_TELNET_PORT = 2049;
   private static final String PHYSICAL_NFS_STORAGE_PREFIX = "PhysicalNfsStorage_";

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
      _physicalNfsStorageProvisionerSpec = new PhysicalNfsStorageProvisionerSpec();
      assemblerSpec.links.add(_physicalNfsStorageProvisionerSpec);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {

      // _physicalNfsStorageProvisionerSpec is added in the init stage but add
      // the line for consistency
      _physicalNfsStorageProvisionerSpec =
            assemblerSpec.links.get(PhysicalNfsStorageProvisionerSpec.class);

      _physicalNfsStorageProvisionerSpec.ip.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_IP));

      _physicalNfsStorageProvisionerSpec.username.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_PHYSICAL_STORAGE_USER));

      _physicalNfsStorageProvisionerSpec.password.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_PHYSICAL_STORAGE_USER));
      _physicalNfsStorageProvisionerSpec.folder.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_KEY_FOLDER));

      //      _physicalNfsStorageProvisionerSpec.name.set(SettingsUtil.getRequiredValue(
      //            testbedSettings,
      //            TESTBED_KEY_NAME));

      _logger.info("Loaded assemblerSpec: " + _physicalNfsStorageProvisionerSpec);
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpecSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      _physicalNfsStorageProvisionerSpec =
            filteredAssemblerSpec.links.get(PhysicalNfsStorageProvisionerSpec.class);

      if (isAssembling) {
         // Only used in assemble command - for other stages use the
         // settings loaded from the settings file.
         final String name =
               PHYSICAL_NFS_STORAGE_PREFIX + RandomStringUtils.randomAlphanumeric(5);
         //         _physicalNfsStorageProvisionerSpec.name.set(name);
         _physicalNfsStorageProvisionerSpec.username.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               PHYSICAL_NFS_STORAGE_USER));
         _physicalNfsStorageProvisionerSpec.password.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               PHYSICAL_NFS_STORAGE_PASSWORD));
         _physicalNfsStorageProvisionerSpec.ip.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               PHYSICAL_NFS_STORAGE_IP));
         final String folder =
               SettingsUtil.getRequiredValue(
                     sessionSettingsReader,
                     PHYSICAL_NFS_STORAGE_PARENT_FOLDER) + "/" + name;
         _physicalNfsStorageProvisionerSpec.folder.set(folder);
      }
   }

   // fake version for physical nfs storage
   @Override
   public String determineResourceVersion() throws Exception {
      return "1.0";
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Retrieve nfs as resource is not implemeted");
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter) throws Exception {

      _logger.info("Start deployment of new physical nfs storage...");

      deployPhysicalNfsStorage(_physicalNfsStorageProvisionerSpec);

      // Publisher info
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_IP,
            _physicalNfsStorageProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_FOLDER,
            _physicalNfsStorageProvisionerSpec.folder.get());
      testbedSettingsWriter.setSetting(TESTBED_KEY_NAME, _physicalNfsStorageProvisionerSpec.name.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_PHYSICAL_STORAGE_USER,
            _physicalNfsStorageProvisionerSpec.username.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_PHYSICAL_STORAGE_PASSWORD,
            _physicalNfsStorageProvisionerSpec.password.get());

      _logger.info("Saving test bed connection data...");

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {

      boolean isHealthy = true;

      Socket socket = null;
      try {
         socket =
               new Socket(_physicalNfsStorageProvisionerSpec.ip.get(),
                     NFS_SERVER_TELNET_PORT);
      } catch (IOException e) {
         isHealthy = false;
      } finally {
         if (socket != null) {
            socket.close();
         }
      }

      return isHealthy;

   }

   @Override
   public void destroyTestbed() throws Exception {
      _logger.info("Delete physical storage ...");

      SshConnectionManager sshConnectionManager =
            new SshConnectionManager(_physicalNfsStorageProvisionerSpec.ip.get(),
                  _physicalNfsStorageProvisionerSpec.username.get(),
                  _physicalNfsStorageProvisionerSpec.password.get());

      String cmd =
            String.format(
                  DESTROY_PHYSICAL_NFS_STORAGE_CMD,
                  _physicalNfsStorageProvisionerSpec.folder.get());

      try {
         // delete result file
         sshConnectionManager.getDefaultConnection().executeSshCommand(
               cmd,
               new StringWriter(),
               new StringWriter());
      } finally {
         sshConnectionManager.close();
      }
   }

   @Override
   public int providerWeight() {
      return 0;
   }

   /**
    * Method that creates a folder on the physical stroage provider that will
    * be used as nfs storage in testing
    *
    * @param _physicalNfsStorageProvisionerSpec
    */
   private static void deployPhysicalNfsStorage(
         PhysicalNfsStorageProvisionerSpec _physicalNfsStorageProvisionerSpec)
               throws Exception {

      SshConnectionManager sshConnectionManager =
            new SshConnectionManager(_physicalNfsStorageProvisionerSpec.ip.get(),
                  _physicalNfsStorageProvisionerSpec.username.get(),
                  _physicalNfsStorageProvisionerSpec.password.get());

      try {
         sshConnectionManager.getDefaultConnection().executeSshCommand(
               String.format(
                     DEPLOY_PHYSICAL_NFS_STORAGE_CMD,
                     _physicalNfsStorageProvisionerSpec.folder.get()),
                     new StringWriter(),
                     new StringWriter());
      } finally {
         sshConnectionManager.close();
      }
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      // it should be base nfs storage
      return NfsStorageProvider.class;
   }

}
