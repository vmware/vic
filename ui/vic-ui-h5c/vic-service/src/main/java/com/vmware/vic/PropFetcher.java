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
package com.vmware.vic;

import java.security.KeyManagementException;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.xml.ws.BindingProvider;
import javax.xml.ws.handler.MessageContext;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

import com.vmware.utils.ssl.ThumbprintHostNameVerifier;
import com.vmware.utils.ssl.ThumbprintTrustManager;
import com.vmware.vic.model.ContainerVm;
import com.vmware.vic.model.VirtualContainerHostVm;
import com.vmware.vim25.DynamicProperty;
import com.vmware.vim25.InvalidPropertyFaultMsg;
import com.vmware.vim25.ManagedObjectReference;
import com.vmware.vim25.ObjectContent;
import com.vmware.vim25.ObjectSpec;
import com.vmware.vim25.OptionValue;
import com.vmware.vim25.PropertyFilterSpec;
import com.vmware.vim25.PropertySpec;
import com.vmware.vim25.RetrieveOptions;
import com.vmware.vim25.RetrieveResult;
import com.vmware.vim25.RuntimeFaultFaultMsg;
import com.vmware.vim25.ServiceContent;
import com.vmware.vim25.TraversalSpec;
import com.vmware.vim25.VimPortType;
import com.vmware.vim25.VimService;
import com.vmware.vim25.VirtualMachineConfigInfo;
import com.vmware.vise.data.query.PropertyValue;
import com.vmware.vise.data.query.ResultItem;
import com.vmware.vise.security.ClientSessionEndListener;
import com.vmware.vise.usersession.ServerInfo;
import com.vmware.vise.usersession.UserSession;
import com.vmware.vise.usersession.UserSessionService;
import com.vmware.vise.vim.data.VimObjectReferenceService;

public class PropFetcher implements ClientSessionEndListener {
	private static final Log _logger = LogFactory.getLog(PropFetcher.class);
	private static VimPortType _vimPort = initializeVimPort();
	private static final String[] VIC_VM_TYPES = {"isVCH", "isContainer"};
	private static final String EXTRACONFIG_VCH_PATH = "init/common/name";
	private static final String EXTRACONFIG_CONTAINER_PATH = "common/name";
	private static final String SERVICE_INSTANCE = "ServiceInstance";
	private static final Set<String> _thumbprints = new HashSet<String>();
	private final UserSessionService _userSessionService;
	private final VimObjectReferenceService _vimObjectReferenceService;

	private static VimPortType initializeVimPort() {
		VimService vimService = new VimService();
		return vimService.getVimPort();
	}

	static {
		HostnameVerifier hostNameVerifier = new ThumbprintHostNameVerifier();
		HttpsURLConnection.setDefaultHostnameVerifier(hostNameVerifier);

		javax.net.ssl.TrustManager[] tms = new javax.net.ssl.TrustManager[1];
		javax.net.ssl.TrustManager tm = new ThumbprintTrustManager();
		tms[0] = tm;
		javax.net.ssl.SSLContext sc = null;

		try {
			sc = javax.net.ssl.SSLContext.getInstance("SSL");
		} catch (NoSuchAlgorithmException e) {
			_logger.error(e);
		}

		javax.net.ssl.SSLSessionContext sslsc = sc.getServerSessionContext();
		sslsc.setSessionTimeout(0);
		try {
			sc.init(null, tms, null);
		} catch (KeyManagementException e) {
			_logger.error(e);
		}
		javax.net.ssl.HttpsURLConnection.setDefaultSSLSocketFactory(
		sc.getSocketFactory());
	}

	public PropFetcher(
			UserSessionService userSessionService,
			VimObjectReferenceService vimObjectReferenceService) {
		if (userSessionService == null ||
			vimObjectReferenceService == null) {
			throw new IllegalArgumentException("constructor argument cannot be null");
		}
		_userSessionService = userSessionService;
		_vimObjectReferenceService = vimObjectReferenceService;
	}

	/**
	 * Get VMs belonging to a given vApp object reference.
	 * @param objRef
	 * @param isVch
	 * @return ResultItem object containing either VCH VM(s) or Container VM(s)
	 *         based on the isVch boolean value
	 * @throws InvalidPropertyFaultMsg
	 * @throws RuntimeFaultFaultMsg
	 */
	public ResultItem getVmsBelongingToMor(Object objRef, boolean isVch)
			throws InvalidPropertyFaultMsg, RuntimeFaultFaultMsg {
		List<PropertyValue> pvList = new ArrayList<PropertyValue>();
		ResultItem resultItem = new ResultItem();
		resultItem.resourceObject = objRef;

		String entityType = _vimObjectReferenceService.getResourceObjectType(objRef);
		String entityName = _vimObjectReferenceService.getValue(objRef);
		String serverGuid = _vimObjectReferenceService.getServerGuid(objRef);

		ManagedObjectReference mor = new ManagedObjectReference();
		mor.setType(entityType);
		mor.setValue(entityName);

		ServiceContent service = getServiceContent(serverGuid);
		if (service == null) {
			_logger.error("Failed to retrieve ServiceContent!");
			return null;
		}

		ManagedObjectReference viewMgrRef = service.getViewManager();
		List<String> vmList = new ArrayList<String>();
		vmList.add("VirtualMachine");
		ManagedObjectReference cViewRef = _vimPort.createContainerView(
				viewMgrRef,
				mor,
				vmList,
				true);

		PropertySpec propertySpec = new PropertySpec();
		propertySpec.setType("VirtualMachine");
		List<String> pSpecPathSet = propertySpec.getPathSet();
		pSpecPathSet.add("name");
		pSpecPathSet.add("summary");
		pSpecPathSet.add("overallStatus");
		pSpecPathSet.add("runtime.powerState");
		pSpecPathSet.add("config.extraConfig");

		// set the root traversal spec
		TraversalSpec tSpec = new TraversalSpec();
		tSpec.setName("traverseEntities");
		tSpec.setPath("view");
		tSpec.setSkip(false);
		tSpec.setType("ContainerView");

		// set objectspec and attach the root traversal spec
		ObjectSpec objectSpec = new ObjectSpec();
		objectSpec.setObj(cViewRef);
		objectSpec.setSkip(Boolean.TRUE);
		objectSpec.getSelectSet().add(tSpec);

		// set traversal node for VirtualApp->VirtualMachine
		TraversalSpec tSpecVappVm = new TraversalSpec();
		tSpecVappVm.setType("VirtualApp");
		tSpecVappVm.setPath("vm");
		tSpecVappVm.setSkip(false);
		tSpec.getSelectSet().add(tSpecVappVm);

		PropertyFilterSpec propertyFilterSpec = new PropertyFilterSpec();
		propertyFilterSpec.getPropSet().add(propertySpec);
		propertyFilterSpec.getObjectSet().add(objectSpec);

		List<PropertyFilterSpec> propertyFilterSpecs = new ArrayList<PropertyFilterSpec>();
		propertyFilterSpecs.add(propertyFilterSpec);
		RetrieveOptions ro = new RetrieveOptions();
		// TODO: pagination
		//ro.setMaxObjects(arg0);

		RetrieveResult props = _vimPort.retrievePropertiesEx(
				service.getPropertyCollector(),
				propertyFilterSpecs,
				ro);
		if (props != null) {
			for (ObjectContent objC : props.getObjects()) {
				// each managed object reference found will be added to resultItem.properties
				PropertyValue pv = new PropertyValue();
				pv.propertyName = "vm";
				if (isVch) {
					pv.value = new VirtualContainerHostVm(objC, serverGuid);
				} else {
					pv.value = new ContainerVm(objC, serverGuid);
				}

				pvList.add(pv);
			}
		}
		resultItem.properties = pvList.toArray(new PropertyValue[]{});

		return resultItem;
	}

	/**
	 * Compute custom VM properties isContainer and isVCH
	 * @param objRef
	 * @return ResultItem object containing PropertyValue[] for the
	 *         the custom VM properties
	 * @throws InvalidPropertyFaultMsg
	 * @throws RuntimeFaultFaultMsg
	 */
	public ResultItem getVmProperties(Object objRef) throws
			InvalidPropertyFaultMsg, RuntimeFaultFaultMsg {
		ResultItem resultItem = new ResultItem();
		resultItem.resourceObject = objRef;
		String entityType = _vimObjectReferenceService.getResourceObjectType(objRef);
		String entityName = _vimObjectReferenceService.getValue(objRef);
		String serverGuid = _vimObjectReferenceService.getServerGuid(objRef);

		ManagedObjectReference vmMor = new ManagedObjectReference();
		vmMor.setType(entityType);
		vmMor.setValue(entityName);

		VirtualMachineConfigInfo config = null;

		// initialize properties isVCH and isContainer
		PropertyValue pv_is_vch = new PropertyValue();
		pv_is_vch.resourceObject = objRef;
		pv_is_vch.propertyName = VIC_VM_TYPES[0];
		pv_is_vch.value = false;

		PropertyValue pv_is_container = new PropertyValue();
		pv_is_container.resourceObject = objRef;
		pv_is_container.propertyName = VIC_VM_TYPES[1];
		pv_is_container.value = false;

		ServiceContent service = getServiceContent(serverGuid);
		if (service == null) {
			_logger.error("Failed to retrieve ServiceContent!");
			return null;
		}

		PropertySpec propertySpec = new PropertySpec();
		propertySpec.setAll(Boolean.FALSE);
		propertySpec.setType("VirtualMachine");
		propertySpec.getPathSet().add("config");

		ObjectSpec objectSpec = new ObjectSpec();
		objectSpec.setObj(vmMor);
		objectSpec.setSkip(Boolean.FALSE);

		PropertyFilterSpec propertyFilterSpec = new PropertyFilterSpec();
		propertyFilterSpec.getPropSet().add(propertySpec);
		propertyFilterSpec.getObjectSet().add(objectSpec);

		List<PropertyFilterSpec> propertyFilterSpecs = new ArrayList<PropertyFilterSpec>();
		propertyFilterSpecs.add(propertyFilterSpec);

		List<ObjectContent> objectContents = _vimPort.retrieveProperties(service.getPropertyCollector(), propertyFilterSpecs);
		if (objectContents != null) {
			for (ObjectContent content : objectContents) {
				List<DynamicProperty> dps = content.getPropSet();
				if (dps != null) {
					for (DynamicProperty dp : dps) {
						config = (VirtualMachineConfigInfo) dp.getVal();

						List<OptionValue> extraConfigs = config.getExtraConfig();
						for(OptionValue option : extraConfigs) {

							if(option.getKey().equals(EXTRACONFIG_CONTAINER_PATH)) {
								pv_is_container.value = true;
								break;
							}

							if(option.getKey().equals(EXTRACONFIG_VCH_PATH)) {
								pv_is_vch.value = true;
								break;
							}
						}
					}
				}
			}
		}

		resultItem.properties = new PropertyValue[] {pv_is_vch, pv_is_container};

		return resultItem;
	}

	/**
	 * Get ServerInfo with the given serverGuid
	 * @param serverGuid
	 * @return ServerInfo object corresponding to the specified serverGuid
	 */
	private ServerInfo getServerInfoObject(String serverGuid) {
		UserSession userSession = _userSessionService.getUserSession();

		for (ServerInfo sinfo : userSession.serversInfo) {
			if (sinfo.serviceGuid.equalsIgnoreCase(serverGuid)) {
				return sinfo;
			}
		}
		return null;
	}

	/**
	 * Set thumbprint from the ServerInfo object
	 * @param sinfo
	 */
	private void setThumbprint(ServerInfo sinfo) {
		String thumbprint = sinfo.thumbprint;
		if (thumbprint != null) {
			_thumbprints.add(thumbprint.replaceAll(":", "").toLowerCase());
		}
		ThumbprintTrustManager.setThumbprints(_thumbprints);
	}

	/**
	 * Get ServerContent object with the given serverGuid
	 * @param serverGuid
	 * @return ServiceContent object corresponding to the specified serverGuid
	 */
	private ServiceContent getServiceContent(String serverGuid) {
		ServerInfo serverInfoObject = getServerInfoObject(serverGuid);
		setThumbprint(serverInfoObject);
		String sessionCookie = serverInfoObject.sessionCookie;
		String serviceUrl = serverInfoObject.serviceUrl;

		List<String> values = new ArrayList<String>();
		values.add("vmware_soap_session=" + sessionCookie);
		Map<String, List<String>> reqHeadrs = new HashMap<String, List<String>>();
		reqHeadrs.put("Cookie", values);

		Map<String, Object> reqContext = ((BindingProvider) _vimPort).getRequestContext();
		reqContext.put(BindingProvider.ENDPOINT_ADDRESS_PROPERTY, serviceUrl);
		reqContext.put(BindingProvider.SESSION_MAINTAIN_PROPERTY, true);
		reqContext.put(MessageContext.HTTP_REQUEST_HEADERS, reqHeadrs);

		final ManagedObjectReference svgInstanceRef = new ManagedObjectReference();
		svgInstanceRef.setType(SERVICE_INSTANCE);
		svgInstanceRef.setValue(SERVICE_INSTANCE);

		ServiceContent serviceContent = null;
		try {
			serviceContent = _vimPort.retrieveServiceContent(svgInstanceRef);
		} catch (RuntimeFaultFaultMsg e) {
			_logger.error("getServiceContent error: " + e);
		}

		return serviceContent;
	}

	@Override
	public void sessionEnded(String clientId) {
		_logger.info("Logging out client session - " + clientId);
	}

}
