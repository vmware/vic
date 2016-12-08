package com.vmware.vsphere.client.automation.srv.exception;

/**
 * Server object was not found
 */
public class ObjectNotFoundException extends Exception {
   private static final long serialVersionUID = 3950005462766841155L;

   /**
    * Initializes new ObjectNotFoundException
    *
    * @param message
    *           describe the object which is not found.
    */
   public ObjectNotFoundException(String message) {
      super(message);
   }
}
