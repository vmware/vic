/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.pbm.profile.CapabilityBasedProfile.ProfileCategoryEnum;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec for storage policies.
 */
public class StoragePolicySpec extends ManagedEntitySpec {

   /**
    * Description of the storage policy.
    */
   public DataProperty<String> description;

   /**
    * Profile category of the storage policy.
    */
   public DataProperty<ProfileCategoryEnum> profileCategory;

   /**
    * Rule-sets collection contained in the storage policy.
    */
   public DataProperty<StoragePolicyRuleSetSpec> rulesets;
}