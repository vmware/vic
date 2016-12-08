/** Copyright 2014 VMWare, Inc. All rights reserved. -- VMWare Confidential */

package com.vmware.client.automation.util.testreporter.racetrack;


/**
 * This interface is used for all implementation that allows users to interact with
 * the Racetrack database.
 */
public interface IRacetrack {

   /**
    * Log a comment in racetrack for this test.
    *
    * @param comment a test case comment for this test.
    *
    * @throws TestException thrown if anything goes wrong.
    * @deprecated use {@link IRacetrack#testCaseComment(String, String)} instead
    */
   @Deprecated
   void testCaseComment(String comment) throws Exception;

   /**
    * Add verification info to racetrack for this test.
    *
    * @param description The description of the verification performed.
    * @param actual      The actual result of the test
    * @param expected    The expected result of the test
    * @param result      The result of the verification true if it passed.
    *
    * @throws TestException thrown if anything goes wrong.
    * @deprecated use {@link IRacetrack#testCaseVerification(String, String, String, String, boolean)}
    */
   @Deprecated
   public void testCaseVerification(String description, String actual, String expected,
         boolean result) throws Exception;


   /**
    * Log a comment in racetrack for this test.
    *
    * @param comment a test case comment for this test.
    *
    * @throws TestException thrown if anything goes wrong.
    */
   void testCaseComment(String testCaseId, String comment) throws Exception;

   /**
    * Add verification info to racetrack for this test.
    *
    * @param description The description of the verification performed.
    * @param actual      The actual result of the test
    * @param expected    The expected result of the test
    * @param result      The result of the verification true if it passed.
    *
    * @throws TestException thrown if anything goes wrong.
    */
   public void testCaseVerification(String testCaseId, String description,
         String actual, String expected, boolean result) throws Exception;
}
