/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.util;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.datamodel.RecentTask;
import com.vmware.client.automation.common.datamodel.RecentTaskFilter;
import com.vmware.client.automation.common.datamodel.RecentTaskFilterResult;
import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.PropertiesUtil;
import com.vmware.client.automation.vcuilib.commoncode.GlobalFunction;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.TASK_STATE;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.util.CommonUtils;

/**
 * A util class for working with the recent tasks
 */
public class TasksUtil {

   private static final Logger _logger = LoggerFactory
         .getLogger(BaseView.class);

   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private static final String ID_RECENT_TASKS_GRID = "recentTasksGridView";
   private static final String ID_RECENT_TASKS_ALL_PROPERTIES = "dataProvider.source.[].{all}";

   private static final String RECENT_TASKS_DATE_FORMAT = "EEE MMM dd HH:mm:ss 'GMT'Z yyyy";
   private static final String RECENT_TASKS_PROPERTY_NAME = "name";
   private static final String RECENT_TASKS_PROPERTY_TARGET = "targetName";
   private static final String RECENT_TASKS_PROPERTY_STATUS = "status";
   private static final String RECENT_TASKS_PROPERTY_INITIATOR = "initiator";
   private static final String RECENT_TASKS_PROPERTY_QUEUED_FOR = "timeInQueue";
   private static final String RECENT_TASKS_PROPERTY_START_TIME = "actualStartTime";
   private static final String RECENT_TASKS_PROPERTY_COMPLETION_TIME = "completedStartTime";
   private static final String RECENT_TASKS_PROPERTY_SERVER = "vcServer";

   /**
    * RegEx for finding properties of a recent task.
    *
    * First string is replaced with property name. Second string is either
    * regular expression for retrieval of property value or a specific property
    * value.
    *
    * Matches: "\nstatusSummary:Completed\n"
    */
   private static final String RECENT_TASKS_PROPERTY_REGEX = "\\n%s:%s\\n";
   /**
    * RegEx for retrieval of property value.
    *
    * Matches: any symbol
    */
   private static final String RECENT_TASKS_PROPERTY_RETRIEVE_REGEX = "(.*?)";
   /**
    * RegEx for finding a combination of properties and their values in a recent
    * task string.
    *
    * Matches both:
    * "\nstatusSummary:Completed\n...\nname:Power Off virtual machine\n"
    * "\nname:Power Off virtual machine\n...\nstatusSummary:Completed\n"
    */
   private static final String RECENT_TASKS_PROPERTY_PREFILTER_REGEX = "(?=.*?(%s))";

   public static long TASK_SEARCH_WAIT_PERIOD_STEP_IN_MILLISEC = 4000;

   /**
    * @Deprecated by TaskSpec.java::TaskStatus
    *
    *             Enum for the recent task states
    */
   @Deprecated
   public enum TaskStatus {
      RUNNING, SUCCESS, FAILURE
   }

   @Deprecated
   private static final String DIVIDER = " /// ";

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    *             Gets the status for a recent task
    *
    * @param taskName
    *           - the name of the task
    * @param targetName
    *           - the task target
    * @return the status for a recent task
    */
   @Deprecated
   public static TaskStatus getTaskStatus(String taskName, String targetName) {
      String[] allTasks = GlobalFunction.getAllTasks(BrowserUtil.flashSelenium);

      for (String task : allTasks) {
         String taskParts[] = task.split(DIVIDER);
         String name = taskParts[0];
         String target = taskParts[1];
         String status = taskParts[2];

         if (name.equals(taskName) && target.equals(targetName)) {
            return toTaskStatus(status);
         }
      }

      SUITA.Factory.UI_AUTOMATION_TOOL.audit.snapshotAppScreen(
            "NoTasksWereFoundForTaskName", "NoTasksWereFoundForTaskName'"
                  + taskName + "' and target: " + targetName);
      throw new RuntimeException("No tasks were found for task name '"
            + taskName + "' and target: " + targetName);
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    *             Gets the info message of a recent task
    *
    * @param taskName
    *           - the name of the task
    * @param targetName
    *           - the task target
    * @return the info message of a recent task
    */
   @Deprecated
   public static String getTaskInfo(String taskName, String targetName) {
      String[] allTasks = GlobalFunction.getAllTasks(BrowserUtil.flashSelenium);

      for (String task : allTasks) {
         String taskParts[] = task.split(DIVIDER);
         String name = taskParts[0];
         String target = taskParts[1];
         String info = taskParts[6];

         if (name.equals(taskName) && target.equals(targetName)) {
            return info;
         }
      }

      throw new RuntimeException("No tasks were found for task name '"
            + taskName + "' and target: " + targetName);
   }

   @Deprecated
   private static TaskStatus toTaskStatus(String taskState) {
      if (taskState.equals(TestConstants.TASK_STATUS_IN_PROGRESS)
            || taskState.equals(TestConstants.TASK_STATUS_RUNNING)) {
         return TaskStatus.RUNNING;
      } else if (taskState.equals(TestConstants.TASK_STATUS_SUCCESS)) {
         return TaskStatus.SUCCESS;
      } else if (taskState.equals(TestConstants.TASK_STATUS_FAILED)) {
         return TaskStatus.FAILURE;
      } else {
         throw new IllegalArgumentException("Can not convert TASK_STATE: "
               + taskState);
      }
   }

   @Deprecated
   private static TASK_STATE toTaskState(TaskStatus taskStatus) {
      switch (taskStatus) {
      case RUNNING:
         return TASK_STATE.TASK_RUNNING;
      case SUCCESS:
         return TASK_STATE.TASK_COMPLETED;
      case FAILURE:
         return TASK_STATE.TASK_FAILED;

      default:
         throw new IllegalArgumentException("Can not convert TaskStatus: "
               + taskStatus);
      }
   }

   /**
    * @Deprecated by
    *             TasksUtil.java::waitForRecentTaskToMatchFilter(RecentTaskFilter
    *             )
    *
    *             Searches the "Recent tasks" portlet for a running task with
    *             particular name and waits for it to complete.
    *
    * @param taskName
    *           The name of the task
    * @param targetName
    *           The target object of the task
    * @param timeout
    *           The timeout of the entire wait process if it is not reached the
    *           expected state or error
    * @param taskStartedAfter
    *           Look only for tasks started after this time
    */
   @Deprecated
   public static void waitForTaskToComplete(String taskName, String targetName,
         long timeout, long... taskStartedAfter) {
      GlobalFunction.waitUntilTasksState(taskName, targetName,
            TestConstants.TASK_STATE.TASK_COMPLETED, timeout,
            TimeUnit.SECONDS.toMillis(3), taskStartedAfter);
   }

   /**
    * Searches the "Recent tasks" portlet for tasks matching single filter
    * criteria. Waits until any task matches search criteria or the timeout
    * specified by maxTimeInQueue exceeds. If no maxTimeInQueue value is set for
    * filter, value of SUITA.Environment.getBackendJobMid() is assumed for
    * default.
    *
    * @param filter
    * @return true if at least one task matching filter criteria was found in
    *         given time period, false - otherwise
    */
   public boolean waitForRecentTaskToMatchFilter(final RecentTaskFilter filter) {
      return waitForRecentTaskToMatchFilter(Collections.singletonList(filter));
   }

   /**
    * Searches the "Recent tasks" portlet for tasks matching specified filter
    * criteria. For every filter is waited until matching task is found or time
    * expires. When no maxTimeInQueue value is specified for the filter, the
    * value of SUITA.Environment.getBackendJobMid() is assumed for default wait
    * period.
    *
    * @param filters
    *           - list of filters
    * @return true - if all filters have at least one matching task, false - if
    *         some or all of the filters have no matching tasks, or unexpected
    *         behavior has occurred
    */
   public boolean waitForRecentTaskToMatchFilter(
         final List<RecentTaskFilter> filters) {

      if (null == filters) {
         _logger.error("Invalid arguments: null value for filters");
         return false;
      }

      // get list of all not fulfilled yet filters
      List<RecentTaskFilter> currentFilters = new ArrayList<RecentTaskFilter>();
      for (RecentTaskFilter filter : filters) {
         if (null == filter) {
            _logger.error("Invalid arguments: null value for filter");
            return false;
         }
         currentFilters.add(filter);
      }

      // wait for matching tasks
      long waitStartTime = System.currentTimeMillis();
      while (currentFilters.size() > 0) {

         // wait for tasks to match filters
         CommonUtils.sleep(TASK_SEARCH_WAIT_PERIOD_STEP_IN_MILLISEC);
         new BaseView().refreshPage();

         // read recent tasks by filters
         List<RecentTaskFilterResult> results = getRecentTasksByFilter(currentFilters);
         if (null == results || results.size() != currentFilters.size()) {
            _logger.error("Unexpected result from getRecentTasksByFilter");
            return false;
         }

         // check for tasks found by filters
         for (RecentTaskFilterResult result : results) {
            if (null == result) {
               _logger
                     .error("Unexpected result values from getRecentTasksByFilter");
               return false;
            }
            RecentTaskFilter filter = result.getFilter();

            // remove filters that have matching tasks
            if (result.getMatchingTasks().size() > 0) {
               currentFilters.remove(currentFilters.indexOf(filter));
               continue;
            }

            // check whether max wait time for filter has expired
            if (System.currentTimeMillis() - waitStartTime > filter
                  .getMaxWaitTime()) {
               _logger.info("Maximum wait time for filter has expired");
               UI.audit.snapshotAppScreen(SubToolAudit.getFPID(),
                     "RECENT_TASK_MATCH");
               return false;
            }
         }
      }

      _logger.info("All filters have been satisfied");
      return true;
   }

   /**
    * Retrieves all tasks from Recent Tasks bar in UI and filters them by given
    * criteria. If null value is passed for filter, all visible in the UI tasks
    * are returned.
    *
    * @param filter
    * @return list of RecentTask objects matching search criteria
    */
   public RecentTaskFilterResult getRecentTasksByFilter(
         final RecentTaskFilter filter) {
      List<RecentTaskFilterResult> results = getRecentTasksByFilter(Collections
            .singletonList(filter));
      if (null == results || results.size() != 1) {
         _logger.error("Unexpected result from getRecentTasksByFilter");
         return null;
      }
      return results.get(0);
   }

   /**
    * Retrieves all tasks from Recent Tasks bar in UI and filters them by given
    * criteria. Groups the results by filter. If null value is passed for
    * filters, all visible in the UI tasks are returned.
    *
    * @param filters
    *           - search criteria
    * @return list of search results by filters, one result object for every
    *         filter; if null value is passed for filters, the list contains one
    *         element
    */
   public List<RecentTaskFilterResult> getRecentTasksByFilter(
         List<RecentTaskFilter> filters) {

      if (null == filters) {
         filters = Collections.singletonList(null);
      }

      // prepare output slots for every filter
      List<RecentTaskFilterResult> results = new ArrayList<RecentTaskFilterResult>();
      for (int i = 0; i < filters.size(); i++) {
         results.add(new RecentTaskFilterResult(filters.get(i)));
      }

      // get all tasks from UI and pre-filter
      List<Pattern> patterns = new ArrayList<Pattern>();
      for (RecentTaskFilter filter : filters) {
         Pattern pattern = formRecentTaskFilterRexEx(filter);
         patterns.add(pattern);
      }
      List<String> uiTasks = getRecentTasksFromUI(patterns);
      if (null == uiTasks)
         return results;

      // form recent tasks list
      for (String uiTask : uiTasks) {
         _logger
               .debug(String.format("Filtering recent task string: %s", uiTask));
         RecentTask task = formRecentTaskFromString(uiTask);
         for (RecentTaskFilterResult result : results) {
            RecentTaskFilter filter = result.getFilter();
            if (null == filter || filter.checkMatch(task)) {
               result.getMatchingTasks().add(task);
            }
         }
      }

      return results;
   }

   /**
    * Forms regular expression pattern for RecentTaskFilter
    *
    * @param filter
    * @return regular expression pattern matching RecentTaskFilter
    */
   private Pattern formRecentTaskFilterRexEx(RecentTaskFilter filter) {
      if (null == filter || null == filter.getTask())
         return null;

      String filterRegEx = "";

      // form filter regular expression
      RecentTask taskFilter = filter.getTask();

      // with task name
      filterRegEx += formRecentTaskPropertyFilterRegEx(taskFilter.name,
            RECENT_TASKS_PROPERTY_NAME);

      // with task target
      filterRegEx += formRecentTaskPropertyFilterRegEx(taskFilter.target,
            RECENT_TASKS_PROPERTY_TARGET);

      // with task status
      filterRegEx += formRecentTaskPropertyFilterRegEx(taskFilter.status,
            RECENT_TASKS_PROPERTY_STATUS);

      // with task initiator
      filterRegEx += formRecentTaskPropertyFilterRegEx(taskFilter.initiator,
            RECENT_TASKS_PROPERTY_INITIATOR);

      // with task server
      filterRegEx += formRecentTaskPropertyFilterRegEx(taskFilter.server,
            RECENT_TASKS_PROPERTY_SERVER);

      if (filterRegEx.isEmpty())
         return null;

      return Pattern.compile(filterRegEx, Pattern.DOTALL);
   }

   /**
    * Forms regular expression string that can be used to find task string
    * representation with given value for property.
    *
    * @param filterValue
    *           - property value to filter by it
    * @param propertyName
    *           - property which task is being filtered by
    * @return regular expression string
    */
   private String formRecentTaskPropertyFilterRegEx(String filterValue,
         String propertyName) {

      String propertyRegEx = "";

      if (null != filterValue && !filterValue.isEmpty()) {
         // form property filter pattern string
         propertyRegEx = String.format(RECENT_TASKS_PROPERTY_PREFILTER_REGEX,
               String.format(RECENT_TASKS_PROPERTY_REGEX, propertyName,
                     filterValue));
      }

      return propertyRegEx;
   }

   /**
    * Retrieves all properties of all recent tasks records present in the UI
    *
    * @param patterns
    *           - list of pre-filter patterns
    * @return list of strings, each representing one record from Recent Tasks
    *         bar
    */
   private List<String> getRecentTasksFromUI(final List<Pattern> patterns) {
      List<String> uiTasks = new ArrayList<String>();

      // retrieve all properties values
      List<String> propList = new ArrayList<String>();
      propList.add(ID_RECENT_TASKS_ALL_PROPERTIES);

      Map<String, String[]> requestedProps = PropertiesUtil.getProperties(
            ID_RECENT_TASKS_GRID, propList.toArray(new String[0]));

      // get recent tasks properties as strings
      String[] propertyStrings = (requestedProps != null) ? requestedProps
            .get(ID_RECENT_TASKS_ALL_PROPERTIES) : null;

      // no tasks found
      if (null == propertyStrings)
         return uiTasks;

      // form list with pre-filter
      if (null != patterns && patterns.size() > 0) {
         for (String taskString : propertyStrings) {
            taskString = "\n" + taskString + "\n";
            _logger.debug(String.format("Pre-filtering recent task string: %s",
                  taskString));
            boolean matchesAnyPattern = false;
            for (Pattern pattern : patterns) {
               _logger.debug(String.format("Matching pattern: %s", pattern));
               if (null == pattern || pattern.matcher(taskString).find()) {
                  matchesAnyPattern = true;
                  break;
               }
            }

            if (matchesAnyPattern) {
               uiTasks.add(taskString);
            }
         }
      }
      // form list with no pre-filter
      else {
         for (String taskString : propertyStrings) {
            taskString = "\n" + taskString + "\n";
            uiTasks.add(taskString);
         }
      }

      return uiTasks;
   }

   /**
    * Constructs RecentTask by given string containing all its properties,
    * aligned on separate lines and separated from value by symbol ":".
    *
    * @param propertyString
    * @return RecentTask with corresponding properties
    */
   private RecentTask formRecentTaskFromString(String propertyString) {

      RecentTask task = new RecentTask();

      SimpleDateFormat sdf = new SimpleDateFormat(RECENT_TASKS_DATE_FORMAT);
      String returnedValue = null;

      // task name
      task.name = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_NAME);

      // task target
      task.target = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_TARGET);

      // task status
      task.status = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_STATUS);

      // task initiator
      task.initiator = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_INITIATOR);

      // task queued for
      returnedValue = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_QUEUED_FOR);
      if (returnedValue != null) {
         task.queuedFor = new Long(returnedValue);
      }

      // task start time
      returnedValue = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_START_TIME);
      if (returnedValue != null) {
         try {
            task.startTime = sdf.parse(returnedValue).getTime();
         } catch (ParseException pe) {
            _logger.warn("Unable to parse task start time", pe.getMessage());
         }
      }

      // task completion time
      returnedValue = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_COMPLETION_TIME);
      if (returnedValue != null) {
         try {
            task.completionTime = sdf.parse(returnedValue).getTime();
         } catch (ParseException pe) {
            _logger.warn("Unable to parse task completion time",
                  pe.getMessage());
         }
      }

      // task server
      task.server = getRecentTaskPropertyValue(propertyString,
            RECENT_TASKS_PROPERTY_SERVER);

      return task;
   }

   /**
    * Retrieves property value by given recent task string and property name
    *
    * @param propertyString
    *           - string containing property name and property value
    * @param propertyName
    *           - property name, which value has to be extracted
    * @return
    */
   private String getRecentTaskPropertyValue(String propertyString,
         String propertyName) {

      // form property value pattern
      Pattern pattern = Pattern.compile(String.format(
            RECENT_TASKS_PROPERTY_REGEX, propertyName,
            RECENT_TASKS_PROPERTY_RETRIEVE_REGEX));

      // find pattern in propertyString
      Matcher matcher = pattern.matcher(propertyString);

      // return property value
      return ((matcher.find()) ? matcher.group(1) : null);
   }
}