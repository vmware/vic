/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Spec defining Virtual Center.
 */
public class VcSpec extends ManagedEntitySpec {

   /**
    * URL of the vSphere Client.
    */
   public DataProperty<String> vscUrl;

   /**
    * SSO user.
    */
   public DataProperty<String> ssoLoginUsername;

   /**
    * SSO password.
    */
   public DataProperty<String> ssoLoginPassword;

   @Deprecated
   /**
    * Use ssoLoginUsername or respective serviceSpec.
    */
   public DataProperty<String> loginUsername;

   @Deprecated
   /**
    * Use ssoLoginPassword or respective serviceSpec.
    */
   public DataProperty<String> loginPassword;


}
