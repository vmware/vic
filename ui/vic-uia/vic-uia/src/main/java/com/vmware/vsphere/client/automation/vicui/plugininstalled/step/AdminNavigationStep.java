package com.vmware.vsphere.client.automation.vicui.plugininstalled.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.components.navigator.step.NGCNavigationStep;
import com.vmware.vsphere.client.automation.vicui.plugininstalled.spec.AdminNavigationSpec;

public class AdminNavigationStep extends NGCNavigationStep {
	@Override
	   public void prepare() throws Exception {
	      _locationSpec = getSpec().get(AdminNavigationSpec.class);

	      if (_locationSpec == null) {
	         throw new IllegalArgumentException(
	               "The required AdminNavigationSpec is not set.");
	      }

	      if (Strings.isNullOrEmpty(_locationSpec.path.get())) {
	         throw new IllegalArgumentException("The path is not set.");
	      }
	   }

	   // TestWorkflowStep methods

	   @Override
	   protected void retrieveLocationSpec(WorkflowSpec filteredWorkflowSpec) {
	      _locationSpec = filteredWorkflowSpec.get(AdminNavigationSpec.class);
	      if(_locationSpec == null) {
	         _logger.info("Prepare for navigation to the admin page.");
	         _locationSpec = new AdminNavigationSpec();
	      }
	   }
}
