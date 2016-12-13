/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;

/**
 * Common workflow step for creation of one or more vSphere clusters via
 * the VC sdk.
 * <br>To use this step in automation tests:
 * <li>In the <code>prepare()</code> method of the
 * <code>BaseTestWorkflow</code> test, create one or more
 * <code>ClusterSpec</code> instances and link them to the test spec.
 * <li>Append a  <code>CreateClusterStep</code> instance to the test/prerequisite
 *  workflow composition.
 */
public class CreateClusterStep extends BaseWorkflowStep {

   private List<ClusterSpec> _clustersToCreate;
   private List<ClusterSpec> _clustersToClean;

   @Override
   /**
    * @inheritDoc
    */
   public void prepare() throws Exception {
      if (Strings.isNullOrEmpty(getTitle())) {
         setTitle("Create clusters");
      }

      _clustersToCreate = getSpec().links.getAll(ClusterSpec.class);

      if (_clustersToCreate == null || _clustersToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'ClusterSpec' instances");
      }

      _clustersToClean = new ArrayList<ClusterSpec>();
   }

   @Override
   /**
    * @inheritDoc
    */
   public void execute() throws Exception {
      for (ClusterSpec clusterSpec : _clustersToCreate) {
         if (!ClusterBasicSrvApi.getInstance().createCluster(clusterSpec)) {
            throw new Exception(String.format(
                  "Unable to create cluster '%s'",
                  clusterSpec.name.get()));
         } else {
            _clustersToClean.add(clusterSpec);
         }
      }
   }

   @Override
   /**
    * @inheritDoc
    */
   public void clean() throws Exception {
      for (ClusterSpec clusterSpec : _clustersToClean) {
         ClusterBasicSrvApi.getInstance().deleteClusterSafely(clusterSpec);
      }
   }

   // TestWorkflowStep methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _clustersToCreate = filteredWorkflowSpec.links.getAll(ClusterSpec.class);
      if (_clustersToCreate == null || _clustersToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatacenterSpec' instances");
      }

      _clustersToClean = new ArrayList<ClusterSpec>();

   }
}
