/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.srv.createinventory;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreClusterSrvApi;

/**
 * Test the Datastore Cluster API operations exposed in DsClusterSrvApi class.
 * The test scenario is:<br>
 * 1. Login to VC and check that the datastore cluster specified in the spec does
 * not exist.<br>
 * 2. Create a datastore cluster and check it is created successfully<br>
 * 3. Delete the datastore cluster.<br>
 */
public class DsClusterOperationsTest extends BaseTestWorkflow {

    @Override
    public void initSpec() {
        BaseSpec spec = new BaseSpec();
        setSpec(spec);

        DatacenterSpec dcSpec = SpecFactory.getSpec(DatacenterSpec.class, testBed.getCommonDatacenterName(), null);

        DatastoreClusterSpec dsClusterSpec = SpecFactory.getSpec(DatastoreClusterSpec.class, dcSpec);

        spec.links.add(dcSpec, dsClusterSpec);
    }

    @Override
    public void composePrereqSteps(WorkflowComposition composition) {
        // Nothing to do
    }

    @Override
    public void composeTestSteps(WorkflowComposition composition) {

        // Validate datastore cluster check existence operation
        composition.appendStep(new BaseWorkflowStep() {

            @Override
            public void execute() throws Exception {
                DatastoreClusterSpec dsClusterSpec = getSpec().links.get(DatastoreClusterSpec.class);

                boolean dsClusterExists = DatastoreClusterSrvApi.getInstance().checkDsClusterExists(dsClusterSpec);
                verifyFatal(TestScope.BAT, !dsClusterExists, "Verify the datastore cluster doesn't exist.");
            }
        }, "Validating datastore cluster check existence operation.");

        // Validate create datastore cluster operation
        composition.appendStep(new BaseWorkflowStep() {

            @Override
            public void execute() throws Exception {
                DatastoreClusterSpec dsClusterSpec = getSpec().links.get(DatastoreClusterSpec.class);

                boolean createDsClusterOperation = DatastoreClusterSrvApi.getInstance().createDatastoreCluster(dsClusterSpec);
                verifyFatal(
                    TestScope.BAT,
                    createDsClusterOperation,
                    "Verify the DatastoreClusterSrvApi.createDatastoreCluster operation.");

                boolean dsClusterExists = DatastoreClusterSrvApi.getInstance().checkDsClusterExists(dsClusterSpec);
                verifyFatal(TestScope.BAT, dsClusterExists, "Verify the datastore cluster is created.");
            }
        },
            "Validating create datastore cluster operation.");

        // Validate delete datastore cluster operation
        composition.appendStep(new BaseWorkflowStep() {

            @Override
            public void execute() throws Exception {
                DatastoreClusterSpec dsClusterSpec = getSpec().links.get(DatastoreClusterSpec.class);

                boolean dsClusterDeleted = DatastoreClusterSrvApi.getInstance().deleteDsCluster(dsClusterSpec);
                verifyFatal(
                    TestScope.BAT,
                    dsClusterDeleted,
                    "Verify the DatastoreClusterSrvApi.removeDsCluster operation.");

                boolean dsClusterExists = DatastoreClusterSrvApi.getInstance().checkDsClusterExists(dsClusterSpec);
                verifyFatal(TestScope.BAT, !dsClusterExists, "Verify the datastore cluster is deleted.");
            }
        },
            "Validating datastore cluster delete operation.");
    }

    @Override
    @Test
    @TestID(id = "0")
    public void execute() throws Exception {
        super.execute();
    }

}
