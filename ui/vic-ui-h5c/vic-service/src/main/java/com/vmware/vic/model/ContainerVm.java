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

import com.vmware.vic.model.constants.Container;
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
			Container.VM_EXTRACONFIG_CONTAINER_KEY;
	private static final String EXTRACONFIG_IMAGE_NAME_KEY =
			Container.VM_EXTRACONFIG_IMAGE_NAME_KEY;
	private static final String EXTRACONFIG_PORTMAPPING_KEY =
			Container.VM_EXTRACONFIG_PORTMAPPING_KEY;
	private static final String VM_CONTAINERNAME_KEY =
	        Container.VM_CONTAINERNAME_KEY;
	private static final String VM_IMAGENAME_KEY =
	        Container.VM_IMAGENAME_KEY;
	private static final String VM_PORTMAPPING_KEY =
	        Container.VM_PORTMAPPING_KEY;
	private static final String PARENT_OBJ_NAME_KEY =
            Container.PARENT_NAME_KEY;
	private String _containerName = null;
	private String _parentObjectName = null;
	private String _imageName = null;
	private String _portMapping = null;

	public ContainerVm(ObjectContent objContent, String serverGuid) {
		super(objContent, serverGuid);
		processDynamicProperties(objContent.getPropSet());
	}

	/**
	 * Getter for Docker Container's name
	 */
	public String getContainerName() {
		return _containerName;
	}

	/**
	 * Getter for Parent Object's name
	 */
	public String getParentObjectName() {
	    return _parentObjectName;
	}

	/**
	 * Getter for Docker Container's imageName
	 */
	public String getImageName() {
		return _imageName;
	}

	/**
	 * Getter for Docker Container's portMapping
	 */
	public String getPortMapping() {
		return _portMapping;
	}

	/**
	 * Property getter
	 * @param property : property to retrieve
	 */
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
		} else if (VM_CONTAINERNAME_KEY.equals(property)) {
			return _containerName;
		} else if (PARENT_OBJ_NAME_KEY.equals(property)) {
            return _parentObjectName;
        } else if (VM_IMAGENAME_KEY.equals(property)) {
			return _imageName;
		} else if (VM_PORTMAPPING_KEY.equals(property)) {
			return _portMapping;
		} else if (VM_KEY_RESOURCECONFIG.equals(property)) {
            return _resourceConfig;
        } else if (VM_KEY_RESOURCEPOOL.equals(property)) {
            return _resourcePool;
        }
		return UNSUPPORTED_PROPERTY;
	}

	/**
	 * Process DynamicProperty[] and extract information
	 * needed for the ContainerVm model
	 * @param dpsList : DynamicProperty list from ObjectContent.getPropSet()
	 */
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

	/**
	 * Return ManagedObjectReference's value portion
	 * @return the 'value' of the ManagedObjectReference object
	 */
	public String getMorValue() {
	    String[] splitIdString = this.getId().split("/");
	    return splitIdString[1];
	}

	/**
	 * Set _parentObjectName which is the name of this VM's
	 * parent object (VirtualApp or ResourcePool)
	 * @param name
	 */
	public void setParentName(String name) {
	    _parentObjectName = name;
	}
}
