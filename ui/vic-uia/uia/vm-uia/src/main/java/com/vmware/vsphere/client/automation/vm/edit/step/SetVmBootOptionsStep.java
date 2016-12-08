/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.createvm.view.CustomizeHardwarePage;
import com.vmware.vsphere.client.automation.vm.lib.messages.VmHardwareMessages;
import com.vmware.vsphere.client.automation.vm.edit.spec.VmBootOptionsSpec;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * Step that fills the Boot VM options in Customize Hardware page
 */
public class SetVmBootOptionsStep extends CommonUIWorkflowStep {

   protected VmBootOptionsSpec bootOptions;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      super.prepare(filteredWorkflowSpec);
      bootOptions = filteredWorkflowSpec.get(VmBootOptionsSpec.class);

      ensureNotNull(bootOptions, "No VmBootOptions object was linked to the spec.");
      ensureNotNull(bootOptions.firmware, "No Firmware was linked to the VmBootOptions.");
   }

   @Override
   public void execute() throws Exception {
      CustomizeHardwarePage customizePage = new CustomizeHardwarePage();

      customizePage.waitForLoadingProgressBar();
      customizePage.selectCustomizeHardwareTab(I18n.get(VmHardwareMessages.class).vmOptionsTab());
      customizePage.waitForLoadingProgressBar();
      customizePage.expandBootOptionsSection();
      customizePage.waitForLoadingProgressBar();
      // Set firmware
      customizePage.setFirmwareBootOption(bootOptions.firmware.get());
      customizePage.waitForLoadingProgressBar();
      customizePage.setSecurityBootOption(bootOptions.securityBoot.get());
      customizePage.waitForLoadingProgressBar();
   }
}
