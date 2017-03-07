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

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.vmware.vic.model.ModelObject;
import com.vmware.vic.model.Root;
import com.vmware.vic.model.RootInfo;
import com.vmware.vic.model.VmQueryResult;
import com.vmware.vic.utils.ConfigLoader;
import com.vmware.vim25.InvalidPropertyFaultMsg;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim25.RuntimeFaultFaultMsg;
import com.vmware.vise.data.Constraint;
import com.vmware.vise.data.PropertySpec;
import com.vmware.vise.data.ResourceSpec;
import com.vmware.vise.data.query.CompositeConstraint;
import com.vmware.vise.data.query.Conjoiner;
import com.vmware.vise.data.query.DataService;
import com.vmware.vise.data.query.ObjectIdentityConstraint;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.QuerySpec;
import com.vmware.vise.data.query.RequestSpec;
import com.vmware.vise.data.query.Response;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.data.query.ResultSet;
import com.vmware.vise.data.query.ResultSpec;
import com.vmware.vise.vim.data.VimObjectReferenceService;

public class ObjectStore {
	private Root _currentRootObject;
	private final DataService _dataService;
	private final ModelObjectUriResolver _modelObjectUriResolver;
	private final ConfigLoader _configLoader;
	private final PropFetcher _propFetcher;
	private final VimObjectReferenceService _objectRefService;
	private final static String CONFIG_PROP_FILES = "configs.properties";
	private final static String UI_VERSION_CONFIG_KEY = "uiVersion";
	private final static String ROOT_TYPE = VicUIDataAdapter.ROOT_TYPE;
	private final static String VCH_TYPE = VicUIDataAdapter.VCH_TYPE;
	private final static String CONTAINER_TYPE = VicUIDataAdapter.CONTAINER_TYPE;
	private final static String VIC_VAPP_IDENTIFIER =
			"vSphere Integrated Containers";
	private final static String[] VAPP_PROPERTIES_TO_FETCH = new String[]{
			"name", "overallStatus", "summary.config.entity",
			"summary.quickStats", "summary.config.cpuAllocation",
			"summary.config.memoryAllocation", "vAppConfig.annotation", "vm"};
	private final static String RESOURCE_ID_FOR_ALL_VICVMS = "vic/ALL";
	private final static Log _logger = LogFactory.getLog(ObjectStore.class);

	public ObjectStore(
			DataService dataService,
			ModelObjectUriResolver modelObjectUriResolver,
			PropFetcher propFetcher,
			VimObjectReferenceService objectRefService) throws IOException {
		if (dataService == null ||
			modelObjectUriResolver == null ||
			propFetcher == null ||
			objectRefService == null) {
			throw new IllegalArgumentException("constructor arg cannot be null");
		}
		_dataService = dataService;
		_modelObjectUriResolver = modelObjectUriResolver;
		_currentRootObject = null;
		_configLoader = new ConfigLoader(CONFIG_PROP_FILES);
		_propFetcher = propFetcher;
		_objectRefService = objectRefService;
	}

	/**
	 * Initialize the ObjectStore
	 * @throws IOException
	 */
	public void init() {
		_currentRootObject = new Root(
				new RootInfo(new String[]{
						_configLoader.getProp(UI_VERSION_CONFIG_KEY)}), 0, 0);
	}

	public void destroy() {
		// nothing to clean up yet
	}

	/**
	 * @return the Root object that contains the # of VirtualContainerHostVms
	 *         and # of ContainerVms
	 * @throws RuntimeFaultFaultMsg
	 * @throws InvalidPropertyFaultMsg
	 */
	private Root getRootObject() {
		synchronized(_currentRootObject) {
			List<ResultItem> vicVappResultItems = getVicVappResultItems();
			int numberOfVchs = vicVappResultItems.size();
			int numberOfContainers = 0;

			// calculates the number of vch vms and container vms
			for (ResultItem ri : vicVappResultItems) {
				for (PropertyValue pv : ri.properties) {
					if (pv.propertyName == "vm") {
						Object[] vmRefs = (Object[])pv.value;
						// # of container VMs = All VMs - VCH VM
						numberOfContainers += (vmRefs.length - 1);
					}
				}
			}
			Root rootObj = new Root(
					new RootInfo(new String[]{
							_configLoader.getProp(UI_VERSION_CONFIG_KEY)}),
					numberOfVchs,
					numberOfContainers);
			_currentRootObject = rootObj;

			return _currentRootObject;
		}
	}

	/**
	 * Get Root, VirtualContainerHostVm or ContainerVm based on URI
	 * @param uri
	 * @return All VCH VMs if uri relates to vic:VirtualContainerHostVm:vic/ALL.
               Otherwise returns the specified VCH VM.
               If the resourceType is vic:ContainerVm, returns all Container VMs.
               Also returns parent vApp's information.
	 */
	public ModelObject getObj(URI uri) {
		String resourceType = _modelObjectUriResolver.getResourceType(uri);

		if (ROOT_TYPE.equals(resourceType)) {
			return getRootObject();
		} else if (VCH_TYPE.equals(resourceType)) {
			return getVms(uri, true);
		} else if (CONTAINER_TYPE.equals(resourceType)) {
			return getVms(uri, false);
		}
		return null;
	}

	/**
	 * Get Root model's URI. This is to be used with Simple Constraint
	 * @return URI for Root model
	 */
	public URI getRootUri() {
		try {
			URI uri = new URI("urn", String.format("%s:%s:%s/%s",
					"vic", ROOT_TYPE, "vic", "vic-root"), null);
			return uri;
		} catch (URISyntaxException e) {
			_logger.error(e.getMessage());
			return null;
		}
	}

	/**
	 * Get vApp(s) and return VMs
	 * @param uri
	 * @param isVch
	 * @return VmQueryResult containing results for VCH VMs or Container VMs
	 */
	private VmQueryResult getVms(URI uri, boolean isVch) {
		String resourceId = _modelObjectUriResolver.getId(uri);
		List<ResultItem> vicVappResultItems = null;

		if (RESOURCE_ID_FOR_ALL_VICVMS.equals(resourceId)) {
			vicVappResultItems = getVicVappResultItems();
		} else {
			vicVappResultItems = getVicVappResultItems(uri);
		}

		Map<Object, ResultItem> vAppVmsMap = new HashMap<Object, ResultItem>();
		for (ResultItem ri : vicVappResultItems) {
			ResultItem result = null;
			Object vAppObjReference = null;
			for (PropertyValue pv : ri.properties) {
				if (pv.propertyName == "summary.config.entity") {
					try {
						vAppObjReference = pv.value;
						result = _propFetcher.getVmsBelongingToMor(
								vAppObjReference, isVch);
					} catch (InvalidPropertyFaultMsg | RuntimeFaultFaultMsg e) {
						_logger.error(e.getMessage());
					}
				}
			}
			ResultItem rri = new ResultItem();
			List<PropertyValue> pvList = new ArrayList<PropertyValue>();
			PropertyValue vAppPropertyValues = new PropertyValue();
			vAppPropertyValues.propertyName = "vAppPropertyValues";
			vAppPropertyValues.value = ri.properties;
			pvList.add(vAppPropertyValues);
			PropertyValue vmsResultItem = new PropertyValue();
			vmsResultItem.propertyName = "vmsResultItem";
			vmsResultItem.value = result;
			pvList.add(vmsResultItem);
			rri.properties = pvList.toArray(new PropertyValue[]{});

			vAppVmsMap.put(vAppObjReference, rri);
		}
		VmQueryResult vmQueryResult = new VmQueryResult(
				uri,
				vAppVmsMap,
				_objectRefService);

		return vmQueryResult;
	}

	/**
	 * Get all vApp ManagedObjectReferences
	 * @return ArrayList containing ResultItem objects for
               vApp ManagedObjectReferences
	 */
	private List<ResultItem> getVicVappResultItems() {
		PropertySpec propertySpec = new PropertySpec();
		propertySpec.propertyNames = VAPP_PROPERTIES_TO_FETCH;
		propertySpec.type = "VirtualApp";
		PropertySpec[] propertySpecs = new PropertySpec[]{propertySpec};

		Constraint constraint = new Constraint();
		constraint.targetType = "VirtualApp";
		List<Constraint> propertyConstraints = new ArrayList<Constraint>();
		propertyConstraints.add(constraint);
		CompositeConstraint cConstraint = new CompositeConstraint();
		cConstraint.conjoiner = Conjoiner.OR;
		cConstraint.nestedConstraints = propertyConstraints.toArray(new Constraint[]{});
		List<ResultItem> riList = getResultItemsFromDataService(cConstraint, propertySpecs);
		List<ResultItem> filtered = filterVicVapps(riList);

		return filtered;
	}

	/**
	 * Get vApp specified by URI
	 * @param uri
	 * @return ArrayList containing ResultItem object for the
	 *         specified vApp ManagedObjectReference
	 */
	private List<ResultItem> getVicVappResultItems(URI uri) {
		PropertySpec propertySpec = new PropertySpec();
		propertySpec.propertyNames = VAPP_PROPERTIES_TO_FETCH;
		propertySpec.type = "VirtualApp";
		PropertySpec[] propertySpecs = new PropertySpec[]{propertySpec};

		ObjectIdentityConstraint oic = new ObjectIdentityConstraint();
		oic.targetType = "VirtualApp";

		ManagedObjectReference vAppMor = new ManagedObjectReference(
				oic.targetType,
				_modelObjectUriResolver.getObjectId(uri),
				_modelObjectUriResolver.getServerGuid(uri));

		oic.target = vAppMor;
		CompositeConstraint cConstraint = new CompositeConstraint();
		cConstraint.conjoiner = Conjoiner.OR;
		cConstraint.targetType = "VirtualApp";
		cConstraint.nestedConstraints = new Constraint[]{oic};
		List<ResultItem> riList = getResultItemsFromDataService(cConstraint, propertySpecs);

		return riList;
	}

	/**
	 * Call vSphere Client DataService to fetch information on vSphere objects
	 * @param constraint
	 * @param propertySpecs
	 * @return ArrayList containing ResultItem objects for the retrieved data
	 */
	private List<ResultItem> getResultItemsFromDataService(
			CompositeConstraint constraint,
			PropertySpec[] propertySpecs) {

		QuerySpec qSpec = new QuerySpec();
		qSpec.resourceSpec = new ResourceSpec();
		qSpec.resourceSpec.constraint = constraint;
		qSpec.resourceSpec.propertySpecs = propertySpecs;
		qSpec.resultSpec = new ResultSpec();

		RequestSpec reqSpec = new RequestSpec();
		reqSpec.querySpec = new QuerySpec[]{qSpec};

		Response response = _dataService.getData(reqSpec);

		if ((response == null) || (response.resultSet.length == 0)) {
			return new ArrayList<ResultItem>();
		}

		for (ResultSet rs : response.resultSet) {
			if (rs.error != null) {
				_logger.error(rs.error.getMessage());
				rs.error.printStackTrace();
				continue;
			}
		}

		return Arrays.asList(response.resultSet[0].items);
	}

	/**
	 * Filter out vApps that are not a vSphere Integrated Container appliance
	 * @param resultItems
	 * @return ArrayList containing ResultItems containing all vApps
	 *         that are vSphere Integrated Container appliances
	 */
	private List<ResultItem> filterVicVapps(List<ResultItem> resultItems) {
		List<ResultItem> results = new ArrayList<ResultItem>();
		if (resultItems == null) {
			return results;
		}

		for (ResultItem ri : resultItems.toArray(new ResultItem[]{})) {
			for (PropertyValue pv : ri.properties) {
				if (pv.propertyName == "vAppConfig.annotation") {
					if (VIC_VAPP_IDENTIFIER.equals(pv.value)) {
						results.add(ri);
					}
				}
			}
		}
		return results;
	}
}
