/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.view;

import com.vmware.client.automation.components.navigator.MultiPageDialogNavigator;

/**
 * This class represent the parent for the cluster multi-page dialogs.
 * It provides ability to navigate to the dialog steps.
 */
public class EditClusterPageNavigator extends MultiPageDialogNavigator {

    private static final String DRS_PAGE_ID = "step_drsConfigForm";
    private static final String HA_PAGE_ID = "step_haConfigForm";

    /**
     * Goes to the vSphere DRS edit settings page.
     *
     * @return  true if the navigation to the specified page was successful
     */
    public boolean goToDrsPage() {
        return new MultiPageDialogNavigator().goToPage(DRS_PAGE_ID);
    }

    /**
     * Goes to the vSphere HA edit settings page.
     *
     * @return  true if the navigation to the specified page was successful
     */
    public boolean goToHaPage() {
        return new MultiPageDialogNavigator().goToPage(HA_PAGE_ID);
    }
}
