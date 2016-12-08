/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Container class for TIWO dialog title properties.
 */
public class DialogTiwoTitleSpec extends BaseSpec {

   /**
    * Specifies TIWO dialog title property.
    */
   public DataProperty<String> dialogTitle;
}
