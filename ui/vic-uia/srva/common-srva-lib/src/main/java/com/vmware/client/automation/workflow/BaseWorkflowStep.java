/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.workflow;

import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.exception.MultiException;

/**
 * Base implementation of a <code>WorkflowStep</code>. All workflow steps
 * should inherit this class.
 *
 * Only the <code>execute</code> method is required to be implemented in the
 * derived class. Unless it is a very simple inline step, use of the
 * <code>prepare</code> to retrieve the spec required for the step.
 *
 * NOTE: The class inherits the FakeTestWorkflowStep which holds the
 * implementations for the new Test Workflow model. Here should be set code that
 * is working only with the old work flow model. Once the tests are migrate3d to
 * the new model this class will be deleted.
 */
public abstract class BaseWorkflowStep extends FakeTestWorkflowStep
      implements WorkflowStep/** Legacy test workflow */ {

   protected static final Logger _logger =
         LoggerFactory.getLogger(BaseWorkflowStep.class);

   private String _title = "";

   private String[] _specTagsFilter;

   private BaseSpec _spec;

   private boolean _skip = false;

   private TestScope _testScope = TestScope.FULL;

   // Accumulates caught exceptional events from the current step, to be
   // escalated at the end of the test
   private MultiException failureCauses = new MultiException();

   @Override
   public MultiException getFailureCauses() {
      failureCauses = new MultiException();
      List<RuntimeException> errorList = super.getFailedValidations();
      for (RuntimeException runtimeException : errorList) {
         failureCauses.add(runtimeException);
      }
      return failureCauses;
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

   @Override
   public void setSkip(boolean skip) {
      _skip = skip;
   }

   @Override
   public boolean getSkip() {
      return _skip;
   }

   @Override
   public void setTestScope(TestScope runScope) {
      _testScope = runScope;
   }

   @Override
   public TestScope getTestScope() {
      return _testScope;
   }

   @Override
   public void prepare() throws Exception {
      //The method implementation here should stay empty.
   }

   /**
    * {@inheritDoc}
    *
    * This method should always be implemented in the derived class.
    */
   @Override
   public abstract void execute() throws Exception;

   @Override
   public void clean() throws Exception {
      //The method implementation here should stay empty.
   }

   /**
    * {@inheritDoc}
    *
    * Make validation which workflow model is used.
    * If the workflo model is the TestWorkflow model the hasTestScope of the
    * super class is invoked.
    */
   @Override
   protected boolean hasTestScope(TestScope requiredTestScope) {
      // when _spec is set we are in the old workflow model
      if(_spec != null) {
         return requiredTestScope.getScopeNumber() <= getTestScope().getScopeNumber();
      }
      return super.hasTestScope(requiredTestScope);
   }
}
