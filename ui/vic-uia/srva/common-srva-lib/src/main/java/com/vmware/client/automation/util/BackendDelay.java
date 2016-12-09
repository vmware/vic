package com.vmware.client.automation.util;

import com.vmware.client.automation.delay.Delay;

/**
 * Defines backend delay constants.
 */
public enum BackendDelay {
   SMALL(2), MEDIUM(4), LARGE(10);

   private int minutes;

   private BackendDelay(int minutes) {
      this.minutes = minutes;
   }

   /**
    * Retrieve duration for the constant.
    *
    * @return  duration of the constant
    */
   public long getDuration() {
      return  Delay.timeout.forMinutes(this.minutes).getDuration();
   }
}
