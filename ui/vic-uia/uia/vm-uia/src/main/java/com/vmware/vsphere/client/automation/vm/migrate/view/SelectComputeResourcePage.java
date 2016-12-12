/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.view;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.vsphere.client.automation.vm.common.VmUtil;
import com.vmware.vsphere.client.automation.vm.migrate.spec.MigrateVmSpec;

/**
 * UI model of the 'Select Compute Resource' step
 * from 'Migrate VM' wizard
 */
public class SelectComputeResourcePage extends WizardNavigator {

    private static final IDGroup GRID_ID = IDGroup.toIDGroup("tiwoDialog/list");
    private static final String NAME_COLUMN = VmUtil
        .getLocalizedString("migrateVm.wizard.targetResourceGrid.nameColumn");

    /**
     * Select destination host from the object selector grid.
     *
     * @param hostSpec target host spec
     * @return true if the selection was successful, false otherwise
     */
    public boolean selectEntity(MigrateVmSpec migrateSpec) {
        switch (migrateSpec.targetEntityType.get()) {
        case HOSTS:
            return GridControl.selectEntity(getGrid(), NAME_COLUMN, migrateSpec.targetEntity.get().name.get());
        default:
            throw new IllegalArgumentException("Unknown targetEntityType passed!");
        }
    }

    // ---------------------------------------------------------------------------
    // Private methods

    /**
     * Finds and returns the advanced data grid on 'Select Compute Resource' wizard step.
     */
    private AdvancedDataGrid getGrid() {
        return GridControl.findGrid(IDGroup.toIDGroup(GRID_ID));
    }
}
