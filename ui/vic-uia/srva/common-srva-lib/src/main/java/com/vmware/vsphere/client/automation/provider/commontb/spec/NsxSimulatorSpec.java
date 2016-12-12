/** Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

public class NsxSimulatorSpec extends EntitySpec {
   /**
    * Hostname or IP address of the host
    */
   public DataProperty<String> hostName;

   /**
    * Username for the host
    */
   public DataProperty<String> userName;

   /**
    * The password corresponding to the username
    */
   public DataProperty<String> password;

   /**
    * Build number of the nsx simulator
    */
   public DataProperty<String> buildNumber;
}
