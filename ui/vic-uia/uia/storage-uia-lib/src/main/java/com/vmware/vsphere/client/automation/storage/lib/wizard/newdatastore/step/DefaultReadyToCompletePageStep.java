/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import java.util.List;

import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.DatastoreBasicSrvApi;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core.NewDatastoreWizardStep;

/**
 * NewDatastoreWizardStep implementation for Ready to complete page in the new
 * datastore wizard
 */
public class DefaultReadyToCompletePageStep extends NewDatastoreWizardStep {

    @UsesSpec()
    protected List<TaskSpec> taskSpecs;

    @Override
    public void executeWizardOperation(WizardNavigator wizardNavigator) {
        // Finish the wizard
        boolean finishWizard = wizardNavigator.finishWizard();
        verifyFatal(TestScope.BAT, finishWizard, "Verify wizard is closed");

        // Wait for triggered recent tasks to complete
        boolean allTasksCompleted = true;
        for (TaskSpec taskSpec : taskSpecs) {
           final boolean taskCompleted = new TasksUtil().waitForRecentTaskToMatchFilter(new RecentTaskFilter(taskSpec));
           allTasksCompleted = allTasksCompleted && taskCompleted;
           verifySafely(taskCompleted, String.format("Task with name '%s' not found", taskSpec.name.get()));
        }
        verifySafely(allTasksCompleted, "Verifying the completion of all the triggered tasks");
    }

    /**
     * Delete the newly created datastore.
     */
    @Override
    public void clean() throws Exception {
        DatastoreBasicSrvApi.getInstance().deleteDatastoreSafely(datastoreSpec);
    }
}