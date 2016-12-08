/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.createinventory;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.lang.RandomStringUtils;
import org.testng.annotations.Test;

import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.step.AddHostStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateClusterStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateDatacenterFolderStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateDatacenterStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateDatastoreFolderStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateHostFolderStep;
import com.vmware.vsphere.client.automation.srv.common.step.CreateNetworkFolderStep;
import com.vmware.vsphere.client.automation.srv.createinventory.spec.CreateInventorySpec;

/**
 * This is a Test Class that creates an inventory in VC consisting of
 *  Datacenter Folder
 *     Datacenter1
 *      Cluster1
 *      Host Folder
 *      Storage Folder
 *      Network Folder
 *  Datacenter2
 *   Standalone Host
 *   Cluster2
 *    Clustered Host
 *
 *   Test creates the inventory described and then Connect/Disconnects Standalone Host and Enters/Exits Maintenance Mode for Clustered host
 */

public class CreateInventoryTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      CreateInventorySpec createInventorySpec = new CreateInventorySpec();
      setSpec(createInventorySpec);

      // Folder for the datacenter in folder
      FolderSpec dcFolder = new FolderSpec();
      dcFolder.name.set(RandomStringUtils.randomAlphanumeric(5));
      dcFolder.type.set(FolderType.DATACENTER);

      // Datacenter that will be created in a folder
      DatacenterSpec datacenter1 = new DatacenterSpec();
      datacenter1.name.set(RandomStringUtils.randomAlphanumeric(5));
      // Set its parent the folder
      datacenter1.parent.set(dcFolder);

      // Specification for the cluster in datacenter in folder
      ClusterSpec cluster1 = new ClusterSpec();
      cluster1.name.set(RandomStringUtils.randomAlphanumeric(5));
      // Set the parent DC
      cluster1.parent.set(datacenter1);

      // Specification for the host folder
      FolderSpec hostFolder = new FolderSpec();
      hostFolder.name.set(RandomStringUtils.randomAlphanumeric(5));
      hostFolder.type.set(FolderType.HOST);
      // Set the parent datacenter in folder
      hostFolder.parent.set(datacenter1);

      // Specification for the network folder
      FolderSpec networkFolder = new FolderSpec();
      networkFolder.name.set(RandomStringUtils.randomAlphanumeric(5));
      networkFolder.type.set(FolderType.NETWORK);
      // Set the parent datacenter in folder
      networkFolder.parent.set(datacenter1);

      // Specification for the storage folder
      FolderSpec storageFolder = new FolderSpec();
      storageFolder.name.set(RandomStringUtils.randomAlphanumeric(5));
      storageFolder.type.set(FolderType.STORAGE);
      // Set the parent datacenter in folder
      storageFolder.parent.set(datacenter1);

      // Set the value for the datacenter in the inventory
      DatacenterSpec datacenter2 = new DatacenterSpec();
      datacenter2.name.set(RandomStringUtils.randomAlphanumeric(5));

      // CLuster in the Datacenter
      ClusterSpec cluster2 = new ClusterSpec();
      cluster2.name.set(RandomStringUtils.randomAlphanumeric(5));
      // Parent Datacenter
      cluster2.parent.set(datacenter2);

      List<HostSpec> hosts = new ArrayList<HostSpec>();
      for (String host : testBed.getHosts(2)) {
         hosts.add(HostUtil.buildHostSpec(host,
               testBed.getESXAdminUsername(), testBed.getESXAdminPasssword(),
               443, null));
      }
      // Standalone host to be added
      // Set its parent the datacenter already defined
      hosts.get(0).parent.set(datacenter2);

      // Clustered host specification
      // Set the parent cluster
      hosts.get(1).parent.set(cluster2);

      createInventorySpec.links.add(dcFolder, datacenter1, cluster1, hostFolder, storageFolder, networkFolder, datacenter2, cluster2, hosts.get(0), hosts.get(1));

   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {

      // Step 1: Create DC folder
      composition.appendStep(new CreateDatacenterFolderStep(),
            "Create Datacenter folder");

      // Step 2: Create Datacenters
      composition.appendStep(new CreateDatacenterStep(), "Create Datacenters");

      // Step 3: Create Clusters
      composition.appendStep(new CreateClusterStep(), "Create Clusters");

      // Step 4: Add Hosts
      composition.appendStep(new AddHostStep(), "Add Hosts");

      //TODO need to address folder hierarchies and spec identities
      // Step 5: Create Host Folder
      composition.appendStep(new CreateHostFolderStep(), "Create Host folder");

      // Step 6: Create Network Folder
      composition.appendStep(new CreateNetworkFolderStep(),
            "Create Network folder");

      // Step 7: Create Datastore Folder
      composition.appendStep(new CreateDatastoreFolderStep(),
            "Create Datastore folder");

      // Step 8: Enter Maintenance Mode Clustered host
      composition.appendStep(new BaseWorkflowStep() {
         private HostSpec _host = null;

         @Override
         public void prepare() throws Exception {
            List<HostSpec> hosts = getSpec().links.getAll(HostSpec.class);
            for (HostSpec tempHost : hosts) {
               if (tempHost.parent.get() instanceof ClusterSpec) {
                  _host = tempHost;
                  break;
               }
            }

            if (_host == null) {
               throw new IllegalArgumentException("No clustered host found");
            }
         }

         @Override
         public void execute() throws Exception {
               HostBasicSrvApi.getInstance().enterMaintenanceMode(_host);
         }
      }, "Enter Maintenance Mode Clustered host");

      // Step 9: Exit Maintenance Mode Clustered host
      composition.appendStep(new BaseWorkflowStep() {
         private HostSpec _host = null;

         @Override
         public void prepare() throws Exception {
            List<HostSpec> hosts = getSpec().links.getAll(HostSpec.class);
            for (HostSpec tempHost : hosts) {
               if (tempHost.parent.get() instanceof ClusterSpec) {
                  _host = tempHost;
                  break;
               }
            }

            if (_host == null) {
               throw new IllegalArgumentException("No clustered host found");
            }
         }

         @Override
         public void execute() throws Exception {
               HostBasicSrvApi.getInstance().exitMaintenanceMode(_host);
         }
      }, "Exit Maintenance Mode Clustered host");

      // Step 10: Disconnect Standalone host
      composition.appendStep(new BaseWorkflowStep() {
         private HostSpec _host = null;

         @Override
         public void prepare() throws Exception {
            List<HostSpec> hosts = getSpec().links.getAll(HostSpec.class);
            for (HostSpec tempHost : hosts) {
               if (tempHost.parent.get() instanceof DatacenterSpec) {
                  _host = tempHost;
                  break;
               }
            }

            if (_host == null) {
               throw new IllegalArgumentException("No standalone host found");
            }
         }

         @Override
         public void execute() throws Exception {
               HostBasicSrvApi.getInstance().disconnectHost(_host);
         }
      }, "Disconnect Standalone host");

      // Step 11: Connect Standalone host
      composition.appendStep(new BaseWorkflowStep() {
         private HostSpec _host = null;

         @Override
         public void prepare() throws Exception {
            List<HostSpec> hosts = getSpec().links.getAll(HostSpec.class);
            for (HostSpec tempHost : hosts) {
               if (tempHost.parent.get() instanceof DatacenterSpec) {
                  _host = tempHost;
                  break;
               }
            }

            if (_host == null) {
               throw new IllegalArgumentException("No standalone host found");
            }
         }

         @Override
         public void execute() throws Exception {
               HostBasicSrvApi.getInstance().connectHost(_host);
         }
      }, "Connect Standalone host");

   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      super.execute();
   }
}
