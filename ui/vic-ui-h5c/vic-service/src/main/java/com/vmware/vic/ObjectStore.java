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
import java.util.HashMap;
import java.util.Map;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.vmware.vic.model.ModelObject;
import com.vmware.vic.model.Root;
import com.vmware.vic.model.RootInfo;
import com.vmware.vic.model.VmQueryResult;
import com.vmware.vic.utils.ConfigLoader;
import com.vmware.vim25.InvalidPropertyFaultMsg;
import com.vmware.vim25.RuntimeFaultFaultMsg;
import com.vmware.vise.data.query.DataService;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.ResultItem;
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
			ResultItem vchsRi = _propFetcher.getVicVms(true);
			ResultItem containersRi = _propFetcher.getVicVms(false);

			int numberOfVchs = vchsRi.properties.length;
			int numberOfContainers = containersRi.properties.length;

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
		ResultItem vmsRi = _propFetcher.getVicVms(isVch);
		Map<String, ModelObject> resultsMap = new HashMap<String, ModelObject>();
		for (PropertyValue pv : vmsRi.properties) {
		    if (pv.propertyName == "vm") {
		        ModelObject mo = (ModelObject)pv.value;
		        resultsMap.put(mo.getId(), mo);
		    }
		}

		VmQueryResult vmQueryResult = new VmQueryResult(
				resultsMap,
				_objectRefService);

		return vmQueryResult;
	}

	public URI createUri(String type, String id) {
		return _modelObjectUriResolver.createUri(type, id);
	}

}
