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
package com.vmware.vic.model;

import java.net.URI;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

//import org.apache.commons.logging.Log;
//import org.apache.commons.logging.LogFactory;

import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.vim.data.VimObjectReferenceService;

public class VmQueryResult extends ModelObject {
	private final URI _objectRef;
	private final Map<String, Map<String, Object>> _vAppVmsMap;
	private int _match = 0;

	public VmQueryResult(
			URI objectRef,
			Map<Object, ResultItem> vAppMap,
			VimObjectReferenceService vimObjectReferenceService) {
		Set<Object> vAppVmsMapKeys = vAppMap.keySet();
		Map<String, Map<String, Object>> vAppVmsMap = new HashMap<String, Map<String, Object>>();

		if (vAppMap == null ||
			vimObjectReferenceService == null) {
			throw new IllegalArgumentException("constructor argument cannot be empty!");
		}

		for (Object objRef : vAppVmsMapKeys) {
			String vAppServerGuid = vimObjectReferenceService.getServerGuid(objRef);
			String vAppResourceId = vimObjectReferenceService.getValue(objRef);
			String vAppUid = vAppServerGuid + "/" + vAppResourceId;

			ResultItem ri = vAppMap.get(objRef);
			Map<String, Object> vAppDetailsMap = processVappVmsResultItem(ri);
			vAppVmsMap.put(vAppUid, vAppDetailsMap);
		}
		_vAppVmsMap = vAppVmsMap;
		_objectRef = objectRef;
	}

	/**
	 * Format ResultItem data suitable for client use
	 * @param ri
	 * @return HashMap<String, Object> object where String is either "vms" or
	 *         "vAppPropertyValues" and Object is their respective data
	 */
	private Map<String, Object> processVappVmsResultItem(ResultItem ri) {
		Map<String, Object> vAppDetailsMap = new HashMap<String, Object>();
		List<ModelObject> vmList = new ArrayList<ModelObject>();

		for (PropertyValue vAppVmsPv : ri.properties) {
			if ("vmsResultItem".equals(vAppVmsPv.propertyName)) {
				ResultItem vmsResultItem = (ResultItem)vAppVmsPv.value;
				for (PropertyValue vmPv : vmsResultItem.properties) {
					ModelObject vm = null;
					boolean shouldAddToList = false;
					if (vmPv.value instanceof ContainerVm) {
						vm = (ContainerVm)vmPv.value;
						shouldAddToList = vm.getProperty("containerName") != null;
					} else if (vmPv.value instanceof VirtualContainerHostVm) {
						vm = (VirtualContainerHostVm)vmPv.value;
						shouldAddToList = vm.getProperty("clientIp") != null;
					}
					// container vm does not have clientIp set up. vch vm does.
					if (shouldAddToList) {
						vmList.add(vm);
						_match += 1;
					}
				}
			} else if ("vAppPropertyValues".equals(vAppVmsPv.propertyName)) {
				for (PropertyValue pv : (PropertyValue[])vAppVmsPv.value) {
					vAppDetailsMap.put(pv.propertyName, pv.value);
				}
			}
		}
		vAppDetailsMap.put("vms", vmList);
		return vAppDetailsMap;
	}

	@Override
	public Object getProperty(String property) {
		if ("match".equals(property)) {
			return _match;
		} else if ("results".equals(property)) {
			return _vAppVmsMap;
		}
		return null;
	}

}
