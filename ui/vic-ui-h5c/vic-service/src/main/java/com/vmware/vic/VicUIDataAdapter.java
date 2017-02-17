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

*/
package com.vmware.vic;

import java.util.List;
import java.util.Set;
import java.net.URI;
import java.util.ArrayList;
import java.util.HashSet;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.vmware.vic.model.ModelObject;
import com.vmware.vim25.InvalidPropertyFaultMsg;
import com.vmware.vim25.RuntimeFaultFaultMsg;
import com.vmware.vise.data.Constraint;
import com.vmware.vise.data.PropertySpec;
import com.vmware.vise.data.ResourceSpec;
import com.vmware.vise.data.query.Comparator;
import com.vmware.vise.data.query.CompositeConstraint;
import com.vmware.vise.data.query.DataProviderAdapter;
import com.vmware.vise.data.query.ObjectIdentityConstraint;
import com.vmware.vise.data.query.PropertyConstraint;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.QuerySpec;
import com.vmware.vise.data.query.RelationalConstraint;
import com.vmware.vise.data.query.RequestSpec;
import com.vmware.vise.data.query.Response;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.data.query.type;
import com.vmware.vise.vim.data.VimObjectReferenceService;
import com.vmware.vise.data.query.ResultSet;

@type("vic:Root,vic:VirtualContainerHostVm,vic:ContainerVm,vic:VmQueryResult")
public class VicUIDataAdapter implements DataProviderAdapter {
	private static final Log _logger = LogFactory.getLog(VicUIDataAdapter.class);
	public static final String ROOT_TYPE = "vic:Root";
	public static final String VCH_TYPE = "vic:VirtualContainerHostVm";
	public static final String CONTAINER_TYPE = "vic:ContainerVm";
	public static final String VMQUERY_RESULT_TYPE = "vic:VmQueryResult";

	private final VimObjectReferenceService _objRefService;
	private final ObjectStore _objectStore;

	public VicUIDataAdapter(
			VimObjectReferenceService objRefService,
			ObjectStore objectStore
			) {
		if (objRefService == null ||
            objectStore == null) {
			throw new IllegalArgumentException("Constructor arguments cannot be null");
		}
		_objRefService = objRefService;
		_objectStore = objectStore;
	}

	/**
	 * Extends vSphere Client's DataService
	 */
	@Override
	public Response getData(RequestSpec request) throws Exception {
		if (request == null) {
			throw new IllegalArgumentException("request should be non-null");
		}

		QuerySpec[] querySpecs = request.querySpec;
		List<ResultSet> results = new ArrayList<ResultSet>(querySpecs.length);

		for (QuerySpec querySpec : querySpecs) {
			ResultSet resultset = processQuerySpec(querySpec);
			results.add(resultset);
		}

		Response response = new Response();
		response.resultSet = results.toArray(new ResultSet[] {});
		return response;
	}

	/**
	 * Process QuerySpec to get requested data
	 * @param querySpec
	 * @return ResultSet containing requested object(s)
	 */
	private ResultSet processQuerySpec(QuerySpec querySpec) {
		ResultSet rs = new ResultSet();
		if (!validateQuery(querySpec)) {
			return rs;
		}

		List<ResultItem> items = processConstraint(
				querySpec.resourceSpec.constraint,
				querySpec.resourceSpec.propertySpecs);
		rs.queryName = querySpec.name;
		rs.totalMatchedObjectCount = (items != null) ? items.size() : 0;
		rs.items = items.toArray(new ResultItem[items.size()]); // TODO: consider pagination

		return rs;
	}

	/**
	 * Processes various types of Constraint for consumption with ResultSet
	 * @param constraint
	 * @return ResultItem list containing query results
	 */
	private List<ResultItem> processConstraint(
			Constraint constraint,
			PropertySpec[] propertySpecs) {
		List<ResultItem> results = null;
		if (constraint instanceof ObjectIdentityConstraint) {
			ObjectIdentityConstraint oic = (ObjectIdentityConstraint)constraint;
			results = processObjectIdentityConstraint(oic, propertySpecs);
		} else if (constraint instanceof CompositeConstraint) {
			_logger.warn("CompositeConstraint is unsupported");
		} else if (constraint instanceof PropertyConstraint) {
			_logger.warn("PropertyConstraint is unsupported");
		} else if (constraint instanceof RelationalConstraint) {
			_logger.warn("RelationalConstraint is unsupported");
		} else if (isSimpleConstraint(constraint)) {
			results = processSimpleConstraint(constraint, propertySpecs);
		}

		if (results == null) {
			results = new ArrayList<ResultItem>();
		}
		return results;
	}

	private boolean isSimpleConstraint(Constraint constraint) {
		return constraint.getClass().getSimpleName().equals(
				Constraint.class.getSimpleName());
	}

	/**
	 * Process an ObjectIdentityConstraint where constraint.target is a
	 * specific object for which we need to return requested properties
	 * @throws RuntimeFaultFaultMsg
	 * @throws InvalidPropertyFaultMsg
	 */
	private List<ResultItem> processObjectIdentityConstraint(
			ObjectIdentityConstraint oic,
			PropertySpec[] propertySpecs) {
		List<ResultItem> items = new ArrayList<ResultItem>();
		URI objectUri = toURI(oic.target);
		if (objectUri != null) {
			ModelObject mo = _objectStore.getObj(objectUri);
			if (mo == null) {
				_logger.error("ModelObject for " + objectUri + " does not exist!");
				return items;
			}
			ResultItem ri = createResultItem(mo, objectUri, propertySpecs);
			if (ri != null) {
				items.add(ri);
			}
		}
		return items;
	}

	private List<ResultItem> processSimpleConstraint(
			Constraint constraint,
			PropertySpec[] propertySpecs) {
		List<ResultItem> items = new ArrayList<ResultItem>();

		// simple constraint is used in collection view and collection view
		// only needs to show vic:Root
		if (ROOT_TYPE.equals(constraint.targetType)) {
			URI objectUri = _objectStore.getRootUri();
			ModelObject mo = _objectStore.getObj(objectUri);
			if (mo == null) {
				_logger.error("ModelObject for " + objectUri + " does not exist!");
				return items;
			}
			ResultItem ri = createResultItem(mo, objectUri, propertySpecs);
			if (ri != null) {
				items.add(ri);
			}
		}
		return items;
	}

	/**
	 * Extract requested properties from the given ModelObject
	 * @param mo
	 * @param uri
	 * @param propertySpecs
	 * @return ResultItem object containing PropertyValues for the given
	 *         ModelObject and propertySpecs
	 */
	private ResultItem createResultItem(
			ModelObject mo,
			URI uri,
			PropertySpec[] propertySpecs) {
		String[] requestedPropertyNames = getPropertyNames(propertySpecs);
		try {
			if (mo == null) {
				throw new IllegalArgumentException(
						"ModelObject not found for " + uri.toString());
			}
			ResultItem ri = new ResultItem();
			ri.resourceObject = uri;
			List<PropertyValue> pvs = new ArrayList<PropertyValue>(
					propertySpecs.length);
			for (String reqPropertyName : requestedPropertyNames) {
				Object value = mo.getProperty(reqPropertyName);
				if (value != ModelObject.UNSUPPORTED_PROPERTY) {
					PropertyValue pv = new PropertyValue();
					pv.propertyName = reqPropertyName;
					pv.resourceObject = uri;
					pv.value = value;
					pvs.add(pv);
				}
			}
			ri.properties = pvs.toArray(new PropertyValue[]{});
			return ri;
		} catch (Exception ex) {
			_logger.error("Error getting the ResultItem for " + uri, ex);
			return null;
		}
	}

	/**
	 * Extract property names from PropertySpec[]
	 * @param propertySpecs
	 * @return String[] containing property names
	 */
	private String[] getPropertyNames(PropertySpec[] propertySpecs) {
		Set<String> properties = new HashSet<String>();
		if (propertySpecs != null) {
			for (PropertySpec pSpec : propertySpecs) {
				for (String pName : pSpec.propertyNames) {
					properties.add(pName);
				}
			}
		}
		return properties.toArray(new String[]{});
	}

	/**
	 * Cast Object of type URI into URI object
	 * @param object
	 * @return (URI) object
	 */
	private URI toURI(Object object) {
		if (!(object instanceof URI)) {
			return null;
		}
		return (URI)object;
	}

	/**
	 * Verify QuerySpec
	 * @param qs
	 * @return true if valid, otherwise false
	 */
	private boolean validateQuery(QuerySpec qs) {
		if (qs == null) {
			return false;
		}

		ResourceSpec resourceSpec = qs.resourceSpec;
		if (resourceSpec == null) {
			return false;
		}

		return validateConstraint(resourceSpec.constraint);
	}

	/**
	 * Validate resourceSpec's constraint
	 * @param constraint
	 * @return true if valid, otherwise false
	 */
	private boolean validateConstraint(Constraint constraint) {
		if (constraint == null) {
			return false;
		}

		if (constraint instanceof ObjectIdentityConstraint) {
			String sourceType = _objRefService.getResourceObjectType(
					((ObjectIdentityConstraint)constraint).target);
			return isSupportedType(sourceType);
		} else if (constraint instanceof CompositeConstraint) {
	         CompositeConstraint cc = (CompositeConstraint)constraint;
	         for (Constraint c : cc.nestedConstraints) {
	            if (!validateConstraint(c)) {
	               return false;
	            }
	         }
	         return true;
		} else if (constraint instanceof PropertyConstraint) {
			return isSupportedType(constraint.targetType) &&
					((PropertyConstraint)constraint).comparator
					.equals(Comparator.TEXTUALLY_MATCHES);
		} else if (isSimpleConstraint(constraint)) {
			return isSupportedType(constraint.targetType);
		}

		_logger.error("querySpec constraint is not supported: " +
				constraint.getClass().getName());
		return false;
	}

	/**
	 * Validate targetType of Constraint object
	 * @param type
	 * @return true if supported, otherwise false
	 */
	private boolean isSupportedType(String type) {
		return VCH_TYPE.equals(type) ||
				CONTAINER_TYPE.equals(type) ||
				ROOT_TYPE.equals(type) ||
				VMQUERY_RESULT_TYPE.equals(type);
	}
}
