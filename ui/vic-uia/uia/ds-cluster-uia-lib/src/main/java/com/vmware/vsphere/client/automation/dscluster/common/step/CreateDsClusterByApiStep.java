/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreClusterSrvApi;

/**
 * Common workflow step for creation of one or more
 * vSphere datastore clusters via the VC sdk.
 * <br>To use this step in automation tests:
 * <li>In the <code>prepare()</code> method of the
 * <code>BaseTestWorkflow</code> test, create one or more
 * <code>DatastoreClusterSpec</code> instances and link them to the test spec.
 * <li>Append a <code>CreateDsClusterStep</code> instance to the test/prerequisite
 *  workflow composition.
 */
public class CreateDsClusterByApiStep extends BaseWorkflowStep {

   private List<DatastoreClusterSpec> _dsClustersToCreate;
   private List<DatastoreClusterSpec> _dsClustersToClean;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      _dsClustersToCreate = filteredWorkflowSpec.links.getAll(DatastoreClusterSpec.class);

      ensureNotNull(_dsClustersToCreate, "The spec has no links to 'DatastoreClusterSpec' instances");

      _dsClustersToClean = new ArrayList<DatastoreClusterSpec>();
   }

   @Override
   public void execute() throws Exception {
      for (DatastoreClusterSpec dsClusterSpec : _dsClustersToCreate) {
         DatastoreClusterSrvApi.getInstance().createDatastoreCluster(
               dsClusterSpec);
         _dsClustersToClean.add(dsClusterSpec);
      }
   }

   @Override
   public void clean() throws Exception {
      for (DatastoreClusterSpec dsClusterSpec : _dsClustersToClean) {
         DatastoreClusterSrvApi.getInstance().deleteDsClusterSafely(dsClusterSpec);
      }
   }
}
