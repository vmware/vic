/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.step;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

/**
 * Step to move a VM<br>
 * Requires <code>VmSpec</code> and one or two of the following:<br>
 * <br>
 * <code>HostSpec</code><br>
 * <code>DatastoreSpec</code>
 */
public class MoveVMStep extends BaseWorkflowStep {
   VmSpec _vm;
   HostSpec _hostDestination;
   DatastoreSpec _datastoreDestination;
   DatastoreSpec _datastoreSource;

   /**
    * @inheritDoc
    */
   @Override
   public void prepare() throws Exception {
      _vm = getSpec().links.get(VmSpec.class);
      _hostDestination = getSpec().links.get(HostSpec.class);
      _datastoreDestination = getSpec().links.get(DatastoreSpec.class);

      // verify VM spec
      VmSrvApi.getInstance().validateVmSpec(_vm);

      if ((_hostDestination == null) && (_datastoreDestination == null)) {
         throw new Exception("HostSpec and DatastoreSpec are both null."
               + "One of them must be set.");
      }
   }

   @Override
   /**
    * @inheritDoc
    */
   public void execute() throws Exception {
      _datastoreSource = VmSrvApi.getInstance().getVmDatastore(_vm);

      boolean isMoved = VmSrvApi.getInstance().migrateVm(_vm, _hostDestination, _datastoreDestination);
      verifyFatal(TestScope.BAT, isMoved, "The VM has been moved");
   }

   @Override
   public void clean() throws Exception {
      if (_datastoreDestination != null) {
         VmSrvApi.getInstance().migrateVm(_vm, null, _datastoreSource);
      }
   }
}
