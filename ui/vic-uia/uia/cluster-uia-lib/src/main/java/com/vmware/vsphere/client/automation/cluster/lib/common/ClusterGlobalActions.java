/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.common;

import com.vmware.suitaf.apl.IDGroup;

/**
 * Holder for cluster action ID constants.
 */
public class ClusterGlobalActions {

   // Sub Menus
   public static final IDGroup AI_NEW_VM_SUBMENU =
         IDGroup.toIDGroup("afContextMenu.newVm");
   public static final IDGroup AI_NEW_VAPP_SUBMENU =
         IDGroup.toIDGroup("afContextMenu.newVApp");
   public static final IDGroup AI_STORAGE_SUBMENU =
         IDGroup.toIDGroup("afContextMenu.storage");
   public static final IDGroup AI_HOST_PROFILES_SUBMENU =
         IDGroup.toIDGroup("afContextMenu.hostProfiles");
   public static final IDGroup AI_VDC_SUBMENU =
         IDGroup.toIDGroup("afContextMenu.vdc");

   //---------------------------------------------------------------------------
   // IDs of cluster actions

   // Rename cluster action menu id.
   public static final IDGroup AI_RENAME_CLUSTER =
         IDGroup.toIDGroup("vsphere.core.cluster.renameAction");

   // Delete cluster action menu id.
   public static final IDGroup AI_DELETE_CLUSTER =
         IDGroup.toIDGroup("vsphere.core.cluster.deleteAction");

   // Add host to cluster action menu id.
   public static final IDGroup AI_ADD_HOST_TO_CLUSTER =
         IDGroup.toIDGroup("vsphere.core.host.addAction");

   // Detach cluster from vDC action menu id.
   // NOTE: The menu is visible if the vCDe plugin is presented
   public static final IDGroup AI_DETACH_FROM_VDC =
         IDGroup.toIDGroup("vsphere.core.cluster.detachFromVdcAction");

	// Disable PBM for cluster
	public static final IDGroup AI_DISABLE_PBM_FOR_CLUSTER = IDGroup
			.toIDGroup("vsphere.core.cluster.disablePolicyBasedPlacement");

	// Enable PBM for cluster
	public static final IDGroup AI_ENABLE_PBM_FOR_CLUSTER = IDGroup
			.toIDGroup("vsphere.core.cluster.enablePolicyBasedPlacement");

   // Create cluster action menu id.
   public static final IDGroup AI_CREATE_CLUSTER = IDGroup
         .toIDGroup("vsphere.core.cluster.createAction");
}
