/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotEmpty;

import java.util.List;
import java.util.ArrayList;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Step for reconfiguring a managed entity. This class should be inherited when
 * an entity will be reconfigured during a test
 *
 * Extend the step when it is needed to reconfigure an entity. During the
 * preparation phase (in 'prepare' method) all configurations are found that
 * also have assigned the property 'reconfigurableConfigVersion' for the
 * entitySpec. Each class that extends this step need to additionally filter the
 * list of old and new configuration specs, in order to get the spec for the
 * correct entity type.
 */
public abstract class ReconfigureManagedEntityStep extends CommonUIWorkflowStep {

   protected List<ManagedEntitySpec> _originalManagedEntitySpecs;
   protected List<ManagedEntitySpec> _reconfiguredManagedEntitySpecs;

   /**
    * Method to get the new spec of entity that will be applied
    *
    * @param specClass
    *           - class of the spec of the entity
    */
   public ManagedEntitySpec getNewSpec(
         Class<? extends ManagedEntitySpec> specClass) {

      for (ManagedEntitySpec _newSpec : _reconfiguredManagedEntitySpecs) {
         if (_newSpec.getClass().equals(specClass)) {
            return _newSpec;
         }
      }

      return null;
   }

   /**
    * Method to get the old spec of entity that will be edited
    *
    * @param specClass
    *           - class of the spec of the entity
    */
   public ManagedEntitySpec getOldSpec(
         Class<? extends ManagedEntitySpec> specClass) {

      for (ManagedEntitySpec _oldSpec : _originalManagedEntitySpecs) {
         if (_oldSpec.getClass().equals(specClass)) {
            return _oldSpec;
         }
      }

      return null;
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      _reconfiguredManagedEntitySpecs = new ArrayList<ManagedEntitySpec>();
      _originalManagedEntitySpecs = new ArrayList<ManagedEntitySpec>();

      for (ManagedEntitySpec entitySpec : filteredWorkflowSpec
            .getAll(ManagedEntitySpec.class)) {

         if (entitySpec.reconfigurableConfigVersion.isAssigned()) {
            if (entitySpec.reconfigurableConfigVersion
                  .get()
                  .getValue()
                  .equals(
                        ManagedEntitySpec.ReconfigurableConfigSpecVersion.NEW
                              .getValue())) {
               _reconfiguredManagedEntitySpecs.add(entitySpec);

            }

            if (entitySpec.reconfigurableConfigVersion
                  .get()
                  .getValue()
                  .equals(
                        ManagedEntitySpec.ReconfigurableConfigSpecVersion.OLD
                              .getValue())) {
               _originalManagedEntitySpecs.add(entitySpec);

            }
         }

      }

      ensureNotEmpty(_originalManagedEntitySpecs, String.format(
            "%s spec for original configuration is missing.",
            ManagedEntitySpec.class));
      ensureNotEmpty(_reconfiguredManagedEntitySpecs, String.format(
            "%s spec for reconfigured configuration is missing.",
            ManagedEntitySpec.class));

   }

}
