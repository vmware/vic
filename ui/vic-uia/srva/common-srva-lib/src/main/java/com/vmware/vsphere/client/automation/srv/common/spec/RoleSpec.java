/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * This class will represent a role with privileges
 */
public class RoleSpec extends ManagedEntitySpec {

   /**
    * Used for PermissionsSpec
    * Administrator = -1
    * Anonymous = -4
    * No Access = -5
    * Read-Only = -2
    * View = -3
    */
   public enum UserRole {
       ADMINISTARTOR(-1), ANONYMOUS(-4), NO_ACCESS(-5), READ_ONLY(-2), VIEW(-3);

       private final int _roleId;
       private final RoleSpec _roleSpec;

       UserRole(int value) {
           this._roleId = value;
           this._roleSpec = new RoleSpec();
           this._roleSpec.roleId.set(this._roleId);
       }

       private int getRoleId(){
           return this._roleId;
       }

       public static UserRole getUserRole(int roleId) {
           for (UserRole userRole : values()){
               if (userRole.getRoleId() == roleId){
                   return userRole;
               }
           }

           _logger.warn("There is no role with id:" + roleId);
           return null;
       }

       public RoleSpec getRoleSpec() {
          return _roleSpec;
       }
   }

   /**
    * Role id form enum UserRole, if left empty, then it is a custom role
    * and the UserSrvClass will look for it by the name of the role
    */
    public DataProperty<Integer> roleId;

    /**
     * IDs of the assigned privileges to the role
     */
    public DataProperty<String> privilegeIds;

    /**
     * IDs of the removed from all possible privileges of the role
     */
    public DataProperty<String> exceptPrivilegeIds;

}
