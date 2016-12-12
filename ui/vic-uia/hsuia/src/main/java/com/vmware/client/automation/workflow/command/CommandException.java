/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.command;

/**
 * The class defines command workflow exception.
 */
public class CommandException extends Exception {

   private static final long serialVersionUID = 193232930715504093L;

   public CommandException(String message) {
      super(message);
   }

   public CommandException(String message, Throwable cause) {
      super(message, cause);
   }
}
