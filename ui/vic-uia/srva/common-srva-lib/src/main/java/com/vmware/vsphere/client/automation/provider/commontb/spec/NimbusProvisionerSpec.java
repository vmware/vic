/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec that is describes provisioner over Nimbus infrastructure.
 * The properties - Nimbus pod, vm name and owner dines a Nimbus vm.
 */
public class NimbusProvisionerSpec extends EntitySpec {
   /**
    * Nimbus pod
    */
   public DataProperty<String> pod;

   /**
    * Nimbus VM name
    */
   public DataProperty<String> vmName;
   public DataProperty<String> vmOwner;

}
