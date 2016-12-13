/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.createinventory;

import java.util.List;

import org.testng.annotations.Test;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatacenterBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * The test is used to verify and fix the common testbed setup.
 * It will be executed by the test setup and test execution systems to verify
 * the setup before running list of tests.
 * NOTE: The test is special it has no clean methods it will try to fix the
 * setup but if it fails the setup will not be cleaned up.
 *
 * Steps:
 * 1. Validates common DC existence and creates it if it doesn't exist.
 * 2. Validates common cluster existence in the common DC and creates it if it
 * doesn't exist.
 * 3. Validates common host is connected in the common cluster and
 * that is it is not in maintenance mode.
 * 4. Validates that the free host is not already attached to the common DC.
 * If so disconnects it.
 * 4. Due to PR 1087863 Validate/create vDC which should be listed last in the
 * Object Navigator.
 * TODO: Once the PR is resolved remove the step from the code.
 */
public class CommonTestbedSetupValidatorTest extends BaseTestWorkflow {

   private static final  String FREE_HOST_TAG = "FREE_HOST";
   private static final String COMMON_HOST_TAG = "COMMON_HOST";

   @Override
   public void initSpec() {
      BaseSpec commonTestbedSpec = new BaseSpec();
      setSpec(commonTestbedSpec);

      // Define common datacenter spec
      DatacenterSpec commonDatacenterSpec = new DatacenterSpec();
      commonDatacenterSpec.name.set(testBed.getCommonDatacenterName());
      // link the datacenter
      commonTestbedSpec.links.add(commonDatacenterSpec);

      // Define common cluster spec
      ClusterSpec commonClusterSpec = new ClusterSpec();
      commonClusterSpec.name.set(testBed.getCommonClusterName());
      commonClusterSpec.drsEnabled.set(true);
      commonClusterSpec.parent.set(commonDatacenterSpec);
      // link the cluster
      commonTestbedSpec.links.add(commonClusterSpec);

      // Define common host spec
      HostSpec commonHostSpec = new HostSpec();
      commonHostSpec.name.set(testBed.getCommonHost());
      commonHostSpec.parent.set(commonClusterSpec);
      commonHostSpec.port.set(443);
      commonHostSpec.userName.set(testBed.getESXAdminUsername());
      commonHostSpec.password.set(testBed.getESXAdminPasssword());
      commonHostSpec.tag.set(COMMON_HOST_TAG);
      // link the host
      commonTestbedSpec.links.add(commonHostSpec);

      // Create a spec for host added in the common DC.
      // It is done when using the Nimbus vc-ui testbed.
      // For such cases the host should be removed before
      // running the tests
      HostSpec freeHostSpec = new HostSpec();
      String hostIp = "";
      try {
         List<String> hostList = testBed.getHosts(1);
         if(hostList.size() == 1) {
            hostIp = hostList.get(0);
         }
      } catch (RuntimeException e) {
         if (e.getMessage().contains("No ESX hosts were found")) {
            hostIp = "";
         } else {
            throw e;
         }
      }
      freeHostSpec.name.set(hostIp);
      freeHostSpec.parent.set(commonDatacenterSpec);
      freeHostSpec.port.set(443);
      freeHostSpec.userName.set(testBed.getESXAdminUsername());
      freeHostSpec.password.set(testBed.getESXAdminPasssword());
      freeHostSpec.tag.set(FREE_HOST_TAG);
      // link the host
      commonTestbedSpec.links.add(freeHostSpec);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // nothing to do here
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Validate/create common datacenter
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            DatacenterSpec commonDatacenterSpec = getSpec().links.get(DatacenterSpec.class);
            boolean dcFound = DatacenterBasicSrvApi.getInstance().checkDatacenterExists(commonDatacenterSpec);
            if(!dcFound) {
               verifyFatal(TestScope.BAT,
                     DatacenterBasicSrvApi.getInstance().createDatacenter(commonDatacenterSpec),
                     "Common datacenter is creation operation!");

               dcFound = DatacenterBasicSrvApi.getInstance().checkDatacenterExists(commonDatacenterSpec);
            }
            verifyFatal(TestScope.BAT, dcFound,
                  "Common datacenter is found in the system!");
         }

      }, "Validate/create common datacenter.");

      // Validate/create common cluster
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            ClusterSpec commonClusterSpec = getSpec().links.get(ClusterSpec.class);
            boolean clusterFound = ClusterBasicSrvApi.getInstance().checkClusterExists(commonClusterSpec);
            if(!clusterFound) {
               verifyFatal(TestScope.BAT,
                     ClusterBasicSrvApi.getInstance().createCluster(commonClusterSpec),
                     "Common cluster is creation operation!");

               clusterFound = ClusterBasicSrvApi.getInstance().checkClusterExists(commonClusterSpec);
            }
            verifyFatal(TestScope.BAT, clusterFound,
                  "Common cluster is found in the system!");

            // Make sure DRS is enabled on the cluster
            if (!ClusterBasicSrvApi.getInstance().isDrsEnabled(commonClusterSpec)) {
               ClusterBasicSrvApi.getInstance().reconfigureCluster(commonClusterSpec, commonClusterSpec);
            }
            verifyFatal(TestScope.BAT,
                  ClusterBasicSrvApi.getInstance().isDrsEnabled(commonClusterSpec),
                  "DRS is enabled on the common cluster.");
         }

      }, "Validate/create common cluster.");

      // Validate/create common host is connected and operational
      composition.appendStep(new BaseWorkflowStep() {

         @Override
         public void execute() throws Exception {
            HostSpec commonHostSpec = getSpec().links.get(HostSpec.class);
            if(!HostBasicSrvApi.getInstance().checkHostExists(commonHostSpec)) {
               verifyFatal(TestScope.BAT,
                     HostBasicSrvApi.getInstance().addHost(commonHostSpec, true),
                     "Verify common host is added to the common cluster!");
            }

            // Try to connect host. Task complete fast if already connected.
            HostBasicSrvApi.getInstance().connectHost(commonHostSpec);

            verifyFatal(TestScope.BAT,
                  HostBasicSrvApi.getInstance().isConnected(commonHostSpec),
                  "Verify common host is connected!");

            // Try to exit from MM. Task complete fast if not in MM.
            HostBasicSrvApi.getInstance().exitMaintenanceMode(commonHostSpec);

            verifyFatal(TestScope.BAT,
                  !HostBasicSrvApi.getInstance().isInMaintenanceMode(commonHostSpec),
                  "Verify cluster is not in maintanance mode!");

         }

            }, "Validate/create common host is connected and operational.",
            TestScope.BAT, new String[] { COMMON_HOST_TAG });

      // Remove the free host if already attached by other test or nimbus
      // testbed
      composition.appendStep(
            new BaseWorkflowStep() {
               HostSpec _freeHostSpec;
               @Override
               public void prepare() throws Exception {
                  _freeHostSpec = getSpec().links.get(HostSpec.class);
               }

               @Override
               public void execute() throws Exception {

                  if (_freeHostSpec.name.get().isEmpty()) {
                     // No host provided no need to check for it!
                     return;
                  }

                  HostBasicSrvApi.getInstance().deleteHostSafely(_freeHostSpec);

                  verifyFatal(TestScope.BAT,
                        !HostBasicSrvApi.getInstance().checkHostExists(_freeHostSpec),
                        "Verify host is diconnected!");

               }

            }, "The free host is not already attached",
            TestScope.BAT,
            new String[] { FREE_HOST_TAG });
   }

   @Override
   @Test(groups={"infra"})
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
