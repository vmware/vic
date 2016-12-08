/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.connector;

/**
 * Interface to define object that keeps connection and provides
 * method for detecting and managing its state.
 * @param <T>
 */
public interface TestbedConnector {

   /**
    * Connect the <T> client.
    */
   public abstract void connect();

   /**
    * Check the connection status.
    * @return true if the connection is active.
    */
   public abstract boolean isAlive();

   /**
    * Disconnect client.
    */
   public abstract void disconnect();

   /**
    * Get connection client T
    * @return T
    */
   public abstract <T extends Object> T getConnection();

}
