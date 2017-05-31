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
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.NavigableMap;
import java.util.Set;
import java.util.TreeMap;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.vmware.vic.model.ContainerVm;
import com.vmware.vic.model.ModelObject;
import com.vmware.vic.model.VirtualContainerHostVm;
import com.vmware.vic.model.constants.BaseVm;
import com.vmware.vic.model.constants.Container;
import com.vmware.vic.model.constants.Vch;
import com.vmware.vic.utils.VicVmComparator;
import com.vmware.vise.data.Constraint;
import com.vmware.vise.data.PropertySpec;
import com.vmware.vise.data.ResourceSpec;
import com.vmware.vise.data.query.CompositeConstraint;
import com.vmware.vise.data.query.Conjoiner;
import com.vmware.vise.data.query.DataService;
import com.vmware.vise.data.query.ObjectIdentityConstraint;
import com.vmware.vise.data.query.ObjectReferenceService;
import com.vmware.vise.data.query.OrderingCriteria;
import com.vmware.vise.data.query.OrderingPropertySpec;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.QuerySpec;
import com.vmware.vise.data.query.RelationalConstraint;
import com.vmware.vise.data.query.RequestSpec;
import com.vmware.vise.data.query.Response;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.data.query.ResultSet;
import com.vmware.vise.data.query.ResultSpec;
import com.vmware.vise.data.query.SortType;

/**
 * General Query utility class for the DataService
 *
 */
public class QueryUtil {

   private static ObjectReferenceService _objectReferenceService;
   private static final Log _logger = LogFactory.getLog(QueryUtil.class);

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

   /**
    * Implementation of DataAccessController's /list API
    * @throws Exception 
    */
   @SuppressWarnings("unchecked")
static ResultSet getListData(
         DataService dataService, String targetType, String[] requestedProperties,
         int offset, int maxResultCount, String[] sortParams,
         String[] filterParams) throws Exception {

      Constraint constraint = new Constraint();
      constraint.targetType = targetType;
      QuerySpec query = QueryUtil.buildQuerySpec(constraint, requestedProperties);
      ResultSpec rs = new ResultSpec();
      rs.offset = offset;
      rs.maxResultCount = maxResultCount;

      if (sortParams != null && sortParams.length == 2) {
         String sortField = sortParams[0];
         String sortDirection = sortParams[1];

         OrderingPropertySpec orderPropSpec = new OrderingPropertySpec();
         orderPropSpec.type = constraint.targetType;
         orderPropSpec.propertyNames = new String[] {sortField};
         orderPropSpec.orderingType = SortType.ASCENDING;
         if ("desc".equalsIgnoreCase(sortDirection)) {
            orderPropSpec.orderingType = SortType.DESCENDING;
         }
         rs.order = new OrderingCriteria();
         rs.order.orderingProperties = new OrderingPropertySpec[] { orderPropSpec };
      }
      query.resultSpec = rs;

      RequestSpec requestSpec = new RequestSpec();
      requestSpec.querySpec = new QuerySpec[]{query};

      // Execute data queries to get counts.
      Response response = dataService.getData(requestSpec);
      ResultSet[] resultSetArray = response.resultSet;

      // return an empty resultset upon when resultSet array is empty
      if (resultSetArray.length == 0) {
          ResultSet emptyResultSet = new ResultSet();
          emptyResultSet.totalMatchedObjectCount = 0;
          emptyResultSet.items = new ResultItem[] {};
          _logger.warn("No entry for response.resultSet");
          return emptyResultSet;
      }

      ResultSet resultSet = resultSetArray[0];
      resultSet.totalMatchedObjectCount = 0;
      if (resultSet.error != null) {
          throw resultSet.error;
      }

      int totalMatchedObjectCount = 0;
      if (resultSet.items.length > 0) {
          PropertyValue[] pvs = resultSet.items[0].properties;
          for (PropertyValue pv : pvs) {
              if ("match".equals(pv.propertyName)) {
                  resultSet.totalMatchedObjectCount +=
                      (int)pv.value;
              } else if ("results".equals(pv.propertyName)) {
                  // returns # of items found
                  Map<String, Object> adjustedResults = adjustItems(
                      (HashMap<String, ModelObject>) pv.value,
                      offset,
                      maxResultCount,
                      resultSet.totalMatchedObjectCount,
                      rs.order,
                      filterParams
                  );
                  pv.value = adjustedResults.get("results");
                  totalMatchedObjectCount =
                      (int)adjustedResults.get("totalMatchedObjectCount");
              }
          }

          if (filterParams != null && filterParams.length > 0) {
              resultSet.totalMatchedObjectCount = totalMatchedObjectCount;
          }
      }

      return resultSet;
   }

   /**
    * Handles pagination, sort and filter
    *
    * @param vAppVmsMap: vApp-VMs map
    * @param offset: index to include from
    * @param maxResultCount: page size
    * @param totalMatchedObjectCount: # of all items
    * @param order: OrderingCriteria instance
    * @param filterParams: filtering criteria
    * @return adjusted Map
    */
   private static Map<String, Object> adjustItems(
		   Map<String, ModelObject> vmsMap,
		   int offset,
		   int maxResultCount,
		   int totalMatchedObjectCount,
		   OrderingCriteria order,
		   String[] filterParams) {

	   // map to return
	   Map<String, Object> resultMap = new HashMap<String, Object>();

	   // NavigableMap to contain sorted TreeMap
	   NavigableMap<String, Object> results = new TreeMap<String, Object>();

	   // # of records matching the provided filters
	   // this would replace totalMatchedObjectCount ONLY if there is at least
	   // one filter active
	   int numOfAllRecordsMatchingFilters = 0;

	   // return empty results for
	   // 1. offset is smaller than 0
	   // 2. offset is greater than or equal to total # of matched objects
	   // 3. page size is 0
	   if (offset < -1 ||
		   offset >= totalMatchedObjectCount ||
		   maxResultCount == 0) {
		   resultMap.put("totalMatchedObjectCount", 0);
		   resultMap.put("results", results);
		   return resultMap;
	   }

	   // filters results based on the given criteria
	   if (filterParams != null && filterParams.length > 0) {
		   Set<String> keys = vmsMap.keySet();
		   Map<String, ModelObject> filteredVmsMap =
	           new HashMap<String, ModelObject>();

		   for (String key : keys) {
			   if (!vmsMap.containsKey(key)) {
				   continue;
			   }

			   // filters out VMs not matching the criteria
			   if (vmMatchesFilterCriteria(vmsMap.get(key), filterParams)) {
			       filteredVmsMap.put(key, vmsMap.get(key));
			       numOfAllRecordsMatchingFilters++;
			   }
		   }

		   // updates totalmatchedObjectCount and vmsMap
		   vmsMap = filteredVmsMap;
		   totalMatchedObjectCount = numOfAllRecordsMatchingFilters;
	   }

	   int endIndex = offset + maxResultCount - 1;
	   if (totalMatchedObjectCount == 0) {
		   endIndex = totalMatchedObjectCount;
	   } else if ((endIndex >= totalMatchedObjectCount) || (maxResultCount < 0)) {
		   endIndex = totalMatchedObjectCount - 1;
	   }

	   if (order != null && order.orderingProperties.length > 0) {
		   SortType sortType = order.orderingProperties[0].orderingType;
		   String orderBy = order.orderingProperties[0].propertyNames[0];
		   results = new TreeMap<String, Object>(
				   new VicVmComparator(
						   vmsMap,
						   orderBy,
						   sortType.equals(SortType.DESCENDING)
					   )
			   );
		   results.putAll(vmsMap);
	   }

	   // returns the sorted/filtered/paginated sub map
	   resultMap.put("totalMatchedObjectCount", totalMatchedObjectCount);

	   if (results.size() == 0) {
		   resultMap.put("results", results);
	   } else {
		   String[] vAppKeys = results.keySet().toArray(new String[]{});
		   resultMap.put("results",
			   results.subMap(vAppKeys[offset], true, vAppKeys[endIndex], true));
	   }

	   return resultMap;
   }

   private static boolean vmMatchesFilterCriteria(
		   ModelObject mo, String[] filterParams) {
	   for (String filterParam : filterParams) {
		   String[] params = filterParam.split("=");
		   if (!getVmPropValue(mo, params[0])
				   .contains(params[1].toLowerCase().trim())) {
			   return false;
		   }
	   }
	   return true;
   }

   private static String getVmPropValue(ModelObject mo, String property) {
	   if (mo instanceof VirtualContainerHostVm) {
		   if (BaseVm.VM_NAME.equals(property)) {
			   return ((VirtualContainerHostVm) mo).getName().toLowerCase();
		   } else if (Vch.VM_VCH_IP.equals(property)) {
			   return ((VirtualContainerHostVm) mo).getClientIp().toLowerCase();
		   } else if (BaseVm.VM_OVERALL_STATUS.equals(property)) {
			   return ((VirtualContainerHostVm) mo).getOverallStatus()
					   .toLowerCase();
		   }
		   return null;
	   } else if (mo instanceof ContainerVm) {
		   if (Container.VM_CONTAINERNAME_KEY.equals(property)) {
		       return ((ContainerVm) mo).getContainerName().toLowerCase();
		   } else if (BaseVm.Runtime.VM_POWERSTATE_BASENAME.equals(property)) {
		       return ((ContainerVm) mo).getPowerState().toLowerCase();
		   } else if (BaseVm.VM_GUESTMEMORYUSAGE.equals(property)) {
		       return Integer.toString(((ContainerVm) mo).getGuestMemoryUsage());
		   } else if (BaseVm.VM_OVERALLCPUUSAGE.equals(property)) {
		       return Integer.toString(((ContainerVm) mo).getOverallCpuUsage());
		   } else if (BaseVm.VM_COMMITTEDSTORAGE.equals(property)) {
		       return Long.toString(((ContainerVm) mo).getCommittedStorage());
		   } else if (Container.VM_PORTMAPPING_KEY.equals(property)) {
		       String pm = ((ContainerVm) mo).getPortMapping();
		       return pm != null ? pm.toLowerCase() : "";
		   } else if (Container.PARENT_NAME_KEY.equals(property)) {
               return ((ContainerVm) mo).getParentObjectName();
           } else if (BaseVm.VM_NAME.equals(property)) {
		       return ((ContainerVm) mo).getName().toLowerCase();
		   } else if (Container.VM_IMAGENAME_KEY.equals(property)) {
		       return ((ContainerVm) mo).getImageName().toLowerCase();
		   }
	   }
	   return null;
   }
}
