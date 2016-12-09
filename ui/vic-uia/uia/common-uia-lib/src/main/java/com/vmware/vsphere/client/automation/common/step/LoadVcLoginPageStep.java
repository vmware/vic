/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 * The step loads in the web browser the vc url to execute test.
 */
public class LoadVcLoginPageStep extends CommonUIWorkflowStep {

   private VcSpec _vcSpec;

   @Override
   public void prepare() throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Load VC login page Step");
      }

      _vcSpec = getSpec().links.get(VcSpec.class);

      if (_vcSpec == null) {
         throw new IllegalArgumentException("The required VcSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_vcSpec.vscUrl.get())) {
         throw new IllegalArgumentException("The url is not set.");
      }
   }

   @Override
   public void execute() throws Exception {
      UI.browser.open(_vcSpec.vscUrl.get());
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Load VC login page Step");
      }

      _vcSpec = filteredWorkflowSpec.links.get(VcSpec.class);

      if (_vcSpec == null) {
         throw new IllegalArgumentException("The required VcSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_vcSpec.vscUrl.get())) {
         throw new IllegalArgumentException("The url is not set.");
      }
   }
}
