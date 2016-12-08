/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.migrate.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Spec for Migrate VM
 *
 */
public class MigrateVmSpec extends BaseSpec {

    /**
     * Types of VM creation through UI, e.g. clone existing or create brand new
     *
     */
    public static enum TargetMigrationTypes {
        HOSTS, CLUSTERS, RESOURCE_POOLS, VAPPS
    };

    /**
     * Target entity for the migration
     */
    public DataProperty<ManagedEntitySpec> targetEntity;

    /**
     * Type of the target entity (host, datastore, etc)
     */
    public DataProperty<TargetMigrationTypes> targetEntityType;
}