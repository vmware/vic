/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for Virtual Center folder properties. Properties necessary
 * for folder creation are included.
 *
 */
public class FolderSpec extends ManagedEntitySpec {
   /**
    * Type of the folder - datacenter, vm, host, network, storage
    */
   public DataProperty<FolderType> type;
}
