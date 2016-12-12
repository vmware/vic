/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.step.AssignStoragePolicyToVmByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.AssignTagToDatastoreByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateBackingCategoryAndTagsByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateDatastoreStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateStoragePolicyByApiStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;

/**
 * Test the Vm Policy vAPI operations exposed in VmPolicySrvApi class.
 * The test scenario is:
 * 1. Create the policy
 * 2. Verify policy is created
 * 3. Create a tag
 * 4. Verify the tag is created
 * 5. Attach resources to the tag
 * 6. Verify the attached resources
 * 7. Delete the tags and verify they are deleted
 * 8. Delete policy and verify it is deleted
 * 9. Create a backing tag category
 */
public class VmPolicyServiceOperationsTest extends BaseTestWorkflow {

   private static final String TAG_VM1 = "vm1";
   private static final String TAG_VM2 = "vm2";
   private static final String TAG_EVERYTHING_ELSE = "everythingElse";

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      // Datacenter that has a host
      DatacenterSpec datacenter = new DatacenterSpec();
      datacenter.name.set(testBed.getCommonDatacenterName());
      datacenter.tag.set(TAG_EVERYTHING_ELSE);

      ClusterSpec cluster = SpecFactory.getSpec(ClusterSpec.class, testBed.getCommonClusterName(), datacenter);
      cluster.tag.set(TAG_EVERYTHING_ELSE);

      HostSpec host = HostUtil.buildHostSpec(
            testBed.getCommonHost(),
            testBed.getESXAdminUsername(),
            testBed.getESXAdminPasssword(),
            443,
            cluster
            );
      host.tag.set(TAG_EVERYTHING_ELSE);

      // VMs
      VmSpec vm1 = SpecFactory.getSpec(VmSpec.class, host);
      vm1.tag.set(TAG_VM1);

      VmSpec vm2 = SpecFactory.getSpec(VmSpec.class, host);
      vm2.tag.set(TAG_VM2);

      DatastoreSpec datastore =
            SpecFactory.getSpec(DatastoreSpec.class, "GoldDatastore", host);
      datastore.type.set(DatastoreType.NFS);
      datastore.remoteHost.set(CommonUtil.getLocalizedString("web.server.ip"));
      datastore.remotePath.set(CommonUtil.getLocalizedString("datastore.profile1"));
      datastore.tag.set(TAG_EVERYTHING_ELSE);

      BackingCategorySpec categorySpec = SpecFactory.getSpec(BackingCategorySpec.class);
      categorySpec.description.set(SpecFactory.buildUniqueDesc());
      categorySpec.tag.set(TAG_EVERYTHING_ELSE);

      BackingTagSpec datastoreTag = SpecFactory.getSpec(BackingTagSpec.class);
      datastoreTag.description.set(SpecFactory.buildUniqueDesc());
      datastoreTag.tag.set(TAG_EVERYTHING_ELSE);

      StoragePolicySpec storagePolicySpec = SpecFactory.getSpec(StoragePolicySpec.class);
      storagePolicySpec.description.set(SpecFactory.buildUniqueDesc());
      storagePolicySpec.tag.set(TAG_EVERYTHING_ELSE);

      spec.links.add(datacenter, host, vm1, vm2,
            datastore, categorySpec, datastoreTag, storagePolicySpec);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Create the VMs
      composition.appendStep(new CreateVmByApiStep(), "Create the VMs");

      // Create a backing category and tag
      composition.appendStep(new CreateBackingCategoryAndTagsByApiStep());

      // Create the datastore
      composition.appendStep(new CreateDatastoreStep());

      // Assign tags to the datastore
      composition.appendStep(new AssignTagToDatastoreByApiStep());

      // Create Storage policy by API
      composition.appendStep(new CreateStoragePolicyByApiStep());
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {
      // Step 5. Add a storage policy to VM1
      composition.appendStep(new AssignStoragePolicyToVmByApiStep(), "Assign storage policy to VM", TestScope.BAT,
            new String[] {TAG_VM1, TAG_EVERYTHING_ELSE});
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }

}
