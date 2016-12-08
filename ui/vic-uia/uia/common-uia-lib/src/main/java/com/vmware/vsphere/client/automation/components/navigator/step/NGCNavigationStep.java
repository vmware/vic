/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.NGCLocationSpec;

/**
 * NavigationStep is a step similar to NGCNavigateStep but with the following diffs:
 * 1. 'create' method is removed.
 * 2. '_locationSpec' is extracted from the test spec (i.e. via getSpec()).
 *
 * These changes are introduced in order to make the step similar in use as all the
 * regular steps. Tagging is successfully applied now via the framework mechanics
 * applied when calling 'getSpec()'.
 */
public class NGCNavigationStep extends CommonUIWorkflowStep {

   protected static final Logger _logger = LoggerFactory.getLogger(NGCNavigationStep.class);

   protected NGCLocationSpec _locationSpec;

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(NGCLocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required NGCLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   @Override
   public void execute() throws Exception {
      boolean navigationResult = NGCNavigator.getInstance().navigateTo(_locationSpec);

      verifyFatal(
            TestScope.FULL,
            navigationResult,
            "Verifying navigation result for: " + _locationSpec.path.get()
      );
   }

   // TestWorkflowStep methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      retrieveLocationSpec(filteredWorkflowSpec);
      _locationSpec.populateEntityName();
      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required NGCLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }

   /**
    * Retrive the location spec from the test spec. The method is invoked by the
    * prepare method that initialize the data for the step.
    * It might be implemented by the sub- classes to retrieve specific entity
    * spec.
    * @param filteredWorkflowSpec
    */
   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
      _locationSpec = filteredWorkflowSpec.get(NGCLocationSpec.class);
   }
}
