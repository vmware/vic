/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.common.DatacenterGlobalActions;
import com.vmware.vsphere.client.automation.common.view.CreateNewDatacenterPage;

/**
 * Launch new datacenter page.
 * This step expects a specific view to be opened
 * vCenter -> datacenters list view
 *
 * Operations performed by this step:
 *  1. Click on create new datacenter button
 *  2. Verify that create new datacenter page is opened
 */
public class LaunchNewDatacenterPageStep extends CommonUIWorkflowStep {

   @Override
   public void execute() throws Exception {

      // Launch create new datacenter page
      ActionNavigator.invokeFromActionsMenu(IDGroup
            .toIDGroup(DatacenterGlobalActions.AI_CREATE_DATACENTER));

      CreateNewDatacenterPage createPage = new CreateNewDatacenterPage();
      createPage.waitForDialogToLoad();

      // Verify that create new datacenter page is opened
      verifyFatal(TestScope.BAT, createPage.isOpen(),
            "Create new datacenter page is opened");
   }
}
