/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.command;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.common.WorkflowController;
import com.vmware.client.automation.workflow.explorer.WorkflowBrowser;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;

/**
 * The controller executes a series of commands a provided in the batch.
 */
public class CommandController extends WorkflowController {

   private List<String> _commandsBatch;

   private Map<String, Class<? extends WorkflowCommand>> _commandsClassesMap;

   private List<WorkflowCommand> _commandsToExecute;

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(CommandController.class);

   @Override
   public WorkflowRegistry getRegistry() {
      // TODO Auto-generated method stub
      return null;
   }

   /**
    * Load the command batch to be executed and initialize all executable
    * commands in the project.
    * @param commandStringsBatch
    *    list of commands to be executed by the command controller.
    */
   public void initialize(List<String> commandStringsBatch) {
      // Validate the input is set
      if (commandStringsBatch == null || commandStringsBatch.size() == 0) {
         throw new IllegalArgumentException(
               "No batch of commands is provided. At least one command must be added.");
      }

      // Get a copy of the command strings.
      _commandsBatch = new ArrayList<String>();
      for (String commandString : commandStringsBatch) {
         _commandsBatch.add(commandString);
      }

      // Load all the commands known in the system.
      Set<Class<? extends WorkflowCommand>> workflowCommands = WorkflowBrowser
            .findExecutableCommands();
      _commandsClassesMap = new HashMap<String, Class<? extends WorkflowCommand>>();
      for (Class<? extends WorkflowCommand> commandClass : workflowCommands) {
         _commandsClassesMap.put(
               commandClass.getAnnotation(WorkflowCommandAnnotation.class)
                     .commandName(), commandClass);
      }
   }

   /**
    * Parse each command and the provided parameters.
    * After that invoke the command prepare method.
    * @throws CommandException
    */
   public void prepare() throws CommandException {
      _commandsToExecute = new ArrayList<WorkflowCommand>();

      // split the commands batch into commands.
      for (String commandString : _commandsBatch) {
         if (Strings.isNullOrEmpty(commandString.trim())) {
            // Skip if empty line
            continue;
         }

         // Get command
         String[] params = commandString.split(" ");
         String commandName = params[0].toLowerCase();
         if (commandName.startsWith("#")) {
            continue;
         }

         // Find command
         if (!_commandsClassesMap.containsKey(commandName)) {
            throw new CommandException(
                  String.format(
                        "%s is not recognized command. No commands from this batch will be executed.",
                        commandName));
         }

         // Load command
         Class<? extends WorkflowCommand> commandClass =
               _commandsClassesMap.get(commandName);
         WorkflowCommand commandInstace = null;
         try {
            commandInstace = commandClass.newInstance();
         } catch (Exception e) {
            throw new CommandException(
                  String.format(
                        "Cannot get an instance of command %s. No commands from this batch will be executed.",
                        commandName), e);
         }

         // Prepare command
         String[] commandParams = Arrays.copyOfRange(params, 1, params.length);
         try {
            commandInstace.prepare(commandParams);
         } catch (Exception e) {
            throw new CommandException(
                  String.format(
                        "Cannot validate input parameters of command %s. No commands from this batch will executed",
                        commandName), e);
         }

         _commandsToExecute.add(commandInstace);
      }
   }

   /**
    * Execute the batch of commands in the order they are provided.
    * @throws Exception
    */
   public void execute() throws Exception {
      for (WorkflowCommand command : _commandsToExecute) {
         try {
            _logger.info("============= Execute CMD: " + command.getCommandName());
            command.execute();
            _logger.info("============= END CMD: " + command.getCommandName());
         } catch (Exception e) {
            _logger.error("============= ERROR During Execute CMD: "
                  + command.getCommandName());
            throw new Exception("Command Error", e);
         }
      }
   }
}
