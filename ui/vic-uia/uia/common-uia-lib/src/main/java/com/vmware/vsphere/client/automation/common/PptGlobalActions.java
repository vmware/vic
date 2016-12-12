/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common;

import com.vmware.suitaf.apl.IDGroup;

/**
 * Holder for PPT action ID constants.
 */
public class PptGlobalActions {

   //TODO: I suppose that this should be gone
   // "All Virtual Infrastructure Actions" action menu id for policy tag
   public static final IDGroup AI_TAG_ALL_ACTIONS =
         IDGroup.toIDGroup("afContextMenu.Virtual Infrastructure");

   // "Create Tag" action menu id for policy tag
   public static final IDGroup AI_TAG_CREATE =
         IDGroup.toIDGroup("vsphere.core.uds.addPolicyTagActionGlobal");

   // "Rename" action menu id
   public static final IDGroup AI_RENAME =
         IDGroup.toIDGroup("vsphere.core.uds.renameAction");

   // "Edit Settings" action menu id for policy tag
   public static final IDGroup AI_TAG_EDIT_SETTINGS =
         IDGroup.toIDGroup("vsphere.core.uds.tag.editAction");

   // "Delete" subaction menu id for policy tag
   public static final IDGroup AI_TAG_DELETE =
         IDGroup.toIDGroup("vsphere.core.uds.tag.deleteAction");

   // "Remove from objects" subaction menu id for policy tag
   public static final IDGroup AI_TAG_DEMOTE =
         IDGroup.toIDGroup("vsphere.core.uds.tag.removeAction");

   // "Manage Hosts and Clusters" subaction menu id for policy tag
   public static final IDGroup AI_TAG_MANAGE_HOSTS_AND_CLUSTERS =
         IDGroup.toIDGroup("vsphere.core.policy.editPolicyTagResourcesAction");

   // "Remove from objects" action menu id for policy
   public static final IDGroup AI_POLICY_DEMOTE =
         IDGroup.toIDGroup("vsphere.core.uds.removeActionGlobal");

   // "Delete" action menu id for policy
   public static final IDGroup AI_POLICY_DELETE =
         IDGroup.toIDGroup("vsphere.core.uds.deleteActionGlobal");

   // "Edit Resources" action menu id for policy
   public static final IDGroup AI_POLICY_EDIT_RESOURCES =
         IDGroup.toIDGroup("vsphere.core.uds.configurePolicyResourcesGlobal");

   // "New tag" action menu id
   public static final IDGroup AI_NEW_TAG =
         IDGroup.toIDGroup("vsphere.core.uds.addPolicyTagActionGlobal");

   // "Promote tag" action menu id
   public static final IDGroup AI_POLICY_PROMOTE_TAG =
         IDGroup.toIDGroup("vsphere.core.uds.promoteTags");

   // "Edit" action menu for backing category id
   public static final IDGroup AI_BACKING_CATEGORY_EDIT =
         IDGroup.toIDGroup("vsphere.core.tagging.editCategoryAction");

   // "Delete" action menu for backing category id
   public static final IDGroup AI_BACKING_CATEGORY_DELETE =
         IDGroup.toIDGroup("vsphere.core.tagging.deleteCategoryAction");

   // "Edit" action menu for backing tag id
   public static final IDGroup AI_BACKING_TAG_EDIT =
         IDGroup.toIDGroup("vsphere.core.tagging.editTagAction");

   // "Delete" action menu for backing tag id
   public static final IDGroup AI_BACKING_TAG_DELETE =
         IDGroup.toIDGroup("vsphere.core.tagging.deleteTagAction");

   // "Edit policies" action menu id
   public static final IDGroup AI_VM_EDIT_POLICIES =
         IDGroup.toIDGroup("vsphere.core.vm.editPolicies");
}
