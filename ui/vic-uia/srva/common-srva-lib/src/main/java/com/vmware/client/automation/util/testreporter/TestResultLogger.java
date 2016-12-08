/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.client.automation.util.testreporter;

/**
 * Interface for test logging to external systems
 *
 */
public interface TestResultLogger {

    /**
     * Connects to the external system for test logging.
     *
     * @return
     */
    public boolean connect();

    /**
     * Opens a new test run. It usually contains one or more testcases.
     *
     * @param testSetSpec TestSetSpec providing info for the test set to be executed
     * @return id of the created test set
     */
    public String startTestSet(TestSetSpec testSetSpec);

    /**
     * Ends a test set
     *
     * @return true if it successfully finished
     */
    public boolean endTestSet();

    /**
     * Logs a single test comment
     *
     * @param comment
     */
    public void logTestComment(String comment);

    /**
     * Logs a verification point comment, containing actual value and expected value
     *
     * @param description - verification description
     * @param actualValue - actual value received
     * @param expectedValue - expected value
     */
    public void logTestVerification(final String description, final String actualValue, final String expectedValue);


    /**
     * Starts test case logging
     *
     * @param methodName - name of the test case method
     * @param packageName - test case package
     * @param description - test case description
     * @return test case id
     */
    public String startTestCase(final String methodName, final String packageName, final String description);

    /**
     * Ends test case logging
     *
     * @param testCaseResult - final result of the test (passed, failed, etc.)
     */
    public void endTestCase(final String testCaseResult);

    /**
     * Uploads a screenshot file to the external logging system by passing the file name and path
     *
     * @param screenshotFile - full path and file name
     */
    public void postScreenshot(final String screenshotFile);

    /**
     * Uploads a screenshot file to the external logging system, using the default path for storing screenshots
     */
    public void postScreenshot();

    /**
     * Gets connection state
     * @return true if connection is established
     */
    public boolean getConnected();

    /**
     * Gets the current test case id
     *
     * @return test case id
     */
    public String getTestCaseId();
}
