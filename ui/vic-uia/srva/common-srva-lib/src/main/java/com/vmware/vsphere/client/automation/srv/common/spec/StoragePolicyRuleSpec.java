/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.pbm.capability.CapabilityInstance;
import com.vmware.vim.binding.pbm.capability.CapabilityMetadata;
import com.vmware.vim.binding.pbm.capability.ConstraintInstance;
import com.vmware.vim.binding.pbm.capability.PropertyInstance;

/**
 * A spec for a storage policy rule ({@link CapabilityInstance}). A rule belongs
 * to a certain namespace and has multiple properties ({@link PropertyInstance}
 * ). Rules are contained in both components ({@link StoragePolicyComponentSpec}
 * ) and policies ({@link StoragePolicySpec}).
 */
public class StoragePolicyRuleSpec extends BaseSpec {
   public static enum ProviderNamespace {
      Replication("com.vmware.vim.xVP.replication"), VmWareVmCrypt(
            "vmwarevmcrypt"), Compression("com.vmware.vim.xVP.compression"), SPM(
            "spm"), Caching("com.vmware.vim.xVP.caching"), Encryption(
            "com.vmware.vim.xVP.encryption"), Persistence2(
            "com.vmware.vim.xVP.persistence2"), Persistence1(
            "com.vmware.vim.xVP.persistence1"), Tag(
            "http://www.vmware.com/storage/tag");
      public final String ns;

      ProviderNamespace(String ns) {
         this.ns = ns;
      }
   }

   /**
    * Id of the storage policy rule (corresponds to
    * {@link CapabilityInstance#id}).
    */
   public DataProperty<String> id;
   /**
    * Label of the storage policy rule. If a storage policy rule contains a
    * single {@link PropertyInstance}, only this label is shown.
    */
   public DataProperty<String> label;
   /**
    * Namespace this rule belongs to.
    */
   public DataProperty<ProviderNamespace> namespace;
   /**
    * Properties belonging to this rule.
    */
   public DataProperty<List<PropertyInstance>> properties;

   public StoragePolicyRuleSpec() {
   }

   public StoragePolicyRuleSpec(String id, String label) {
      this.id.set(id);
      this.label.set(label);
   }

   public StoragePolicyRuleSpec(String id, String label, ProviderNamespace ns) {
      this(id, label);
      this.namespace.set(ns);
   }

   @SuppressWarnings("unchecked")
   public StoragePolicyRuleSpec(String id, String label, ProviderNamespace ns,
         PropertyInstance property) {
      this(id, label, ns);
      final List<PropertyInstance> properties = new ArrayList<>();
      properties.add(property);
      this.properties.set(properties);
   }

   @SuppressWarnings("unchecked")
   public StoragePolicyRuleSpec(String id, String label, ProviderNamespace ns,
         List<PropertyInstance> properties) {
      this(id, label, ns);
      this.properties.set(properties);
   }

   /**
    * Applies the values of the custom fields of the input spec (source) to this
    * spec (target). Custom fields are only those fields, which are explicitly
    * defined in the class, i.e. without the inherited ones.
    *
    * @param ruleSpec
    *           source spec, whose field values to be applied onto this (target)
    *           spec
    */
   @SuppressWarnings("unchecked")
   public void apply(StoragePolicyRuleSpec ruleSpec) {
      this.id.set(ruleSpec.id.get());
      this.label.set(ruleSpec.label.get());
      this.namespace.set(ruleSpec.namespace.get());
      this.properties.set(ruleSpec.properties.get());
   }

   public List<PropertyInstance> getProperties() {
      return properties.get();
   }

   /**
    * Transforms the current spec into a {@link CapabilityInstance} instance.
    *
    * @return a capability instance with the same namespace, id and set of
    *         properties (values)
    */
   public CapabilityInstance asCapabilityInstance() {
      CapabilityInstance capabilityInstance = new CapabilityInstance();
      capabilityInstance.setId(new CapabilityMetadata.UniqueId(
            namespace.get().ns, id.get()));
      ConstraintInstance constraintInstance = new ConstraintInstance(properties
            .get().toArray(new PropertyInstance[properties.get().size()]));
      capabilityInstance
            .setConstraint(new ConstraintInstance[] { constraintInstance });
      return capabilityInstance;
   }

   @Override
   public String toString() {
      return String.format("{id: %s, ns: %s, label: %s}", id, namespace, label);
   }
}
