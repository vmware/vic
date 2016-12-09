/**
 * Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.common.step.ReconfigureManagedEntityStep;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ClusterBasicSrvApi;

/**
 * Method that allows cluster reconfiguration based on spec versioning.
 */
public class ReconfigureClusterByApiStep extends ReconfigureManagedEntityStep {

    @Override
    public void execute() throws Exception {
        _reconfiguredManagedEntitySpecs.get(0).name.set(_originalManagedEntitySpecs.get(0).name.get());
        _reconfiguredManagedEntitySpecs.get(0).parent.set(_originalManagedEntitySpecs.get(0).parent.get());

        ClusterBasicSrvApi.getInstance().reconfigureCluster(
            (ClusterSpec) _originalManagedEntitySpecs.get(0),
            (ClusterSpec) _reconfiguredManagedEntitySpecs.get(0));
    }

    @Override
    public void clean() throws Exception {
        ClusterBasicSrvApi.getInstance().reconfigureCluster(
            (ClusterSpec) _reconfiguredManagedEntitySpecs.get(0),
            (ClusterSpec) _originalManagedEntitySpecs.get(0));
    }
}
