/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.util.testreporter;


/**
 * Connection details for the external test logger system
 *
 */
public interface TestLoggerConnectSpec {

   /**
    * Returns the URL of the host where test result logger application resided.
    *
    * @return     the logger page URL
    */
   public String getTestLoggerURL();

   /**
    * Get threaded logging state.
    *
    * @return     true if threaded logging is enabled
    */
   public boolean getThreadedLogging();
}
