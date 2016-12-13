package com.vmware.vsphere.client.automation.vicui.common;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

public class VicVcEnvSpec extends EntitySpec {
	public DataProperty<String> vcVersion;
	public DataProperty<String> vchVmName;
	public DataProperty<String> containerVmName;
}
