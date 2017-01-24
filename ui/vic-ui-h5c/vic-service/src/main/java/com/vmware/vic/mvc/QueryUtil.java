/*

Copyright 2017 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Mac OS script starting an Ant build of the current flex project
Note: if Ant runs out of memory try defining ANT_OPTS=-Xmx512M

*/

package com.vmware.vic.mvc;

import java.lang.reflect.Array;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import com.vmware.vise.data.Constraint;
import com.vmware.vise.data.PropertySpec;
import com.vmware.vise.data.ResourceSpec;
import com.vmware.vise.data.query.CompositeConstraint;
import com.vmware.vise.data.query.Conjoiner;
import com.vmware.vise.data.query.DataService;
import com.vmware.vise.data.query.ObjectIdentityConstraint;
import com.vmware.vise.data.query.ObjectReferenceService;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.QuerySpec;
import com.vmware.vise.data.query.RelationalConstraint;
import com.vmware.vise.data.query.RequestSpec;
import com.vmware.vise.data.query.Response;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.data.query.ResultSet;

/**
 * General Query utility class for the DataService
 *
 */
public class QueryUtil {

   private static ObjectReferenceService _objectReferenceService;

   public static void setObjectReferenceService(
            ObjectReferenceService objectReferenceService) {
      _objectReferenceService = objectReferenceService;
   }

   /**
    * Helper method to get one or more properties for a given object.
    *
    * @param dataService
    *    The data service instance that will be used for retrieving
    *    the property value.
    *
    * @param obj
    *    Object whose properties are required.
    *
    * @param properties
    *    The names of the object properties you need to retrieve.
    *
    * @return An array of <code>PropertyValue</code> instances. Each instance should
    *    be the value of one of the properties from the <code>properties</code> array.
    *    The results however can be shuffled: The first
    *    <code>PropertyValue</code> from the result does not necessarily correspond
    *    to the first property name in <code>properties</code> array.
    *
    * @see #getProperties(DataService, Object[], String[])
    */
   public static PropertyValue[] getProperties(
         DataService dataService,
         Object obj,
         String[] properties) throws Exception {
      return getProperties(dataService, new Object[] {obj}, properties);
   }

   /**
    * Helper method to get one or more properties for a set of given objects.
    *
    * Note: In case if only some of the requested properties are retrieved
    * (for some of the objects), instead of throwing error, the partial result
    * will be returned.
    *
    * @param dataService
    *    DataService instance to use for retrieving the properties of the object.
    *
    * @param objs
    *    An array of object instances whose properties are required.
    *
    * @param properties
    *    Name of the properties to be requested.
    *
    * @return
    *    Values of the requested properties. The array is flat: this means that
    *    the property values for the different objects will be mixed in it, in no
    *    particular order.
    *
    * @throws Exception
    *    Throws <code>Exception</code> if invalid query parameters are passed.
    *    For e.g. if <code>objs</code> or <code>properties</code>
    *    is null or empty.
    *
    *    Throws <code>Exception</code> if an empty result is retrieved
    *    for the query.
    */
   public static PropertyValue[] getProperties(
         DataService dataService,
         Object[] objs,
         String[] properties) throws Exception {
      if (objs == null || objs.length == 0 ||
            properties == null || properties.length == 0) {
         throw new Exception("Invalid parameters for getProperties");
      }

      Object obj = objs[0];

      QuerySpec query = buildQuerySpec(objs, properties);

      query.name = _objectReferenceService.getUid(obj) + ".properties";

      ResultSet resultSet = getData(dataService, query);

      ArrayList<PropertyValue> result = new ArrayList<PropertyValue>();
      if (resultSet != null) {
         ResultItem[] items = resultSet.items;
         if (items != null && items.length > 0 && items[0] != null) {
            for (ResultItem item : items) {
               for (PropertyValue v : item.properties) {
                  v.resourceObject = item.resourceObject;
                  result.add(v);
               }
            }
         }
      }

      // Check if we have at least the partial result, if so
      // then we return the result else we throw
      if (result.isEmpty() && resultSet.error != null) {
         throw resultSet.error;
      }

      // Return the result. It will be empty array if that property does not exist.
      return toArray(result, PropertyValue.class);
   }


   @SuppressWarnings({ "unchecked" })
   public static <T> T[] toArray(Collection<T> collection, Class<T> elementType) {
      T[] result = null;

      if (collection == null) {
         result = (T[]) Array.newInstance(elementType, 0);
         return result;
      }

      T[] copy = (T[]) Array.newInstance(elementType, collection.size());
      result = collection.toArray(copy);
      return result;
   }


   /**
    * A shortcut method to <code>getPropertiesForRelatedObjects</code> for the case
    * where just one property is requested.
    *
    * <p> Note that the value returned is still an array of properties, cause the
    * <code>relationship</code> may resolve to more than one related object.
    *
    * @see #getPropertiesForRelatedObjects(DataService, Object, String, String, String[])
    */
   public static PropertyValue[] getPropertyForRelatedObjects(
         DataService dataService,
         Object object,
         String relationship,
         String targetType,
         String property) throws Exception {
      return getPropertiesForRelatedObjects(
            dataService, object, relationship, targetType, new String[] {property});
   }


   /**
    * Return the requested <code>properties</code> on the objects which are in
    * <code>relationship</code> with the original <code>object</code>.
    *
    * <p> The following example will demonstrate how to retrieve all datastore
    * names of the datastores where the VM is located:
    * {@code
    * PropertyValue[] names = QueryUtil.getPropertiesForRelatedObjects(
    *    dataService, vm, "datastore", "Datastore", new String[] {"name"});
    * if (names != null) {
    *    for (PropetyValue propValue : propValues) {
    *       System.out.println(propValue.value);
    *    }
    * }
    * }
    *
    * @param dataService
    *   The data service to query about the properties.
    *
    * @param object
    *   The root object for the relationship.
    *
    * @param relationship
    *   The relationship string. A string referencing a property of the
    *   <code>object</code> parameter.
    *
    * @param targetType
    *   The type of the object that the <code>relationship</code> will return.
    *   In general you can put <code>null</code> here if you are unsure what objects
    *   are returned, but if you know the type, please do specify it. Passing
    *   <code>null</code> may not always work.
    *
    * @param properties
    *   An array of string referencing properties on the objects returned when the
    *   <code>relationship</code> is resolved.
    *
    * @return
    *   An array of <code>PropetyValue</code> instances or null if the
    *   <code>relationship</code> cannot be resolved properly or the properties
    *   requested are missing.
    *
    * @throws Exception
    */
   public static PropertyValue[] getPropertiesForRelatedObjects(
            DataService dataService,
            Object obj,
            String relationship,
            String targetType,
            String[] properties) throws Exception {
      if (obj == null || properties == null || properties.length == 0) {
         throw new Exception("invalid parameters in getPropertiesForRelatedObjects");
      }

      // If no relation is given, consider the "properties" as being requested
      // directly on "obj".
      if (relationship == null || relationship.length() == 0) {
         return getProperties(dataService, obj, properties);
      }

      // create object identity constraint to match the given server object
      ObjectIdentityConstraint objectConstraint =
            createObjectIdentityConstraint(obj);

      // create relational constraint for the given relationship to the source object
      RelationalConstraint relationalConstraint =
            createRelationalConstraint(relationship,
                                       objectConstraint,
                                       true, // true is important here.
                                       targetType);

      QuerySpec query = buildQuerySpec(relationalConstraint, properties);
      query.name = _objectReferenceService.getUid(obj) + "." + relationship + ".properties";

      ResultSet resultSet = getData(dataService, query);

      // Check if we have at least one result.
      ArrayList<PropertyValue> result = new ArrayList<PropertyValue>();
      if (resultSet != null && resultSet.items != null) {
         for (ResultItem item : resultSet.items) {
            if (item != null && item.properties != null) {
               for (PropertyValue propValue : item.properties) {
                  if (propValue != null) {
                     result.add(propValue);
                  }
               }
            }
         }
      }

      if (result.isEmpty() && resultSet.error != null) {
         throw resultSet.error;
      }

      // There was no error reported, so assume that property does not exist
      // or the value is null
      return toArray(result, PropertyValue.class);
   }

   /**
    * Helper method to call DataService with single QuerySpec. Uses injected
    * DataService
    *
    * @param dataService
    *           The data service instance.
    * @param query
    *           QuerySpec for the data-service
    * @return ResultSet from the data-service
    *
    * @throws Exception
    *    Throws <code>Exception</code> if an empty result is retrieved
    *    for the query.
    */
   public static ResultSet getData(
         DataService dataService,
         QuerySpec query) throws Exception {

      RequestSpec requestSpec = new RequestSpec();
      requestSpec.querySpec = new QuerySpec[] { query };

      Response response = new Response();
      response  = dataService.getData(requestSpec);

      ResultSet[] retVal = response.resultSet;
      if (retVal == null || retVal.length == 0 || retVal[0] == null) {
         throw new Exception("Empty result");
      }
      return retVal[0];
   }

   /**
    * Helper utility for building simplest data-service QuerySpec: one object
    * and several properties of this object.
    *
    * @param entity
    *           - managed object of interest
    * @param properties
    *           - names of properties
    * @return - QuerySpec to feed into DataService
    */
   public static QuerySpec buildQuerySpec(Object entity,
         String[] properties) {
      ObjectIdentityConstraint oc = new ObjectIdentityConstraint();
      oc.target = entity;

      String targetType = _objectReferenceService.getResourceObjectType(entity);
      Set<String> targetTypes = new HashSet<String>();
      targetTypes.add(targetType);

      QuerySpec query = buildQuerySpec(oc, properties, targetTypes);
      return query;
   }

   public static QuerySpec buildQuerySpec(
         Object[] entities,
         String[] properties) {
      if (entities.length == 1) {
         return buildQuerySpec(entities[0], properties);
      }
      CompositeConstraint cc = new CompositeConstraint();
      cc.conjoiner = Conjoiner.OR;
      Constraint[] nestedConstraints = new Constraint[entities.length];
      Set<String> targetTypes = new HashSet<String>();
      String targetType = null;

      for (int index=0; index < entities.length; index++) {
         ObjectIdentityConstraint oc = new ObjectIdentityConstraint();
         oc.target = entities[index];
         nestedConstraints[index] = oc;

         targetType = _objectReferenceService.getResourceObjectType(oc.target);
         targetTypes.add(targetType);
      }
      cc.nestedConstraints = nestedConstraints;

      QuerySpec query = buildQuerySpec(cc, properties, targetTypes);
      return query;
   }

   /**
    * Helper utility for building simplest QuerySpec that requests a few properties
    * on the objects that are identified with the <code>constraint</code>.
    *
    * @param constraint
    *   A constraint which will define the set of object being considered for
    *   property retrieval.
    *
    * @param properties
    *   The properties we want to retrieve.
    *
    * @return
    *   A <code>QuerySpec</code> which you can directly feed into a getData()
    *   method of the <code>DataService</code>.
    */
   public static QuerySpec buildQuerySpec(Constraint constraint,
                                          String[] properties) {
      QuerySpec query =  buildQuerySpec(
            constraint,
            properties,
            null /* target types not known */);
      return query;
   }

   /**
    * Helper utility for building simplest QuerySpec that requests a few properties
    * on the objects that are identified with the <code>constraint</code>.
    *
    * @param constraint
    *   A constraint which will define the set of object being considered for
    *   property retrieval.
    *
    * @param properties
    *   The properties we want to retrieve.
    *
    * @param targetTypes
    *   For each target type in the set, a <code>PropertySpec</code> is created.
    *   If targetTypes is null, then a single <code>PropertySpec</code> is created
    *   with type unset.
    *
    * @return
    *   A <code>QuerySpec</code> which you can directly feed into a getData()
    *   method of the <code>DataService</code>.
    */
   public static QuerySpec buildQuerySpec(
         Constraint constraint,
         String[] properties,
         Set<String> targetTypes) {
      QuerySpec query = new QuerySpec();
      ResourceSpec resourceSpec = new ResourceSpec();
      resourceSpec.constraint = constraint;

      List<PropertySpec> pSpecs = new ArrayList<PropertySpec>();
      if (targetTypes != null) {
         for (String targetType : targetTypes) {
            PropertySpec propSpec = createPropertySpec(properties, targetType);
            pSpecs.add(propSpec);
         }
      } else {
         PropertySpec propSpec = createPropertySpec(properties, null);
         pSpecs.add(propSpec);
      }

      resourceSpec.propertySpecs = pSpecs.toArray(new PropertySpec[]{});
      query.resourceSpec = resourceSpec;

      return query;
   }


   /**
    * Creates a RelationalConstraint
    *
    * @param relationship
    *    relationship to traverse
    *
    * @param constraintOnRelatedObject
    *    constraint on the objects related to the targeted objects by the
    *    given relationship
    *
    * @param hasInverseRelation
    *    indicates whether the constraint given by constraintOnRelatedObject applies to the
    *    source, as opposed to the target of the relationship given in this instance
    *
    * @param targetType
    *    type of the objects to retrieve
    *
    * @return The created RelationalConstraint.
    */
   public static RelationalConstraint createRelationalConstraint(
         String relationship,
         Constraint constraintOnRelatedObject,
         Boolean hasInverseRelation,
         String targetType) {
      RelationalConstraint rc = new RelationalConstraint();
      rc.relation = relationship;
      rc.hasInverseRelation = hasInverseRelation;
      rc.constraintOnRelatedObject = constraintOnRelatedObject;
      rc.targetType = targetType;
      return rc;
   }


   /**
    * Creates an ObjectIdentityConstraint
    *
    * @param entity
    *    Object to be looked up
    *
    * @return The created ObjectIdentityConstraint.
    */
   public static ObjectIdentityConstraint createObjectIdentityConstraint(
         Object entity) {
      ObjectIdentityConstraint oc = new ObjectIdentityConstraint();
      oc.target = entity;
      oc.targetType = _objectReferenceService.getResourceObjectType(entity);
      return oc;
   }

   /**
    * Utility to create a property spec.
    */
   private static PropertySpec createPropertySpec(
         String[] properties, String targetType) {
      PropertySpec propSpec = new PropertySpec();
      propSpec.type = targetType;
      propSpec.propertyNames = properties;
      return propSpec;
   }

}
