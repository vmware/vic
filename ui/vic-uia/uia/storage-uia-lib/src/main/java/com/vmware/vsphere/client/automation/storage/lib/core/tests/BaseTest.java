package com.vmware.vsphere.client.automation.storage.lib.core.tests;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.ProvidersManifestInitializer;

/**
 * {@link NGCTestWorkflow} implementation providing basic functionality to
 * request provider published entities
 */
public abstract class BaseTest extends NGCTestWorkflow {

   private WorkflowProcessor setupProcessor;
   private WorkflowProcessor executionProcessor;

   /**
    * {@inheritDoc}
    */
   @Override
   public final void initSpec(final WorkflowSpec testSpec,
         final TestBedBridge testbedBridge) {

      List<ISpecInitializer> specInitializers = new ArrayList<>();
      setSpecInitializers(specInitializers);
      for (ISpecInitializer specInitializer : specInitializers) {
         specInitializer.initSpec(testSpec, testbedBridge);
      }

      new ProvidersManifestInitializer(this).initSpec(testSpec, testbedBridge);

      super.initSpec(testSpec, testbedBridge);

      initializeWorkflowSpecs(testSpec);
      setupProcessor = new WorkflowProcessor(testSpec);
      composeSetupSteps(setupProcessor);
      executionProcessor = new WorkflowProcessor(testSpec);
      composeExecutionSteps(executionProcessor);
   }

   /**
    * Push specs to {@link WorkflowSpec}
    */
   protected void initializeWorkflowSpecs(final WorkflowSpec testSpec) {
      // Default implementation has no specs to push
   }

   protected void setSpecInitializers(
         final List<ISpecInitializer> specInitializers) {
      // Default implementation has no spec initializers
   }

   // TODO: Finalize this once all the derived classes are migrated to use the
   // new compose methods
   @Override
   public void composePrereqSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composePrereqSteps(flow);
      setupProcessor.compose(flow);
   }

   // TODO: Finalize this once all the derived classes are migrated to use the
   // new compose methods
   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      super.composeTestSteps(flow);
      executionProcessor.compose(flow);
   }

   // TODO: Abstract this once all the derived classes are migrated to use the
   // new compose methods
   protected void composeSetupSteps(IWorkflow workflow) {
   }

   // TODO: Abstract this once all the derived classes are migrated to use the
   // new compose methods
   protected void composeExecutionSteps(IWorkflow workflow) {
   }
}
