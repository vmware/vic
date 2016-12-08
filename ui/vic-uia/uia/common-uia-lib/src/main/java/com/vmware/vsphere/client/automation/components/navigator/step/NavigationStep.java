/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * NavigationStep is a step similar to NGCNavigateStep but with the following diffs:
 * 1. 'create' method is removed.
 * 2. '_locationSpec' is extracted from the test spec (i.e. via getSpec()).
 *
 * These changes are introduced in order to make the step similar in use as all the
 * regular steps. Tagging is successfully applied now via the framework mechanics
 * applied when calling 'getSpec()'.
 */
public class NavigationStep extends BaseWorkflowStep {
   private LocationSpec _locationSpec;

   @Override
   public void prepare() throws Exception {
      _locationSpec = getSpec().get(LocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required VcdLocationSpec is not set.");
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

      _locationSpec = filteredWorkflowSpec.get(LocationSpec.class);

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required VcdLocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }
}
