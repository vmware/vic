/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

/**
 * Exception thrown when the same spec is registered more than once.
 *
 */
public class DuplicateSpecFoundException extends Exception {

   private static final long serialVersionUID = 0;

   public DuplicateSpecFoundException(String message) {
      super(message);
   }
}
