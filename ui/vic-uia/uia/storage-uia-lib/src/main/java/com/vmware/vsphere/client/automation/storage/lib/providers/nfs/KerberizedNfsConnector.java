package com.vmware.vsphere.client.automation.storage.lib.providers.nfs;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;

import org.apache.commons.lang.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;

/**
 * TestbedConnector implementation for a Kerberized NFS
 *
 * @see TestbedConnector
 */
public class KerberizedNfsConnector implements TestbedConnector {

   private static final int KERBERIZED_NFS_SERVER_TCP_PORT = 2049;
   private static final int TCP_CONNECTION_TIMEOUT = 10000;
   private static final Logger _logger = LoggerFactory
         .getLogger(KerberizedNfsConnector.class);

   private final ServiceSpec serviceConnectionSpec;

   /**
    * Initializes new instance of KerberizedNfsConnector
    *
    * @param serviceConnectionSpec
    *           the kerberized nfs service connection spec
    */
   public KerberizedNfsConnector(ServiceSpec serviceConnectionSpec) {
      this.serviceConnectionSpec = serviceConnectionSpec;
   }

   @Override
   public void connect() {
      // No need to connect to the NFS server from the automation codebase.
   }

   @Override
   public boolean isAlive() {
      Socket socket = new Socket();
      try {
         socket.connect(
               new InetSocketAddress(serviceConnectionSpec.endpoint.get(),
                     KERBERIZED_NFS_SERVER_TCP_PORT), TCP_CONNECTION_TIMEOUT);
         return true;
      } catch (IOException e) {
         _logger.debug("Unable to connect to TCP port "
               + KERBERIZED_NFS_SERVER_TCP_PORT
               + "; Kerberized NFS server might be dead", e);
      } finally {
         try {
            socket.close();
         } catch (IOException e) {
            _logger.debug(
                  "Unable to close Kerberized NFS Server TCP connection", e);
         }
      }
      return false;

   }

   @Override
   public void disconnect() {
      // No need to disconnect to the NFS server from the automation codebase.
   }

   @Override
   public <T> T getConnection() {
      throw new NotImplementedException();
   }

}
