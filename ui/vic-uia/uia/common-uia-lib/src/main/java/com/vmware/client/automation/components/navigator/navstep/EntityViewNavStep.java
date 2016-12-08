/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.navstep;

import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.util.CommonUtils;

/**
 * A <code>NavigationStep</code> for selecting an entity view in the object
 * navigator such as Hosts, Clusters, Virtual Machines and etc.
 */
public class EntityViewNavStep extends BaseNavStep {

   private String _uId;

   /**
    * The constructor defines a mapping between the navigation identifier
    * and the UID in the object navigator.
    *
    * @param nId
    *    @see <code>NavigationStep</code>.
    *
    * @param uId
    *    UID of the entity view item in the object navigator.
    */
   public EntityViewNavStep(String nId, String uId) {
      super(nId);

      _uId = uId;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) {
      UI.component.click(IDGroup.toIDGroup(_uId));

      // This is a workaround for waiting the main page to load properly. After the
      // click, the object navigator is loaded quickly with objects, the Refresh
      // button disappears and then the entities grid starts loading which
      // invokes another automatically triggered refresh. Tests should wait for
      // the entities data grid to load => this sleep() hides the obj. navigator
      // status refresh.
      CommonUtils.sleep(5000L);
   }
}
