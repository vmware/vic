package com.vmware.suitaf.util;

/**
 * This class is used in places where a caught exception must be captured,
 * analyzed, some remedies implemented, logging done and eventually
 * re-thrown
 *
 * @author dkozhuharov
 */
public class FailureCauseHolder {
    private Throwable cause = null;
    private boolean isRuntimeException = false;
    private boolean isError = false;
    private boolean isNull = true;

    public static final FailureCauseHolder forCause(Throwable cause) {
        FailureCauseHolder fch = new FailureCauseHolder();
        fch.setCause(cause);
        return fch;
    }

    /**
     * Method to set the cause of the failure. Must be used in the catch
     * clause where the error is cought
     *
     * @param cause
     */
    public void setCause(Throwable cause) {
        this.cause = cause;
        this.isRuntimeException = (cause instanceof RuntimeException);
        this.isError = (cause instanceof Error);
        this.isNull = (cause == null);
    }
    /** Retrieve the captured failure cause
     * @return
     */
    public Throwable getCause() { return this.cause; }
    /** Checks if the failure cause is instance of {@link RuntimeException}
     * @return
     */
    public boolean isRuntimeException() { return isRuntimeException; }
    /** Checks if the failure cause is instance of {@link Error}
     * @return
     */
    public boolean isError() { return isError; }
    /** Checks if the failure cause is missing
     * @return
     */
    public boolean isNull() { return isNull; }
    /** Re-throws the failure cause if present. */
    public void escalateCause() {
        if (isRuntimeException)  throw (RuntimeException) cause;
        if (isError)             throw (Error) cause;
    }
}