/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

/**
 * Signals that an operation orchestrated by a {@link ProviderWorkflowController}
 * has failed. The controller will attempt to restore the participating resources in
 * their original state the success of which will be indicated a property called
 * gracefullyRestored. If the participating resources are not gracefully restored,
 * they should be discarded from further use as # might be in dirty state.
 * 
 * This exception will also hold the immediate exception, if any, which has caused
 * the failure.
 */
public class ProviderControllerException extends Exception {

   /**
    * 
    */
   private static final long serialVersionUID = 1L;

   public ProviderControllerException(String message) {
      super(message);
   }

   public ProviderControllerException(String message, Throwable cause) {
      super(message, cause);
   }

}
