/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Spec for task.
 *
 * Used to specify criteria for filtering recent tasks. For example: when
 * expecting task with certain name to be found completed in a certain time
 * duration.
 */
public class TaskSpec extends BaseSpec {

   /**
    * Name of the task.
    */
   public DataProperty<String> name;

   /**
    * Target of the task.
    */
   public DataProperty<ManagedEntitySpec> target;

   /**
    * Status of the task.
    */
   public DataProperty<TaskStatus> status;

   /**
    * Initiator of the task.
    */
   public DataProperty<String> initiator;

   /**
    * Maximum allowed milliseconds for the task to stay in the queue.
    */
   public DataProperty<Long> maxTimeInQueue;

   /**
    * Time after which the task was started.
    */
   public DataProperty<Long> startTimeFrom;

   /**
    * Server of the task.
    */
   public DataProperty<String> server;

   /**
    * Represents RecentTask status
    */
   public enum TaskStatus {
      RUNNING(CommonUtil.getLocalizedString("recenttasks.status.Running")), COMPLETED(
            CommonUtil.getLocalizedString("recenttasks.status.Completed")), FAILED(
            CommonUtil.getLocalizedString("recenttasks.status.Failed"));

      private TaskStatus(String name) {
         _name = name;
      }

      private String _name;

      public String getName() {
         return _name;
      }
   }

}
