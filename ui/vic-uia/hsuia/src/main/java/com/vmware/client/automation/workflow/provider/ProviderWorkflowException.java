/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;


/**
 * Signals that an operation in {@link ProviderWorkflow} has failed. 
 * 
 * This exception will also hold the immediate exception, if any, which has caused
 * the failure. Only errors that occurred in provider implementations could be 
 * wrapped.
 * 
 * Inner exceptions should be analyzed for known information to construct
 * fine-grained report.
 */
public class ProviderWorkflowException extends Exception {

   private static final long serialVersionUID = 6204538698431763608L;

   /**
    * Default constructor. Both explanation message and cause are required parameters. 
    * 
    * @param message
    *    Explanation message
    *    
    * @param cause
    *    Cause that came from provider's implementation
    */
   public ProviderWorkflowException(String message, Throwable cause) {
      super(message, cause);
   }

}
