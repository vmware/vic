/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.assertions.Assertion;
import com.vmware.client.automation.assertions.AssertionUtil;
import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.workflow.test.TestWorkflowStep;

/**
 * Base workflow step class. All steps should extend it.
 * Its goal is to provide step common functionality - verify methods,
 * test scope and etc.
 */
public abstract class CommonWorkflowStep implements TestWorkflowStep {

   private static final Logger _logger =
         LoggerFactory.getLogger(CommonWorkflowStep.class);

   public static final String FAIL_VERIFICATION_PREFFIX = "FAILED: ";
   public static final String PASS_VERIFICATION_PREFFIX = "PASSED: ";

   // These packages do not provide information for the validation failure
   // and will be removed from the stack trace
   private static final String WORKFLOW_PACKAGE = "com.vmware.client.automation.workflow";
   private static final String TESTNG_PACKAGE = "org.testng";

   // List of non-fatal errors found during the step execution
   private final List<RuntimeException> notFatalErrorsList =
         new ArrayList<RuntimeException>();

   // Default step scope
   private TestScope _testScope = TestScope.FULL;

   @Override
   public List<RuntimeException> getFailedValidations() {
      return notFatalErrorsList;
   }

   @Override
   public void setStepTestScope(TestScope testScope) {
      _testScope = testScope;
   }

   @Override
   public TestScope getStepTestScope() {
      return _testScope;
   }

   // ---------------------------------------------------------------------------
   // Verify Fatal

   /**
    * Check if expected condition is met. If not an exception will be thrown
    * and the test execution will be stopped. The result is logged.
    *
    * Use this method to implement a check that would require accessing UI
    * object state that is already loaded by the automation test. As the object
    * state is already available, this will not affect the execution time in case
    * the verification is not required by the test scope.
    *
    * @param requiredTestScope   the required test scope
    * @param verifiedCondition   a condition to check
    * @param message             description of the performed verification
    * @return                    true - if the verification has passed successfully
    */
   protected boolean verifyFatal(TestScope requiredTestScope, boolean verifiedCondition,
         String message) {
      final boolean res = verifiedCondition;
      return verifyFatal(
            requiredTestScope,
            new TestScopeVerification() {
               @Override
               public boolean verify() {
                  return res;
               }
            },
            message
         );
   }

   /**
    * Check if expected condition is met. If not an exception will be thrown
    * and the test execution will be stopped. The result is logged.
    *
    * @param verifiedCondition   a condition to check
    * @param message             description of the performed verification
    * @return                    true - if the verification has passed successfully
    */
   protected boolean verifyFatal(boolean verifiedCondition, String message) {
      return verifyFatal(TestScope.MINIMAL, verifiedCondition, message);
   }

   /**
    * Check if expected condition is met. If not an exception will be thrown
    * and the test execution will be stopped. The result is logged.
    *
    * Use this method to implement a check that would require accessing a UI
    * object that is not loaded by the automation test yet. The latter will
    * be performed only if the required test scope is satisfied. Thus the
    * overall test execution time will be reduced.
    *
    * @param requiredTestScope   the required test scope
    * @param verification        a verification to be checked
    * @param message             description of the performed verification
    * @return                    true - if the verification has passed successfully
    */
   protected boolean verifyFatal(TestScope requiredTestScope,
         TestScopeVerification verification, String message) {
      return verify(requiredTestScope, null, verification, message, false);
   }

   /**
    * Check if expected condition is met. If not an exception will be thrown
    * and the test execution will be stopped. The result is logged.
    *
    * @param verification  a verification to be checked
    * @param message       description of the performed verification
    * @return              true - if the verification has passed successfully
    */
   protected boolean verifyFatal(TestScopeVerification verification, String message) {
      return verify(TestScope.MINIMAL, null, verification, message, false);
   }


   /**
    * Check if expected condition is met. If not an exception will be thrown
    * and the test execution will be stopped. The result is logged.
    *
    * @param requiredTestScope      the required test scope
    * @param assertion              a condition to check
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verifyFatal(TestScope requiredTestScope, Assertion assertion) {
      return verify(requiredTestScope, assertion, null, null, false);
   }

   /**
    * Check if expected condition is met. If not an exception will be thrown
    * and the test execution will be stopped. The result is logged.
    *
    * @param assertion              a condition to check
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verifyFatal(Assertion assertion) {
      return verifyFatal(TestScope.MINIMAL, assertion);
   }

   // ---------------------------------------------------------------------------
   // Verify Safely

   /**
    * Check if expected condition is met. The result is logged.
    *
    * The verification is performed only if the required test scope is
    * satisfied.
    *
    * Use this method to implement a check that would require accessing a UI
    * object that is not loaded by the automation test yet. The latter will
    * be performed only if the required test scope is satisfied. Thus the
    * overall test execution time will be reduced.
    *
    * @param requiredTestScope   the required test scope
    * @param verification        a verification to be checked
    * @param message             description of the performed verification
    * @return                    true - if the verification has passed successfully
    */
   protected boolean verifySafely(TestScope requiredTestScope,
         TestScopeVerification verification, String message) {
      return verify(requiredTestScope, null, verification, message, true);
   }

   /**
    * Check if expected condition is met. The result is logged.
    *
    * @param verification     a verification to be performed
    * @param message          description of the performed verification
    * @return                 true - if the verification has passed successfully
    */
   protected boolean verifySafely(TestScopeVerification verification, String message) {
      return verify(TestScope.MINIMAL, null, verification, message, true);
   }

   /**
    * Check if expected condition is met. The result is logged.
    *
    * The actual verification should be done before this method is called.
    *
    * Use this method to implement a check that would require accessing UI
    * object state that is already loaded by the automation test. As the object
    * state is already available, this will not affect the execution time in case
    * the verification is not required by the test scope.
    *
    * @param requiredTestScope      the required test scope
    * @param verifiedCondition      a condition to check
    * @param message                description of the performed verification
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verifySafely(TestScope requiredTestScope,
         boolean verifiedCondition, String message) {
      final boolean res = verifiedCondition;
      return verifySafely(requiredTestScope, new TestScopeVerification() {
         @Override
         public boolean verify() {
            return res;
         }
      }, message);
   }

   /**
    * Check if expected condition is met. The result is logged.
    *
    * @param verifiedCondition      a condition to check
    * @param message                description of the performed verification
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verifySafely(boolean verifiedCondition, String message) {
      return verifySafely(TestScope.MINIMAL, verifiedCondition, message);
   }

   /**
    * Check if expected condition is met. The result is logged.
    *
    * @param requiredTestScope      the required test scope
    * @param assertion              a condition to check
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verifySafely(TestScope requiredTestScope, Assertion assertion) {
      return verify(requiredTestScope, assertion, null, null, true);
   }

   /**
    * Check if expected condition is met. The result is logged.
    *
    * @param assertion              a condition to check
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verifySafely(Assertion assertion) {
      return verifySafely(TestScope.MINIMAL, assertion);
   }

   /**
    * Helper method that performs the actual verification work of all
    * 'verify' methods.
    *
    * A verification is performed if expected condition is met. Based
    * on the <code>verifySafely</code> parameter, if the condition is
    * not met, exception might be thrown. The result is always
    * logged.
    *
    * The verification is performed only if the required test scope is
    * satisfied.
    *
    * @param requiredTestScope      the required test scope
    * @param assertion              assertion to verify. Either this or
    *                               verification should be passed
    * @param verification           a verification to be performed. Either this or
    *                               assertion should be passed
    * @param message                description of the performed verification.
    * @param verifySafely           If true, the result will be just logged.
    *                               If false, an exception will be thrown.
    * @return                       true - if the verification has passed successfully
    */
   protected boolean verify(TestScope requiredTestScope, Assertion assertion,
         TestScopeVerification verification, String message, boolean verifySafely) {
      assertVerificationParams(assertion, verification);

      // check if the scope of the verification is satisfied
      if (hasTestScope(requiredTestScope)) {
         boolean result;

         // if TestScopeVerification is used convert it to EqualsAssertion
         if (verification != null) {
            assertion = convertToAssertion(verification, message);
         }

         // perform the verification
         result = performVerification(assertion, verifySafely);

         return result;
      } else {
         return true;
      }
   }

   /**
    * Perform a verification.
    *
    * @param assertion        verification that will be performed
    * @param verifySafely     if it is safe or fatal one
    * @return                 the result of the verification
    */
   protected boolean performVerification(Assertion assertion, boolean verifySafely) {
      boolean result;

      try {
         if (verifySafely) {
            result = AssertionUtil.assertSafely(assertion);
            if(!result) {
              // Save the exception for later escalation
               RuntimeException verificationEx =
                     createVerificationException(assertion.getDescription());
               notFatalErrorsList.add(verificationEx);
            }
         } else {
            result = AssertionUtil.assertFatal(assertion);
         }
      } catch (AssertionError ae) {
         result = false;
         throw createVerificationException(ae.getMessage());
      }

      return result;
   }

   /**
    * Indicate if the step is running on a test scope that satisfies
    * (is higher or equal to) the required one.
    *
    * Use this method to determine if a code path that requires specific
    * test scope should be executed.
    *
    * @param requiredTestScope      the required test scope
    * @return                       true - if the step test scope is greater or
    *                               equal to the required one
    */
   protected boolean hasTestScope(TestScope requiredTestScope) {
      return requiredTestScope.getScopeNumber() <= getStepTestScope().getScopeNumber();
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Create RuntimeException with the specified error message.
    *
    * As addition the method removes from the exception the stack traces that are
    * for the WORKFLOW_PACKAGE and TESTNG_PACKAGE as they does not provide any
    * useful information for the test validation fails.
    *
    * @param errorMessage     exception message.
    * @return                 RuntimeException for the specified message
    */
   private RuntimeException createVerificationException(String errorMessage) {
      RuntimeException result = new RuntimeException(errorMessage);
      StackTraceElement[] stackTraceElements = result.getStackTrace();
      List<StackTraceElement> editedStackTraces = new ArrayList<StackTraceElement>();
      for (StackTraceElement stackTraceElement : stackTraceElements) {
         if(!stackTraceElement.getClassName().startsWith(WORKFLOW_PACKAGE) &&
            !stackTraceElement.getClassName().startsWith(TESTNG_PACKAGE) ) {
            editedStackTraces.add(stackTraceElement);
         }
      }

      result.setStackTrace(
            editedStackTraces.toArray(new StackTraceElement[editedStackTraces.size()]));
      return result;
   }

   private void assertVerificationParams(Assertion assertion,
         TestScopeVerification verification) {
      if (!(assertion == null ^ verification == null)) {
         throw new IllegalArgumentException("Inappropriate assertion was passed");
      }
   }

   /**
    * Converts TestScopeVerification to Assertion.
    *
    * @param verification        verification to be converted
    * @return                    converted Assertion
    * @throws RuntimeException   if the conversion failed
    */
   private Assertion convertToAssertion(TestScopeVerification verification,
         String message) throws RuntimeException {
      try {
         return new EqualsAssertion(verification.verify(), true, message);
      } catch (Exception e) {
         _logger.error(String.format(
               "Error occurred while performing verification: %s\n%s",
               e.getMessage(),
               e.getStackTrace()));
         throw new RuntimeException(e);
      }
   }
}
