package com.vmware.vsphere.client.automation.common.annotations;

import java.lang.reflect.Array;
import java.lang.reflect.Field;
import java.lang.reflect.ParameterizedType;
import java.util.List;
import java.util.Set;

import org.reflections.ReflectionUtils;

import com.google.common.base.Predicate;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.hsua.common.datamodel.PropertyBox;

/**
 * Initialization implementation for classes with @
 * {@link UsesSpec} annotated fields
 */
public class UsesSpecFieldsInitializer {

   /**
    * Error in initialization of @{@link UsesSpec} annotated field
    */
   public static class UsesSpecInitializationException extends
         RuntimeException {

      private static final long serialVersionUID = 1L;
      private final String initializationMessage;

      /**
       * Initializes new instance of @{@link UsesSpecInitializationException}
       *
       * @param message
       *           exception message
       * @param field
       *           the field for which the error occurs
       */
      UsesSpecInitializationException(String message, Field field) {
         this(message, field, null);
      }

      /**
       * Initializes new instance of @{@link UsesSpecInitializationException}
       *
       * @param message
       *           exception message
       * @param field
       *           the field for which the error occurs
       * @param causedBy
       *           the cause for this exception
       */
      public UsesSpecInitializationException(String message, Field f,
            Throwable causedBy) {
         super(String.format("Error initializing field %s in %s. %s",
               f.getName(), f.getDeclaringClass().getCanonicalName(), message),
               causedBy);
         this.initializationMessage = message;
      }

      public String getInitializationMessage() {
         return this.initializationMessage;
      }
   }

   /**
    * The instance for initialization
    */
   private final Object instance;

   /**
    * Initializes new instance of @{@link UsesSpecFieldsInitializer}
    *
    * @param instance
    *           the instance to be initialized
    */
   public UsesSpecFieldsInitializer(Object instance) {
      this.instance = instance;

   }

   /**
    * Performs the initialization of the fields for current instance from the
    * provided spec supplier
    *
    * @param specSupplier
    *           the spec container to supply the spec values
    */
   public void initializeFields(WorkflowSpec specSupplier) {
      for (Field field : findWorkflowAnnotatedFields()) {
         Object spec = null;

         final Class<?> fieldType = field.getType();
         if (fieldType.isArray()) {
            // For fields:
            // @WorkflowSpecAnnotation
            // private SomeSpecClass[] arrayOfSpecsVariableName;
            spec = extractArraySpecForField(specSupplier, field);
         } else if (fieldType.isAssignableFrom(List.class)) {
            // For fields:
            // @WorkflowSpecAnnotation
            // private List<SomeSpecClass> listOfSpecsVariableName;
            spec = extractListOfSpecsForField(specSupplier, field);
         } else {
            // For fields:
            // @WorkflowSpecAnnotation
            // private SomeSpecClass variableName;
            spec = extractRawSpecForField(specSupplier, field);
         }

         if (spec == null) {
            throw new UsesSpecInitializationException(
                  "No spec built for field.", field);
         }

         setField(field, this.instance, spec);
      }
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
   private static void setField(Field field, Object fieldHolder, Object value) {
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

   /**
    * Extracts from a given WorkflowSpec the provided spec for a given field
    *
    * @param specSupplier
    * @param field
    * @return
    */
   @SuppressWarnings("unchecked")
   private PropertyBox extractRawSpecForField(WorkflowSpec specSupplier,
         Field field) {
      Class<?> fieldType = field.getType();
      if (!PropertyBox.class.isAssignableFrom(fieldType)) {
         throw new UsesSpecInitializationException(
               "The field should inherit from PropertyBox or should be Array or List<> of classes inherited from property box.",
               field);
      }

      List<? extends PropertyBox> spec = specSupplier.getAll((Class<? extends PropertyBox>) field.getType());

      // The filed should be populated only if there is only one spec from the corresponding type
      if (spec.size() == 1) {
         return spec.get(0);
      } else if (spec.isEmpty()) {
         throw new UsesSpecInitializationException(
            "There isn't spec with the corresponding type.", field);
      } else {
         throw new UsesSpecInitializationException(
            "Several specs were found for field. Declare the field as array or List<T> (super of List<T>) to get all of them or filter the workflow to get single.",
            field);
      }
   }

   /**
    * Extracts from a given WorkflowSpec List of SomeSpecClass instances
    * provided for a given Field
    *
    * @param specSupplier
    * @param field
    * @return
    */
   @SuppressWarnings("unchecked")
   private List<? extends PropertyBox> extractListOfSpecsForField(
         WorkflowSpec specSupplier, Field field) {
      ParameterizedType listType = (ParameterizedType) field.getGenericType();

      Class<?> listElementsType = (Class<?>) listType.getActualTypeArguments()[0];

      if (!PropertyBox.class.isAssignableFrom(listElementsType)) {
         throw new UsesSpecInitializationException(
               "The type of the list elements should inherit from PropertyBox",
               field);
      }

      return specSupplier
            .getAll((Class<? extends PropertyBox>) listElementsType);
   }

   /**
    * Extracts array of specs for a field
    *
    * @param specContainer
    * @param field
    * @return
    */
   private Object extractArraySpecForField(WorkflowSpec specContainer,
         Field field) {
      Class<?> arrayElementsType = field.getType().getComponentType();

      if (!PropertyBox.class.isAssignableFrom(arrayElementsType)) {
         throw new UsesSpecInitializationException(
               "The type of the array elements should inherit PropertyBox",
               field);
      }

      @SuppressWarnings("unchecked")
      Class<? extends PropertyBox> arrayElementsTypeAsPropertyBox = (Class<? extends PropertyBox>) field
            .getType().getComponentType();

      Object result = toArrayOfType(field.getType().getComponentType(),
            specContainer.getAll(arrayElementsTypeAsPropertyBox));

      return result;
   }

   /**
    * Extracts all fields with any access modifier in the instance class
    * hierarchy which are annotated with @ {@link UsesSpec}
    * annotation
    *
    * @return
    */
   @SuppressWarnings("unchecked")
   private Set<Field> findWorkflowAnnotatedFields() {
      return ReflectionUtils.getAllFields(instance.getClass(),
            new Predicate<Field>() {
               @Override
               public boolean apply(Field input) {
                  return input
                        .isAnnotationPresent(UsesSpec.class);
               }
            });
   }

   /**
    * Converts List<T> to T[]
    *
    * @param elementType
    *           the Class intance for T
    * @param list
    *           the list to be converted
    * @return
    */
   @SuppressWarnings("unchecked")
   private static <L> L[] toArrayOfType(Class<L> elementType, List<?> list) {
      Object array = Array.newInstance(elementType, 1);
      return list.toArray((L[]) array);

   }

}
