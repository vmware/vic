/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.test.utils.vcuicat;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.List;

import org.testng.annotations.Test;

import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.provider.commontb.HostProvider;
import com.vmware.vsphere.client.automation.provider.commontb.VcProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * We need this class as Stateless testbed deploy leaves ESX hosts inside VC.
 * The purpose of this code is to workaround this problem. Connect to the VC and
 * to remove the hosts so CommonTestBedProvider can continue to work as
 * expected.
 */
public class VcHostDisconnection extends BaseTestWorkflow {

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      if (isProviderRun) {
         super.invokeTestExecuteCommand();
      } else {
         super.execute();
      }
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBedProvider = testbedBridge.requestTestbed(
            VcProvider.class, true);
      // TODO: kiliev Fix the test to connect to VC and retrive all hosts and
      // then disconnect them
      VcSpec vcSpec = testBedProvider
            .getPublishedEntitySpec(VcProvider.DEFAULT_ENTITY);

      // Spec for the datacenter
      DatacenterSpec datacenterSpec = new DatacenterSpec();
      // This is the default name of the datacenter
      // when stateless esx tesdtbed is deployed
      datacenterSpec.name.set("Untitled");
      datacenterSpec.parent.set(vcSpec);

      // Spec for the Host1
      HostSpec hostSpec1 = testbedBridge.requestTestbed(HostProvider.class,
            true).getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      hostSpec1.parent.set(datacenterSpec);

      // Spec for the Host2
      HostSpec hostSpec2 = testbedBridge.requestTestbed(HostProvider.class,
            true).getPublishedEntitySpec(HostProvider.DEFAULT_ENTITY);
      hostSpec2.parent.set(datacenterSpec);

      testSpec.add(vcSpec, datacenterSpec, hostSpec1, hostSpec2);
   }

   @Override
   public void composeTestSteps(
         WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      flow.appendStep("Remove all hosts from the VC", new BaseWorkflowStep() {
         private List<HostSpec> _hostsToRemove;

         // TestWorkflowStep methods
         @Override
         public void prepare(WorkflowSpec filteredWorkflowSpec)
               throws Exception {
            _hostsToRemove = filteredWorkflowSpec.getAll(HostSpec.class);
         }

         public void execute() throws Exception {
            for (HostSpec host : _hostsToRemove) {
               try {
                  // The hosts names are IP addresses but are connected by
                  // hostname in the VC. We need to change the host name from ip
                  // to hostname.
                  InetAddress addr;
                  addr = InetAddress.getByName(host.name.get());
                  String hostName = addr.getHostName();
                  host.name.set(hostName);
               } catch (UnknownHostException e) {
               }
               _logger.info(String.format("Remove host '%s'", host.name.get()));
               HostBasicSrvApi.getInstance().deleteHostSafely(host);
            }
         }

      });
   }

   @Override
   public void initSpec() {
      // This method is here just because BaseTestWorkflow requires its
      // implementation.
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // This method is here just because BaseTestWorkflow requires its
      // implementation.
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {
      // This method is here just because BaseTestWorkflow requires its
      // implementation.
   }
}
