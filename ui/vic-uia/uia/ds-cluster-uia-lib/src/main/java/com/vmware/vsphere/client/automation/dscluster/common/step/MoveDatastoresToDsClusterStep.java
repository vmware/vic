/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreClusterSrvApi;

/**
 * Common test step that moves datastores to datastore cluster.
 */
public class MoveDatastoresToDsClusterStep extends BaseWorkflowStep {

   private DatastoreClusterSpec _dsCluster;
   private List<DatastoreSpec> _datastores;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      _dsCluster = filteredWorkflowSpec.links.get(DatastoreClusterSpec.class);
      ensureNotNull(_dsCluster, "The spec has no links to 'DatastoreClusterSpec' instances");

      _datastores = filteredWorkflowSpec.links.getAll(DatastoreSpec.class);
      if (CollectionUtils.isEmpty(_datastores)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatastoreSpec' instances");
      }
   }

   @Override
   public void prepare() throws Exception {

      _dsCluster = getSpec().links.get(DatastoreClusterSpec.class);
      if (_dsCluster == null) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatastoreClusterSpec' instances");
      }

      _datastores = getSpec().links.getAll(DatastoreSpec.class);
      if (CollectionUtils.isEmpty(_datastores)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatastoreSpec' instances");
      }
   }

   @Override
   public void execute() throws Exception {
      DatastoreClusterSrvApi.getInstance().moveDatastoresToDsCluster(_datastores, _dsCluster);
   }

   @Override
   public void clean() throws Exception {
      // No cleanup needed, the datastore will be automatically released, when
      // the datastore cluster gets deleted at cleanup of
      // CreateDatastoreClusterStep()
   }
}
