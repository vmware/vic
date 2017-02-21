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

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.vmware.vise.data.query.DataServiceExtensionRegistry;
import com.vmware.vise.data.query.PropertyProviderAdapter;
import com.vmware.vise.data.query.PropertyRequestSpec;
import com.vmware.vise.data.query.ResultSet;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.data.query.TypeInfo;

public class VicUIPropertyProvider implements PropertyProviderAdapter {
	private static final Log _logger = LogFactory.getLog(VicUIPropertyProvider.class);
	private static final String[] VIC_VM_TYPES = {"isVCH", "isContainer"};
	private final PropFetcher _propFetcher;

	public VicUIPropertyProvider(
			DataServiceExtensionRegistry extensionRegistry,
			PropFetcher propFetcher) {
		TypeInfo vmTypeInfo = new TypeInfo();
		vmTypeInfo.type = "VirtualMachine";
		vmTypeInfo.properties = VIC_VM_TYPES;
		TypeInfo[] providerTypes = new TypeInfo[] { vmTypeInfo };

		_propFetcher = propFetcher;
		extensionRegistry.registerDataAdapter(this, providerTypes);
	}

	@Override
	public ResultSet getProperties(PropertyRequestSpec propertyRequest) {
		ResultSet resultSet = new ResultSet();

		try {
			List<ResultItem> resultItems = new ArrayList<ResultItem>();

			for (Object objRef : propertyRequest.objects) {
				ResultItem resultItem = _propFetcher.getVmProperties(objRef);
				if (resultItem != null) {
					resultItems.add(resultItem);
				}
			}

			resultSet.items = resultItems.toArray(new ResultItem[] {});

		} catch (Exception e) {
			_logger.error("VicUIServiceImpl.getProperties error: " + e);
		}

		return resultSet;
	}
}