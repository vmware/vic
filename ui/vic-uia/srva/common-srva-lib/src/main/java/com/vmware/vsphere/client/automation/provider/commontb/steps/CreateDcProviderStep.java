/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.steps;

import java.util.Set;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.SpecTraversalUtil;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatacenterBasicSrvApi;

/**
 * Provider work flow step to create a Data Center.
 */
public class CreateDcProviderStep implements ProviderWorkflowStep {

   private DatacenterSpec _datacenterToCreate;

   @Override
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssmblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {

      if (isAssembling) {
         Set<DatacenterSpec> setOfDatacenters = SpecTraversalUtil
               .getAllSpecsFromContainerTree(filteredAssmblerSpec,
                     DatacenterSpec.class);
         _datacenterToCreate = setOfDatacenters.iterator().next();
         // _datacenterToCreate =
         // filteredAssmblerSpec.links.get(DatacenterSpec.class);
      } else {
         _datacenterToCreate = filteredPublisherSpec.links
               .get(DatacenterSpec.class);
      }
   }

   @Override
   /**
    * Delete DC specified by the DatacenterSpec.
    */
   public void disassemble() throws Exception {
      DatacenterBasicSrvApi.getInstance().deleteDatacenterSafely(_datacenterToCreate);
   }

   @Override
   /**
    * Return true if Data Center described by the DatacenterSpec exists.
    */
   public boolean checkHealth() throws Exception {
      return DatacenterBasicSrvApi.getInstance().checkDatacenterExists(_datacenterToCreate);
   }

   @Override
   /**
    * Check if DC specified by the DatacenterSpec exist. If it does not exist
    * create one.
    */
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      if (!DatacenterBasicSrvApi.getInstance().checkDatacenterExists(_datacenterToCreate)) {
         if (!DatacenterBasicSrvApi.getInstance().createDatacenter(_datacenterToCreate)) {
            throw new Exception(String.format(
                  "Unable to create datacenter '%s'",
                  _datacenterToCreate.name.get()));
         }
      }

   }

}
