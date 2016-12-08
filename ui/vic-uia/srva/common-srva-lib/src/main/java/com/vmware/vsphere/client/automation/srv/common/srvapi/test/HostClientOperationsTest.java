/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.servicespec.HostServiceSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostClientBasicSrvApi;

/**
 * Test the Host Client API operations exposed in HostClientSrvApi class. The test scenario
 * is: 1. Login to host 2. Remove local datastore 3. Readd local datastore
 */
public class HostClientOperationsTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      // host service spec
      HostServiceSpec hostServiceSpec = new HostServiceSpec();
      hostServiceSpec.endpoint.set(testBed.getCommonHost());
      hostServiceSpec.username.set(testBed.getESXAdminUsername());
      hostServiceSpec.password.set(testBed.getESXAdminPasssword());

      // Common host whose local datastore will be deleted
      HostSpec host =
            HostUtil.buildHostSpec(
                  testBed.getCommonHost(),
                  testBed.getESXAdminUsername(),
                  testBed.getESXAdminPasssword(),
                  443,
                  null);
      host.service.set(hostServiceSpec);

      // Datastore to be added
      DatastoreSpec localDatastoreSpec = new DatastoreSpec();
      localDatastoreSpec.name.set(testBed.getLocalDatastoreName());
      localDatastoreSpec.type.set(DatastoreType.VMFS);
      localDatastoreSpec.parent.set(host);
      localDatastoreSpec.service.set(hostServiceSpec);

      spec.links.add(localDatastoreSpec, host);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Nothing to do
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Validate delete local datastore of an ESX host operation.
      composition.appendStep(new BaseWorkflowStep() {
         HostSpec _hostSpec;

         @Override
         public void prepare() throws Exception {
            _hostSpec = getSpec().links.get(HostSpec.class);

         }

         @Override
         public void execute() throws Exception {

            boolean isLocalDsDeleted =
                  HostClientBasicSrvApi.getInstance().deleteAllVmfsDatastores(_hostSpec);
            verifyFatal(
                  TestScope.BAT,
                  isLocalDsDeleted,
                  "Verify the local datastore of ESX host is successfully deleted!");
         }
      },
            "Validating delete all vmfs datastores operation.");

      // Validate create local datastore of an ESX host operation.
      composition.appendStep(new BaseWorkflowStep() {
         DatastoreSpec _datastoreSpec;

         @Override
         public void prepare() throws Exception {
            _datastoreSpec = getSpec().links.get(DatastoreSpec.class);

         }

         @Override
         public void execute() throws Exception {

            boolean isLocalDsCreated =
                  HostClientBasicSrvApi.getInstance().createVmfsDatastore(_datastoreSpec);
            verifyFatal(
                  TestScope.BAT,
                  isLocalDsCreated,
                  "Verify the local datastore of ESX host is successfully created!");
         }
      },
            "Validating create local datastore operation.");


      // Validate delete local datastore of an ESX host operation.
      composition.appendStep(new BaseWorkflowStep() {
         DatastoreSpec _datastoreSpec;

         @Override
         public void prepare() throws Exception {
            _datastoreSpec = getSpec().links.get(DatastoreSpec.class);

         }

         @Override
         public void execute() throws Exception {

            boolean isLocalDsDeleted =
                  HostClientBasicSrvApi.getInstance().deleteVmfsDatastore(_datastoreSpec);
            verifyFatal(
                  TestScope.BAT,
                  isLocalDsDeleted,
                  "Verify the local datastore of ESX host is successfully deleted!");
         }

         @Override
         public void clean() throws Exception {

            boolean isLocalDsCreated =
                  HostClientBasicSrvApi.getInstance().createVmfsDatastore(_datastoreSpec);
            verifyFatal(
                  TestScope.BAT,
                  isLocalDsCreated,
                  "Verify the local datastore of ESX host is successfully created!");
         }

      },
            "Validating delete local datastore operation.");

      // Validate enable vmotion on an ESX host operation
      composition.appendStep(new BaseWorkflowStep() {
         HostSpec _hostSpec;

         @Override
         public void prepare() throws Exception {
            _hostSpec = getSpec().links.get(HostSpec.class);

         }

         @Override
         public void execute() throws Exception {
            boolean isVmotionEnabled = HostClientBasicSrvApi.getInstance().enableVmotion(_hostSpec);
            verifyFatal(
                  TestScope.BAT,
                  isVmotionEnabled,
                  "Verify vmotion is enabled on host.");
         }
      },
            "Validating the enabling of vmotion on host.");

      // Validate disable vmotion on an ESX host operation
      composition.appendStep(new BaseWorkflowStep() {
         HostSpec _hostSpec;

         @Override
         public void prepare() throws Exception {
            _hostSpec = getSpec().links.get(HostSpec.class);

         }

         @Override
         public void execute() throws Exception {
            boolean isVmotionDisabled = HostClientBasicSrvApi.getInstance().disableVmotion(_hostSpec);
            verifyFatal(
                  TestScope.BAT,
                  isVmotionDisabled,
                  "Verify vmotion is disabled on host.");
         }
      },
            "Validating the disabling of vmotion on host.");
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }

}
