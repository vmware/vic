/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * The class holds a location path for navigating through the standard
 * client navigation model.
 *
 * A standard UI navigation model represents the following navigation
 * elements: root home view, specific home view, entity view, entity,
 * primary tab, secondary tab, and table of content.
 *
 * Example: home.root/home.consumers/entity.orgs/@myorg/
 * entity.l1tab.manage/org.l2tab.settings/org.toc.general
 *
 * Currently the path elements should be exposed through constants in the
 * navigators and their operations should be implemented as navigation steps.
 *
 * In the future the path elements may be replaced with actual client
 * extensions' identifiers.
 */
public class LocationSpec extends BaseSpec{

   public static final String PATH_SEPARATOR = "/";

   public static final String ENTITY_IDENTIFIER = "@";

   public static final String CHILD_ENTITY_IDENTIFIER = "#";

   public static final String CHILD_ENTITY_GRID_ID_DELIMITER = "$";

   public static final String CHILD_ENTITY_GRID_ID_PATH_SEPARATOR = ">";
   /**
    * A path corresponding to the requested target destination the client
    * should be navigated to.
    *
    * Path elements are separated by '/'. Entity names are preceded by '@'.
    *
    * No wildcards are supported. A path cannot be relative to the current
    * location, but should always start from the home root.
    *
    * Example: home.root/home.consumers/entity.orgs/@myorg/
    * entity.l1tab.manage/org.l2tab.settings/org.toc.general
    */
   public DataProperty<String> path;
}