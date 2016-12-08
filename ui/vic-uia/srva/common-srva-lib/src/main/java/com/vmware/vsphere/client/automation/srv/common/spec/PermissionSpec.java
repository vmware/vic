/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * This class will represent a user with permissions.
 */
public class PermissionSpec extends ManagedEntitySpec {

   /**
    * True if this represents permissions of a group,
    * and false if this represents the permissions of
    * a user
    */
   public DataProperty<Boolean> group;

   /**
    * propagate permission to children
    */
   public DataProperty<Boolean> propagate;

    /**
     * Role assigned to the permission
     */
    public DataProperty<RoleSpec> role;

}
