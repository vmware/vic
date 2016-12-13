/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.storage.lib.core.steps;

import java.util.ArrayList;
import java.util.List;

import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.HostBasicSrvApi;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;

public class AttachIscsiTargetStep extends EnhancedBaseWorkflowStep {

   @UsesSpec()
   private List<HostSpec> _hostsToAttachDatastoreTo;
   private final List<HostSpec> _hostsToRemoveDatastoreFrom = new ArrayList<>();

   @Override
   public void execute() throws Exception {
      for (HostSpec hostSpec : _hostsToAttachDatastoreTo) {
         if (!HostBasicSrvApi.getInstance().createIscsiAdapter(hostSpec)) {
            throw new Exception(String.format(
                  "Unable to add iSCSI adapter to host '%s'",
                  hostSpec.name.get()));
         }
         _hostsToRemoveDatastoreFrom.add(hostSpec);
      }

   }

   @Override
   public void clean() throws Exception {
      for (HostSpec host : _hostsToRemoveDatastoreFrom) {
         HostBasicSrvApi.getInstance().destroyIscsiAdapter(host);
      }
   }
}
