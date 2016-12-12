/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider.command;

import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;

/**
 * This class implements command that calls the specified provider for verification
 * of a test bed's health. No modification of the test bed will be made.
 * 
 * For more details on how test bed commands work, @see BaseTestBedCommand
 * 
 * Command:
 *    testbed-checkhealth package.provider_class path_to_testbed_settings_key_value_file
 *
 * Example:
 *    testbed-checkhealth com.vmware.vsphere.client.automation.provider.HostProvider /common-root/unique-location/providerName.settings
 */
@WorkflowCommandAnnotation(commandName = "testbed-checkhealth")
public class CheckTestBedCommand extends BaseTestBedCommand {

   @Override
   /**
    * Validates the provided test bed settings file can be read.
    */
   protected void runCustomValidation() {
      validateInputFile(testbedSettingsFilePath);
   };

   @Override
   /**
    * Calls the health checking API of the controller.
    */
   protected void runController(ProviderWorkflowController controller)
         throws ProviderControllerException, ProviderWorkflowException {
      controller.checkHealth();   
   }
}
