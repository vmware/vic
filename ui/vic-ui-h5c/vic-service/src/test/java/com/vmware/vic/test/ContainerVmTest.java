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

import static org.junit.Assert.*;

import java.net.URI;

import org.junit.Test;

import com.vmware.vic.ModelObjectUriResolver;
import com.vmware.vic.model.ContainerVm;
import com.vmware.vic.model.ModelObject;
import com.vmware.vic.model.constants.BaseVm;
import com.vmware.vic.model.constants.Container;
import com.vmware.vim25.ManagedEntityStatus;
import com.vmware.vim25.VirtualMachinePowerState;

public class ContainerVmTest extends Common {
    private static final String VM_TYPE_CONTAINERVM = "vic:ContainerVm";
	private ContainerVm _vm;

	private ContainerVm createMockContainerVmWithoutPortMapping() {
		ContainerVm vm = getMockedContainerVm(
				"server1",
				"id1",
				"container vm 1",
				ManagedEntityStatus.GREEN,
				VirtualMachinePowerState.POWERED_ON,
				"container without portmapping info",
				"busybox",
				null);
		return vm;
	}

	private ContainerVm createMockContainerVmWithPortMapping() {
		ContainerVm vm = getMockedContainerVm(
				"server2",
				"id2",
				"container vm 2",
				ManagedEntityStatus.GRAY,
				VirtualMachinePowerState.SUSPENDED,
				"container with portmapping info",
				"nginx",
				"8080:80/tcp");
		return vm;
	}

	@Test
	public void testGetTypeForContainerVmWithoutPortMapping() {
		_vm = createMockContainerVmWithoutPortMapping();
		assertNotNull(_vm);
		assertEquals(VM_TYPE_CONTAINERVM, _vm.getType());
	}

	@Test
	public void testGetIdForContainerVmWithoutPortMapping() {
		_vm = createMockContainerVmWithoutPortMapping();
		assertEquals("server1/id1", _vm.getId());
	}

	@Test
	public void testGetUriForContainerVmWithoutPortMapping() {
		_vm = createMockContainerVmWithoutPortMapping();
		ModelObjectUriResolver uriResolver = new ModelObjectUriResolver();
		URI uri = _vm.getUri(uriResolver);
		assertEquals(String.format(
		        "urn:vic:%s:%s", VM_TYPE_CONTAINERVM, "server1/id1"),
		        uri.toString());
	}

	@Test
	public void testGetPropertyForContainerVmWithoutPortMapping() {
		_vm = createMockContainerVmWithoutPortMapping();
		assertTrue(_vm.getProperty(BaseVm.VM_NAME).equals("container vm 1"));
		assertTrue(_vm.getProperty(BaseVm.VM_OVERALL_STATUS)
		        .equals(ManagedEntityStatus.GREEN));
		assertTrue(_vm.getProperty(BaseVm.Runtime.VM_POWERSTATE_FULLPATH)
	            .equals(VirtualMachinePowerState.POWERED_ON));
		assertTrue(_vm.getProperty(Container.VM_CONTAINERNAME_KEY)
		        .equals("container without portmapping info"));
		assertTrue(_vm.getProperty(Container.VM_IMAGENAME_KEY)
		        .equals("busybox"));
		assertTrue(_vm.getProperty(BaseVm.VM_OVERALLCPUUSAGE).equals(1000));
		assertTrue(_vm.getProperty(BaseVm.VM_GUESTMEMORYUSAGE).equals(500));
		assertTrue(_vm.getProperty(BaseVm.VM_COMMITTEDSTORAGE)
		        .equals((long)123456789));
		assertTrue(_vm.getProperty("iDontExist")
		        .equals(ModelObject.UNSUPPORTED_PROPERTY));
	}

	@Test
	public void testGettersForContainerVmWithoutPortMapping() {
		_vm = createMockContainerVmWithoutPortMapping();
		assertEquals(_vm.getName(), "container vm 1");
		assertEquals(_vm.getOverallStatus(), "GREEN");
		assertEquals(_vm.getPowerState(), "POWERED_ON");
		assertEquals(_vm.getContainerName(), "container without portmapping info");
		assertEquals(_vm.getPortMapping(), null);
		assertEquals(_vm.getImageName(), "busybox");
		assertEquals(_vm.getOverallCpuUsage(), 1000);
		assertEquals(_vm.getGuestMemoryUsage(), 500);
		assertEquals(_vm.getCommittedStorage(), (long)123456789);
	}

	@Test
	public void testGetTypeForContainerVmWithPortMapping() {
		_vm = createMockContainerVmWithPortMapping();
		assertNotNull(_vm);
		assertEquals("vic:ContainerVm", _vm.getType());
	}

	@Test
	public void testGetIdForContainerVmWithPortMapping() {
		_vm = createMockContainerVmWithPortMapping();
		assertEquals("server2/id2", _vm.getId());
	}

	@Test
	public void testGetUriForContainerVmWithPortMapping() {
		_vm = createMockContainerVmWithPortMapping();
		ModelObjectUriResolver uriResolver = new ModelObjectUriResolver();
		URI uri = _vm.getUri(uriResolver);
		assertEquals("urn:vic:vic:ContainerVm:server2/id2", uri.toString());
	}

	@Test
	public void testGetPropertyForContainerVmWithPortMapping() {
		_vm = createMockContainerVmWithPortMapping();
		assertTrue(_vm.getProperty(BaseVm.VM_NAME).equals("container vm 2"));
		assertTrue(_vm.getProperty(BaseVm.VM_OVERALL_STATUS)
		        .equals(ManagedEntityStatus.GRAY));
		assertTrue(_vm.getProperty(BaseVm.Runtime.VM_POWERSTATE_FULLPATH)
		        .equals(VirtualMachinePowerState.SUSPENDED));
		assertTrue(_vm.getProperty(Container.VM_CONTAINERNAME_KEY)
		        .equals("container with portmapping info"));
		assertTrue(_vm.getProperty(Container.VM_IMAGENAME_KEY).equals("nginx"));
		assertTrue(_vm.getProperty(Container.VM_PORTMAPPING_KEY)
		        .equals("8080:80/tcp"));
		assertTrue(_vm.getProperty(BaseVm.VM_OVERALLCPUUSAGE).equals(1000));
		assertTrue(_vm.getProperty(BaseVm.VM_GUESTMEMORYUSAGE).equals(500));
		assertTrue(_vm.getProperty(BaseVm.VM_COMMITTEDSTORAGE)
		        .equals((long)123456789));
		assertTrue(_vm.getProperty("iDontExist")
		        .equals(ModelObject.UNSUPPORTED_PROPERTY));
	}

	@Test
	public void testGettersForContainerVmWithPortMapping() {
		_vm = createMockContainerVmWithPortMapping();
		assertEquals(_vm.getName(), "container vm 2");
		assertEquals(_vm.getOverallStatus(), "GRAY");
		assertEquals(_vm.getPowerState(), "SUSPENDED");
		assertEquals(_vm.getContainerName(), "container with portmapping info");
		assertEquals(_vm.getImageName(), "nginx");
		assertEquals(_vm.getPortMapping(), "8080:80/tcp");
		assertEquals(_vm.getOverallCpuUsage(), 1000);
		assertEquals(_vm.getGuestMemoryUsage(), 500);
		assertEquals(_vm.getCommittedStorage(), (long)123456789);
	}
}
