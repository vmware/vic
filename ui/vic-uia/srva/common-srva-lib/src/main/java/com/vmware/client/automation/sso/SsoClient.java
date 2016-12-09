/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.sso;

import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.util.concurrent.Executors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.cis.cm.client.ComponentManagerClient;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.exception.SsoException;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.servicespec.VcServiceSpec;
import com.vmware.client.automation.util.QueryClientWrapper;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.cis.cm.SearchCriteria;
import com.vmware.vim.binding.cis.cm.ServiceEndPoint;
import com.vmware.vim.binding.cis.cm.ServiceInfo;
import com.vmware.vim.binding.cis.cm.ServiceType;
import com.vmware.vim.binding.impl.cis.cm.SearchCriteriaImpl;
import com.vmware.vim.binding.pbm.ServiceInstance;
import com.vmware.vim.binding.pbm.ServiceInstanceContent;
import com.vmware.vim.binding.pbm.profile.ProfileManager;
import com.vmware.vim.binding.sso.admin.PrincipalManagementService;
import com.vmware.vim.binding.vim.ServiceDirectory.ServiceEndpoint;
import com.vmware.vim.binding.vim.version.version11;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.query.client.QueryAuthenticationManager;
import com.vmware.vim.query.client.impl.ClientImpl;
import com.vmware.vim.sso.client.SamlToken;
import com.vmware.vim.sso.client.TokenSpec;
import com.vmware.vim.sso.client.TokenSpec.DelegationSpec;
import com.vmware.vim.sso.client.exception.InvalidTokenException;
import com.vmware.vim.vmomi.client.Client;
import com.vmware.vim.vmomi.client.http.HttpClientConfiguration;
import com.vmware.vim.vmomi.client.http.HttpConfiguration;
import com.vmware.vim.vmomi.client.http.ThumbprintVerifier;
import com.vmware.vim.vmomi.client.http.impl.AllowAllThumbprintVerifier;
import com.vmware.vim.vmomi.client.http.impl.HttpConfigurationImpl;
import com.vmware.vim.vmomi.core.RequestContext;
import com.vmware.vim.vmomi.core.Stub;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vim.vmomi.core.impl.RequestContextImpl;
import com.vmware.vim.vmomi.core.types.VmodlContext;
import com.vmware.vim.vmomi.core.types.VmodlType;
import com.vmware.vim.vmomi.core.types.VmodlTypeMap;
import com.vmware.vise.vim.cm.ComponentManagerService;
import com.vmware.vise.vim.commons.ServiceEndpointType;
import com.vmware.vise.vim.commons.vcservice.LoginSpec;
import com.vmware.vise.vim.commons.vcservice.ServiceEndpointEx;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vise.vim.commons.vcservice.VcServiceConnectionInfo;
import com.vmware.vise.vim.commons.vcservice.impl.VcServiceImpl;
import com.vmware.vise.vim.security.sso.SsoAdminService;
import com.vmware.vise.vim.security.sso.SsoService;
import com.vmware.vise.vim.security.sso.SsoUtil;
import com.vmware.vise.vim.security.sso.exception.SsoServiceException;
import com.vmware.vise.vim.security.sso.impl.SsoAdminServiceImpl;
import com.vmware.vise.vim.security.sso.impl.SsoCmLocatorImpl;
import com.vmware.vise.vim.security.sso.impl.SsoServiceImpl;

/**
 * Used for authenticating against SSO.
 */
// TODO: rename to SSOconnection or similar as it keeps the sso connection
// instance
public class SsoClient {

   private static final Logger _logger = LoggerFactory.getLogger(SsoClient.class);

   /** The last acquired SSO token. rado */
   private SamlToken _token;

   /** The SSO service instance. */
   private SsoService _ssoService;

   private final ServiceSpec serviceSpec;

   /** The SSO Admin service instance. */
   private SsoAdminService _ssoAdminService;

   /** Default SSO token duration in seconds - 12 hours. */
   private static final long DEFAULT_TOKEN_DURATION = 43200;

   private static final String CM_URL_TEMPLATE = "https://%s/cm/sdk";

   private static final String VC_SERVICE_URL_FORMAT = "https://%s/sdk";

   private static final String VC_SERVICE_PBM_ENDPOINT = "/pbm/sdk";

   private static final String ENDPOINT_TYPE_INV_SVC = "com.vmware.cis.inventory";

   // vc service connection
   private VcService _vcService = null;

   private VcServiceConnectionInfo _connInfo = null;

   static {
      // Initializes vmodl types in a non-osgi container.
      // Internally it uses spring ClassPathXmlApplicationContext to load
      // those types from file com/vmware/vim/binding/ssocontext.xml.
      // This file should be available in runtime, and it is provided by
      // sso-adminserver-client-bindings jar.
      // VmodlContext.initContext(new String[] { "com.vmware.vim.binding.sso",
      // "com.vmware.vim.binding.vim", "com.vmware.vim.binding.pbm" });

      // Initializes vmodl types in a non-osgi container.
      // Internally it uses spring ClassPathXmlApplicationContext to load
      // those types from file com/vmware/vim/binding/ssocontext.xml.
      // This file should be available in runtime, and it is provided by
      // sso-adminserver-client-bindings jar.
      VmodlContext.initContext(new String[] { "com.vmware.vim.binding.vim",
            "com.vmware.vim.binding.pbm", "com.vmware.vim.binding.sso", "com.vmware.vim.binding.sms" });
   }

   public SsoClient(ServiceSpec spec) throws SsoException {

      this.serviceSpec = spec;

      if (serviceSpec instanceof VcServiceSpec) {
         String cmServiceUrl =
               String.format(CM_URL_TEMPLATE, serviceSpec.endpoint.get());

         try {
            ComponentManagerClient cmClient = getComponentManager();
            ServiceInfo ssoServiceInfo = cmClient.lookupSso(true);

            ComponentManagerService cmService =
                  new DefaultComponentManagerServiceImpl(cmServiceUrl, ssoServiceInfo);
            DefaultKeystoreService keystoreService =
                  new DefaultKeystoreService(ssoServiceInfo);
            SsoCmLocatorImpl ssoCmLocator =
                  new SsoCmLocatorImpl(cmService, keystoreService);

            _ssoService = new SsoServiceImpl(ssoCmLocator, true, null, keystoreService);
         } catch (Exception e) {
            throw new SsoException(e);
         }

      }

   }

   /**
    * Acquires an SSO token for a given user with given password.
    *
    * @param username
    *           the username to be used
    * @param password
    *           the password to be used
    * @throws SsoException
    */
   public void connect() throws SsoException {
      try {
         _token = null;

         DelegationSpec delegationspec = new DelegationSpec(true);
         TokenSpec tokenSpec = SsoUtil.buildTokenSpec(DEFAULT_TOKEN_DURATION, delegationspec, true);

         _token = _ssoService.acquireToken(
               serviceSpec.username.get(), serviceSpec.password.get(), tokenSpec);

         if (_ssoAdminService != null) {
            // invalidating the cached sso admin service because the token has changed
            _ssoAdminService = null;
         }
      } catch (SsoServiceException e) {
         throw new SsoException(e);
      }
   }

   public boolean isAlive() {
      try {
         return !_ssoService.checkTokenLifetime(_token);
      } catch (java.lang.IllegalArgumentException e) {
         return false;
      } catch (InvalidTokenException e) {
         return false;
      }
   }

   public void disconnect() {
      // TODO: It is registry job
      // if (_vcService != null && isVcConnectionAlive()) {
      // _logger.info("Log out of the VC server");
      // _vcService.logout();
      // }
      //
      // for (VcService service : _vcServicesByUrl.values()) {
      // if (isVcConnectionAlive(service)) {
      // _logger.info("Log out of the VC server");
      // ((VcServiceImpl) service).logout();
      // }
      // }
      throw new RuntimeException("Not implemented!");
   }

   // getters for different services
   // TODO rkovachev move them to proper util class.

   /**
    * Returns an instance of {@link SsoAdminService} authenticated with the last
    * issued token. If there is no such token, the returned admin service will
    * allow calling only methods that do not require authentication.
    *
    * The {@code SsoAdminService} exposes SSO managing services via getters, for
    * example {@code SsoAdminService.getDomainManagementService()}.
    *
    * @return an instance of {@link SsoAdminService}
    * @throws SsoException
    * @see #loginUsingPassword
    * @see com.vmware.vise.vim.security.sso.SsoAdminService
    */
   public SsoAdminService getSsoAdminService() throws SsoException {
      if (_ssoAdminService == null) {
         try {
            _ssoAdminService =
                  new SsoAdminServiceImpl(_ssoService.getServerInfo(), _ssoService,
                        _token);
         } catch (Exception e) {
            throw new SsoException(e);
         }
      }
      return _ssoAdminService;
   }

   /**
    * Retrieves principal management for specific vCenter instance.
    * It provides CRUD operations for SSO users and groups.
    *
    * @return retrieved principal management
    *
    * @throws SsoException
    * @throws SsoServiceException
    */
   public PrincipalManagementService getPrincipalManagement()
         throws SsoServiceException, SsoException {

      return getSsoAdminService().getPrincipalManagementService();
   }

   /**
    * Method that retrieves the endpoint url from a service in component manager
    * 
    * @param serviceSpec - spec for vcenter connection
    * @param svcType - type of the service, i.e. inventory service, dataservice, etc.
    * @param typeId - endpoint, whose url is needed
    * @return - specified endpoint url
    * @throws Exception - in case connection to vcenter is not possible, the service
    *            is not found or the endpoint is invalid
    */
   public String getEndPointUrlByServiceTypeFromCm(ServiceType svcType, String typeId) throws Exception {
      ComponentManagerClient cmClient = getComponentManager();
      cmClient.login(this.getLastToken(), this.getPrivateKey());
      // lookup inventory service
      SearchCriteria searchCriteria = new SearchCriteriaImpl();
      searchCriteria.setServiceType(svcType);
      ServiceInfo[] svcInfos = cmClient.lookup(searchCriteria);

      String expectedHostName = com.vmware.client.automation.util.SsoUtil.resolveHost(serviceSpec.endpoint.get());

      // We may have multuiple service infos in multi-vc environment
      ServiceInfo svcInfo = null;
      if (svcInfos.length == 0) {
         throw new Exception("Invalid service info received for service");
      } else if (svcInfos.length == 1) {
         svcInfo = svcInfos[0];
      } else {
         for (ServiceInfo info : svcInfos) {
            for (ServiceEndPoint endpoint : info.getServiceEndPoints()) {
               if (endpoint.getUrl().getHost().equals(expectedHostName)) {
                  svcInfo = info;
                  break;
               }
            }
            if (svcInfo != null) {
               break;
            }
         }

         if (svcInfo == null) {
            throw new Exception("Could not find expected service info");
         }
      }

      for (ServiceEndPoint svcEndPoint : svcInfo.getServiceEndPoints()) {
         if (svcEndPoint.getEndPointType().getTypeId().equals(typeId)) {
            return svcEndPoint.getUrl().toString();
         }
      }
      throw new IllegalArgumentException("No such endpoint found!");
   }

   public VcService getVcService() throws VcException {

      if (_vcService == null) {
         try {
            _vcService = new VcServiceImpl(buildLoginSpec());
            _connInfo = ((VcServiceImpl) _vcService).login();
         } catch (Exception e) { // Thrown by the login method
            throw new VcException(e.getMessage(), e.getCause());
         }
      }

      if (!_connInfo.getConnectionState()) {
         throw new VcException("Connection to the VC not established.");
      }

      return _vcService;
   }

   // TODO: fix it to use the CM provided by the Platform
   public ComponentManagerClient getComponentManager() throws Exception {
      KeyStore trustStore = null;
      ThumbprintVerifier verifier = new AllowAllThumbprintVerifier();
      String cmUrl =
            String.format(String.format(CM_URL_TEMPLATE, serviceSpec.endpoint.get()));

      // Initialize component manager client
      try {
         return new ComponentManagerClient(new URI(cmUrl), trustStore, verifier);

      } catch (URISyntaxException e) {
         _logger.error("Unable to create CM client!", e);
         throw new Exception(e);
      }
   }

   /**
    * Creates a query client using and SSO connection to the VC.
    *
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return the client
    */
   // TODO: refactor it to get properly the query service
   public com.vmware.vim.query.client.Client getQueryClient() throws Exception {
      // ComponentManagerClient cmClient = this.getComponentManager();
      // cmClient.login(this.getLastToken(), this.getPrivateKey());
      //
      // ServiceInfo invSvcInfo = cmClient.lookupInventory(true);
      // invSvcInfo.getServiceEndPoints();
      // ServiceEndPoint invSvcEndpoint =
      // getServiceEndpointByType(invSvcInfo, ENDPOINT_TYPE_INV_SVC);
      //
      // ServiceEndpointEx serviceEndpoint =
      // (ServiceEndpointEx) getInventoryServiceEndpoint();
      //
      // String searchUrl = invSvcEndpoint.getUrl().toString();
      // serviceEndpoint =
      // new ServiceEndpointEx(serviceEndpoint.getLduGuid(),
      // serviceEndpoint.getServiceEndpointType(), serviceEndpoint.getKey(),
      // serviceEndpoint.getInstanceUuid(), serviceEndpoint.getInstanceName(),
      // serviceEndpoint.getVcInstanceId(), serviceEndpoint.getProtocol(),
      // invSvcEndpoint.getUrl().toString(),
      // serviceEndpoint.getSslThumbprint(), serviceEndpoint.getCertificate());
      //
      // com.vmware.vise.search.LoginSpec spec = new com.vmware.vise.search.LoginSpec();
      // spec.webServicesUrl = searchUrl;
      // spec.loginMethod = LoginMethod.SAML_TOKEN;
      // spec.ssoToken = getLastToken();
      // spec.clientPrivateKey = getPrivateKey();
      //
      // // We are using dirty reflection calls here because the devs have the API
      // // methods with either "default" visibility
      // // or "private" visibility.
      // Class<?> queryClientUtil =
      // Class.forName("com.vmware.vise.search.transport.impl.QueryClientUtil");
      //
      // Method doGetClient =
      // queryClientUtil.getDeclaredMethod(
      // "doGetClient",
      // ServiceEndpoint.class,
      // com.vmware.vise.search.LoginSpec.class);
      //
      // doGetClient.setAccessible(true);
      // com.vmware.vim.query.client.Client queryClient =
      // (com.vmware.vim.query.client.Client) doGetClient.invoke(
      // queryClientUtil,
      // serviceEndpoint,
      // spec);
      //
      // return new QueryClientWrapper(queryClient);

      ComponentManagerClient cmClient = getComponentManager();
      cmClient.login(this.getLastToken(), this.getPrivateKey());
      ServiceInfo invSvcInfo = cmClient.lookupInventory(true);

      ServiceEndPoint invSvcEndpoint = getServiceEndpointByType(invSvcInfo, ENDPOINT_TYPE_INV_SVC);

      URI queryServiceUrl = invSvcEndpoint.getUrl();

      HttpConfiguration httpConfig = new HttpConfigurationImpl();
      httpConfig.setThumbprintVerifier(ThumbprintVerifier.Factory.createAllowAllThumbprintVerifier());

      @SuppressWarnings("deprecation")
      com.vmware.vim.query.client.Client clientHack = new ClientImpl(queryServiceUrl, httpConfig, null, null, null);

      authenticateClient(clientHack);
      return new QueryClientWrapper(clientHack);
   }

   private void authenticateClient(com.vmware.vim.query.client.Client client) throws Exception {
      QueryAuthenticationManager authMgr = client.getAuthenticationManager();

      authMgr.loginBySamlToken(getLastToken(), getPrivateKey());
   }

   /**
    * Returns the last acquired SSO token.
    *
    * @return the last acquired token by calling {@code loginUsingPassword} or {@code null} if the last call failed
    */
   public com.vmware.vim.sso.client.SamlToken getLastToken() {
      return _token;
   }

   /**
    * Returns the SSO service private key.
    *
    * @return the HoK private key used for current SSO connection
    */
   public PrivateKey getPrivateKey() {
      return _ssoService.getHokPrivateKey();
   }

   /**
    * Creates PBM service client by using existing VC connection. If VC
    * connection is not present it will initialize one.
    *
    * @return PBM service client instance
    * @throws VcException
    *            if login to the VC server fails
    */
   private com.vmware.vim.vmomi.client.Client getPbmVsliClient() throws VcException {
      // Obtain PBM service location
      URI serviceLocation = getPbmServiceLocation(getVcService());

      // Create http configuration
      HttpConfiguration httpConfiguration = new HttpConfigurationImpl();

      // Set the SSL thumbprint verifier
      ThumbprintVerifier thumbprintVerifier = new AllowAllThumbprintVerifier();
      httpConfiguration.setThumbprintVerifier(thumbprintVerifier);

      HttpClientConfiguration httpClientConfiguration =
            HttpClientConfiguration.Factory.newInstance();
      httpClientConfiguration.setExecutor(Executors.newFixedThreadPool(10));
      httpClientConfiguration.setHttpConfiguration(httpConfiguration);

      // Create client
      return com.vmware.vim.vmomi.client.Client.Factory.createClient(
            serviceLocation,
            com.vmware.vim.binding.pbm.version.version1.class,
            httpClientConfiguration);
   }

   /**
    * Obtains profile manager.
    *
    * @return initialized profile manager
    * @throws VcException
    *            if connection to VC does not succeeds
    */
   public ProfileManager getProfileManager() throws VcException {
      Client client = this.getPbmVsliClient();
      ServiceInstanceContent content = getServiceInstanceContent(getPbmService(client));
      return getManagedObject(content.getProfileManager(), client);
   }

   /**
    * Obtains PBM service from the PBM connection client.
    *
    * @param client
    *           initialized client that connects to PBM folder endpoint
    * @return the obtained PBM service
    * @throws VcException
    *            if connection to the VC cannot be established
    */
   private ServiceInstance getPbmService(Client client) throws VcException {
      return createStub("PbmServiceInstance", "ServiceInstance", client);
   }

   /**
    * Obtains PBM service instance content.
    *
    * @param pbmService
    *           PBM service, whose content will be retrieved
    * @return ` the obtained service content
    */
   private ServiceInstanceContent getServiceInstanceContent(ServiceInstance pbmService) {

      ServiceInstanceContent content = null;
      BlockingFuture<ServiceInstanceContent> future =
            new BlockingFuture<ServiceInstanceContent>();
      pbmService.getContent(future);
      try {
         content = future.get();
      } catch (Exception e) {
         _logger.error("Failed to obtain PBM service content!");
         throw new RuntimeException(e);
      }
      return content;
   }

   private ServiceEndPoint getServiceEndpointByType(ServiceInfo svcInfo,
         String endpointType) {
      for (ServiceEndPoint endpoint : svcInfo.getServiceEndPoints()) {
         if (endpoint.getEndPointType() != null
               && endpoint.getEndPointType().getTypeId() != null
               && endpoint.getEndPointType().getTypeId().equals(endpointType)) {
            return endpoint;
         }
      }
      // nothing found
      throw new RuntimeException("Inventory Service endpoint was not found");
   }

   private <T extends ManagedObject> T getManagedObject(ManagedObjectReference moRef,
         Client client) throws VcException {

      RequestContext sessionContext = new RequestContextImpl();
      sessionContext.put("vcSessionCookie", VcServiceUtil.getVcService(serviceSpec)
            .getConnectionInfo().getSessionCookie());

      VmodlTypeMap typeMap = VmodlTypeMap.Factory.getTypeMap();
      VmodlType vmodlType = typeMap.getVmodlType(moRef.getType());
      @SuppressWarnings("unchecked")
      Class<T> typeClass = (Class<T>) vmodlType.getTypeClass();

      T result = client.createStub(typeClass, moRef);
      // It is required that we pass the VC session cookie , otherwise
      // the service will complain about an invalid session.
      ((Stub) result)._setRequestContext(sessionContext);
      return result;
   }

   @SuppressWarnings("unchecked")
   private <T extends ManagedObject> T createStub(String moRefType, String moRefId,
         Client client) throws VcException {
      ManagedObjectReference moRef = new ManagedObjectReference(moRefType, moRefId);
      return (T) getManagedObject(moRef, client);
   }

   /**
    * Obtains PBM service endpoint URL from existing VC connection.
    *
    * @param vcService
    *           initialized VC connection
    */
   // TODO: RR: refactor it
   private URI getPbmServiceLocation(VcService vcService) {

      // Retrieves URL of the VC.
      String vcServiceUrl = vcService.getServiceUrl();

      if (Strings.isNullOrEmpty(vcServiceUrl)) {
         _logger.error("getPbmServiceLocation: Failed to retrieve the VC URL.");
         throw new RuntimeException("Cannot retrieve VC URL from VC connection.");
      }

      // Constructs the PBM service URL by replacing the relative
      // path to the web services with the PBM service directory.
      try {
         URL vcUrl = new URL(vcServiceUrl);
         URL pbmUrl =
               new URL(vcUrl.getProtocol(), vcUrl.getHost(), vcUrl.getPort(),
                     VC_SERVICE_PBM_ENDPOINT);
         return pbmUrl.toURI();
      } catch (URISyntaxException e) {
         _logger.error("wrong URI to pbm directory");
         throw new IllegalArgumentException(e);
      } catch (MalformedURLException e) {
         _logger.error("wrong URI to pbm directory");
         throw new IllegalArgumentException(e);
      }
   }

   /**
    * Creates a LoginSpec instance with the specified VC connection settings
    */
   // TODO: use cm instead of creating service url.
   private LoginSpec buildLoginSpec() {
      LoginSpec loginSpec = new LoginSpec();
      loginSpec.userName = serviceSpec.username.get();
      loginSpec.password = serviceSpec.password.get();
      loginSpec.serviceUrl =
            String.format(VC_SERVICE_URL_FORMAT, serviceSpec.endpoint.get());
      loginSpec.ignoreSslThumbprint = true;
      loginSpec.vmodlVersion = version11.class;
      // Use default endpoints
      loginSpec.endpoints = new ServiceEndpointEx[0];

      return loginSpec;
   }

   private ServiceEndpoint getInventoryServiceEndpoint() throws VcException {
      // browse the VC endpoints to find the endpoint of the Inventory Service
      for (ServiceEndpoint endpoint : this.getVcService().fetchServiceEndpoints()) {
         // check if this endpoint is inventory service
         if (endpoint instanceof ServiceEndpointEx) {
            ServiceEndpointType endpointType =
                  ((ServiceEndpointEx) endpoint).getServiceEndpointType();
            if (endpointType == ServiceEndpointType.IS) {
               return endpoint;
            }
         }
      }
      // nothing found
      throw new RuntimeException("Inventory Service endpoint was not found");
   }
}
