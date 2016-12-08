/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;

/**
 * Common test step that deletes datastores.
 */
public class DeleteDatastoreStep extends BaseWorkflowStep {

   private List<DatastoreSpec> _datastoresToCreate;
   private List<DatastoreSpec> _datastoresToDelete;

   /**
    * @inheritDoc
    */
   @Override
   public void prepare() throws Exception {

      _datastoresToDelete = getSpec().links.getAll(DatastoreSpec.class);

      if (CollectionUtils.isEmpty(_datastoresToDelete)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatastoreSpec' instances");
      }

      _datastoresToCreate = new ArrayList<DatastoreSpec>();
   }

   /**
    * {@inheritDoc}
    */
   @Override
   public void execute() throws Exception {
      for (DatastoreSpec datastoreSpec : _datastoresToDelete) {
         if (!DatastoreBasicSrvApi.getInstance().deleteDatastoreSafely(datastoreSpec)) {
            throw new Exception(
                  String.format("Unable to delete datastore '%s'",
                        datastoreSpec.name.get()));
         }
         _datastoresToCreate.add(datastoreSpec);
      }
   }

   /**
    * @inheritDoc
    */
   @Override
   public void clean() throws Exception {
      for (DatastoreSpec datastoreSpec : _datastoresToCreate) {
         DatastoreBasicSrvApi.getInstance().createDatastore(datastoreSpec);
      }
   }
}
