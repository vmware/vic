/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util;

import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.concurrent.Executors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.sso.SsoClient;
import com.vmware.vim.binding.impl.vmodl.TypeNameImpl;
import com.vmware.vim.binding.vim.Task;
import com.vmware.vim.binding.vim.TaskInfo;
import com.vmware.vim.binding.vim.TaskInfo.State;
import com.vmware.vim.binding.vim.option.OptionManager;
import com.vmware.vim.binding.vim.option.OptionValue;
import com.vmware.vim.binding.vim.view.ViewManager;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.binding.vmodl.query.InvalidCollectorVersion;
import com.vmware.vim.binding.vmodl.query.InvalidProperty;
import com.vmware.vim.binding.vmodl.query.PropertyCollector;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.Change;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.FilterSpec;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.FilterUpdate;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.ObjectSpec;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.ObjectUpdate;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.PropertySpec;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.SelectionSpec;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.TraversalSpec;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.UpdateSet;
import com.vmware.vim.binding.vmodl.query.PropertyCollector.WaitOptions;
import com.vmware.vim.vmomi.client.common.ProtocolBinding;
import com.vmware.vim.vmomi.client.common.Session;
import com.vmware.vim.vmomi.client.http.HttpClientConfiguration;
import com.vmware.vim.vmomi.client.http.HttpConfiguration;
import com.vmware.vim.vmomi.client.http.ThumbprintVerifier;
import com.vmware.vim.vmomi.client.http.impl.AllowAllThumbprintVerifier;
import com.vmware.vim.vmomi.client.http.impl.HttpConfigurationImpl;
import com.vmware.vise.vim.commons.vcservice.VcService;

/**
 * Provides functionality to connect to a Virtual Center server through the Vim
 * API.
 */
public class VcServiceUtil {

   private static final Logger _logger = LoggerFactory.getLogger(VcServiceUtil.class);

   // ---------------------------------------------------------------------------
   // Private Constants
   private static final String VC_SERVICE_PBM_ENDPOINT = "/pbm/sdk";

   private static final String VC_HOSTNAME_OPTION_KEY = "config.vpxd.hostnameUrl";

   // ---------------------------------------------------------------------------
   // Instance Variables

   // Hide constructor
   private VcServiceUtil() {
   }

   // ---------------------------------------------------------------------------
   // Class Methods

   /**
    * Returns a <code>VcService</code> instance in connected state. Note that
    * the method expects that the VC login settings are preliminary set using
    * the <code>setLoginCredentials()</code> method.
    *
    * @param entitySpec
    *           get reference to specification of the LDU connection details.
    *
    * @return <code>VcService</code> object which might already exist or be
    *         created from scratch depending on the
    *         <code>{@link serviceSpec}</code> argument.
    */
   public static VcService getVcService(EntitySpec entitySpec) {
      return getVcService(entitySpec.service.get());
   }

   /**
    * Returns a <code>VcService</code> instance in connected state. Note that
    * the method expects that the VC login settings are preliminary set using
    * the <code>setLoginCredentials()</code> method.
    *
    * @param serviceSpec
    *           specification of the LDU connection details.
    *
    * @return <code>VcService</code> object which might already exist or be
    *         created from scratch depending on the
    *         <code>{@link serviceSpec}</code> argument.
    */
   public static VcService getVcService(ServiceSpec serviceSpec) {
      TestbedConnector connector = SsoUtil.getConnector(serviceSpec);
      SsoClient ssoClient = connector.getConnection();
      try {
         return ssoClient.getVcService();
      } catch (VcException e) {
         throw new RuntimeException(String.format(
               "Error retrieving vc service for %s", serviceSpec), e);
      }
   }

   /**
    * Waits for the specified task to complete on the VC server. Returns
    * immediately if the task is in state "success" or "error".
    *
    * @param taskMoRef
    *           MoRef of the VC task.
    *
    * @return <code>State</code> instance representing the task state.
    *
    * @throws VcException
    *            If querying the task status fails.
    * @deprecated use waitForTaskCompletion(ManagedObjectReference taskMoRef,
    *             ServiceSpec serviceSpec)
    */
   // @Deprecated
   // public static State waitForTaskCompletion(ManagedObjectReference
   // taskMoRef)
   // throws VcException {
   // return waitForTaskCompletion(taskMoRef, null);
   // }

   /**
    * Waits for the specified task to complete on the VC server. Returns
    * immediately if the task is in state "success" or "error".
    *
    * @param taskMoRef
    *           MoRef of the VC task.
    * @param serviceSpec
    *           specification of the LDU connection details.
    *
    * @return <code>State</code> instance representing the task state.
    *
    * @throws VcException
    *            If querying the task status fails.
    */
   public static State waitForTaskCompletion(ManagedObjectReference taskMoRef,
         ServiceSpec serviceSpec) throws VcException {

      if (taskMoRef == null) {
         throw new IllegalArgumentException("taskMoRef not set");
      }

      Task task = getVcService(serviceSpec).getManagedObject(taskMoRef);

      State taskState = task.getInfo().getState();

      if (!isTaskComplete(taskState)) {
         try {
            taskState = waitForUpdatedTaskState(taskMoRef, getVcService(serviceSpec));
         } catch (Exception e) {
            _logger.error("Exception while quering task state: " + e.getMessage());
            throw new VcException(e.getMessage(), e.getCause());
         }
      }

      return taskState;
   }

   /**
    * Waits for the specifies task to complete on the VC server and checks
    * whether it finished with "success" state.
    *
    * @param taskMoRef
    *           MoRef of the VC task.
    *
    * @return True if the task completed with status "success", false otherwise
    *
    * @throws VcException
    *            If querying the task status fails
    * @deprecated Use waitForTaskSuccess(ManagedObjectReference taskMoRef,
    *             ServiceSpec serviceSpec)
    */
   //   @Deprecated
   //   public static boolean waitForTaskSuccess(ManagedObjectReference taskMoRef)
   //         throws VcException {
   //      return waitForTaskSuccess(taskMoRef, null);
   //   }

   /**
    * Waits for the specifies task to complete on the VC server and checks
    * whether it finished with "success" state.
    *
    * @param taskMoRef
    *           MoRef of the VC task.
    * @param serviceSpec
    *           specification of the LDU connection details.
    *
    * @return True if the task completed with status "success", false otherwise
    *
    * @throws VcException
    *            If querying the task status fails
    */
   public static boolean waitForTaskSuccess(ManagedObjectReference taskMoRef,
         ServiceSpec serviceSpec) throws VcException {
      State taskState = waitForTaskCompletion(taskMoRef, serviceSpec);

      if (taskState.name().equals(TaskInfo.State.success.name())) {
         return true;
      }

      return false;
   }

   public static boolean waitForTaskSuccess(ManagedObjectReference taskMoRef,
         EntitySpec entitySpec) throws VcException {
      return waitForTaskSuccess(taskMoRef, entitySpec.service.get());
   }

   /**
    * Creates PBM service client by using existing VC connection.
    * If VC connection is not present it will initialize one.
    *
    * @return              PBM service client instance
    * @throws VcException  if login to the VC server fails
    */
   public static com.vmware.vim.vmomi.client.Client getPbmVsliClient(
         EntitySpec entitySpec) throws VcException {
      // Obtain PBM service location
      URI serviceLocation = getPbmServiceLocation(getVcService(entitySpec));

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
            com.vmware.vim.binding.pbm.version.versions.PBM_VERSION_NEWEST,
            httpClientConfiguration);
   }

   /**
    * Return the session cookie for the current session
    * @throws VcException
    */
   public static String getSessionCookie(ServiceSpec serviceSpec) throws VcException {

      com.vmware.vim.vmomi.client.Client client =
            getVcService(serviceSpec).getVmomiClient();

      if (client == null) {
         return null;
      }
      try {
         ProtocolBinding binding = client.getBinding();
         if (binding == null) {
            _logger.error("Null protocol binding returned by vc client");
            return null;
         }
         Session session = binding.getSession();
         if (session == null) {
            _logger.error("Null session returned by protocol binding");
            return null;
         }
         String sessionCookie = session.getId();
         return sessionCookie;
      } catch (Exception e) {
         _logger.error("Error in getSessionCookie");
      }
      return null;
   }

   /**
    * Retrieves VC hostname.
    *
    * @param serviceSpec   VC service spec that determines instance connection, optional
    * @return              VC hostname
    * @throws VcException  if unable to connect to the VC
    */
   public static String getVcHostname(ServiceSpec serviceSpec) throws VcException {
      VcService vcService = getVcService(serviceSpec);
      ManagedObjectReference setting =
            vcService.getServiceInstanceContent().getSetting();
      OptionManager optionManager = vcService.getManagedObject(setting);

      OptionValue[] settings = optionManager.getSetting();

      String result = null;
      if (settings != null) {
         for (OptionValue optionVal : settings) {
            if (optionVal.getKey().equals(VC_HOSTNAME_OPTION_KEY)) {
               result = (String) optionVal.getValue();
            }
         }
      }

      return result;
   }


   // ---------------------------------------------------------------------------
   // Private Helpers

   /**
    * Obtains PBM service endpoint URL from existing VC connection.
    *
    * @param vcService  initialized VC connection
    */
   private static URI getPbmServiceLocation(VcService vcService) {

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
    * Waits for the specified task to change its state and returns the updated
    * state.
    */
   private static State waitForUpdatedTaskState(ManagedObjectReference taskMoRef,
         VcService service) throws InvalidProperty, InvalidCollectorVersion {
      Task task = service.getManagedObject(taskMoRef);
      State taskState = task.getInfo().getState();

      _logger.info(String.format(
            "Wait state update of task '%s' initial state is: %s",
            task.getInfo().getName(),
            taskState));

      // Get the property collector
      ManagedObjectReference propCollectorRef =
            service.getServiceInstanceContent().getPropertyCollector();
      PropertyCollector propCollector = service.getManagedObject(propCollectorRef);

      // Configure wait options to block the thread until task state is updated
      WaitOptions waitOption = new WaitOptions();
      waitOption.setMaxWaitSeconds(1);

      propCollector.createFilter(buildTaskStateFilterSpec(taskMoRef, service), true);

      long endTime =
            System.currentTimeMillis() + BackendDelay.LARGE.getDuration();

      String version = "";

      while (!isTaskComplete(taskState)) {
         UpdateSet updateSet = propCollector.waitForUpdatesEx(version, waitOption);

         if (updateSet != null) {
            version = updateSet.getVersion();

            State newState = getUpdatedTaskState(updateSet);
            if (newState != null) {
               taskState = newState;
               _logger.info(String.format(
                     "Task state changed to: '%s'",
                     taskState.name()));
            }
         }

         // Check whether waiting for the task o complete has timed out
         if (System.currentTimeMillis() > endTime) {
            _logger.error(String.format("Cancel waiting for task '%s' ", task.getInfo()
                  .getName()));
            break;
         }
      }

      return taskState;
   }

   /*
    * Constructs a FilterSpec object for retrieval of TaskInfo.State property
    * value with the property collector
    */
   private static FilterSpec buildTaskStateFilterSpec(ManagedObjectReference taskMoRef,
         VcService service) {
      FilterSpec filterSpec = new FilterSpec();

      // Create list view for the task
      ManagedObjectReference taskList[] = new ManagedObjectReference[] { taskMoRef };

      ManagedObjectReference viewMrMoRef = null;
      ViewManager viewManager = null;
      viewMrMoRef = service.getServiceInstanceContent().getViewManager();
      viewManager = service.getManagedObject(viewMrMoRef);

      ManagedObjectReference listViewMoRef = viewManager.createListView(taskList);

      // Create traversal spec to select the task in the view
      TraversalSpec traversalSpec = new TraversalSpec();
      traversalSpec.setName("traverseTasks");
      traversalSpec.setPath("view");
      traversalSpec.setSkip(false);
      traversalSpec.setType(new TypeNameImpl("ListView"));

      // Create object spec to start the traversal
      ObjectSpec objectSpec = new ObjectSpec();
      objectSpec.setObj(listViewMoRef);
      objectSpec.setSkip(true);
      objectSpec.selectSet = new SelectionSpec[] { traversalSpec };

      // Create porperty spec for Task.info.state
      PropertySpec propSpec = new PropertySpec();
      propSpec.setType(new TypeNameImpl("Task"));
      propSpec.setAll(false);
      propSpec.pathSet = new String[] { "info.state" };

      // Set the filter spec properties
      filterSpec.objectSet = new ObjectSpec[] { objectSpec };
      filterSpec.propSet = new PropertySpec[] { propSpec };

      return filterSpec;
   }

   /**
    * Searches in the updateSet a change corresponding to an update of
    * TaskInfo.State property value.
    */
   private static State getUpdatedTaskState(UpdateSet updateSet) {
      FilterUpdate filterUpdates[] = updateSet.getFilterSet();

      for (FilterUpdate filterUpdate : filterUpdates) {
         ObjectUpdate objectUpdates[] = filterUpdate.getObjectSet();

         for (ObjectUpdate objectUpdate : objectUpdates) {
            Change changes[] = objectUpdate.getChangeSet();

            for (Change change : changes) {
               if (change.getName().equals("info.state")) {
                  return (State) change.getVal();
               }
            }
         }
      }

      return null;
   }

   /**
    * Checks whether a task is completed. State "success" and "error" means that
    * the task is completed.
    */
   private static boolean isTaskComplete(State taskState) {
      return taskState.name().equals(State.success.name())
            || taskState.name().equals(State.error.name());
   }

}
