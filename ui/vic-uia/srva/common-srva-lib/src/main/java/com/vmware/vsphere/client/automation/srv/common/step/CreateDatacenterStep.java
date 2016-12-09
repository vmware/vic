package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatacenterBasicSrvApi;

public class CreateDatacenterStep extends BaseWorkflowStep {

   private List<DatacenterSpec> _datacentersToCreate;
   private List<DatacenterSpec> _datacentersToDelete;

   @Override
   /**
    * @inheritDoc
    */
   public void prepare() throws Exception {

      _datacentersToCreate = getSpec().links.getAll(DatacenterSpec.class);
      getSpec().links.getAll(BaseSpec.class);
      if (_datacentersToCreate == null || _datacentersToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatacenterSpec' instances");
      }

      _datacentersToDelete = new ArrayList<DatacenterSpec>();
   }

   @Override
   public void execute() throws Exception {

      for (DatacenterSpec datacenterSpec : _datacentersToCreate) {
         if (!DatacenterBasicSrvApi.getInstance().createDatacenter(datacenterSpec)) {
            throw new Exception(String.format(
                  "Unable to create datacenter '%s'",
                  datacenterSpec.name.get()));
         }
         _datacentersToDelete.add(datacenterSpec);
      }

   }

   @Override
   /**
    * @inheritDoc
    */
   public void clean() throws Exception {
      for (DatacenterSpec datacenterSpec : _datacentersToDelete) {
         DatacenterBasicSrvApi.getInstance().deleteDatacenterSafely(datacenterSpec);
      }
   }

   // TestWorkflowStep methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _datacentersToCreate = filteredWorkflowSpec.links.getAll(DatacenterSpec.class);
      if (_datacentersToCreate == null || _datacentersToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'DatacenterSpec' instances");
      }

      _datacentersToDelete = new ArrayList<DatacenterSpec>();

   }

}
