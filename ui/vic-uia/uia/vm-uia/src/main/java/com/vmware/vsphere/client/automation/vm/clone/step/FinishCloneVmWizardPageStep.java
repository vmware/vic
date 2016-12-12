/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.clone.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;
import com.vmware.vsphere.client.automation.vm.clone.CloneVmFlowStep;
import com.vmware.vsphere.client.automation.vm.clone.spec.CloneVmSpec;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;

/**
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.clone.step.FinishCloneVmWizardPageStep}
 */
@Deprecated
public class FinishCloneVmWizardPageStep extends CloneVmFlowStep {

   private CloneVmSpec _cloneVmSpec;
   private TaskSpec _cloneVmTaskSpec;
   private CustomizeHwVmSpec _originalVmSpec;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _cloneVmSpec = filteredWorkflowSpec.get(CloneVmSpec.class);
      _originalVmSpec = filteredWorkflowSpec.get(CustomizeHwVmSpec.class);
      ensureNotNull(_cloneVmSpec, "cloneVmSpec spec is missing.");
      ensureNotNull(_originalVmSpec, "originalVmSpec spec is missing.");

      // Spec for the clone VM task
      _cloneVmTaskSpec = new TaskSpec();
      _cloneVmTaskSpec.name.set(VmUtil.getLocalizedString("task.cloneVm.name"));
      _cloneVmTaskSpec.status.set(TaskSpec.TaskStatus.COMPLETED);
      _cloneVmTaskSpec.target.set(_originalVmSpec);
   }

   @Override
   public void execute() throws Exception {
      new WizardNavigator().waitForLoadingProgressBar();
      boolean finishWizard = new WizardNavigator().finishWizard();
      verifyFatal(finishWizard, "Verify wizard is closed");
      // Wait for recent task to complete
      boolean isTaskFound = new TasksUtil().waitForRecentTaskToMatchFilter(
            new RecentTaskFilter(_cloneVmTaskSpec));

      verifyFatal(isTaskFound,
            String.format("Verifying task '%s' has reached status '%s'",
                  _cloneVmTaskSpec.name.get(), _cloneVmTaskSpec.status.get()));
   }

   /**
    * Delete the newly created vm and log if cleanup is not successful.
    */
   @Override
   public void clean() throws Exception {
      _logger.info("VM to clean after the step: " + _cloneVmSpec.name.get());
      if (VmSrvApi.getInstance().isVmPoweredOn(_cloneVmSpec)) {
         VmSrvApi.getInstance().powerOffVm(_cloneVmSpec);
      }
      VmSrvApi.getInstance().deleteVmSafely(_cloneVmSpec);
   }
}
