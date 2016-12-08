/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.step;

import java.util.List;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.servicespec.SeleniumServiceSpec;
import com.vmware.client.automation.util.SeleniumUtil;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.provider.connector.SeleniumConnector;

/**
 * Common workflow step for connecting a web browser.
 *
 * Currently the step opens a new browser instance of Internet Explorer and
 * launches the passed URL on it.
 */
public class ConnectBrowserStep extends BaseWorkflowStep {

   private SeleniumServiceSpec _seleniumServiceSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare() throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Connect Browser Step");
      }

      _seleniumServiceSpec = getSpec().links.get(SeleniumServiceSpec.class);

      if (_seleniumServiceSpec == null) {
         throw new IllegalArgumentException("The required BrowserSpec is not set.");
      }

   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      SeleniumConnector connector =
            SeleniumUtil.getSeleniumConnector(_seleniumServiceSpec);
      connector.startBrowser();
   }

   @Override
   public void clean() {
      SeleniumUtil.getSeleniumConnector(_seleniumServiceSpec).disconnect();
   }


   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      List<EntitySpec> specList = filteredWorkflowSpec.links.getAll(EntitySpec.class);
      for (EntitySpec entitySpec : specList) {
         if (entitySpec.service.get() instanceof SeleniumServiceSpec) {
            _seleniumServiceSpec = (SeleniumServiceSpec) entitySpec.service.get();
         }
      }

      if (_seleniumServiceSpec == null) {
         throw new IllegalArgumentException(
               "The required SeleniumServiceSpec is not set.");
      }

   }
}
