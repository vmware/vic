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

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import java.util.ArrayList;
import java.util.List;

import com.vmware.vic.model.ContainerVm;
import com.vmware.vic.model.VirtualContainerHostVm;
import com.vmware.vim25.ArrayOfOptionValue;
import com.vmware.vim25.DynamicProperty;
import com.vmware.vim25.ManagedEntityStatus;
import com.vmware.vim25.ManagedObjectReference;
import com.vmware.vim25.ObjectContent;
import com.vmware.vim25.OptionValue;
import com.vmware.vim25.VirtualMachinePowerState;
import com.vmware.vim25.VirtualMachineQuickStats;
import com.vmware.vim25.VirtualMachineStorageSummary;
import com.vmware.vim25.VirtualMachineSummary;

public class Common {
	protected VirtualContainerHostVm getMockedVirtualContainerHostVm(
			String serverGuid,
			String vmObjectId,
			String vmName,
			ManagedEntityStatus overallStatus,
			VirtualMachinePowerState powerState,
			String clientIpBase64Encoded) {
		List<DynamicProperty> dps = new ArrayList<DynamicProperty>();
		DynamicProperty dpName = new DynamicProperty();
		dpName.setName("name");
		dpName.setVal(vmName);
		dps.add(dpName);

		DynamicProperty dpOverallStatus = new DynamicProperty();
		dpOverallStatus.setName("overallStatus");
		dpOverallStatus.setVal(overallStatus);
		dps.add(dpOverallStatus);

		DynamicProperty dpPowerState = new DynamicProperty();
		dpPowerState.setName("runtime.powerState");
		dpPowerState.setVal(powerState);
		dps.add(dpPowerState);

		DynamicProperty dpSummary = new DynamicProperty();
		VirtualMachineSummary vmSummary = new VirtualMachineSummary();
		VirtualMachineQuickStats vmQuickStatsMock = mock(VirtualMachineQuickStats.class);
		when(vmQuickStatsMock.getGuestMemoryUsage()).thenReturn(500);
		when(vmQuickStatsMock.getOverallCpuUsage()).thenReturn(1000);
		vmSummary.setQuickStats(vmQuickStatsMock);
		VirtualMachineStorageSummary vmStorageSummaryMock = mock(VirtualMachineStorageSummary.class);
		when(vmStorageSummaryMock.getCommitted()).thenReturn((long)123456789);
		vmSummary.setStorage(vmStorageSummaryMock);

		dpSummary.setName("summary");
		dpSummary.setVal(vmSummary);
		dps.add(dpSummary);

		DynamicProperty dpConfigExtraConfig = new DynamicProperty();
		ArrayOfOptionValue ovArrayMock = mock(ArrayOfOptionValue.class);
		List<OptionValue> ovList = new ArrayList<OptionValue>();
		OptionValue ovClientIpKey = new OptionValue();
		ovClientIpKey.setKey("guestinfo.vice..init.networks|client.assigned.IP");
		ovClientIpKey.setValue(clientIpBase64Encoded);
		ovList.add(ovClientIpKey);
		when(ovArrayMock.getOptionValue()).thenReturn(ovList);
		dpConfigExtraConfig.setName("config.extraConfig");
		dpConfigExtraConfig.setVal(ovArrayMock);
		dps.add(dpConfigExtraConfig);

		ManagedObjectReference mor = new ManagedObjectReference();
		mor.setType("VirtualMachine");
		mor.setValue(vmObjectId);

		ObjectContent objContent = mock(ObjectContent.class);
		when(objContent.getObj()).thenReturn(mor);
		when(objContent.getPropSet()).thenReturn(dps);

		VirtualContainerHostVm vm = new VirtualContainerHostVm(objContent, serverGuid);
		return vm;
	}

	protected ContainerVm getMockedContainerVm(
			String serverGuid,
			String vmObjectId,
			String vmName,
			ManagedEntityStatus overallStatus,
			VirtualMachinePowerState powerState,
			String containerName,
			String imageName,
			String portMapping) {
		// mock ObjectContent object and its members
		List<DynamicProperty> dps = new ArrayList<DynamicProperty>();
		DynamicProperty dpName = new DynamicProperty();
		dpName.setName("name");
		dpName.setVal(vmName);
		dps.add(dpName);

		DynamicProperty dpOverallStatus = new DynamicProperty();
		dpOverallStatus.setName("overallStatus");
		dpOverallStatus.setVal(overallStatus);
		dps.add(dpOverallStatus);

		DynamicProperty dpPowerState = new DynamicProperty();
		dpPowerState.setName("runtime.powerState");
		dpPowerState.setVal(powerState);
		dps.add(dpPowerState);

		DynamicProperty dpSummary = new DynamicProperty();
		VirtualMachineSummary vmSummary = new VirtualMachineSummary();
		VirtualMachineQuickStats vmQuickStatsMock = mock(VirtualMachineQuickStats.class);
		when(vmQuickStatsMock.getGuestMemoryUsage()).thenReturn(500);
		when(vmQuickStatsMock.getOverallCpuUsage()).thenReturn(1000);
		vmSummary.setQuickStats(vmQuickStatsMock);
		VirtualMachineStorageSummary vmStorageSummaryMock = mock(VirtualMachineStorageSummary.class);
		when(vmStorageSummaryMock.getCommitted()).thenReturn((long)123456789);
		vmSummary.setStorage(vmStorageSummaryMock);

		dpSummary.setName("summary");
		dpSummary.setVal(vmSummary);
		dps.add(dpSummary);

		DynamicProperty dpConfigExtraConfig = new DynamicProperty();
		ArrayOfOptionValue ovArrayMock = mock(ArrayOfOptionValue.class);
		List<OptionValue> ovList = new ArrayList<OptionValue>();

		OptionValue ovContainerNameKey = new OptionValue();
		ovContainerNameKey.setKey("guestinfo.vice./common/name");
		ovContainerNameKey.setValue(containerName);
		ovList.add(ovContainerNameKey);

		OptionValue ovImageNameKey = new OptionValue();
		ovImageNameKey.setKey("guestinfo.vice./repo");
		ovImageNameKey.setValue(imageName);
		ovList.add(ovImageNameKey);

		if (portMapping != null) {
			OptionValue ovPortMappingKey = new OptionValue();
			ovPortMappingKey.setKey("guestinfo.vice./networks|bridge/ports~");
			ovPortMappingKey.setValue(portMapping);
			ovList.add(ovPortMappingKey);
		}

		when(ovArrayMock.getOptionValue()).thenReturn(ovList);
		dpConfigExtraConfig.setName("config.extraConfig");
		dpConfigExtraConfig.setVal(ovArrayMock);
		dps.add(dpConfigExtraConfig);

		ManagedObjectReference mor = new ManagedObjectReference();
		mor.setType("VirtualMachine");
		mor.setValue(vmObjectId);

		ObjectContent objContent = mock(ObjectContent.class);
		when(objContent.getObj()).thenReturn(mor);
		when(objContent.getPropSet()).thenReturn(dps);

		ContainerVm vm = new ContainerVm(objContent, serverGuid);
		return vm;
	}
}
