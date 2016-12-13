/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.BaseDialogNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.suitaf.util.CommonUtils;
import com.vmware.vsphere.client.automation.common.spec.DialogTiwoTitleSpec;

/**
 * Common workflow steps that can be used for maximizing TIWO dialog
 * once it is minimized.
 */
public class MaximizeTiwoDialogStep extends BaseWorkflowStep {

   DialogTiwoTitleSpec _tiwoDialogTitleSpec;

   @Override
   public void prepare() throws Exception {
      _tiwoDialogTitleSpec = getSpec().get(DialogTiwoTitleSpec.class);
      if (_tiwoDialogTitleSpec == null) {
         throw new IllegalArgumentException(
               "Maximize TIWO step requires DialogTiwoTitleSpec"
            );
      }
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) {
      _tiwoDialogTitleSpec = filteredWorkflowSpec.get(DialogTiwoTitleSpec.class);
      if (_tiwoDialogTitleSpec == null) {
         throw new IllegalArgumentException(
               "Maximize TIWO step requires DialogTiwoTitleSpec"
            );
      }
   }

   @Override
   public void execute() throws Exception {
      BaseDialogNavigator dialogNavigator = new BaseDialogNavigator();

      dialogNavigator.restore(_tiwoDialogTitleSpec.dialogTitle.get());

      //Sometimes web driver is too fast
      CommonUtils.sleep(500L);
      dialogNavigator.waitForDialogToLoad();
      dialogNavigator.waitApplySavedDataProgressBar();

      //Sometimes web driver is too fast
      CommonUtils.sleep(3000L);
   }
}
