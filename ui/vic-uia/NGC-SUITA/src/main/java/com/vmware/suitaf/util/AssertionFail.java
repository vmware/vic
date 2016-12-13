package com.vmware.suitaf.util;

/**
 * This error is used to mark that specific framework assertion has been made
 * and a procedure for state logging has been performed. It allows to avoid
 * double test failure handling.
 *
 * @author dkozhuharov
 */
@SuppressWarnings("serial")
public class AssertionFail extends RuntimeException {
    public AssertionFail() {
        super();
    }
    public AssertionFail(String message, Throwable cause) {
        super(message, cause);
    }
    public AssertionFail(String message) {
        super(message);
    }
    public AssertionFail(Throwable cause) {
        super(cause);
    }
}