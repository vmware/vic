package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Date;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @Deprecated by TasksUtil.java
 *
 * Recent task implementation.
 *
 * NOTE: this class is a copy of the one from VCUI-QE-LIB
 */
@Deprecated
public class RecentTasks {

   private static final Logger logger = LoggerFactory.getLogger(RecentTasks.class);

   /**
    * The name of the task
    */
   private String taskName;

   /**
    * The target of the task
    */
   private String taskTarget;

   /**
    * The status of the task
    */
   private String status;

   /**
    * The progress of the task
    */
   private String progress;

   /**
    * The start Time of the task
    */
   private Long startTime;

   /**
    * The end time of the task
    */
   private Long endTime;


   private static final String DIVIDER = " /// ";
   private static final String DATE_PATTERN = "EEE MMM dd HH:mm:ss 'GMT'Z yyyy";

   // TODO Need to find a way to parse the Time zone and year after it
   private static final String DATE_FORMAT = "EEE MMM dd HH:mm:ss";

   public RecentTasks() {
   }

   public RecentTasks(String recentTaskName, String recentTaskTarget, Long recentTaskTime) {
      this.setTaskName(recentTaskName);
      this.setTaskTime(recentTaskTime);
      this.setTaskTarget(recentTaskTarget);
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    * Returns a task based on the string representation of a task. The client
    * gets the task as:
    * taskName + divider + targetEntity + divider + status + divider + progress
    *  + divider + startTime + divider + completeTime; where divider is a '/'
    * @param task
    * @return
    */
   public static RecentTasks fromString(String task) {
      if (task == null) {
         throw new RuntimeException("The String representation of the task is null");
      }
      String taskParts[] = task.split(DIVIDER);

      RecentTasks recentTask = new RecentTasks();
      recentTask.taskName = taskParts[0];
      recentTask.taskTarget = taskParts[1];
      recentTask.status = taskParts[2];
      recentTask.progress = taskParts[3];
      recentTask.startTime = parseDate(taskParts[4]).getTime();

      //the end time might not be present if the task is currently running
      //so check for it before setting it
      if (taskParts.length > 5 && taskParts[5] != null && !taskParts[5].equals("null")) {
         recentTask.endTime = parseDate(taskParts[5]).getTime();
      }

      return recentTask;
   }

   /**
    * Parses the date as a String and returns its value as a date.
    * @param date
    * @return
    */
   private static Date parseDate(String date) {
      DateFormat dateFormat = new SimpleDateFormat(DATE_PATTERN);
      try {
         return dateFormat.parse(date);
      } catch (ParseException e) {
         logger.error("There was an error parsing the recent task date:" + date);
         return null;
      }
   }

   /**
    * Gets recent task's name
    *
    * @return String with task name
    */
   public String getTaskName() {
      return taskName;
   }

   /**
    * Gets recent task time
    *
    * @return Long object with the datetime
    */
   public Long getTaskTime() {
      return startTime;
   }

   /**
    * Gets recent task's target
    *
    * @return Task target
    */
   public String getTaskTarget() {
      return taskTarget;
   }

   public String getStatus() {
      return status;
   }

   public String getProgress() {
      return progress;
   }

   /**
    * Sets task time of recent task
    *
    * @param taskTime
    */
   private void setTaskTime(Long taskTime) {
      this.startTime = taskTime;
   }

   /**
    * Sets task target of recent task
    *
    * @param taskTarget
    */
   private void setTaskTarget(String taskTarget) {
      this.taskTarget = taskTarget;
   }

   /**
    * Sets recent task name
    *
    * @param taskName
    */
   private void setTaskName(String taskName) {
      this.taskName = taskName;
   }

   private interface IRecentTaskFilter {
      public boolean isOk(String task);
   }

   private class TaskByAttribute implements IRecentTaskFilter {
      private String recentTaskAttribute = null;

      public TaskByAttribute(String recentTaskAttribute) {
         this.recentTaskAttribute = recentTaskAttribute;
      }

      @Override
      public boolean isOk(String task) {
         return task.contains(recentTaskAttribute);
      }
   }

   private class TasksByTaskNameAndTarget implements IRecentTaskFilter {
      private String recentTaskName = null;
      private String recentTaskTarget = null;

      public TasksByTaskNameAndTarget(String recentTaskName, String recentTaskTarget) {
         this.recentTaskName = recentTaskName;
         this.recentTaskTarget = recentTaskTarget;
      }

      @Override
      public boolean isOk(String task) {
         return task.contains(recentTaskName) && task.contains(recentTaskTarget);
      }
   }

   private class TasksByTaskNameAndTargets implements IRecentTaskFilter {
      private String recentTaskName = null;
      private ArrayList<String> recentTaskTargets = null;

      public TasksByTaskNameAndTargets(String recentTaskName,
            ArrayList<String> recentTaskTargets) {
         this.recentTaskName = recentTaskName;
         this.recentTaskTargets = recentTaskTargets;
      }

      @Override
      public boolean isOk(String task) {
         boolean found = false;
         for (int i = 0; i < recentTaskTargets.size(); i++) {
            if (task.contains(recentTaskName) && task.contains(recentTaskTargets.get(i))) {
               found = true;
            }
         }
         return found;
      }
   }

   private static ArrayList<RecentTasks> getRecentTasks(IRecentTaskFilter filter) {
      ArrayList<RecentTasks> relevantTasks = new ArrayList<RecentTasks>();
      String[] allTasks = GlobalFunction.getAllTasks(flashSelenium);
      RecentTasks rt = null;
      for (String s : allTasks) {
         if (filter.isOk(s)) {
            SimpleDateFormat sdf = new SimpleDateFormat(DATE_PATTERN);
            try {
               rt =
                     new RecentTasks(s.split(" /// ")[0], s.split(" /// ")[1], sdf
                           .parse(s.split(" /// ")[4]).getTime());
               relevantTasks.add(rt);
            } catch (ParseException e) {
               e.printStackTrace();
            }
         }
      }
      return relevantTasks;
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    * Returns all recent tasks matching the attribute passed, i.e. task name
    *
    * @param recentTaskAttribute
    * @return List with all recent tasks found by the specified criteria
    */
   @Deprecated
   public ArrayList<RecentTasks> getRecentTasksByAttribute(String recentTaskAttribute) {
      return getRecentTasks(new TaskByAttribute(recentTaskAttribute));
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    * Gets recent tasks by specified target name and task name
    *
    * @param recentTaskName
    * @param recentTaskTarget
    * @return List with all recent tasks found by the specified criteria
    */
   @Deprecated
   public ArrayList<RecentTasks> getRecentTasksByTaskNameAndTarget(
         String recentTaskName, String recentTaskTarget) {
      return getRecentTasks(new TasksByTaskNameAndTarget(recentTaskName,
            recentTaskTarget));
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    * Gets recent tasks by specified task name and multiple target names
    *
    * @param recentTaskName
    * @param recentTaskTargets - ArrayList with target names
    * @return List with all recent tasks found by the specified criteria
    */
   @Deprecated
   public ArrayList<RecentTasks> getRecentTasksByTaskNameAndTargets(
         String recentTaskName, ArrayList<String> recentTaskTargets) {
      return getRecentTasks(new TasksByTaskNameAndTargets(recentTaskName,
            recentTaskTargets));
   }

   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    * Sorts recent tasks in ascending order, latest at the end
    *
    * @param recentTasksToSort
    * @return Sorted by date ArrayList with recent tasks
    */
   @Deprecated
   public ArrayList<RecentTasks> sortRecentTasks(ArrayList<RecentTasks> recentTasksToSort) {
      RecentTasksComparator rtComparator = new RecentTasksComparator();
      Collections.sort(recentTasksToSort, rtComparator);
      return recentTasksToSort;
   }
}
