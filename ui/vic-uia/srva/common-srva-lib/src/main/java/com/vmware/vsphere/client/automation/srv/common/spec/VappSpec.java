/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;


/**
 * Container class for Virtual Center vApp properties.
 *
 */
public class VappSpec extends ManagedEntitySpec {
   public DataProperty<VmSpec> vmList;
}
