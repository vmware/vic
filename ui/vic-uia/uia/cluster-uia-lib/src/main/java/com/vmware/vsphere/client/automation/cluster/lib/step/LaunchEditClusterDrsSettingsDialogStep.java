/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.step;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import com.vmware.client.automation.workflow.CommonUIWorkflowStep;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.cluster.lib.view.ClusterDrsSettingsView;
import com.vmware.vsphere.client.automation.cluster.lib.view.EditClusterDrsSettingsPage;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;

/**
 * Launch edit vSphere DRS settings dialog.
 * This step expects a specific view to be opened
 * cluster -> manage -> settings -> vSphere DRS
 *
 * Operation performed by this step:
 *  1. Click on edit vSphere DRS button
 *  2. Verify the edit cluster dialog has been opened
 *  3. Verify vSphere DRS tab is opened
 */
public class LaunchEditClusterDrsSettingsDialogStep extends CommonUIWorkflowStep {

    private static final String DIALOG_TITLE_FORMAT =
            CommonUtil.getLocalizedString("editCluster.dialog.title.format");
    private ClusterSpec _clusterSpec;

    @Override
    public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
       _clusterSpec = filteredWorkflowSpec.get(ClusterSpec.class);

       ensureNotNull(_clusterSpec, "ClusterSpec spec is missing.");
    }

    @Override
    public void execute() throws Exception {

        // launch edit cluster settings dialog
        new ClusterDrsSettingsView().clickEditDrsSettingsButton();

        final EditClusterDrsSettingsPage editPage = new EditClusterDrsSettingsPage();

        editPage.waitForDialogToLoad();

        // Verify the edit dialog is opened
        verifyFatal(
                editPage.isOpen(),
                "Edit dialog is opened."
            );

        // Verify the title of the dialog
        verifySafely(new TestScopeVerification() {

                @Override
                public boolean verify() throws Exception {
                    String expectedTitle =
                            String.format(DIALOG_TITLE_FORMAT, _clusterSpec.name.get());

                    return editPage.getTitle().equals(expectedTitle);
                }
            }, "Verify the title of the edit-cluster dialog");
    }
}
