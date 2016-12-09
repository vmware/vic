/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.vm.createvm.view.CustomizeHardwarePage;
import com.vmware.vsphere.client.automation.vm.lib.messages.VmHardwareMessages;
import com.vmware.vsphere.client.automation.vm.edit.spec.UpgradeVmSpec;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * Step that fills the Upgrade VM options in Customize Hardware page
 */
public class SetUpgradeHardwareOptionsStep extends CommonUIWorkflowStep {

   private UpgradeVmSpec _upgradeSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      super.prepare(filteredWorkflowSpec);

      _upgradeSpec = filteredWorkflowSpec.get(UpgradeVmSpec.class);

      ensureNotNull(_upgradeSpec, "VmSpec object is missing.");
      ensureAssigned(_upgradeSpec.scheduleUpdate, "Ensure that scheduleUpdate is added");
   }

   @Override
   public void execute() throws Exception {
      CustomizeHardwarePage customizePage = new CustomizeHardwarePage();

      customizePage.waitForLoadingProgressBar();
      customizePage.selectCustomizeHardwareTab(I18n.get(VmHardwareMessages.class).vmVirtualHardwareTab());
      customizePage.waitForLoadingProgressBar();
      // Set schedule upgrade configuration
      customizePage.setScheduleUpgrade(_upgradeSpec.scheduleUpdate.get());
      if (_upgradeSpec.compatibilityVersion.isAssigned()) {
         customizePage.expandUpgradeSection();
         customizePage.waitForLoadingProgressBar();
         customizePage.setCompatibilityVersion(_upgradeSpec.compatibilityVersion.get());
      }
   }
}
