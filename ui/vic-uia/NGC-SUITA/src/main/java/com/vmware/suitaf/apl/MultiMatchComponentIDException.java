package com.vmware.suitaf.apl;

/**
 * It's been taken into account that a component identifier
 * could in some cases match more than one component on screen. This is
 * accepted and handled in some methods. Others, which don't have meaningful
 * scenario for multiple matches will throw this exception.
 *
 * @author dkozhuharov
 */
@SuppressWarnings("serial")
public class MultiMatchComponentIDException extends RuntimeException {
    public MultiMatchComponentIDException(String message, Throwable cause) {
        super(message, cause);
    }
    public MultiMatchComponentIDException(String message) {
        super(message);
    }
    public MultiMatchComponentIDException(Throwable cause) {
        super(cause);
    }
}