/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.common.view.EntityView;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Performs verification of entity name via UI.
 */
public class VerifyEntityNameByUiStep extends CommonUIWorkflowStep {

   private ManagedEntitySpec _entitySpec;
   private Class<? extends ManagedEntitySpec> _specClass;

   /**
    * Use the default constructor in combination with
    * tags to filter specific ManagedEntitySpec
    *
    * @param specClass - class of the spec of the entity.
    */
   @Deprecated
   public VerifyEntityNameByUiStep(Class<? extends ManagedEntitySpec> specClass) {
      _specClass = specClass;
   }

   public VerifyEntityNameByUiStep() {
      _specClass = ManagedEntitySpec.class;
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _entitySpec = filteredWorkflowSpec.get(_specClass);

      ensureNotNull(_entitySpec, String.format("%s spec is missing.", _specClass));
   }

   @Override
   public void execute() throws Exception {
      verifyFatal(new EqualsAssertion(
         new EntityView().getEntityName(),
         _entitySpec.name.get(),
         String.format(
            "Verifying by UI the name for %s",
            _entitySpec.name.get())));
   }
}
