package com.vmware.vsphere.client.automation.storage.lib.core.steps;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;
import com.vmware.vsphere.client.automation.common.step.EnhancedBaseWorkflowStep;

public abstract class BasePrerequisiteStep extends EnhancedBaseWorkflowStep {

   /**
    * Get a vmodl service from a mor and a service spec
    *
    * @param serviceMor
    * @param serviceSpec
    *
    * @return
    * @throws Exception
    */
   // TODO This wrapper is added for easier refactor. Extract this method to
   // an API layer
   protected <T extends ManagedObject> T getService(
         ManagedObjectReference serviceMor, ServiceSpec serviceSpec)
         throws Exception {

      return ManagedEntityUtil.getManagedObjectFromMoRef(serviceMor,
            serviceSpec);
   }
}
