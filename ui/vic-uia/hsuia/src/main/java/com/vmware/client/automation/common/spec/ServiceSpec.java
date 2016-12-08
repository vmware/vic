/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Class to define service connection details - LDU name, username and password.
 * The ServiceSpec will be linked to a common <code>ManagedEntitySpec</code> instance to specify the LDU
 * and credentials to be used for the API operation.
 * The XXXSrvApi classes will use the ServiceSpec to obtain service connection
 * with respective LDU from the XXXServiceUtil.
 * If no ServiceSpec is linked to the XXXSpec the XXXServiceUtil provides the
 * default service connection specified by the common testbed.
 * TODO: VcServiceUtil to work with the ServiceSpec. Only VcdeServiceUtil is
 * working with the ServiceSpec
 */
public class ServiceSpec extends BaseSpec {
   private static final String NOT_ASSIGNED = "NOT_ASSIGNED";

   /**
    * LDU name or IP
    */
   public DataProperty<String> endpoint;

   /**
    * Username to be used to login to the service.
    */
   public DataProperty<String> username;

   /**
    * Password to be used to login to the service.
    */
   public DataProperty<String> password;

   /**
    * LDU name or IP
    *
    * TODO: Remove it once the auto-merge from the vsphere-2015 branch is deactivated.
    */
   @Deprecated
   public DataProperty<String> lduName;

   /**
    * Username to be used to login to the service.
    *
    * TODO: Remove it once the auto-merge from the vsphere-2015 branch is deactivated.
    */
   @Deprecated
   public DataProperty<String> localUsername;

   /**
    * Password to be used to login to the service.
    *
    * TODO: Remove it once the auto-merge from the vsphere-2015 branch is deactivated.
    */
   @Deprecated
   public DataProperty<String> localPassword;

   /**
    * Username to be used to login to the service.
    *
    * TODO: Remove it once the auto-merge from the vsphere-2015 branch is deactivated.
    */
   @Deprecated
   public DataProperty<String> ssoUsername;

   /**
    * Password to be used to login to the service.
    *
    * TODO: Remove it once the auto-merge from the vsphere-2015 branch is deactivated.
    */
   @Deprecated
   public DataProperty<String> ssoPassword;

   /**
    * Thumbprint.
    *
    * TODO: Remove it once the auto-merge from the vsphere-2015 branch is deactivated.
    */
   @Deprecated
   public DataProperty<String> thumbprint;

   @Override
   public String toString() {

      //TODO: Remove printing of depreccated params once the auto-merge
      // from vsphere-2015 branch is deactivated.
      return String.format(
            "%s: lduName->%s, endpoint->%s, username->%s, password->%s, localUsername->%s, localPassword->%s, ssoUsername->%s, " +
            "ssoPassword->%s, thumbprint->%s",
            this.getClass().getSimpleName(),
            lduName.isAssigned() ? lduName.get() : NOT_ASSIGNED,
            endpoint.isAssigned() ? endpoint.get() : NOT_ASSIGNED,
            username.isAssigned() ? username.get() : NOT_ASSIGNED,
            password.isAssigned() ? password.get() : NOT_ASSIGNED,
            localUsername.isAssigned() ? localUsername.get() : NOT_ASSIGNED,
            localPassword.isAssigned() ? localPassword.get() : NOT_ASSIGNED,
            ssoUsername.isAssigned() ? ssoUsername.get() : NOT_ASSIGNED,
            ssoPassword.isAssigned() ? ssoPassword.get() : NOT_ASSIGNED,
            thumbprint.isAssigned() ? thumbprint.get() : NOT_ASSIGNED);
   }
}
