/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec representing the properties needed for creating snapshot using Nimbus
 */
public class VmSnapshotSpec extends EntitySpec {
   /**
    * Spec of the vm on which the snapshot will be made
    */
   public DataProperty<NimbusProvisionerSpec> vmSpec;

   /**
    * Snapshot name
    */
   public DataProperty<String> snapshotName;

}
