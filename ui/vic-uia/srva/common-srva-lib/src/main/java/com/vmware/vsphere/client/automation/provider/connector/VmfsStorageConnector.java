/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.connector;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.VmfsStorageServiceSpec;

/**
 * Class that is used to connect to a vmfs storage
 */
public class VmfsStorageConnector implements TestbedConnector {

   private static final Logger _logger = LoggerFactory.getLogger(VmfsStorageConnector.class);

   private final ServiceSpec _connectorSpec;
   private static final int VMFS_SERVER_TCP_PORT = 3260;
   private static final int TCP_CONNECTION_TIMEOUT = 10000;

   public VmfsStorageConnector(VmfsStorageServiceSpec connectorSpec) {
      _connectorSpec = connectorSpec;
   }

   @Override
   public void connect() {
      // Nothing to do here
   }

   @Override
   public boolean isAlive() {
      boolean isAlive = true;

      // TODO: Find a common location to this code. Probably most of the connectors will use it
      Socket socket = new Socket();
      try {
         socket.connect(
            new InetSocketAddress(_connectorSpec.endpoint.get(), VMFS_SERVER_TCP_PORT),
            TCP_CONNECTION_TIMEOUT);
      } catch (IOException e) {
         _logger.debug("Unable to connect to TCP port " + VMFS_SERVER_TCP_PORT + "; VMFS server might be dead", e);
         isAlive = false;
      } finally {
         try {
            socket.close();
         } catch (IOException e) {
            _logger.debug("Unable to close VMFS Server TCP connection", e);
         }
      }

      return isAlive;
   }

   @Override
   public void disconnect() {
      // Nothing to do here

   }

   @Override
   public <T extends Object> T getConnection() {
      // Nothing to do here
      return null;
   }
}
