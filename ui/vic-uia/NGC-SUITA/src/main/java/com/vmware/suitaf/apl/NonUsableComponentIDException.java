package com.vmware.suitaf.apl;

/**
 * This exception is thrown by the factory method {@link
 * AutomationPlatformLink#getComponentID(ComponentIDTypes, String...)} when
 * it could not convert the generic ID format to internal representation.
 * <br><br>
 * @author dkozhuharov
 */
@SuppressWarnings("serial")
public class NonUsableComponentIDException extends RuntimeException {
    public NonUsableComponentIDException(String message, Throwable cause) {
        super(message, cause);
    }
    public NonUsableComponentIDException(String message) {
        super(message);
    }
    public NonUsableComponentIDException(Throwable cause) {
        super(cause);
    }
}