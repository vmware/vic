/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.BackingTagsBasicSrvApi;

/**
 * Creates an backing category
 * TODO: make it work with multiple specs
 */
public class CreateBackingCategoryByApiStep extends BaseWorkflowStep {

   private BackingCategorySpec _backingCategorySpec;

   @Override
   public void prepare() {
      _backingCategorySpec = getSpec().links.get(BackingCategorySpec.class);

      if (_backingCategorySpec == null) {
         throw new IllegalArgumentException("BackingCategorySpec is missing.");
      }
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(getTestScope(), BackingTagsBasicSrvApi.getInstance().createBackingCategory(_backingCategorySpec),
            "Verifying backing category creation");
   }

   @Override
   public void clean() throws Exception {
      verifySafely(getTestScope(), BackingTagsBasicSrvApi.getInstance().deleteBackingCategorySafely(_backingCategorySpec),
            "Verify backing category deleted");
   }

}
