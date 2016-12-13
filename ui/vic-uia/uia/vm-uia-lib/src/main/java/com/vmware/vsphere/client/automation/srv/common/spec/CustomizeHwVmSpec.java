/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Container class for extending the {@link VmSpec} with additional properties
 * for Customizing the Hardware.
 */
public class CustomizeHwVmSpec extends VmSpec {

   /**
    * Property for HDD capacity/size
    */
   public DataProperty<HddSpec> hddList;

   /**
    * Property that specifies CD/DVD Drive to be added, modified, or deleted
    */
   public DataProperty<CdDvdDriveSpec> cdDvdDriveList;

   /**
    * Property that specifies Shared PCI Device to be added, modified, or
    * deleted
    */
   public DataProperty<SharedPciDeviceSpec> sharedPciDeviceList;
}
