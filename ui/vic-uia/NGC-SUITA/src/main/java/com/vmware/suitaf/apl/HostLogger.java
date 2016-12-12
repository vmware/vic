package com.vmware.suitaf.apl;

/**
 * This is a generic logger interface. It will be used as a call-back interface
 * for the {@link AutomationPlatformLink} implementors. It will allow the
 * implementor to join its log output in the SUITA framework logging stream.
 * Additionally the framework code receives the control of the log level.
 * <br><br>
 * The methods of the interface are deliberately left without an
 * {@link Exception} parameter. In any case of execution failures an Exception
 * must be thrown.
 *
 * @author dkozhuharov
 */
public interface HostLogger {

    /**
     * Send debug-type messages for logging
     * @param message - the message string.
     */
    public void debug(String message);

    /**
     * Send info-type messages for logging
     * @param message - the message string.
     */
    public void info(String message);

    /**
     * Send warn-type messages for logging
     * @param message - the message string.
     */
    public void warn(String message);

    /**
     * Send error-type messages for logging
     * @param message - the message string.
     */
    public void error(String message);

    /**
     * Send messages for plain logging
     * @param message - the message string.
     */
    public void dump(String message);
}
