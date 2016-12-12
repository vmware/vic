/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.NGCLocationSpec;

/**
 * Common workflow step for performing standard VCD navigation operations.
 * @deprecated use NavigationStep instead.
 */
@Deprecated
public class NGCNavigateStep extends BaseWorkflowStep {

   private NGCLocationSpec _locationSpec;

   // Hide the constructor. The factory method should be used instead.
   private NGCNavigateStep() {};

   /**
    * Create an instance of <code>VcdNavigateStep</code> and initialize it with
    * <code>VcdLocationSpec</code>.
    *
    * @param locationSpec
    *    A <code>VcdLocationSpec</code>.
    *
    * @return
    *    A <code>NGCNavigateStep</code>.
    */
   public static NGCNavigateStep create(NGCLocationSpec locationSpec) {
      if (locationSpec == null) {
         throw new IllegalArgumentException("Required VcdLocationSpec is not set.");
      }

      NGCNavigateStep step = new NGCNavigateStep();
      //      step.setSpec(locationSpec);
      step._locationSpec = locationSpec;

      return step;
   }

   @Override
   public void prepare() throws Exception {
      // workaround
      //      setSpec(_locationSpec);
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("NGC Navigate Step");
      }

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
      verifyFatal(TestScope.BAT, navigationResult, "Verifying navigation result for: " + _locationSpec.path.get());
      getStepTestScope();
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      if (_locationSpec == null) {
         throw new IllegalArgumentException(
               "The required LocationSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
         throw new IllegalArgumentException("The path is not set.");
      }
   }
}
