/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed.model;

import java.util.EnumSet;

/**
 * Specifies the different types of available fixtures
 * One Fixtures item contains all the Inventory objects that that Fixtures item
 * is supposed to denote. For example, the NGC_COMMON fixture item contains a
 * Datacenter with a Cluster in it and a Host in the Cluster. It also contains
 * a shared Datastore. All other Fixtures items contain in themselves the
 * NGC_COMMON item and respectively the Inventory it describes. They only build
 * on top of it, i.e. CONTENT_LIBRARY adds on top of NGC_COMMON Inventory a
 * content library managed entity.
 */
@Deprecated
public enum Fixtures {
   NGC_COMMON(),
   CONTENT_LIBRARY(
         FixtureEntities.CONTENT_LIBRARY_LOCAL,
         FixtureEntities.CONTENT_LIBRARY_PUBLISHED,
         FixtureEntities.CONTENT_LIBRARY_SUBSCRIBED
   ),
   VDC(FixtureEntities.VDC_VDC),
   NEW_PROVISIONING(
         FixtureEntities.VDC_VDC,
         FixtureEntities.CONTENT_LIBRARY_LOCAL
      );
   // PBM,
   // LEGACY_PROVISIONING);

   private final EnumSet<FixtureEntities> _fixtureEntities;

   // initializes the items from the common inventory
   private Fixtures() {
      _fixtureEntities = EnumSet.of(
            FixtureEntities.VC,
            FixtureEntities.NGC_COMMON_CLUSTERED_HOST,
            FixtureEntities.NGC_COMMON_CLUSTER,
            FixtureEntities.NGC_COMMON_DATASTORE,
            FixtureEntities.NGC_COMMON_DATACENTER
         );
   }

   /**
    * Initializes inventory items on top of common inventory,
    * i.e. each enum item contains in itself the NGC_COMMON item
    *
    * @param entities
    */
   private Fixtures(FixtureEntities... entities) {
      this();
      for (FixtureEntities fe : entities) {
         _fixtureEntities.add(fe);
      }
   }

   /**
    * Gets all the inventory entities in a Fixture Item
    * @return EnumSet of FixtureEntities contained in the enum item
    */
   public EnumSet<FixtureEntities> getFixtureEntities() {
      return _fixtureEntities;
   }

   /**
    * Checks if the Fixture entity is contained in the Fixtures items
    * @param fixtureEntity - the entity for which to check
    * @return true if it is contained, otherwise false
    */
   public boolean containsFixtureEntity(FixtureEntities fixtureEntity) {
      for (FixtureEntities fe : getFixtureEntities()) {
         if (fe.equals(fixtureEntity)) {
            return true;
         }
      }
      return false;
   }
}
