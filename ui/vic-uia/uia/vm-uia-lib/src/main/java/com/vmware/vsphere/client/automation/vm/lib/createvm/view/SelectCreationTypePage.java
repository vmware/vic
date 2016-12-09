/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.lib.createvm.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.vm.lib.createvm.spec.CreateVmSpec.VmCreationType;

/**
 * Select creation type for the VM
 */
public class SelectCreationTypePage extends WizardNavigator {
    private static final String VM_CREATION_LIST = "optionsList";
    private static String PAGE_TITLE = CommonUtil
          .getLocalizedString("createVm.wizard.creationTypePage.title");

    /**
     * Validates that the view is present on the screen, before executing any
     * actions on it.
     */
    public void validate() {
        if (!verifyPageIndex()) {
            throw new IllegalStateException("Unexpected dialog title. Expected was: " + PAGE_TITLE);
        }
    }

    /**
     * Checks if the Select Creation Type page is open, by verifying the page index.
     *
     * @return True if the page index corresponds to Select Creation Type page, false
     *         otherwise
     */
    public boolean verifyPageIndex() {
        waitForDialogToLoad();
        return getCurrentlyActivePage().equals(1);
    }

    /**
     * Selects creation type for a new vm wizard, e.g. clone existing or create brand new
     *
     * @param creationType
     */
    public void selectCreationType(VmCreationType creationType) {
        UI.component.selectByIndex(creationType.ordinal(), VM_CREATION_LIST);
    }
}
