/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.step;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.step.NavigationStep;
import com.vmware.vsphere.client.automation.components.tree.spec.TreeLocationSpec;

/**
 * This is the base tree navigation step. It extracts the TreeLocationSpec from
 * the workflow and navigates to it using NGCNavigator.
 */
public class TreeNavigationStep extends NavigationStep {

   protected static final Logger _logger = LoggerFactory
         .getLogger(TreeNavigationStep.class);

   protected TreeLocationSpec _locationSpec;

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(TreeLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required TreeLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   @Override
   public void execute() throws Exception {
      boolean navigationResult = NGCNavigator.getInstance().navigateTo(
            _locationSpec);

      verifyFatal(TestScope.FULL, navigationResult,
            "Verifying navigation result for: " + _locationSpec.path.get());
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      retrieveLocationSpec(filteredWorkflowSpec);
      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required TreeLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   /**
    * Retrive the location spec from the test spec. The method is invoked by the
    * prepare method that initialize the data for the step. It might be
    * implemented by the sub- classes to retrieve specific entity spec.
    *
    * @param filteredWorkflowSpec
    */
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(TreeLocationSpec.class);
   }
}