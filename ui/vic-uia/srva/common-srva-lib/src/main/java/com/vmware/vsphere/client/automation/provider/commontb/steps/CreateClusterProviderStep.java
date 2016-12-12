/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.steps;

import java.util.Set;

import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.SpecTraversalUtil;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;

/**
 * Provider work flow step to create a Cluster.
 */
public class CreateClusterProviderStep implements ProviderWorkflowStep {

   private ClusterSpec _clusterToCreate;

   @Override
   /**
    * Init the cluster spec of the cluster to be created.
    */
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssmblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) {
      if (isAssembling) {
         Set<ClusterSpec> setOfDatacenters = SpecTraversalUtil
               .getAllSpecsFromContainerTree(filteredAssmblerSpec,
                     ClusterSpec.class);
         _clusterToCreate = setOfDatacenters.iterator().next();

         // _clusterToCreate =
         // filteredAssmblerSpec.links.get(ClusterSpec.class);
      } else {
         _clusterToCreate = filteredPublisherSpec.links.get(ClusterSpec.class);
      }
   }

   @Override
   /**
    * Check if the cluster already exist if not creates one.
    */
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      if (!ClusterBasicSrvApi.getInstance().checkClusterExists(_clusterToCreate)) {
         if (!ClusterBasicSrvApi.getInstance().createCluster(_clusterToCreate)) {
            throw new Exception(String.format("Unable to create cluster '%s'",
                  _clusterToCreate.name.get()));
         }
      }

   }

   @Override
   /**
    * Return true if cluster specified by the step spec exists.
    */
   public boolean checkHealth() throws Exception {
      return ClusterBasicSrvApi.getInstance().checkClusterExists(_clusterToCreate);
   }

   @Override
   /**
    * Deletes the cluster specified by the step spec.
    */
   public void disassemble() throws Exception {
      ClusterBasicSrvApi.getInstance().deleteClusterSafely(_clusterToCreate);
   }

}
