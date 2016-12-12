package com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage;

import com.vmware.vsphere.client.automation.storage.lib.core.annotations.ProviderManifest;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.UsesProvider;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProvider;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProvider.MockProviderSpec;

public class AnnotatedClasses {

   /**
    * Class with allowed usage of @{@link ProviderManifest} annotation
    */
   @ProviderManifest(@UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class))
   public static class ManifestWithSingleProvider {
   }

   /**
    * Class with allowed usage of @{@link ProviderManifest} annotation
    */
   @ProviderManifest({
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class),
         @UsesProvider(id = "mock.provider.instance.2", entity = MockProvider.ENTITY_1, clazz = MockProvider.class) })
   public static class ManifestWithMultipleProviders {

   }

   /**
    * Class with disallowed usage of @{@link ProviderManifest} annotation. The
    * issue is that we say that we want twice the same provider entity from the
    * same provider instance
    */
   @ProviderManifest({
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class),
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class) })
   public static class ManifestDuplicatingId {
   }

   /**
    * Class with disallowed usage of @{@link ProviderManifest} similar to
    * {@link ManifestDuplicatingId} but the duplication is done when deriving
    */
   @ProviderManifest(@UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class))
   public static class ManifestDuplicatingIdsWhenDeriving extends
         ManifestWithSingleProvider {
   }

   /**
    * Class with disallowed usage of @{@link ProviderManifest} similar to
    * {@link ManifestDuplicatingId} but the duplication is done with mixed
    * {@link ProviderManifest} on class level and {@link UsesProvider} on field
    */
   public static class ManifestDuplicatingIdsWithField extends
         ManifestWithSingleProvider {
      @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class)
      private MockProviderSpec provider1Entity1;
   }

   /**
    * Class with allowed usage of @{@link ProviderManifest} annotation
    */
   @ProviderManifest({
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class),
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_2, clazz = MockProvider.class) })
   public static class ManifestAndFields {

      @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_3, clazz = MockProvider.class)
      private MockProviderSpec provider1Entity3;

      public MockProviderSpec getProvider1Entity3() {
         return provider1Entity3;
      }
   }

   /**
    * Class with allowed usage of @{@link ProviderManifest} and @
    * {@link UsesProvider} annotations
    */
   @ProviderManifest({
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class),
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_2, clazz = MockProvider.class),
         @UsesProvider(id = "mock.provider.instance.2", entity = MockProvider.ENTITY_2, clazz = MockProvider.class) })
   public static class ManifestAndFieldsMultipleProviders {

      @UsesProvider(id = "mock.provider.instance.2", entity = MockProvider.ENTITY_1, clazz = MockProvider.class)
      private MockProviderSpec provider2Entity1;

      @UsesProvider(id = "mock.provider.instance.3", entity = MockProvider.ENTITY_1, clazz = MockProvider.class)
      private MockProviderSpec provider3Entity1;

      public MockProviderSpec getProvider2Entity1() {
         return provider2Entity1;
      }

      public MockProviderSpec getProvider3Entity1() {
         return provider3Entity1;
      }
   }

   /**
    * Class with allowed usage of @{@link UsesProvider} annotation
    */
   @ProviderManifest({
         @UsesProvider(id = "mock.provider.instance.1", entity = MockProvider.ENTITY_3, clazz = MockProvider.class),
         @UsesProvider(id = "mock.provider.instance.2", entity = MockProvider.ENTITY_3, clazz = MockProvider.class) })
   public static class DerivedManifestAndFieldsMultipleProviders extends
         ManifestAndFieldsMultipleProviders {

      @UsesProvider(id = "mock.provider.instance.4", entity = MockProvider.ENTITY_1, clazz = MockProvider.class)
      private MockProviderSpec provider4Entity1;

      public MockProviderSpec getProvider4Entity1() {
         return provider4Entity1;
      }
   }

   /**
    * Class with allowed usage of @{@link UsesProvider} annotation
    */
   public static class FieldsSingleProvider {

      @UsesProvider(id = "mock.provider.usage.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class)
      private MockProviderSpec provider1Entity1;

      @UsesProvider(id = "mock.provider.usage.1", entity = MockProvider.ENTITY_2, clazz = MockProvider.class)
      public MockProviderSpec provider1Entity2;

      @UsesProvider(id = "mock.provider.usage.1", entity = MockProvider.ENTITY_3, clazz = MockProvider.class)
      protected MockProviderSpec provider1Entity3;

      public MockProviderSpec getProvider1Entity1() {
         return provider1Entity1;
      }

      public MockProviderSpec getProvider1Entity2() {
         return provider1Entity2;
      }

      public MockProviderSpec getProvider1Entity3() {
         return provider1Entity3;
      }
   }

   /**
    * Class with allowed usage of @{@link UsesProvider} annotation
    */
   public static class FieldsMultipleProviders {

      @UsesProvider(id = "mock.provider.usage.1", entity = MockProvider.ENTITY_1, clazz = MockProvider.class)
      private MockProviderSpec provider1Entity1;

      @UsesProvider(id = "mock.provider.usage.1", entity = MockProvider.ENTITY_2, clazz = MockProvider.class)
      public MockProviderSpec provider1Entity2;

      @UsesProvider(id = "mock.provider.usage.2", entity = MockProvider.ENTITY_2, clazz = MockProvider.class)
      protected MockProviderSpec provider1Entity3;

      public MockProviderSpec getProvider1Entity1() {
         return provider1Entity1;
      }

      public MockProviderSpec getProvider1Entity2() {
         return provider1Entity2;
      }

      public MockProviderSpec getProvider1Entity3() {
         return provider1Entity3;
      }
   }
}
