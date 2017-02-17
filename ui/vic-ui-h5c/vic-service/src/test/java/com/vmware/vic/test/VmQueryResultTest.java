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
package com.vmware.vic.test;

import static org.junit.Assert.assertEquals;

import org.junit.Test;

import com.vmware.vic.model.ContainerVm;
import com.vmware.vic.model.VirtualContainerHostVm;
import com.vmware.vic.model.VmQueryResult;
import com.vmware.vim25.ManagedEntityStatus;
import com.vmware.vim25.VirtualMachinePowerState;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.vim.data.VimObjectReferenceService;

import static org.mockito.Mockito.*;

import java.net.URI;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.Assert.*;

public class VmQueryResultTest extends Common {
	private ResultItem mockResultItem() {
		ResultItem ri = new ResultItem();
		List<PropertyValue> pvList = new ArrayList<PropertyValue>();
		PropertyValue pvVmsResultItem = new PropertyValue();
		pvVmsResultItem.propertyName = "vmsResultItem";
		pvVmsResultItem.value = getMockedVmsResultItem();
		pvList.add(pvVmsResultItem);

		PropertyValue pvVappPropertyValues = new PropertyValue();
		pvVappPropertyValues.propertyName = "vAppPropertyValues";
		pvVappPropertyValues.value = getVappPropertyValuesArray();
		pvList.add(pvVappPropertyValues);

		ri.properties = pvList.toArray(new PropertyValue[]{});
		return ri;
	}

	private ResultItem getMockedVmsResultItem() {
		ResultItem ri = new ResultItem();
		List<PropertyValue> pvList = new ArrayList<PropertyValue>();

		PropertyValue pvVchVm = new PropertyValue();
		pvVchVm.propertyName = "vm";
		pvVchVm.value = getMockedVchVm();
		pvList.add(pvVchVm);

		PropertyValue pvContainerVm = new PropertyValue();
		pvContainerVm.propertyName = "vm";
		pvContainerVm.value = getMockedContainerVm();
		pvList.add(pvContainerVm);

		ri.properties = pvList.toArray(new PropertyValue[]{});
		return ri;
	}

	private VirtualContainerHostVm getMockedVchVm() {
		return getMockedVirtualContainerHostVm(
				"server3",
				"id5",
				"vch vm",
				ManagedEntityStatus.GREEN,
				VirtualMachinePowerState.POWERED_ON,
				"ChFtuw==");
	}

	private ContainerVm getMockedContainerVm() {
		return getMockedContainerVm(
				"server5",
				"id7",
				"container vm",
				ManagedEntityStatus.GREEN,
				VirtualMachinePowerState.POWERED_ON,
				"container name",
				"nginx:alpine",
				"8088:80/tcp");
	}

	private PropertyValue[] getVappPropertyValuesArray() {
		// not necessary so we're returning an empty array
		List<PropertyValue> pvList = new ArrayList<PropertyValue>();
		return pvList.toArray(new PropertyValue[]{});
	}

	/**
	 * Set up vApp where there is one vApp which has one VCH VM
	   and one Container VM in it.
	 */
	@Test
	public void prepareVmQueryResult() {
		URI uri = null;
		Map<Object, ResultItem> vAppMap = new HashMap<Object, ResultItem>();

		Object keyObj = new Object();
		ResultItem ri = mockResultItem();
		vAppMap.put(keyObj, ri);

		VimObjectReferenceService objRefService = mock(VimObjectReferenceService.class);
		when(objRefService.getServerGuid(keyObj)).thenReturn("server-1");
		when(objRefService.getValue(keyObj)).thenReturn("vapp-id-1");

		// create an instance based on mocked data
		VmQueryResult vmQueryResult = new VmQueryResult(uri, vAppMap, objRefService);
		assertNotNull(vmQueryResult);
		assertEquals(vmQueryResult.getProperty("match"), 2);
	}
}
