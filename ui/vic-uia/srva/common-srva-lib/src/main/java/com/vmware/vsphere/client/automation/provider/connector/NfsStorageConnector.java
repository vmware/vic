/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.connector;

import java.io.IOException;
import java.net.Socket;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.NfsStorageServiceSpec;

/**
 * Class that is used to connect to a nfs storage
 */
public class NfsStorageConnector implements TestbedConnector {

   private static final Logger _logger =
         LoggerFactory.getLogger(NfsStorageConnector.class);

   private final ServiceSpec _connectorSpec;
   private static final int NFS_SERVER_TELNET_PORT = 2049;

   public NfsStorageConnector(NfsStorageServiceSpec connectorSpec) {
      _connectorSpec = connectorSpec;
   }

   @Override
   public void connect() {
      // Nothing to do here
   }

   @Override
   public boolean isAlive() {
      boolean isAlive = true;

      Socket socket = null;
      try {
         socket = new Socket(_connectorSpec.endpoint.get(), NFS_SERVER_TELNET_PORT);
      } catch (IOException e) {
         _logger.debug("Unable to connect to port " + NFS_SERVER_TELNET_PORT + "; NFS server might be dead");
         isAlive = false;
      } finally {
         if (socket != null) {
            try {
               socket.close();
            } catch (IOException e) {
               _logger.debug("Unable to close NFS Server telnet connection");
               e.printStackTrace();
            }
         }
      }

      return isAlive;
   }

   @Override
   public void disconnect() {
      // Nothing to do here

   }

   @SuppressWarnings("unchecked")
   @Override
   public Object getConnection() {
      // Nothing to do here
      return null;
   }
}
