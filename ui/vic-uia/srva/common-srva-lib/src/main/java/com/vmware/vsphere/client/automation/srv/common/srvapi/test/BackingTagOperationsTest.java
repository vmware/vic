/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import org.apache.commons.lang.RandomStringUtils;
import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.srvapi.BackingTagsBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.step.CreateBackingCategoryAndTagsByApiStep;

/**
 * Test the Backing Tags API operations exposed in TaggingSrvApi class.
 */
public class BackingTagOperationsTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      DatacenterSpec datacenterSpec = SpecFactory.getSpec(DatacenterSpec.class,
            testBed.getCommonDatacenterName(), null);

      ClusterSpec clusterSpec = SpecFactory.getSpec(ClusterSpec.class, testBed.getCommonClusterName(), datacenterSpec);

      HostSpec hostSpec = HostUtil.buildHostSpec(testBed.getCommonHost(),
            testBed.getESXAdminUsername(), testBed.getESXAdminPasssword(), 443,
            clusterSpec);

      // Define backing category to be used in the test
      BackingCategorySpec backingCategorySpec = SpecFactory.getSpec(BackingCategorySpec.class);
      backingCategorySpec.description.set(RandomStringUtils.randomAlphanumeric(10));

      // Define backing tag to be used in the test
      BackingTagSpec backingtagSpec = SpecFactory.getSpec(BackingTagSpec.class);
      backingtagSpec.description.set(RandomStringUtils.randomAlphanumeric(10));

      spec.links.add(backingCategorySpec, backingtagSpec);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // No prereqs needed
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {
      // 1. Create a backing category and tag
      composition.appendStep(new CreateBackingCategoryAndTagsByApiStep());

      // 2. Verify category and tag creation
      composition.appendStep(new BaseWorkflowStep() {
         private BackingCategorySpec _backingCategory;
         private BackingTagSpec _backingTag;

         @Override
         public void prepare() {
            _backingCategory = getSpec().links.get(BackingCategorySpec.class);
            _backingTag = getSpec().links.get(BackingTagSpec.class);
         }

         @Override
         public void execute() throws Exception {
            verifyFatal(TestScope.BAT, BackingTagsBasicSrvApi.getInstance().checkBackingCategoryExists(_backingCategory),
                  "Verify backing category is created");
            verifyFatal(TestScope.BAT, BackingTagsBasicSrvApi.getInstance().checkBackingTagExists(_backingCategory, _backingTag),
                  "Verify backing tag is created");
         }
      });

      // 3. Delete the backing tag and category and verify they are deleted
      composition.appendStep(new BaseWorkflowStep() {
         private BackingCategorySpec _backingCategory;
         private BackingTagSpec _backingTag;

         @Override
         public void prepare() {
            _backingCategory = getSpec().links.get(BackingCategorySpec.class);
            _backingTag = getSpec().links.get(BackingTagSpec.class);
         }

         @Override
         public void execute() throws Exception {
            // Delete tag
            BackingTagsBasicSrvApi.getInstance().deleteBackingTagSafely(_backingCategory, _backingTag);
            verifyFatal(TestScope.BAT, !BackingTagsBasicSrvApi.getInstance().checkBackingTagExists(_backingCategory, _backingTag),
                  "Verify backing category tag is deleted");

            // Delete category
            BackingTagsBasicSrvApi.getInstance().deleteBackingCategorySafely(_backingCategory);
            verifyFatal(TestScope.BAT, !BackingTagsBasicSrvApi.getInstance().checkBackingCategoryExists(_backingCategory),
                  "Verify backing category is deleted");
         }
      });
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
