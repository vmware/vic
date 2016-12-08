/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.lib.createvm.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

/**
 * Spec for new VM creation
 *
 */
public class CreateVmSpec extends VmSpec {
    /**
     * Types of VM creation through UI, e.g. clone existing or create brand new
     *
     */
    public static enum VmCreationType {
        CREATE_NEW_VM, DEPLOY_FROM_TEMPLATE, CLONE_EXISTING_VM, CLONE_TO_TEMPLATE, CLONE_TEMPLATE_TO_TEMPLATE, CONVERT_TEMPLATE_TO_VM
    };

    /**
     * VM creation type
     */
    public DataProperty<VmCreationType> creationType;

    /**
     * Destination compute resource for the VM
     */
    public DataProperty<ManagedEntitySpec> computeResource;

    /**
     * Datastore cluster on which to place the VM
     */
    public DataProperty<DatastoreClusterSpec> datastoreCluster;
}