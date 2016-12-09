package com.vmware.vsphere.client.automation.common.step;

import java.lang.reflect.Method;
import java.util.Arrays;
import java.util.List;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.TestScopeVerification;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vise.store.InvalidArgumentException;

/**
 * Common step that goes through multiple wizard steps by clicking next button.
 * No changes in the wizard steps are done - default settings are left.
 * This step expects the wizard to be opened.
 * The first view class that has to be passed is reached after first click on Next button.
 *
 * The expected verify title method name is 'verifyPageTitle'.
 * If the view does not have such method defined the execution will fail.
 *
 * Operation performed by this step:
 *  * Click next button
 *  * Verify expected wizard step title
 *  * Repeat for every passed wizard step view
 */
public class GoThroughWizardStepsWithDefaultSettings extends BaseWorkflowStep {

   private static final String DEFAULT_VERIFY_METHOD_NAME = "verifyPageTitle";

   private final List<Class<? extends WizardNavigator>> _viewClasses;

   public GoThroughWizardStepsWithDefaultSettings(
         Class<? extends WizardNavigator>... viewClasses) {
      this._viewClasses = Arrays.asList(viewClasses);
   }

   @Override
   public void prepare() throws Exception {
      if (CollectionUtils.isEmpty(_viewClasses)) {
         throw new InvalidArgumentException(
               "No views were passed to GoThroughWizardStepsWithDefaultSettings test step"
            );
      }
   }

   @Override
   public void execute() throws Exception {

      for (final Class<? extends WizardNavigator> viewClass : _viewClasses) {
         verifyFatal(
               TestScope.FULL,
               new WizardNavigator().gotoNextPage(),
               "Verifying that next button was successfully clicked."
            );

         verifyFatal(
               TestScope.UI,
               new TestScopeVerification() {

                  @Override
                  public boolean verify() throws Exception {

                     try {
                        Method method = viewClass.getMethod(DEFAULT_VERIFY_METHOD_NAME);

                        if (method.getReturnType().equals(boolean.class)) {
                           return (Boolean) method.invoke(null);
                        } else {
                           _logger.error(String.format(
                                "verifyPageTitle method of class %s doesn't return boolean",
                                 viewClass.getName())
                              );
                           return false;
                        }
                     } catch (NoSuchMethodException e) {
                        _logger.error(
                              String.format(
                                    "No verifyPageTitle method was found for class: %s",
                                    viewClass.getName()
                                 )
                           );
                        return false;
                     }
                  }
               },
               String.format("Verifying wizard step %s page title", viewClass.getName())
            );
      }
   }

   // TestWorkflowStep  methods

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      if (CollectionUtils.isEmpty(_viewClasses)) {
         throw new InvalidArgumentException(
               "No views were passed to GoThroughWizardStepsWithDefaultSettings test step"
            );
      }
   }
}
