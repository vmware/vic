/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;

/**
 * Mounts a datastore to a host via the API.
 */
public class MountDatastoreByApiStep extends BaseWorkflowStep {

   private DatastoreSpec _datastoreToAttachSpec;

   /**
    * {@inheritDoc}
    */
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _datastoreToAttachSpec = filteredWorkflowSpec.get(DatastoreSpec.class);
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      DatastoreBasicSrvApi.getInstance().createDatastore(_datastoreToAttachSpec);
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void clean() throws Exception {
      DatastoreBasicSrvApi.getInstance().unmountDatastore(_datastoreToAttachSpec);
   }
}
