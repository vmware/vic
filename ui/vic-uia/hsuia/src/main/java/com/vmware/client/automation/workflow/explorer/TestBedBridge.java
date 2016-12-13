/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import com.vmware.client.automation.workflow.provider.ProviderWorkflow;

/**
 * The interface provides an API for loading test bed resources from tests or
 * test bed providers.
 */
public interface TestBedBridge {

   /**
    * Calling this method will publish a request for test bed configuration and allocate
    * an access point to the <code>EntitySpecs</code> provided in this test bed.
    * 
    * The requested test bed will not be used in other workflows, while it is
    * still used with this one.
    * 
    * @param testbedProviderClass
    * 	Identification of the test bed configuration.
    * 
    * @param isShared
    *    mark if the testbed will be used as common testbed and may be shared
    *    in other tests.
    * 
    * @return
    * 	<code></code> Note: This should be somehow tied to the provider spec.
    */
   public TestbedSpecConsumer requestTestbed(
         Class<? extends ProviderWorkflow> testbedProviderClass, boolean isShared);
}
