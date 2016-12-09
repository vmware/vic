/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.registry.command;

import java.util.Properties;

import com.vmware.client.automation.workflow.command.CommandException;
import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;


/**
 * The <code>SetSessionRegistryCommand</code> class implements a registry command
 * that sets the session (global) readonly key-value settings into the system.
 * 
 * Typically all other commands require prior initialization of the registry,
 * so this one should be called in very beginning of the batch of command.
 * 
 * Command:
 *    registry-set-session path_to_session_key_value_file
 *
 * Example:
 *    registry-set-session /common-root/unique-location/session.settings
 */
@WorkflowCommandAnnotation(commandName = "registry-set-session")
public class SetSessionRegistryCommand extends BaseRegistryCommand {

   private String _sessionSettingsFilePath;

   @Override
   public void prepare(String[] commandParams) {
      validateParameters(
            commandParams,
            new String[] {"sessionSettingsFilePath"});

      _sessionSettingsFilePath = commandParams[0];
   }

   @Override
   public void execute() throws CommandException {
      Properties settings = loadKVSettings(_sessionSettingsFilePath);
      getRegistryInitializer().setSessionSettings(settings);
   }
}
