/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;

/**
 * A <code>LocationSpec</code> suitable for modelling a standard navigation in
 * the datacenter related tests.
 */
public class DatacenterLocationSpec extends NGCLocationSpec {

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(
          String datacenterName,
          String primaryTabNId,
          String secondaryTabNId,
          String tocTabNid) {

       super(NGCNavigator.NID_HOME_VCENTER,
             NGCNavigator.NID_VCENTER_DCS,
             datacenterName,
             primaryTabNId,
             secondaryTabNId,
             tocTabNid);
    }

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(
          DatacenterSpec datacenterSpec,
          String primaryTabNId,
          String secondaryTabNId,
          String tocTabNid) {

       super(NGCNavigator.NID_HOME_VCENTER,
             NGCNavigator.NID_VCENTER_DCS,
             datacenterSpec,
             primaryTabNId,
             secondaryTabNId,
             tocTabNid);
    }

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(
            String datacenterName,
            String primaryTabNId,
            String secondaryTabNId) {
       this(datacenterName, primaryTabNId, secondaryTabNId, null);
    }

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(
            DatacenterSpec datacenterSpec,
            String primaryTabNId,
            String secondaryTabNId) {
       this(datacenterSpec, primaryTabNId, secondaryTabNId, null);
    }

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(String datacenterName, String primaryTabNId) {
       this(datacenterName, primaryTabNId, null, null);
    }

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(String datacenterName) {
       this(datacenterName, null, null, null);
    }

    /**
     * Build a location path based on the provided datacenter navigation identifiers.
     */
    public DatacenterLocationSpec(DatacenterSpec datacenterSpec) {
       this(datacenterSpec, null, null, null);
    }

    /**
     * Build a location that will navigate the UI to the datacenter entity view.
     */
    public DatacenterLocationSpec() {
       super(NGCNavigator.NID_HOME_VCENTER, NGCNavigator.NID_VCENTER_DCS);
    }
}
