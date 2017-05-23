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

package com.vmware.vic.utils;

import java.util.Comparator;
import java.util.Map;

import com.vmware.vic.model.ContainerVm;
import com.vmware.vic.model.ModelObject;
import com.vmware.vic.model.VirtualContainerHostVm;
import com.vmware.vic.model.constants.BaseVm;
import com.vmware.vic.model.constants.Container;
import com.vmware.vic.model.constants.Vch;

/**
 * Comparator for VIC VMs
 */
public class VicVmComparator implements Comparator<String> {
	private Map<String, ModelObject> base;
	private String compareBy;
	private boolean reverse;

	public VicVmComparator(
			Map<String, ModelObject> base,
			String compareBy,
			boolean reverse
		) {
		this.base = base;
		this.compareBy = compareBy;
		this.reverse = reverse;
	}

	/**
	 * Comparator.compare() implementation for ModelObject
	 * @param a : objectId of VM 1
	 * @param b : objectId of VM 2
	 */
	@Override
	public int compare(String a, String b) {
		try {
		    ModelObject mo_A = (ModelObject) base.get(a);
	        ModelObject mo_B = (ModelObject) base.get(b);
			int result = 0;
			result = getStringPropertyFromVm(mo_A)
				.compareTo(getStringPropertyFromVm(mo_B));
			if (result == 0) {
			    result = mo_A.hashCode() - mo_B.hashCode();
			}

			return result * (this.reverse ? -1 : 1);
		} catch (IndexOutOfBoundsException e) {
			return 0;
		}
	}

	/**
	 * Retrieve string property value from a ModelObject
	 * @param mo : ModelObject instance
	 * @return property value for the property name specified by compareBy
	 */
	private String getStringPropertyFromVm(ModelObject mo) {
		if (mo instanceof VirtualContainerHostVm) {
			if (BaseVm.ID.equals(compareBy)) {
				return ((VirtualContainerHostVm) mo).getId();
			} else if (BaseVm.VM_NAME.equals(compareBy)) {
				return ((VirtualContainerHostVm) mo).getName();
			} else if (Vch.VM_VCH_IP.equals(compareBy)) {
				return ((VirtualContainerHostVm) mo).getClientIp();
			} else if (BaseVm.VM_OVERALL_STATUS.equals(compareBy)) {
				return ((VirtualContainerHostVm) mo).getOverallStatus();
			}
		} else if (mo instanceof ContainerVm) {
			if (BaseVm.ID.equals(compareBy)) {
			    return ((ContainerVm) mo).getId();
			} else if (Container.VM_CONTAINERNAME_KEY.equals(compareBy)) {
			    return ((ContainerVm) mo).getContainerName();
			} else if (BaseVm.Runtime.VM_POWERSTATE_BASENAME.equals(compareBy)) {
			    return ((ContainerVm) mo).getPowerState();
			} else if (BaseVm.VM_GUESTMEMORYUSAGE.equals(compareBy)) {
			    return Integer.toString(((ContainerVm) mo).getGuestMemoryUsage());
			} else if (BaseVm.VM_OVERALLCPUUSAGE.equals(compareBy)) {
                return Integer.toString(((ContainerVm) mo).getOverallCpuUsage());
			} else if (BaseVm.VM_COMMITTEDSTORAGE.equals(compareBy)) {
                return Long.toString(((ContainerVm) mo).getCommittedStorage());
			} else if (Container.VM_PORTMAPPING_KEY.equals(compareBy)) {
			    String pm = ((ContainerVm) mo).getPortMapping();
                return pm != null ? pm : "";
			} else if (BaseVm.VM_NAME.equals(compareBy)) {
			    return ((ContainerVm) mo).getName();
			} else if (Container.VM_IMAGENAME_KEY.equals(compareBy)) {
			    return ((ContainerVm) mo).getImageName();
			}
		}
		
		return null;
	}

}
