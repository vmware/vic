/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.cluster.lib.spec.EditClusterSpec;
import com.vmware.vsphere.client.automation.cluster.lib.view.EditClusterDrsSettingsPage;

/**
 * Operations included in this step:
 *  1. validate edit cluster spec
 *  2. Checks edit cluster spec for any DRS settings and sets all properties
 */
public class EditClusterDrsSettingsStep extends CommonUIWorkflowStep {

   private EditClusterSpec _editClusterSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _editClusterSpec = filteredWorkflowSpec.get(EditClusterSpec.class);

      ensureNotNull(_editClusterSpec, "EditClusterSpec spec is missing.");

      if (_editClusterSpec.drsAutomationLevel.isAssigned()
            && (!_editClusterSpec.drsEnabled.isAssigned() || !_editClusterSpec.drsEnabled
                  .get())) {
         throw new IllegalArgumentException(
               "EditClusterSpec is not configured correctly! "
                     + "DRS Automation Level cannot be used if DRS is not enabled.");
      }

      if (_editClusterSpec.vmDrsAutomationLevelEnabled.isAssigned()
            && _editClusterSpec.vmDrsAutomationLevelEnabled.get()
            && (!_editClusterSpec.drsEnabled.isAssigned() || !_editClusterSpec.drsEnabled
                  .get())) {
         throw new IllegalArgumentException(
               "EditClusterSpec is not configured correctly! "
                     + "Individual VM DRS Automation Level cannot be enabled if DRS is not.");
      }
   }

   @Override
   public void execute() throws Exception {
      EditClusterDrsSettingsPage editPage = new EditClusterDrsSettingsPage();

      if (_editClusterSpec.drsEnabled.isAssigned()) {
         editPage.setDrsEnabled(_editClusterSpec.drsEnabled.get());
      }

      if (_editClusterSpec.drsAutomationLevel.isAssigned()) {
         editPage.setDrsAutomationLevel(_editClusterSpec.drsAutomationLevel.get());
      }

      if (_editClusterSpec.vmDrsAutomationLevelEnabled.isAssigned()) {
         editPage.setVmDrsAutomationLevelEnabled(
               _editClusterSpec.vmDrsAutomationLevelEnabled.get()
            );
      }

      //TODO: implement missing DRS settings - migration threshold,
      // power management automation level, advanced options

      // Check Advanced DRS Policies optionBoxes if specified
      if (_editClusterSpec.advancedEnforceEvenDistributionEnabled.isAssigned()) {
         editPage.setEnforceEvenDistribution(_editClusterSpec.advancedEnforceEvenDistributionEnabled.get());
      }
      if (_editClusterSpec.advancedConsumedMemoryEnabled.isAssigned()) {
         editPage.setConsumedMemory(_editClusterSpec.advancedConsumedMemoryEnabled.get());
      }
      if (_editClusterSpec.advancedCPUOverCommitmentEnabled.isAssigned()) {
         editPage.setCpuOvercommitment(_editClusterSpec.advancedCPUOverCommitmentEnabled.get());
         if (_editClusterSpec.advancedCPUOverCommitmentEnabled.get() &&
               _editClusterSpec.advancedCPUOverCommitmentValue.isAssigned()) {
            editPage.setCpuOvercommitmentValue(_editClusterSpec.advancedCPUOverCommitmentValue.get());
         }
      }

      boolean editPageIsClosed = editPage.clickOk();
      verifyFatal(editPageIsClosed, "Verify edit page is closed");
   }
}
