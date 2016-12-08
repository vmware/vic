package com.vmware.hsua.common.datamodel;

import java.lang.annotation.Annotation;
import java.lang.annotation.Documented;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.reflect.Field;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


/**
 * This class is the basic implementation of a property container for the
 * {@link PropertyBox} classes. It will be extended in the declarations
 * of the {@link PropertyBox} child classes.
 * <br><br>
 * A property is a value container that allows context aware usage of the
 * value. A property implementation can hold one or more uniform values in a
 * list.
 * <br><br>
 * The child implementations of the {@link AbstractProperty} abstract class
 * can define special value-set transformations as implementation of the
 * {@link AbstractProperty#setTransformer(Object)} abstract function. Value-set
 * transformations could depend on the {@link Annotation}s assigned to
 * the property.
 * <br><br>
 * The Property contains a boolean flag that distinguish between assigned
 * and not-assigned states. This allows each {@link PropertyBox} container
 * to be used partially - i.e. only properties with assigned values are
 * intended for use, so others could be ignored/left default.
 *
 * @author dkozhuharov
 *
 * @param <E> - type of the value that will be contained in the property.
 */
public abstract class AbstractProperty<E> {

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(AbstractProperty.class);

   /** This field holds references to the contained values */
   protected List<E> values = new ArrayList<E>();
   /** This field holds reference to the containing object */
   protected PropertyBox holderPBox = null;
   /** This field holds references to the containing class field */
   protected Field holderField = null;
   /** This field holds the assignment state of the property */
   protected boolean isAssigned = false;

   /** This method returns reference to the containing object */
   public final PropertyBox holderPBox() {
      return holderPBox;
   }
   /** This method returns references to the containing class field */
   public final Field holderField() {
      return holderField;
   }
   /**
    * This method returns a reference to a particular {@link Annotation}
    * instance declared for the holding {@link Field} when give the
    * annotation's class.
    * @param <T> - a type parameter to ensure that the type of the passed
    * annotation class is the same as the type of the returned annotation
    * instance.
    * @param annotationClass - the class of the annotation to be retrieved.
    * @return {@link Annotation} instance declared for the holding
    * {@link Field} or <b>null</b> if no such type of annotation was declared.
    */
   public <T extends Annotation> T annotation(Class<T> annotationClass) {
      return holderField.getAnnotation(annotationClass);
   }

   /**
    * This constructor is intentionally made protected because only the
    * property box class must be able to create instances.
    */
   protected AbstractProperty() {
   }

   /**
    * This method initializes a property instance. It sets a back reference
    * from current Property instance to the {@link Field} instance that
    * holds it and to the object that contains it.
    *
    * @param holderInstance - the {@link PropertyBox} instance that
    * contains this property
    * @param holderField - the {@link Field} instance that holds this
    * property instance.
    * @param holderAnnotations - the full list of {@link Annotation}s
    * assigned to this property is provided so that child property
    * implementations can access and use them as context data.
    */
   protected final void init(PropertyBox holderPBox, Field holderField,
         Annotation[] holderAnnotations) {
      this.holderPBox = holderPBox;
      this.holderField = holderField;
   }

   /**
    * This setter method allows setting of one or more values.
    * Setting of new value(s) releases the former property content.<br>
    * Calling of this method without values puts the property in
    * not-assigned state.
    *
    * @param value - the value that will be contained in the Property
    * of the type {@link E}
    */
   public final void set(E...value) {
      set(Arrays.asList(value));
   }

   /**
    * This setter method allows setting of one or more values given in
    * a {@link Iterable} container. The provided values are stored and
    * later returned in the order provided by the parameter's iterator.<br>
    * Setting of new value(s) releases the former property content.
    *
    * @param valueList - a collection with values of type {@link E}
    */
   public final void set(Iterable<E> valueList) {
      values.clear();
      isAssigned = false;

      for (E value : valueList) {
         setValue(value);
      }
   }

   /**
    * This method allows loading of one or more {@link String} values given in
    * a vararg container. The provided values are converted using
    * a load-conversion class provided through the {@link Load} annotation.
    * The converted values are stored in the {@link AbstractProperty} in
    * the order provided in the parameter.<br>
    * Loading of new value(s) releases the former property content.<br>
    * @param value - an array with string values
    * @throws RuntimeException if {@link Load} annotation
    * was not given for the property or if the conversion class provided
    * does not have appropriate <b><tt>valueOf</tt></b> method implemented.
    */
   public final void load(String...value) {
      load(Arrays.asList(value));
   }

   /**
    * This method allows loading of one or more {@link String} values given in
    * a {@link Iterable} container. The provided values are converted using
    * a load-conversion class provided through the {@link Load} annotation.
    * The converted values are stored in the {@link AbstractProperty} in
    * the order provided by the parameter's iterator.<br>
    * Loading of new value(s) releases the former property content.<br>
    * @param valueList - a collection with string values
    * @throws RuntimeException if {@link Load} annotation
    * was not given for the property or if the conversion class provided
    * does not have appropriate <b><tt>valueOf</tt></b> method implemented.
    */
   @SuppressWarnings("unchecked")
   public final void load(Iterable<String> valueList) {
      values.clear();
      isAssigned = false;

      // Check if loading was preconfigured for this property
      Load loadAnn = this.annotation(Load.class);
      if (loadAnn == null) {
         throw new RuntimeException(
               "Loading for property <" + getPropertyName() + "> " +
               "failed. Must have '@Load' annotation.");
      }
      Class<? extends E> loadClass = (Class<? extends E>) loadAnn.value();

      for (String value : valueList) {
         try {
            setValue(PropertyBoxCommonUtils.valueOf(loadClass, value));
         } catch (Exception e) {
            throw new RuntimeException(
                  "Loading of <" + value + "> " +
                        "as type <" + loadClass.getSimpleName() + "> " +
                        "in property <" + getPropertyName() + "> " +
                        "failed. Check its '@Load' annotation.", e);
         }
      }
   }

   /**
    * This internal method implements the single value setting logic. In
    * the common case the value parameter is transformed using the
    * {@link AbstractProperty#setTransformer(Object)} implementation and
    * then is added to the value container.<br>
    * In the special case when the value is of type {@link PropertyBox}
    * the property cross-linking scenario is executed.
    *
    * @param value - a single values of type {@link E} to be added
    */
   protected void setValue(E value) {
      // Fire the extender's implementation of value transformation
      // and add the transformed value
      values.add(setTransformer(value));
      // Mark this property as populated
      isAssigned = true;
   }

   /**
    * This abstract method will be called just before a new value
    * accommodation in the property. Implementors are given the chance to
    * add some value transformations at value setting time.
    *
    * @param value - the value that is being set
    * @return - the transformed value
    */
   protected abstract E setTransformer(E value);

   /**
    * This method performs copying of the value of one Property to another
    * one without activation of the set-transformation logic.
    *
    * NOTE: If the source property is not assigned, properties will be bind
    * by reference. This type of mapping is required in provides, where
    * late property assignment is done.
    *
    * @param sourceProperty
    */
   public final void set(AbstractProperty<E> sourceProperty) {

      if (sourceProperty.isAssigned) {
         values.clear();
         values.addAll(sourceProperty.values);
         isAssigned = sourceProperty.isAssigned;
      } else {
         // if the source property is not assigned
         // then we want to map the whole reference
         // this extension is needed when there is late property binding
         // an example of such need is in provides which first claim resource
         // then empty specs are prepared, which are filled with data after
         // the resources are assigned to the specs
         try {
            Field field = this.holderPBox.getClass().getField(this.holderField.getName());
            field.set(this.holderPBox, sourceProperty);
         } catch (Exception e) {
            _logger.error("Unable to set the passed property.");
            throw new IllegalArgumentException("Passed proprty is not suitable for mapping!");
         }
      }
   }

   /**
    * This method returns if any value is assigned to the Property.
    * @return <b>true</b> if a value is assigned
    */
   public final boolean isAssigned() {
      return isAssigned;
   }
   protected void assertNotEmpty() {
      if (!isAssigned()) {
         if(holderField.getName().equals("service")) {
            _logger.error("!!! NOT ASSIGNED SERVICE SPEC - there is something wrong if the test is run using TestWorkflowController !!!");
         } else {
            throw new RuntimeException("The property '" +
                  getPropertyName() + "' is empty.");
         }
      }
   }

   /**
    * This method retrieves the first of the stored values.<br>
    * @return a value of type {@link E}
    */
   public final E get() {
      assertNotEmpty();
      // Extremely ugly temporary workaround for merging to new workflow model
      // TODO: rkovachev track it and revert that change once the tests are migrated.
      if(values.size() == 0) {
         return null;
      }
      return values.get(0);
   }

   /**
    * Returns the reference to the property object.
    * This method is specially designed to be used in cases when
    * two properties has to be linked in order post-initial initialisation
    * to take place.
    *
    * @return  a reference to the property object
    */
   public final AbstractProperty<E> getUnassignedValue() {
      if (this.isAssigned()) {
         throw new RuntimeException("It is not possible to retrieve unassigned "
               + "value in case of assigned propery.");
      }

      return this;
   }

   /**
    * This method creates a new {@link List} instance containing all stored
    * values. The order of the values is preserved.<br>
    * @return list of values of type {@link E}
    */
   public final List<E> getAll() {
      return new ArrayList<E>(values);
   }

   /**
    * This method returns the name of this property instance. It is
    * composed from the simple name of the property box class and the name
    * of the field that holds the reference to this instance.<br>
    * Ment to be used in logging.
    *
    * @return descriptive name of the property.
    */
   public final String getPropertyName() {
      String holderName = holderField.getDeclaringClass().getSimpleName();
      String propertyName = holderField.getName();
      final String classSuffix = "Info";

      if (holderName.endsWith(classSuffix)) {
         holderName = holderName.substring(
               0, holderName.length() - classSuffix.length());
      }

      return holderName + "." + propertyName;
   }

   @Override
   public String toString() {
      if (!isAssigned) {
         return getPropertyName() + ":<not_assigned>";
      } else if (values.size() == 1) {
         return getPropertyName() + ":" + get().toString();
      } else {
         return getPropertyName() + ":" + values.toString();
      }
   }

   /**
    * This annotation will be used to mark the fields that are to receive
    * values from an external property file. Each annotation accepts one
    * "value" parameter. It must be the class to which the String value
    * from the property file must be converted before it is assigned.
    * <br><br>
    * @author dkozhuharov
    */
   @Retention(RetentionPolicy.RUNTIME)
   @Documented
   public static @interface Load {
      Class<?> value();
   }

   /**
    * Returns a string representation of the stored value, by unwrapping it
    * and calling its toString() method.
    *
    * @return String representation of the stored value
    */
   public String stringValue() {
      return this.get().toString();
   }
}
