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

import javax.xml.bind.DatatypeConverter;

import com.vmware.vic.model.constants.Vch;
import com.vmware.vim25.ArrayOfOptionValue;
import com.vmware.vim25.DynamicProperty;
import com.vmware.vim25.ManagedEntityStatus;
import com.vmware.vim25.ObjectContent;
import com.vmware.vim25.OptionValue;
import com.vmware.vim25.ResourceConfigSpec;
import com.vmware.vim25.VirtualMachinePowerState;
import com.vmware.vim25.VirtualMachineSummary;

public class VirtualContainerHostVm extends VicBaseVm {
	private static final String EXTRACONFIG_CLIENT_IP_KEY =
			Vch.VM_EXTRACONFIG_CLIENT_IP_KEY;
	private static final String EXTRACONFIG_DOCKER_PERSONALITY_ARGS_KEY =
            Vch.VM_EXTRACONFIG_DOCKER_PERSONALITY_ARGS_KEY;
	private static final String DOCKER_ENGINE_SERVER_TLS_PORT = "2376";
	private static final String VM_KEY_IS_USING_TLS = Vch.VM_IS_USING_TLS;
	private String _clientIp = null;
	private boolean _isUsingTls = true;

	public VirtualContainerHostVm(ObjectContent objContent, String serverGuid) {
		super(objContent, serverGuid);
		processDynamicProperties(objContent.getPropSet());
	}

	/**
	 * Getter for VCH VM's IP
	 */
	public String getClientIp() {
		return _clientIp;
	}

	/**
	 * Getter for isUsingTls
	 */
	public boolean getIsUsingTls() {
	    return _isUsingTls;
	}

	/**
	 * Process DynamicProperty[] and extract information
	 * needed for the VirtualContainerHostVm model
	 * @param dpsList : DynamicProperty list from ObjectContent.getPropSet()
	 */
	@Override
	protected void processDynamicProperties(List<DynamicProperty> dpsList) {
		for (DynamicProperty dp : dpsList) {
			if (dp.getName().equals(VM_KEY_NAME)) {
				_vmName = (String)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_GUESTFULLNAME)) {
			    _guestFullName = (String)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_OVERALL_STATUS)) {
				_overallStatus = (ManagedEntityStatus)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_POWERSTATE)) {
				_powerState = (VirtualMachinePowerState)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_SUMMARY)) {
				processVmSummary((VirtualMachineSummary)dp.getVal());
			} else if (dp.getName().equals(VM_KEY_RESOURCECONFIG)) {
			    _resourceConfig = (ResourceConfigSpec)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_RESOURCEPOOL)) {
			    _resourcePool = (Object)dp.getVal();
			} else if (dp.getName().equals(VM_KEY_CONFIG_EXTRACONFIG)) {
				processExtraConfig((ArrayOfOptionValue)dp.getVal());
			}
		}
	}

	/**
	 * Extract VCH IP from config.extraConfig. Also determine
	 * if VCH appliance uses TLS by looking for a string pattern
	 * "2376" in docker engine server arguments
	 * @param ovs
	 */
	private void processExtraConfig(ArrayOfOptionValue ovs) {
		if (ovs != null) {
			for (OptionValue ov : ovs.getOptionValue()) {
				String key = ov.getKey();
				if (EXTRACONFIG_CLIENT_IP_KEY.equals(key)) {
					byte[] decoded = DatatypeConverter.parseBase64Binary((String)ov.getValue());
					StringBuilder sb = new StringBuilder();
					for (int i = 0; i < decoded.length; i++) {
						sb.append((decoded[i] << 24) >>> 24);
						if (i < decoded.length - 1) {
							sb.append(".");
						}
					}
					_clientIp = sb.toString();
					continue;
				}

				if (EXTRACONFIG_DOCKER_PERSONALITY_ARGS_KEY.equals(key)) {
				    _isUsingTls = ((String)ov.getValue()).indexOf(
				            DOCKER_ENGINE_SERVER_TLS_PORT) > -1;
				    continue;
				}
			}
		}
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
		} else if (VM_KEY_CLIENT_IP.equals(property)) {
			return _clientIp;
		} else if (VM_KEY_OVERALLCPUUSAGE.equals(property)) {
			return _overallCpuUsage;
		} else if (VM_KEY_GUESTMEMORYUSAGE.equals(property)) {
			return _guestMemoryUsage;
		} else if (VM_KEY_COMMITTEDSTORAGE.equals(property)) {
			return _committedStorage;
		} else if (VM_KEY_IS_USING_TLS.equals(property)) {
		    return _isUsingTls;
		} else if (VM_KEY_RESOURCECONFIG.equals(property)) {
		    return _resourceConfig;
		} else if (VM_KEY_RESOURCEPOOL.equals(property)) {
		    return _resourcePool;
		}
		return UNSUPPORTED_PROPERTY;
	}
}
