/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.suitaf;

import static com.vmware.suitaf.SUITA.Factory.aplx;
import static com.vmware.suitaf.util.CommonUtils.stripNonAlphaNumericChars;

import java.io.File;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.AutomationPlatformLinkExt;
import com.vmware.suitaf.apl.Category;
import com.vmware.suitaf.apl.ComponentDumpLevels;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.util.CallStackAnalyser;
import com.vmware.suitaf.util.CallStackAnalyser.Analysis;
import com.vmware.suitaf.util.CallStackAnalyser.CodeZone;
import com.vmware.suitaf.util.Condition;
import com.vmware.suitaf.util.FailureCauseHolder;

/**
 * This is the sub-tool {@link SubToolAudit} with the following specs:
 * <li> <b>Function Type:</b> STATE RETRIEVAL
 * <li> <b>Description:</b> Contains commands that retrieve and log the
 * current state of the UI and of the test-application itself.
 * <li> <b>Based on APL:</b>
 * {@link AutomationPlatformLinkExt#dumpComponents(
 * ComponentMatcher, int, ComponentDumpLevels, Category...)},
 * {@link AutomationPlatformLinkExt#captureBitmap(String)}
 * <li> <b>Auxiliary APL:</b>
 * {@link AutomationPlatformLink#getComponentMatcher(IDGroup...)}
 */
public class SubToolAudit extends BaseSubTool {

   // external reporting test case screenshot prefix
   // this is an unique name for detecting the screenshot logging
   private static final String EXTERNAL_REPORT_SCREENSHOT_PREFIX = "SCREENSHOT";
   protected static final Logger _logger =
         LoggerFactory.getLogger(SubToolAudit.class);

    public SubToolAudit(UIAutomationTool uiAutomationTool) {
        super(uiAutomationTool);
    }

    // =================================================================
    // Application Screen snapshot tool
    // =================================================================
    /**
     * State Retrieval Method that makes a snapshot of of the screen and
     * saves it in a file at the host OS.
     * @param logTicket - a text marker for the logging of the current snapshot
     * @param t - a failure cause Exception that will be used to generate the
     * file name
     */
    public void snapshotAppScreen(String logTicket, Throwable t) {
        snapshotAppScreen(logTicket, t.getMessage());
    }
    /**
     * State Retrieval Method that makes a snapshot of of the screen and
     * saves it in a file at the host OS.
     * @param logTicket - a text marker for the logging of the current snapshot
     * @param stateDescr - description of the state that has been snapped.
     * @return String with the screenshot filename location
     */
    public String snapshotAppScreen(String logTicket, String stateDescr) {
        String screenShotName = stripNonAlphaNumericChars(logTicket + " " + stateDescr) + ".png";
        // Replace empty spaces with underscores in the file name.
        screenShotName = screenShotName.replace(" ", "_");

        String errorImageLocation = createImageLocation(screenShotName);


        // If the image storage web server log link to it.
        if(!SUITA.Environment.getErrorImageWebURL().isEmpty()) {
           _logger.info("SCREENSHOT URL: " + createImageURL(screenShotName));
        }

        try {
            aplx().captureBitmap(errorImageLocation);
        } catch (Throwable e) {
           _logger.error("FAILED screenshot to: " + errorImageLocation, e);
        }

        _logger.info(EXTERNAL_REPORT_SCREENSHOT_PREFIX + ": {}", errorImageLocation);

        return errorImageLocation;
    }

    /**
     * Helper Method that normalizes the input text to produce a bitmap
     * picture file name.
     * @param imagePrefixFree text explaining the file (e.g. error message)
     * @return  a bitmap picture file name.
     */
    private static String createImageLocation(String imageName) {
        String imageDir = SUITA.Environment.getErrorImageDir();
        if (! imageDir.endsWith("" + File.separatorChar)) {
            imageDir += File.separatorChar;
        }
        return imageDir + imageName;
    }

    /**
     * Return URL to the stored screen shot location.
     * @param imageName  name of the created screen shot.
     * @return URL pointing to the image.
     */
    private static String createImageURL(String imageName) {
       return SUITA.Environment.getErrorImageWebURL() + "/" + imageName;
    }

    // =================================================================
    // Specific auditing routines
    // =================================================================
    public static final String F_P_PREFFIX = "FAIL";
    /**
     * Helper Method that generates Failure Package prefixes to be used for
     * marking all log entries of one failure
     * @return a unique Failure Package Prefix
     */
    public static final String getFPID() {
        return F_P_PREFFIX + String.format("%014d", System.currentTimeMillis());
    }

    /**
     * Static Constant Field holding a Code Zone definition for the
     * SUITA framework
     */
    public static final CodeZone FRAMEWORK_CODE_ZONE =
        (new CodeZone("SUITAF"))
        .setInclusion(Package.getPackage("com.vmware.suitaf"));
    /**
     * Static Constant Field holding a Code Zone definition for the
     * Assertion Sub Tool
     */
    public static final CodeZone ASSERTION_CODE_ZONE =
        (new CodeZone("ASSERT"))
        .setInclusion(SubToolAssertFor.class);
    /**
     * Static Constant Field holding a Code Zone definition for the
     * Automation Platform Link interface implementation
     */
    public static final CodeZone APL_CODE_ZONE =
        (new CodeZone("APL"))
        .setInclusion(Package.getPackage("com.vmware.suitaf.apl"));

    /**
     * Complex State Handling Method that handles APL interface failures
     * @param t - exception that has been thrown by an APL method call
     */
    public void aplFailure(Throwable t) {
        FailureCauseHolder fch = FailureCauseHolder.forCause(t);
        String logTicket = getFPID();

        Analysis anl = CallStackAnalyser.process(
                t.getStackTrace(), FRAMEWORK_CODE_ZONE, APL_CODE_ZONE);
        _logger.error(logTicket + " [" + anl.getExitPoint().getInvokerPoint()
                + "] THROWS: " + t);
        _logger.error(logTicket + " " + anl.getEntryPoint(APL_CODE_ZONE));
        _logger.error(logTicket + " " + anl.getEntryPoint(FRAMEWORK_CODE_ZONE));

        this.snapshotAppScreen(logTicket, t);

        fch.escalateCause();
    }

    /**
     * Complex State Handling Method that handles the outcome of an assertion.
     * Depending on the Throwable parameter:
     * <li> if is <b>null</b> - a success message is logged giving the asserted
     * fact, where in the test the assertion was made and what was the state
     * of the {@link Condition} that was used for evaluation.
     * <li> if not <b>null</b> - an assertion failure is logged giving the
     * asserted fact, where in the test the assertion was made and what was
     * the state of the {@link Condition} that was used for evaluation. Then
     * a screen snapshot is initiated followed by a dump of the screen
     * component tree. At the end an exception is thrown to cause the test
     * to fail.
     * @param t - this Throwable is either null or holds the exception generated
     * by the failed condition evaluation.
     * @param assertedFact - text that describes the action that has been
     * asserted
     * @param assertionCondition - the state of the condition in a text form.
     */
    public void assertionOutcome(
            Throwable t, String assertedFact, String assertionCondition) {
        FailureCauseHolder fch = FailureCauseHolder.forCause(t);

        Analysis anl = CallStackAnalyser.process(
                Thread.currentThread().getStackTrace(),
                ASSERTION_CODE_ZONE);
        if (fch.isNull()) {
           _logger.info("PASSED: {}",
                 assertedFact
                    + " AT: "
                    + anl.getEntryPoint(ASSERTION_CODE_ZONE).getInvokerPoint()
                    + " ASSESSMENT: "
                    + assertionCondition);
        }
        else {
            String logTicket = getFPID();
            _logger.error("FAILED: {}",
                   "[" + logTicket + "] " + assertedFact
                    + " AT: "
                    + anl.getEntryPoint(ASSERTION_CODE_ZONE).getInvokerPoint()
                    + " ASSESSMENT: "
                    + assertionCondition);
            this.snapshotAppScreen(logTicket, t);

            fch.escalateCause();
        }
    }
}
