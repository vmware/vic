/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.view.EntityNavigationTreeView;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Verifies if we are currently navigated to the desired entity.
 */
public class VerifyEntityNavigationStep extends BaseWorkflowStep {

   private ManagedEntitySpec _entity;

   @Override
   public void prepare() throws Exception {
      _entity = getSpec().links.get(ManagedEntitySpec.class);

      if (_entity == null) {
         throw new IllegalArgumentException("ManagedEntitySpec not found.");
      }
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(
            getTestScope(),
            new EntityNavigationTreeView().getFocusedEntityName().equals(_entity.name.get()),
            "Verifying navigation to item " + _entity.name.get() + "is successful");
   }
}
