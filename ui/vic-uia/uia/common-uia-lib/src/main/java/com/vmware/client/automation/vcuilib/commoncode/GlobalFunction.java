/**
 * Copyright 2012 VMWare, Inc. All rights reserved. -- VMWare Confidential
 */

package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_MANAGE_MAIN_TAB_CONTAINER;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_MONITOR_MAIN_TAB_CONTAINER;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_NO_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_RECENT_TASKS_ICON_ON_SIDE_BAR;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_RELATED_ITEMS_MAIN_TAB_CONTAINER;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_TASKS_LIST;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_YES_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_YES_NO_DIALOG;
import static com.vmware.client.automation.vcuilib.commoncode.TestBaseUI.captureSnapShot;
import static com.vmware.client.automation.vcuilib.commoncode.TestBaseUI.verifySafely;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_ONE_MINUTE;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.MANAGE_TAB;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.MONITOR_TAB;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.RELATED_ITEMS_TAB;
import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import java.util.Arrays;
import java.util.Date;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.TASK_STATE;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.Timeout;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.Alert;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.componentframework.controls.mx.custom.RecentTasksList;
import com.vmware.flexui.componentframework.controls.mx.custom.ViClientPermanentTabBar;
import com.vmware.flexui.selenium.MethodCallUtil;

/**
 * Class with general utility methods in the UI.
 *
 * NOTE: this class is a partial copy of the one from VCUI-QE-LIB
 */
public class GlobalFunction {

   private static final Logger logger = LoggerFactory.getLogger(GlobalFunction.class);

   /**
    * Returning the current method name, where it is calling from
    *
    * @return the current method name, where it is calling from
    */
   public static String getMethodName() {
      return Thread.currentThread().getStackTrace()[2].getMethodName();
   }

   /**
    * Returns the main tab container id value for a given main tab name
    *
    * @param mainTabName
    * @return
    */
   public static String convertMainTabNameToExtensionId(String mainTabName) {
      String mainTabExtId = null;
      if (mainTabName.equals(MONITOR_TAB)) {
         mainTabExtId = ID_MONITOR_MAIN_TAB_CONTAINER;
      } else if (mainTabName.equals(MANAGE_TAB)) {
         mainTabExtId = ID_MANAGE_MAIN_TAB_CONTAINER;
      } else if (mainTabName.equals(RELATED_ITEMS_TAB)) {
         mainTabExtId = ID_RELATED_ITEMS_MAIN_TAB_CONTAINER;
      }
      return mainTabExtId;
   }

   /**
    * This method will allow the tab navigation through the index
    *
    * @param index Index of tab to navigate
    * @param flashSelenium flashSelenium of the application
    */
   public static void tabNavigate(String index, FlashSelenium flashSelenium) {
      tabNavigate(IDConstants.ID_TAB_NAVIGATOR, index, flashSelenium);
   }

   /**
    * This method will allow the tab navigation through the index
    *
    * @param index Index of tab to navigate
    * @param flashSelenium flashSelenium of the application
    * @param tabID ID of the tab
    */
   public static void tabNavigate(String tabID, String index, FlashSelenium flashSelenium) {
      logger.info("Entering method: " + GlobalFunction.getMethodName());
      ViClientPermanentTabBar tabNavigation =
            new ViClientPermanentTabBar(tabID, flashSelenium);
      tabNavigation.waitForElementEnable(DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE * 10);
      tabNavigation.selectTabAt(index, DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE);

      try {
         int retry = 0;
         String selectedIndex = tabNavigation.getProperty("selectedIndex");
         while (!selectedIndex.equals(index) && retry < 10) {
            logger.info("Selected index is actually: " + selectedIndex
                  + " and expected index is: " + index);
            logger.info("Retry #" + retry);
            Thread.sleep(TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE);
            tabNavigation.selectTabAt(index, DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE);
            selectedIndex = tabNavigation.getProperty("selectedIndex");
            retry++;
         }
         if (!selectedIndex.equals(index) && retry == 10) {
            logger.info("Timed out!");
            TestBaseUI.captureSnapShot("tabNavigate_failures");
            throw new Exception("Timed out");
         } else {
            logger.info("Tab navigation is successful & tabIndex is " + index);
            tabNavigation
                  .waitForElementEnable(DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE * 10);
         }
      } catch (Exception e) {
         logger.error("Tab Navigation is not successful");
         e.printStackTrace();
      }
   }

   /**
    * @Deprecated by TasksUtil.java::waitForRecentTaskToMatchFilter(RecentTaskFilter)
    *
    * This method is waiting for task to reach the expected state How this
    * method works is firstly it tries to find the task name, which is the same
    * as the expected one Opens the Recent tasks list if it's not opened If
    * expected state is "running": - It will match if the found state is
    * "running" - It will match also if the found state is "failed" or
    * "success", but with a warning If the expected state is "success": - It
    * will match if the found state is "success" - It will match also if the
    * found state is "failed", but with a warning - But it will wait until the
    * time out or "success" or "failed" status is reached, when status is
    * "running" If the expected state is "failed": - It will match if the found
    * state is "success", but with a error warning - It will match if the found
    * state is "failed" - But it will wait until the time out or "success" or
    * "failed" status is reached, when status is "running"
    *
    * @param taskName The task name to be waited
    * @param targetName The entity name as a target of a task
    * @param taskState The task state that the waited task is expected
    * @param timeout The timeout of the entire wait process if it is not reached
    *           the expected state or error
    * @param waitInterval The interval between each check
    * @param flashSelenium flashSelenium of the current application browser
    * @param taskStarteAfter It will only look at the task that have been
    *           started after this date
    */
   @Deprecated
   public static void waitUntilTasksState(String taskName, String targetName,
         TASK_STATE taskState, long timeout, long waitInterval,
         FlashSelenium flashSelenium, long... taskStartedAfter) {
      logger.info("Entering method: " + GlobalFunction.getMethodName());
      String[] tasks = null;
      RecentTasks expectedTask = null;
      long startTime = System.currentTimeMillis();
      long currentTimeElapsed = 0;
      logger.info("Task State in string is " + taskState.toString());
      logger.info("Task Name is " + taskName);
      logger.info("Target Name is " + targetName);
      boolean found = false;

      if (MethodCallUtil.getVisibleOnPath(
            flashSelenium,
            IDConstants.ID_RECENT_TASKS_PORTLET)) {
         logger.info("Recent Tasks portlet already exists, so there is no need to invoke...");
      } else {
         logger.warn("Recent Tasks portlet doesn't exist, invoking RecentTasks portlet...");
         Button recentTasksButton =
               new Button(IDConstants.ID_RECENT_TASKS_ICON_ON_SIDE_BAR, flashSelenium);
         // If appSidebar is collapsed click on the Recent tasks icon
         if (MethodCallUtil.getVisibleOnPath(
               flashSelenium,
               ID_RECENT_TASKS_ICON_ON_SIDE_BAR)) {
            recentTasksButton.click(DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE);
         }
      }

      int refreshUIAfterAttempts = 0;
      while (found == false && currentTimeElapsed < timeout) {
         tasks = getAllTasks(flashSelenium);
         logger.info("tasks: " + Arrays.toString(tasks));
         for (int i = 0; i < tasks.length; i++) {
            if (tasks[i].contains(taskName)
                  && tasks[i].split(" /// ")[1].equals(targetName)) {

               // if we are only need to look at tasks that started after a
               // certain time
               if (taskStartedAfter != null && taskStartedAfter.length > 0) {
                  Date afterDate = new Date(taskStartedAfter[0]);
                  logger.info("We are looking for tasks that started after:" + afterDate);
                  RecentTasks foundTask = RecentTasks.fromString(tasks[i]);
                  if (foundTask.getTaskTime() < taskStartedAfter[0]) {
                     logger.info("Found the task but it started before:" + afterDate
                           + ". It started at: " + new Date(foundTask.getTaskTime())
                           + ". Discarding it.");
                     continue;
                  }
               } else {
                  logger.info("not looking for tasks that start after some time");
               }

               logger.info("Found it! Task no: " + (i + 1));
               expectedTask = RecentTasks.fromString(tasks[i]);
               break;
            }
         }

         if (expectedTask != null) {
            String taskStatus = expectedTask.getStatus();
            logger.info("Found task status is :" + taskStatus);
            switch (taskState) {
               case TASK_RUNNING:
                  if ((taskStatus.equals(TestConstants.TASK_STATUS_IN_PROGRESS))
                        || (taskStatus.equals(TestConstants.TASK_STATUS_RUNNING))) {
                     found = true;
                  } else if (taskStatus.equals(TestConstants.TASK_STATUS_SUCCESS)
                        || taskStatus.equals(TestConstants.TASK_STATUS_FAILED)) {
                     throw new AssertionError("It finds the task, but with status "
                           + taskStatus + " instead of " + taskState.toString());
                  }
                  break;
               case TASK_FAILED:
                  if (taskStatus.equals(TestConstants.TASK_STATUS_SUCCESS)) {
                     throw new AssertionError(
                           "Found unexpected status of the task, expected - "
                                 + taskState.toString() + " found - " + taskStatus);
                  } else if (taskStatus.equals(TestConstants.TASK_STATUS_FAILED)) {
                     found = true;
                  } else if (taskStatus.equals(TestConstants.TASK_STATUS_IN_PROGRESS)) {
                     // Go to the next loop
                  }
                  break;
               case TASK_COMPLETED:
                  if (taskStatus.equals(TestConstants.TASK_STATUS_FAILED)) {
                     throw new AssertionError(
                           "Found unexpected status of the task, expected - "
                                 + taskState.toString() + " found - " + taskStatus);
                  } else if (taskStatus.equals(TestConstants.TASK_STATUS_SUCCESS)) {
                     found = true;
                  } else if (taskStatus.equals(TestConstants.TASK_STATUS_IN_PROGRESS)) {
                     // Go to the next loop
                     logger.info("Still In Progress");
                  } else if (taskStatus.equals(TestConstants.TASK_STATUS_RUNNING)) {
                     // Go to the next loop
                     logger.info("Still Running");
                  }
                  break;
            }
         }
         if (found) {
            break;
         } else {
            refreshUIAfterAttempts++;
            if (refreshUIAfterAttempts > 5) {
               refreshUI();
            }
            try {
               Thread.sleep(waitInterval);
            } catch (Exception e) {
               e.printStackTrace();
            }
            currentTimeElapsed = System.currentTimeMillis() - startTime;
         }
      }
      if (found == false) {
         captureSnapShot("Task_Not_Found");
         throw new AssertionError("Unable to find the task " + taskName
               + " with the given task state " + taskState);
      }
   }

   /**
    * @Deprecated by TasksUtil.java::waitForRecentTaskToMatchFilter(RecentTaskFilter)
    *
    * Overloaded method of waitUntilTasksState (without flashSelenium parameter
    * input). Thus by default, it will be using the flashSelenium from
    * BrowserUtil
    *
    * @param taskName The task name to be waited
    * @param targetName The entity name as a target of a task
    * @param taskState The task state that the waited task is expected
    * @param timeout The timeout of the entire wait process if it is not reached
    *           the expected state or error
    * @param waitInterval The interval between each check
    * @param taskStarteAfter It will only look at the task that have been
    *           started after this date
    */
   @Deprecated
   public static void waitUntilTasksState(String taskName, String targetName,
         TASK_STATE taskState, long timeout, long waitInterval, long... taskStartedAfter) {
      GlobalFunction.waitUntilTasksState(
            taskName,
            targetName,
            taskState,
            timeout,
            waitInterval,
            flashSelenium,
            taskStartedAfter);
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    * This is a wrapper call to getAllTasks from RecentTasks portlet to be
    * passed through seleniumflexapi
    *
    * @param flashSelenium flashSelenium of the current application browser
    * @return String[] tasks
    */
   @Deprecated
   public static String[] getAllTasks(FlashSelenium flashSelenium) {
      RecentTasksList recentTasksList =
            new RecentTasksList(IDConstants.ID_TASKS_GRID, flashSelenium);
      if (!recentTasksList.isVisibleOnPath()) {
         recentTasksList = new RecentTasksList(ID_TASKS_LIST, flashSelenium);
      }
      String[] tasks = recentTasksList.getRecentTasks();
      if(tasks == null) {
         tasks = new String[0];
      }
      return tasks;
   }

   /**
    * Wait for a component to be populated and visible till specified time.
    *
    * @param String - Id of the component to be identified on the UI.
    * @param int - retry count.
    */
   public static void waitBeforeComponentVisible(String componentId, int retryNo)
         throws Exception {
      UIComponent component = new UIComponent(componentId, flashSelenium);
      int retry = 0;
      // Until the value not assigned to it, It return -1 value.
      while (!component.getVisible() && retry < retryNo) {
         retry++;
         logger.info("The component is yet not visible, retry #" + retry);
         try {
            Thread.sleep(TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE);
         } catch (Exception e) {
            e.printStackTrace();
         }
         if (retry == retryNo) {
            logger.warn("The component is still empty string at the ");
         }
      }
   }

   /**
    * Method to keep the thread waiting until the Advanced datagrid loaded
    *
    * @param tableID : Specify the AdvanceDataGrid ID which we want to pass
    * @param retryNo : no. of iteration after loop exits
    */
   public static boolean waitforUIComponentVisible(String tableID, int retryNo) {
      int retry = 0;
      boolean flag = true;
      try {
         logger.info("Visibility: "
               + MethodCallUtil.getVisibleOnPath(flashSelenium, tableID));
         while (!MethodCallUtil.getVisibleOnPath(flashSelenium, tableID).equals(true)
               && retry < retryNo) {
            retry++;
            logger.info("UIComponent is still not visible, retry #" + retry);
            Thread.sleep(TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE * 1);
            if (retry == retryNo) {
               logger.warn("UIComponent is still not visible at the "
                     + "last retry, giving up...");
               flag = false;
            }
         }
      } catch (Exception e) {
         e.printStackTrace();
      }
      return flag;
   }

   public static void waitForProgressBar(String progressBarId, long timeout,
         FlashSelenium flashSelenium) {
      UIComponent progressBar = null;
      try {
         progressBar = new UIComponent(progressBarId, flashSelenium);
      } catch (AssertionError ae) {
         // do nothing looks like the progress bar disappeared fast
         logger.error("Do nothing looks like the progress bar disappeared fast. "
               + "Printing the stacktrace below:");
         ae.printStackTrace();
      }
      if (progressBar.getVisible()) {
         // wait till the time progress bar is visible
         long startTime = System.currentTimeMillis();
         while (progressBar.getVisible()) {
            if ((System.currentTimeMillis() - startTime) > timeout) {
               logger.info("Timing out after waiting for " + timeout / 1000 + " seconds.");
               return;
            }
         }
      }
   }

   /**
    * Allows the user to click Yes or No in a confirmation dialog and
    * optionally do some verification.
    *
    * @param clickYes if true it clicks Yes, otherwise it clicks No
    * @param expectedMessage the message that is expected for verification. If
    * <code>null</code>, no verification of the dialog text is made
    */
   public static void handleAndVerifyConfirmationDialog(boolean clickYes,
         String expectedMessage) {

      Alert confirmDialog = new Alert(IDConstants.ID_CONFIRMATION_DIALOG, flashSelenium);
      confirmDialog.waitForElementEnable(Timeout.TEN_SECONDS.getDuration());

      if (expectedMessage != null) {
         verifySafely(
               confirmDialog.getText(),
               expectedMessage,
               "Check confirmation message");
      }

      String buttonId =
            clickYes ? IDConstants.ID_CONFIRM_YES_LABEL
                  : IDConstants.ID_CONFIRM_NO_LABEL;

      Button button =
            new Button(IDConstants.ID_CONFIRMATION_DIALOG + "/" + buttonId,
                  flashSelenium);
      button.click();
   }

   /**
    * Allows the user to click Yes or No in a confirmation dialog
    * @param clickYes If true, Yes will be clicked, if false, No will be clicked
    */
   public static void handleConfirmationDialog(boolean clickYes) {
      handleAndVerifyConfirmationDialog(clickYes, null);
   }

   /**
    * This function allows the user to click a YesNo message box and optionally do some verification.
    *
    * @param clickYes if true it clicks yes, otherwise it clicks false
    * @param doVerification whether to verify the messagebox text
    * @param expectedMessage the message that is expected for verification
    * @throws Exception
    */
   public static void verifyAndHandleYesNoMsgBox(boolean clickYes,
         boolean doVerification, String expectedMessage) {
      Alert yesNoMsgBox = new Alert(ID_YES_NO_DIALOG, flashSelenium);
      yesNoMsgBox.waitForElementEnable(Timeout.TEN_SECONDS.getDuration());

      if (doVerification) {
         verifySafely(
               yesNoMsgBox.getText(),
               expectedMessage,
               "Confirmation message is correct.");
      }

      String buttonId = ID_NO_BUTTON;

      if (clickYes) {
         buttonId = ID_YES_BUTTON;
      }

      Button confirmationButton = new Button(buttonId, flashSelenium);
      confirmationButton.click();
   }

   /**
    * This function will refresh the UI using the refresh button and waits for
    * the node in the Invetory to be visivble
    *
    * @param nodeId
    * @return void
    */
   public static void refreshUI(String nodeId) {
      try {
         refreshUI();
         InvTree.getInstance().waitForNodeReady(nodeId, 30);
      } catch (Exception ex) {
         logger.error("Failure while waiting for node to be ready");
         ex.printStackTrace();
      }
   }

   /**
    * This function will refresh the UI using the refresh button.
    *
    * @return void
    */
   public static void refreshUI() {
      try {
         Button refreshButton = new Button(IDConstants.ID_REFRESH_BUTTON, flashSelenium);
         refreshButton.click(TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE);
         refreshButton.waitForElementEnable(DEFAULT_TIMEOUT_ONE_MINUTE);
         Thread.sleep(DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE * 3);
      } catch (Exception ex) {
         logger.error("Failure in recognizing refresh button on UI.");
         ex.printStackTrace();
      }
   }
}
