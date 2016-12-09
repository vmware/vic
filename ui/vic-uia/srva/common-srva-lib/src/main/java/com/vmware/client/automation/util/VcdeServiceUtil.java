/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util;

import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;

import com.vmware.cis.authz.Permission;
import com.vmware.cis.authz.PermissionStub;
import com.vmware.cis.authz.PrivilegeStub;
import com.vmware.cis.authz.RoleStub;
import com.vmware.cis.authz.sessions.SessionManager;
import com.vmware.cis.cm.client.CisServiceType;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.sso.SsoClient;
import com.vmware.content.Library;
import com.vmware.content.LocalLibrary;
import com.vmware.content.SubscribedLibrary;
import com.vmware.content.library.Item;
import com.vmware.content.library.item.UpdateSession;
import com.vmware.content.library.item.updatesession.File;
import com.vmware.vapi.bindings.Service;
import com.vmware.vapi.bindings.StubConfiguration;
import com.vmware.vapi.bindings.StubFactory;
import com.vmware.vapi.cis.authn.ProtocolFactory;
import com.vmware.vapi.cis.authn.SecurityContextFactory;
import com.vmware.vapi.core.ApiProvider;
import com.vmware.vapi.core.ExecutionContext.SecurityContext;
import com.vmware.vapi.protocol.ProtocolConnection;
import com.vmware.vcenter.ovf.ExportSession;
import com.vmware.vise.vim.security.accesscontrol.AuthorizationService;
import com.vmware.vise.vim.security.accesscontrol.impl.AuthorizationServiceImpl;

/**
 * Util class to provide ability to connect to vCDe services.
 */
public class VcdeServiceUtil {

   // vCDe services URL templates
   // TODO rkovachev: Get the VAPI endpoint from the CM
   private static final String SERVICE_URL_TEMPLATE = "https://%s:443/site/api";

   private static final String INVSVC_VAPI_TYPE_ID = "com.vmware.cis.inventory.vapi";

   /**
    * Return <code>SubscribedLibrary</code> service instance. The
    * <code>ServiceSpec</code> parameter provides the LDU connection details to
    * use. If SSO connection to it is already established uses it. Otherwise
    * establish connection to the specified LDU.
    *
    * @param spec
    *           specify the service LDU connection details. If the
    *           <code>spec</code> is null uses the default login details set by
    *           the <code>setLoginCredentials</code>.
    * @return <code>LocalLibraryService</code> service instance.
    * @throws Exception
    *            if the SSO connection can not be established.
    */
   public static SubscribedLibrary getSubscribedLibrary(ServiceSpec serviceSpec)
         throws Exception {

      return getService(SubscribedLibrary.class, serviceSpec);
   }

   /**
    * Return <code>LocalLibrary</code> service instance. The
    * <code>ServiceSpec</code> parameter provides the LDU connection details to
    * use. If SSO connection to it is already established uses it. Otherwise
    * establish connection to the specified LDU.
    *
    * @param spec
    *           specify the service LDU connection details. If the
    *           <code>spec</code> is null uses the default login details set by
    *           the <code>setLoginCredentials</code>.
    * @return <code>LocalLibrary</code> service instance.
    * @throws Exception
    *            if the SSO connection can not be established.
    */
   public static LocalLibrary getLocalLibraryService(ServiceSpec serviceSpec)
         throws Exception {
      return getService(LocalLibrary.class, serviceSpec);
   }

   /**
    * Return <code>LibraryService</code> service instance. The
    * <code>ServiceSpec</code> parameter provides the LDU connection details to
    * use. If SSO connection to it is already established uses it. Otherwise
    * establish connection to the specified LDU.
    *
    * @param spec
    *           specify the service LDU connection details. If the
    *           <code>spec</code> is null uses the default login details set by
    *           the <code>setLoginCredentials</code>.
    * @return <code>Library</code> service instance.
    * @throws Exception
    *            if the SSO connection can not be established.
    */
   public static Library getLibraryService(ServiceSpec serviceSpec)
         throws Exception {
      return getService(Library.class, serviceSpec);
   }


   /**
    * Return <code>LibraryItemService</code> service instance. The
    * <code>ServiceSpec</code> parameter provides the LDU connection details to
    * use. If SSO connection to it is already established uses it. Otherwise
    * establish connection to the specified LDU.
    *
    * @param spec
    *           specify the service LDU connection details. If the
    *           <code>spec</code> is null uses the default login details set by
    *           the <code>setLoginCredentials</code>.
    * @return <code>LibraryItemService</code> service instance.
    * @throws Exception
    *            if the SSO connection can not be established.
    */
   public static Item getLibraryItemService(ServiceSpec serviceSpec)
         throws Exception {
      return getService(Item.class, serviceSpec);
   }

   /**
    * Gets update session instance.
    * The service code parameter provides the LDU connection details to use.
    * If SSO connection to it is already established uses it.
    * Otherwise establish connection to the specified LDU.
    *
    * @param spec       specify the service LDU connection details. If the spec is null
    *                   uses the default login details set by the setLoginCredentials
    * @return           created update session instance
    * @throws Exception if the SSO connection can not be established
    */
   public static UpdateSession getUpdateSessionService(ServiceSpec serviceSpec)
         throws Exception {
      return getService(UpdateSession.class, serviceSpec);
   }

   /**
    * Gets update file session instance.
    * The service code parameter provides the LDU connection details to use.
    * If SSO connection to it is already established uses it.
    * Otherwise establish connection to the specified LDU.
    *
    * @param spec       specify the service LDU connection details. If the spec is null
    *                   uses the default login details set by the setLoginCredentials
    * @return           created update file session instance
    * @throws Exception if the SSO connection can not be established
    */
   public static File getUpdateSessionFileService(ServiceSpec serviceSpec)
         throws Exception {
      return getService(File.class, serviceSpec);
   }

   /**
    * Return <code>ExportSession</code> service instance.
    * The <code>ServiceSpec</code> parameter provides the LDU connection details to use.
    * If SSO connection to it is already established uses it.
    * Otherwise establish connection to the specified LDU.
    *
    * @param spec       specify the service LDU connection details.
    *                   If the <code>spec</code> is null uses the default login details
    *                   set by the <code>setLoginCredentials</code>
    * @return           OVF exporter service instance
    * @throws Exception if the SSO connection can not be established
    */
   public static ExportSession getOvfExportService(ServiceSpec serviceSpec) throws Exception {
      return getService(ExportSession.class, serviceSpec);
   }

   /**
    * Method that gets the Permission Manager
    *
    * @param serviceSpec - specification for connection to vcenter
    * @return permission manager
    * @throws Exception - if there is error in connection
    */
   public static Permission getPermissionManager(
         ServiceSpec serviceSpec) throws Exception {

      //TODO lgrigorova: check if we need to add a logout method for the session
      char[] sessionId = getVapiSessionManager(serviceSpec).create();

      // authenticating to vapi through session id - cannot be done through saml
      StubConfiguration stubConfig = getStubConfig(sessionId);

      // getting of Api Provider from connection to vapi invsvc
      ApiProvider apiProvider = getConnection(getVapiInvsvcUrlFromCm(serviceSpec)).getApiProvider();

      return new PermissionStub(apiProvider, stubConfig);
   }

   /**
    * Creates a stub for the specified service interface.
    *
    * @param ssoClient
    *           SSO client connected to the needed LDU
    * @param serviceClass
    *           service that need to be return
    * @param endPointUrl
    *           end point of the service
    * @return VAPI service stub.
    * @throws KeyStoreException
    * @throws NoSuchAlgorithmException
    * @throws KeyManagementException
    */
   private static <T extends Service> T getService(
         Class<T> serviceClass, String endPointUrl, SsoClient ssoConnector)
               throws KeyManagementException, NoSuchAlgorithmException,
               KeyStoreException {

      SecurityContext context = SecurityContextFactory
            .createSamlSecurityContext(
                  ssoConnector.getLastToken(), ssoConnector.getPrivateKey());

      return getVapiService(serviceClass, context, endPointUrl);
   }

   /**
    * Creates a stub for the specified service interface depending on Security Context.
    *
    * @param ssoClient
    *           SSO client connected to the needed LDU
    * @param serviceClass
    *           service that need to be return
    * @param endPointUrl
    *           end point of the service
    * @return VAPI service stub.
    * @throws KeyStoreException
    * @throws NoSuchAlgorithmException
    * @throws KeyManagementException
    */
   private static <T extends Service> T getVapiService(Class<T> serviceClass, SecurityContext context,
         String endPointUrl) throws KeyManagementException, NoSuchAlgorithmException, KeyStoreException {

      StubFactory stubFactory = new StubFactory(
            getConnection(endPointUrl).getApiProvider());
      StubConfiguration stubConfig = new StubConfiguration();
      stubConfig.setSecurityContext(context);

      return stubFactory.createStub(serviceClass, stubConfig);
   }


   /**
    * Method that gets the SessionManager for vapi invsvc through saml
    * @param serviceSpec - service spec for connection to vcenter
    * @return session manager for vapi
    * @throws Exception if there is error in connection
    */
   private static SessionManager getVapiSessionManager(ServiceSpec serviceSpec) throws Exception {
      String invServiceEndpointUrl = getVapiInvsvcUrlFromCm(serviceSpec);
      SsoClient ssoConnector = (SsoClient) SsoUtil.getConnector(serviceSpec).getConnection();
      return getService(SessionManager.class, invServiceEndpointUrl, ssoConnector);
   }

   /**
    * Method creates a Stub Configuration with session security context. It can
    * be used only by session-aware services.
    * @param sessionId - session id from connection established by the vapi session manager
    * @return stub configuration by session id
    */
   private static StubConfiguration getStubConfig(char[] sessionId) {
      SecurityContext context = SecurityContextFactory.
            createSessionSecurityContext(sessionId);

      StubConfiguration stubConfig = new StubConfiguration();
      stubConfig.setSecurityContext(context);

      return stubConfig;
   }

   /**
    * Method that gets the inventory service url from component manager
    * @param serviceSpec - service specification for login to vcenter
    * @return invsvc url
    * @throws Exception - if there is error in login
    */
   private static String getVapiInvsvcUrlFromCm(ServiceSpec serviceSpec) throws Exception {
      SsoClient ssoConnector = (SsoClient) SsoUtil.getConnector(serviceSpec).getConnection();
      return ssoConnector.getEndPointUrlByServiceTypeFromCm(
            CisServiceType.IS.getServiceType(), INVSVC_VAPI_TYPE_ID);
   }

   // TODO: rkovachev Re write it to use VapiConnectionManagerImpl from the platform
   private static <T extends Service> T getService(Class<T> serviceClass, ServiceSpec serviceSpec) throws Exception {

      // VapiConnectionManagerImpl connManager = new VapiConnectionManagerImpl(
      // this.getComponentManage().lookupComponentManager(false),
      // this._ssoService, serviceClass, solutionUser, endpoint);

      String endPointUrl = createVapiEndpointURL(serviceSpec);

      SsoClient ssoConnector = (SsoClient) SsoUtil.getConnector(serviceSpec).getConnection();

      return getService(serviceClass, endPointUrl, ssoConnector);
   }


   /**
    * Constructs service end point URL path.
    *
    * @param spec
    *           defines the LDU to be used
    * @return end point URL string
    */
   private static String createVapiEndpointURL(ServiceSpec spec) {
      return String.format(SERVICE_URL_TEMPLATE, spec.endpoint.get());
   }


   /**
    * Method that  gets a connection by specified endpoint
    * @param endPointUrl - url to connect to
    * @return - connection
    * @throws KeyManagementException
    * @throws NoSuchAlgorithmException
    * @throws KeyStoreException
    */
   private static ProtocolConnection getConnection(String endPointUrl)
         throws KeyManagementException, NoSuchAlgorithmException, KeyStoreException {
      // Set a HttpClient with self signed trust strategy SSL context and
      // with all host name verifier.
      // In that way the UI automation does not need to have the SSL
      // certificate.
      FakeSSLProtocolConnectionFactory fakeSSLProtocolConnectionFactory =
            new FakeSSLProtocolConnectionFactory();
      ProtocolFactory protFactory = new ProtocolFactory(
            fakeSSLProtocolConnectionFactory);

      return protFactory.getHttpConnection(endPointUrl, null, null);
   }

}
