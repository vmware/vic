/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.edit.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.List;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HddSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.vm.lib.view.VmSummaryTabView;

/**
 * Step for verifying that the edit hdd configuration of the VM is applied under
 * the VM Summary Tab
 */
public class VerifyAddHddStep extends CommonUIWorkflowStep {

   private CustomizeHwVmSpec _vmSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vmSpec = filteredWorkflowSpec.get(CustomizeHwVmSpec.class);

      ensureNotNull(_vmSpec, "VmSpec object is missing.");
   }

   /** Generate expected capacity value in UI
    * @param capacity - value of capacity
    * @param capacityType - GB, MB
    * @return
    */
   private String generateExpectedCapacity(String capacity,
         String capacityType) {
      StringBuilder newString = new StringBuilder();
      newString.append(capacity);
      newString.append(" ");
      newString.append(capacityType);
      return newString.toString();

   }

   @Override
   public void execute() throws Exception {

      VmSummaryTabView summaryTab = new VmSummaryTabView();
      summaryTab.expandVmHardwarePortlet();
      int disksNumber = _vmSpec.hddList.getAll().size();

      List<HddSpec> hddSpecs = _vmSpec.hddList.getAll();
      // In the spec for added hdd there is only 1 hdd - the added one
      HddSpec hddSpecAdded = _vmSpec.hddList.getAll().get(0);
      String hddAdded = "", hddAddedCapacity = "";
      // Verify values of hdd spec for added hdd
      hddAdded = summaryTab.getHddLabel(hddSpecs.indexOf(hddSpecAdded));
      hddAddedCapacity = summaryTab
            .getHddCapacityLabel(hddSpecs.indexOf(hddSpecAdded));

      verifySafely(TestScope.BAT, hddSpecAdded.name.get().equals(hddAdded),
            "Hdd added is available. Expected: "
                  + hddSpecAdded.getPropertyFields() + " but found: "
                  + hddAdded);
      String expectedCapacity = generateExpectedCapacity(hddSpecAdded.hddCapacity.get(), hddSpecAdded.hddCapacityType.get());
      verifySafely(TestScope.BAT, expectedCapacity.equals(hddAddedCapacity),
            "Hdd added capacity is correct. Expected: " + expectedCapacity
                  + " but found: " + hddAddedCapacity);

      if(_vmSpec.datastore.isAssigned()){
         String hddDatastoreLabel = summaryTab.getHddDatastoreLabel(hddSpecs.indexOf(hddSpecAdded));
         String expectedHddDsLabel = _vmSpec.datastore.get().name.get();

         verifySafely(TestScope.UI, expectedHddDsLabel.equals(hddDatastoreLabel),
               "Hdd datastore label is correct. Expected: " + expectedHddDsLabel
                     + " but found: " + hddDatastoreLabel);
      }

   }

}