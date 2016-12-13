/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import org.apache.commons.lang.RandomStringUtils;
import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VAppSrvApi;

/**
 * Test the vApp API operations exposed in vAppSrvApi class. The test scenario
 * is: 1. Login to VC and check that the vApp specified in the spec does not
 * exist. 2. Create the vApp and check it is created 3. Delete the vApp
 */
public class VappOperationsTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);
      DatacenterSpec dcSpec = new DatacenterSpec();
      dcSpec.name.set(testBed.getCommonDatacenterName());
      ClusterSpec clSpec = SpecFactory.getSpec(ClusterSpec.class,
            testBed.getCommonClusterName(), dcSpec);
      HostSpec hostSpec = HostUtil.buildHostSpec(testBed.getCommonHost(),
            testBed.getESXAdminUsername(), testBed.getESXAdminPasssword(), 443,
            clSpec);
      VappSpec vApp = new VappSpec();
      vApp.name.set("VappApiTest_" + RandomStringUtils.randomAlphanumeric(5));
      // TODO Here we cannot make sure to get a host that is standalone or is in
      // a DRS enabled cluster - framework limitation (precondition for creating
      // a vApp)
      vApp.parent.set(hostSpec);
      VmSpec vm = SpecFactory.getSpec(VmSpec.class, "VappApiTestVM_"
            + RandomStringUtils.randomAlphanumeric(5), hostSpec);
      vm.guestId.set("winVistaGuest");
      vApp.vmList.set(vm);
      spec.links.add(vApp);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Nothing to do
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Validate check existence vApp operation.
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            VappSpec vappSpec = getSpec().links.get(VappSpec.class);

            boolean vappExists = VAppSrvApi.getInstance().checkVAppExists(vappSpec);
            verifyFatal(TestScope.BAT, !vappExists,
                  "Verify the vApp doesn't exist in the inventory!");
         }
      }, "Validating vApp check existence operation.");

      // Validate create vApp operation.
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            VappSpec vappSpec = getSpec().links.get(VappSpec.class);

            boolean createVappOperation = VAppSrvApi.getInstance().createVApp(vappSpec);
            verifyFatal(TestScope.BAT, createVappOperation,
                  "Verify the VappSrvApi.createVapp!");

            boolean vappExists = VAppSrvApi.getInstance().checkVAppExists(vappSpec);
            verifyFatal(TestScope.BAT, vappExists,
                  "Verify the vApp exists in the inventory!");
         }
      }, "Validating create vApp operation.");

      // Validate power on / off vApp operation

      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            VappSpec vappSpec = getSpec().links.get(VappSpec.class);

            boolean powerOnVappOperation = VAppSrvApi.getInstance().powerOnVapp(vappSpec);
            verifyFatal(TestScope.BAT, powerOnVappOperation,
                  "Verify the VappSrvApi.powerOnVapp!");
            boolean powerOffVappOperation = VAppSrvApi.getInstance().powerOffVapp(vappSpec,
                  true);
            verifyFatal(TestScope.BAT, powerOffVappOperation,
                  "Verify the VappSrvApi.powerOffVapp safely!");
         }
      }, "Validating power on/ off vApp operation.");

      // Validate delete vApp operation
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            VappSpec vappSpec = getSpec().links.get(VappSpec.class);

            boolean vappDeleted = VAppSrvApi.getInstance().deleteVApp(vappSpec);
            verifyFatal(TestScope.BAT, vappDeleted,
                  "Verify the VappSrvApi.deleteVapp!");

            vappDeleted = VAppSrvApi.getInstance().deleteVAppSafely(vappSpec);
            verifyFatal(TestScope.BAT, vappDeleted,
                  "Verify the VappSrvApi.deleteVappSafely!");

            boolean vappExists = VAppSrvApi.getInstance().checkVAppExists(vappSpec);
            verifyFatal(TestScope.BAT, !vappExists,
                  "Verify the vApp is deleted!");
         }
      }, "Validating vApp delete operation.");
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }

}
