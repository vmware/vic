/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;

/**
 * Common test step that creates datastores.
 */
public class CreateDatastoreStep extends BaseWorkflowStep {

   private List<DatastoreSpec> _datastoresToCreate;
   private List<DatastoreSpec> _datastoresToDelete;

   /**
    * @inheritDoc
    */
   @Override
   public void prepare() throws Exception {

      _datastoresToCreate = getSpec().links.getAll(DatastoreSpec.class);

      if (CollectionUtils.isEmpty(_datastoresToCreate)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatastoreSpec' instances");
      }

      _datastoresToDelete = new ArrayList<DatastoreSpec>();
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      for (DatastoreSpec datastoreSpec : _datastoresToCreate) {
         if (!DatastoreBasicSrvApi.getInstance().createDatastore(datastoreSpec)) {
            throw new Exception(String.format(
                  "Unable to create datastore '%s'",
                  datastoreSpec.name.get()));
         }
         _datastoresToDelete.add(datastoreSpec);
      }
   }

   /**
    * @inheritDoc
    */
   @Override
   public void clean() throws Exception {
      for (DatastoreSpec datastoreSpec : _datastoresToDelete) {
         DatastoreBasicSrvApi.getInstance().deleteDatastoreSafely(datastoreSpec);
      }
   }

   // TestWorkflowStep methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _datastoresToCreate = filteredWorkflowSpec.getAll(DatastoreSpec.class);

      if (CollectionUtils.isEmpty(_datastoresToCreate)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatastoreSpec' instances");
      }

      _datastoresToDelete = new ArrayList<DatastoreSpec>();

   }
}
