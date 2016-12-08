/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.exception;

/**
 * Wraps exceptions occurred while working with the SSO.
 */
public class SsoException extends Exception {

   public SsoException(Exception cause) {
      super(cause);
   }
}