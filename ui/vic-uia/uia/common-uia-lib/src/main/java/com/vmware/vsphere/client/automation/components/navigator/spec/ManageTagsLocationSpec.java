/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * A <code>LocationSpec</code> suitable for modeling the standard navigation
 * to a library > manage > tags view.
 */
public class ManageTagsLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided library navigation
    * identifiers.
    */
   public ManageTagsLocationSpec(String homeViewNid, String entityViewNid,
         String gridId, String libraryName, String primaryTabNId,
         String secondaryTabNId, String libraryItemName,
         String itemPrimaryTabNId, String itemSecondaryTabNId,
         String itemTocTabNid) {
      super(homeViewNid, entityViewNid, libraryName, primaryTabNId,
            secondaryTabNId, null, convertIdToPathId(gridId) + libraryItemName,
            itemPrimaryTabNId, itemSecondaryTabNId, itemTocTabNid);
   }

   public ManageTagsLocationSpec(String entityViewNid, String gridId,
         String libraryName, String primaryTabNId, String secondaryTabNId,
         String libraryItemName, String itemPrimaryTabNId,
         String itemSecondaryTabNId, String itemTocTabNid) {
      super(NGCNavigator.NID_HOME_VCENTER, entityViewNid, libraryName,
            primaryTabNId, secondaryTabNId, null, convertIdToPathId(gridId)
                  + libraryItemName, itemPrimaryTabNId, itemSecondaryTabNId,
            itemTocTabNid);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName,
         String primaryTabNId, String secondaryTabNId, String libraryItemName,
         String itemPrimaryTabNId, String itemSecondaryTabNId,
         String itemTocTabNid) {
      super(NGCNavigator.NID_HOME_LIBRARIES, null,
            libraryName, primaryTabNId, secondaryTabNId, null,
            null, itemPrimaryTabNId,
            itemSecondaryTabNId, itemTocTabNid);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName,
         String primaryTabNId, String secondaryTabNId, String libraryItemName,
         String itemPrimaryTabNId, String itemSecondaryTabNId) {

      this(gridId, libraryName, primaryTabNId, secondaryTabNId,
            libraryItemName, itemPrimaryTabNId, itemSecondaryTabNId, null);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName,
         String primaryTabNId, String secondaryTabNId, String libraryItemName,
         String itemPrimaryTabNId) {
      this(gridId, libraryName, primaryTabNId, secondaryTabNId,
            libraryItemName, itemPrimaryTabNId, null, null);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName,
         String primaryTabNId, String secondaryTabNId, String libraryItemName) {
      this(gridId, libraryName, primaryTabNId, secondaryTabNId,
            libraryItemName, null, null, null);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName,
         String primaryTabNId, String secondaryTabNId) {
      this(gridId, libraryName, primaryTabNId, secondaryTabNId, null, null,
            null, null);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName,
         String primaryTabNId) {
      this(gridId, libraryName, primaryTabNId, null, null, null, null, null);
   }

   public ManageTagsLocationSpec(String gridId, String libraryName) {
      this(gridId, libraryName, null, null, null, null, null, null);
   }

   public ManageTagsLocationSpec(String gridId) {
      this(gridId, null, null, null, null, null, null, null);
   }

   public ManageTagsLocationSpec() {
      this(null, null, null, null, null, null, null, null);
   }
}
