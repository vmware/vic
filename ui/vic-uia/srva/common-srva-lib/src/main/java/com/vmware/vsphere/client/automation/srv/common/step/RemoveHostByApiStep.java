/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;
import org.apache.commons.collections4.CollectionUtils;

import java.util.ArrayList;
import java.util.List;

/**
 * A class that is used for removing all the hosts, specified in a test,
 * from the VC Inventory
 */
public class RemoveHostByApiStep extends BaseWorkflowStep {

   private List<HostSpec> _hostsToRemove;
   private List<HostSpec> _hostsToAdd;

   @Override
   public void prepare() throws Exception {

      // Get all Hosts' specs
      _hostsToRemove = getSpec().links.getAll(HostSpec.class);

      if (CollectionUtils.isEmpty(_hostsToRemove)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'HostSpec' instances");
      }

      _hostsToAdd = new ArrayList<HostSpec>();
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _hostsToRemove = filteredWorkflowSpec.getAll(HostSpec.class);
      TestSpecValidator.ensureNotEmpty(_hostsToRemove, "Host Spec is null");
      _hostsToAdd = new ArrayList<>();
   }

   @Override
   public void execute() throws Exception {
      for (HostSpec host: _hostsToRemove) {
         if (!HostBasicSrvApi.getInstance().deleteHostSafely(host)) {
            String errorTemplate = "Unable to remove host '%s'";
            String hostIp = host.name.get();
            String errorMessage = String.format(errorTemplate, hostIp);
            throw new Exception(errorMessage);
         }
         _hostsToAdd.add(host);
      }
   }

   @Override
   public void clean() throws Exception {
      for (HostSpec hostSpec : _hostsToAdd) {
         if (!HostBasicSrvApi.getInstance().addHost(hostSpec, true)) {
            String errorMessage = "Unable to add host: %s";
            String hostIp = hostSpec.name.get();
            throw new RuntimeException(String.format(errorMessage, hostIp));
         }
      }
   }
}
