/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.pbm.profile.SubProfileCapabilityConstraints.SubProfile;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec for storage policy rule-sets. Each rule-set
 * contains both rules ({@link StoragePolicyRuleSpec})
 * and components ({@link StoragePolicyComponentSpec}).
 * A rule-set as known to the UI is roughly equivalent to a {@link SubProfile}.
 */
public class StoragePolicyRuleSetSpec extends ManagedEntitySpec {
   public DataProperty<StoragePolicyRuleSpec> rules;
   public DataProperty<StoragePolicyComponentSpec> components;
}
