package com.vmware.client.automation.vcuilib.commoncode;

import java.util.Comparator;

/**
 * Recent task comparator.
 *
 * NOTE: this class is a copy of the one from VCUI-QE-LIB
 */
public class RecentTasksComparator implements Comparator<RecentTasks> {

   @Override
   public int compare(RecentTasks rt1, RecentTasks rt2) {
      if (rt1.getTaskTime() instanceof Long) {
         Long rt1Object = rt1.getTaskTime();
         Long rt2Object = rt2.getTaskTime();
         return rt1Object.compareTo(rt2Object);
      }
      return -1;
   }
}
