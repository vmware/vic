/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree.spec;

import com.google.common.base.Strings;
import com.vmware.client.automation.components.navigator.Navigator;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.components.tree.EntityTypes;

/**
 * A <code>LocationSpec</code> suitable for modeling tree based navigation as
 * implemented in NGC.
 *
 * Instead of building the path <code>String</code>, navigation steps should be
 * provided as stand alone navigation identifiers. They are internally assembled
 * into a navigation path.
 */
public class TreeLocationSpec extends LocationSpec {

   protected EntityTypes _entityType;
   protected ManagedEntitySpec _entity;

   /**
    * Build a location path based on the provided navigation identifiers (NID
    * constants) located in <code>NGCNavigator</code>.
    */
   public TreeLocationSpec(String entityViewNId, ManagedEntitySpec entitySpec,
         String primaryTabNId, String secondaryTabNId, String tocTabNId) {

      _entity = entitySpec;
      StringBuilder text = new StringBuilder();

      // Append the default home root
      text.append(Navigator.NID_HOME_ROOT);

      if (!Strings.isNullOrEmpty(entityViewNId)) {
         appendPathElement(text, entityViewNId);
      }

      if (entitySpec != null) {
         appendPathElement(text, ENTITY_IDENTIFIER + entitySpec.name.get());
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

      path.set(text.toString());
   }

   /**
    * Build a location path based only on entity view. This will open one of the
    * four tree related tabs and will not select any elements
    *
    * @param entityViewNId
    */
   public TreeLocationSpec(String entityViewNId) {
      this(entityViewNId, null, null, null, null);
   }

   /**
    * Build a location path based on the provided navigation identifiers
    * reaching up to the secondary tab. The entitySpec specifies the entity
    * object to navigate to. The method may be used when id needed to be
    * navigated to requested resources which name is not present at the
    * initialization stage of the test.
    *
    * @param homeViewNId
    * @param entityViewNId
    * @param entitySpec
    *           Entity object to navigate to.
    * @param primaryTabNId
    * @param secondaryTabNId
    * @param tocTabNId
    */
   public TreeLocationSpec(String homeViewNId, String entityViewNId,
         ManagedEntitySpec entitySpec, String primaryTabNId,
         String secondaryTabNId, String tocTabNId) {
      this(entityViewNId, entitySpec, primaryTabNId, secondaryTabNId, tocTabNId);
   }

   /**
    * Append a path element to the location path.
    *
    * @param text
    *           A <code>StringBuilder</code>
    *
    * @param pathElement
    *           A path element
    */
   private void appendPathElement(StringBuilder text, String pathElement) {
      if (text.length() > 0) {
         text.append(PATH_SEPARATOR);
         text.append(pathElement);
      }
   }

   /**
    * Get the tree location spec entity type
    *
    * @return A String representation of the entity type
    */
   public String getEntityType() {
      return _entityType.getEntityType();
   }

   /**
    * Get the ManagedEntitySpec for the entity we are navigating to
    *
    * @return The ManagedEntitySpec for the entity
    */
   public ManagedEntitySpec getEntitySpec() {
      return _entity;
   }
}