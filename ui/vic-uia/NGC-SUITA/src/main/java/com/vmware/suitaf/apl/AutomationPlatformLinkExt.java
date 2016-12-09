package com.vmware.suitaf.apl;

import java.util.Map;

/**
 * This interface extends the base {@link AutomationPlatformLink} interface.
 * The extensibility functions added here are in support of the test development
 * and test exploitation.
 *
 * @author dkozhuharov
 */
public interface AutomationPlatformLinkExt extends AutomationPlatformLink {

    // ==============================================================
    //    (APL-10) Link setup and recovery;
    // ==============================================================

    /**
     * This method is similar to the method:
     * {@link AutomationPlatformLink#resetLink(Map)}. The difference is that
     * this method after connecting to the Automation Agent tries to gain
     * control of an already opened browser.
     * <br><br>
     * This method will be used for screen inspection tasks.
     *
     * @param linkInitParams - a map of key-value pairs that represent a
     * complete and meaningful set of parameter for the initialization of
     * the current APL implementation.
     * @param hostLogger - a call-back interface to allow controlled logging
     * in the main SUITA log stream.
     */
    void attachLink(Map<String, String> linkInitParams, HostLogger hostLogger);



    // ==============================================================
    //    (APL-11) Browser control;
    // ==============================================================

    // Browser programmability functions
    /**
     * This method should provide java-script code execution in the browser.
     * @param code - java-script source code
     */
    void executeJavaScript(String code);

    // ==============================================================
    //    (APL-13) Screen inspection (parent/child component relations
    //             and property enumeration)
    // ==============================================================

    /**
     * This method will make a snapshot of the desktop and store it in a
     * bitmap format in a file.
     *
     * @param fileName - the full file name where to store the snapshot
     */
    void captureBitmap(String fileName);





}
