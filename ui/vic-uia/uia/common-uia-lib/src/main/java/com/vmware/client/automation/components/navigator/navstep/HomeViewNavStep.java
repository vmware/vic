/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator.navstep;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.components.navigator.Navigator;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.delay.Delay;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * The <code>HomeViewNavStep</code> is able to select a home root container
 * for groups in the inventories, monitoring, administration, and other
 * categories list on the main home view. The action can be used from the
 * main home view.
 *
 * It is also capable of jumping to the main home
 * view from any location.
 */
public class HomeViewNavStep extends BaseNavStep {

   private static final Logger _logger = LoggerFactory.getLogger(HomeViewNavStep.class);

   private static final String HOME_VIEW_ID = "TreeNodeItem_vsphere.core.navigator.home";

   private String _uId;
   private String _nId;

   /**
    * The constructor initializes the <code>HomeViewNavStep</code> with
    * navigation identifier and UI control identifier that points to a
    * home container group.
    *
    * @param nId
    *    Navigation identifier.
    *
    * @param uId
    *    UI control UID.
    */
   public HomeViewNavStep(String nId, String uId) {
      super(nId);

      _uId = uId;
      _nId = nId;
   }

   @Override
   /**
    * {@inheritDoc}
    * NOTE: The locationSpec parameter is not used. The Home ViewNavigation operate only with the navigation on the
    * Home navigation menu and the location spec is not needed.
    */
   public void doNavigate(LocationSpec locationSpec) {
      // Flag if the step navigates to the home screen.
      boolean isHomeNavigaiton = this._nId.equals(Navigator.NID_HOME_ROOT);
      boolean isValidNavigation = true;
      // Number of times to try to navigate to target location if first navigation operation fail.
      int retryCount = 1;
      // Validation timeout in seconds
      long validationTimeout = Delay.timeout.forSeconds(5).getDuration();
      do {
         if(UI.condition.isFound(IDGroup.toIDGroup(_uId)).estimate()) {
            _logger.info("Click on navigation item: " + _uId);
            UI.component.click(IDGroup.toIDGroup(_uId));
         }

         if(isHomeNavigaiton) {
            _logger.info("Wait is found: " + HOME_VIEW_ID);
            isValidNavigation = UI.condition.isFound(IDGroup.toIDGroup(HOME_VIEW_ID)).await(validationTimeout);
         } else {
            _logger.info("Wait is vanished: " + _uId);
            isValidNavigation = UI.condition.notFound(IDGroup.toIDGroup(_uId)).await(validationTimeout)
                  || UI.component.property.get(Property.CURRENT_STATE, _uId).equals("selected");
         }

         _logger.info("isValidNavigation = " + isValidNavigation);

         // Fail navigation step
         if(!isValidNavigation && retryCount == 0) {
            UI.audit.snapshotAppScreen("FAILED TO NAVIGATE TO " + _nId, "");
            throw new RuntimeException("Failed to navigate to " + _nId
                  + " in the navigation tree!");

         }
         retryCount--;
      } while(!isValidNavigation);

   }
}
