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

import java.util.List;

import com.vmware.vic.model.constants.BaseVm;
import com.vmware.vic.model.constants.Vch;
import com.vmware.vim25.DynamicProperty;
import com.vmware.vim25.ManagedEntityStatus;
import com.vmware.vim25.ManagedObjectReference;
import com.vmware.vim25.ObjectContent;
import com.vmware.vim25.ResourceConfigSpec;
import com.vmware.vim25.VirtualMachinePowerState;
import com.vmware.vim25.VirtualMachineQuickStats;
import com.vmware.vim25.VirtualMachineStorageSummary;
import com.vmware.vim25.VirtualMachineSummary;

public abstract class VicBaseVm extends ModelObject {
	protected static final String VM_KEY_NAME = BaseVm.VM_NAME;
	protected static final String VM_KEY_OVERALL_STATUS =
	        BaseVm.VM_OVERALL_STATUS;
	protected static final String VM_KEY_POWERSTATE =
	        BaseVm.Runtime.VM_POWERSTATE_FULLPATH;
	protected static final String VM_KEY_SUMMARY = BaseVm.VM_SUMMARY;
	protected static final String VM_KEY_GUESTFULLNAME =
	        BaseVm.Config.VM_GUESTFULLNAME;
	protected static final String VM_KEY_CONFIG_EXTRACONFIG =
	        BaseVm.Config.VM_EXTRACONFIG;
	protected static final String VM_KEY_RESOURCECONFIG =
	        BaseVm.VM_RESOURCECONFIG;
	protected static final String VM_KEY_RESOURCEPOOL =
	        BaseVm.VM_RESOURCEPOOL;
	protected static final String VM_KEY_CLIENT_IP = Vch.VM_CLIENT_IP;
	protected static final String VM_KEY_OVERALLCPUUSAGE =
	        BaseVm.VM_OVERALLCPUUSAGE;
	protected static final String VM_KEY_GUESTMEMORYUSAGE =
	        BaseVm.VM_GUESTMEMORYUSAGE;
	protected static final String VM_KEY_COMMITTEDSTORAGE =
	        BaseVm.VM_COMMITTEDSTORAGE;
	protected final ManagedObjectReference _objectRef;
	protected String _vmName = null;
	protected String _guestFullName = null;
	protected ResourceConfigSpec _resourceConfig = null;
	protected Object _resourcePool = null;
	protected int _overallCpuUsage;
	protected int _guestMemoryUsage;
	protected long _committedStorage;
	protected VirtualMachinePowerState _powerState = null;
	protected ManagedEntityStatus _overallStatus = null;

	public VicBaseVm(
			ObjectContent objContent,
			String serverGuid) {
		if (objContent == null) {
			throw new IllegalArgumentException("constructor argument cannot be null");
		}
		_objectRef = objContent.getObj();
		this.setId(serverGuid + "/" + _objectRef.getValue());
	}

	abstract protected void processDynamicProperties(List<DynamicProperty> dpsList);

	/**
	 * Process VirtualMachineSummary to extract quickStats and storage info
	 * @param summary
	 */
	protected void processVmSummary(VirtualMachineSummary summary) {
		VirtualMachineQuickStats quickStats = summary.getQuickStats();
		if (quickStats != null) {
			_overallCpuUsage = quickStats.getOverallCpuUsage();
			_guestMemoryUsage = quickStats.getGuestMemoryUsage();
		}

		VirtualMachineStorageSummary vmStorageSummary = summary.getStorage();
		if (vmStorageSummary != null) {
			_committedStorage = vmStorageSummary.getCommitted();
		}
	}

	public String getName() {
		return _vmName;
	}

	public String getOverallStatus() {
		return _overallStatus.toString();
	}

	public String getPowerState() {
		return _powerState.toString();
	}
	
	public String getGuestFullName() {
	    return _guestFullName;
	}

	public int getOverallCpuUsage() {
		return _overallCpuUsage;
	}

	public int getGuestMemoryUsage() {
		return _guestMemoryUsage;
	}

	public long getCommittedStorage() {
		return _committedStorage;
	}
	
	public ResourceConfigSpec getResourceConfig() {
	    return _resourceConfig;
	}
	
	public Object getResourcePool() {
	    return _resourcePool;
	}

}
