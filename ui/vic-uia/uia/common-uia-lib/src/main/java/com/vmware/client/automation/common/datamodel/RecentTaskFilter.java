/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.datamodel;

import com.vmware.client.automation.common.spec.TaskSpec;
import com.vmware.suitaf.SUITA;

/**
 * Represents filter class containing all required properties by which a Recent
 * Task can be searched in the list retrieved from the UI
 */
public class RecentTaskFilter {

   private static final long DEFAULT_MAX_WAIT_TIME_IN_MILLISEC = SUITA.Environment
         .getBackendJobMid();

   /*******************************************************************
    * CONSTRUCTORS
    *******************************************************************/

   /**
    * Construct from scratch
    */
   public RecentTaskFilter() {
   }

   /**
    * Construct from TaskSpec
    *
    * Use this constructor when finding or waiting for tasks by given Task
    * specification.
    *
    * @param spec
    */
   public RecentTaskFilter(TaskSpec spec) {
      if (null != spec) {
         RecentTask task = new RecentTask();
         // name
         if (spec.name != null && spec.name.isAssigned()) {
            task.name = spec.name.get();
         }
         // target
         if (spec.target != null && spec.target.isAssigned()
               && spec.target.get().name != null
               && spec.target.get().name.isAssigned()) {
            task.target = spec.target.get().name.get();
         }
         // status
         if (spec.status != null && spec.status.isAssigned()) {
            task.status = spec.status.get().getName();
         }
         // initiator
         if (spec.initiator != null && spec.initiator.isAssigned()) {
            task.initiator = spec.initiator.get();
         }
         // server
         if (spec.server != null && spec.server.isAssigned()) {
            task.server = spec.server.get();
         }
         this.setTask(task);

         // max time in queue for task
         if (spec.maxTimeInQueue != null && spec.maxTimeInQueue.isAssigned()) {
            this.maxTimeInQueue = spec.maxTimeInQueue.get().longValue();
         }

         // start time from
         if (spec.startTimeFrom != null && spec.startTimeFrom.isAssigned()) {
            this.startTimeFrom = spec.startTimeFrom.get().longValue();
         }
      }
   }

   /**
    * Construct from RecentTask
    *
    * Use this constructor when you have retrieved information about a task in
    * RecentTasks and you want to continue monitoring the state of this task.
    *
    * @param task
    */
   public RecentTaskFilter(RecentTask task) {
      this.setTask(task);
   }

   /*******************************************************************
    * MEMBER VARIABLES
    *******************************************************************/

   /**
    * RecentTask object for searching by concrete values
    */
   private RecentTask _task;

   /**
    * Maximum allowed milliseconds for the task to stay in the queue.
    */
   public Long maxTimeInQueue;

   /**
    * Time after which the task was started.
    */
   public Long startTimeFrom;

   /*******************************************************************
    * GETTERS AND SETTERS
    *******************************************************************/

   /**
    * Getter method for _task
    *
    * @return _task value
    */
   public RecentTask getTask() {
      return _task;
   }

   /**
    * Setter method for _task
    *
    * @param task
    */
   public void setTask(RecentTask task) {
      if (null != task) {
         this._task = (RecentTask) task.copy();
      }
   }

   /*******************************************************************
    * ADDITIONAL METHODS
    *******************************************************************/

   /**
    * Calculates filter's maximum time for waiting for certain task in the UI to
    * reach filter criteria
    *
    * @return
    */
   public long getMaxWaitTime() {
      if (null != maxTimeInQueue)
         return maxTimeInQueue.longValue();
      return DEFAULT_MAX_WAIT_TIME_IN_MILLISEC;
   }

   /**
    * Check whether given RecentTask matches current filter criteria
    *
    * @param task
    * @return true - task matches filter criteria, false - otherwise
    */
   public boolean checkMatch(RecentTask task) {
      // empty task
      if (null == task)
         return false;

      if (null != _task) {

         // name mismatch
         if (null != _task.name && !_task.name.equals(task.name))
            return false;

         // target mismatch
         if (null != _task.target && !_task.target.equals(task.target))
            return false;

         // status mismatch
         if (null != _task.status && !_task.status.equals(task.status))
            return false;

         // initiator mismatch
         if (null != _task.initiator && !_task.initiator.equals(task.initiator))
            return false;

         // queued for mismatch
         if (null != _task.queuedFor && !_task.queuedFor.equals(task.queuedFor))
            return false;

         // start time mismatch
         if (null != _task.startTime && !_task.startTime.equals(task.startTime))
            return false;

         // completion time mismatch
         if (null != _task.completionTime
               && !_task.completionTime.equals(task.completionTime))
            return false;

         // server mismatch
         if (null != _task.server && !_task.server.equals(task.server))
            return false;
      }

      // task took too much time
      if (null != maxTimeInQueue
            && maxTimeInQueue.compareTo(task.queuedFor) < 0)
         return false;

      // task was started too early
      if (null != startTimeFrom && startTimeFrom.compareTo(task.startTime) > 0)
         return false;

      return true;
   }
}
