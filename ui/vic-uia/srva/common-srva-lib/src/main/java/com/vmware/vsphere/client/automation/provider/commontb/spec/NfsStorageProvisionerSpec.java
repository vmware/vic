/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * NFS storage provider spec
 */
public class NfsStorageProvisionerSpec extends NimbusProvisionerSpec {

    /**
     * Nimbus VM user
     */
    public DataProperty<String> user;

    /**
     * Nimbus VM password
     */
    public DataProperty<String> password;

    /**
     * IP address of the storage
     */
    public DataProperty<String> ip;

    /**
     * IPv4 address of the storage
     */
    public DataProperty<String> ipv4;

    /**
     * IPv6 address of the storage
     */
    public DataProperty<String> ipv6;

    /**
     * The folder that is shared on the storage
     */
    public DataProperty<String> folder;

    /**
     * The name of the storage
     */
    public DataProperty<String> name;
}
