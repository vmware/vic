/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider.command;

import com.vmware.client.automation.workflow.command.BaseCommand;
import com.vmware.client.automation.workflow.command.CommandException;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowContext;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;

/**
 * The <code>BaseTestBedCommand</code> class is as base class that holds all the common
 * implementation for the test bed commands.
 * 
 * Typically each test bed command requires specification of provider class and test bed
 * settings file.
 * 
 * The provider class is a string representing the canonical name (package + class name)
 * of the provider that shall be used by the command.
 * 
 * The test bed configuration is provide in a key-value plain text file, which
 * holds of the data for the test bed instance such as network access, authentication,
 * identification of the hosting system and other.
 * 
 * The acquisition and validation of these parameters is done in common prepare() method.
 * The mechanism for initializing the provider's controller is also one and the same for
 * each command and is implemented in execute().
 * 
 * Specific command classes should implement their own logic for custom validation (if any
 * by optionally overriding runCustomValidation()) and the calls to the controller's API
 * by implementing runController().
 * 
 * Command:
 *    testbed-[command] package.provider_class path_to_testbed_settings_key_value_file
 *    
 * Example:
 *    testbed-[command] com.vmware.vsphere.client.automation.provider.HostProvider /common-root/unique-location/providerName.settings
 */
public abstract class BaseTestBedCommand extends BaseCommand {
   
   protected String providerClassName;
   protected String testbedSettingsFilePath;

   @Override
   public void prepare(String[] commandParams) {
      validateParameters(
            commandParams,
            new String[] {"providerClassName", "testbedSettingsFilePath"});
       
      providerClassName = commandParams[0];
      testbedSettingsFilePath = commandParams[1];
   }

   @Override
   public void execute() throws CommandException {
      // Get the provider class
   	Class<? extends ProviderWorkflow> providerClass = getProviderClass(providerClassName);
      
      // Run custom validation
      runCustomValidation();

      WorkflowRegistry registry = getRegistry();
      ProviderWorkflowContext context = null;
      ProviderWorkflowController controller = null;
      try {
      //   context = registry.registerWorkflowInstance(providerClassName);
         context = registry.registrerWorkflowContext(providerClass);
         controller = ProviderWorkflowController.create(context, testbedSettingsFilePath);
         
         // Ask the specific command to run the controller's API
         runController(controller);
         
         //  controller.saveReport(_testBedSettingsFilePath); // TODO: Something to store the properties into the file.
      } catch (Exception e) {
         throw new CommandException(
               String.format("Execution failed on command %s", getCommandName()), e);
      } finally {
         // Make sure the context is always properly removed.
         registry.unregistrerWorkflowContext(context);
      }
   }
   
   /**
    * Override this method to perform custom validation of the command parameters.
    */
   protected void runCustomValidation() {
   };
   
   /**
    * Implement this method to perform command specific API calls to the controller.
    * 
    * @param controller
    *       Provider controller
    *
    * @throws ProviderControllerException
    * @throws ProviderWorkflowException
    */
   protected abstract void runController(ProviderWorkflowController controller)
         throws ProviderControllerException, ProviderWorkflowException;
   
   
   @SuppressWarnings("unchecked")
   /**
    * Returns provider class for a given canonical provider name.
    * 
    * @param providerClassName
    * 	Canonical provider class name
    * @return
    * 	Provider class
    * 
    * @throws CommandException
    * 	This exception is thrown when the class cannot be found or is not provider class.
    */
   private Class<? extends ProviderWorkflow> getProviderClass(
   		String providerClassName) throws CommandException {   	
   	try {
	      return 
	      		(Class<? extends ProviderWorkflow>) getRegistry().getRegisteredWorkflowClass(
	      				providerClassName, ProviderWorkflow.class);
      } catch (ClassNotFoundException e) {
      	throw new CommandException(
      			String.format(
      					"Cannot find provider class: %s",
      					providerClassName),
      			e);
      }
   }
}
