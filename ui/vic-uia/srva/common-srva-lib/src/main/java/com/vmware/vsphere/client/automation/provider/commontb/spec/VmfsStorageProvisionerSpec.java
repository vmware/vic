/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * VMFS shared storage provider spec.
 */
public class VmfsStorageProvisionerSpec extends NimbusProvisionerSpec {

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
}
