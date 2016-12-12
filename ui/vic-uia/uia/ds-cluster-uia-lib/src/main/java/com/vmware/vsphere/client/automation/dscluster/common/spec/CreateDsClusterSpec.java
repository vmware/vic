/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.dscluster.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;

/**
 * Container for test-specific data for create datastore cluster tests.
 *
 */
public class CreateDsClusterSpec extends BaseSpec {

    /**
     * Host parents of the datastores that will be members of the datastore cluster
     */
    public DataProperty<ManagedEntitySpec> datastoresParents;

    /**
     * Members of the datastore cluster
     */
    public DataProperty<DatastoreSpec> datastores;
}
