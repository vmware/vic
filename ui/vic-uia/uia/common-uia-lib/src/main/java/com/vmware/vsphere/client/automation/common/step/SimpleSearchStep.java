/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.common.view.SimpleSearchView;

/**
 * Executes a simple search for any entity given
 */
public class SimpleSearchStep extends BaseWorkflowStep {

   private ManagedEntitySpec _entity;

   @Override
   public void prepare() {
      _entity = getSpec().links.get(ManagedEntitySpec.class);

      if (_entity == null) {
         throw new IllegalArgumentException("ManagedEntitySpec object is missing.");
      }
   }

   @Override
   public void execute() throws Exception {
      SimpleSearchView searchView = new SimpleSearchView();
      searchView.setSearchText(_entity.name.get());
      searchView.clickSearchButton();
   }
}
