package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * Common step for clicking back button in wizard navigators
 * in order to reach specific view.
 *
 * Operation performed by this step:
 * * Click the back button
 * * Verify expected wizard step title
 */
public class ClickBackAndReach extends ClickButtonAndReach {

   public ClickBackAndReach(Class<? extends WizardNavigator> viewClass) {
      this(viewClass, DEFAULT_VERIFY_METHOD_NAME);
   }

   public ClickBackAndReach(Class<? extends WizardNavigator> viewClass,
         String verifyMethodName) {
      this._viewClass = viewClass;
      this._verifyMethodName = verifyMethodName;
      this._navButtonType = WizardNavigationButton.BACK;
   }
}
