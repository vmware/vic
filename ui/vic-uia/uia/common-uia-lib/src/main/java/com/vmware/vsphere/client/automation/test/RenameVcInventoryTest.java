/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vim.binding.vim.Datacenter;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;
import com.vmware.vsphere.client.automation.common.workflow.NGCTestWorkflow;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatacenterBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;

/**
 * Adds a suffix to the entities from the VC inventory
 */
public class RenameVcInventoryTest extends NGCTestWorkflow {

   @Override
   @Test(description = "Adds a suffix to the entities from VC inventory")
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer commonTestbed = testbedBridge.requestTestbed(CommonTestBedProvider.class, true);

      // VC
      VcSpec vcSpec = commonTestbed.getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      // Datacenter
      DatacenterSpec datacenterSpec = commonTestbed.getPublishedEntitySpec(CommonTestBedProvider.DC_ENTITY);

      // Cluster
      ClusterSpec clusterSpec = commonTestbed.getPublishedEntitySpec(CommonTestBedProvider.CLUSTER_ENTITY);

      // // Common datastore
      // DatastoreSpec saDatastoreSpec = commonTestbed.getPublishedEntitySpec(CommonTestBedProvider.LOCAL_DS_ENTITY);
      // saDatastoreSpec.tag.set(TAG_STANDALONE_DATASTORE);

      testSpec.add(vcSpec, datacenterSpec, clusterSpec);

      super.initSpec(testSpec, testbedBridge);
   }

   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      flow.appendStep("Rename the entities",
            new EnhancedBaseWorkflowStep() {
               @UsesSpec
               protected VcSpec _vc;

               @UsesSpec
               protected DatacenterSpec _datacenter;

               @UsesSpec
               protected ClusterSpec _cluster;

               @Override
               public void execute() throws Exception {
                  String vcIp = _vc.service.get().endpoint.get();

                  verifyFatal(renameDatastores('-' + vcIp),
                        "Verify datastores are renamed");

                  ClusterSpec clusterNewName = new ClusterSpec();
                  clusterNewName.name.set(_cluster.name.get() + '-' + vcIp);
                  verifyFatal(ClusterBasicSrvApi.getInstance().
                        renameCluster(_cluster, clusterNewName),
                        "Verify cluster is renamed");

                  verifyFatal(DatacenterBasicSrvApi.getInstance().
                        renameDatacenter(_datacenter, _datacenter.name.get() + '-' + vcIp),
                        "Verify datacenter is renamed");
               }

               private boolean renameDatastores(String suffix) throws Exception {
                  Datacenter datacenter = ManagedEntityUtil.getManagedObject(_datacenter, _datacenter.service.get());

                  boolean result = true;
                  ManagedObjectReference[] datastoresMoRefs = datacenter.getDatastore();
                  for (ManagedObjectReference datastoreMoRef : datastoresMoRefs) {
                     Datastore datastore =
                           ManagedEntityUtil.getManagedObjectFromMoRef(datastoreMoRef, _datacenter.service.get());
                     ManagedObjectReference taskMoRef = datastore.rename(datastore.getName() + suffix);

                     // success or failure of the task
                     result = result && VcServiceUtil.waitForTaskSuccess(taskMoRef, _datacenter.service.get());
                  }
                  return result;
               }
            });
   }
}
