/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.DvsSpec;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation in
 * the distributed switch related tests.
 */
public class DvsLocationSpec extends NGCLocationSpec {

    /**
     * Build a location path based on the provided cluster navigation identifiers.
     */
    public DvsLocationSpec(String dvsName, String primaryTabNId, String secondaryTabNId, String tocTabNid) {

        super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DVSES, dvsName, primaryTabNId, secondaryTabNId,
            tocTabNid);
    }

    /**
     * Build a location path based on the provided cluster navigation identifiers.
     */
    public DvsLocationSpec(DvsSpec dvsSpec, String primaryTabNId, String secondaryTabNId, String tocTabNid) {

        super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DVSES, dvsSpec, primaryTabNId, secondaryTabNId,
            tocTabNid);
    }

    /**
     * Build a location path based on the provided distributed switch navigation identifiers.
     */
    public DvsLocationSpec(String dvsName, String primaryTabNId, String secondaryTabNId) {
        this(dvsName, primaryTabNId, secondaryTabNId, null);
    }

    /**
     * Build a location path based on the provided distributed switch navigation identifiers.
     */
    public DvsLocationSpec(DvsSpec dvsSpec, String primaryTabNId, String secondaryTabNId) {
        this(dvsSpec, primaryTabNId, secondaryTabNId, null);
    }

    /**
     * Build a location path based on the provided distributed switch navigation identifiers.
     */
    public DvsLocationSpec(String dvsName, String primaryTabNId) {
        this(dvsName, primaryTabNId, null, null);
    }

    /**
     * Build a location path based on the provided distributed switch navigation identifiers.
     */
    public DvsLocationSpec(DvsSpec dvsSpec, String primaryTabNId) {
        this(dvsSpec, primaryTabNId, null, null);
    }


    /**
     * Build a location path based on the provided distributed switch navigation identifiers.
     */
    public DvsLocationSpec(String dvsName) {
        this(dvsName, null, null, null);
    }

    /**
     * Build a location path based on the provided distributed switch navigation identifiers.
     */
    public DvsLocationSpec(DvsSpec dvsSpec) {
        this(dvsSpec, null, null, null);
    }

    /**
     * Build a location that will navigate the UI to the distributed switch entity view.
     */
    public DvsLocationSpec() {
        super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DVSES);
    }
}
