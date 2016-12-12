/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.servicespec;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Define VC service spec.
 *
 */
public class VcServiceSpec extends ServiceSpec {

   /**
    * Parameter that should be set to false if it is vc connection
    */
   public DataProperty<Boolean> isHostClient;

   /**
    * isHostClient should be False
    */
   public VcServiceSpec() {
      isHostClient.set(Boolean.FALSE);
   }

}
