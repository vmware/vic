package com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage;

import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.common.annotations.MockedSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.ProviderManifest;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.ProvidersManifestInitializer;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.ProvidersManifestInitializer.UsesProviderAnnotationException;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.UsesProvider;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.DerivedManifestAndFieldsMultipleProviders;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.FieldsMultipleProviders;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.FieldsSingleProvider;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestAndFields;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestAndFieldsMultipleProviders;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestDuplicatingId;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestDuplicatingIdsWhenDeriving;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestDuplicatingIdsWithField;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestWithMultipleProviders;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.AnnotatedClasses.ManifestWithSingleProvider;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProvider;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProvider.MockProviderSpec;
import com.vmware.vsphere.client.automation.storage.lib.core.annotations.providerusage.mock.MockProviderTestBedBridge;

public class UsesProviderInitializerTests {
   private WorkflowSpec workflowSpec;
   private MockProviderTestBedBridge testbedBridge;

   /**
    * Counts the amount of {@link MockProviderSpec} in a {@link WorkflowSpec}
    * with a particular entityId (see {@link MockProviderSpec#entityId})
    *
    * @param workflowSpec
    * @param entityId
    * @return
    */
   private static void assertSpecsCount(WorkflowSpec workflowSpec,
         String entityId, int expectedCount) {
      int counter = 0;
      for (MockProviderSpec spec : workflowSpec.getAll(MockProviderSpec.class)) {
         if (spec.entityId.equals(entityId)) {
            counter++;
         }
      }

      Assert.assertEquals(String.format(
            "%s specs with id %s found in workflow", expectedCount, entityId),
            expectedCount, counter);
   }

   @Before
   public void initWokrflowSpec() {
      workflowSpec = new WorkflowSpec();
      testbedBridge = new MockProviderTestBedBridge();
   }

   /**
    * Does not affects classes with no annotations
    */
   @Test
   public void noAnnotationsUsed() {
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            new Object());

      initializer.initSpec(this.workflowSpec, this.testbedBridge);

      // No specs pushed to workflow
      Assert.assertTrue("No specs go to workflow",
            this.workflowSpec.getAll(MockedSpec.class).isEmpty());

      // No providers requests made
      Assert.assertEquals("Provider instances used to initialize", 0,
            testbedBridge.getTotalRequestsForProvider());
   }

   /**
    * Can provide specs from a single provider instance only using
    * {@link ProviderManifest} See {@link ProvidersManifestInitializer}
    */
   @Test
   public void manifestOnlySingleProvider() {
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            new ManifestWithSingleProvider());

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // Single provider call
      Assert.assertEquals("Provider instances used to initialize", 1,
            testbedBridge.getTotalRequestsForProvider());

      // Single spec in the workflow
      Assert.assertEquals("Amount of specs pushed to workflow", 1,
            this.workflowSpec.getAll(MockProviderSpec.class).size());

      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 1);
   }

   /**
    * Can request entities from more than one provider instances using the
    * {@link ProviderManifest}
    */
   @Test
   public void manifestOnlyMultipleProviders() {
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            new ManifestWithMultipleProviders());

      initializer.initSpec(this.workflowSpec, this.testbedBridge);

      // Requested 2 providers
      Assert.assertEquals("Provider instances used to initialize", 2,
            testbedBridge.getTotalRequestsForProvider());

      // Pushed 2 specs to workflow
      Assert.assertEquals("Total amount of specs pushed to workflow", 2,
            this.workflowSpec.getAll(MockProviderSpec.class).size());

      // Entity_1 specs pushed to workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 2);
   }

   /**
    * Descriptive exception is throw when user is forcing single provider
    * instance but is annotating with the same entity key twice
    */
   @Test(expected = UsesProviderAnnotationException.class)
   public void manifestDuplicatingIds() {
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            new ManifestDuplicatingId());

      try {
         initializer.initSpec(this.workflowSpec, this.testbedBridge);
      } catch (UsesProviderAnnotationException e) {
         Assert.assertEquals("Correct exception message",
               "Dublicate usage of entity key mock.entity.1",
               e.shortDescription);
         throw e;
      }
   }

   /**
    * Descriptive exception is throw when user is forcing single provider
    * instance but is annotating with the same entity key twice
    */
   @Test(expected = UsesProviderAnnotationException.class)
   public void manifestDuplicatingIdsInDerivedClass() {
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            new ManifestDuplicatingIdsWhenDeriving());

      try {
         initializer.initSpec(this.workflowSpec, this.testbedBridge);
      } catch (UsesProviderAnnotationException e) {
         Assert.assertEquals("Correct exception message",
               "Dublicate usage of entity key mock.entity.1",
               e.shortDescription);
         throw e;
      }
   }

   /**
    * Descriptive exception is throw when user is forcing single provider
    * instance but is annotating with the same entity key twice
    */
   @Test(expected = UsesProviderAnnotationException.class)
   public void manifestDuplicatingIdsInField() {
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            new ManifestDuplicatingIdsWithField());

      try {
         initializer.initSpec(this.workflowSpec, this.testbedBridge);
      } catch (UsesProviderAnnotationException e) {
         Assert.assertEquals("Correct exception message",
               "Dublicate usage of entity key mock.entity.1",
               e.shortDescription);
         throw e;
      }
   }

   /**
    * Can use {@link UsesProvider} annotation on both - fields and
    * {@link ProviderManifest}
    */
   @Test
   public void manifestAndFieldsSingleProviderInstance() {
      ManifestAndFields annotatedClassInstance = new ManifestAndFields();
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            annotatedClassInstance);

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // Single request for provider
      Assert.assertEquals("Provider instances used to initialize", 1,
            testbedBridge.getTotalRequestsForProvider());

      // The class fields are assigned
      Assert.assertEquals("The field is assigned correctly",
            MockProvider.ENTITY_3,
            annotatedClassInstance.getProvider1Entity3().entityId);

      // Total specs pushed to workflow
      Assert.assertEquals("Amount of entities pushed to workflow", 3,
            this.workflowSpec.getAll(MockProviderSpec.class).size());

      // Entities by entity keys pushed in workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 1);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_2, 1);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_3, 1);
   }

   /**
    * Can request entities from multiple provider instance when using both
    * fields and {@link ProviderManifest}
    */
   @Test
   public void manifestAndFieldsMultipleProviders() {
      ManifestAndFieldsMultipleProviders annotatedClassInstance = new ManifestAndFieldsMultipleProviders();
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            annotatedClassInstance);

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // Single request for provider
      Assert.assertEquals("Provider instances used to initialize", 3,
            testbedBridge.getTotalRequestsForProvider());

      // The class fields are assigned
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider2Entity1().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider3Entity1().entityId);

      // Total specs pushed to workflow
      Assert.assertEquals("Amount of entities pushed to workflow", 5,
            this.workflowSpec.getAll(MockProviderSpec.class).size());

      // Entities by entity keys pushed in workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 3);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_2, 2);
   }

   /**
    * The {@link UsesProvider} can be applied to a base type
    */
   @Test
   public void baseClassAnnotationsOnly() {
      ManifestAndFieldsMultipleProviders annotatedClassInstance = new ManifestAndFieldsMultipleProviders() {
      };
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            annotatedClassInstance);

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // 3 request for provider
      Assert.assertEquals("Provider instances used to initialize", 3,
            testbedBridge.getTotalRequestsForProvider());

      // Total specs pushed to workflow
      Assert.assertEquals("Amount of entities pushed to workflow", 5,
            this.workflowSpec.getAll(MockProviderSpec.class).size());

      // The class fields are assigned
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider2Entity1().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider3Entity1().entityId);

      // Entities by entity keys pushed in workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 3);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_2, 2);
   }

   /**
    * The {@link UsesProvider} can be applied on base class and on deriving
    * class
    */
   @Test
   public void baseClassAndDerivedClassAnnotations() {
      DerivedManifestAndFieldsMultipleProviders annotatedClassInstance = new DerivedManifestAndFieldsMultipleProviders();
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            annotatedClassInstance);

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // 4 request for provider
      Assert.assertEquals("Provider instances used to initialize", 4,
            testbedBridge.getTotalRequestsForProvider());

      // Total specs pushed to workflow
      Assert.assertEquals("Amount of entities pushed to workflow", 8,
            this.workflowSpec.getAll(MockProviderSpec.class).size());

      // Base class fields hold specs for the correct entity id
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider2Entity1().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider3Entity1().entityId);

      // Derived class fields hold specs for the correct entity id
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider4Entity1().entityId);

      // Entities by entity keys pushed in workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 4);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_2, 2);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_3, 2);
   }

   /**
    * The {@link UsesProvider} can be applied to class fields
    */
   @Test
   public void fieldsAnnotationsSingleProvider() {
      FieldsSingleProvider annotatedClassInstance = new FieldsSingleProvider();
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            annotatedClassInstance);

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // Single provider instance requested
      Assert.assertEquals("Provider instances used to initialize", 1,
            testbedBridge.getTotalRequestsForProvider());

      // The fields are assigned correctly
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider1Entity1().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_2,
            annotatedClassInstance.getProvider1Entity2().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_3,
            annotatedClassInstance.getProvider1Entity3().entityId);

      // Entities by entity keys pushed in workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 1);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_2, 1);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_3, 1);

   }

   /**
    * Multiple provider instances can be requested via the {@link UsesProvider}
    * annotation applied to a different fields
    */
   @Test
   public void fieldsAnnotationsMultipleProviders() {
      FieldsMultipleProviders annotatedClassInstance = new FieldsMultipleProviders();
      ProvidersManifestInitializer initializer = new ProvidersManifestInitializer(
            annotatedClassInstance);

      initializer.initSpec(this.workflowSpec, this.testbedBridge);
      // 2 provider instances requested
      Assert.assertEquals("Provider instances used to initialize", 2,
            testbedBridge.getTotalRequestsForProvider());

      // The fields are assigned correctly
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_1,
            annotatedClassInstance.getProvider1Entity1().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_2,
            annotatedClassInstance.getProvider1Entity2().entityId);
      Assert.assertEquals("Spec for field is correct", MockProvider.ENTITY_2,
            annotatedClassInstance.getProvider1Entity3().entityId);

      // Entities by entity keys pushed in workflow
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_1, 1);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_2, 2);
      assertSpecsCount(this.workflowSpec, MockProvider.ENTITY_3, 0);

   }

}
