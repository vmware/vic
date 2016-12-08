package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpecFieldsInitializer;

/**
 * Enhanced work flow step implementation which allows successor classes to use
 * the {@link UsesSpec}
 */
public abstract class EnhancedBaseWorkflowStep extends CommonUIWorkflowStep {

   private final UsesSpecFieldsInitializer specFiledsInitializer = new UsesSpecFieldsInitializer(
         this);

   /**
    * {@inheritDoc}
    */
   @Override
   public final void prepare(WorkflowSpec filteredWorkflowSpec)
         throws Exception {

      specFiledsInitializer.initializeFields(filteredWorkflowSpec);

   }

}