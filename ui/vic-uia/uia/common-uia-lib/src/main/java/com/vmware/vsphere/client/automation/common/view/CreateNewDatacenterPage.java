/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.components.navigator.SinglePageDialogNavigator;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Represents the create new datacenter page.
 */
public class CreateNewDatacenterPage extends SinglePageDialogNavigator {

    private static final IDGroup NAME_TF_ID = IDGroup.toIDGroup("inputLabel");
    private static final IDGroup NAVIGATION_TREE_ID = IDGroup.toIDGroup("navTree");

    /**
     * Type datacenter name in the name text field.
     *
     * @param datacenterName       datacenter name to be typed in
     */
    public void setDatacenterName(String datacenterName) {
        UI.component.value.set(datacenterName, NAME_TF_ID);
    }

    /**
     * Select vCenter or folder component.
     */
    public void selectVcenterOrFolder(int index) {
        UI.component.selectByIndex(index, NAVIGATION_TREE_ID);
    }
}
