package com.vmware.vsphere.client.automation.storage.lib.core.specs.spbm;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

public class StoragePolicyErrorMessageSpec extends EntitySpec {

   /**
    * Error message
    */
   public DataProperty<String> expectedMessage;
}
