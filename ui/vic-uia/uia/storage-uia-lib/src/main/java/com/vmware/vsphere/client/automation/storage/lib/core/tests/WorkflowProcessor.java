package com.vmware.vsphere.client.automation.storage.lib.core.tests;

import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Set;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.common.WorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;

/**
 * Workflow which tags all of the specs which should be visible by the Step
 */
class WorkflowProcessor implements IWorkflow {

   private static class StepData {
      public final String title;
      public final WorkflowStep step;
      public final String tag;

      public StepData(String title, WorkflowStep step, String tag) {
         this.title = title;
         this.step = step;
         this.tag = tag;
      }
   }

   private static int nameCounter = 0;

   private final BaseSpec containerSpec;
   private final List<StepData> steps = new LinkedList<>();

   public WorkflowProcessor(BaseSpec containerSpec) {
      this.containerSpec = containerSpec;
   }

   @Override
   public void append(String title, WorkflowStep step,
         BaseSpec... specsWhitelist) {
      String tag = null;
      if (specsWhitelist != null && specsWhitelist.length != 0) {
         tag = generateUniqueTagName();
         assignTags(tag, specsWhitelist);
      }

      steps.add(new StepData(title, step, tag));
   }

   /**
    * Compose all added step into the {@link WorkflowStepsSequence}
    *
    * @param flow
    */
   public void compose(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      for (StepData stepData : steps) {
         if (stepData.tag == null) {
            flow.appendStep(stepData.title, stepData.step);
         } else {
            flow.appendStep(stepData.title, stepData.step,
                  new String[] { stepData.tag });
         }
      }
   }

   private void assignTags(String tag, BaseSpec[] specsWhitelist) {
      Set<BaseSpec> diffClassSpecsWhitelist = new HashSet<>(
            containerSpec.links.getAll(BaseSpec.class));
      for (BaseSpec spec : specsWhitelist) {
         addSpecTag(spec, tag);

         // Remove all of the available specs associated with the same class
         // only the specs with are provided in the specsWhitelist should be
         // tagged
         List<?> sameClassSpec = containerSpec.links.getAll(spec.getClass());
         diffClassSpecsWhitelist.removeAll(sameClassSpec);
      }

      for (BaseSpec spec : diffClassSpecsWhitelist) {
         addSpecTag(spec, tag);
      }
   }

   private void addSpecTag(BaseSpec spec, String tag) {
      List<String> specTags = spec.tag.getAll();
      specTags.add(tag);
      spec.tag.set(specTags);
   }

   private static synchronized String generateUniqueTagName() {
      return "ProcessorTag_" + (nameCounter++);
   }
}
