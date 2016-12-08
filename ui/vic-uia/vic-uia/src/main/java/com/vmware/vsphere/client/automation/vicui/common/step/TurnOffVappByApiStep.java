package com.vmware.vsphere.client.automation.vicui.common.step;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VAppSrvApi;
import com.vmware.vsphere.client.automation.srv.common.srvapi.VmSrvApi;

public class TurnOffVappByApiStep extends BaseWorkflowStep {
   private List<VappSpec> _vAppsToPowerOn;
   private List<VappSpec> _vAppsToPowerOff;

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vAppsToPowerOff = filteredWorkflowSpec.getAll(VappSpec.class);

      if (CollectionUtils.isEmpty(_vAppsToPowerOff)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VappSpec' instances");
      }

      _vAppsToPowerOn = new ArrayList<VappSpec>();
   }

   @Override
   public void execute() throws Exception {
      for (VappSpec vapp : _vAppsToPowerOff) {
         if (VAppSrvApi.getInstance().powerOffVapp(vapp, true)) {
            _vAppsToPowerOn.add(vapp);
            for (VmSpec vm:vapp.vmList.getAll()){
                verifyFatal(TestScope.FULL, VmSrvApi.getInstance().waitForVmPowerState(vm, false),
                    "Verifying VM was successfully powered off");
            }
         } else {
            throw new Exception(
                  String.format(
                        "Unable to power on vApp with name '%s'", vapp.name.get()));
         }
      }
   }

   @Override
   public void clean() throws Exception {
//      for (VappSpec vapp : _vAppsToPowerOn) {
//         VAppSrvApi.getInstance().powerOnVapp(vapp);
//      }
   }
}
