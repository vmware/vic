/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.step;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.URLSpec;
import com.vmware.client.automation.workflow.CommonUIWorkflowStep;

/**
 * Common workflow step for connecting a web browser.
 *
 * Currently the step opens a new browser instance of Internet Explorer and
 * launches the passed URL on it.
 */
public class LoadURLStep extends CommonUIWorkflowStep {

   private URLSpec _URLSpec;

   @Override
   public void prepare() throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Connect Browser Step");
      }

      _URLSpec = getSpec().links.get(URLSpec.class);

      if (_URLSpec == null) {
         throw new IllegalArgumentException(
               "The required BrowserSpec is not set.");
      }

      if (Strings.isNullOrEmpty(_URLSpec.url.get())) {
         throw new IllegalArgumentException("The url is not set.");
      }
   }

   @Override
   public void execute() throws Exception {
      UI.browser.open(_URLSpec.url.get());
   }
}
