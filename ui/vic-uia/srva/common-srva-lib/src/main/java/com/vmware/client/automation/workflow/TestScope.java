/** Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow;

/**
 * The test scope defines what kind of verifications are performed and paths
 * taken respectively while the test is executed.
 *
 * A single test might be designed to handle more than a single scope.
 * It is encouraged that a test performs multiple checks without doing a lot
 * of traversals over the screen.
 *
 * Currently the available scopes are BAT, SAT, and UI. More scopes can
 * be added in the future.
 */
public enum TestScope {

   /**
    * The MINIMAL scope includes test steps that are completely mandatory to be
    * executed in order the test to be executed correctly.
    *
    * This scope tests performs only minimal set of verifications. In most cases it will
    * verify the most important points. All prerequisite steps have to be marked with
    * this test scope in order the test to start executing.
    */
   MINIMAL(1),

   /**
    * BAT scope includes testing of straight scenarios with strictly
    * minimal verification.
    *
    * An example of such test is creating a VDC through the UI and verifying
    * only through the API that the VDC is created. No verifications through
    * the UI should be made. This will ensure in the fastest possible way
    * that the basic functionality works without additionally checking if
    * the automatic refresh will be triggered. The latter has a workaround
    * and couldn't be considered a blocker.
    */
   BAT(10),

   /**
    * SAT scope includes testing of straight scenarios with verification
    * as defined in per the DoD (Definition of Done) criteria.
    *
    * An example of such test is creating a VDC through the UI and verifying
    * both through the API and the UI that the VDC is created. Doing both of
    * this is important as it will provide clear information in the triaging
    * phase about the cause and the importance of the potentially discovered
    * problem.
    */
   SAT(20),

   /**
    * UI scope includes testing of all kinds of scenarios and performing all
    * kinds of verifications.
    */
   UI(30),

   /**
    * FULL scope includes all previous scopes. If a scope is not specified
    * a test should be run in this scope.
    */
   FULL(10000);



   private int _scopeNumber;

   /**
    * Initialize the <code>TestRunScope</code> with a number.
    *
    * A higher number also includes scopes defined with lower numbers.
    *
    * @param scopeNumber
    *    Scope number.
    */
   TestScope(int scopeNumber) {
      _scopeNumber  = scopeNumber;
   }

   /**
    * Scope number.
    */
   public int getScopeNumber() {
      return _scopeNumber;
   }
}
