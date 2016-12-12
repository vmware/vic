/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.common;

/**
 * Enum defining the states of a workflow.
 */
public enum WorkflowPhaseState {
   BLOCKED, READY_TO_START, IN_PROGRESS, PASSED, FAILED, SKIPPED
}
