/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Host provider spec
 */
public class HostProvisionerSpec extends NimbusProvisionerSpec {

   public DataProperty<String> ip;
   public DataProperty<String> user;
   public DataProperty<String> password;

   /**
    * Build info
    */
   public DataProperty<String> product;
   public DataProperty<String> branch;
   public DataProperty<String> buildNumber;
}
