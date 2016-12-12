/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.BackingTagsBasicSrvApi;

/**
 * Creates an backing tag
 */
public class CreateBackingCategoryAndTagsByApiStep extends BaseWorkflowStep {

   private BackingCategorySpec _category;
   private List<BackingTagSpec> _tags;

   @Override
   public void prepare() {
      _category = getSpec().get(BackingCategorySpec.class);
      _tags = getSpec().getAll(BackingTagSpec.class);

      if (_category == null) {
         throw new IllegalArgumentException("BackingCategorySpec is missing.");
      }
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _category = filteredWorkflowSpec.get(BackingCategorySpec.class);
      if (_category == null) {
         throw new IllegalArgumentException("BackingCategorySpec object is missing.");
      }

      _tags = filteredWorkflowSpec.getAll(BackingTagSpec.class);
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(TestScope.BAT, BackingTagsBasicSrvApi.getInstance().createBackingCategory(_category),
            "Verify create backing category");

      for (BackingTagSpec tag: _tags) {
         verifyFatal(TestScope.BAT, BackingTagsBasicSrvApi.getInstance().createBackingTag(_category, tag),
               "Verify create backing tag");
      }
   }

   @Override
   public void clean() throws Exception {
      for (BackingTagSpec tag: _tags) {
         verifySafely(TestScope.BAT, BackingTagsBasicSrvApi.getInstance().deleteBackingTagSafely(_category, tag),
               "Verify delete backing tag");
      }

      verifySafely(TestScope.BAT, BackingTagsBasicSrvApi.getInstance().deleteBackingCategorySafely(_category),
            "Verify delete backing category");
   }

}
