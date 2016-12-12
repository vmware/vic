package com.vmware.client.automation.workflow.test;


/**
 * Signals that an operation in {@link TestWorkflow} has failed.
 *
 * This exception will also hold the immediate exception, if any, which has caused
 * the failure. Only errors that occurred in provider implementations could be
 * wrapped.
 *
 * Inner exceptions should be analyzed for known information to construct
 * fine-grained report.
 */
public class TestWorkflowException extends Exception {

   private static final long serialVersionUID = 1L;

   /**
    * Default constructor. Both explanation message and cause are required parameters.
    *
    * @param message
    *    Explanation message
    *
    * @param cause
    *    Cause that came from provider's implementation
    */
   public TestWorkflowException(String message, Throwable cause) {
      super(message, cause);
   }
}
