/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.flexui.componentframework.controls.mx.Label;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Implements the Summary page of a Cluster
 */
public class HostProfilesListView extends BaseView {

   private static final String ID_SELECTED_SET_DATAGRID = "selectedSetDataGrid";
   private static final String ID_HOST_PROFILE_TREE_ITEM = ID_SELECTED_SET_DATAGRID + "/text=%s";

   private static final IDGroup ID_LOADING_PROGRESS_BAR = IDGroup
        .toIDGroup(ID_SELECTED_SET_DATAGRID + "/loadingProgressBar");
   private static final IDGroup AI_EDIT_HOST_PROFILE = IDGroup
         .toIDGroup("vsphere.core.hostprofile.editHostProfileAction");
   private static final IDGroup AI_EXPORT_HOST_PROFILE = IDGroup
               .toIDGroup("vsphere.core.hostprofile.exportAction");
   private static final IDGroup AI_EDIT_HOST_CUSTOMIZATIONS = IDGroup
         .toIDGroup("vsphere.core.hostprofile.editHostCustomizationsAction");
   private static final IDGroup AI_EXPORT_HOST_CUSTOMIZATIONS = IDGroup
         .toIDGroup("vsphere.core.hostprofile.exportHostCustomizationsAction");
   private static final IDGroup AI_REMEDIATE_HOST_PROFILE = IDGroup
         .toIDGroup("vsphere.core.hostprofile.remediateHostsProfileAction");
   private static final IDGroup AI_ATTACH_DETACH_TO_HOST_PROFILE = IDGroup
         .toIDGroup("vsphere.core.hostprofile.attachAction");
   private static final IDGroup AI_CHECK_COMPLIANCE_HOST_PROFILE = IDGroup
         .toIDGroup("vsphere.core.hostProfile.checkComplianceAction");
   private static final IDGroup AI_COPY_STTNGS_TO_HOST_PROFILES = IDGroup
         .toIDGroup("vsphere.core.hostprofile.copyHostProfileSettingsAction");
   private static final String AI_IMPORT_HOST_PROFILE =
         "appBody/dataGridToolbar/vsphere.core.hp."
         + "importProfileActionGlobal/button";

   /**
    * Invokes Import Host Profile from the toolbar.
    */
   public void invokeImportHostProfileFromToolbar() {
      UI.component.click(AI_IMPORT_HOST_PROFILE);
   }
   /**
    * Copy settings to host profiles from actions menu of Host Profile
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void invokeCopySttngsToHostProfilesActionsMenu(String hostProfileName) {
      invokeActionFromContextMenu(hostProfileName, AI_COPY_STTNGS_TO_HOST_PROFILES);
   }

   /**
    * Check host profile compliance from actions menu of Host Profile
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void checkComplianceHostProfileActionsMenu(String hostProfileName) {
      invokeActionFromContextMenu(hostProfileName, AI_CHECK_COMPLIANCE_HOST_PROFILE);
   }

   /**
    * Invoke Edit Host Profile wizard from actions menu of Host Profile
    *
    * @param hostProfileName - the name of the host profile to click
    */
   public void invokeEditSettingsHostProfileActionsMenu(String hostProfileName) {
      ActionNavigator.invokeFromActionsMenu(AI_EDIT_HOST_PROFILE);
   }

   /**
    * The method verifies if a remediate menu item is enabled in the host profile
    * context menu
    *
    * @param hostProfileName - the name of the host profile to click
    */
   public boolean isHostProfileRemediateContextMenuEnabled(String hostProfileName) {
      return isContextMenuActionEnabled(hostProfileName, AI_REMEDIATE_HOST_PROFILE);
   }

   /**
    * Invoke Edit Host Profile wizard from context menu of Host Profile
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void invokeExportHostCustomizationsContextMenu(String hostProfileName) {
      invokeActionFromContextMenu(hostProfileName, AI_EXPORT_HOST_CUSTOMIZATIONS);
   }

   /**
    * Invoke Export Host profile dialog from the context menu of a Host Profile.
    * @param hostProfileName the name of the host profile to invoke on.
    */
   public void invokeExportHostProfileContextMenu(String hostProfileName) {
      _logger.debug(String.format("Invoking Export Host Profile context menu "
                                  + "item for %s", hostProfileName));
      invokeActionFromContextMenu(hostProfileName,AI_EXPORT_HOST_PROFILE);
   }

   /**
    * Invoke Edit Host Profile wizard from context menu of Host Profile
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void invokeEditHostCustomizationsContextMenu(String hostProfileName) {
      invokeActionFromContextMenu(hostProfileName, AI_EDIT_HOST_CUSTOMIZATIONS);
   }

   /**
    * Invoke Remediate wizard from context menu of Host Profile
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void invokeRemediateHostProfileContextMenu(String hostProfileName) {
      invokeActionFromContextMenu(hostProfileName, AI_REMEDIATE_HOST_PROFILE);
   }

   /**
    * Invoke Attach/Detach Hosts/Clusters to Host profile wizard from context menu of Host
    * Profile
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void invokeAttachDetachToHostProfileContextMenu(String hostProfileName) {
      invokeActionFromContextMenu(hostProfileName, AI_ATTACH_DETACH_TO_HOST_PROFILE);
   }

   /**
    * Invoke right click on host profile entity
    *
    * @param hostProfileName - the name of the host profile to select
    */
   public void rightClickOnEntity(String hostProfileName) {
      // TODO move this to SUITA, PR 1405967
      Label tree_item =
            new Label(String.format(ID_HOST_PROFILE_TREE_ITEM, hostProfileName),
                  BrowserUtil.flashSelenium);
      tree_item.rightMouseClick();
   }

   /**
    * Verify whether the Edit Host Customizations menu item is disabled.
    *
    * @return true if the menu item is disabled
    */
   public boolean isEditHostCustomizationsMenuItemEnabled() {
      return ActionNavigator.isMenuActionEnabled(AI_EDIT_HOST_CUSTOMIZATIONS);
   }

   // private methods
   private void invokeActionFromContextMenu(String hostProfileName, IDGroup action) {
      UI.condition.notFound(ID_LOADING_PROGRESS_BAR).await(
         SUITA.Environment.getBackendJobLarge());

      rightClickOnEntity(hostProfileName);

      ActionNavigator.invokeMenuAction(action);
   }

   // private methods
   private boolean isContextMenuActionEnabled(String hostProfileName, IDGroup action) {
      rightClickOnEntity(hostProfileName);

      return ActionNavigator.isMenuActionEnabled(action);
   }
}
