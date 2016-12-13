package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CONFIRMATION_DIALOG;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CONFIRM_YES_LABEL;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_ERROR_DIALOG;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_ERROR_WARNING_OK;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_WARNING_DIALOG;
import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;
import static com.vmware.flexui.selenium.BrowserUtil.selenium;

import java.text.SimpleDateFormat;
import java.util.Calendar;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.Assert;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.flexui.componentframework.controls.mx.Alert;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.flexui.selenium.MethodCallUtil;

/**
 * Test base class from VCUI-QE-LIB.
 *
 * NOTE: this class is a partial copy of the one from VCUI-QE-LIB
 */
public class TestBaseUI {

   private static final Logger _logger = LoggerFactory.getLogger(TestBaseUI.class);

   // The reason why we have separate AMF Connection here is because
   // Login is using BaseState without logging in, therefore it has
   // to establish its own AMFConnection
   protected static String serverName = null;
   protected static String userName = null;
   protected static String password = null;
   public static String multiVcStr = null;
   public static String[] multiVcArr = null;
   protected String testId = null;
   protected static String loginAppName = null; // Name of App Under Test
   protected static String clientAppName = null;
   protected String jsessionId = null;
   protected static String useSSL = null;
   protected static String clientKeyStroke = null;
   protected static String clientUrl = null;
   protected static String amfURL = null;
   protected static boolean remoteTest = true;
   public static boolean errorFlag = false;
   public static AssertionError firstAssertionError = null;
   public static String firstScreenshot = null;
   public static StringBuffer errorLog = new StringBuffer();
   public static int browser_pixel_diff_x = 0;
   public static int browser_pixel_diff_y = 0;
   public static long default_timeout = 0;
   public static String ESXVersion = null;
   public static String LEGACY_HOST_IP = null;
   public static String esx_password = null;
   public static String path_nas = null;
   public static String ip_nas = null;
   public static String DATE_PATTERN = null;
   protected static int maxRetriesLogout = 10;
   protected static boolean verifyInventoryBaseStateBeforeClass = false;
   protected static boolean verifyInventoryBaseStateAfterMethod = false;
   protected int testCounter = 0;
   protected static int reloginConst = 0;
   protected static int resetClientConst = 0;

   // to find and print elapsed time
   protected long startTime;
   protected long endTime;
   protected long elapsedTime;
   protected Calendar calendar;
   protected SimpleDateFormat sdf;
   protected static final String DATE_FORMAT_NOW = "yyyy/MM/dd HH:mm:ss";

   // TODO move these to TestConstants
   private static final String BROWSER_TYPE_FIREFOX = "firefox";
   private static final String BROWSER_TYPE_IEXPLORE = "iexplore";
   private static final String CERTIFICATE_ERROR_TITLE_IE = "Certificate Error";
   private static final String TEXT_ID = "id";
   private static final String EQUALS_SIGN = "=";
   private static final String ID_OVERRIDE_CERT_LINK_BUTTON_IE = "overridelink";

   private static String currentTestId = null;

   /**
    * Verify Test Result and print the Report Log accordingly.
    *
    * @param actualValue Object
    * @param expectedValue Object
    * @param description String
    * @param testId String - prefix for the screenshot to take
    */
   public static void verifySafely(Object actualValue, Object expectedValue,
         String description) {
      verifySafely(actualValue, expectedValue, description, null);
   }

   /**
    * A safe wrapper for AssertTrue (from Junit or TestNg). This will log the
    * result of the operation.
    * @param actual
    * @param description
    */
   public static void verifyTrueSafely(Boolean actual, String description) {
      verifySafely(actual, Boolean.TRUE, description);
   }

   /**
    * A safe wrapper for AssertFalse (from Junit or TestNg). This will log the
    * result of the operation.
    * @param actual
    * @param description
    */
   public static void verifyFalseSafely(Boolean actual, String description) {
      verifySafely(actual, Boolean.FALSE, description);
   }

   /**
    * A safe wrapper for AssertNotNull (from Junit or TestNg). This will log the
    * result of the operation.
    * @param actual
    * @param description
    */
   public static void verifyNotNullSafely(Object actual, String description) {
      verifyTrueSafely(actual != null, description);
   }

   /**
    * A safe wrapper for AssertTrue (from Junit or TestNg). This will log the
    * result of the operation. If this fails the whole test will be failed.
    * @param actual
    * @param description
    */
   public static void verifyTrueFatal(Boolean actual, String description) {
      verifyFatal(actual, Boolean.TRUE, description);
   }

   /**
    * A safe wrapper for AssertFalse (from Junit or TestNg). This will log the
    * result of the operation. If this fails the whole test will be failed.
    * @param actual
    * @param description
    */
   public static void verifyFalseFatal(Boolean actual, String description) {
      verifyFatal(actual, Boolean.FALSE, description);
   }

   /**
    * A safe wrapper for AssertNotNull (from Junit or TestNg). This will log the
    * result of the operation. If this fails the whole test will be failed.
    * @param actual
    * @param description
    */
   public static void verifyNotNullFatal(Object actual, String description) {
      verifyTrueFatal(actual != null, description);
   }

   /**
    * Verify Test Result and print the Report Log accordingly. Private method
    * that also sets the prefix for the captured screenshot
    *
    * @param actualValue Object
    * @param expectedValue Object
    * @param description String
    * @param testId String - prefix for the screenshot to take
    */
   public static void verifySafely(Object actualValue, Object expectedValue,
         String description, String testId) {
      try {
         Assert.assertEquals(actualValue, expectedValue, description);
         _logger.info("Success: " + description + ". value: " + actualValue);

      } catch (AssertionError ae) {
         _logger.error("Failed: " + description + ". actual: " + actualValue + ", expected: "
               + expectedValue);

         errorFlag = true;
         // Keep the first exception that occurred so that we re-throw it in
         // throwAssertionErrorOnFailure()
         if (firstAssertionError == null) {
            firstAssertionError = ae;
         }
         _logger.error(ae.getMessage(), ae);
         String snapShot = captureSnapShot(testId);
         // Keep the screenshot for the first exception that occurred so that we
         // add it to the message in throwAssertionErrorOnFailure()
         if (firstScreenshot == null) {
            firstScreenshot = snapShot;
         }
         errorLog.append("\n" + description + ": " + ae.getMessage() + " : " + snapShot);
      }
   }

   /**
    * Throws the exception safely, i.e. throws when throwAssertionErrorOnFailure() is
    * called. For use in catch blocks.
    *
    * @param exception Throwable object
    * @param testId String prefix for the screenshot to capture
    */
   public static void throwSafely(Throwable exception, String testId) {
      if (exception == null) {
         throw new AssertionError("Exception passed to throwSafely is null");
      }
      if (firstAssertionError == null) {
         firstAssertionError = new AssertionError(exception.getMessage());
         firstAssertionError.setStackTrace(exception.getStackTrace());
      }
      String exceptionMessage = exception.getMessage();
      if (exceptionMessage == null) {
         exceptionMessage = "Blank exception message";
      }
      _logger.error("Throw safely: " + exceptionMessage, exception);
      verifySafely(true, false, exceptionMessage, testId);
   }

   /**
    * Verify test Results and raise AssertionError on Failure.
    * Similar to verifySafely, but will not continue on Failures.
    *
    * @param actual
    * @param expected
    * @param description
    */
   public static void verifyFatal(Object actual, Object expected, String description) {
      // In order for verifyFatal() to do only the current verification we reset
      // the errorFlag before doing the verification.
      // Without this if any of the *previous* verifications done using
      // verifySafely() has failed and had set the errorFlag then verifyFatal()
      // will also fail, but that will be wrong as we want to fail only if the
      // *current* verification done in verifyFatal() fails
      boolean isErrorFlagSet = errorFlag;
      errorFlag = false;
      verifySafely(actual, expected, description);
      throwAssertionErrorOnFailure();
      // The following won't be reached if the above verifySafely() failed and
      // throwAssertionErrorOnFailure() throws an AssertionError
      errorFlag = isErrorFlagSet;
   }

   /**
    * This method detects and dismisses modal dialog and
    * throws AssertionError if errorFlag is true.
    * This method should be call at the end of each test case.
    * @throws AssertionError
    */
   public static void throwAssertionErrorOnFailure() throws AssertionError {
      // implicitly detect and dismiss modal dialog at the end of each test to
      // report the test as failure if there is an unexpected modal dialog
      try {
         if (detectAndHandleModalDialog()) {
            errorFlag = true;
            errorLog.append("ModalDialog detected");
         }
      } catch (Throwable e) {
         _logger.error("Caught exception while detecting and handling the modal "
               + "dialog: " + e.getMessage(), e);
      }
      if (errorFlag) {
         if (firstAssertionError != null) {
            AssertionError ae =
                  new AssertionError(firstAssertionError.getMessage() + " : "
                        + firstScreenshot);
            ae.setStackTrace(firstAssertionError.getStackTrace());
            _logger.error("Contents of errorLog: " + errorLog);
            // Throw the first assertion error that occurred so that TestNG
            // trace shows the trace for the first test issue instead of the
            // stacktrace of the throwAssertionErrorOnFailure() method
            throw ae;
         } else {
            throw new AssertionError(errorLog);
         }
      }
   }

   /**
    * Wrapper call to the BrowserUtil.captureSnapshot()
    * All calls within VCUIQA-FLEX-UI should use this function and
    * not call BrowserUtil.captureSnapshot() directly.
    *
    * @param String sName - ScreenShot Name
    *
    * @return String - ScreenShot Location
    */
   public static String captureSnapShot(String sName) {
      String scrnShot = BrowserUtil.captureSnapshot(sName, "test-output/screenShots");
      _logger.info("viClient ScreenShot captured : " + scrnShot);
      return scrnShot;
   }


   /**
    * Helper Method to detect error/warning/confirmation modal dialog and dismiss it.
    * Implicitly call in throwAssertionErrorOnFailure which will be called only
    * once per test.
    *
    */
   private static boolean detectAndHandleModalDialog() throws AssertionError {
      boolean isModalDialogDetected = false;
      Button okButton = null;
      Alert modalDialog = null;
      _logger.info("Checking for Modal dialog...");
      try {
         flashSelenium = new FlashSelenium(selenium, TestBaseUI.clientAppName);
         if (MethodCallUtil.getVisibleOnPath(flashSelenium, ID_ERROR_DIALOG)) {
            isModalDialogDetected = true;
            _logger.info("Flex error dialog detected!");
            modalDialog = new Alert(ID_ERROR_DIALOG, flashSelenium);
            _logger.info("Flex error dialog message reported: " + modalDialog.getText());
            captureSnapShot("flexErrorPopup_");
            // Click yes on the unexpected Flex error dialog
            okButton =
                  new Button(ID_ERROR_DIALOG + "/" + ID_ERROR_WARNING_OK, flashSelenium);
            okButton.click();
         } else if (MethodCallUtil.getVisibleOnPath(flashSelenium, ID_WARNING_DIALOG)) {
            isModalDialogDetected = true;
            _logger.info("Flex warning dialog detected!");
            modalDialog = new Alert(ID_WARNING_DIALOG, flashSelenium);
            _logger.info("Flex warning dialog message reported: " + modalDialog.getText());
            captureSnapShot("flexWarningPopup_");
            // Click yes on the unexpected Flex error dialog
            okButton =
                  new Button(ID_WARNING_DIALOG + "/" + ID_ERROR_WARNING_OK,
                        flashSelenium);
            okButton.click();
         } else if (MethodCallUtil.getVisibleOnPath(
               flashSelenium,
               ID_CONFIRMATION_DIALOG)) {
            isModalDialogDetected = true;
            _logger.info("Flex confirmation dialog detected!");
            modalDialog = new Alert(ID_CONFIRMATION_DIALOG, flashSelenium);
            _logger.info("Flex confirmation dialog message reported: "
                  + modalDialog.getText());
            captureSnapShot("flexConfirmationPopup_");
            // Click yes on the unexpected Flex error dialog
            okButton =
                  new Button(ID_CONFIRMATION_DIALOG + "/" + ID_CONFIRM_YES_LABEL,
                        flashSelenium);
            okButton.click();
         }
      } catch (Exception e) {
         errorLog.append("\n" + "Caught Exception while dismissing modal dialog:"
               + e.getMessage());
         throw new AssertionError(errorLog);
      }
      _logger.info("No modal dialog detected...");
      return isModalDialogDetected;
   }

   /**
    * Returns the test id of the current test case. This is something that
    * usually you shouldn't need to know and is used for hacking, so use this
    * method with caution.
    *
    * @return String testId of the current test case or null if not set
    */
   public static String getCurrentTestId() {
      return currentTestId;
   }
}
