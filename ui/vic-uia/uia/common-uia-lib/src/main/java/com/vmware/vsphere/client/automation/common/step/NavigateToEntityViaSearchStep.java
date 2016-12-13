/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.common.view.HomeView;

/**
 * Common workflow step for navigating to entity through search
 */
public class NavigateToEntityViaSearchStep extends BaseWorkflowStep {

   protected String _name;

   @Override
   public void prepare() {
      ManagedEntitySpec entity = getSpec().links.get(ManagedEntitySpec.class);

      if (entity == null) {
         throw new IllegalArgumentException("No links to ManagedEntitySpec found");
      }

      _name = entity.name.get();
   }

   @Override
   public void execute() throws Exception {
      HomeView homeView = new HomeView();

      homeView.setSearchText(_name);
      homeView.clickSearchButton();
      boolean isClicked = homeView.clickFirstResult();
      verifyFatal(TestScope.FULL, isClicked, "Clicking on first search result");

   }
}
