/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.workflow;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.exception.MultiException;

/**
 * Base implementation of a test workflow composed of <code>WorkflowStep</code>
 * steps. It handles the composition, test spec initialization, step execution
 * and clean up of a workflow.
 *
 * The entry points are based on the TestNG's beforeClass, test method and
 * afterClass, so the test workflow has the same stages and can be launched and
 * monitored as a standard TestNG test.
 *
 * The <code>TestWorkflow</code> also implements the <code>WorkflowStep</code>
 * interface. Therefore a <code>TestWorkflow</code> can be added as sub-workflow
 * through a <code>WorkflowStep</code> in another <code>TestWorkflow</code>/
 *
 * All classes implementing end to end test should inherit this class.
 *
 * NOTE:
 *  BaseTestWorkflow            - a legacy workflow
 *  Extends TestWorkflowAdapter - Test command workflow
 */
public abstract class BaseTestWorkflow extends TestWorkflowAdapter {

   protected static final Logger _logger = LoggerFactory.getLogger(BaseTestWorkflow.class);

   private String _title = "";

   private BaseSpec _spec;

   private String[] _specTagsFilter;

   private boolean _skip = false;

   private WorkflowComposition _prereqComposition;

   private WorkflowComposition _testComposition;

   private List<WorkflowStep> _stepsToExecute;

   // Accumulates caught exceptional events, to be escalated at the end of the
   // test
   private final MultiException failureCauses = new MultiException();


   @Override
   public MultiException getFailureCauses() {
      MultiException me = failureCauses;
      return me;
   }

   @Override
   public void setSkip(boolean skip) {
      _skip = skip;
   }

   @Override
   public boolean getSkip() {
      return _skip;
   }

   @Override
   public void setSpec(BaseSpec spec) {
      _spec = spec;
   }

   @Override
   public BaseSpec getSpec() {
      return _spec;
   }

   @Override
   public void setSpecTagsFilter(String[] tagsFilter) {
      _specTagsFilter = tagsFilter;
   }

   @Override
   public String[] getSpecTagsFilter() {
      return _specTagsFilter;
   }

   @Override
   public void setTitle(String title) {
      _title = title;
   }

   @Override
   public String getTitle() {
      String title;
      if (Strings.isNullOrEmpty(_title)) {
         // Use the class name for a title.
         title = this.getClass().getSimpleName();
         if (Strings.isNullOrEmpty(title)) {
            // Use the base class name for a title. Covers the case when
            // anonymous step is used.
            title = this.getClass().getSuperclass().getSimpleName();
         }
      } else {
         // Return the title set outside.
         title = _title;
      }

      return title;
   }

   public List<WorkflowStep> getPrereqSteps() {
      return _prereqComposition.getAllSteps();
   }

   public List<WorkflowStep> getTestSteps() {
      return _testComposition.getAllSteps();
   }

   @Override
   public void setTestScope(TestScope testScope) {
      this._testScope = testScope;
   }

   @Override
   public TestScope getTestScope() {
      return _testScope;
   }

   /**
    * Initialize the test spec and build the workflow steps.
    */
   @Override
   public void prepare() {
      initSpec();

      _prereqComposition = new WorkflowComposition(_testScope);
      _testComposition = new WorkflowComposition(_testScope);

      composePrereqSteps(_prereqComposition);
      composeTestSteps(_testComposition);

      _stepsToExecute = new ArrayList<WorkflowStep>();
   }

   /**
    * Start executing the workflow.
    */
   @Override
   public void execute() throws Exception {
      run();
   }

   /**
    * Perform clean up of the workflow steps.
    */
   @Override
   public void clean() {
      doClean();
   }

   @Override
   public void beforeClass() {
      super.beforeClass();

      // Prepare the workflow test spec and steps.
      prepare();
   }

   @Override
   public void afterClass() {
      super.afterClass();

      // Clean up after the steps are complete.
      clean();
   }

   /**
    * Implement this method to initialize the workflow test spec.
    */
   public abstract void initSpec();

   /**
    * Implement this method to build the list of the workflow steps building the
    * prerequisites. They will be executed in the order they are provided.
    *
    * @param composition
    *           A reference to the steps container.
    */
   public abstract void composePrereqSteps(WorkflowComposition composition);

   /**
    * Implement this method to build the list of the workflow steps related to
    * the actual test. They will be executed in the order they are provided.
    *
    * @param composition
    *           A reference to the steps container.
    */
   public abstract void composeTestSteps(WorkflowComposition composition);

   public void skipSteps(int[] stepsIdx) {

   }

   /**
    * Run the clean up stage of all steps that were executed or scheduled for
    * execution.
    */
   private void doClean() {
      if (_stepsToExecute == null || _stepsToExecute.size() == 0) {
         // Nothing to do if no workflow steps are run.
         return;
      }

      if (!_isCleanupEnabled) {
         return;
      }

      // Clean up
      WorkflowStep step = null;
      for (int i = _stepsToExecute.size() - 1; i >= 0; --i) {
         step = _stepsToExecute.get(i);
         try {
            _logger.info(
                  String.format(
                        "Starting clean up of test workflow step %s.",
                        step.getTitle()));

            step.clean();

            _logger.info(
                  String.format(
                        "Clean up of test workflow step %s finished successfully",
                        step.getTitle()));
         } catch (Throwable t) {
            _logger.error(
                  String.format(
                        "Error cleaning up test workflow step %s.\n%s",
                        step.getTitle(),
                        t.getMessage()
                     )
               );
         }
      }
   }

   /**
    * Run sequentially all steps marked as runnable in the workflow. The method
    * invokes the prepare and execute stages of the steps.
    */
   private void run() throws Exception {
      runComposition(_prereqComposition, 0);
      runComposition(_testComposition, 5000);
   }

   /**
    * Run a composition of steps sequentially. The method invokes the
    * prepare and execute stages of the steps.
    *
    * @param composition
    *    A <code>WorkflowComposition</code>.
    *
    * @param separationTimeout
    *    Waiting timeout in milliseconds before running the execute method
    *    of the next step.
    */
   private void runComposition(WorkflowComposition composition, int separationTimeout)
         throws Exception {
      List<WorkflowStep> workflowSteps = composition.getRunnableSteps();
      if (workflowSteps.size() == 0) {
         _logger.warn("No workflow steps are defined.");
         return;
      }

      // Prepare and execute the steps
      for (WorkflowStep step : workflowSteps) {

         // Check if test step requires compatible scope with the
         // one the test in run on in order to determine if the
         // step should be executed.
         if (step.getTestScope().getScopeNumber() <= getTestScope().getScopeNumber()) {
            // Set the test scope to the step.
            step.setTestScope(getTestScope());
         } else {
            // Skip the step if it requires higher scope than the test
            // is run on.
            continue;
         }

         // Init the step spec
         step.setSpec(getDeepClonedSpec(getSpec(), step.getSpecTagsFilter()));

         // Sync the data of the input parameters
         syncParamMapData(composition.getInputParamsMap(step));

         try {
            _logger.info(
                  String.format(
                        "Preparing test workflow step %s.",
                        step.getTitle()));

            step.prepare();

            _logger.info(
                  String.format(
                        "test workflow step %s prepared successfully.",
                        step.getTitle()));

         } catch (Throwable t) {
            _logger.error(
                  String.format(
                        "Error executing prepare step %s.\n%s",
                        step.getTitle(), t));
            //TODO: Implement ITestListener that sets the outcome to SKIP in this case
            throw t;
         }

         _stepsToExecute.add(step);

         try {
            _logger.info(
                  String.format(
                        "==================== Executing test workflow step %s.",
                        step.getTitle()));

            // Call the "validate" method on the view assigned to this step
            if (step.getClass().isAnnotationPresent(View.class)) {
               // Get the "view" class
               View viewAnnotation = step.getClass().getAnnotation(View.class);
               Class<?> viewClass = viewAnnotation.value();

               _logger.warn("Calling the validate method on the associated view: "
                     + viewClass.getSimpleName());

               // Call the "validate" method
               Method validateMethod =
                     viewClass.getMethod(viewAnnotation.validateMethodName());
               try {
                  validateMethod.invoke(viewClass.newInstance());
               } catch (InvocationTargetException e) {
                  if (e.getTargetException() != null) {
                     throw new IllegalStateException("Associated view validation failed: " + viewClass.getSimpleName(),
                           e.getTargetException());
                  } else {
                     throw new IllegalStateException("Associated view validation failed: " + viewClass.getSimpleName(),
                           e);
                  }
               }
            }

            step.execute();

            // Sync the data of the output parameters
            syncParamMapData(composition.getOutputParamsMap(step));

            _logger.info(String.format(
                  "Test Workflow step %s executed successfully.",
                  step.getTitle()));

            // Gather the exceptions from the current step for later escalation
            failureCauses.add(step.getFailureCauses());
         } catch (Throwable t) {
            _logger.error(
                  String.format(
                        "Error executing test workflow step %s.\n%s",
                        step.getTitle(),
                        t));

            // Gather the exceptions from the current step and throw them
            failureCauses.add(step.getFailureCauses());
            failureCauses.add(t);
            failureCauses.ifExceptionThrow();
         }
      }

      // Throw any accumulated exceptions after all steps get
      // executed, to cause failed status of the test
      failureCauses.ifExceptionThrow();
   }

   /**
    * Sync the destination data with the source data in a spec parameter map.
    *
    * @param paramMaps
    *           A <code>List</code> of <code>SpecParamMap</code> items.
    */
   private void syncParamMapData(List<WorkflowComposition.SpecParamMap> paramMaps) {
      if (paramMaps != null) {
         for (WorkflowComposition.SpecParamMap paramMap : paramMaps) {
            paramMap.syncData();
         }
      }
   }

   /**
    * Make a deep clone of a spec including the linked specs. The routime will
    * also apply tags filter, if specified, cloning only specs that have at
    * least one matching tag with the filter. Tag filtering is not performed on
    * the root spec, which is the workflow spec.
    *
    * @param sourceSpec
    *           A reference to the spec
    * @param tagsFilter
    *           An array of tags to use as filter
    * @return Cloned spec.
    */
   private BaseSpec getDeepClonedSpec(BaseSpec sourceSpec, String[] tagsFilter)
         throws IllegalArgumentException, SecurityException, InstantiationException,
         IllegalAccessException, InvocationTargetException, NoSuchMethodException {

      return getDeepClonedSpecInternal(sourceSpec, tagsFilter, new ArrayList<BaseSpec>());
   }

   /**
    * A helper method that implements the spec cloning as described in the
    * <code>getDeepClonedSpec</code> routine.
    *
    * @param sourceSpec
    *           A reference to the spec
    * @param tagsFilter
    *           An array of tags to use as filter
    * @param processedSourceSpecs
    *           A list of the already cloned source specs. The
    *           <code>PropertyBox</code> links automatically provide references
    *           to all the parent specs, which should be excluded from this
    *           cloning.
    *
    * @return Cloned spec.
    */
   private BaseSpec getDeepClonedSpecInternal(BaseSpec sourceSpec, String[] tagsFilter,
         List<BaseSpec> processedSourceSpecs) throws IllegalArgumentException,
         SecurityException, InstantiationException, IllegalAccessException,
         InvocationTargetException, NoSuchMethodException {

      // Empty spec. there's nothing to process.
      if (sourceSpec == null) {
         return null;
      }

      // Skip already processed specs. This is necessary because spec links
      // provide
      // automatically references to their container specs.
      if (processedSourceSpecs.contains(sourceSpec)) {
         return null;
      }

      boolean cloneSpec = true;

      // Try to match tags if it's not the workflow container assuming the root
      // is always the workflow container.
      if (processedSourceSpecs.size() > 0 && requireTagMatchEval(tagsFilter)
            && !hasTagMatch(sourceSpec, tagsFilter)) {
         cloneSpec = false;
      }

      processedSourceSpecs.add(sourceSpec);

      BaseSpec newSpec = null;

      if (cloneSpec) {
         // Initialize the new spec from the correct type
         newSpec = sourceSpec.getClass().newInstance();

         // Copy the data properties
         newSpec.copy(sourceSpec);

         // Copy the linked specs
         List<BaseSpec> linkedSourceSpecs = sourceSpec.links.getAll(BaseSpec.class);
         for (BaseSpec linkedSourceSpec : linkedSourceSpecs) {

            BaseSpec newLinkedSpec =
                  getDeepClonedSpecInternal(
                        linkedSourceSpec,
                        tagsFilter,
                        processedSourceSpecs);

            if (newLinkedSpec != null) {
               newSpec.links.add(newLinkedSpec);
            }
         }
      }

      return newSpec;
   }

   /**
    * Check if at least one of the tags in the filter matches a tag in the spec.
    *
    * @param sourceSpec
    *           A spec reference
    * @param tagsFilter
    *           An array of tags to use as filter.
    *
    * @return true - at least one filter tag has a match in the spec; false -
    *         otherwise;
    */
   private boolean hasTagMatch(BaseSpec spec, String[] tagsFilter) {
      // Check if the spec has the required tags and it's not the workflow spec.
      if (spec.tag.isAssigned()) {
         boolean tagFound = false;

         List<String> specTags = spec.tag.getAll();
         for (String tagFilter : tagsFilter) {
            if (specTags.contains(tagFilter)) {
               tagFound = true;
               break;
            }
         }

         return tagFound;
      }

      return false;
   }

   /**
    * Check if evaluation for matching tags between the spec and the tags filter
    * should be performed.
    *
    * @param tagsFilter
    *           An array of tags to use as filter
    *
    * @return true - the tags filter is not empty; false - otherwise;
    */
   private boolean requireTagMatchEval(String[] tagsFilter) {
      return tagsFilter != null && tagsFilter.length > 0;
   }
}
