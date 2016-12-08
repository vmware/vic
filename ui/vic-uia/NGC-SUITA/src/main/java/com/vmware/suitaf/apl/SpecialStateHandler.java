package com.vmware.suitaf.apl;

public interface SpecialStateHandler {

    void passOptionalData(Object...optionals);
    /**
     * The implementation of this method must "know" when, where and how this
     * "special state" appears. It should return <b>true</b> only if the state
     * could be securely identified and the implementor have a handling routine
     * implemented.
     * <br><br>
     * E.g. if the "special state" is a security warning dialog, the implementor
     * must be able to see and recognize this dialog and be able to click the
     * appropriate button to allow test scenario continuation. If the dialog
     * is browser-specific (as will be in most cases) - the implementor must
     * check the currently used browser. It must return <b>false</b> in case
     * the "special state" is irrelevant for the current browser.
     *
     * @param sState - enum value determining the "special state" being
     * checked (see {@link SpecialStates})
     * @param optional - {@link IDPart} property filters that could be
     * needed for the state recognition. Depends on the specific state.
     * @return - <b>true</b> if the state was found present and could be
     * handled.
     */
    boolean stateRecognize();

    /**
     * The implementation of this method must be able to safely handle the
     * "special state" to allows continuation of the currently executing test.
     * If handling attempt fails - runtime exception should be thrown.
     *
     * @param sState - enum value determining the "special state" being
     * checked (see {@link SpecialStates})
     * @param optional - {@link IDPart} property filters that could be
     * needed for the state recognition. Depends on the specific state.
     */
    void stateHandle();

}
