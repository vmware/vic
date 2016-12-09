/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.util.Properties;

import com.vmware.client.automation.workflow.provider.ProviderWorkflow;

/**
 * Interface for initializing a registry.
 */
public interface RegistryInitializationBridge {

   /**
    * Set the global registry settings - i.e. Nimhus details, build numbers and etc.
    * @param settings
    */
   public void setSessionSettings(Properties settings);

   /**
    * Add testbed to the registry.
    * @param providerWorkflow  provider type to which will be mapped the testbed.
    * @param testbedFilePath   path to the tresbed file.
    * @param testbedSettings   loaded from testbed file properties.
    */
   public void addTestBed(ProviderWorkflow providerWorkflow, String testbedFilePath,
         Properties testbedSettings);
}
