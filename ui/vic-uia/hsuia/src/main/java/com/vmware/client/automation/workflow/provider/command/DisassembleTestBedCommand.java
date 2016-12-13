/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider.command;

import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowException;

/**
 * This class implements command that calls the specified provider for
 * destruction of a test bed.
 * 
 * For more details on how test bed commands work, @see BaseTestBedCommand
 * 
 * Command:
 *    testbed-disassemble package.provider_class path_to_testbed_settings_key_value_file
 *    
 * If the command is successful, the test bed settings will be emptied. On failure to
 * cleanly destroy test bed, its settings file will be left intact.
 *
 * Example:
 *    testbed-disassemble com.vmware.vsphere.client.automation.provider.HostProvider /common-root/unique-location/providerName.settings
 */
@WorkflowCommandAnnotation(commandName = "testbed-disassemble")
public class DisassembleTestBedCommand extends BaseTestBedCommand {

   @Override
   /**
    * Validates the provided test bed settings file can be read
    */
   protected void runCustomValidation() {
      validateInputFile(testbedSettingsFilePath);
   };

   @Override
   /**
    * Calls the disassembling API of the controller.
    */
   protected void runController(ProviderWorkflowController controller)
         throws ProviderControllerException, ProviderWorkflowException {
      controller.disassemble();   
   }
}
