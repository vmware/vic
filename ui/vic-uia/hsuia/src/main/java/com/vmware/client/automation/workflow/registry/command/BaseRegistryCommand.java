/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.registry.command;

import com.vmware.client.automation.workflow.command.BaseCommand;
import com.vmware.client.automation.workflow.explorer.RegistryInitializationBridge;

/**
* The <code>BaseRegistryCommand</code> implements common utilities
* and validation methods required by most of the registry manipulation commands.
* 
* By extending the <code>BaseCommand</code>, The class implements
* <code>WorkflowCommand</code>, which is required for all classes that must be
* recognized as command classes.
* 
* Concrete registry manipulation command implementations could extend this class
* for convenience, but this is not required by the system.
*/
public abstract class BaseRegistryCommand extends BaseCommand {

   /**
    * Provide reference to a specialized registry interface. 
    * 
    * TODO: Be more specific when the final registry interfaces are set.
    */
   protected RegistryInitializationBridge getRegistryInitializer() {
      return getRegistry();
   }

}
