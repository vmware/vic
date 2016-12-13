/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.steps;

import java.util.List;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;

/**
 * Provider work flow step to attach a datastore.
 */
public class AttachStorageProviderStep implements ProviderWorkflowStep {

   private List<DatastoreSpec> _datastoresToAttachSpec;

   @Override
   /**
    * Init the datastore spec with details of the datastore to be attached.
    */
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssmblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      if (isAssembling) {
         _datastoresToAttachSpec = filteredAssmblerSpec.links.getAll(DatastoreSpec.class);
      } else {
         _datastoresToAttachSpec = filteredPublisherSpec.links.getAll(DatastoreSpec.class);
      }

   }

   @Override
   /**
    * Check if the datastore specified by the step spec data is attached. If not attach it.
    */
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      for (DatastoreSpec dsSpec : _datastoresToAttachSpec) {
         if (!DatastoreBasicSrvApi.getInstance().checkDatastoreExists(dsSpec)) {
            if (!DatastoreBasicSrvApi.getInstance().createDatastore(dsSpec)) {
               throw new Exception(String.format(
                     "Unable to add datastore '%s'",
                     dsSpec.name.get()));
            }
         }
      }
   }

   @Override
   /**
    * Return true if the datastore specified is connected.
    */
   public boolean checkHealth() throws Exception {
      // TODO: lgrigorova how should check health be realized for a datastore
      return true;
   }

   @Override
   /**
    * Remove the datastore from the inventory.
    */
   public void disassemble() throws Exception {
      for (DatastoreSpec dsSpec : _datastoresToAttachSpec) {
         DatastoreBasicSrvApi.getInstance().deleteDatastoreSafely(dsSpec);
      }
   }

}
