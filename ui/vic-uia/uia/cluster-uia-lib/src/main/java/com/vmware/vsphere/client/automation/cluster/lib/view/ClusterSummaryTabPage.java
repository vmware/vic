/**
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.cluster.lib.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Represents the elements and supported actions with them on
 * the following location: Cluster -> Summary
 */
public class ClusterSummaryTabPage extends BaseView{
   private static final IDGroup ID_PORTLET_CLUSTER_VSPHERE_DRS = IDGroup
         .toIDGroup("vsphere.core.cluster.summary.drsDetailsView.chrome");
   private static final IDGroup ID_PORTLET_CLUSTER_VSPHERE_HA = IDGroup
         .toIDGroup("vsphere.core.cluster.summary.haPortletView.chrome");
   private static final IDGroup ID_VSPHERE_DRS_BADGE = IDGroup
         .toIDGroup("drsBadge");
   private static final IDGroup ID_VSPHERE_HA_BADGE = IDGroup
         .toIDGroup("haBadge");

   /**
    * Checks whether the DRS portlet is visible
    *
    * @return true if it is visible
    */
   public boolean isRunDrsPortletVisible() {
      waitForPageToRefresh();
      return UI.condition.isFound(ID_PORTLET_CLUSTER_VSPHERE_DRS).estimate();
   }

   /**
    * Checks whether the DRS badge is visible
    *
    * @return true if it is visible
    */
   public boolean isDrsBadgeVisible() {
      waitForPageToRefresh();
      return UI.condition.isFound(ID_VSPHERE_DRS_BADGE).estimate();
   }

   /**
    * Checks whether the HA portlet is visible
    *
    * @return true if it is visible
    */
   public boolean isHaPortletVisible() {
      waitForPageToRefresh();
      return UI.condition.isFound(ID_PORTLET_CLUSTER_VSPHERE_HA).estimate();
   }

   /**
    * Checks whether the HA badge is visible
    *
    * @return true if it is visible
    */
   public boolean isHaBadgeVisible() {
      waitForPageToRefresh();
      return UI.condition.isFound(ID_VSPHERE_HA_BADGE).estimate();
   }
}
