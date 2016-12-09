/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.connector;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.servicespec.HostServiceSpec;
import com.vmware.client.automation.sso.SsoClient;

/**
 * Class that provides host client connector.
 */
public class HostConnector implements TestbedConnector {

   private final ServiceSpec _connectorSpec;
   private final SsoClient _ssoClient;

   /**
    * Constructor for HostConnector that initializes SsoCLient and the connection spec
    * @param connectorSpec - should be HostServcieSpec
    * @throws SsoException - it is actually thrown only in the case of VcConnection
    */
   public HostConnector(HostServiceSpec connectorSpec) throws SsoException {
      _connectorSpec = connectorSpec;
      _ssoClient = new SsoClient(_connectorSpec);
   }

   @Override
   public void connect() {
      try {
         _ssoClient.getVcService();
      } catch (VcException e) {
         throw new RuntimeException("Host connection failed!", e);
      }
   }

   @Override
   public boolean isAlive() {
      // TODO Auto-generated method stub
      return true;
   }

   @Override
   public void disconnect() {
      // TODO Auto-generated method stub

   }

   @SuppressWarnings("unchecked")
   @Override
   public SsoClient getConnection() {
      return _ssoClient;
   }

}
