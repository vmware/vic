/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ResourcePoolBasicSrvApi;

/**
 * Test the Resource Pool API operations exposed in ResourcePoolSrvApi class.
 * The test scenario is:<br>
 * 1. Login to VC and check that the resource pool specified in the spec does
 * not exist.<br>
 * 2. Create the resource pool and check it is created 3. Delete the resource
 * pool.<br>
 */
public class ResourcePoolOperationsTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      DatacenterSpec dcSpec = SpecFactory.getSpec(DatacenterSpec.class,
            testBed.getCommonDatacenterName(), null);

      ClusterSpec clSpec = SpecFactory.getSpec(ClusterSpec.class,
            testBed.getCommonClusterName(), dcSpec);

      ResourcePoolSpec resPoolSpec = SpecFactory.getSpec(
            ResourcePoolSpec.class, clSpec);

      spec.links.add(resPoolSpec);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Nothing to do
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Validate resource pool check existence operation
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            ResourcePoolSpec resPoolSpec = getSpec().links
                  .get(ResourcePoolSpec.class);

            boolean resPoolExists = ResourcePoolBasicSrvApi.getInstance()
                  .checkResourcePoolExists(resPoolSpec);
            verifyFatal(TestScope.BAT, !resPoolExists,
                  "Verify the resource pool doesn't exist.");
         }
      }, "Validating resource pool check existence operation.");

      // Validate create resource pool operation
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            ResourcePoolSpec resPoolSpec = getSpec().links
                  .get(ResourcePoolSpec.class);

            boolean createResPoolOperation = ResourcePoolBasicSrvApi.getInstance()
                  .createResourcePool(resPoolSpec);
            verifyFatal(TestScope.BAT, createResPoolOperation,
                  "Verify the ResourcePoolSrvApi.getInstance().createResourcePool operation.");

            boolean resPoolExists = ResourcePoolBasicSrvApi.getInstance()
                  .checkResourcePoolExists(resPoolSpec);
            verifyFatal(TestScope.BAT, resPoolExists,
                  "Verify the resource pool is created.");
         }
      }, "Validating create resource pool operation.");

      // Validate delete resource pool operation
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            ResourcePoolSpec resPoolSpec = getSpec().links
                  .get(ResourcePoolSpec.class);

            boolean resPoolDeleted = ResourcePoolBasicSrvApi.getInstance()
                  .deleteResourcePool(resPoolSpec);
            verifyFatal(TestScope.BAT, resPoolDeleted,
                  "Verify the ResourcePoolSrvApi.getInstance().deleteResourcePool operation.");

            resPoolDeleted = ResourcePoolBasicSrvApi.getInstance()
                  .deleteResourcePoolSafely(resPoolSpec);
            verifyFatal(TestScope.BAT, resPoolDeleted,
                  "Verify the ResourcePoolSrvApi.getInstance().deleteResourcePoolSafely operation.");

            boolean resPoolExists = ResourcePoolBasicSrvApi.getInstance()
                  .checkResourcePoolExists(resPoolSpec);
            verifyFatal(TestScope.BAT, !resPoolExists,
                  "Verify the resource pool is deleted.");
         }
      }, "Validating resource pool delete operation.");
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
