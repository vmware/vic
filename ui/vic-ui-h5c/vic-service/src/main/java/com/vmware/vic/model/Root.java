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

public class Root extends ModelObject {
	private final String _uiVersion;
	private int _vchVmsLen;
	private int _containerVmsLen;

	public Root(RootInfo rootInfo, int vchVmsLen, int containerVmsLen) {
		this.setId("vic/vic-root");
		_uiVersion = rootInfo.uiVersion;
		_vchVmsLen = vchVmsLen;
		_containerVmsLen = containerVmsLen;
	}

	public String getName() {
		return "vSphere Integrated Containers";
	}

	public String getUiVersion() {
		return _uiVersion;
	}

	public int getVchVmsLen() {
		return _vchVmsLen;
	}

	public int getContainerVmsLen() {
		return _containerVmsLen;
	}

	@Override
	public Object getProperty(String property) {
		if ("uiVersion".equals(property)) {
			return _uiVersion;
		} else if ("vchVmsLen".equals(property)) {
			return _vchVmsLen;
		} else if ("containerVmsLen".equals(property)) {
			return _containerVmsLen;
		} else if ("id".equals(property)) {
			return this.getId();
		} else if ("name".equals(property)) {
			return "Root";
		}
		return UNSUPPORTED_PROPERTY;
	}

	@Override
	public String toString() {
		return "uiVersion: " + _uiVersion + ", vchVms.length: " + _vchVmsLen + ", containerVms.length: " + _containerVmsLen;
	}
}
