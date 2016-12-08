package com.vmware.client.automation.util;

import com.vmware.client.automation.delay.Delay;

/**
 * Defines UI delay constants.
 */
public enum UiDelay {
   PAGE_LOAD_TIMEOUT(30), UI_OPERATION_TIMEOUT(10);

   private int seconds;

   private UiDelay(int seconds) {
      this.seconds = seconds;
   }

   /**
    * Retrieve duration for the constant.
    *
    * @return  duration of the constant
    */
   public long getDuration() {
      return  Delay.timeout.forSeconds(this.seconds).getDuration();
   }
}
