/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.lang.reflect.Modifier;
import java.util.HashSet;
import java.util.Set;

import org.reflections.Reflections;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.workflow.command.WorkflowCommand;
import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.common.Workflow;

/**
 * The class provides methods to browse the class path and search for commands
 * and workflows.
 * TODO rkovachev :make the browsed packages configurable
 */
public class WorkflowBrowser {

   private final static String WORKFLOW_PACKAGES = "com.vmware";
   private final static String COMMAND_PACKAGES = "com.vmware.client.automation";

   private static final Logger _logger = LoggerFactory.getLogger(WorkflowBrowser.class);

   /**
    * Discover all workflows in the loaded class paths that represent or extend
    * the specified workflow class inside the package scope.
    *
    * @return
    */
   public static Set<Class<? extends Workflow>> findExecutableWorkflows(/*Class<? extends Workflow> workflowBaseClass*/) {
      Reflections reflections = new Reflections(WORKFLOW_PACKAGES);

      Set<Class<? extends Workflow>> allWorkflowClasses =
            reflections.getSubTypesOf(Workflow.class);
      Set<Class<? extends Workflow>> instWorkflowClasses =
            new HashSet<Class<? extends Workflow>>();

      for (Class<? extends Workflow> workflowClass : allWorkflowClasses) {
         try {
            _logger.debug(workflowClass.getCanonicalName());
            if (Modifier.isAbstract(workflowClass.getModifiers())) {
               _logger.debug("ABSTRACT");
               continue;
            }
            try {
               workflowClass.newInstance();
            } catch (NoClassDefFoundError error) {
               _logger.error("NoClassDefFoundError:" + error.getMessage());
            } catch (UnsatisfiedLinkError uerror) {
               _logger.error("UnsatisfiedLinkError:" +uerror.getMessage());
            }
            instWorkflowClasses.add(workflowClass);
         } catch (InstantiationException e) {
            // Cannot instantiate. Skip this workflow class.
            continue;
         } catch (IllegalAccessException e) {
            // Cannot instantiate. Skip this worklfow class
            continue;
         }
      }

      return instWorkflowClasses;

      //		String packageScope = WORKFLOW_PACKAGES;
      //
      //      Reflections reflections =
      //            new Reflections(
      //                  packageScope, new SubTypesScanner(false), ClasspathHelper.forClassLoader());
      //
      //      Set<Class<? extends Workflow>> allSubTypes = reflections.getSubTypesOf(Workflow.class);
      //
      //      Set<Class<? extends Workflow>> result = new HashSet<Class<? extends Workflow>>();
      //
      //      boolean includeClass = false;
      //      boolean traverseSubPackages = true;
      //      for (Class<? extends Workflow> testClass : allSubTypes) {
      //         if (traverseSubPackages) {
      //            includeClass = testClass.getPackage().getName().startsWith(packageScope);
      //         } else {
      //            includeClass = testClass.getPackage().getName().equals(packageScope);
      //         }
      //         if (includeClass) {
      //            result.add(testClass);
      //         }
      //         includeClass = false;
      //      }
      //
      //		return result;
   }

   /**
    * Browse the classpath and reutrn set of commands. NOTE: It browse the
    * package specified by COMMAND_PACKAGES
    *
    * @return
    */
   public static Set<Class<? extends WorkflowCommand>> findExecutableCommands() {
      Reflections reflections =  new Reflections(COMMAND_PACKAGES);
      Set<Class<? extends WorkflowCommand>> allCommandClasses = reflections
            .getSubTypesOf(WorkflowCommand.class);

      Set<Class<? extends WorkflowCommand>> result = new HashSet<Class<? extends WorkflowCommand>>();
      for (Class<? extends WorkflowCommand> commandClass : allCommandClasses) {
         if (commandClass.isAnnotationPresent(WorkflowCommandAnnotation.class)) {
            result.add(commandClass);
         }
      }

      return result;
   }

   // TODO: Build a single private findClasses command, which is used by the other finders.
   // TODO: Cache the loaded classes, so large searches are not repeated.
   // TODO: Fine grain the package scopes.
}
