/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.workflow;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.exception.MultiException;

/**
 * The interface defines the API of a test workflow step.
 * @deprecated see the TestWorkflowController
 */
@Deprecated
public interface WorkflowStep {

   /**
    * Test spec setter.
    */
   public void setSpec(BaseSpec spec);

   /**
    * Test spec getter.
    *
    * It should always return a reference to the correct derived spec type.
    */
   public BaseSpec getSpec();

   /**
    * Setter of a list of <code>String</code> tags by which the spec is
    * filtered.
    */
   public void setSpecTagsFilter(String[] tagsFilter);

   /**
    * Getter of a list of <code>String</code> tags by which the spec is
    * filtered.
    */
   public String[] getSpecTagsFilter();

   /**
    * Step title setter.
    */
   public void setTitle(String title);

   /**
    * Step title getter.
    */
   public String getTitle();

   /**
    * Indicate if the step will be skipped. Set to true to skip the step.
    * By default it should be set to false.
    */
   public void setSkip(boolean skip);

   /**
    * Indicate if the step will be skipped.
    *
    * @return
    *    If true, the step will be skipped.
    */
   public boolean getSkip();

   /**
    * Indicate the run scope on which the step is executed.
    */
   public void setTestScope(TestScope testScope);

   /**
    * Indicate the run scope on which the step is executed.
    *
    * Use this indicator to check if a verification should be made or code
    * path executed based on the run scope they require.
    */
   public TestScope getTestScope();

   /**
    * Gets the accumulated causes of failure for later escalation.
    *
    * @return MultiException wrapper of failure causes
    */
   public MultiException getFailureCauses();

   /**
    * Use this method to complete the test spec initialization or other
    * data related preparation. It is called before the <code>execute</code>
    * method of the same step.
    *
    * The implementation of this method is optional.
    *
    * Note: Do not use it to perform mutation operations.
    */
   public void prepare() throws Exception;

   /**
    * Use this method to perform the actual execution and, if required,
    * verification of the test step. It is called after the
    * <code>prepare</code> method of the same step.
    *
    * The implementation of this method is required.
    */
   public void execute() throws Exception;


   /**
    * Use this method to revert the changes committed to the backend in
    * the <code>execute</code> method of the same step.
    *
    * The implementation of this method is required only of the
    * <code>execute</code> method is design to commit a change to the
    * backend.
    *
    * The method is called in the clean-up phase backwards on each workflow
    * step.
    */
   public void clean() throws Exception;

}
