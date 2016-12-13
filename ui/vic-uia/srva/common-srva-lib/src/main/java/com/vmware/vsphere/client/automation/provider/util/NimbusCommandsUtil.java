/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.util;

import java.io.IOException;
import java.io.StringWriter;

import org.json.JSONException;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.ssh.SshConnection;
import com.vmware.client.automation.util.ssh.SshConnectionManager;
import com.vmware.vsphere.client.automation.provider.commontb.spec.HostProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NfsStorageProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusServiceSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.VcProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.VmSnapshotSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.VmfsStorageProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.XVpProvisionerSpec;

/**
 * Utility class that provides ability to execute basic Nimbus commands deploy, destoy vm and etc.
 * More details for Nimbus can be founded here: https://wiki.eng.vmware.com/ToolsEng/Projects/Nimbus/
 */
public class NimbusCommandsUtil {

   private static String DEPLOY_CLOUDVM_CMD =
         "USER=%s NIMBUS=%s /mts/git/bin/nimbus-vcvadeploy --cpus 8 --memory 16384 --vcvaBuild %s %s --result /tmp/%s.json --useQaNgc --lease 20";
   private static String DEPLOY_CISWIN_CMD =
         "USER=%s NIMBUS=%s /mts/git/bin/nimbus-vcdeploy-cat --vpxBuild %s %s --result /tmp/%s.json --useQaNgc --lease 20";
   private static String DEPLOY_ESX_CMD =
         "USER=%s NIMBUS=%s /mts/git/bin/nimbus-esxdeploy --preferLinkedClones --useVHV --nics 6 --disk 10000000 --disk 10000000 --memory 8192 --cpus 4"
         + " %s %s --result /tmp/%s.json --lease 20 --usePrepared";
   private static String DEPLOY_NFS_STORAGE_CMD =
         "USER=%s NIMBUS=%s /mts/git/bin/nimbus-nfsdeploy --result /tmp/%s.json %s --lease 20";
   private static String DEPLOY_VMFS_STORAGE_CMD =
      "USER=%s NIMBUS=%s /mts/git/bin/nimbus-iscsideploy --result /tmp/%s.json %s --lun 10:2 --lease 20";
   private static String DESTROY_CLOUDVM_CMD =
         "USER=%s NIMBUS=%s /mts/git/bin/nimbus-ctl -p nimbus/%s kill %s";
   private static String GET_VM_INFO_CMD = "cat /tmp/%s.json";
   private static String NIMBUS_NFS_DELETE_RESULT_FILE_CMD_KEY = "rm -f -r %s";
   private static String REVERT_VM_SNAPSHOT_CMD = "USER=%s NIMBUS=%s /mts/git/bin/nimbus-ctl --snapshot %s revert-snapshot %s";
   private static String CREATE_VM_SNAPSHOT_CMD = "USER=%s NIMBUS=%s /mts/git/bin/nimbus-ctl --snapshot %s --snapshotIncludeMemory create-snapshot %s";
   private static String DEPLOY_XVP_CMD = "USER=%s NIMBUS=%s /mts/git/bin/nimbus-ovfdeploy -p xvp_config=%s -p xvp_version=%s %s %s --result /tmp/%s.json --lease 20";

   // Common for all Nimbus vms
   private static final String NIMBUS_RESULT_JSON_POD_KEY = "pod";
   private static final String NIMBUS_RESULT_JSON_NAME_KEY = "name";
   private static final String NIMBUS_RESULT_JSON_IP_KEY = "ip";
   private static final String NIMBUS_RESULT_JSON_PASSWORD_KEY = "password";
   private static final String TMP_FOLDER = "/tmp/";
   private static final String JSON_EXT = ".json";

   // NFS storage specific
   private static final String NIMBUS_RESULT_JSON_IP4_KEY = "ip4";
   private static final String NIMBUS_RESULT_JSON_IP6_KEY = "ip6";

   // VC specific
   private static final String NIMBUS_RESULT_JSON_VIM_USER_KEY = "vimUsername";
   private static final String NIMBUS_RESULT_JSON_VIM_PASSWORD_KEY = "vimPassword";

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(NimbusCommandsUtil.class);

   public static void deployVc(VcProvisionerSpec vcProvisionerSpec) {

      NimbusServiceSpec nimbusServiceSpec =
            ((NimbusServiceSpec) vcProvisionerSpec.service.get());

      String nimbusVmName = vcProvisionerSpec.vmName.get();

      String command = DEPLOY_CISWIN_CMD;
      if (vcProvisionerSpec.product.get().contains("cloudvm")) {
         command = DEPLOY_CLOUDVM_CMD;
      }

      String deployCmd =
            String.format(
                  command,
                  nimbusServiceSpec.deployUser.get(),
                  nimbusServiceSpec.pod.get(),
                  vcProvisionerSpec.buildNumber.get(),
                  nimbusVmName,
                  nimbusVmName);

      _logger.info("Deploy CMD: " + deployCmd);

      JSONObject resutlJson = deployVM(vcProvisionerSpec, deployCmd);

      //      String jsonResult =
      //            "{\"username\":\"root\",\"password\":\"vmware\","
      //                  + "\"name\":\"mts-automation-VC_VM_cloudvm_sb-3598864_95782\","
      //                  + "\"selected_computer_resource\":\"Team_Labs\",\"pod\":\"rila-1\","
      //                  + "\"vimUsername\":\"administrator@vsphere.local\",\"vimPassword\":\"Admin!23\","
      //                  + "\"ip\":\"10.26.253.106\",\"ip4\":\"10.26.253.106\",\"ip6\":null,"
      //                  + "\"systemPNID\":\"sof2-rila-6-106.eng.vmware.com\",\"deploy_status\":\"success\"}";


      try {
         vcProvisionerSpec.pod.set(resutlJson.getString(NIMBUS_RESULT_JSON_POD_KEY));
         vcProvisionerSpec.ip.set(resutlJson.getString(NIMBUS_RESULT_JSON_IP_KEY));
         vcProvisionerSpec.password.set(resutlJson
               .getString(NIMBUS_RESULT_JSON_PASSWORD_KEY));
         vcProvisionerSpec.vmName.set(resutlJson.getString(NIMBUS_RESULT_JSON_NAME_KEY));

         vcProvisionerSpec.user.set(resutlJson
               .getString(NIMBUS_RESULT_JSON_VIM_USER_KEY));

         vcProvisionerSpec.password.set(resutlJson
               .getString(NIMBUS_RESULT_JSON_VIM_PASSWORD_KEY));


      } catch (JSONException e) {
         // TODO Auto-generated catch block
         e.printStackTrace();
      }


   }

   /**
    *
    * @param hostProvisionerSpec
    * @param buildNimber
    */
   public static void deployHost(HostProvisionerSpec hostProvisionerSpec) {

      NimbusServiceSpec nimbusServiceSpec =
            ((NimbusServiceSpec) hostProvisionerSpec.service.get());

      String nimbusVmName = hostProvisionerSpec.vmName.get();

      String deployCmd =
            String.format(
                  DEPLOY_ESX_CMD,
                  nimbusServiceSpec.deployUser.get(),
                  nimbusServiceSpec.pod.get(),
                  nimbusVmName,
                  hostProvisionerSpec.buildNumber.get(),
                  nimbusVmName);

      _logger.info("Deploy CMD: " + deployCmd);

      JSONObject resutlJson = deployVM(hostProvisionerSpec, deployCmd);

      //      String jsonResult = "{\"name\":\""
      //      + _hostProvisionerSpec.vmName.get()
      //      + "\",\"mac\":\"00:50:56:96:75:8e\",\"pod\":\"rila-1\",\"ip\":\"10.26.253.66\",\"password\":\"ca$hc0w\"}";

      try {
         hostProvisionerSpec.pod.set(resutlJson.getString(NIMBUS_RESULT_JSON_POD_KEY));
         hostProvisionerSpec.ip.set(resutlJson.getString(NIMBUS_RESULT_JSON_IP_KEY));
         hostProvisionerSpec.password.set(resutlJson
               .getString(NIMBUS_RESULT_JSON_PASSWORD_KEY));
         hostProvisionerSpec.vmName.set(resutlJson
               .getString(NIMBUS_RESULT_JSON_NAME_KEY));

      } catch (JSONException e) {
         // TODO Auto-generated catch block
         e.printStackTrace();
      }

   }

   /**
    * Deploys a VM in the Nimbus infrastructure that has a configured NFS server. For more information,
    * check Nimbus wiki.
    * @param nfsStorageProviderSpec
    */
   public static void deployNfsStorage(NfsStorageProvisionerSpec nfsStorageProviderSpec) {

      NimbusServiceSpec nimbusServiceSpec =
            ((NimbusServiceSpec) nfsStorageProviderSpec.service.get());
      String nimbusVmName = nfsStorageProviderSpec.vmName.get();

      String deployCmd =
            String.format(
                  DEPLOY_NFS_STORAGE_CMD,
                  nimbusServiceSpec.deployUser.get(),
                  nimbusServiceSpec.pod.get(),
                  nimbusVmName,
                  nimbusVmName);

      _logger.info("Deploy CMD: {}", deployCmd);

      JSONObject resultJson = deployVM(nfsStorageProviderSpec, deployCmd);

      try {
         nfsStorageProviderSpec.pod
         .set(resultJson.getString(NIMBUS_RESULT_JSON_POD_KEY));
         nfsStorageProviderSpec.ip.set(resultJson.getString(NIMBUS_RESULT_JSON_IP_KEY));
         nfsStorageProviderSpec.ipv4.set(resultJson
               .getString(NIMBUS_RESULT_JSON_IP4_KEY));
         nfsStorageProviderSpec.ipv6.set(resultJson
               .getString(NIMBUS_RESULT_JSON_IP6_KEY));
         nfsStorageProviderSpec.vmName.set(resultJson
               .getString(NIMBUS_RESULT_JSON_NAME_KEY));

      } catch (JSONException e) {
         e.printStackTrace();
      }
   }

   /**
    * Deploys a VM in the Nimbus infrastructure that has a configured VMFS server. For more information,
    * check Nimbus wiki
    *
    * @param nfsStorageProviderSpec
    */
   public static void deployVmfsStorage(VmfsStorageProvisionerSpec nfsStorageProviderSpec) {

      NimbusServiceSpec nimbusServiceSpec = ((NimbusServiceSpec) nfsStorageProviderSpec.service.get());
      String nimbusVmName = nfsStorageProviderSpec.vmName.get();

      String deployCmd = String.format(
         DEPLOY_VMFS_STORAGE_CMD,
         nimbusServiceSpec.deployUser.get(),
         nimbusServiceSpec.pod.get(),
         nimbusVmName,
         nimbusVmName);

      _logger.info("Deploy CMD: {}", deployCmd);

      JSONObject resultJson = deployVM(nfsStorageProviderSpec, deployCmd);

      try {
         nfsStorageProviderSpec.pod.set(resultJson.getString(NIMBUS_RESULT_JSON_POD_KEY));
         nfsStorageProviderSpec.ip.set(resultJson.getString(NIMBUS_RESULT_JSON_IP_KEY));
         nfsStorageProviderSpec.ipv4.set(resultJson.getString(NIMBUS_RESULT_JSON_IP4_KEY));
         nfsStorageProviderSpec.ipv6.set(resultJson.getString(NIMBUS_RESULT_JSON_IP6_KEY));
         nfsStorageProviderSpec.vmName.set(resultJson.getString(NIMBUS_RESULT_JSON_NAME_KEY));
      } catch (JSONException e) {
         e.printStackTrace();
      }
   }

   /**
    * Deploys xVP Storage Provider in the Nimbus. For more information, check
    * Nimbus wiki.
    * @param xVpProvisionerSpec
    */
   public static void deployXvp(XVpProvisionerSpec xVpProvisionerSpec) {
      NimbusServiceSpec nimbusServiceSpec = ((NimbusServiceSpec) xVpProvisionerSpec.service
            .get());
      String replConfig = xVpProvisionerSpec.replConfig.get();
      String vpVersion = xVpProvisionerSpec.version.get();
      String nimbusVmName = xVpProvisionerSpec.vmName.get();
      String url = xVpProvisionerSpec.url.get();

      String deployCmd = String.format(DEPLOY_XVP_CMD,
           nimbusServiceSpec.deployUser.get(),
           nimbusServiceSpec.pod.get(),
           replConfig,
           vpVersion,
           nimbusVmName,
           url,
           nimbusVmName);


      _logger.info("Deploy CMD: {}", deployCmd);

      JSONObject resultJson = deployVM(xVpProvisionerSpec, deployCmd);

      try {
         xVpProvisionerSpec.pod.set(resultJson.getString(NIMBUS_RESULT_JSON_POD_KEY));
         xVpProvisionerSpec.ip.set(resultJson.getString(NIMBUS_RESULT_JSON_IP_KEY));
         xVpProvisionerSpec.vmName.set(resultJson.getString(NIMBUS_RESULT_JSON_NAME_KEY));
      } catch (JSONException e) {
         e.printStackTrace();
      }

   }

   /**
    * Destroy VM deployed in Nimbus infrastructure.
    *
    * @param nimbusProvisionerSpec
    */
   public static void destroyVM(NimbusProvisionerSpec nimbusProvisionerSpec) {

      String destroyVmCmd =
            String.format(
                  DESTROY_CLOUDVM_CMD,
                  nimbusProvisionerSpec.vmOwner.get(),
                  nimbusProvisionerSpec.pod.get(),
                  nimbusProvisionerSpec.vmOwner.get(),
                  nimbusProvisionerSpec.vmName.get());

      _logger.info("Destroy VM CMD: " + destroyVmCmd);

      SshConnection sshConnection =
            establishSshConnection((NimbusServiceSpec) nimbusProvisionerSpec.service
                  .get());
      try {
         // Exec destroy command
         execSshCommand(sshConnection, destroyVmCmd);
      } finally {
         sshConnection.close();
      }
   }

   /**
    * Method that deletes the json result file on Nimbus; it expects it is in
    * /tmp/'nimbusProvisionerSpec.name.get()'.json
    *
    * @param nimbusProvisionerSpec
    *            - it takes the nimbus service spec in order to connect to
    *            nimbus and name, as it presumes this is the name of the temp
    *            file
    */
   public static void deleteResultFile(NimbusProvisionerSpec nimbusProvisionerSpec) {
      NimbusServiceSpec nimbusSvc =
            (NimbusServiceSpec) nimbusProvisionerSpec.service.get();
      SshConnectionManager sshConnectionManager =
            new SshConnectionManager(nimbusSvc.endpoint.get(), nimbusSvc.username.get(),
                  nimbusSvc.password.get());

      try {
         // delete result file
         execSshCommand(
               sshConnectionManager.getDefaultConnection(),
               String.format(NIMBUS_NFS_DELETE_RESULT_FILE_CMD_KEY, TMP_FOLDER
                     + nimbusProvisionerSpec.vmName.get() + JSON_EXT));
      } finally {
         sshConnectionManager.close();
      }
   }

   /**
    * Creates snapshot of Virtual Machine on Nimbus
    *
    * @param nimbusProvisionerSpec
    *           - it takes the nimbus service spec in order to connect to nimbus
    * @param vmSnapshotSpec
    *           - spec representing the properties of the vm snapshot
    */
   public static void createVmSnapshot(
         NimbusProvisionerSpec nimbusProvisionerSpec,
         VmSnapshotSpec vmSnapshotSpec) {
      String createSnapshotCmd = String.format(CREATE_VM_SNAPSHOT_CMD,
            nimbusProvisionerSpec.vmOwner.get(), vmSnapshotSpec.vmSpec.get().pod.get(),
            vmSnapshotSpec.snapshotName.get(), vmSnapshotSpec.vmSpec.get().vmName.get());

      _logger.info("Create snapshot of VM CMD: " + createSnapshotCmd);

      SshConnection sshConnection = establishSshConnection((NimbusServiceSpec) nimbusProvisionerSpec.service
            .get());
      try {
         execSshCommand(sshConnection, createSnapshotCmd);
      } finally {
         sshConnection.close();
      }
   }

   /**
    * Reverts snapshot of Virtual Machine on Nimbus
    *
    * @param nimbusProvisionerSpec
    *           - it takes the nimbus service spec in order to connect to nimbus
    * @param vmSnapshotSpec
    *           - spec representing the properties of the vm snapshot
    */
   public static void revertVmSnapshot(
         NimbusProvisionerSpec nimbusProvisionerSpec,
         VmSnapshotSpec vmSnapshotSpec) {
      String revertSnapshotCmd = String.format(REVERT_VM_SNAPSHOT_CMD,
            nimbusProvisionerSpec.vmOwner.get(), vmSnapshotSpec.vmSpec.get().pod.get(),
            vmSnapshotSpec.snapshotName.get(), vmSnapshotSpec.vmSpec.get().vmName.get());

      _logger.info("Revert snapshot of VM CMD: " + revertSnapshotCmd);

      SshConnection sshConnection = establishSshConnection((NimbusServiceSpec) nimbusProvisionerSpec.service
            .get());
      try {
         execSshCommand(sshConnection, revertSnapshotCmd);
      } finally {
         sshConnection.close();
      }
   }

   /**
    * Execute an ssh command with result
    * @param sshConnection
    * @param command
    * @return
    */
   private static String execSshCommand(SshConnection sshConnection, String command) {
      StringWriter writer = new StringWriter();
      StringWriter errorWriter = new StringWriter();
      try {
         sshConnection.executeSshCommand(command, writer, errorWriter);
      } catch (IOException e) {
         // TODO Auto-generated catch block
         e.printStackTrace();
      }
      _logger.info("INFO:" + writer.toString());
      _logger.info("ERROR:" + errorWriter.toString());

      return writer.toString();
   }

   //

   private static SshConnection establishSshConnection(
         NimbusServiceSpec nimbusServiceSpec) {
      SshConnectionManager sshConnectionManager =
            new SshConnectionManager(nimbusServiceSpec.endpoint.get(),
                  nimbusServiceSpec.username.get(), nimbusServiceSpec.password.get());

      return sshConnectionManager.getDefaultConnection();
   }

   private static JSONObject deployVM(NimbusProvisionerSpec nimbusProvisionerSpec,
         String command) {

      JSONObject jsonObject = new JSONObject();
      SshConnection sshConnection =
            establishSshConnection((NimbusServiceSpec) nimbusProvisionerSpec.service
                  .get());

      try {
         // Exec deploy command
         execSshCommand(sshConnection, command);

         // Read JSON result
         String listJSONCmd =
               String.format(GET_VM_INFO_CMD, nimbusProvisionerSpec.vmName.get());
         _logger.info("Read JSON Result CMD: " + listJSONCmd);

         String jsonResult = execSshCommand(sshConnection, listJSONCmd);
         _logger.info("Nimbus Result JSON: " + jsonResult);

         JSONTokener tokener = new JSONTokener(jsonResult);
         try {
            jsonObject = new JSONObject(tokener);
         } catch (JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
         }
      } finally {
         sshConnection.close();
      }

      return jsonObject;
   }
}
