/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for opaque network properties.
 * The parent of opaque network should be datacenter.
 */
public class OpaqueNetworkSpec extends ManagedEntitySpec {

   /**
    * Opaque network type
    */
   public DataProperty<String> type;

}
