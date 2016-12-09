/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.util.testreporter.racetrack;

/**
 * The log event abstract class that all log events should extend.
 * Info wiki - https://wiki.eng.vmware.com/RacetrackWebServices
 * Racetrack currently could log several types of events (check the wiki above), basically they are:
 * - info messages
 * - test verifications (in case of failure it uploads screenshot)
 * - log exceptions
 */
public abstract class LogEvent {

   protected RacetrackLogger _racetrackLogger;

   /**
    * Creates basic log event.
    *
    * @param logger     racetrackLogger that will be used
    */
   public LogEvent(RacetrackLogger logger) {
      this._racetrackLogger = logger;
   }

   /**
    * Implementations of this method should execute the actual log action.
    */
   public abstract void execute();
}
