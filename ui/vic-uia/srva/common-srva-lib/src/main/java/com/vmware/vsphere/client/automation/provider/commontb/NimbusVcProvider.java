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
import com.vmware.vsphere.client.automation.provider.commontb.spec.NimbusServiceSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.VcProvisionerSpec;
import com.vmware.vsphere.client.automation.provider.util.NimbusCommandsUtil;

/**
 * Class that provides the ability to deploy, check and destroy CludVM/CisWin
 * on a Nimbus environment.
 */
public class NimbusVcProvider extends BaseVcProvider {

   private static final String VSC_URL_TEMPLATE = "https://%s/vsphere-client/";

   // Session info
   private static final String RESOURCE_KEY_PRODUCT = "resource.vc.product";
   private static final String RESOURCE_KEY_BRANCH = "resource.vc.branch";
   private static final String RESOURCE_KEY_BUILDNUM = "resource.vc.buildNumber";

   public static final String RESOURCE_KEY_SSO_USER = "resource.vc.sso.admin.user";
   public static final String RESOURCE_KEY_SSO_PASSWORD =
         "resource.vc.sso.admin.password";

   // Testbed assembler settings
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME =
         "testbed.assembler.nimbus.vm.name";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER =
         "testbed.assembler.nimbus.vm.owner";
   private static final String TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD =
         "testbed.assembler.nimbus.vm.pod";
   private static final String TESTBED_ASSEMBLER_KEY_PRODUCT =
         "testbed.resource.product";
   private static final String TESTBED_ASSEMBLER_KEY_BUILD_NUMBER =
         "testbed.resource.build.number";
   private static final String TESTBED_ASSEMBLER_KEY_BRANCH = "testbed.resource.branch";

   // Assembler Spec
   private VcProvisionerSpec vcProvisionerSpec;


   // logger
   private static final Logger _logger = LoggerFactory.getLogger(NimbusVcProvider.class);

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {

      vcProvisionerSpec = new VcProvisionerSpec();
      assemblerSpec.links.add(vcProvisionerSpec);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {

      // actually the _vcProvisionerSpec is added in the init stage but add the
      // line for consistency
      vcProvisionerSpec = assemblerSpec.links.get(VcProvisionerSpec.class);

      // NOTE: Nimbus service is set in the prepare stage - check preapreForOperations

      vcProvisionerSpec.ip.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_VC_IP));

      vcProvisionerSpec.user.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_VC_USER));

      vcProvisionerSpec.password.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_VC_PASSWORD));

      vcProvisionerSpec.vmName.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME));

      vcProvisionerSpec.vmOwner.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER));

      vcProvisionerSpec.pod.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD));

      // Load product build info
      vcProvisionerSpec.branch.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_BRANCH));

      vcProvisionerSpec.product.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_PRODUCT));

      vcProvisionerSpec.buildNumber.set(SettingsUtil.getRequiredValue(
            testbedSettings,
            TESTBED_ASSEMBLER_KEY_BUILD_NUMBER));

      _logger.info("Loaded assemblerSpec: " + vcProvisionerSpec);
   }

   @Override
   public void prepareForOperations(PublisherSpec filteredPublisherSpecSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      vcProvisionerSpec = filteredAssemblerSpec.links.get(VcProvisionerSpec.class);

      if (isAssembling) {
         // Only used in assemble command - for other stages use the
         // settings loaded from the settings file.
         vcProvisionerSpec.product.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               RESOURCE_KEY_PRODUCT));
         vcProvisionerSpec.branch.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               RESOURCE_KEY_BRANCH));
         vcProvisionerSpec.buildNumber.set(SettingsUtil.getRequiredValue(
               sessionSettingsReader,
               RESOURCE_KEY_BUILDNUM));
         vcProvisionerSpec.vmName.set("VC_VM_" + vcProvisionerSpec.product.get() + "_"
               + vcProvisionerSpec.buildNumber.get() + "_"
               + RandomStringUtils.randomNumeric(5));
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

      vcProvisionerSpec.service.set(nimbusSpec);
   }

   @Override
   public String determineResourceVersion() throws Exception {
      // Version product-branch-buildnumber
      return String.format(
            "%s-%s-%s",
            vcProvisionerSpec.product.get(),
            vcProvisionerSpec.branch.get(),
            vcProvisionerSpec.buildNumber.get());
   }

   @Override
   public void retrieveResource() throws Exception {
      throw new RuntimeException("Retrieve vc as resource is not yet implemeted");
   }

   @Override
   public void deployTestbed(SettingsWriter testbedSettingsWriter) throws Exception {

      _logger.info("Start deployment of new VC...");

      NimbusCommandsUtil.deployVc(vcProvisionerSpec);


      // Assembler info
      // Nimbus provider specific info
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_NAME,
            vcProvisionerSpec.vmName.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_POD,
            vcProvisionerSpec.pod.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_NIMBUS_VM_OWNER,
            ((NimbusServiceSpec) vcProvisionerSpec.service.get()).username.get());

      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_VC_IP,
            vcProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_VC_USER,
            vcProvisionerSpec.user.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_VC_PASSWORD,
            vcProvisionerSpec.password.get());

      // Publisher info
      testbedSettingsWriter.setSetting(TESTBED_KEY_ENDPOINT, vcProvisionerSpec.ip.get());

      //TODO: rkovachev - we put IP instead of hostname
      // later we have to extract hostname from nimbus results - PR1371376
      testbedSettingsWriter.setSetting(TESTBED_KEY_NAME, vcProvisionerSpec.ip.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_VSC_URL,
            String.format(VSC_URL_TEMPLATE, vcProvisionerSpec.ip.get()));
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_USERNAME,
            vcProvisionerSpec.user.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_PASSWORD,
            vcProvisionerSpec.password.get());

      // Build info
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_PRODUCT,
            vcProvisionerSpec.product.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_BRANCH,
            vcProvisionerSpec.branch.get());
      testbedSettingsWriter.setSetting(
            TESTBED_ASSEMBLER_KEY_BUILD_NUMBER,
            vcProvisionerSpec.buildNumber.get());

      _logger.info("Saving test bed connection data...");
   }

   @Override
   public boolean checkTestbedHealth() throws Exception {
      // TODO: rkovachev implement provider specific check.
      // The check health for now is just checking the VC service is working.
      // It is done by the runEstablishConnections() method part of the provider work flow.
      _logger.info("Health checking for: " + vcProvisionerSpec.ip.toString() + " "
            + vcProvisionerSpec.pod.toString() + " is DONE by the assign connection "
            + "and connect methods of the provider workflow!");
      return true;
   }

   @Override
   public void destroyTestbed() throws Exception {
      _logger.info("Delete nimbus VM...");

      NimbusCommandsUtil.destroyVM(vcProvisionerSpec);
   }

   @Override
   public int providerWeight() {
      return 4;
   }


   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      return BaseVcProvider.class;
   }
}
