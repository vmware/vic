package com.vmware.vsphere.client.automation.storage.lib.core.tests;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.SerializationUtils;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;

/**
 * Representation of a published by provider entity/entities. This class is
 * Immutable
 */
public class ProviderPublication implements ISpecInitializer, Serializable {

   private static final long serialVersionUID = 1L;

   /**
    * Resource usage for the provider publications
    */
   public static enum ProviderPublicationUsage {
      SHARED(true), EXCLUSIVE(false);

      public final boolean isShared;

      ProviderPublicationUsage(final boolean isShared) {
         this.isShared = isShared;
      }

   }

   /**
    * The provider class
    */
   private final Class<? extends ProviderWorkflow> providerClass;

   /**
    * Is sharing of current published entities allowed
    */
   private final ProviderPublicationUsage usage;

   /**
    * The keys with which the provider published the entities
    */
   private final List<String> publishedEntityKeys = new ArrayList<>();

   /**
    * Tags to be applied on entities
    */
   private String tagToApply;

   /**
    * Initialize new instance of ProviderPublication
    *
    * @param providerClass
    *           the provider implementation
    * @param isShared
    *           is sharing allowed for the current provider entities
    * @param tagToApply
    *           the tag to be applied for the provider published entities
    * @param publishedEntityKeys
    *           the entity key(s)
    */
   public ProviderPublication(Class<? extends ProviderWorkflow> providerClass,
         ProviderPublicationUsage usage, String... publishedEntityKeys) {
      this.providerClass = providerClass;
      this.usage = usage;
      if (publishedEntityKeys == null || publishedEntityKeys.length == 0) {
         throw new RuntimeException(
               "Illegal usage of ProviderPublication. The provider publication instance should have at least one entity key. Provided"
                     + ArrayUtils.toString(publishedEntityKeys));
      }
      this.publishedEntityKeys.addAll(Arrays.asList(publishedEntityKeys));
   }

   /**
    * Initialize new instance of ProviderPublication
    *
    * @param providerClass
    *           the provider implementation
    * @param publishedEntityKeys
    *           the entity key(s)
    */
   public ProviderPublication(Class<? extends ProviderWorkflow> providerClass,
         String... publishedEntityKeys) {
      this(providerClass, ProviderPublicationUsage.SHARED, publishedEntityKeys);
   }

   /**
    * Returns new instance of {@link ProviderPublication} with a tag
    *
    * @param tag
    * @return
    */
   public ProviderPublication withTag(String tag) {
      ProviderPublication result = SerializationUtils.clone(this);
      result.tagToApply = tag;
      return result;
   }

   /**
    * Returns new instance of {@link ProviderPublication} with additional entity
    * keys
    *
    * @param entityKeys
    * @return
    */
   public ProviderPublication withEntityKeys(String... entityKeys) {
      if (entityKeys == null || entityKeys.length == 0) {
         throw new IllegalArgumentException(
               "The entity keys can not be null or empty array");
      }
      ProviderPublication result = SerializationUtils.clone(this);
      result.publishedEntityKeys.addAll(Arrays.asList(entityKeys));

      return result;
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      final TestbedSpecConsumer testBedSpecConsumer = testbedBridge
            .requestTestbed(this.providerClass, this.usage.isShared);

      for (String publishedEntityKey : publishedEntityKeys) {
         EntitySpec publishedEntity = testBedSpecConsumer
               .getPublishedEntitySpec(publishedEntityKey);
         if (this.tagToApply != null) {
            publishedEntity.tag.set(tagToApply);
         }
         testSpec.add(publishedEntity);
      }
   }
}