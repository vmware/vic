/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider.command;

import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;

/**
 * This class implements command that calls the specified provider for
 * construction of a test bed. If the construction is not successful, the provider
 * will attempt to destroy the test bed back from the point where it was successful.
 * 
 * For more details on how test bed commands work, @see BaseTestBedCommand
 * 
 * Command:
 *    testbed-assemble package.provider_class path_to_testbed_settings_key_value_file
 *    
 * If the command is successful, the test bed settings will remain written
 * in the path_to_testbed_settings_key_value_file after the command's execution.
 * 
 * On failure to cleanly destroy test bed, which was not successfully constructed,
 * the intermittent content of the test bed settings file will be left intact.
 *
 * Example:
 *    testbed-assemble com.vmware.vsphere.client.automation.provider.HostProvider /common-root/unique-location/providerName.settings
 */
@WorkflowCommandAnnotation(commandName = "testbed-assemble")
public class AssembleTestBedCommand extends BaseTestBedCommand {

   @Override
   /**
    * Validates the file path is writable.
    */
   protected void runCustomValidation() {
      createAndValidateOutputFile(testbedSettingsFilePath);
   };

   @Override
   /**
    * Calls the assembling API of the controller.
    */
   protected void runController(ProviderWorkflowController controller)
         throws ProviderControllerException, ProviderWorkflowException {
      controller.assemble();   
   }
}
