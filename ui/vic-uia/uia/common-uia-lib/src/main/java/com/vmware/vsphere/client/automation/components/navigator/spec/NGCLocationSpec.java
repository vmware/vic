/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.google.common.base.Strings;
import com.vmware.client.automation.components.navigator.Navigator;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A <code>LocationSpec</code> suitable for modeling a standard navigation as
 * implemented in NGC.
 *
 * Instead of building the path <code>String</code>, navigation steps should be
 * provided as stand alone navigation identifier. They are internally assembled
 * into a navigation path.
 */
public class NGCLocationSpec extends LocationSpec {

   // Entity ID place holder in the navigation path string
   private static String ENTITY_ID_HOLDER = "ENTITY_ID_HOLDER";

   // Entity to navigate to
   public DataProperty<ManagedEntitySpec> entity;

   /**
    * Build a location path based on the provided navigation identifiers
    * (NID constants) located in <code>NGCNavigator</code>.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNId,
         String childEntityName,
         String childPrimaryTabNId,
         String childSecondaryTabNId,
         String childTocTabNId) {

      StringBuilder text = new StringBuilder();

      // Append the default home root.
      text.append(Navigator.NID_HOME_ROOT);

      if (!Strings.isNullOrEmpty(homeViewNId)) {
         appendPathElement(text, homeViewNId);
      }

      if (!Strings.isNullOrEmpty(entityViewNId)) {
         appendPathElement(text, entityViewNId);
      }

      if (!Strings.isNullOrEmpty(entityName)) {
         appendPathElement(text, ENTITY_IDENTIFIER + entityName);
      }

      if (!Strings.isNullOrEmpty(primaryTabNId)) {
         appendPathElement(text, primaryTabNId);
      }

      if (!Strings.isNullOrEmpty(secondaryTabNId)) {
         appendPathElement(text, secondaryTabNId);
      }

      if (!Strings.isNullOrEmpty(tocTabNId)) {
         appendPathElement(text, tocTabNId);
      }

      if (!Strings.isNullOrEmpty(childEntityName)) {
         appendPathElement(text, CHILD_ENTITY_IDENTIFIER + childEntityName);
      }

      if (!Strings.isNullOrEmpty(childPrimaryTabNId)) {
         appendPathElement(text, childPrimaryTabNId);
      }

      if (!Strings.isNullOrEmpty(childSecondaryTabNId)) {
         appendPathElement(text, childSecondaryTabNId);
      }

      if (!Strings.isNullOrEmpty(childTocTabNId)) {
         appendPathElement(text, childTocTabNId);
      }

      path.set(text.toString());
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNId,
         String childEntityName,
         String childPrimaryTabNId,
         String childSecondaryTabNId) {
      this(homeViewNId,
            entityViewNId,
            entityName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNId,
            childEntityName,
            childPrimaryTabNId,
            childSecondaryTabNId,
            null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNId,
         String childEntityName,
         String childPrimaryTabNId) {
      this(homeViewNId,
            entityViewNId,
            entityName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNId,
            childEntityName,
            childPrimaryTabNId,
            null, null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNId,
         String childEntityName) {
      this(homeViewNId,
            entityViewNId,
            entityName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNId,
            childEntityName,
            null, null, null);
   }


   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNId) {
      this(homeViewNId,
            entityViewNId,
            entityName,
            primaryTabNId,
            secondaryTabNId,
            tocTabNId,
            null, null, null, null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    * The entitySpec specifies the entity object to navigate to. The method may
    * be used when id needed to be navigated to requested resources which name
    * is not present at the initialization stage of the test.
    * @param homeViewNId
    * @param entityViewNId
    * @param entitySpec     Entity object to navigate to.
    * @param primaryTabNId
    * @param secondaryTabNId
    * @param tocTabNId
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         ManagedEntitySpec entitySpec,
         String primaryTabNId,
         String secondaryTabNId,
         String tocTabNId) {
      this(homeViewNId,
            entityViewNId,
            ENTITY_ID_HOLDER,
            primaryTabNId,
            secondaryTabNId,
            tocTabNId,
            null, null, null, null);

      this.entity.set(entitySpec);
   }


   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId,
         String secondaryTabNId) {
      this(homeViewNId,
            entityViewNId,
            entityName,
            primaryTabNId,
            secondaryTabNId,
            null, null, null, null, null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab.
    */
   public NGCLocationSpec(
         String homeViewNId,
         String entityViewNId,
         String entityName,
         String primaryTabNId) {
      this(homeViewNId, entityViewNId, entityName, primaryTabNId,
           null, null, null, null, null, null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the selected entity.
    */
   public NGCLocationSpec(String homeViewNId, String entityViewNId, String entityName) {
      this(homeViewNId, entityViewNId, entityName,
           null, null, null, null, null, null, null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the entity view
    */
   public NGCLocationSpec(String homeViewNId, String entityViewNId) {
      this(homeViewNId, entityViewNId, null, null, null, null, null, null, null, null);
   }

   /**
    * Set the entity identifier(entity name) at the navigation path if the
    * entity spec is assigned.
    * The entity identifier is used for navigating to a specified entity pages.
    */
   public void populateEntityName() {
      if(entity.isAssigned()) {
         String pathToEdit = path.get();
         pathToEdit = pathToEdit.replace(ENTITY_ID_HOLDER, entity.get().name.get());
         path.set(pathToEdit);
      }
   }

   /**
    * Append a path element to the location path.
    *
    * @param text
    *    A <code>StringBuilder</code>.
    *
    * @param pathElement
    *    A path element.
    */
   private void appendPathElement(StringBuilder text, String pathElement) {
      if (text.length() > 0) {
         text.append(PATH_SEPARATOR);
         text.append(pathElement);
      }
   }

   /**
    * Converts an object's ID to a path ID that Navigator can parse
    *
    * @param id the object id such as parentId/childId
    * @return the converted path ID string
    */
   protected static String convertIdToPathId(String id) {
      if (id == null) {
         return null;
      } else {
         return id.replace(PATH_SEPARATOR, CHILD_ENTITY_GRID_ID_PATH_SEPARATOR).concat(CHILD_ENTITY_GRID_ID_DELIMITER);
      }
   }
}
