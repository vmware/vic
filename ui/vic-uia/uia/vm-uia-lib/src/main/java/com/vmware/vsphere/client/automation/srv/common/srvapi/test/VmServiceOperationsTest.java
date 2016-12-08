/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi.test;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.SpecFactory;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;
import com.vmware.vsphere.client.automation.srv.common.step.CreateVmByApiStep;

/**
 * Test the API operations exposed in VmSrvApi class.
 * The test scenario is:
 * 1. Create the VM and check it is created.
 * 2. Delete the VM and check that it is not present.
 * 3. Power on VM and check if it is powered on
 * 4. Power off VM and check if it is powered off
 */
public class VmServiceOperationsTest extends BaseTestWorkflow {

   public static String TAG_VM_PREREQ = "vmPrereq";

   @Override
   public void initSpec() {
      BaseSpec spec = new BaseSpec();
      setSpec(spec);

      // Datacenter that has a host
      DatacenterSpec datacenter = new DatacenterSpec();
      datacenter.name.set(testBed.getCommonDatacenterName());

      HostSpec host = HostUtil.buildHostSpec(
            testBed.getHosts(1).get(0),
            testBed.getESXAdminUsername(),
            testBed.getESXAdminPasssword(),
            443,
            null
            );
      host.parent.set(datacenter);

      VmSpec vm1 = SpecFactory.getSpec(VmSpec.class, host);
      vm1.tag.set(TAG_VM_PREREQ);

      VmSpec vm2 = SpecFactory.getSpec(VmSpec.class, host);

      VmSpec powerOnVmSpec = SpecFactory.getSpec(VmSpec.class, host);
      powerOnVmSpec.guestId.set("winVistaGuest");
      powerOnVmSpec.tag.set(TAG_VM_PREREQ);

      VmSpec powerOffVmSpec = SpecFactory.getSpec(VmSpec.class, host);
      powerOffVmSpec.guestId.set("winVistaGuest");
      powerOffVmSpec.tag.set(TAG_VM_PREREQ);

      spec.links.add(
            datacenter,
            host,
            vm1,
            vm2,
            powerOnVmSpec,
            powerOffVmSpec
            );
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // Prerequisite - create VMs
      composition.appendStep(new CreateVmByApiStep(), "Prerequisite - create VMs", TestScope.BAT, new String[] {TAG_VM_PREREQ});
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Validate create VM operation
      composition.appendStep(
            new BaseWorkflowStep() {

               private VmSpec _vmSpec;

               @Override
               public void prepare() throws Exception {
                  _vmSpec = getSpec().links.getAll(VmSpec.class).get(1);
               }

               @Override
               public void execute() throws Exception {
                  verifyFatal(
                        TestScope.BAT,
                        VmSrvApi.getInstance().createVm(_vmSpec),
                        "Verifying the vmSrvApi.createVm!"
                        );

                  verifyFatal(
                        TestScope.BAT,
                        VmSrvApi.getInstance().checkVmExists(_vmSpec),
                        "Verifying the VM is created!"
                        );
               }
            },
            "Validating create VM operation"
            );

      // Validate delete VM operation
      composition.appendStep(
            new BaseWorkflowStep() {

               private VmSpec _vmSpec;

               @Override
               public void prepare() throws Exception {
                  _vmSpec = getSpec().links.getAll(VmSpec.class).get(0);
               }

               @Override
               public void execute() throws Exception {
                  verifyFatal(
                        TestScope.BAT,
                        VmSrvApi.getInstance().deleteVm(_vmSpec),
                        "Verifying the VmSrvApi.deleteVm!"
                        );

                  verifyFatal(
                        TestScope.BAT,
                        !VmSrvApi.getInstance().checkVmExists(_vmSpec),
                        "Verifying the VM was deleted!");
               }
            },
            "Validating delete VM operation"
            );

      // Validate power on VM operation
      composition.appendStep(
            new BaseWorkflowStep() {

               private VmSpec _vmSpec;

               @Override
               public void prepare() throws Exception {
                  _vmSpec = getSpec().links.getAll(VmSpec.class).get(2);
               }

               @Override
               public void execute() throws Exception {
                  verifyFatal(
                        TestScope.BAT,
                        VmSrvApi.getInstance().powerOnVm(_vmSpec),
                        "Verifying the VmSrvApi.powerOnVm!"
                        );

                  verifyFatal(
                        TestScope.BAT,
                        VmSrvApi.getInstance().isVmPoweredOn(_vmSpec),
                        "Verifying the VM is powered on!");
               }

               @Override
               public void clean() throws Exception {
                  //TODO: rreymer implement when power off operation is available
               }
            },
            "Validating power on VM operation"
            );
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
