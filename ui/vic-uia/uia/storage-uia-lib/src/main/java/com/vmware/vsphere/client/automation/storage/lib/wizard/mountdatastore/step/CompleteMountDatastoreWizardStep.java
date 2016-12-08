package com.vmware.vsphere.client.automation.storage.lib.wizard.mountdatastore.step;

import java.util.List;

import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.SinglePageWizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;

/**
 * {@link NewDatastoreWizardStep} implementation for completion of the mount
 * datastore to additional hosts wizard
 */
public class CompleteMountDatastoreWizardStep extends SinglePageWizardStep {

   // TODO: The hosts associated with this steps may differ from the hosts used
   // on the attached step.
   @UsesSpec
   private List<HostSpec> attachedHosts;

   @UsesSpec
   private DatastoreSpec datastoreSpec;

   @UsesSpec
   private TaskSpec taskSpec;

   @Override
   protected void executeWizardOperation(
         SinglePageDialogNavigator wizardNavigator) throws Exception {
      wizardNavigator.waitForLoadingProgressBar();
      wizardNavigator.clickOk();

      final boolean taskCompleted = new TasksUtil()
            .waitForRecentTaskToMatchFilter(new RecentTaskFilter(taskSpec));
      verifySafely(
            taskCompleted,
            String.format("Task with name '%s' is completed ",
                  taskSpec.name.get()));

   }

   @Override
   public void clean() throws Exception {
      for (HostSpec hostSpec : attachedHosts) {
         DatastoreBasicSrvApi.getInstance().unmountDatastore(datastoreSpec,
               hostSpec);
      }
   }
}
