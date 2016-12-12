/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec that defines Nimbus service.
 */
public class NimbusServiceSpec extends ServiceSpec {

   /**
    * Nimbus pod to be used for deploying operation.
    * If not nto sure what yo set use AUTO
    */
   public DataProperty<String> pod;

   /**
    * User to impersonate to as for nimbus connection is used mts-automaiton user.
    */
   public DataProperty<String> deployUser;
}
