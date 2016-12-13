package com.vmware.vsphere.client.automation.vicui.plugininstalled.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;
import com.vmware.vsphere.client.automation.components.navigator.spec.NGCLocationSpec;

public class AdminNavigationSpec extends NGCLocationSpec {
	
	public AdminNavigationSpec() {
		super(NGCNavigator.NID_HOME_ADMINISTRATION, NGCNavigator.NID_ADMINISTRATION_CLIENT_PLUGINS);
	}
}
