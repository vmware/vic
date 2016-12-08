/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UpdateBackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.BackingTagsBasicSrvApi;

/**
 * Common workflow step that updates backing tag via API.
 */
public class UpdateBackingTagByApiStep extends BaseWorkflowStep {

   private UpdateBackingTagSpec _updateBackingTagSpec;

   @Override
   public void prepare() throws Exception {
      _updateBackingTagSpec = getSpec().get(UpdateBackingTagSpec.class);
      if (_updateBackingTagSpec == null) {
         throw new IllegalArgumentException(
               "The update backing tag spec expects UpdateBackingTagSpec!");
      }
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _updateBackingTagSpec = filteredWorkflowSpec.get(UpdateBackingTagSpec.class);
      if (_updateBackingTagSpec == null) {
         throw new IllegalArgumentException(
               "The update backing tag spec expects UpdateBackingTagSpec!");
      }
   }

   @Override
   public void execute() throws Exception {
      if (!BackingTagsBasicSrvApi.getInstance().updateBackingTag(_updateBackingTagSpec)) {
         throw new Exception(
               String.format(
                     "Unable to update backing tag: '%s'",
                     _updateBackingTagSpec.targetTag.get().name.get()
                     ));
      }
   }

   @Override
   public void clean() throws Exception {
      BackingTagsBasicSrvApi.getInstance().deleteBackingTagSafely(
            _updateBackingTagSpec.category.get(),
            _updateBackingTagSpec.newTargetConfigs.get()
            );
   }
}
