/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.cluster.lib.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.cluster.lib.view.ClusterSummaryTabPage;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;

/**
 * Verifies the DRS portlet and badge presence
 */
public class VerifyVsphereDrsOnSummaryPageStep extends CommonUIWorkflowStep {

   protected ClusterSpec _clusterSpec;
   protected Boolean _expectedDrsState;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
       _clusterSpec = filteredWorkflowSpec.get(ClusterSpec.class);
       ensureNotNull(_clusterSpec, "ClusterSpec spec is missing.");

       _expectedDrsState = _clusterSpec.drsEnabled.get();
       ensureNotNull(_expectedDrsState, "drsEnabled property for ClusterSpec is not found.");
   }

   @Override
   public void execute() throws Exception {
      ClusterSummaryTabPage summaryPage = new ClusterSummaryTabPage();
      boolean areDrsPortletAndBadgeVisible = (summaryPage
            .isRunDrsPortletVisible() && summaryPage.isDrsBadgeVisible());

      verifySafely(new EqualsAssertion(
            areDrsPortletAndBadgeVisible,
            _expectedDrsState,
            String.format(
                  "DRS Portlet and Badge presence is according the expectations for cluster %s",
                  _clusterSpec.name.get())));
   }
}