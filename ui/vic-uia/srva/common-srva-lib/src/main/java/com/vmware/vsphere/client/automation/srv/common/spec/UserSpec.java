/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * This class will represent the VC user.
 */
public class UserSpec extends ManagedEntitySpec {

   public DataProperty<String> username;
   public DataProperty<String> password;
   // TODO: add other user properties - user roles, etc.
}
