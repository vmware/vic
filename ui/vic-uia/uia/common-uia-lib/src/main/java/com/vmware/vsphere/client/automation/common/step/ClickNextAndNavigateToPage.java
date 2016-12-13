/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import java.lang.reflect.Method;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;

/**
 * Common step for clicking next button in wizard navigators
 * in order to reach specific view.
 *
 * Operation performed by this step:
 * 1. Click next button
 * 2. Verify expected wizard step title
 * 3. Continue with above steps until view is reached
 */
public class ClickNextAndNavigateToPage extends BaseWorkflowStep {

   private static final String DEFAULT_VERIFY_METHOD_NAME = "verifyPageTitle";
   private static final WizardNavigator _wizardNavigator = new WizardNavigator();

   private final Class<? extends WizardNavigator> _viewClass;
   private final String _verifyMethodName;

   public ClickNextAndNavigateToPage(Class<? extends WizardNavigator> viewClass) {
      this(viewClass, DEFAULT_VERIFY_METHOD_NAME);
   }

   public ClickNextAndNavigateToPage(Class<? extends WizardNavigator> viewClass, String verifyMethodName) {
      this._viewClass = viewClass;
      this._verifyMethodName = verifyMethodName;
   }

   @Override
   public void execute() throws Exception {
      boolean reached = false;

      _logger.info(String.format("Reach: %s view.", _viewClass.getName()));

      verifyFatal(TestScope.FULL, _wizardNavigator.gotoNextPage(),
            "Verifying that next button was successfully clicked.");

      reached = verify();
      String title = _wizardNavigator.getPageTitle();

      while (!reached && _wizardNavigator.isNextBtnEnabled()) {
         verifyFatal(TestScope.FULL, _wizardNavigator.gotoNextPage(),
               "Verifying that next button was successfully clicked.");
         reached = verify();

         if (title.equals(_wizardNavigator.getPageTitle())) {
            _logger.error("Wizard stays on same page.");
            reached = false;
            break;
         }

         title = _wizardNavigator.getPageTitle();
      }

      verifyFatal(TestScope.FULL, reached, "Verifying that page: " + _viewClass.getSimpleName() + " is reached");
   }

   private boolean verify() throws Exception {

      try {
         Method method = _viewClass.getMethod(_verifyMethodName);

         if (method.getReturnType().equals(boolean.class)) {
            return (Boolean) method.invoke(null);
         } else {
            _logger.error(String.format("verifyPageTitle method of class %s doesn't return boolean",
                  _viewClass.getName()));
            return false;
         }
      } catch (NoSuchMethodException e) {
         _logger.error(String.format("No verifyPageTitle method was found for class: %s", _viewClass.getName()));
         return false;
      }
   }
}
