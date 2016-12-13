package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec representing the properties needed for creating xVp using Nimbus
 */
public class XVpProvisionerSpec extends NimbusProvisionerSpec {
   public DataProperty<String> ip;

   /**
    * Build info
    */
   public DataProperty<String> version;
   public DataProperty<String> replConfig;
   public DataProperty<String> url;

}
