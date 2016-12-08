/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common;

import com.vmware.suitaf.apl.IDGroup;

/**
 * Class that holds the IDs of the actions in All Actions menu related to VMs.
 */
public class VmGlobalActions {

   // Sub Menus
   public static final IDGroup AI_CLONE_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.clone");
   public static final IDGroup AI_TEMPLATE_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.template");
   public static final IDGroup AI_FAULT_TOLERANCE_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.faultTolerance");
   public static final IDGroup AI_VM_POLICIES_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.policies");
   public static final IDGroup AI_COMPATIBILITY_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.compatibility");
   public static final IDGroup AI_TAGS_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.tagsAndCustomAttributes");
   public static final IDGroup AI_ALARMS_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.alarms");
   public static final IDGroup AI_POWER_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.power");
   public static final IDGroup AI_GUEST_OS_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.guestOs");
   public static final IDGroup AI_SNAPSHOTS_SUBMENU = IDGroup
         .toIDGroup("afContextMenu.snapshots");

   // ---------------------------------------------------------------------------
   // IDs of VM invoke actions
   public static final IDGroup AI_MIGRATE_VM = IDGroup
         .toIDGroup("vsphere.core.vm.migrateAction");

   public static final IDGroup AI_CLONE_VM_TO_VDC = IDGroup
         .toIDGroup("vsphere.core.vm.cloneVmToVdcAction");

   public static final IDGroup AI_CLONE_VM_TO_VM = IDGroup
         .toIDGroup("vsphere.core.vm.provisioning.cloneVmToVmAction");

   public static final IDGroup AI_POLICIES = IDGroup
         .toIDGroup("afContextMenu.policies");

   public static final IDGroup AI_ADD_TO_LIBRARY = IDGroup
         .toIDGroup("vsphere.core.vm.addToLibrary");

   public static final IDGroup AI_EDIT_POLICIES = IDGroup
         .toIDGroup("vsphere.core.vm.editPolicies");

   public static final IDGroup AI_CHECK_COMPLIANCE = IDGroup
         .toIDGroup("vsphere.core.vm.provisioning.checkComplianceAction");

   public static final IDGroup AI_REMEDIATE = IDGroup
         .toIDGroup("vsphere.core.vm.provisioning.remediateAction");

   public static final IDGroup AI_EDIT_SETTINGS = IDGroup
         .toIDGroup("vsphere.core.vm.provisioning.editAction");

   public static final IDGroup AI_EDIT_VM_STORAGE_POLICIES = IDGroup
         .toIDGroup("vsphere.core.pbm.storage.manageVmStorageProfilesAction");

   public static final IDGroup AI_CHECK_VM_STORAGE_POLICIES_COMPLIANCE = IDGroup
         .toIDGroup("vsphere.core.pbm.storage.checkVmRollupComplianceAction");

   public static final IDGroup AI_TURN_ON_FT = IDGroup
         .toIDGroup("vsphere.core.vm.turnOnFt");
   // ---------------------------------------------------------------------------
   // IDs of New VM action
   public static final IDGroup AI_NEW_VM = IDGroup
         .toIDGroup("vsphere.core.vm.provisioning.createVmAction");

   public static final IDGroup AI_POWER_OFF_VM = IDGroup
         .toIDGroup("vsphere.core.vm.powerOffAction");

   public static final IDGroup AI_POWER_ON_VM = IDGroup
         .toIDGroup("vsphere.core.vm.powerOnAction");

}
