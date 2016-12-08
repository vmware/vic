/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.command;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Interface for the workflow command operations
 */
public interface WorkflowCommand {

   /**
    * Annotation for a workflow command implementations
    */
   @Retention(RetentionPolicy.RUNTIME)
   @Target(value = ElementType.TYPE)
   public static @interface WorkflowCommandAnnotation {

      /**
       * Provide name for command implementation
       *
       * @return
       */
      public String commandName();
   }

   /**
    * Return the name of the workflow command
    *
    * @return
    */
   public String getCommandName();

   /**
    * Prepare command operation. Parse command params and etc.
    *
    * @param commandParams
    * @throws Exception
    */
   public void prepare(String[] commandParams) throws Exception;

   /**
    * Execute the command
    *
    * @throws Exception
    */
   public void execute() throws Exception;
}
