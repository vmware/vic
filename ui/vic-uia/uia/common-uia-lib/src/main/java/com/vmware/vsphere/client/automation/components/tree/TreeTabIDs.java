/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * An enumeration for the NIDs or each tree tab we use in tree navigation.
 */
public enum TreeTabIDs {
   HOSTS_AND_CLUSTERS(NGCNavigator.NID_HOME_HOSTS_AND_CLUSTERS_TAB), VMS_AND_TEMPLATES(
         NGCNavigator.NID_HOME_VMS_AND_TEMPLATES_TAB), STORAGE(
         NGCNavigator.NID_HOME_STORAGE_TAB), NETWORKING(
         NGCNavigator.NID_HOME_NETWORKING_TAB);

   private String tabID;

   TreeTabIDs(String id) {
      tabID = id;
   }

   public String getTreeTabID() {
      return tabID;
   }
}