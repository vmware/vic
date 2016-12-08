package com.vmware.hsua.common.datamodel;

import java.lang.annotation.Annotation;
import java.lang.reflect.Field;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.hsua.common.datamodel.SparseGraph.LinkType;

/**
 * This class is intended to be a base for classes that contain
 * {@link AbstractProperty} fields i.e. Property Boxes. The function of the
 * property box is similar to that of the java bean classes: they contain
 * structured data, they can be populated and red.<br><br>
 * Besides the functionality of a bean class, the property boxes have the
 * following advantages:
 * <li> Simple definition - a property box class does not have to have setter
 * and getter methods, neither a constructor - just a list of property
 * definitions.
 * <li> Self awareness - the {@link AbstractProperty} fields which are used in
 * the property box classes have links to the property box instance that
 * contains them and to the {@link Field} instance that holds them.
 * <li> Easy PBox linking - there is an easy way to define cross-linked property
 * fields. Cross-linked properties allows each of two PBox classes to have
 * property that refer the other one. After declaring the cross-linked couple -
 * setting a value in one of the properties automatically sets the other one
 * with the corresponding value.
 * <li> Meta data enriched data fields - the access to the {@link Field} holding
 * the Property allows access to the {@link Annotation}s of the field. This
 * allows meta-data to be added to each property field through annotations.
 * <li> Smart setting and usage of data fields - property field's meta data is
 * used to add smart automatic functions at properties setting or usage.
 * <br><br>
 * Check the unit test class NGCDemoDataModelTest for usage examples.
 * <br><br>
 * @author dkozhuharov
 */
public abstract class PropertyBox {

   // logger
   private static final Logger _logger = LoggerFactory.getLogger(PropertyBox.class);

   /**
    * This map allows retrieval of fields for each annotation types used.
    * The key is annotation type and the value is a list of all properties been
    * annotated with it.
    */
   private final Map<Class<?>, List<Field>> annotation2Field =
         new HashMap<Class<?>, List<Field>>();

   /**
    * This constructors instantiates automatically all {@link AbstractProperty} fields
    * declared in the child classes.
    */
   protected PropertyBox() {
      // Loop the Property fields of this class and initialize them
      for (Field field : getPropertyFields()) {
         try {
            // Get the annotations of the current Property field
            Annotation[] annotationsList = field.getAnnotations();

            // Create new Property instance
            AbstractProperty<?> fieldInst =
                  (AbstractProperty<?>) field.getType().newInstance();
            // Assign the new property instance to the Property field
            field.set(this, fieldInst);
            // Initialize the property field
            fieldInst.init(this, field, annotationsList);

            // Fill in the annotation-to-Property map
            for (Annotation annotation : annotationsList) {
               List<Field> list;
               if(annotation2Field.containsKey(annotation.annotationType())) {
                  list = annotation2Field.get(annotation.annotationType());
               } else {
                  list = new ArrayList<Field>();
                  annotation2Field.put(annotation.annotationType(), list);
               }
               list.add(field);
            }

         } catch (Throwable e) {
            throw new RuntimeException(
                  "Initialization failed for property-box field <"
                        + this.getClass().getSimpleName() + "."
                        + field.getName() + ">. Fields must be 'public'.", e);
         }
      }
   }

   /**
    * This method retrieves a list of all fields of type {@link AbstractProperty}
    * @return list of fields of type {@link AbstractProperty}
    */
   public Iterable<Field> getPropertyFields() {
      HashSet<Field> propFields = new HashSet<Field>();
      Class<?> cls = this.getClass();

      while (cls != null && PropertyBox.class.isAssignableFrom(cls)) {
         for (Field field : cls.getDeclaredFields()) {
            if (AbstractProperty.class.isAssignableFrom(field.getType())) {
               propFields.add(field);
            }
         }
         cls = cls.getSuperclass();
      }

      return propFields;
   }


   /**
    * This method copies enumerated list of properties between two instances
    * of a given {@link PropertyBox} child class. Besides the values - the
    * copy operation copies also the <b>isAssigned</b> state of the properties.
    * @param sourceProperties - {@link AbstractProperty} instances to be copied in the
    * corresponding properties of the current {@link PropertyBox} instance. The
    * copying does not invoke the set-transformation logic.
    */
   public void copy(AbstractProperty<?>...sourceProperties) {
      for (AbstractProperty<?> property : sourceProperties) {
         copyProperty(property.holderField, property);
      }
   }

   /**
    * This method copies all of the properties between two instances
    * of a given {@link PropertyBox} child class. Besides the values - the
    * copy operation copies also the <b>isAssigned</b> state of the properties.
    * @param sourcePropertyBox - {@link PropertyBox} child class instance to be
    * used as a property source.
    */
   public void copy(PropertyBox sourcePropertyBox) {
      if (sourcePropertyBox == null) {
         return;
      }
      for (Field propField : getPropertyFields()) {
         try {
            copyProperty(propField, (AbstractProperty<?>) propField.get(sourcePropertyBox));
         } catch (Exception e) {
            _logger.error(
                  "Fail to retrieve property <{}> from class <{}>",
                  propField.getName(),
                  sourcePropertyBox.getClass().getSimpleName());
            throw new RuntimeException("Unable to copy propertyBox fields!");
         }
      }
   }

   @SuppressWarnings("unchecked")
   private <T> void copyProperty(Field propField, AbstractProperty<T> sourceProperty) {
      try {
         AbstractProperty<T> localProperty =
               (AbstractProperty<T>) sourceProperty.holderField.get(this);
         localProperty.set(sourceProperty);
      } catch (Exception e) {
         _logger.warn(
               "Unable to copy failed for property <{}> into class <{}>",
               sourceProperty.getPropertyName(),
               this.getClass().getSimpleName());
         _logger.warn("Will try to match the field again!");

         //TODO: once property boxes are migrated to simple structures this workaround
         // should be removed
         // if the field is different from the destination one
         // try to copy it's value regardless of the fact that the source does not match
         // destination type
         // this workaround is needed in case of getUnassignedValue() use
         try {
            AbstractProperty<T> localProperty =
                  (AbstractProperty<T>) propField.get(this);
            localProperty.set(sourceProperty);
         } catch (IllegalAccessException ex) {
            _logger.error(
                  "Could not copy failed for property <{}> into class <{}>",
                  sourceProperty.getPropertyName(),
                  this.getClass().getSimpleName());
            throw new RuntimeException("Unable to copy property box field!");
         }
      }
   }

   /**
    * Return the first property marked with annotation {@code annotationClass}.
    * @param  annotationClass   annotation type.
    * @return Property being annotated.
    */
   public AbstractProperty<?> getAnnotated(Class<?> annotationClass) {
      List<AbstractProperty<?>> properties = getAllAnnotated(annotationClass);
      return (properties == null ? null : properties.get(0));
   }

   /**
    * Return all properties marked with annotation {@code annotationClass}.
    * @param annotationClass
    * @return List of properties being annotated.
    */
   public List<AbstractProperty<?>> getAllAnnotated(Class<?> annotationClass) {
      List<Field> fields = annotation2Field.get(annotationClass);
      if(fields == null) {
         return null;
      }
      List<AbstractProperty<?>> properties = new ArrayList<AbstractProperty<?>>();
      for (Field field : fields) {
         try {
            properties.add((AbstractProperty<?>)field.get(this));
         } catch (Exception e) {
            _logger.warn("Problem extracting field: " + field.getName() + " Due to: " + e);
         }
      }

      return properties;
   }

   // ======================================================================
   // Linked Property Box accessors
   // ======================================================================

   /**
    * This field is the access point to a link management sub-tool. It allows
    * adding, retrieving and removing of links to other {@link PropertyBox}
    * instances.
    * @see PropertyBoxLinks
    */
   public final PropertyBoxLinks links = new PropertyBoxLinks(this);

   /**
    * This method is a convenience wrapper for the link management sub-tool
    * accessible through the field {@link #links}.<br>
    * The method allows access to one {@link PropertyBox} instance of class
    * given by the parameter <b>linkedPBoxClass</b> which is linked with
    * a link of type {@link LinkType#OUTGOING}.
    * <br><br>
    * <b>NOTE:</b>
    * <br><br>
    * @param <L> - the type of the {@link PropertyBox} instance to be returned
    * @param linkedPBoxClass - a {@link Class} instance representing the
    * type of the {@link PropertyBox} instance to be returned
    * @return the matching {@link PropertyBox} instance or <b>null</b> if no
    * matching instance was found.
    */
   public final <L extends PropertyBox> L get(
         Class<L> linkedPBoxClass) {
      return links.get(linkedPBoxClass, LinkType.values());
   }
   /**
    * This method is a convenience wrapper for the link management sub-tool
    * accessible through the field {@link #links}.<br>
    * The method allows retrieval of all {@link PropertyBox} instances of class
    * given by the parameter <b>linkedPBoxClass</b> which are linked with
    * a link of type {@link LinkType#OUTGOING}.
    * <br><br>
    * @param <L> - the type of the {@link PropertyBox} instances to be returned
    * @param linkedPBoxClass - a {@link Class} instance representing the
    * type of the {@link PropertyBox} instances to be returned
    * @return the list of matching {@link PropertyBox} instances or empty list
    * if no matching instance was found.
    */
   public final <L extends PropertyBox> List<L> getAll(
         Class<L> linkedPBoxClass) {
      return links.getAll(linkedPBoxClass, LinkType.values());
   }

   /**
    * This method is a convenience wrapper for the link management sub-tool
    * accessible through the field {@link #links}.<br>
    * This method adds a link connecting the <b>Managed PBox</b> to one
    * or more {@link PropertyBox} instances given as enumerated list. The
    * type of the links is the default type {@link LinkType#OUTGOING}.
    * <br>
    * @param linkedPBoxes - var arg that could take one or more
    * {@link PropertyBox} instances for linking.
    */
   public void add(PropertyBox...linkedPBoxes) {
      links.add(linkedPBoxes);
   }
}
