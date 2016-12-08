/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

/**
 * The exception is thrown if a spec cannot be found in the implementation in
 * a spec container.
 */
public class SpecNotFoundException extends Exception {

	/**
	 * 
	 */
   private static final long serialVersionUID = 4120928585406950965L;

   public SpecNotFoundException(String message) {
      super(message);
   }
}
