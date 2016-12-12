/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

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
import com.vmware.vsphere.client.automation.provider.commontb.spec.HostProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusServiceSpec;
import com.vmware.vsphere.client.automation.provider.util.NimbusCommandsUtil;

/**
 * Provider that deploys ESX template on Nimbus environment using the
 * nimbus-esxdeploy-ob command.
 * That command deploys ESX that has two issues. The first is that the
 * datastore id is not unique so two ESX deployed with that command will
 * not be able to added in the same DC.
 * The second is that the vmotion traffic is not enabled.
 */
public class NimbusHostProvider extends BaseHostProvider {

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(NimbusHostProvider.class);

   // Assembler Spec
   private HostProvisionerSpec hostProvisionerSpec;

   // Product deploy setting - should be in session settings. It is part of the configuration run.
   private static final String RESOURCE_KEY_PRODUCT = "resource.host.product";
   private static final String RESOURCE_KEY_BRANCH = "resource.host.branch";
   private static final String RESOURCE_KEY_BUILDNUM = "resource.host.buildNumber";

   private static final String HOST_ROOT_USERNAME = "resource.host.root.user";
   private static final String HOST_SERVICE_PORT = "resource.host.service.port";

   private String _hostRootUsername;
   private String _hostServicePort;

   // Assemble info - Nimbus settings

   // Testbed assembler keys
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME =
         "testbed.assembler.nimbus.vm.name";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER =
         "testbed.assembler.nimbus.vm.owner";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD =
         "testbed.assembler.nimbus.vm.pod";

   private static final String TESTBED_ASSEMBLER_KEY_SSH_IP = "testbed.ssh.ip";
   private static final String TESTBED_ASSEMBLER_KEY_SSH_USER = "testbed.ssh.user";
   private static final String TESTBED_ASSEMBLER_KEY_SSH_PASSWORD = "testbed.ssh.pass";

   private static final String TESTBED_ASSEMBLER_KEY_PRODUCT =
         "testbed.resource.product";
   private static final String TESTBED_ASSEMBLER_KEY_BUILD_NUMBER =
         "testbed.resource.build.number";
   private static final String TESTBED_ASSEMBLER_KEY_BRANCH = "testbed.resource.branch";

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
      hostProvisionerSpec = new HostProvisionerSpec();
      assemblerSpec.links.add(hostProvisionerSpec);
   }


   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {

      _logger.info("Start HostProvider assign assemble specs");

      hostProvisionerSpec = assemblerSpec.links.get(HostProvisionerSpec.class);

      // NOTE: Nimbus service is set in the prepare stage - check preapreForOperations

      hostProvisionerSpec.ip.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_SSH_IP));

      hostProvisionerSpec.user.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_SSH_USER));

      hostProvisionerSpec.password.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_SSH_PASSWORD));

      hostProvisionerSpec.vmName.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME));

      hostProvisionerSpec.vmOwner.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER));

      hostProvisionerSpec.pod.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD));

      // Load product build info
      hostProvisionerSpec.branch.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_BRANCH));

      hostProvisionerSpec.product.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_PRODUCT));

      hostProvisionerSpec.buildNumber.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_BUILD_NUMBER));


      _logger.info("Loaded assemblerSpec: " + hostProvisionerSpec);
   }


   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpecSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      hostProvisionerSpec = filteredAssemblerSpec.links.get(HostProvisionerSpec.class);

      if (isAssembling) {
         // Only used in assemble command - for other stages use the
         // settings loaded from the settings file.

         hostProvisionerSpec.product.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               RESOURCE_KEY_PRODUCT));
         hostProvisionerSpec.branch.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               RESOURCE_KEY_BRANCH));
         hostProvisionerSpec.buildNumber.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               RESOURCE_KEY_BUILDNUM));
         hostProvisionerSpec.vmName.set("ESX_VM_" + hostProvisionerSpec.product.get()
               + "_" + hostProvisionerSpec.buildNumber.get() + "_"
               + RandomStringUtils.randomNumeric(5));

         // Load host admin user and host service port as it is not provided by Nimbus result file.
         _hostRootUsername =
               SettingsUtil.getRequiredValue(sessionSettingsReader, HOST_ROOT_USERNAME);

         _hostServicePort =
               SettingsUtil.getRequiredValue(sessionSettingsReader, HOST_SERVICE_PORT);
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

      hostProvisionerSpec.service.set(nimbusSpec);
   }


   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter) throws Exception {

      _logger.info("Start deployment of new host...");

      NimbusCommandsUtil.deployHost(hostProvisionerSpec);

      // TODO: rkovache find a way to detect that the host is up and running.
      // Although the Nimbus script reports that the host is deployed it has not
      // got the host service started and connecting to it return error 503.
      // The workaround is to wait for 30 seconds to be up.
      _logger.info("Wait 30 seconds the Host service to be up!");
      Thread.sleep(30000);

      // Assembler info
      // Nimbus provider specific info
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME,
            hostProvisionerSpec.vmName.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD,
            hostProvisionerSpec.pod.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER,
            ((NimbusServiceSpec) hostProvisionerSpec.service.get()).username.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_SSH_IP,
            hostProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_SSH_USER, _hostRootUsername);
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_SSH_PASSWORD,
            hostProvisionerSpec.password.get());

      // Build info
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_PRODUCT,
            hostProvisionerSpec.product.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_BRANCH,
            hostProvisionerSpec.branch.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_BUILD_NUMBER,
            hostProvisionerSpec.buildNumber.get());

      // Publisher info
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_ENDPOINT,
            hostProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(TESTBED_KEY_USERNAME, _hostRootUsername);
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_PASSWORD,
            hostProvisionerSpec.password.get());
      testbedSettingsWriter.setSetting(TESTBED_KEY_SERVICE_PORT, _hostServicePort);

      _logger.info("Saving test bed connection data...");

   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      // TODO: rkovachev implement it based on the _hostProviderSpec info
      _logger.info("Health checking for: " + hostProvisionerSpec.ip.toString() + " "
            + hostProvisionerSpec.pod.toString());
      return true;
   }

   @Override
   public String determineResourceVersion() throws Exception {
      // Version product-branch-buildnumber
      return String.format(
            "%s-%s-%s",
            hostProvisionerSpec.product.get(),
            hostProvisionerSpec.branch.get(),
            hostProvisionerSpec.buildNumber.get());
   }

   @Override
   public void destroyTestbed() throws Exception {
      _logger.info("Delete nimbus VM...");

      NimbusCommandsUtil.destroyVM(hostProvisionerSpec);
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return BaseHostProvider.class;
   }
}
