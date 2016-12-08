/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import java.lang.reflect.Method;
import java.util.Collections;
import java.util.List;

import com.vmware.client.automation.components.control.ObjectSelectorControl3;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.vsphere.client.automation.common.spec.WizardTiwoRestorationSpec;

/**
 * Basic workflow step that verifies wizard navigator error message
 * after restoration from TIWO.
 */
public class VerifyWizardTiwoRestoreErrorStep extends BaseWorkflowStep {

   private WizardTiwoRestorationSpec _restorationSpec;

   @Override
   public void prepare() throws Exception {
      _restorationSpec = getSpec().get(WizardTiwoRestorationSpec.class);
      if (_restorationSpec == null) {
         throw new IllegalArgumentException(
               "Restoration step requires WizardTiwoRestorationSpec"
            );
      }
   }

   @Override
   public void execute() throws Exception {
      WizardNavigator wizardNavigator = new WizardNavigator();

      wizardNavigator.waitForDialogToLoad();

      // check if dialog is opened
      verifyFatal(
            TestScope.FULL,
            executeViewBooleanMethod("isOpen"),
            "Checked that proper wizard view is opened."
         );

      // check if we are on the right wizard page
      verifyFatal(
            TestScope.FULL,
            executeViewBooleanMethod("verifyPageTitle"),
            "Checked the title of the expected wizard view."
         );

      // make sure that if the page has object selector it is fully loaded
      //TODO: remove all the debugging logs later
      _logger.info("before checking ObjSel3's loading progress bar");
      if (ObjectSelectorControl3.loadingProgressBarPresent()) {
         _logger.info("loading progress bar found on ObjSel3 - before wait to finish");
         ObjectSelectorControl3.waitForObjectSelectorToLoad();
         _logger.info("loading progress bar found on ObjSel3 - after wait to finish");
      }
      _logger.info("after checking the loading progress bar of ObjSel3");

      // check if expected warning is present
      if (_restorationSpec.expectedErrorMessage.isAssigned()) {
         List<String> messages = wizardNavigator.getMessagesFromValidationPanel();

         String globalError = wizardNavigator.getGlobalWizardErrorMessage();
         if (globalError != null) {
            messages.add(globalError);
         }

         verifyFatal(
               TestScope.FULL,
               messages.size() > 0,
               "Checking if there are any error messages in the validation panel." +
               String.format(
                     "Expected errors are:\n %s \n ",
                     _restorationSpec.expectedErrorMessage.getAll()
                  )
            );

         // sort both list of messages in order to compare them easily
         Collections.sort(messages);
         List<String> expectedMessages = _restorationSpec.expectedErrorMessage.getAll();
         Collections.sort(expectedMessages);

         // compare already ordered message lists
         verifyFatal(
               TestScope.FULL,
               messages.equals(expectedMessages),
               "Checking if expected error messages are present and correct." +
               String.format(
                     "Expected errors are:\n %s \n "
                     + "Actual errors are:\n %s \n",
                     _restorationSpec.expectedErrorMessage.getAll(),
                     messages
                  )
            );
      }
   }

   // ---------------------------------------------------------------------------
   // Private methods

   private boolean executeViewBooleanMethod(String methodName) {
      boolean result = false;
      Class<? extends WizardNavigator> viewClass = _restorationSpec.expectedView.get();

      if (viewClass == null) {
         _logger.warn("Expected view is not set so skipping opened view check.");
         return true;
      }

      try {
         Method method = viewClass.getMethod(methodName);
         result = (Boolean) method.invoke(null);
      } catch (NoSuchMethodException e) {
         _logger.error(
               String.format(
                     "No %s method was found for class: %s",
                     methodName,
                     viewClass.getName()
                  )
            );
         result =  false;
      } catch (Exception e) {
         _logger.error(
               String.format(
                     "Unable to execute %s method on class: %s",
                     methodName,
                     viewClass.getName()
                  )
            );
         e.printStackTrace();
         result =  false;
      }

      return result;
   }
}
