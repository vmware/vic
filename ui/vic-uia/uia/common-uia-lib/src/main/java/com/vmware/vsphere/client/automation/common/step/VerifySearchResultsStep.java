/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.lang.NotImplementedException;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.common.view.BaseSearchView;
import com.vmware.vsphere.client.automation.srv.common.spec.AdvancedSearchSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

/**
 * Verify that an entity is in the results
 */
public class VerifySearchResultsStep extends CommonUIWorkflowStep {

   private List<ManagedEntitySpec> _entities;
   private List<ManagedEntitySpec> _negativeEntities;
   private Class<?> _entityClass = null;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      AdvancedSearchSpec searchSpec = filteredWorkflowSpec
            .get(AdvancedSearchSpec.class);

      if (searchSpec != null && searchSpec.searchResults.isAssigned()) {
         _entities = searchSpec.searchResults.getAll();
      } else {
         _entities = filteredWorkflowSpec.getAll(ManagedEntitySpec.class);
      }

      if (searchSpec != null && searchSpec.negativeResults.isAssigned()) {
         _negativeEntities = searchSpec.negativeResults.getAll();
      } else {
         _negativeEntities = new ArrayList<ManagedEntitySpec>();
      }

      if (_entities.isEmpty() && _negativeEntities.isEmpty()) {
         throw new IllegalArgumentException(
               "ManagedEntitySpec object is missing.");
      }

      // Assert that all entities are of the same class
      List<ManagedEntitySpec> allEntities = new ArrayList<ManagedEntitySpec>();
      allEntities.addAll(_entities);
      allEntities.addAll(_negativeEntities);
      for (ManagedEntitySpec entity : allEntities) {
         if (_entityClass == null) {
            _entityClass = entity.getClass();
         } else if (!entity.getClass().equals(_entityClass)) {
            throw new IllegalArgumentException(
                  "All entities must be of the same class.");
         }
      }
   }

   @Override
   public void execute() throws Exception {
      // Switch to the appropriate tab based on the searched entity class
      String tabName;
      if (_entityClass.equals(VmSpec.class)) {
         tabName = CommonUtil.getLocalizedString("tab.vms");
      } else if (_entityClass.equals(ClusterSpec.class)) {
         tabName = CommonUtil.getLocalizedString("tab.clusters");
      } else if (_entityClass.equals(HostSpec.class)) {
         tabName = CommonUtil.getLocalizedString("tab.hosts");
      } else {
         throw new NotImplementedException("Can not handle entities of class: "
               + _entityClass.toString());
      }

      BaseSearchView searchView = new BaseSearchView();

      // Verify the positive entities
      for (ManagedEntitySpec entity : _entities) {
         int idx = searchView.getItemIndex(tabName, entity.name.get());
         verifyFatal(getTestScope(), idx >= 0, "Found " + entity.name.get()
               + " in " + tabName);
      }

      // Verify the negative entities
      for (ManagedEntitySpec entity : _negativeEntities) {
         int idx = searchView.getItemIndex(tabName, entity.name.get());
         verifyFatal(getTestScope(), idx < 0,
               "Could not find entity that we are not supposed to: "
                     + entity.name.get() + " in " + tabName);
      }
   }
}
