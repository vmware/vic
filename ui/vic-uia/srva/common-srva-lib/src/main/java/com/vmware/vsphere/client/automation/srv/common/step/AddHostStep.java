/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;

/**
 * A class that is used for adding all the hosts, specified in a test,
 * to the VC Inventory
 */
public class AddHostStep extends BaseWorkflowStep {

   private List<HostSpec> _hostsToAdd;
   private List<HostSpec> _hostsToRemove;

   @Override
   public void prepare() throws Exception {

      // Get all Hosts' specs
      _hostsToAdd = getSpec().links.getAll(HostSpec.class);

      if (_hostsToAdd == null || _hostsToAdd.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'HostSpec' instances");
      }

      // The delete method needs to know the datacenter that holds the host
      _hostsToRemove = new ArrayList<HostSpec>();
   }

   @Override
   public void execute() throws Exception {
      for (HostSpec hostSpec : _hostsToAdd) {
         if (!HostBasicSrvApi.getInstance().addHost(hostSpec, true)) {
            throw new Exception(String.format("Unable to add host '%s'",
                  hostSpec.name.get()));
         }
         _hostsToRemove.add(hostSpec);
      }
   }

   @Override
   public void clean() throws Exception {
      for (HostSpec host : _hostsToRemove) {
         HostBasicSrvApi.getInstance().deleteHostSafely(host);
      }
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {

      // Get all Hosts' specs
      _hostsToAdd = filteredWorkflowSpec.getAll(HostSpec.class);

      if (_hostsToAdd == null || _hostsToAdd.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'HostSpec' instances");
      }

      // The delete method needs to know the datacenter that holds the host
      _hostsToRemove = new ArrayList<HostSpec>();
   }
}
