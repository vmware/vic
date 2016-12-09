/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.util.testreporter;

/**
 * Contains specification for the test set which is executed, including browser
 * OS, description of the test set, build number, branch, etc.
 *
 */
public interface TestSetSpec {

   /**
    * Returns the browser OS where the test is executing
    *
    * @return Concatenated String with OS name, OS architecture and JRE
    *         architecture
    */
   public String getBrowserOs();

   /**
    * Returns type of the browser used for test execution
    *
    * @return Browser type
    */
   public String getBrowser();

   /**
    * Returns the test owner of the executed test set
    *
    * @return test owner name
    */
   public String getTestOwner();

   /**
    * Returns the application under test name
    *
    * @return product name
    */
   public String getProductName();

   /**
    * Returns the build number of the application under test
    *
    * @return build number of the AUT
    */
   public String getBuildNumber();

   /**
    * Returns the build type of the application (beta, obj)
    *
    * @return build type
    */
   public String getBuildType();

   /**
    * Returns the test type for the run (regression, l18n, etc.)
    *
    * @return test type
    */
   public String getTestType();

   /**
    * Gets the test set description
    *
    * @return test set description
    */
   public String getTestSetDescription();

   /**
    * Returns the application branch used for the test run (main, rel, etc.)
    *
    * @return branch name
    */
   public String getBranch();

   /**
    * Returns the display language of the OS where the test runs
    *
    * @return display language on the OS where test set is executed
    */
   public String getLanguage();

   /**
    * Gets existing result id used for replacing/appending results from current
    * run to existing result id from previous run
    *
    * @return existing result id
    */
   public String getExistingResultId();
}
