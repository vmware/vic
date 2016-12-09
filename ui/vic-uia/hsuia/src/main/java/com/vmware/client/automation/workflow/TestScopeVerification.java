/** Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow;

/**
  * The interface defines a verification method. Implementor classes should
  * use it to perform a single check only if the test scope of the step/test is
  * greater than the test scope of the verification..
  *
  * Making checks this way lets the framework decide if a check should be
  * made or not based on common runtime criteria.
 */
public interface TestScopeVerification {

   /**
    * Implement the check here.
    *
    * @return
    *    True - the verification passed successfully,
    *    False - otherwise
    *
    */
   public boolean verify() throws Exception;
}
