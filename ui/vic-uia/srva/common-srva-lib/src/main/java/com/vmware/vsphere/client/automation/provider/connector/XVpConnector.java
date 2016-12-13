package com.vmware.vsphere.client.automation.provider.connector;

import java.io.IOException;
import java.net.Socket;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.servicespec.XVpServiceSpec;

/**
 * Class that is used to connect to a xVp
 */
public class XVpConnector implements TestbedConnector {

   private static final Logger _logger = LoggerFactory
         .getLogger(XVpConnector.class);

   private final ServiceSpec _connectorSpec;
   private static final int XVP_SERVER_TCP_PORT = 8443;

   public XVpConnector(XVpServiceSpec connectorSpec) {
      _connectorSpec = connectorSpec;
   }

   @Override
   public void connect() {
      // do nothing
   }

   @Override
   public boolean isAlive() {
      boolean isAlive = true;

      Socket socket = null;
      try {
         socket = new Socket(_connectorSpec.endpoint.get(), XVP_SERVER_TCP_PORT);
      } catch (IOException e) {
         _logger.debug("Unable to connect to port " + XVP_SERVER_TCP_PORT
               + "; VP might be dead");
         isAlive = false;
      } finally {
         if (socket != null) {
            try {
               socket.close();
            } catch (IOException e) {
               _logger.debug("Unable to close VP Server tcp connection");
               e.printStackTrace();
            }
         }
      }

      return isAlive;
   }

   @Override
   public void disconnect() {
      // do nothing

   }

   @SuppressWarnings("unchecked")
   @Override
   public Object getConnection() {
      // do nothing
      return null;
   }

}
