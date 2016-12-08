/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.navstep;

import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.suitaf.apl.IDGroup;

/**
 * The <code>TreeTabNavStep</code> selects one of the four tabs for tree
 * navigation - Hosts and Clusters, VMs and Templates, Storage, Networking.
 *
 * The action can be used from the home screen only.
 *
 * The <code>TreeTabNavStep</code> should be registered in
 * <code>NGCNavigator</code>.
 */
public class TreeTabNavStep extends BaseNavStep {

   private String _uId;

   /**
    * The constructor initializes the <code>TreeTabNavStep</code> with the name
    * of the entity to be selected and the ID used in the step registration.
    *
    * @param nid
    *           NID of the step
    * @param uid
    *           the unique id of the tab we want to select
    */
   public TreeTabNavStep(String nid, String uid) {
      super(nid);
      _uId = uid;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) throws Exception {
      UI.component.click(IDGroup.toIDGroup(_uId));
   }
}