/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.exception;

/**
 * Base exception class for exceptions that have occurred in <link>NGCTocTree</link>.
 */
public class TocTreeException extends RuntimeException {

   private static final long serialVersionUID = 1;

   public TocTreeException(String message) {
      super(message);
   }

   public TocTreeException(String message, Throwable throwable) {
      super(message, throwable);
   }
}

