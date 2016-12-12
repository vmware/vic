package com.vmware.suitaf.apl;

/**
 * A Set - enumeration of the known special UI states.<br>
 * A "special state" will be called such state that:
 * <li> blocks test scenario execution.
 * <li> is not part of any test scenario.
 * <li> can be recognized unambiguously.
 * <li> can be handled safely to allow continuation of current test.
 * <br><br>
 * The set is open and could be extended with time.
 *
 * @author dkozhuharov
 */
public enum SpecialStates {
    // Special States definitions
    WIN_IE_CERT_ERROR_OVERRIDE_LINK(
            "Click Certificate error override link on the page."),
    WIN_IE_WEBCERTDLG_CLOSE(
            "Close Web Certificate Dialog based on <JDialog>."),
    WIN_IE_WEBCERTNOTTHRUSTEDDLG_CLOSE(
            "Close Certificate Not Trusted Dialog based on <Window>."),
    WIN_IE_SECURITYWARNINGDIALOG_CLOSE(
            "Close Security Warning Dialog based on <Window>"),
    WIN_IE_CONSOLEWINDOW_CLOSE(
            "Close VM Console Window based on <Window>."),
    WIN_IE_ANYDIALOG_CLOSE(
            "Close Any Dialogs based on <Dialog>."),
    WIN_IE_FILECHOOSER_CLOSE(
            "Close a file-chooser dialog based on <Window>"),
    WIN_IE_PROGRESS_CLOSE(
            "Close A Progress window based on <BrowserWindow>."),
    WIN_IE_CLOSE_CLICK_IGNORE(
            "Close A System dialog based on <PushButton>."),
    ;

    public final String description;
    public final SpecialStateHandler blankHandler;
    SpecialStates(String description) {
        this.description = description;
        this.blankHandler =
            new SpecialStateHandler() {

            @Override
            public void passOptionalData(Object... optionals) {
            }

            @Override
            public void stateHandle() {
                throw new UnsupportedOperationException(
                        "State " + this + " handling is not implemented");
            }

            @Override
            public boolean stateRecognize() {
                return false;
            }
        };
    }
}