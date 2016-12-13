/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

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

/**
 * Common workflow step for powering on vApps via API.
 * <br>To use this step in automation tests:
 * <li>In the <code>initSpec()</code> method of the
 * <code>BaseTestWorkflow</code> test, create a <code>VappSpec</code> instances
 * and link them to the test spec.
 * <li>Append a <code>PowerOnVappByApiStep</code> instance to the test/prerequisite
 *  workflow composition.
 */
public class PowerOnVappByApiStep extends BaseWorkflowStep {
   private List<VappSpec> _vAppsToPowerOn;
   private List<VappSpec> _vAppsToPowerOff;

   @Override
   public void prepare() throws Exception {
      _vAppsToPowerOn = getSpec().links.getAll(VappSpec.class);

      if (CollectionUtils.isEmpty(_vAppsToPowerOn)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }

      _vAppsToPowerOff = new ArrayList<VappSpec>();
   }

   @Override
   public void execute() throws Exception {
      for (VappSpec vapp : _vAppsToPowerOn) {
         if (VAppSrvApi.getInstance().powerOnVapp(vapp)) {
            _vAppsToPowerOff.add(vapp);
            for (VmSpec vm:vapp.vmList.getAll()){
               verifyFatal(TestScope.FULL, VmSrvApi.getInstance().waitForVmPowerState(vm, true),
                     "Verifying VM was successfully powered on");
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
      for (VappSpec vapp : _vAppsToPowerOff) {
         VAppSrvApi.getInstance().powerOffVapp(vapp, true);
      }
   }

   // TestWorkflowStep  methods
   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vAppsToPowerOn = filteredWorkflowSpec.getAll(VappSpec.class);

      if (CollectionUtils.isEmpty(_vAppsToPowerOn)) {
         throw new IllegalArgumentException(
               "The spec has no links to 'VmSpec' instances");
      }

      _vAppsToPowerOff = new ArrayList<VappSpec>();
   }
}
