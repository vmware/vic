/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.workflow;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

import com.google.common.base.Strings;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * A container of the steps in a workflow.
 * @deprecated see the TestWorkflowController
 */
@Deprecated
public class WorkflowComposition {

   private final TestScope _defaultTestScope;

   /**
    * Creates a <code>WorkflowComposition</code> instance
    *
    * @param defaultStepsTestScope The default test scope that will be
    *       assigned to steps for which no explicit test scope is specified.
    */
   public WorkflowComposition(TestScope defaultStepsTestScope) {
      _defaultTestScope = defaultStepsTestScope;
   }

   /**
    * A map between a property from the workflow test spec and
    * the workflow step spec.
    *
    * Provides API for syncing data of the destination property with
    * the data of the source property.
    */
   public class SpecParamMap {

      /**
       * A reference to the source property.
       */
      @SuppressWarnings("rawtypes")
      public DataProperty src;

      /**
       * A reference to the destination property.
       */
      @SuppressWarnings("rawtypes")
      public DataProperty dst;

      /**
       * Copy the data from the source property to the
       * destination property.
       */
      @SuppressWarnings("unchecked")
      public void syncData() {
         if (src == null || dst == null) {
            return;
         }

         dst.set(src.get());
      }
   }

   private final List<WorkflowStep> _workflowSteps = new ArrayList<WorkflowStep>();

   private final HashMap<WorkflowStep, List<SpecParamMap>> _inputParams =
         new HashMap<WorkflowStep, List<SpecParamMap>>();

   private final HashMap<WorkflowStep, List<SpecParamMap>> _outputParams =
         new HashMap<WorkflowStep, List<SpecParamMap>>();


   /**
    * Return all the defined workflow steps.
    *
    * @return
    *    A <code>List</code> of all workflow steps.
    */
   public List<WorkflowStep> getAllSteps() {
      List<WorkflowStep> workflowSteps = new ArrayList<WorkflowStep>();
      workflowSteps.addAll(_workflowSteps);
      return workflowSteps;
   }

   /**
    * Return all steps that are marked to be run. A step can be skipped
    * if its <code>skip</code> flag is true.
    *
    * @return
    *    A <code>List</code> of all workflow steps that will be run.
    */
   public List<WorkflowStep> getRunnableSteps() {
      List<WorkflowStep> runnableSteps = new ArrayList<WorkflowStep>();
      for (WorkflowStep step : _workflowSteps) {
         if (!step.getSkip()) {
            runnableSteps.add(step);
         }
      }
      return runnableSteps;
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * The method will do nothing if the step is already defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    *
    * @param title
    *    Optional step title.
    *
    * @param testScope
    *    Optional test scope. The step will be skipped if the test run doesn't
    *    have the required test scope.
    *
    * @param specTagsFilter
    *    Optional spec tag filter. Only specs containing at least one of the tags will
    *    be made available to the test step.
    */
   public void appendStep(
         WorkflowStep workflowStep, String title, TestScope testScope, String[] specTagsFilter) {

      if (!Strings.isNullOrEmpty(title)) {
         workflowStep.setTitle(title);
      }

      // If test scope is defined, assign it to the step, otherwise set the
      // default test scope
      if (testScope != null) {
         workflowStep.setTestScope(testScope);
      } else {
         workflowStep.setTestScope(_defaultTestScope);
      }

      if (specTagsFilter != null) {
         workflowStep.setSpecTagsFilter(specTagsFilter);
      }

      if (!_workflowSteps.contains(workflowStep)) {
         _workflowSteps.add(workflowStep);
      }
   }


   /**
    * Append a workflow step at the end of the list.
    *
    * The method will do nothing if the step is already defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    *
    * @param title
    *    Optional step title.
    *
    * @param
    *    Optional test scope. The step will be skipped if the test run doesn't
    *    have the required test scope.
    */
   public void appendStep(
         WorkflowStep workflowStep, String title, TestScope testScope) {
      appendStep(workflowStep, title, testScope, null);
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * The method will do nothing if the step is already defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    *
    * @param title
    *    Optional step title.
    */
   public void appendStep(WorkflowStep workflowStep, String title) {
      appendStep(workflowStep, title, null, null);
   }

   /**
    * Append a workflow step at the end of the list.
    *
    * The method will do nothing if the step is already defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    */
   public void appendStep(WorkflowStep workflowStep) {
      appendStep(workflowStep, null, null, null);
   }


   /**
    * Remove a workflow step from te list. Do nothing if the step is
    * not defined in the list.
    *
    * @param workflowStep
    *    A workflow step.
    */
   public void removeStep(WorkflowStep workflowStep) {
      if (_workflowSteps.contains(workflowStep)) {
         _workflowSteps.remove(workflowStep);
      }
   }


   /**
    * Create a mapping between two <code>DataProperty</code> parameters. The
    * value of the source parameter will be transfered to the destination
    * parameter before step execution stage begins.
    *
    * The method is used to specify data to be transfer from the local
    * test spec to the local spec of a common step shared between many workflows.
    *
    * @param workflowStep
    *    The workflow step on which the param map is applied.
    * @param src
    *    Source <code>DataProperty</code> parameter.
    * @param dst
    *    Destination <code>DataProperty</code> parameter.
    */
   public void mapInputParam(
         WorkflowStep workflowStep,
         @SuppressWarnings("rawtypes") DataProperty src,
         @SuppressWarnings("rawtypes") DataProperty dst) {

      List<SpecParamMap> paramMapList =
            _inputParams.get(workflowStep);
      if (paramMapList == null) {
         paramMapList = new ArrayList<SpecParamMap>();
         _inputParams.put(workflowStep, paramMapList);
      }

      SpecParamMap paramMap = new SpecParamMap();
      paramMap.src = src;
      paramMap.dst = dst;
      paramMapList.add(paramMap);
   }

   /**
    * Create a mapping between two <code>DataProperty</code> parameters. The
    * value of the source parameter will be transfered to the destination
    * parameter after step execution stage completes.
    *
    * The method is used to specify data to be transfered from the local spec
    * of a common step shared between many workflows to the local test spec.
    *
    * @param workflowStep
    *    The workflow step on which the param map is applied.
    * @param src
    *    Source <code>DataProperty</code> parameter.
    * @param dst
    *    Destination <code>DataProperty</code> parameter.
    */
   public void mapOutputParam(
         WorkflowStep workflowStep,
         @SuppressWarnings("rawtypes") DataProperty src,
         @SuppressWarnings("rawtypes") DataProperty dst) {

      List<SpecParamMap> paramMapList =
            _outputParams.get(workflowStep);
      if (paramMapList == null) {
         paramMapList = new ArrayList<SpecParamMap>();
         _outputParams.put(workflowStep, paramMapList);
      }

      SpecParamMap paramMap = new SpecParamMap();
      paramMap.src = src;
      paramMap.dst = dst;
      paramMapList.add(paramMap);

   }

   /**
    * Remove all parameters specified for the the workflow step.
    *
    * @param workflowStep
    *    Workflow step.
    */
   public void unMapParams(WorkflowStep workflowStep) {
      _inputParams.put(
            workflowStep,new ArrayList<SpecParamMap>());
      _outputParams.put(
            workflowStep, new ArrayList<SpecParamMap>());
   }

   /**
    * Return all the input parameter mappings for a workflow step.
    *
    * @param workflowStep
    *    A workflow step.
    *
    * @return
    *    A <code>List</code> of <code>SpecParamMap</code> items.
    */
   public List<SpecParamMap> getInputParamsMap(WorkflowStep workflowStep) {
      return _inputParams.get(workflowStep);
   }

   /**
    * Return all the output parameter mappings for a workflow step.
    *
    * @param workflowStep
    *    A workflow step.
    *
    * @return
    *    A <code>List</code> of <code>SpecParamMap</code> items.
    */
   public List<SpecParamMap> getOutputParamsMap(WorkflowStep workflowStep) {
      return _outputParams.get(workflowStep);
   }

}
