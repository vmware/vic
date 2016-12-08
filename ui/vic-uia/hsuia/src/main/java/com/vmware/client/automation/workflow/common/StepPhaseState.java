/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.common;

/**
 * Define the phases of a step.
 */
public enum StepPhaseState {
   BLOCKED,
   READY_TO_START,
   IN_PROGRESS,
   DONE,
   FAILED
}
