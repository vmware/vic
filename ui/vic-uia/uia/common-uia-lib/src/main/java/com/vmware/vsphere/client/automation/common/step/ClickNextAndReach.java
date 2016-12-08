/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import com.vmware.client.automation.components.navigator.WizardNavigator;

/**
 * Common step for clicking next button in wizard navigators
 * in order to reach specific view.
 *
 * Operation performed by this step:
 * 1. Click next button
 * 2. Verify expected wizard step title
 */
public class ClickNextAndReach extends ClickButtonAndReach {

   public ClickNextAndReach(Class<? extends WizardNavigator> viewClass) {
      this(viewClass, DEFAULT_VERIFY_METHOD_NAME);
   }

   public ClickNextAndReach(Class<? extends WizardNavigator> viewClass,
         String verifyMethodName) {
      this._viewClass = viewClass;
      this._verifyMethodName = verifyMethodName;
      this._navButtonType = WizardNavigationButton.NEXT;
   }
}
