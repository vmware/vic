/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.connector;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.sso.SsoClient;
import com.vmware.vim.binding.vim.fault.NotAuthenticated;

/**
 * VC client connector.
 *
 */
public class VcConnector implements TestbedConnector {

   private final ServiceSpec _connectorSpec;
   private SsoClient _ssoClient;

   public VcConnector(VcServiceSpec connectorSpec) throws SsoException {
      _connectorSpec = connectorSpec;
      _ssoClient = new SsoClient(_connectorSpec);
   }

   @Override
   public boolean isAlive() {
      return _ssoClient.isAlive() && isVcSessionAuthenticated();
   }

   private boolean isVcSessionAuthenticated() {
      boolean result = true;
      try {
         _ssoClient.getVcService().getServiceInstance().currentTime();
      } catch (NotAuthenticated | VcException e) {
         result = false;
      }
      return result;
   }

   @SuppressWarnings("unchecked")
   @Override
   public SsoClient getConnection() {
      return _ssoClient;
   }

   @Override
   public void connect() {
      try {
         _ssoClient.connect();
      } catch (SsoException e) {
         throw new RuntimeException(e);
      }
   }

   @Override
   public void disconnect() {
      _ssoClient.disconnect();

   }

   /**
    * Reconnects the connector
    *
    * @throws SsoException
    */
   public void reconnect() throws SsoException {
      _ssoClient = new SsoClient(_connectorSpec);
      connect();
   }

}
