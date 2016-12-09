/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.view;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.datamodel.RecentTask;
import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.datamodel.RecentTaskFilterResult;
import com.vmware.client.automation.common.spec.TaskSpec.TaskStatus;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.util.TasksUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.CommonUtils;

/** Class which represents a single view in the client. */
public class BaseView {

   protected static final Logger _logger = LoggerFactory
         .getLogger(BaseView.class);

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private static final IDGroup REFRESH_BTN = IDGroup
         .toIDGroup("refreshButton");

   private static IDGroup ID_LOGGED_IN_BTN = IDGroup
         .toIDGroup("mainControlBar/loggedInUser");

   private static final LoginView _loginView = new LoginView();

   /**
    * Refreshes the main page of the NGC. If refresh is already running, the
    * method just waits for it to complete.
    */
   public void refreshPage() {
      _logger.info("Refresh main page");

      if (UI.component.exists(REFRESH_BTN)) {
         _logger.info("Click on the Refresh button");
         try {
            UI.component.click(REFRESH_BTN);
         } catch (RuntimeException e) {
            if (e.getMessage().contains(REFRESH_BTN.toString())) {
               _logger.warn("Auto refresh is invoked and the refresh button "
                     + "can not be found to click on it!");
            } else {
               throw e;
            }
         }
      } else {
         _logger.info("Global refresh already invoked and running");
      }

      waitForPageToRefresh();
   }

   /**
    * Waits for the main page of the NGC to refresh. The method assumes that
    * refreshing is complete once the "Loading" spinner is replaced with the
    * "Refresh" button.
    */
   public void waitForPageToRefresh() {
      _logger.debug(">>> Wait main page to refresh:START <<<");

      // set timeout to larger value - backendJobMid - as some pages load longer
      UI.condition.isFound(REFRESH_BTN).await(
            SUITA.Environment.getBackendJobMid());

      _logger.debug(">>> Wait main page to refresh:END <<<");
   }

   /**
    * Searches the "Recent tasks" portlet for running tasks. Filters UI view for
    * current user tasks only and waits until all tasks complete or the timeout
    * specified by SUITA.Environment.getBackendJobMid() exceeds.
    *
    * @return true - if no running tasks are found, that is all tasks are either
    *         completed or have failed, false - if time has expired and there
    *         are still running tasks
    *
    * @Deprecated by
    *             TasksUtil.java::waitForRecentTaskToMatchFilter(RecentTaskFilter
    *             )
    */
   @Deprecated
   public boolean waitForRecentTaskCompletion() {
      _logger.info("Wait for the task running in the recent tasks portlet");

      // define filter
      RecentTask filterTask = new RecentTask();
      filterTask.status = TaskStatus.RUNNING.getName();
      RecentTaskFilter filter = new RecentTaskFilter(filterTask);

      // wait for running tasks
      long waitStartTime = System.currentTimeMillis();
      do {
         // read recent tasks by filter
         RecentTaskFilterResult result = new TasksUtil()
               .getRecentTasksByFilter(filter);

         // check no tasks in running state
         if (null != result && result.getMatchingTasks().size() == 0) {
            return true;
         }

         // wait for tasks to complete
         CommonUtils.sleep(TasksUtil.TASK_SEARCH_WAIT_PERIOD_STEP_IN_MILLISEC);
      } while (System.currentTimeMillis() - waitStartTime <= filter
            .getMaxWaitTime());

      UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "RECENT_TASK_RUNNING");
      return false;
   }

   /**
    * Checks whether the main NGC page is open. The result is based on the
    * visibility of the logged-in button.
    *
    * @return True if the main page is open, false otherwise
    */
   protected boolean isMainPageOpen() {
      // Loading the main page takes more time => longer timeout
      // TODO: rkovachev - decrease the timeout back to 1 minutes after PR
      // 1491962 is resolved
      return isMainPageOpen(SUITA.Environment.getPageLoadTimeout() * 6);
   }

   /**
    * Checks whether the main NGC page is open. The result is based on the
    * visibility of the logged-in button.
    *
    * @param timeout
    *           The time to wait for the main page to open in mili sec.
    *
    * @return True if the main page is open, false otherwise
    */
   protected boolean isMainPageOpen(long timeout) {
      // Loading the main page takes more time => longer timeout
      return UI.condition.isFound(ID_LOGGED_IN_BTN).await(timeout);
   }

   /**
    * Returns the currently logged user.
    *
    */
   protected String getLoggedInUsername() {
      String loggedInUser = UI.component.property.get(Property.TEXT,
            ID_LOGGED_IN_BTN);

      _logger.info(String.format("The logged in user is %s", loggedInUser));

      return loggedInUser.split("@")[0];
   }

   /**
    * Logs out the current user of the NGC client.
    */
   public void logout() throws Exception {
      _logger.info("Logging out of the application.");
      if (isMainPageOpen(0)) {
         UI.component.click(ID_LOGGED_IN_BTN);
         ActionNavigator.invokeLogoutMenuItem();
      }
      if (!_loginView.isLoginPageOpen()) {
         throw new Exception(
               "Log in page is not loaded after clicking logout button");
      }
   }

   /**
    * Logs out the current user of the NGC client.
    */
   public void logoutWithWait() throws Exception {
      _logger.info("Logging out of the application.");
      if (isMainPageOpen(0)) {
         UI.component.click(ID_LOGGED_IN_BTN);
         ActionNavigator.invokeLogoutMenuItem();
         Thread.sleep(5000);

         int counter = 10;
         while (counter > 0) {
            if (!_loginView.isLoginPageOpen()) {
               Thread.sleep(5000);
               counter--;
            } else {
               break;
            }
         }

         if (counter == 0) {
            throw new Exception(
                  "Log in page is not loaded after clicking logout button");
         }
      }
      else if (!_loginView.isLoginPageOpen()) {
         throw new Exception(
               "Log in page is not loaded after clicking logout button");
      }
   }

   /**
    * Filter the tasks portlet to show tasks filtered by the current user.
    *
    * @return true if the filter is selected
    */
   public boolean showMyTasks() {
      IDGroup userTasksFilterMenu = IDGroup.toIDGroup("userFilterButton");
      IDGroup myTasksMenuItem = IDGroup.toIDGroup("automationName=My Tasks");
      String validationString = "My Tasks";

      try {
         _logger.info("Click on filter by taks owner menu");
         UI.component.click(userTasksFilterMenu);
         _logger.info("Find and select the My tasks menu item");
         UI.component.click(myTasksMenuItem);
         _logger.info("Check my tasks");
      } catch (RuntimeException e) {
         _logger.warn(e.getMessage());
         UI.audit
               .snapshotAppScreen(SubToolAudit.getFPID(), "FailToFilterTasks");
         return false;
      }
      return UI.condition.isTrue(
            UI.component.property.get(Property.TEXT, userTasksFilterMenu)
                  .equalsIgnoreCase(validationString)).await(3000);
   }
}
