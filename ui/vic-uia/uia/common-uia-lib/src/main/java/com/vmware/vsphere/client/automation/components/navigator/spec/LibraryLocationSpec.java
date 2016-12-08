/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * A <code>LocationSpec</code> suitable for modeling the standard navigation in
 * the Content Library related tests.
 */
public class LibraryLocationSpec extends NGCLocationSpec {

   /**
    * Build a location path based on the provided library navigation identifiers.
    */
   public LibraryLocationSpec(
         String libraryName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNid) {

      super(
            NGCNavigator.NID_HOME_LIBRARIES,
            null,
            libraryName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNid);
   }

   /**
    * Build a location path based on the provided library navigation identifiers.
    */
   public LibraryLocationSpec(
           String libraryName,
           String primaryTabNId,
           String secondaryTabNId) {
      this(libraryName, primaryTabNId, secondaryTabNId, null);
   }

   /**
    * Build a location path based on the provided library navigation identifiers.
    */
   public LibraryLocationSpec(String libraryName, String primaryTabNId) {
      this(libraryName, primaryTabNId, null, null);
   }

   /**
    * Build a location path based on the provided library navigation identifiers.
    */
   public LibraryLocationSpec(String libraryName) {
      this(libraryName, null, null, null);
   }

   /**
    * Build a location that will navigate the UI to the library entity view.
    */
   public LibraryLocationSpec() {
      this(null, null, null, null);
   }
}



