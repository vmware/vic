/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for distributed virtual portgroup properties. The parent for
 * DVPortgroup should be DVS
 */
public class DvPortgroupSpec extends ManagedEntitySpec {
   /**
    * Active uplinks that should be bind to this portgroup. By default if active
    * uplinks are not set, all uplinks are marked as active
    */
   public DataProperty<String[]> activeUplinks;
}
