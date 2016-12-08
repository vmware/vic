/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Executor;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.sms.ServiceInstance;
import com.vmware.vim.binding.sms.StorageManager;
import com.vmware.vim.binding.sms.Task;
import com.vmware.vim.binding.sms.TaskInfo;
import com.vmware.vim.binding.sms.fault.CertificateNotTrusted;
import com.vmware.vim.binding.sms.provider.Provider;
import com.vmware.vim.binding.sms.provider.ProviderInfo;
import com.vmware.vim.binding.sms.provider.VasaProviderInfo;
import com.vmware.vim.binding.sms.provider.VasaProviderSpec;
import com.vmware.vim.binding.sms.version.version5;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.vmomi.client.Client;
import com.vmware.vim.vmomi.client.http.HttpClientConfiguration;
import com.vmware.vim.vmomi.client.http.ThumbprintVerifier;
import com.vmware.vim.vmomi.client.http.impl.AllowAllThumbprintVerifier;
import com.vmware.vim.vmomi.client.http.impl.HttpConfigurationImpl;
import com.vmware.vim.vmomi.core.Future;
import com.vmware.vim.vmomi.core.RequestContext;
import com.vmware.vim.vmomi.core.Stub;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vim.vmomi.core.impl.RequestContextImpl;
import com.vmware.vim.vmomi.core.types.VmodlType;
import com.vmware.vim.vmomi.core.types.VmodlTypeMap;
import com.vmware.vise.util.concurrent.ThreadPoolFactory;
import com.vmware.vise.util.concurrent.WorkerThreadFactory;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.StorageProviderSpec;

public class StorageProviderBasicSrvApi {

   private static final String SERVICE_INSTANCE_MO_REF_TYPE = "SmsServiceInstance";
   private static final String SERVICE_INSTANCE_MO_REF_ID = "ServiceInstance";
   private static final String VC_SESSION_COOKIE = "vcSessionCookie";
   private static final String SMS_SERVICE_SUBDIR = "/sms/sdk";
   private static final long SMS_TASK_TIMEOUT_IN_MS = 8 * 60 * 1000; // 8 min.

   private static final Logger _logger =
         LoggerFactory.getLogger(StorageProviderBasicSrvApi.class);

   private static Executor _threadPoolExecutor =
         ThreadPoolFactory.newFlexibleThreadPool(
               10,
               new WorkerThreadFactory("sms-service-thread-pool")
               );

   private static StorageProviderBasicSrvApi instance = null;
   protected StorageProviderBasicSrvApi() {}

   /**
    * Get instance of StorageProviderSrvApi.
    *
    * @return  created instance
    */
   public static StorageProviderBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized(StorageProviderBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing StorageProviderSrvApi.");
               instance = new StorageProviderBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Registers new VASA storage provider.
    *
    * @param storageProvider     spec of the storage provider
    * @return                    if the operation was successful or not
    */
   public boolean createStorageProvider(StorageProviderSpec storageProvider) {
      validateStorageProviderSpec(storageProvider);

      try {
         VasaProviderSpec vasaProviderSpec = createVasaProviderSpec(storageProvider);

         // Get the SMS service associated with the VC.
         SmsService smsService = getSmsService(storageProvider);

         // Start task to register the new storage provider
         BlockingFuture<ManagedObjectReference> future =
               new BlockingFuture<ManagedObjectReference>();
         smsService.getStorageManager().registerProvider(vasaProviderSpec, future);
         Task registrationTask = (Task)smsService.getManagedObject(future.get());

         // Wait for task completion
         try {
            waitRegistrationTaskToComplete(registrationTask);
         } catch (CertificateNotTrusted e) {
            // In case of CertificateNotTrusted fault we need to put the certificate reurned by the error and try again.
            byte[] certificateData = getCertificateData(e);
            X509Certificate cert = generateCertificateFromByteArray(certificateData);
            // The SMS service expects the certificate to be passed in as Base64 encoded string.
            byte[] decodedCertificate = Base64.encodeBase64(cert.getEncoded());
            vasaProviderSpec.setCertificate(new String(decodedCertificate));

            // Do the second call
            future = new BlockingFuture<ManagedObjectReference>();
            smsService.getStorageManager().registerProvider(vasaProviderSpec, future);
            ManagedObject mObj = smsService.getManagedObject(future.get());
            Task registrationTaskWithCertificate = (Task) mObj;
            waitRegistrationTaskToComplete(registrationTaskWithCertificate);
         }
      } catch (Exception ex) {
         _logger.error("Failed to register storage provider" + ex);
         return false;
      }

      return true;
   }

   /**
    * Check if specific storage provider is present in the inventory.
    *
    * @param storageProviderSpec
    *           the spec of the provider
    * @return if it can be found in the inventory
    * @throws Exception
    */
   public boolean isStorageProviderPresent(
         StorageProviderSpec storageProviderSpec) throws Exception {
      validateStorageProviderSpec(storageProviderSpec);
      List<ProviderInfo> providers = getStorageProviderInfo(storageProviderSpec);
      for (ProviderInfo provider : providers) {
         if (provider.getName().equals(storageProviderSpec.name.get())) {
            return true;
         }
      }
      return false;
   }

   /**
    * Retrieves storage manager.
    *
    * @return              retrieved storage manager
    * @throws Exception    if storage manager cannot be retrieved
    */
   public StorageManager getStorageManager(EntitySpec entitySpec) throws Exception {
      return getSmsService(entitySpec).getStorageManager();
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Retrieves the VC's correspondent SMS service.
    *
    * @return the VC's {@link SmsService} instance
    * @throws SmsServiceUnableToConnectException, VcException
    */
   private SmsService getSmsService(EntitySpec entitySpec)
         throws SmsServiceUnableToConnectException, VcException {
      VcService vcService = VcServiceUtil.getVcService(entitySpec);

      // Obtain SMS service location
      URI serviceLocation = getSmsServiceLocation(vcService);

      // Create http configuration
      HttpConfigurationImpl httpConfiguration = new HttpConfigurationImpl();

      ThumbprintVerifier thumbprintVerifier = new AllowAllThumbprintVerifier();
      httpConfiguration.setThumbprintVerifier(thumbprintVerifier);

      HttpClientConfiguration httpClientConfiguration =
            HttpClientConfiguration.Factory.newInstance();
      httpClientConfiguration.setExecutor(_threadPoolExecutor );
      httpClientConfiguration.setHttpConfiguration(httpConfiguration);

      // Create client
      Client client = Client.Factory.createClient(
            serviceLocation,
            version5.class,
            httpClientConfiguration);

      // Obtain service content
      ManagedObjectReference serviceMoRef = new ManagedObjectReference(
            SERVICE_INSTANCE_MO_REF_TYPE, SERVICE_INSTANCE_MO_REF_ID);
      ServiceInstance service = client.createStub(ServiceInstance.class, serviceMoRef);

      // It is required that we pass the VC session cookie , otherwise
      // the service will complain about an invalid session.
      RequestContext sessionContext = new RequestContextImpl();
      sessionContext.put(VC_SESSION_COOKIE,
            vcService.getConnectionInfo().getSessionCookie());
      ((Stub)service)._setRequestContext(sessionContext);

      return new SmsServiceImpl(client, service, sessionContext);
   }

   /**
    * Waits storage provider registration task to complete and populates the result
    * and error properties of the {@link OperationResult}.
    * @throws TimeoutException
    * @throws InterruptedException
    * @throws ExecutionException
    * @throws CertificateException
    */
   private void waitRegistrationTaskToComplete(Task task)
         throws CertificateNotTrusted, ExecutionException, InterruptedException, TimeoutException {

      SmsTaskWaiter taskWaiter = new SmsTaskWaiter(task);
      TaskInfo taskInfo = taskWaiter.waitTaskToComplete(SMS_TASK_TIMEOUT_IN_MS);
      // Check task's status
      if (taskInfo.error instanceof CertificateNotTrusted) {
         throw (CertificateNotTrusted) taskInfo.error;
      } else {
         _logger.info("The task finished: " + taskInfo.result);
      }
   }

   /**
    * Converts CertificateNotTrusted exception to byte[] data
    * data object.
    * @throws IOException
    * @throws CertificateException
    * @throws NoSuchAlgorithmException
    */
   private byte[] getCertificateData(CertificateNotTrusted fault)
         throws CertificateException, IOException, NoSuchAlgorithmException {
      // CertificateNotTrusted fault contains the certificate encoded in Base64 string.
      return Base64.decodeBase64(fault.getCertificate().getBytes());
   }

   /**
    * Constructs new X509Certificate from byte array.
    */
   private X509Certificate generateCertificateFromByteArray(byte[] certificate) throws CertificateException,
   IOException {
      CertificateFactory cf = CertificateFactory.getInstance("X.509");
      ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(certificate);
      try {
         X509Certificate cert = (X509Certificate) cf.generateCertificate(byteArrayInputStream);
         return cert;
      } finally {
         byteArrayInputStream.close();
      }
   }

   /**
    * Returns the {@link SmsService} location for the given virtual centre.
    *
    * @param vcService                    current vcService
    * @return                             the {@link SmsService} location.
    * @throws LocationLookupException
    */
   private URI getSmsServiceLocation(VcService vcService) {

      // Retrieves the VC URL.
      String vcServiceUrl = vcService.getServiceUrl();

      if (Strings.isNullOrEmpty(vcServiceUrl)) {
         _logger.error("getSmsServiceLocation: Failed to retrieve the VC URL.");
         throw new IllegalArgumentException("Failed to retrieve the VC URL.");
      }

      // Constructs the SMS service URL by replacing the relative
      // path to the web services with the SMS service directory.
      try {
         URL vcUrl = new URL(vcServiceUrl);
         URL smsUrl = new URL(
               vcUrl.getProtocol(),
               vcUrl.getHost(),
               vcUrl.getPort(),
               SMS_SERVICE_SUBDIR);
         return smsUrl.toURI();
      } catch (URISyntaxException e) {
         _logger.error("URI syntac is not correct! " + e.toString());
         throw new IllegalArgumentException(e);
      } catch (MalformedURLException e) {
         _logger.error("URI is malformed! " + e.toString());
         throw new IllegalArgumentException(e);
      }
   }

   private VasaProviderSpec createVasaProviderSpec(StorageProviderSpec spec){
      // Construct the VasaProviderSpec
      VasaProviderSpec vasaProviderSpec = new VasaProviderSpec();
      vasaProviderSpec.name = spec.name.get();
      vasaProviderSpec.url = spec.providerUrl.get();
      vasaProviderSpec.username = spec.username.get();
      vasaProviderSpec.password = spec.password.get();
      vasaProviderSpec.certificate = null;

      return vasaProviderSpec;
   }


   private void validateStorageProviderSpec(StorageProviderSpec storageProvider) {

      // Validate provider's URL.
      try {
         new URL(storageProvider.providerUrl.get());
      } catch (MalformedURLException e) {
         _logger.error("Failed to register storage provider: invalid provider url" + e);
         throw new IllegalArgumentException("The storage provider URL is not valid");
      }
   }

   /**
    * Helper class that represents a SMS Service associated with the
    * vCenter Server.
    */
   private interface SmsService {

      /**
       * Get the {@link ServiceInstance} of the SMS service.
       * @return {@link ServiceInstance} data object.
       */
      ServiceInstance getServiceInstace();

      /**
       * Gets {@link StorageManager} managed object.
       * @return {@link StorageManager} managed object.
       * @throws SmsServiceUnableToConnectException
       */
      StorageManager getStorageManager() throws SmsServiceUnableToConnectException;

      /**
       * Creates a managed object from the given managed object reference.
       * @param <T> type of the managed object
       * @param moRef managed object reference
       * @return managed object for the given managed object reference
       */
      <T extends ManagedObject> T getManagedObject(ManagedObjectReference moRef);

      /**
       * Closes the connection to the SMS service and releases resources.
       */
      void logout();
   }

   /**
    * Default implementation of {@link SmsService}.
    */
   private class SmsServiceImpl implements SmsService {

      private final Client _vmomiClient;
      private final ServiceInstance _serviceInstance;
      private final VmodlTypeMap _typeMap;
      private StorageManager _storageManager;
      private final RequestContext _sessionContext;

      /**
       * Initializes the SmsServiceImpl instance.
       *
       * @param vmomiClient      the client used to communicate with the SMS service
       * @param serviceInstance  the ServiceInstance of the SMS service
       * @param sessionContext   a RequestContext used to pass the
       *                         session cookie to the SMS service
       */
      public SmsServiceImpl(Client vmomiClient, ServiceInstance serviceInstance,
            RequestContext sessionContext) {
         _vmomiClient = vmomiClient;
         _serviceInstance = serviceInstance;
         _typeMap = VmodlTypeMap.Factory.getTypeMap();
         _sessionContext = sessionContext;
      }

      @Override
      public ServiceInstance getServiceInstace() {
         return _serviceInstance;
      }

      @Override
      public StorageManager getStorageManager() throws SmsServiceUnableToConnectException {
         if (_storageManager == null) {
            BlockingFuture<ManagedObjectReference> future =
                  new BlockingFuture<ManagedObjectReference>();

            _serviceInstance.queryStorageManager(future);
            try {
               _storageManager = getManagedObject(future.get());
            } catch (ExecutionException e) {
               throw new SmsServiceUnableToConnectException(e);
            } catch (InterruptedException e) {
               throw new SmsServiceUnableToConnectException(e);
            }
         }

         return _storageManager;
      }

      @SuppressWarnings("unchecked")
      @Override
      public <T extends ManagedObject> T getManagedObject(ManagedObjectReference moRef) {
         VmodlType vmodlType = _typeMap.getVmodlType(moRef.getType());
         Class<T> typeClass = (Class<T>) vmodlType.getTypeClass();

         T result = _vmomiClient.createStub(typeClass, moRef);
         // It is required that we pass the VC session cookie , otherwise
         // the service will complain about an invalid session.
         ((Stub)result)._setRequestContext(_sessionContext);
         return result;
      }

      @Override
      public void logout() {
         try {
            if (_vmomiClient != null) {
               _vmomiClient.shutdown();
            }
         } catch (Exception ex) {
            _logger.error("Failed to shutdown vlsi client: " + ex.getMessage());
         }
      }
   }

   /**
    * Exception, thrown when VC connection problem occurs.
    */
   public class SmsServiceUnableToConnectException extends Exception {

      private final long serialVersionUID = 1L;

      /**
       * Constructs a new exception with the specified cause.
       *
       * @param cause
       */
      public SmsServiceUnableToConnectException(Throwable cause) {
         super("Service Unavailable", cause);
      }
   }

   /**
    * Helper class used to monitor single {@link com.vmware.vim.binding.sms.Task}
    */
   private class SmsTaskWaiter {

      private final String TASK_STATE_RUNNING = "running";
      private final long SLEEP_INTERVAL_IN_MILLIS = 1000; // 1 sec
      private final long MIN_TIMEOUT_IN_MILLIS = 1000; // 1 sec
      private final long INFINITE_TIMEOUT = -1;

      private final Task _task;

      /**
       * Constructs new instance.
       * @param task
       *    {@link com.vmware.vim.binding.sms.Task} which will be monitored.
       */
      public SmsTaskWaiter(Task task) {
         _task = task;
      }

      /**
       * Waits the SMS task to complete and returns its {@link TaskInfo}.
       *
       * @param timeout
       *       The timeout in milliseconds to wait the task to complete. If the
       * task runs longer than the timeout a TimeoutException is thrown.
       * @throws InterruptedException
       * @throws ExecutionException
       * @throws TimeoutException
       */
      public TaskInfo waitTaskToComplete(long timeout)
            throws ExecutionException, InterruptedException, TimeoutException {
         BlockingFuture<TaskInfo> future = new BlockingFuture<TaskInfo>();
         long startTime = System.currentTimeMillis();
         TaskInfo taskInfo;
         while (true) {
            _task.queryInfo(future);
            if (timeout != INFINITE_TIMEOUT) {
               // Retrieving task's info can take longer than the specified timeout, so
               // we need to set a timeout to the queryInfo server call as well.
               // The correct timeout here is calculated by subtracting the elapsed time
               // from the original timeout.
               long queryInfoTimeout = timeout - (System.currentTimeMillis() - startTime);
               // We don't want negative timeout
               queryInfoTimeout = Math.max(MIN_TIMEOUT_IN_MILLIS, queryInfoTimeout);
               taskInfo = future.get(queryInfoTimeout, TimeUnit.MILLISECONDS);
            } else {
               taskInfo = future.get();
            }
            if (TASK_STATE_RUNNING.equals(taskInfo.state)) {
               long elapsedTimeMillis = System.currentTimeMillis() - startTime;
               if (timeout != INFINITE_TIMEOUT && timeout < elapsedTimeMillis) {
                  throw new TimeoutException("Wait timeout exception!");
               }
               Thread.sleep(SLEEP_INTERVAL_IN_MILLIS);
            } else {
               return taskInfo;
            }
         }
      }
   }

   /**
    * Filter the list of provider info search for vasa provider information
    *
    * @param providersInfo
    * @param storageProviderSpec
    * @return
    * return null in case no vasa provider info is found
    */
   public VasaProviderInfo getVasaStorageProviderInfo(
         List<ProviderInfo> providersInfo,
         StorageProviderSpec storageProviderSpec) {

      for (ProviderInfo providerInfo : providersInfo) {
         if (providerInfo instanceof VasaProviderInfo) {
            VasaProviderInfo currentProviderInfo = (VasaProviderInfo) providerInfo;
            if (storageProviderSpec.name.get().equals(currentProviderInfo.name)) {
               return currentProviderInfo;
            }
         }
      }
      return null;
   }

   /**
    * Get info for VASA storage provider
    *
    * @param storageProviderSpec
    * @return
    * @throws Exception
    */
   public List<ProviderInfo> getStorageProviderInfo(
         StorageProviderSpec storageProviderSpec) throws Exception {
      SmsService smsService = getSmsService(storageProviderSpec);
      Future<ManagedObjectReference[]> future = new BlockingFuture<ManagedObjectReference[]>();

      smsService.getStorageManager().queryProvider(future);
      ManagedObjectReference[] providerMors = future.get();

      List<ProviderInfo> providersInfo = new ArrayList<>();
      if (providerMors != null) {
         for (ManagedObjectReference providerMor : providerMors) {
            Provider provider = smsService.getManagedObject(providerMor);
            Future<ProviderInfo> providerInfoFuture = new BlockingFuture<ProviderInfo>();
            provider.queryProviderInfo(providerInfoFuture);
            ProviderInfo providerInfo = providerInfoFuture.get();
            providersInfo.add(providerInfo);

         }
      }
      return providersInfo;
   }

   /**
    * Unregister VASA storage provider, throws an exception if no storage
    * provider is found.
    *
    * @param storageProviderSpec
    * @return
    * @throws Exception
    */
   public void deleteStorageProvider(StorageProviderSpec storageProviderSpec) throws Exception {
      validateStorageProviderSpec(storageProviderSpec);

         List<ProviderInfo> providersInfo = getStorageProviderInfo(storageProviderSpec);
         VasaProviderInfo vasaProviderInfo = getVasaStorageProviderInfo(
               providersInfo, storageProviderSpec);

         if (vasaProviderInfo == null) {
            throw new RuntimeException(
                  String.format("Failed to retrieve VasaProviderInfo %s",
                        storageProviderSpec));
         }

         // Start task to unregister the new storage provider
         SmsService smsService = getSmsService(storageProviderSpec);
         Future<ManagedObjectReference> unregisterFuture = new BlockingFuture<ManagedObjectReference>();
         smsService.getStorageManager().unregisterProvider(
               vasaProviderInfo.providerId, unregisterFuture);

         Task registrationTask = (Task) smsService
               .getManagedObject(unregisterFuture.get());

         waitRegistrationTaskToComplete(registrationTask);
   }
}
