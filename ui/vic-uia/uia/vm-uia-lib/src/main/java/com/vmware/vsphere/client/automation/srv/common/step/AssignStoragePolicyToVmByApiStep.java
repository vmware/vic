/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmPolicySrvApi;

/**
 * Assigns storage policies to a given VMs using the vAPI.<br>
 * <p>
 * Require <code>VmSpec(s), StoragePolicySpec</code>
 */
public class AssignStoragePolicyToVmByApiStep extends BaseWorkflowStep {
   private StoragePolicySpec _policy;
   private VmSpec _vm;

   @Override
   public void prepare() {
      _vm = getSpec().links.get(VmSpec.class);
      _policy = getSpec().links.get(StoragePolicySpec.class);

      if (_vm == null) {
         throw new IllegalArgumentException("VmSpec not found");
      }

      if (_policy == null) {
         throw new IllegalArgumentException("PlacementPolicySpec not found");
      }
   }

   @Override
   public void execute() throws Exception {
      VmPolicySrvApi.getInstance().addStoragePolicy(_vm, _policy);
   }

   @Override
   public void clean() throws Exception {
      VmPolicySrvApi.getInstance().removeStoragePolicy(_vm, _policy);
   }

}
