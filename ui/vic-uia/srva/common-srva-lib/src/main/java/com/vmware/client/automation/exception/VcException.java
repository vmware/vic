package com.vmware.client.automation.exception;

public class VcException extends Exception {

   private static final long serialVersionUID = -4761393545874520054L;

    public VcException(String message) {
      super("VC Service failure: " + message);
   }

   public VcException(String message, Throwable cause) {
      super("VC Service failure: " + message, cause);
   }

}
