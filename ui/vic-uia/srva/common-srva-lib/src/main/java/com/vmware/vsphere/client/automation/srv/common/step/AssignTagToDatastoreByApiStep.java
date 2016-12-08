/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotEmpty;
import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.Collections;
import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.BackingTagsBasicSrvApi;

/**
 * Assigning tag to datastore.
 */
public class AssignTagToDatastoreByApiStep extends BaseWorkflowStep {

   private BackingCategorySpec _category;
   private BackingTagSpec _tag;
   private List<DatastoreSpec> _datastores;

   @Override
   public void prepare() {
      _category = getSpec().get(BackingCategorySpec.class);
      ensureNotNull(_category, "BackingCategorySpec object is missing.");

      _tag = getSpec().get(BackingTagSpec.class);
      ensureNotNull(_tag, "BackingTagSpec object is missing.");

      _datastores = Collections.singletonList(getSpec().get(DatastoreSpec.class));
      ensureNotEmpty(_datastores, "DatastoreSpec objects are missing.");
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _category = filteredWorkflowSpec.get(BackingCategorySpec.class);
      ensureNotNull(_category, "BackingCategorySpec object is missing.");

      _tag = filteredWorkflowSpec.get(BackingTagSpec.class);
      ensureNotNull(_tag, "BackingTagSpec object is missing.");

      _datastores = Collections.singletonList(filteredWorkflowSpec
            .get(DatastoreSpec.class));
      ensureNotEmpty(_datastores, "DatastoreSpec objects are missing.");
   }

   @Override
   public void execute() throws Exception {
      BackingTagsBasicSrvApi.getInstance().attachResources(_category, _tag, _datastores);
   }

   @Override
   public void clean() throws Exception {
      BackingTagsBasicSrvApi.getInstance().detachResources(_category, _tag, _datastores);
   }
}
