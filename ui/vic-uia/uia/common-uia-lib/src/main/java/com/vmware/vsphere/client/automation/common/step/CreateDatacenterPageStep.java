/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureAssigned;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.view.CreateNewDatacenterPage;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatacenterBasicSrvApi;

/**
 * Creates new datacenter.
 * This step expects new datacenter page to be already opened.
 *
 * Operation performed by this step:
 *  1. Fill new datacenter page with data from the spec
 *  2. Submit the new datacenter page
 *  3. Verify if the page was submitted successfully
 */
public class CreateDatacenterPageStep extends BaseWorkflowStep {

   private DatacenterSpec _datacenterSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _datacenterSpec = filteredWorkflowSpec.get(DatacenterSpec.class);

      ensureNotNull(_datacenterSpec, "DatacenterSpec is missing.");
      ensureAssigned(_datacenterSpec.name,
            "DatacenterSpec does not have name property set.");
   }

   @Override
   public void execute() throws Exception {
      CreateNewDatacenterPage createPage = new CreateNewDatacenterPage();

      // Fill data in new datacenter page
      createPage.setDatacenterName(_datacenterSpec.name.get());

      // Verify if the dialog was submitted successfully
      verifyFatal(TestScope.BAT, createPage.clickOk(),
            "Verifying the create new datacenter page is submitted");
      new BaseView().waitForRecentTaskCompletion();
   }

   @Override
   public void clean() throws Exception {
      DatacenterBasicSrvApi.getInstance().deleteDatacenterSafely(_datacenterSpec);
   }
}
