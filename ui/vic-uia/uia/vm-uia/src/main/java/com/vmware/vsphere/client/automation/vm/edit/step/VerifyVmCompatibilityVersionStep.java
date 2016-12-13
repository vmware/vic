/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.text.MessageFormat;

import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.vm.edit.spec.UpgradeVmSpec;

/**
 * Step for verifying that the VM compatibility version is correct under the VM Summary Tab
 */
public class VerifyVmCompatibilityVersionStep extends CommonUIWorkflowStep {
   private static final String COMPATIBILITY_TEXT_ID = "summary_vmVersion_valueLbl";

   private UpgradeVmSpec _upgradeSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      super.prepare(filteredWorkflowSpec);

      _upgradeSpec = filteredWorkflowSpec.get(UpgradeVmSpec.class);

      ensureNotNull(_upgradeSpec, "VmSpec object is missing.");
      ensureAssigned(_upgradeSpec.compatibilityVersion, "Ensure that compatibilityVersion is added");
      ensureAssigned(_upgradeSpec.vmHardwareVersion, "Ensure that vmHardwareVersion is added");
   }

   @Override
   public void execute() throws Exception {
      String displayedText = UI.component.property.get(Property.TEXT, COMPATIBILITY_TEXT_ID);
      String expectationFormat = "{0} ({1})";
      String expected = MessageFormat.format(
         expectationFormat,
         _upgradeSpec.compatibilityVersion.get(),
         _upgradeSpec.vmHardwareVersion.get());
      verifySafely(new EqualsAssertion(displayedText, expected, "Verify that the VM compatibility Version is correct."));
   }
}
