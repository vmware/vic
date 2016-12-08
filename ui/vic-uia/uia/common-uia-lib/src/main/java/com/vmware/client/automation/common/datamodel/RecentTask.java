/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.datamodel;

/**
 * Represents data for recent task retrieved from UI
 *
 * Reflects data from a row in Recent Tasks bar.
 */
public class RecentTask {

   /**
    * Name of the task.
    */
   public String name;

   /**
    * Target of the task.
    */
   public String target;

   /**
    * Status of the task.
    */
   public String status;

   /**
    * Initiator of the task.
    */
   public String initiator;

   /**
    * Milliseconds the task has stayed in the queue.
    */
   public Long queuedFor;

   /**
    * Start time of the task.
    */
   public Long startTime;

   /**
    * Completion time of the task.
    */
   public Long completionTime;

   /**
    * Server of the task.
    */
   public String server;

   /**
    * Creates a new instance of a RecentTask with identical values for member
    * variables.
    *
    * @return a copy of current object
    */
   public Object copy() {
      RecentTask o = new RecentTask();
      o.name = this.name;
      o.target = this.target;
      o.status = this.status;
      o.initiator = this.initiator;
      if (this.queuedFor != null)
         o.queuedFor = new Long(this.queuedFor.longValue());
      if (this.startTime != null)
         o.startTime = new Long(this.startTime.longValue());
      if (this.completionTime != null)
         o.completionTime = new Long(this.completionTime.longValue());
      o.server = this.server;
      return o;
   }

}
