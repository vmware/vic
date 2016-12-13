/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele.sikulix;

/**
 * Sikuli connection types.
 */
enum ConnectionType {
   /**
    * Used when the target screen is on the local machine
    */
   LOCAL,
   /**
    * Used when the target screen is on a remote VNC machine
    */
   REMOTE_VNC
}
