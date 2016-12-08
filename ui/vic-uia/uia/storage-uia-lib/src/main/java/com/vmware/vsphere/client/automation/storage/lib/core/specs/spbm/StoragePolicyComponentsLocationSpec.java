package com.vmware.vsphere.client.automation.storage.lib.core.specs.spbm;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.NGCLocationSpec;

public class StoragePolicyComponentsLocationSpec extends NGCLocationSpec {

   /**
    * Build a location that will navigate the UI to the storage policy
    * components view.
    */
   public StoragePolicyComponentsLocationSpec() {
      super(NGCNavigator.NID_HOME_RULES_AND_PROFILES,
            NGCNavigator.NID_VM_STORAGE_POLICIES, null,
            NGCNavigator.NID_STORAGE_POLICY_STORAGE_POLICY_COMPONENTS);
   }
}
