/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.BackingTagsBasicSrvApi;

/**
 * Detach host/cluster resources from tags
 */
public class DetachResourcesFromBackingTagsByApiStep extends BaseWorkflowStep {

   private BackingCategorySpec _category;
   private List<BackingTagSpec> _tags;

   @Override
   public void prepare() {
      _category = getSpec().get(BackingCategorySpec.class);
      _tags = getSpec().getAll(BackingTagSpec.class);

      if (_category == null) {
         throw new IllegalArgumentException("BackingCategorySpec object is missing.");
      }
   }

   @Override
   public void execute() throws Exception {
      for (BackingTagSpec tag: _tags) {
         List<ManagedEntitySpec> entities = new ArrayList<ManagedEntitySpec>();
         entities.addAll(tag.taggedObjects.getAll());
         BackingTagsBasicSrvApi.getInstance().detachResources(_category, tag, entities);
      }
   }

   @Override
   public void clean() throws Exception {
      for (BackingTagSpec tag: _tags) {
         List<ManagedEntitySpec> entities = new ArrayList<ManagedEntitySpec>();
         entities.addAll(tag.taggedObjects.getAll());
         BackingTagsBasicSrvApi.getInstance().attachResources(_category, tag, entities);
      }
   }
}
