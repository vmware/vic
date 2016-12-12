/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.steps;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * Provider work flow step to remove and readd local vmfs storage of hosts.
 * This is necessary as the nimbus scripts have a flaw and provide hosts
 * with identical local storages that prevent the addition of more then
 * one such hosts to a single VC Inventory.
 */
public class RemoveReaddHostLocalDatastoreProviderStep implements ProviderWorkflowStep {

   private HostSpec _hostToReaddStorageSpec;

   @Override
   /**
    * Init the host spec with details of the host to be attached.
    */
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssmblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      if (isAssembling) {
         _hostToReaddStorageSpec = filteredAssmblerSpec.links.get(HostSpec.class);
      } else {
         _hostToReaddStorageSpec = filteredPublisherSpec.links.get(HostSpec.class);
      }

   }

   @Override
   /**
    * Remove and readd host's local storages
    */
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      if (HostBasicSrvApi.getInstance().checkHostExists(_hostToReaddStorageSpec)) {
         List<DatastoreSpec> listDatastoreSpecs = HostBasicSrvApi.getInstance().getHostDatastoreSpecs(_hostToReaddStorageSpec);
         List<DatastoreSpec> listDatastoresToReadd = new ArrayList<DatastoreSpec>();

         // find all VMFS datastores
         for (DatastoreSpec dsSpec :listDatastoreSpecs) {
            if (dsSpec.type.get().equals(DatastoreType.VMFS)) {
               listDatastoresToReadd.add(dsSpec);
            }
         }

         // remove all VMFS datastores
         for (DatastoreSpec dsSpec : listDatastoresToReadd) {
            DatastoreBasicSrvApi.getInstance().deleteDatastoreSafely(dsSpec);
         }

         // add again all VMFS datastores
         for (DatastoreSpec dsSpec : listDatastoresToReadd) {
            DatastoreBasicSrvApi.getInstance().createDatastore(dsSpec);
         }
      }
   }

   @Override
   /**
    * Return true if the host specified is connected.
    */
   public boolean checkHealth() throws Exception {
      // nothing here
      return true;
   }

   @Override
   /**
    * Remove the host from the inventory.
    */
   public void disassemble() throws Exception {
      HostBasicSrvApi.getInstance().deleteHostSafely(_hostToReaddStorageSpec);

   }

}
