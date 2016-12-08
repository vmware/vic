/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.pbm.capability.provider.LineOfServiceInfo.LineOfServiceEnum;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicyRuleSpec.ProviderNamespace;

/**
 * A spec for a storage policy component of a specific line of service. A
 * storage policy component contains multiple rules from the same namespace.
 */
public class StoragePolicyComponentSpec extends ManagedEntitySpec {
   /**
    * The <code>providerNamespace</code> should match the namespaces of all the
    * rules contained in the component. This field is only exposed here for
    * convenience. No <i>Tag</i> namespace is allowed for components.
    */
   public DataProperty<ProviderNamespace> providerNamespace;
   /**
    * Strings in English should match {@link LineOfServiceEnum}. Providers
    * publish one or more namespace(s) for a given <code>lineOfService</code> .
    */
   public DataProperty<LineOfServiceEnum> lineOfService;
   /**
    * Collection of rules.
    */
   public DataProperty<StoragePolicyRuleSpec> rules;
   public DataProperty<String> description;
}
