package com.vmware.vsphere.client.automation.common.step;

import java.lang.reflect.Method;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.suitaf.util.AssertionFail;

/**
 * Common step for clicking next button in wizard navigators
 * in order to reach specific view.
 *
 * Operation performed by this step:
 * 1. Click next button
 * 2. Verify expected wizard step title
 */
public abstract class ClickButtonAndReach extends BaseWorkflowStep {

   protected static final String DEFAULT_VERIFY_METHOD_NAME = "verifyPageTitle";
   protected static final WizardNavigator _wizardNavigator = new WizardNavigator();

   protected Class<? extends WizardNavigator> _viewClass;
   protected String _verifyMethodName;
   protected WizardNavigationButton _navButtonType;

   @Override
   public void execute() throws Exception {
      _logger.info(String.format("Reach: %s view.", _viewClass.getName()));

      String currentPageTitle = _wizardNavigator.getPageTitle();

      try {
         switch (_navButtonType) {
            case BACK:
               verifyFatal(TestScope.MINIMAL, _wizardNavigator.gotoPrevPage(),
                     "Verifying that back button was successfully clicked.");
               break;
            case NEXT:
               verifyFatal(TestScope.MINIMAL, _wizardNavigator.gotoNextPage(),
                     "Verifying that next button was successfully clicked.");
               break;
            default:
               // Default is next
               verifyFatal(TestScope.MINIMAL, _wizardNavigator.gotoNextPage(),
                     "Verifying that next button was successfully clicked.");
               break;
         }

      } catch (RuntimeException re) {
         String newPageTitle = _wizardNavigator.getPageTitle();
         if (newPageTitle.equals(currentPageTitle)) {
            // Proceed to verification check if this is the right page to reach
         } else {
            throw re;
         }
      }

      verifyFatal(
            TestScope.UI,
            new TestScopeVerification() {

               @Override
               public boolean verify() throws Exception {

                  try {
                     Method method = _viewClass.getMethod(_verifyMethodName);

                     if (method.getReturnType().equals(boolean.class)) {
                        return (Boolean) method.invoke(_viewClass.newInstance());
                     } else {
                        _logger.error(String.format(
                              "verifyPageTitle method of class %s dosn't return boolean",
                              _viewClass.getName()
                              )
                           );
                        return false;
                     }
                  } catch (NoSuchMethodException e) {
                     _logger.error(
                           String.format(
                                 "No verifyPageTitle method was found for class: %s",
                                 _viewClass.getName()
                              )
                        );
                     return false;

                  }
               }
            },
            String.format("Verifying wizard step %s page title", _viewClass.getName())
         );
   }

   /**
    * Helper class that represents any of the wizard navigation buttons
    * and holds the name of the respective navigation methods.
    */
   protected enum WizardNavigationButton {
      NEXT("gotoNextPage"), BACK("gotoPrevPage");

      private String navigationMethod;

      private WizardNavigationButton(String navigationMethod) {
         this.navigationMethod = navigationMethod;
      }

      /**
       * Gets the name of the navigation method for this button.
       *
       * @return     the name of the nav method
       */
      public String getNavigationMethod() {
         return this.navigationMethod;
      }
   }
}
