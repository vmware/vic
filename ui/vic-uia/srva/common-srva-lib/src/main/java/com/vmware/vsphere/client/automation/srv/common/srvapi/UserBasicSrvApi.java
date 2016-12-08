/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.ExecutionException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.cis.authz.PermissionTypes.CreateSpec;
import com.vmware.cis.authz.PermissionTypes.Info;
import com.vmware.cis.authz.PermissionTypes.UpdateSpec;
import com.vmware.cis.authz.Principal;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.sso.SsoClient;
import com.vmware.client.automation.util.SsoUtil;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.client.automation.util.VcdeServiceUtil;
import com.vmware.vapi.std.DynamicID;
import com.vmware.vim.binding.sso.PrincipalId;
import com.vmware.vim.binding.sso.admin.PersonDetails;
import com.vmware.vim.binding.sso.admin.PrincipalManagementService;
import com.vmware.vim.binding.vim.AuthorizationManager;
import com.vmware.vim.binding.vim.AuthorizationManager.Permission;
import com.vmware.vim.binding.vim.AuthorizationManager.Role;
import com.vmware.vim.binding.vim.fault.AuthMinimumAdminPermission;
import com.vmware.vim.binding.vim.fault.NotFound;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.binding.vmodl.fault.InvalidArgument;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.PermissionSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.RoleSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;

/**
 * An utility class that provides API operation for managing SSO users and
 * groups.
 */
public class UserBasicSrvApi {

   private static final Logger _logger = LoggerFactory
         .getLogger(UserBasicSrvApi.class);

   private static final String GLOBAL_PERMISSION_ID = "global-permission";
   private static final String GLOBAL_PERMISSION_TYPE = "PermissionFolder";

   /**
    * The <code>DynamicID</code> for the <code>AuthzService</code> root object.
    */
   public static final DynamicID AUTH_SERVICE_ROOT;
   public static final DynamicID AUTH_SERVICE_ENTITY;

   static {
      AUTH_SERVICE_ROOT = new DynamicID();
      AUTH_SERVICE_ROOT.setId(GLOBAL_PERMISSION_ID);
      AUTH_SERVICE_ROOT.setType(GLOBAL_PERMISSION_TYPE);
      AUTH_SERVICE_ENTITY = new DynamicID();
   }

   private static UserBasicSrvApi instance = null;

   protected UserBasicSrvApi() {
   }

   /**
    * Get instance of UserSrvApi.
    *
    * @return created instance
    */
   public static UserBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (UserBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing UserSrvApi.");
               instance = new UserBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates user that belongs to the specified domain.
    *
    * @param userSpec user specification that will be used for user creation
    * @return true if the creation was successful, false otherwise
    * @throws Exception if connection to the SSO cannot be established
    */
   public boolean createUser(UserSpec userSpec) throws Exception {
      validateUserSpec(userSpec);
      _logger.info(String.format("Creating user: '%s', pass: ", userSpec.username.get(), userSpec.password.get()));

      SsoClient ssoConnector = SsoUtil.getVcConnector(userSpec).getConnection();

      PrincipalManagementService principalManagement =
            ssoConnector.getPrincipalManagement();

      // TODO: add first and last name
      String[] usernameParts = userSpec.username.get().split("@");
      PrincipalId principalId = principalManagement.createLocalPersonUser(usernameParts[0], new PersonDetails(),
            userSpec.password.get());

      return principalId != null;
   }

   /**
    * Deletes user.
    *
    * @param userSpec specification of the user that will be deleted
    * @throws Exception if something goes wrong, or VC connectivity fails
    */
   public void deleteUser(UserSpec userSpec) throws Exception {
      validateUserSpec(userSpec);
      _logger.info(String.format("Deleting user '%s'", userSpec.username.get()));

      SsoClient ssoConnector = SsoUtil.getVcConnector(userSpec).getConnection();

      PrincipalManagementService principalManagement =
            ssoConnector.getPrincipalManagement();

      String[] usernameParts = userSpec.username.get().split("@");
      principalManagement.deleteLocalPrincipal(usernameParts[0]);
   }

   /**
    * Adds specific user to several groups.
    * If any of the add operations fails the method will finish.
    *
    * @param userSpec the user that will be added to the administrators groups
    * @param groupNames list of group names to which the user will be added
    * @return true if the operation was successful, false otherwise
    * @throws Exception if communication to the SSO cannot be established
    */
   public boolean addUserToGroups(UserSpec userSpec, List<String> groupNames) throws Exception {
      validateUserSpec(userSpec);
      _logger.info(String.format("Adding user '%s' to the groups: %s", userSpec.username.get(), groupNames));

      String[] usernameParts = userSpec.username.get().split("@");
      PrincipalId principal = new PrincipalId(
            usernameParts[0], usernameParts[1]);

      SsoClient ssoConnector = SsoUtil.getVcConnector(userSpec).getConnection();

      PrincipalManagementService principalManagement =
            ssoConnector.getPrincipalManagement();

      for (String group : groupNames) {
         if (!principalManagement.addUserToLocalGroup(principal, group)) {
            _logger.error(String.format("Cloud not add user %s to group %s", userSpec.username.get(), group));
            return false;
         }
      }

      return true;
   }

   /**
    * Method that creates a specified role and if privileges are assigned, assigns them to the role
    *
    * @param roleSpec
    * @return
    * @throws Exception
    */
   public boolean createRole(RoleSpec roleSpec) throws Exception {
      validateRoleSpec(roleSpec);

      if (!roleSpec.privilegeIds.isAssigned()) {
         throw new IllegalArgumentException(
               "Role privilege ids are not assigned!");
      }

      ServiceSpec serviceSpec = roleSpec.service.get();
      VcService service = VcServiceUtil.getVcService(serviceSpec);

      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      // Authorization manager is used for managing permissions
      AuthorizationManager am = service.getManagedObject(amMor);

      // a role is a set of privileges to assign
      List<String> privilegeIdsList = roleSpec.privilegeIds.getAll();
      String[] privilegeIdsArray = new String[privilegeIdsList.size()];

      // result is the custom role id
      int id = am.addRole(roleSpec.name.get(), roleSpec.privilegeIds.getAll().toArray(privilegeIdsArray));

      return id != 0;
   }

   /**
    * Method that creates a role with all possible privileges except the specified ones.
    *
    * @param roleSpec - should contain the name of the role
    * @param privilegesToRemove - list with the names of the roles to remove
    * @return true if successful, false otherwise
    * @throws Exception
    */
   public boolean createRoleAllExceptSpecifiedPrivileges(RoleSpec roleSpec, String... privilegesToRemove)
         throws Exception {
      validateRoleSpec(roleSpec);

      ServiceSpec serviceSpec = roleSpec.service.get();
      VcService service = VcServiceUtil.getVcService(serviceSpec);

      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      // Authorization manager is used for managing permissions
      AuthorizationManager am = service.getManagedObject(amMor);
      Role adminRole = null;
      Role[] roles = am.getRoleList();
      for (Role role : roles) {
         // -1 is the default id of the Administrator Role
         if (role.getRoleId() == -1) {
            adminRole = role;
            break;
         }
      }

      if (adminRole == null) {
         throw new Exception(
               "Cannot retrieve administrator role!");
      }

      List<String> allPrivileges = new LinkedList<String>(
            Arrays.asList(adminRole.getPrivilege()));
      allPrivileges.removeAll(Arrays.asList(privilegesToRemove));
      // result is the custom role id
      int id = am.addRole(roleSpec.name.get(), allPrivileges.toArray(new String[] {}));

      return id != 0;
   }

   /**
    * Method that deletes a specified role by name
    *
    * @param roleSpec - spec of the role to delete
    * @return true if successful, false otherwise
    * @throws Exception - if cannot connect to VC or role is not found by specified
    *            name in spec
    */
   public boolean deleteRole(RoleSpec roleSpec) throws Exception {
      validateRoleSpec(roleSpec);

      ServiceSpec serviceSpec = roleSpec.service.get();
      VcService service = VcServiceUtil.getVcService(serviceSpec);

      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      // Authorization manager is used for managing permissions
      AuthorizationManager am = service.getManagedObject(amMor);

      // in the API roles are manipulated by their ids
      int roleId = getRoleIdByRoleName(roleSpec.name.get(), am);
      try {
         am.removeRole(roleId, true);
         return true;
      } catch (Exception e) {
         _logger.error(e.getMessage());
         _logger.error(Arrays.asList(e.getStackTrace()).toString());
         return false;
      }

   }

   /**
    * Return true if the role specified by the role spec is found in the system.
    * 
    * @param roleSpec
    * @return
    * @throws VcException
    */
   public boolean isRolePresent(RoleSpec roleSpec) throws VcException {
      validateRoleSpec(roleSpec);

      ServiceSpec serviceSpec = roleSpec.service.get();
      VcService service = VcServiceUtil.getVcService(serviceSpec);
      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      AuthorizationManager am = service.getManagedObject(amMor);

      for (Role role : am.getRoleList()) {
         if (role.getName().equals(roleSpec.name.get())) {
            return true;
         }
      }

      return false;
   }

   /**
    * Return array of privileges assigned to the specified role.
    * 
    * @param roleSpec
    * @return
    * @throws VcException
    */
   public String[] getRolePriviliges(RoleSpec roleSpec) throws VcException {
      validateRoleSpec(roleSpec);

      ServiceSpec serviceSpec = roleSpec.service.get();
      VcService service = VcServiceUtil.getVcService(serviceSpec);
      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      AuthorizationManager am = service.getManagedObject(amMor);
      for (Role role : am.getRoleList()) {
         if (role.getName().equals(roleSpec.name.get())) {
            return role.getPrivilege();
         }
      }

      throw new IllegalArgumentException("Could not find the role '" + roleSpec.name.get() + "'");
   }

   /**
    * Method that sets number of permissions on an entity. The permissions should contain id of a role,
    * which may be also a custom role, or if such is not available as in the case of the custom roles,
    * it should contain the role name
    *
    * @param entitySpec - the entity on which to set the permissions
    * @param permissionSpecs - the permissions that will be set
    * @return true if successful, false otherwise
    * @throws Exception - throws exception if the entity is not present in the Inventory, if connection to VC
    *            cannot be established, if specs are invalid, or the role is inexistent
    */
   public boolean setEntityPermissions(ManagedEntitySpec entitySpec, PermissionSpec... permissionSpecs)
         throws Exception {
      _logger.info("Setting permissions on entity: " + entitySpec.name.get());
      boolean result = true;

      for (PermissionSpec permissionSpec : permissionSpecs) {
         validatePermissionSpec(permissionSpec);
      }

      ServiceSpec serviceSpec = entitySpec.service.get();

      for (PermissionSpec permissionSpec : permissionSpecs) {

         CreateSpec createSpec = new CreateSpec();
         ManagedObjectReference entityMor = ManagedEntityUtil
               .getManagedObjectReference(entitySpec);
         // adding serverGuid as part of the id - without it the permissions
         // won't be assigned to the entity
         AUTH_SERVICE_ENTITY.setId(entityMor.getValue() + ":"
               + entityMor.getServerGuid());
         AUTH_SERVICE_ENTITY.setType(entityMor.getType());
         createSpec.setResourceId(AUTH_SERVICE_ENTITY);

         // username ( the user that is created for the test)
         createSpec.setPrincipal(getPrincipal(permissionSpec));
         // id of the role that will be assigned to the user as global
         // permissions
         int roleId;
         if (!permissionSpec.role.get().roleId.isAssigned()) {
            VcService service = VcServiceUtil.getVcService(serviceSpec);
            ManagedObjectReference amMor = service.getServiceInstanceContent()
                  .getAuthorizationManager();
            // Authorization manager is used for managing permissions
            AuthorizationManager am = service.getManagedObject(amMor);
            roleId = getRoleIdByRoleName(permissionSpec.role.get().name.get(),
                  am);
         } else {
            roleId = permissionSpec.role.get().roleId.get();
         }
         createSpec.setRoleId(Arrays.asList(new String[] { Integer
               .toString(roleId) }));
         createSpec.setPropagate(permissionSpec.propagate.get());
         // applying of permissions
         String create = VcdeServiceUtil.getPermissionManager(serviceSpec)
               .create(createSpec);
         result = result && create != null;
      }

      return result;
   }
   
   /**
    * Edits entity permission by changing the role or if it is propagated to
    * children.
    *
    * @param entitySpec
    *           - the entity which permission will be edited
    * @param permissionName
    *           - the name of the edited permission
    * @param isPropagating
    *           - shows whether the permission is propagating or not
    * @param rolesToUpdate
    *           - the new roles that will be set to the permission
    * @throws Exception
    */
   public void editEntityPermissions(ManagedEntitySpec entitySpec,
         String permissionName, boolean isPropagating,
         RoleSpec... rolesToUpdate) throws Exception {
      _logger.info("Editing permissions on entity: " + entitySpec.name.get());

      ServiceSpec serviceSpec = entitySpec.service.get();
      UpdateSpec updateSpec = new UpdateSpec();
      updateSpec.setPropagate(isPropagating);
      List<String> roleIdsList = new LinkedList<>();
      for (RoleSpec role : rolesToUpdate) {
         VcService service = VcServiceUtil.getVcService(serviceSpec);
         ManagedObjectReference amMor = service.getServiceInstanceContent()
               .getAuthorizationManager();
         // Authorization manager is used for managing permissions
         AuthorizationManager am = service.getManagedObject(amMor);
         int roleId = getRoleIdByRoleName(role.name.get(),
               am);
         roleIdsList.add(Integer.toString(roleId));
      }
      updateSpec.setRoleId(roleIdsList);

      VcdeServiceUtil.getPermissionManager(serviceSpec).update(permissionName,
            updateSpec);
   }

   /**
    * Method that removes permissions from an entity. The permissions should
    * contain name property which represents user or group for which the
    * permission is defined.
    *
    * @param entitySpec
    *           - the entity on which to remove the permissions
    * @param permissionSpecs
    *           - the permissions that will be set
    * @return true if successful, false otherwise
    * @throws Exception
    *            - if the entity is not present in the Inventory
    * @throws VcException
    *            - if connection to VC cannot be established
    * @throws NotFound
    *            - if a permission for this entity and user does not exist
    * @throws AuthMinimumAdminPermission
    *            - if this change would leave the system with no Administrator
    *            permission on the root node
    * @throws InvalidArgument
    *            -if the entity does not support removing permissions
    */
   public boolean removeEntityPermission(ManagedEntitySpec entitySpec,
         PermissionSpec permissions) throws Exception {

      _logger.info("Remove permissions from entity: " + entitySpec.name.get());
      boolean result = false;
      ServiceSpec serviceSpec = entitySpec.service.get();

      ManagedObjectReference entityMor = ManagedEntityUtil
            .getManagedObjectReference(entitySpec);
      // adding serverGuid as part of the id - without it the permissions
      // won't be assigned to the entity
      AUTH_SERVICE_ENTITY.setId(entityMor.getValue() + ":"
            + entityMor.getServerGuid());
      AUTH_SERVICE_ENTITY.setType(entityMor.getType());

      com.vmware.cis.authz.Permission permissionManager = VcdeServiceUtil
            .getPermissionManager(serviceSpec);
      List<Info> listDetails = permissionManager.listDetail();
      for (Info listDetail : listDetails) {
         if (listDetail.getResourceId().equals(AUTH_SERVICE_ENTITY)
               && listDetail.getPrincipal().equals(getPrincipal(permissions))) {
            permissionManager.delete(listDetail.getId());
            result = true;
         }
      }
      return result;

   }

   /**
    * Method that sets global permissions for the user and the role described in the PermissionSpec
    * on the root entity of global permissions
    *
    * @param permissionSpecs - set of permissions to apply
    * @throws Exception - in case of error when connecting to vcenter or inventory service
    */
   public boolean setGlobalPermissions(PermissionSpec... permissionSpecs) throws Exception {
      boolean result = true;

      for (PermissionSpec permissionSpec : permissionSpecs) {
         validatePermissionSpec(permissionSpec);
      }

      ServiceSpec serviceSpec = permissionSpecs[0].service.get();

      // for each permission CreateSpec is created and the permission is created on AUTHZ root
      for (PermissionSpec permissionSpec : permissionSpecs) {
         CreateSpec createSpec = new CreateSpec();
         // root of global permissions
         createSpec.setResourceId(AUTH_SERVICE_ROOT);
         // username ( the user that is created for the test)
         createSpec.setPrincipal(getPrincipal(permissionSpec));
         // id of the role that will be assigned to the user as global permissions
         int roleId;
         if (!permissionSpec.role.get().roleId.isAssigned()) {
            VcService service = VcServiceUtil.getVcService(serviceSpec);
            ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
            // Authorization manager is used for managing permissions
            AuthorizationManager am = service.getManagedObject(amMor);
            roleId = getRoleIdByRoleName(permissionSpec.role.get().name.get(), am);
         } else {
            roleId = permissionSpec.role.get().roleId.get();
         }
         createSpec.setRoleId(Arrays.asList(new String[] { Integer.toString(roleId) }));
         createSpec.setPropagate(permissionSpec.propagate.get());

         // applying of permissions
         result = result && VcdeServiceUtil.getPermissionManager(serviceSpec).create(createSpec) != null;
      }

      return result;
   }

   /**
    * Checks if a role is assigned to a user
    * 
    * @param userSpec - the user
    * @param roleSpec - the role
    * @return true if the role is assigned to the user
    */
   public boolean isRoleInUser(UserSpec userSpec, RoleSpec roleSpec) throws Exception {
      ServiceSpec serviceSpec = userSpec.service.get();

      VcService vcService = VcServiceUtil.getVcService(serviceSpec);
      ManagedObjectReference amMor = vcService.getServiceInstanceContent().getAuthorizationManager();
      AuthorizationManager authorizationManager = vcService.getManagedObject(amMor);

      String permissinId = getGlobalPermissionIdByUserName(serviceSpec, getPrincipal(userSpec));
      List<String> roleIds = VcdeServiceUtil.getPermissionManager(serviceSpec).get(permissinId).getRoleId();

      String expectedRoleId = String.valueOf(getRoleIdByRoleName(roleSpec.name.get(), authorizationManager));
      for (String roleId : roleIds) {
         if (roleId.equals(expectedRoleId)) {
            return true;
         }
      }

      return false;
   }

   /**
    * Method that deletes global permissions for the user and the role described in the PermissionSpec
    * on the root entity of global permissions
    *
    * @param permissionSpecs - set of permissions to delete
    * @throws Exception - in case of error when connecting to vcenter or inventory service
    */
   public void deleteGlobalPermissions(PermissionSpec... permissionSpecs) throws Exception {

      for (PermissionSpec permissionSpec : permissionSpecs) {
         validatePermissionSpec(permissionSpec);
      }

      ServiceSpec serviceSpec = permissionSpecs[0].service.get();

      // for each permission CreateSpec is created and the permission is created on AUTHZ root
      for (PermissionSpec permissionSpec : permissionSpecs) {
         String permissionId = getGlobalPermissionIdByUserName(serviceSpec, getPrincipal(permissionSpec));

         // deleting of permissions
         VcdeServiceUtil.getPermissionManager(serviceSpec).delete(permissionId);
      }
   }

   /**
    * Method that sets permissions on root folder for specified user and role
    *
    * @param permissionSpec
    * @return true if successful, otherwise false
    * @throws Exception
    */
   public boolean setRootFolderPermissions(PermissionSpec... permissionSpecs) throws Exception {

      for (PermissionSpec permissionSpec : permissionSpecs) {
         validatePermissionSpec(permissionSpec);
      }

      ServiceSpec serviceSpec = permissionSpecs[0].service.get();
      VcService service = VcServiceUtil.getVcService(serviceSpec);
      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      AuthorizationManager am = service.getManagedObject(amMor);
      Permission[] permissions = new Permission[permissionSpecs.length];
      int i = 0;

      for (PermissionSpec permissionSpec : permissionSpecs) {
         RoleSpec roleSpec = permissionSpec.role.get();
         Permission permission = new Permission(FolderBasicSrvApi.getInstance().getRootFolder(serviceSpec)._getRef(),
               permissionSpec.name.get(), permissionSpec.group.get(), roleSpec.roleId.get(),
               permissionSpec.propagate.get());
         permissions[i] = permission;
         i++;
      }

      BlockingFuture<Void> bf = new BlockingFuture<Void>();
      am.setEntityPermissions(FolderBasicSrvApi.getInstance().getRootFolder(serviceSpec)._getRef(), permissions, bf);

      try {
         bf.get();
      } catch (ExecutionException e) {
         _logger.error("Error setting permissions!");
         e.printStackTrace();
         return false;
      } catch (InterruptedException e) {
         _logger.error("Error setting permissions!");
         e.printStackTrace();
         return false;
      }

      return true;
   }

   /**
    * Method that gets the permissions of the root folder in vc inventory
    * 
    * @return an array of PermissionSpecs
    * @throws Exception
    */
   public PermissionSpec[] getRootFolderPermissions(ServiceSpec serviceSpec) throws Exception {

      VcService service = VcServiceUtil.getVcService(serviceSpec);
      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      AuthorizationManager am = service.getManagedObject(amMor);

      return getPermissionSpecsOfManagedObject(FolderBasicSrvApi.getInstance().getRootFolder(serviceSpec), am);
   }

   /**
    * Method that gets the permissions of the root folder in vc inventory
    * 
    * @return an array of PermissionSpecs
    * @throws Exception
    */
   public PermissionSpec[] getEntityPermissions(ManagedEntitySpec entitySpec) throws Exception {

      VcService service = VcServiceUtil.getVcService(entitySpec);
      ManagedObject entity = ManagedEntityUtil.getManagedObject(entitySpec, null);
      ManagedObjectReference amMor = service.getServiceInstanceContent().getAuthorizationManager();
      AuthorizationManager am = service.getManagedObject(amMor);

      return getPermissionSpecsOfManagedObject(entity, am);
   }

   // ---------------------------------------------------------------------------
   // Private methods

   private Principal getPrincipal(PermissionSpec permissionSpec) throws Exception {

      Principal principal = new Principal();
      String[] userNameParts = permissionSpec.name.get().split("@");

      // TODO lgrigorova: if username doesn't have domain add default domain
      if (userNameParts.length != 2) {
         throw new IllegalArgumentException(
               "User name has to be in the format <name>@<domain>!");
      }

      String userName = userNameParts[1].toUpperCase() + "\\" + userNameParts[0];
      principal.setUserName(userName);
      Principal.Type type = permissionSpec.group.get() ? Principal.Type.GROUP : Principal.Type.USER;
      principal.setType(type);

      return principal;
   }

   private Principal getPrincipal(UserSpec userSpec) throws Exception {
      Principal principal = new Principal();
      String[] userNameParts = userSpec.name.get().split("@");

      // TODO: if username doesn't have domain add default domain
      if (userNameParts.length != 2) {
         throw new IllegalArgumentException("User name has to be in the format <name>@<domain>!");
      }

      String userName = userNameParts[1].toUpperCase() + "\\" + userNameParts[0];
      principal.setUserName(userName);
      principal.setType(Principal.Type.USER);

      return principal;
   }

   private PermissionSpec[] getPermissionSpecsOfManagedObject(ManagedObject entity, AuthorizationManager am) {
      // reference to the entity, and true for inherited permissions retrieval
      Permission[] permissions = am.retrieveEntityPermissions(entity._getRef(), true);

      PermissionSpec[] permissionSpecs = new PermissionSpec[permissions.length];
      int i = 0;

      for (Permission permission : permissions) {
         PermissionSpec permissionSpec = new PermissionSpec();
         permissionSpec.name.set(permission.getPrincipal());
         permissionSpec.group.set(permission.isGroup());
         RoleSpec roleSpec = new RoleSpec();
         int roleId = permission.getRoleId();
         roleSpec.name.set(getRoleNameByRoleId(roleId, am));
         roleSpec.roleId.set(roleId);
         permissionSpec.role.set(roleSpec);
         permissionSpec.propagate.set(permission.isPropagate());

         permissionSpecs[i] = permissionSpec;
         i++;
      }

      return permissionSpecs;

   }

   private void validateUserSpec(UserSpec userSpec) {
      if (userSpec == null) {
         throw new IllegalArgumentException(
               "User spec is not specified!");
      }

      if (!userSpec.username.isAssigned()) {
         throw new IllegalArgumentException(
               "User name is not assigned!");
      }

      if (!userSpec.password.isAssigned()) {
         throw new IllegalArgumentException(
               "User password is not assigned");
      }

      String[] usernameParts = userSpec.username.get().split("@");

      if (usernameParts.length != 2) {
         throw new IllegalArgumentException(
               "Username is not correct: " + usernameParts);
      }
   }

   private void validateRoleSpec(RoleSpec roleSpec) {
      if (roleSpec == null) {
         throw new IllegalArgumentException(
               "Role spec is not specified!");
      }

      if (!roleSpec.name.isAssigned()) {
         throw new IllegalArgumentException(
               "Role name is not assigned!");
      }

   }

   private void validatePermissionSpec(PermissionSpec permSpec) {
      if (permSpec == null) {
         throw new IllegalArgumentException(
               "Permission Spec is not specified!");
      }

      boolean valid = true;
      valid = valid && permSpec.name.isAssigned();
      valid = valid && permSpec.propagate.isAssigned();
      valid = valid && permSpec.group.isAssigned();
      valid = valid
            && (permSpec.role.isAssigned() && (permSpec.role.get().roleId.isAssigned() || permSpec.role.get().name
                  .isAssigned()));

      if (!valid) {
         throw new IllegalArgumentException(
               "The Permission Spec is not valid!");
      }
   }

   private int getRoleIdByRoleName(String roleName, AuthorizationManager am) {
      Role[] roles = am.getRoleList();
      for (Role role : roles) {
         if (role.getName().equals(roleName)) {
            return role.getRoleId();
         }
      }

      throw new IllegalArgumentException(
            "No such role found: " + roleName);
   }

   private String getRoleNameByRoleId(int roleId, AuthorizationManager am) {
      Role[] roles = am.getRoleList();
      for (Role role : roles) {
         if (role.getRoleId() == roleId) {
            return role.getName();
         }
      }

      throw new IllegalArgumentException(
            "No such role found: " + roleId);
   }

   // Method that gets the global permission id by the username
   private String getGlobalPermissionIdByUserName(ServiceSpec serviceSpec, Principal principal) throws Exception {
      String globalPermissionId = null;

      List<Info> permissionInfoList = VcdeServiceUtil.getPermissionManager(serviceSpec).listDetail();

      for (Info permissionInfo : permissionInfoList) {
         if (permissionInfo.getPrincipal().equals(principal)) {
            globalPermissionId = permissionInfo.getId();
         }
      }

      if (Strings.isNullOrEmpty(globalPermissionId)) {
         throw new IllegalArgumentException(
               "No such permission in Global Permissions");
      }

      return globalPermissionId;
   }

}
