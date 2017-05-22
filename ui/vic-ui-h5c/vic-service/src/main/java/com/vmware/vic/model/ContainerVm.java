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

import com.vmware.vim25.ArrayOfOptionValue;
import com.vmware.vim25.DynamicProperty;
import com.vmware.vim25.ManagedEntityStatus;
import com.vmware.vim25.ObjectContent;
import com.vmware.vim25.OptionValue;
import com.vmware.vim25.ResourceConfigSpec;
import com.vmware.vim25.VirtualMachinePowerState;
import com.vmware.vim25.VirtualMachineSummary;

public class ContainerVm extends VicBaseVm {
	private static final String EXTRACONFIG_CONTAINER_NAME_KEY =
			"common/name";
	private static final String EXTRACONFIG_IMAGE_NAME_KEY =
			"guestinfo.vice./repo";
	private static final String EXTRACONFIG_PORTMAPPING_KEY =
			"guestinfo.vice./networks|bridge/ports~";
	private static final String VM_KEY_CONTAINERNAME = "containerName";
	private static final String VM_KEY_IMAGENAME = "imageName";
	private static final String VM_KEY_PORTMAPPING = "portMapping";
	private String _containerName = null;
	private String _imageName = null;
	private String _portMapping = null;

	public ContainerVm(ObjectContent objContent, String serverGuid) {
		super(objContent, serverGuid);
		processDynamicProperties(objContent.getPropSet());
	}

	public String getContainerName() {
		return _containerName;
	}

	public String getImageName() {
		return _imageName;
	}

	public String getPortMapping() {
		return _portMapping;
	}

	@Override
	public Object getProperty(String property) {
		if ("objectRef".equals(property)) {
			return _objectRef;
		} else if (VM_KEY_NAME.equals(property)) {
			return _vmName;
		} else if (VM_KEY_OVERALL_STATUS.equals(property)) {
			return _overallStatus;
		} else if (VM_KEY_POWERSTATE.equals(property)) {
			return _powerState;
		} else if (VM_KEY_GUESTFULLNAME.equals(property)) {
            return _guestFullName;
		} else if (VM_KEY_OVERALLCPUUSAGE.equals(property)) {
			return _overallCpuUsage;
		} else if (VM_KEY_GUESTMEMORYUSAGE.equals(property)) {
			return _guestMemoryUsage;
		} else if (VM_KEY_COMMITTEDSTORAGE.equals(property)) {
			return _committedStorage;
		} else if (VM_KEY_CONTAINERNAME.equals(property)) {
			return _containerName;
		} else if (VM_KEY_IMAGENAME.equals(property)) {
			return _imageName;
		} else if (VM_KEY_PORTMAPPING.equals(property)) {
			return _portMapping;
		} else if (VM_KEY_RESOURCECONFIG.equals(property)) {
            return _resourceConfig;
        } else if (VM_KEY_RESOURCEPOOL.equals(property)) {
            return _resourcePool;
        }
		return UNSUPPORTED_PROPERTY;
	}

	@Override
	protected void processDynamicProperties(List<DynamicProperty> dpsList) {
		for (DynamicProperty dp : dpsList) {
			if (dp.getName().equals(VM_KEY_NAME)) {
				_vmName = (String)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_OVERALL_STATUS)) {
				_overallStatus = (ManagedEntityStatus)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_POWERSTATE)) {
				_powerState = (VirtualMachinePowerState)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_SUMMARY)) {
				processVmSummary((VirtualMachineSummary)dp.getVal());
			} else if (dp.getName().equals(VM_KEY_CONFIG_EXTRACONFIG)) {
				processExtraConfig((ArrayOfOptionValue)dp.getVal());
			} else if (dp.getName().equals(VM_KEY_RESOURCECONFIG)) {
                _resourceConfig = (ResourceConfigSpec)dp.getVal();
            } else if (dp.getName().equals(VM_KEY_RESOURCEPOOL)) {
                _resourcePool = (Object)dp.getVal();
            }
		}
	}

	/**
	 * Extract Container information from config.extraConfig
	 * @param ovs
	 */
	private void processExtraConfig(ArrayOfOptionValue ovs) {
		if (ovs != null) {
			for (OptionValue ov : ovs.getOptionValue()) {
				String key = ov.getKey();
				if (EXTRACONFIG_CONTAINER_NAME_KEY.equals(key)) {
					_containerName = (String)ov.getValue();
				} else if (EXTRACONFIG_IMAGE_NAME_KEY.equals(key)) {
					_imageName = (String)ov.getValue();
				} else if (EXTRACONFIG_PORTMAPPING_KEY.equals(key)) {
					_portMapping = (String)ov.getValue();
				}
			}
		}
	}
}
