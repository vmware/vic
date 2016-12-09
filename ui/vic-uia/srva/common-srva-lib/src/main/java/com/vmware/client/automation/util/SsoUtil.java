/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.util;

import java.net.InetAddress;
import java.net.UnknownHostException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.servicespec.HostServiceSpec;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.vsphere.client.automation.provider.connector.HostConnector;
import com.vmware.vsphere.client.automation.provider.connector.VcConnector;

/** A class for managing SSO connections to the VC */
public class SsoUtil {

   private static final Logger _logger = LoggerFactory.getLogger(SsoUtil.class);

   // Default login details used for backward compativility to the composition
   // workflow model
   // TODO: rkovachev remove it once the test are migrated to new model
   private static VcServiceSpec _defaultServiceSpec = null;

   private static VcConnector staticConnector = null;

   /**
    * Set default login details. The service spec should be provided for each
    * requested spec in the testbed.
    * NOTE: Keep it just for the backward compatibility to the System test
    * team.
    * TODO: remove it once the the tests are migrated to the new workflow model
    * 
    * @param serverIp default server IP
    * @param username default username
    * @param password default password
    **/
   @Deprecated
   public static void setLoginCredentials(String serverIp, String username,
         String password) {
      _defaultServiceSpec = new VcServiceSpec();
      _defaultServiceSpec.endpoint.set(serverIp);
      _defaultServiceSpec.username.set(username);
      _defaultServiceSpec.password.set(password);
   }

   /**
    * Check if the <code>spec</code> is null. If so return for usage the default
    * service spec using the connection details set by <code>setLoginCredentials</code>. If the default service spec is
    * not
    * initialized the method will throw RuntimeException.
    *
    * @param spec provides the LDU connection details.
    * @return ServiceSpec defines the LDU connection details for obtaining
    *         service.
    * */
   @Deprecated
   // Use the spec from TestSpec
   public static ServiceSpec getServiceSpec(ServiceSpec spec) {
      if (spec != null) {
         return spec;
      }
      if (_defaultServiceSpec == null) {
         // When running with the new testworkflow default service spec is not set
         _logger.debug("Default service spec is not initialized!");
      }
      return _defaultServiceSpec;
   }

   /**
    * Provide connector for the specified service spec.
    * 
    * @param serviceSpec define the service connection details
    * @return connector object for the provided service spec.
    *
    *         NOTE: For backward compatibility the method will work even the provided
    *         spec is null. That is temporary and will be removed once the tests are
    *         migrated to the new test workflow model.
    *
    *         TODO: rkovachev remove the code that handle the case for null serviceSpec.
    */
   public static TestbedConnector getConnector(ServiceSpec serviceSpec) {
      // NOTE: Temporary workaround till the existign tests are migrated to the
      // new model.
      // TODO: rkovachev remove the code once the tests are migrated.
      if (serviceSpec == null) {
         // load the default service spec if no spec is provided.
         serviceSpec = getServiceSpec(null);
      }

      // Registry get connector for ServiceSpec serviceSpec.
      // Get connector for VC
      if (serviceSpec instanceof VcServiceSpec) {
         // Get SSO from spec using the Spec and Registry
         VcServiceSpec vcSericeSpec = new VcServiceSpec();
         vcSericeSpec.copy(serviceSpec);
         // when no default servie spec is set the code is running in the new
         // provider/test workflow controller model.
         if (_defaultServiceSpec == null) {
            // Load the connection from the registry
            VcConnector connector =
                  (VcConnector) WorkflowRegistry.getRegistry().getActiveTestbedConnection(
                        serviceSpec);

            // The connector is not registered in the registry. Try to init now.
            // It is the case in the provider chech health and disassemble flows.
            if (connector == null) {
               try {
                  connector = new VcConnector(vcSericeSpec);
               } catch (SsoException e) {
                  throw new RuntimeException("Failed to initialize VC Connector!", e);
               }
            }

            if (!connector.isAlive()) {
               connector.connect();
            }
            return connector;
         } else {
            // Load the connection based on the default service spec.
            if (staticConnector == null) {
               try {
                  VcServiceSpec vcServiceSpec = new VcServiceSpec();
                  vcServiceSpec.copy(serviceSpec);

                  staticConnector = new VcConnector(vcServiceSpec);
               } catch (SsoException ssoEexception) {
                  throw new RuntimeException(
                        "Error during creating SSO conneciton", ssoEexception);
               }
            }
            staticConnector.connect();
            return staticConnector;
         }
      } else if (serviceSpec instanceof HostServiceSpec) {
         // VcConnector connector =
         // (VcConnector) WorkflowRegistry.getRegistry().getActiveTestbedConnection(serviceSpec);

         // Init connector based on the provide spec. This is temporary solution
         // for testing Hostclient connection.
         // TODO: rkovachev, lgrigorova
         // Remove the static init once the test is fixed to in the new model
         // if (connector == null) {
         // Get connector for Host Client
         try {
            HostServiceSpec hostServiceSpec = new HostServiceSpec();
            hostServiceSpec.copy(serviceSpec);

            return new HostConnector(hostServiceSpec);
         } catch (SsoException e) {
            throw new RuntimeException("Fail to initialize Host Connector");
         }
         // }

         // return connector;
      } else {
         throw new RuntimeException("Ivalid serviceSpec type." + serviceSpec.toString());
      }
   }

   /**
    * Return VcConnector for the specified entity spec.
    * 
    * @param entitySpec VC entity object
    * @return
    * @throws Exception
    */
   public static VcConnector getVcConnector(EntitySpec entitySpec) throws Exception {
      ServiceSpec serviceSpec = entitySpec.service.get();
      if (serviceSpec instanceof VcServiceSpec) {
         return (VcConnector) getConnector(serviceSpec);
      } else {
         throw new Exception("Service spec is not a VcServiceSpec!");
      }
   }

   /**
    * Resolves a host from a given IP
    *
    * @param ip - the IP to resolve
    * @return The host name
    */
   public static String resolveHost(String ip) {
      InetAddress ipAddress;
      try {
         ipAddress = InetAddress.getByName(ip);
      } catch (UnknownHostException e) {
         throw new RuntimeException(e);
      }
      return ipAddress.getHostName();
   }

   // Private methods below

}
