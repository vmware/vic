/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.recenttasks;

import java.util.ArrayList;

import com.vmware.client.automation.vcuilib.commoncode.RecentTasks;

/**
 * @Deprecated by TasksUtil.java
 */
@Deprecated
public class RecentTask extends RecentTasks
{
   /**
    * @Deprecated by TasksUtil.java::getRecentTasksByFilter(RecentTaskFilter)
    *
    * Gets the time of the most recent task by task name and target name
    *
    * @param
    * @return Most recent task time in long
    */
   @Deprecated
   public static long getRecentTaskTimeByNameAndTarget (String taskName, String targetName) {
      ArrayList<RecentTasks> recentTasks = new RecentTasks()
            .getRecentTasksByTaskNameAndTarget(taskName, targetName);
      if (recentTasks == null || recentTasks.isEmpty()) {
         return -1;
      }
         return RecentTask.getMostRecentTask(recentTasks).getTaskTime();
   }

   private static RecentTasks getMostRecentTask(ArrayList<RecentTasks> tasks){
      ArrayList<RecentTasks> sortedTasks = new RecentTask().sortRecentTasks(tasks);
      return sortedTasks.get(sortedTasks.size() - 1);
   }
}
