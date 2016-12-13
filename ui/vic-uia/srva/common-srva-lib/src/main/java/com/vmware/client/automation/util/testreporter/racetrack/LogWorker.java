/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.util.testreporter.racetrack;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * The main worker thread which takes care of executing log events that are present in the queue.
 */
public class LogWorker implements Runnable {

   // Class logger.
   private static final Logger _logger = LoggerFactory.getLogger(LogWorker.class);

   // The queue that holds all of the pending log events.
   private static BlockingQueue<LogEvent> eventQueue =
         new LinkedBlockingQueue<LogEvent>();
   // Current test case id.
   private static String currentTestCaseId;

   /**
    * Adds an event to the queue.
    *
    * @param event
    */
   public static void queueEvent(final LogEvent event) {
      eventQueue.add(event);
   }

   /**
    * Sets the current test case id.
    *
    * @param newTestCaseId
    */
   public static void setCurrentTestCaseId(final String newTestCaseId) {
      LogWorker.currentTestCaseId = newTestCaseId;
   }

   /**
    * Returns the current test case id.
    *
    * @return
    */
   public static String getCurrentTestCaseId() {
      return LogWorker.currentTestCaseId;
   }

   @Override
   public void run() {
      LogEvent event = null;

      while (true) {
         try {
            // Try to get an event from the queue. Block if none is available.
            event = eventQueue.take();
         } catch (InterruptedException e) {
            // We received a request to terminate the thread, thus go ahead and terminate.
            _logger
                  .info("[Event Thread] Interrupt request received. Thread will now terminate.");
            break;
         }
         // Execute log event.
         event.execute();
      }
   }
}
