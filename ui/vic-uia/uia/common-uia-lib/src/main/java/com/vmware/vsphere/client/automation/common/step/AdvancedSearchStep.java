/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.view.AdvancedSearchView;
import com.vmware.vsphere.client.automation.srv.common.spec.AdvancedSearchSpec;

/**
 * Executes an "advanced search" for the given search spec. The step enters
 * search criteria values and clicks the search button. Then the results are
 * loaded in the results view.
 */
public class AdvancedSearchStep extends CommonUIWorkflowStep {

   private AdvancedSearchSpec _searchSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _searchSpec = filteredWorkflowSpec.get(AdvancedSearchSpec.class);

      if (_searchSpec == null) {
         throw new IllegalArgumentException(
               "AdvancedSearchSpec object is missing.");
      }
   }

   @Override
   public void execute() throws Exception {
      AdvancedSearchView searchView = new AdvancedSearchView();
      searchView.selectSearchForType(_searchSpec.entityType.get());

      if (_searchSpec.propertyName.isAssigned()) {
         searchView.selectCriteriaProperty(0, _searchSpec.propertyName.get());
      }

      if (_searchSpec.operator.isAssigned()) {
         searchView.selectOperator(0, _searchSpec.operator.get());
      }

      if (_searchSpec.propertyValue.isAssigned()) {
         searchView.setValue(0, _searchSpec.propertyValue.get());
      }

      if (_searchSpec.compliance.isAssigned()) {
         searchView.selectCriteriaCompliance(0, _searchSpec.compliance.get());
      }

      searchView.clickSearchButton();
   }
}
