/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.test;


/**
 * Signals that an operation orchestrated by a {@link TestWorkflowController}
 * has failed. The controller will attempt to restore the participating resources in
 * their original state the success of which will be indicated a property called
 * gracefullyRestored. If the participating resources are not gracefully restored,
 * they should be discarded from further use as # might be in dirty state.
 *
 * This exception will also hold the immediate exception, if any, which has caused
 * the failure.
 */
public class TestControllerException extends Exception {

   /**
    *
    */
   private static final long serialVersionUID = 1L;

   //  public boolean gracefullyRestored;

   public TestControllerException(String message) {
      super(message);
   }

   public TestControllerException(String message, Throwable cause) {
      super(message, cause);
   }
}
