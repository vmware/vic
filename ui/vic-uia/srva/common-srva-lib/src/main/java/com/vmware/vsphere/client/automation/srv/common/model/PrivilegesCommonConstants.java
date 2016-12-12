/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.model;

/**
 * This class represents an enum of privileges. Each privilege is described by
 * its name, privGroupName and privId. Set of privileges create a role. For
 * assigning privileges only privId is used.
 */
public class PrivilegesCommonConstants {

   public enum Privileges {

      SYSTEM_VIEW("View", "System", "System.View"),
      SYSTEM_READ("Read", "System", "System.Read"),
      VAPP_IMPORT("Import", "VApp", "VApp.Import"),
      ALLOCATE_SPACE("AllocateSpace", "Datastore", "Datastore.AllocateSpace"),
      ADD_NEW_DISK("AddNewDisk", "VirtualMachine.Config", "VirtualMachine.Config.AddNewDisk"),
      ASSIGN_NETWORK("Assign", "Network", "Network.Assign"),
      ADD_LIBRARY_ITEM("AddLibraryItem", "ContentLibrary", "ContentLibrary.AddLibraryItem"),
      CREATE_LOCAL_LIBRARY("CreateLocalLibrary", "ContentLibrary", "ContentLibrary.CreateLocalLibrary"),
      CREATE_SUBSCRIBED_LIBRARY("CreateSubscribedLibrary", "ContentLibrary", "ContentLibrary.CreateSubscribedLibrary"),
      DELETE_LIBRARY_ITEM("DeleteLibraryItem", "ContentLibrary", "ContentLibrary.DeleteLibraryItem"),
      DELETE_LIBRARY("DeleteLocalLibrary", "ContentLibrary", "ContentLibrary.DeleteLocalLibrary"),
      DELETE_SUBSCRIBED_LIBRARY("DeleteSubscribedLibrary", "ContentLibrary", "ContentLibrary.DeleteSubscribedLibrary"),
      DOWNLOAD_SESSION("DownloadSession", "ContentLibrary", "ContentLibrary.DownloadSession"),
      EVICT_LIBRARY_ITEM("EvictLibraryItem", "ContentLibrary", "ContentLibrary.EvictLibraryItem"),
      SYNCH_LIBRARY("SyncLibrary", "ContentLibrary", "ContentLibrary.SyncLibrary"),
      UPDATE_SUBSCRIBED_LIBRARY("UpdateSubscribedLibrary", "ContentLibrary", "ContentLibrary.UpdateSubscribedLibrary"),
      UPDATE_LIBRARY("UpdateLibrary", "ContentLibrary", "ContentLibrary.UpdateLibrary"),
      PROBE_LIBRARY_SUBSCRIPTION("ProbeSubscription", "ContentLibrary", "ContentLibrary.ProbeSubscription"),
      EDIT_TAG_CATEGORY("EditCategory", "InventoryService.Tagging", "InventoryService.Tagging.EditCategory"),
      EDIT_TAG("EditTag", "InventoryService.Tagging", "InventoryService.Tagging.EditTag"),
      ASSIGN_UNASSIGN_TAG("AttachTag", "InventoryService.Tagging", "InventoryService.Tagging.AttachTag"),
      CREATE_TAG_CATEGORY("CreateCategory", "InventoryService.Tagging", "InventoryService.Tagging.CreateCategory"),
      CREATE_TAG("CreateTag", "InventoryService.Tagging", "InventoryService.Tagging.CreateTag"),
      DELETE_TAG_CATEGORY("DeleteCategory", "InventoryService.Tagging", "InventoryService.Tagging.DeleteCategory"),
      DELETE_TAG("DeleteTag", "InventoryService.Tagging", "InventoryService.Tagging.DeleteTag"),
      MANAGE_CUSTOM_ATTRIBUTE("ManageCustomFields", "Global", "Global.ManageCustomFields"),
      SET_CUSTOM_ATTRIBUTE("SetCustomField", "Global", "Global.SetCustomField");

      private final String _name;
      private final boolean _onParent = false;
      private final String _privGroupName;
      private final String _privId;

      private Privileges(String name, String privGroupName, String privId) {
         this._name = name;
         this._privGroupName = privGroupName;
         this._privId = privId;
      }

      public String getName() {
         return this._name;
      }

      public String getGroupName() {
         return this._privGroupName;
      }

      public String getId() {
         return this._privId;
      }

      public boolean isOnParent() {
         return _onParent;
      }
   }
}
