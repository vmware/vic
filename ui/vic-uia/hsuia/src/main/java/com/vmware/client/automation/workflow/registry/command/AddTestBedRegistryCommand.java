/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.registry.command;

import java.util.Properties;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.command.CommandException;
import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.explorer.SessionSettingsConstants;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowContext;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;

/**
 * The <code>AddTestBedRegistryCommand</code> class implements a registry command
 * that adds a test bed configuration to the registry.
 *
 * The test bed configuration is provide in a key-value plain text file, which
 * holds of the data for the test bed instance such as network access, authentication,
 * identification of the hosting system and other.
 *
 * Command:
 *    registry-add-testbed path_to_testbed_settings_key_value_file
 *
 * Example:
 *    registry-add-testbed /common-root/unique-location/providerName.settings
 */
@WorkflowCommandAnnotation(commandName = "registry-add-testbed")
public class AddTestBedRegistryCommand extends BaseRegistryCommand {

   protected static final Logger _logger = LoggerFactory.getLogger(AddTestBedRegistryCommand.class);

   private String _testbedSettingsFilePath;

   @Override
   public void prepare(String[] commandParams) {
      validateParameters(
            commandParams,
            new String[] {"testbedSettingsFilePath"});

      _testbedSettingsFilePath = commandParams[0];
      _logger.info("Testbed settings file: " + _testbedSettingsFilePath);
   }

   @SuppressWarnings("unchecked")
   @Override
   public void execute() throws CommandException {
      _logger.info("Register testbed from: " + _testbedSettingsFilePath);
      // Find provider.id
      Properties settings = loadKVSettings(_testbedSettingsFilePath);
      String providerClassName = settings.getProperty(SessionSettingsConstants.KEY_PROVIDER_ID);
      if (Strings.isNullOrEmpty(providerClassName)) {
         throw new IllegalArgumentException(
               String.format("Required key %s not set in %s",
                     SessionSettingsConstants.KEY_PROVIDER_ID,
                     _testbedSettingsFilePath));
      }

      // Validate the provider exists -> registry
      ProviderWorkflowContext context = null;
      try {
         // Find provider class
         Class<ProviderWorkflow> providerClass =
               (Class<ProviderWorkflow>)getRegistry().getRegisteredWorkflowClass(
                     providerClassName, ProviderWorkflow.class);

         // Get context for the provider operation
         context = getRegistry().registrerWorkflowContext(providerClass);
         ProviderWorkflowController controller =
               ProviderWorkflowController.create(context, _testbedSettingsFilePath);

         // Register testbed
         controller.register();
      } catch (ClassNotFoundException e) {
         throw new CommandException(
               String.format("Command cannot find provider class: %s", providerClassName),e);
      } catch (ProviderWorkflowException | ProviderControllerException e1) {
         throw new CommandException(
               String.format("Could not register testbed: %s", providerClassName),e1);
      } finally {
         getRegistry().unregistrerWorkflowContext(context);
      }
   }
}
