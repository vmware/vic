/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.servicespec;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Define Host service spec.
 */
public class HostServiceSpec extends ServiceSpec {

   /**
    * Parameter that should be set to true if it is host client connection
    */
   public DataProperty<Boolean> isHostClient;

   /**
    * isHostClient should be True
    */
   public HostServiceSpec() {
      isHostClient.set(Boolean.TRUE);
   }

}
