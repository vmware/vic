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
import com.vmware.vic.model.ModelObject;
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

	/**
	 * Set up vApp where there is one vApp which has one VCH VM
	   and one Container VM in it.
	 */
	@Test
	public void vmQueryResultForVch() {
		Map<String, ModelObject> vmsMap = new HashMap<String, ModelObject>();
		vmsMap.put("id5", getMockedVchVm());

		VimObjectReferenceService objRefService = mock(VimObjectReferenceService.class);

		// create an instance based on mocked data
		VmQueryResult vmQueryResult = new VmQueryResult(vmsMap, objRefService);
		assertNotNull(vmQueryResult);
		assertEquals(vmQueryResult.getProperty("match"), 1);
	}
	
	/**
     * Set up vApp where there is one vApp which has one VCH VM
       and one Container VM in it.
     */
    @Test
    public void vmQueryResultForContainer() {
        Map<String, ModelObject> vmsMap = new HashMap<String, ModelObject>();
        vmsMap.put("id7", getMockedContainerVm());

        VimObjectReferenceService objRefService = mock(VimObjectReferenceService.class);

        // create an instance based on mocked data
        VmQueryResult vmQueryResult = new VmQueryResult(vmsMap, objRefService);
        assertNotNull(vmQueryResult);
        assertEquals(vmQueryResult.getProperty("match"), 1);
    }
}
