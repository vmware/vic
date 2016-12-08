package com.vmware.vsphere.client.automation.storage.lib.core.annotations;

import java.lang.reflect.Field;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.reflections.ReflectionUtils;

import com.google.common.base.Predicate;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpecFieldsInitializer.UsesSpecInitializationException;
import com.vmware.vsphere.client.automation.storage.lib.core.tests.ProviderPublication.ProviderPublicationUsage;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpecFieldsInitializer;

/**
 * Implementation of class fields and workflow initialization for
 * {@link UsesProvider} and {@link ProviderManifest} annotations
 */
public class ProvidersManifestInitializer {

   /**
    * {@link RuntimeException} implementation for errors caused by the
    * {@link UsesProvider} and/or {@link ProviderManifest} errors
    */
   public class UsesProviderAnnotationException extends RuntimeException {

      private static final long serialVersionUID = 1L;

      public final String shortDescription;

      public UsesProviderAnnotationException(UsesProvider annotation,
            Class<?> clazz, String shortDescription) {
         super(String.format(
               "%s in type %s.\nAnnotation value: %s\nFull class name: %s",
               shortDescription, clazz.getSimpleName(), annotation,
               clazz.getCanonicalName()));
         this.shortDescription = shortDescription;
      }
   }

   /**
    * Data model of request for {@link ProviderWorkflow} provided entity
    */
   private static final class ProviderRequest {

      /**
       * The class which implements the {@link ProviderWorkflow}
       */
      private final Class<? extends ProviderWorkflow> providerClass;

      /**
       * Map of Entity keys for {@link ProviderWorkflow} entities and fields to
       * which should ne assigned
       */
      private final Map<String, Field> entityRequests;

      /**
       * The usage of the provided entities
       */
      public ProviderPublicationUsage usage = ProviderPublicationUsage.SHARED;

      /**
       * Initializes new instance of {@link ProviderRequest}
       *
       * @param providerClass
       *           the provider class
       */
      ProviderRequest(Class<? extends ProviderWorkflow> providerClass) {
         this.providerClass = providerClass;
         this.entityRequests = new HashMap<>();
      }

      /**
       * Checks if entity with a specific key is already requested
       *
       * @param entityId
       * @return
       */
      boolean isRequested(String entityId) {
         return entityRequests.containsKey(entityId);
      }

      /**
       * Adds new entity to the entity request
       *
       * @param entityId
       *           the providers entity id
       * @param assigningField
       *           the field to be assigned
       */
      void addEntityId(String entityId, Field assigningField,
            ProviderPublicationUsage usage) {
         this.entityRequests.put(entityId, assigningField);
         if (usage == ProviderPublicationUsage.EXCLUSIVE) {
            // If we have any of the entities requested with exclusive lock the
            // whole provider instance should be requested with exclusive lock
            usage = ProviderPublicationUsage.EXCLUSIVE;
         }
      }
   }

   /**
    * Instance to be initialized
    */
   private final Object instance;

   /**
    * Provider requests
    */
   private Map<String, ProviderRequest> providerRequests;

   /**
    * Initializes new instance of {@link ProvidersManifestInitializer}
    *
    * @param instance
    *           the object to be initialized
    */
   public ProvidersManifestInitializer(Object instance) {
      this.instance = instance;
   }

   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      initializeProviderRequests();
      handleProviderRequests(testSpec, testbedBridge);
   }

   /**
    * Initializes the provider request for current {@link #instance} to be
    * initialized
    */
   private void initializeProviderRequests() {
      this.providerRequests = new HashMap<>();

      List<ProviderManifest> typeHierarchyProviderManifests = extractManifestAnnotationFromTypeHierarchy();

      for (ProviderManifest providerManifest : typeHierarchyProviderManifests) {
         for (UsesProvider usesProviderAnnotation : providerManifest.value()) {
            addProviderRequest(usesProviderAnnotation, null);
         }
      }

      for (Field field : findUseProviderAnnotatedFields()) {
         addProviderRequest(field.getAnnotation(UsesProvider.class), field);
      }
   }

   /**
    * Extracts all the fields in a type hierarchy which have the
    * {@link ProviderManifest} annotation
    *
    * @return
    */
   @SuppressWarnings("unchecked")
   private List<ProviderManifest> extractManifestAnnotationFromTypeHierarchy() {
      List<ProviderManifest> result = new ArrayList<ProviderManifest>();
      Set<Class<?>> typesWithManifestAnnotation = ReflectionUtils
            .getAllSuperTypes(this.instance.getClass(),
                  new Predicate<Class<?>>() {
                     @Override
                     public boolean apply(Class<?> input) {
                        return input
                              .isAnnotationPresent(ProviderManifest.class);
                     }
                  });

      for (Class<?> manifestType : typesWithManifestAnnotation) {
         result.add(manifestType.getAnnotation(ProviderManifest.class));
      }

      return result;
   }

   /**
    * Handles the provider requests for current {@link #instance}
    *
    * @param testSpec
    * @param testbedBridge
    */
   private void handleProviderRequests(WorkflowSpec testSpec,
         TestBedBridge testbedBridge) {
      for (ProviderRequest request : providerRequests.values()) {
         final TestbedSpecConsumer testBedSpecConsumer = testbedBridge
               .requestTestbed(request.providerClass, request.usage.isShared);
         for (String entityKey : request.entityRequests.keySet()) {
            EntitySpec publishedEntity = testBedSpecConsumer
                  .getPublishedEntitySpec(entityKey);

            if (publishedEntity == null) {
               throw new RuntimeException(String.format(
                     "No entity spec registerd from provider %s with key %s",
                     request.providerClass.getCanonicalName(), entityKey));
            }

            Field field = request.entityRequests.get(entityKey);

            if (field != null) {
               assignField(publishedEntity, field);
            }

            testSpec.add(publishedEntity);
         }
      }
   }

   /**
    * Adds a provider request
    *
    * @param annotation
    * @param field
    */
   private void addProviderRequest(UsesProvider annotation, Field field) {
      if (!providerRequests.containsKey(annotation.id())) {
         providerRequests.put(annotation.id(),
               new ProviderRequest(annotation.clazz()));
      }

      ProviderRequest request = providerRequests.get(annotation.id());
      if (!request.providerClass.equals(annotation.clazz())) {
         throw new UsesProviderAnnotationException(
               annotation,
               this.instance.getClass(),
               String.format(
                     "Dubious provider class. The id %s is used with provider %s but is requested for %s",
                     annotation.id(), request.providerClass.getName(),
                     annotation.clazz().getName()));
      }

      if (request.isRequested(annotation.entity())) {
         throw new UsesProviderAnnotationException(annotation,
               this.instance.getClass(), String.format(
                     "Dublicate usage of entity key %s", annotation.entity()));
      }

      request.addEntityId(annotation.entity(), field, annotation.usage());
   }

   /**
    * Extracts all fields with any access modifier in the instance class
    * hierarchy which are annotated with @ {@link WorkflowSpecAnnotation}
    * annotation
    *
    * @return
    */
   @SuppressWarnings("unchecked")
   private Set<Field> findUseProviderAnnotatedFields() {
      return ReflectionUtils.getAllFields(instance.getClass(),
            new Predicate<Field>() {
               @Override
               public boolean apply(Field input) {
                  return input.isAnnotationPresent(UsesProvider.class);
               }
            });
   }

   /**
    * Assigns {@link EntitySpec} field of {@link #instance}
    *
    * @param publishedEntity
    *           the value to be assigned
    * @param field
    *           the field to be assigned
    */
   private void assignField(EntitySpec publishedEntity, Field field) {
      Class<?> fieldType = field.getType();
      if (!publishedEntity.getClass().isAssignableFrom(fieldType)) {
         throw new RuntimeException(
               String.format(
                     "The spec found for field %s is of type %s and can not be assigned to field of type %s",
                     field.getName(), publishedEntity.getClass()
                           .getCanonicalName(), field.getType().getName()));
      }

      setFieldIgnoreAccessModifier(field, this.instance, publishedEntity);
   }

   /**
    * Sets the value of a field for instnace
    *
    * @param field
    *           the field of the instance to be set
    * @param fieldHolder
    *           the instance to be used for the set
    * @param value
    *           the required value
    * @throws IllegalArgumentException
    * @throws IllegalAccessException
    */
   private static void setFieldIgnoreAccessModifier(Field field,
         Object fieldHolder, Object value) {
      boolean isAccessible = field.isAccessible();
      try {
         field.setAccessible(true);
         field.set(fieldHolder, value);
      } catch (IllegalArgumentException | IllegalAccessException e) {
         throw new UsesSpecInitializationException(
               "Exception thrown while initializing property.", field, e);
      } finally {
         field.setAccessible(isAccessible);
      }
   }

}
