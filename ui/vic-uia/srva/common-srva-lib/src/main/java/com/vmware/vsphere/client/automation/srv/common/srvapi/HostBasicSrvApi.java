/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import com.vmware.vsphere.client.automation.srv.common.HostUtil;
import com.vmware.vsphere.client.automation.srv.common.spec.*;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.BackendDelay;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.ClusterComputeResource;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.ComputeResource.ConfigSpec;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.HostSystem.ConnectionState;
import com.vmware.vim.binding.vim.Network;
import com.vmware.vim.binding.vim.Task;
import com.vmware.vim.binding.vim.host.ConfigInfo;
import com.vmware.vim.binding.vim.host.ConnectSpec;
import com.vmware.vim.binding.vim.host.DatastoreSystem;
import com.vmware.vim.binding.vim.host.FirewallSystem;
import com.vmware.vim.binding.vim.host.HostAccessManager.LockdownMode;
import com.vmware.vim.binding.vim.host.HostBusAdapter;
import com.vmware.vim.binding.vim.host.InternetScsiHba;
import com.vmware.vim.binding.vim.host.InternetScsiHba.SendTarget;
import com.vmware.vim.binding.vim.host.InternetScsiHba.StaticTarget;
import com.vmware.vim.binding.vim.host.NetworkInfo;
import com.vmware.vim.binding.vim.host.StorageDeviceInfo;
import com.vmware.vim.binding.vim.host.StorageSystem;
import com.vmware.vim.binding.vim.host.VirtualNic;
import com.vmware.vim.binding.vim.host.VirtualNicManager;
import com.vmware.vim.binding.vim.host.VirtualNicManager.NetConfig;
import com.vmware.vim.binding.vim.host.VirtualNicManagerInfo;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * Class that implements the host operations through the API - add,
 * connect/disconnect, enter/exit maintenance mode, check existence, delete
 */
public class HostBasicSrvApi {

   private static final Logger _logger = LoggerFactory
         .getLogger(HostBasicSrvApi.class);

   private static final String ICSI_DRIVER_TYPE = "iscsi_vmk";
   private static final String FIREWALL_SERVICE_NTPCLIENT = "ntpClient";

   // VIM type of Clusters
   private static final String CLUSTER_COMPUTE_RESOURCE_VIM_TYPE = "ClusterComputeResource";

   // used as a boolean condition for the wait state method
   public enum HostState {
      CONNECTED, DISCONNECTED, ENTER_MAINTENANCE_MODE, EXIT_MAINTENANCE_MODE, NOT_RESPONDING
   }

   private static HostBasicSrvApi instance = null;

   protected HostBasicSrvApi() {
   }

   /**
    * Get instance of HostSrvApi.
    *
    * @return created instance
    */
   public static HostBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (HostBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing HostSrvApi.");
               instance = new HostBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a standalone host in a specified datacenter or a clustered host in
    * a specified cluster and waits the host to get connected
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be added, the
    *           parent should contain either the DatacenterSpec of the target
    *           datacenter for the standalone host, or the ClusterSpec of the
    *           target cluster of the clustered host. Note that the ClusterSpec
    *           should also contain its parent DatacenterSpec
    * @param asConnected
    *           Flag to specify whether or not the host should be connected as
    *           soon as it is added. The creation operation fails if a
    *           connection attempt is made and fails.
    * @return True if the creation was successful, false otherwise.
    * @throws Exception
    */
   public boolean addHost(HostSpec hostSpec, boolean asConnected)
         throws Exception {
      if (hostSpec.parent.get() instanceof DatacenterSpec) {
         return addStandaloneHost(hostSpec, null, asConnected, null,
               (int) BackendDelay.MEDIUM.getDuration() / 1000);
      }

      return addClusteredHost(hostSpec, asConnected, null, null,
            (int) BackendDelay.MEDIUM.getDuration() / 1000);
   }

   /**
    * Checks whether the specified host is in Maintenance Mode
    *
    * @param hostSpec
    *           <code>HostSpec</code> instance representing the host
    *
    * @return True if the host is in Maintenance Mode, otherwise false
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean isInMaintenanceMode(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);
      // Get the host
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      return host.getRuntime().isInMaintenanceMode();
   }

   /**
    * Enters a host into maintenance mode and waits for the state to settle
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be added, the
    *           parent should contain the ClusterSpec
    *
    * @return True if the host entered maintenance mode successfully, false
    *         otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   public boolean enterMaintenanceMode(HostSpec hostSpec) throws Exception {
      return enterMaintenanceMode(hostSpec, false,
            (int) BackendDelay.MEDIUM.getDuration() / 1000, 0);
   }

   /**
    * Exits a host from maintenance mode
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be added, the
    *           parent should contain the ClusterSpec
    *
    * @return True if the host exited maintenance mode successfully, false
    *         otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   public boolean exitMaintenanceMode(HostSpec hostSpec) throws Exception {
      return exitMaintenanceMode(hostSpec,
            (int) BackendDelay.SMALL.getDuration() / 1000, 0);
   }

   /**
    * Checks whether the specified host is connected
    *
    * @param hostSpec
    *           <code>HostSpec</code> instance representing the host
    *
    * @return True if the host is connected, otherwise false
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean isConnected(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);
      // Get the host
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      return host.getRuntime().getConnectionState()
            .equals(ConnectionState.connected) ? true : false;
   }

   /**
    * Create iSCSI adapter and add dynamic target to it.
    *
    * @param hostSpec
    *           the spec of the host to which iSCSI adapter will be added
    * @return true if the operation was successful
    * @throws Exception
    */
   public boolean createIscsiAdapter(HostSpec hostSpec) throws Exception {

      // validate if the host spec has iSCSI server IP
      if (!hostSpec.iscsiServerIp.isAssigned()) {
         _logger.error("Host spec does not have iSCSI server specified.");
         throw new IllegalArgumentException(
               "Host spec does not have iSCSI server specified.");
      }

      StorageSystem hostStorageSystem = getHostStorageSystem(hostSpec);

      _logger.info("About to create iSCSI adapter.");
      hostStorageSystem.updateSoftwareInternetScsiEnabled(true);

      // find the name of the iSCSI adapter
      String iscsiAdapterName = null;
      StorageDeviceInfo storageDeviceInfo = hostStorageSystem
            .getStorageDeviceInfo();
      for (HostBusAdapter adapter : storageDeviceInfo.getHostBusAdapter()) {
         if (adapter.getDriver().equals(ICSI_DRIVER_TYPE)) {
            iscsiAdapterName = adapter.getDevice();
         }
      }

      // check if the adapter was created successfully
      if (iscsiAdapterName == null) {
         _logger.error("The iSCSI adapter was not created!");
         return false;
      }

      SendTarget sendTarget = new SendTarget();
      sendTarget.setAddress(hostSpec.iscsiServerIp.get());
      SendTarget[] sendTargets = new SendTarget[] { sendTarget };

      // add dynamic target to iSCSI adapter
      _logger.info(String.format(
            "About to add dynamic target to iSCSI addapter - server IP: %s",
            hostSpec.iscsiServerIp.get()));
      hostStorageSystem.addInternetScsiSendTargets(iscsiAdapterName,
            sendTargets);
      hostStorageSystem.rescanAllHba();
      return true;
   }

   /**
    * Destroy iSCSI adapter and remove dynamic and static target from it.
    *
    * @param hostSpec
    *           the spec of the host to which iSCSI adapter will be removed
    * @return true if the operation was successful
    * @throws Exception
    */
   public boolean destroyIscsiAdapter(HostSpec hostSpec) throws Exception {

      // validate if the host spec has iSCSI server IP
      if (!hostSpec.iscsiServerIp.isAssigned()) {
         _logger.error("Host spec does not have iSCSI server specified.");
         throw new IllegalArgumentException(
               "Host spec does not have iSCSI server specified.");
      }

      StorageSystem hostStorageSystem = getHostStorageSystem(hostSpec);

      // find the name of the iSCSI adapter
      String iscsiAdapterName = null;
      List<StaticTarget> staticTargetsForDeleteion = new ArrayList<StaticTarget>();
      StorageDeviceInfo storageDeviceInfo = hostStorageSystem
            .getStorageDeviceInfo();

      for (HostBusAdapter adapter : storageDeviceInfo.getHostBusAdapter()) {

         if (adapter.getDriver().equals(ICSI_DRIVER_TYPE)
               && adapter instanceof InternetScsiHba) {
            iscsiAdapterName = adapter.getDevice();

            // Find the static target for deletion
            InternetScsiHba iscsiHba = (InternetScsiHba) adapter;
            for (StaticTarget staticTarget : iscsiHba.configuredStaticTarget) {
               if (staticTarget.getAddress().equals(
                     hostSpec.iscsiServerIp.get())) {
                  staticTargetsForDeleteion.add(staticTarget);
               }
            }
         }
      }

      // check if the adapter exists
      if (iscsiAdapterName == null) {
         _logger.error("The iSCSI adapter was not found!");
         throw new RuntimeException("The iSCSI adapter was not found!");
      }

      // remove both dynamic and static target to iSCSI adapter
      _logger.info(String.format(
            "About to remove dynamic target to iSCSI addapter - server IP: %s",
            hostSpec.iscsiServerIp.get()));

      SendTarget sendTarget = new SendTarget();
      sendTarget.setAddress(hostSpec.iscsiServerIp.get());
      SendTarget[] sendTargets = new SendTarget[] { sendTarget };

      hostStorageSystem.removeInternetScsiSendTargets(iscsiAdapterName,
            sendTargets);

      // We have also to remove the static target as it is added automatically
      StaticTarget[] staticTargets = new StaticTarget[staticTargetsForDeleteion
            .size()];
      staticTargetsForDeleteion.toArray(staticTargets);

      hostStorageSystem.removeInternetScsiStaticTargets(iscsiAdapterName,
            staticTargets);
      hostStorageSystem.rescanAllHba();

      return true;
   }

   StorageSystem getHostStorageSystem(HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

      StorageSystem storageSystem = ManagedEntityUtil
            .getManagedObjectFromMoRef(host.getConfigManager()
                  .getStorageSystem(), hostSpec.service.get());
      return storageSystem;
   }

   /**
    * Connect a disconnected host
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be connected
    *
    * @return True if the host was connected successfully, false otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   public boolean connectHost(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ConnectSpec connectSpec = buildConnectSpec(hostSpec);
      ManagedObjectReference task = host.reconnect(connectSpec, null);

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(task, hostSpec)
            && waitForHostState(HostState.CONNECTED, host,
                  (int) BackendDelay.MEDIUM.getDuration() / 1000);
   }

   /**
    * Disconnect a connected host
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be
    *           disconnected
    *
    * @return True if the host was disconnected successfully, false otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   public boolean disconnectHost(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ManagedObjectReference task = host.disconnect();

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(task, hostSpec)
            && waitForHostState(HostState.DISCONNECTED, host,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified host exists.
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be checked
    *           for existence
    *
    * @return True if the host exists, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean checkHostExists(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);

      _logger.info(String.format("Checking whether host '%s' exists",
            hostSpec.name.get()));

      try {
         ManagedEntityUtil.getManagedObject(hostSpec);
      } catch (ObjectNotFoundException e) {
         return false;
      }
      return true;
   }

   /**
    * Deletes the specified host from the inventory and wait the specified
    * timeout to let the entity get really deleted, returns false if the timeout
    * is over or the entity cannot be deleted
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be deleted
    * @return True if the deletion was successful, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails
    */
   public boolean deleteHost(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);

      _logger.info(String.format("Deleting host '%s", hostSpec.name.get()));

      // Get the host
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ManagedObjectReference task = null;
      if (host.getParent().getType().equals(CLUSTER_COMPUTE_RESOURCE_VIM_TYPE)) {
         // clustered host has to be disconnected before removing it from the VC
         if (isConnected(hostSpec)) {
            disconnectHost(hostSpec);
         }
         task = host.destroy();
      } else {
         ComputeResource parentComputeResource = VcServiceUtil.getVcService(
               hostSpec).getManagedObject(host.getParent());
         task = parentComputeResource.destroy();
      }

      // success or failure of the task
      return VcServiceUtil.waitForTaskSuccess(task, hostSpec)
            && ManagedEntityUtil.waitForEntityDeletion(hostSpec,
                  (int) BackendDelay.SMALL.getDuration() / 1000);
   }

   /**
    * Checks whether the specified host exists and deletes it.
    *
    * @param hostSpec
    *           <code>HostSpec</code> instance representing the host to be
    *           deleted.
    *
    * @return True if the host doesn't exist or if the host is deleted
    *         successfully, false otherwise.
    *
    * @throws Exception
    *            If login to the VC service fails, or the specified folder
    *            doesn't exist.
    */
   public boolean deleteHostSafely(HostSpec hostSpec) throws Exception {

      if (checkHostExists(hostSpec)) {
         return deleteHost(hostSpec);
      }

      // Positive taskSucceeded if the host doesn't exist
      return true;
   }

   /**
    * Retrieves list of host datastore names.
    *
    * @param hostSpec
    *           the spec of the host that will be searched
    * @deprecated use {@link HostBasicSrvApi#getDatastores(HostSpec)}
    */
   @Deprecated
   // TODO: MOVE to the place where such custom logic is needed. This method
   // does projection from Datastore to String and has nothing to do with Host
   // API
   public List<String> getHostDatastores(HostSpec hostSpec) {
      List<String> result = new ArrayList<String>();

      for (Datastore datastore : getDatastores(hostSpec)) {
         result.add(datastore.getName());
      }

      return result;
   }

   /**
    * Retrieve the datastores for a given host
    *
    * @param hostSpec
    *           the spec for the host
    * @return list of {@link Datastore} associated with the host
    */
   public List<Datastore> getDatastores(HostSpec hostSpec) {
      VcService vcService = VcServiceUtil.getVcService(hostSpec);

      List<Datastore> result = new ArrayList<Datastore>();
      for (ManagedObjectReference datastoreRef : this.getHostSystem(hostSpec)
            .getDatastore()) {
         Datastore datastore = vcService
               .getManagedObject(datastoreRef);
         result.add(datastore);
      }

      return result;
   }

   /**
    * Retrieve the datastore system for a host
    *
    * @param hostSpec
    * @return {@link DatastoreSystem} associated with the host
    */
   public DatastoreSystem getDatastoreSystem(HostSpec hostSpec) {
      HostSystem host = getHostSystem(hostSpec);

      VcService service = VcServiceUtil.getVcService(hostSpec);
      return service.getManagedObject(host.getConfigManager()
            .getDatastoreSystem());
   }

   /**
    * Retrieves list of host datastore specs.
    *
    * @param hostSpec
    *           the spec of the host that will be searched
    * @throws RuntimeException
    *            if login to the VC server fails or host is not present
    * @deprecated use {@link HostBasicSrvApi#getDatastores(HostSpec)}
    */
   @Deprecated
   public List<DatastoreSpec> getHostDatastoreSpecs(HostSpec hostSpec) {
      HostSystem host = null;
      VcService service = null;
      try {
         service = VcServiceUtil.getVcService(hostSpec);
         host = ManagedEntityUtil.getManagedObject(hostSpec);
      } catch (Exception e) {
         String errorMessage = String.format("Cannot retrieve host %s",
               hostSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }

      List<DatastoreSpec> datastoreSpecs = new ArrayList<DatastoreSpec>();
      for (ManagedObjectReference datastoreRef : host.getDatastore()) {
         Datastore datastore = service.getManagedObject(datastoreRef);
         DatastoreSpec tempDatastoreSpec = SpecFactory.getSpec(
               DatastoreSpec.class, datastore.getName(), hostSpec);
         tempDatastoreSpec.type.set(DatastoreType.valueOf(datastore
               .getSummary().getType()));
         tempDatastoreSpec.service.set(hostSpec.service.get());
         datastoreSpecs.add(tempDatastoreSpec);
      }

      return datastoreSpecs;
   }

   /**
    * Retrieves list of host network names.
    *
    * @param hostSpec
    *           the spec for the host that will be searched
    * @throws RuntimeException
    *            if login to the VC server fails or host is not present
    */
   public List<String> getNetworkNames(HostSpec hostSpec) {
      List<String> networkNames = new ArrayList<String>();
      HostSystem host = null;
      VcService service = null;
      ServiceSpec serviceSpec = hostSpec.service.get();

      try {
         service = VcServiceUtil.getVcService(serviceSpec);
         host = ManagedEntityUtil.getManagedObject(hostSpec, serviceSpec);
      } catch (Exception e) {
         String errorMessage = String.format("Cannot retrieve host %s",
               hostSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }

      for (ManagedObjectReference networkMor : host.getNetwork()) {
         Network network = service.getManagedObject(networkMor);
         networkNames.add(network.getName());
      }

      return networkNames;
   }

   /**
    * Method that enables vmotion if it is not set
    *
    * @param hostSpec
    *           - specification of the host on which to enable vmotion
    * @return true if successfully set or already set, false if not set and no
    *         available vnics found
    * @throws Exception
    */
   public boolean enableVmotion(HostSpec hostSpec) throws Exception {
      return setVmotion(true, getHostSystem(hostSpec),
            VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that disables vmotion if it is not set
    *
    * @param hostSpec
    *           - specification of the host on which to enable vmotion
    * @return true if successfully disabled or already disabled, false otherwise
    * @throws Exception
    */
   public boolean disableVmotion(HostSpec hostSpec) throws Exception {
      return setVmotion(false, getHostSystem(hostSpec),
            VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that checks the vMotion state on a host
    *
    * @param expectedState
    *           - the state we expect
    * @param hostSpec
    *           - specification of the host on which to check vMotion
    * @return true if the state is expected, false otherwise
    * @throws Exception
    */
   public boolean checkVmotionState(boolean expectedState, HostSpec hostSpec)
         throws Exception {
      return expectedState == getVmotionState(getHostSystem(hostSpec),
            VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that enables fault tolerance logging if it is not set
    *
    * @param hostSpec
    *           - specification of the host on which to enable fault tolerance
    *           logging
    * @return true if successfully set or already set, false if not set and no
    *         available vnics found
    * @throws Exception
    */
   public boolean enableFtLogging(HostSpec hostSpec) throws Exception {
      return setFtLogging(true, getHostSystem(hostSpec),
            VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that disables fault tolerance logging if it is not set
    *
    * @param hostSpec
    *           - specification of the host on which to enable fault tolerance
    *           logging
    * @return true if successfully disabled or already disabled, false otherwise
    * @throws Exception
    */
   public boolean disableFtLogging(HostSpec hostSpec) throws Exception {
      return setFtLogging(false, getHostSystem(hostSpec),
            VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that checks fault tolerance logging state on a host
    *
    * @param expectedState
    *           - the state we expect
    * @param hostSpec
    *           - specification of the host on which to check fault tolerance
    *           logging
    * @return true if the state is expected, false otherwise
    * @throws Exception
    */
   public boolean checkFtLoggingState(boolean expectedState, HostSpec hostSpec)
         throws Exception {
      return expectedState == getFtLoggingState(getHostSystem(hostSpec),
            VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that opens the NTP Client ports in the firewall
    *
    * @param hostSpec
    *           - specification of the host to be modified
    * @throws Exception
    *            if host spec is invalid, the host doesn't exist or there is a
    *            problem with ComplianceCheckErrorMessageMakeCompliant
    */
   public void enableNtpClient(final HostSpec hostSpec) throws Exception {
      setFirewallRuleSet(hostSpec, FIREWALL_SERVICE_NTPCLIENT, true);
   }

   /**
    * Method that closes the NTP Client ports in the firewall
    *
    * @param hostSpec
    *           - specification of the host to be modified
    * @throws Exception
    *            if host spec is invalid, the host doesn't exist or there is a
    *            problem with ComplianceCheckErrorMessageMakeCompliant
    */
   public void disableNtpClient(final HostSpec hostSpec) throws Exception {
      setFirewallRuleSet(hostSpec, FIREWALL_SERVICE_NTPCLIENT, false);
   }

   public String getHostVersion(final HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      return host.getConfig().getProduct().getVersion();
   }

   /**
    * Method that waits for a host to enter specified state or for the timeout to pass
    * @param hostSpec - host
    * @param hostState - state and timeout
    * @return true if successful, false otherwise
    * @throws Exception
     */
   public boolean waitForHostToEnterState(HostSpec hostSpec,
                                    HostStateSpec hostState) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

      // wait for timeout to get host in state
      int retries = hostState.numOfRetries.get();
      HostState state = getHostState(hostState.state.get());
      while (!isHostInState(host, state) && retries > 0) {
         Thread.sleep(1000);
         retries--;
      }
      // if timeout hasn't passed return true
      if (retries > 0) {
         return true;
      }

      _logger.error("Timeout has passed without reaching the expected state: "
              + state.name());
      return false;
   }

   /**
    * Creates a standalone host in a datacenter
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be added, the
    *           parent should contain the DatacenterSpec
    * @param compResSpec
    *           Optionally specify the configuration for the compute resource
    *           that will be created to contain the host.
    * @param asConnected
    *           Flag to specify whether or not the host should be connected as
    *           soon as it is added. The creation operation fails if a
    *           connection attempt is made and fails.
    * @param license
    *           Provide a licenseKey or licenseKeyType
    * @param timeout
    *           - timeout in seconds, i.e. 60 is 60 seconds, to wait for the
    *           specified state to be reached, needed only if host is being
    *           added as connected
    * @return True if the creation was successful, false otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   protected boolean addStandaloneHost(HostSpec hostSpec, ConfigSpec compResSpec,
         boolean asConnected, String license, int timeout) throws Exception {
      validateHostSpec(hostSpec);

      _logger.info(String.format("Adding host '%s' to datacenter '%s'",
            hostSpec.name, hostSpec.parent.get().name.get()));

      Folder folder = getHostFolder((DatacenterSpec) hostSpec.parent.get());
      ConnectSpec connectSpec = buildConnectSpec(hostSpec);
      // Set ssl thumbprint
      String sslThumbprint = DatacenterBasicSrvApi.getInstance()
            .validateSslThumbprint(hostSpec,
                  (DatacenterSpec) hostSpec.parent.get());
      if (sslThumbprint != null) {
         connectSpec.setSslThumbprint(sslThumbprint);
      }
      // Add host
      ManagedObjectReference task = folder.addStandaloneHost(connectSpec,
            compResSpec, asConnected, license);
      // success or failure of the task
      boolean taskSucceeded = VcServiceUtil.waitForTaskSuccess(task, hostSpec);

      if (taskSucceeded && asConnected) {
         VcService service = VcServiceUtil.getVcService(hostSpec);
         // get the added host system
         Task addHostTask = service.getManagedObject(task);
         ManagedObjectReference crMor = (ManagedObjectReference) addHostTask
               .getInfo().getResult();

         // Workaround for waitForTaskSuccess() returning true prematurely,
         // before the Add host task is completed - PR 1423973
         int retry = 0;
         while ((crMor == null) && (retry < timeout)) {
            crMor = (ManagedObjectReference) addHostTask.getInfo().getResult();
            Thread.sleep(1000);
            retry++;
         }

         ComputeResource cr = service.getManagedObject(crMor);
         ManagedObjectReference hostMor = cr.getHost()[0];
         HostSystem host = service.getManagedObject(hostMor);
         // if host is added as connected wait with 1 min timeout to have it
         // really connected as task returns success earlier
         boolean isConnected = waitForHostState(HostState.CONNECTED, host,
               timeout);
         // exit host of Maintenance Mode if it is in
         if (isConnected && host.getRuntime().isInMaintenanceMode()) {
            return exitMaintenanceMode(hostSpec);
         }
         return isConnected;
      }
      return taskSucceeded;
   }

   /**
    * Retrieves the host folder of specified datacenter
    *
    * @param datacenterSpec
    *           - specification of the datacenter
    * @throws Exception
    *            If login to the VC service fails
    */
   protected Folder getHostFolder(DatacenterSpec datacenterSpec) throws Exception {
      return FolderBasicSrvApi.getInstance().getHostFolder(datacenterSpec);
   }

   /**
    * Creates a host into cluster
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be added, the
    *           parent should contain the ClusterSpec. Note that the ClusterSpec
    *           should contain the DatacenterSpec of its parent datacenter
    * @param asConnected
    *           - whether the host to be connected
    * @param resPool
    *           - whether to create a new child resource pool in the cluster
    * @param license
    *           Provide a licenseKey or licenseKeyType
    * @param timeout
    *           - timeout in seconds, i.e. 60 is 60 seconds, to wait for the
    *           specified state to be reached, needed only if host is being
    *           added as connected
    * @return True if the creation was successful, false otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   protected boolean addClusteredHost(HostSpec hostSpec, boolean asConnected,
         ManagedObjectReference resPool, String license, int timeout)
         throws Exception {
      validateHostSpec(hostSpec);

      _logger.info(String.format("Adding host '%s' to cluster '%s'",
            hostSpec.name.get(), hostSpec.parent.get().name.get()));

      ClusterSpec clusterSpec = (ClusterSpec) hostSpec.parent.get();

      ClusterComputeResource cluster = ManagedEntityUtil
            .getManagedObject(clusterSpec);
      ConnectSpec connectSpec = buildConnectSpec(hostSpec);
      // Set ssl thumbprint
      DatacenterSpec datacenterSpec = (DatacenterSpec) clusterSpec.parent.get();
      String sslThumbprint = DatacenterBasicSrvApi.getInstance()
            .validateSslThumbprint(hostSpec, datacenterSpec);
      if (sslThumbprint != null) {
         connectSpec.setSslThumbprint(sslThumbprint);
      }
      // add host
      ManagedObjectReference task = cluster.addHost(connectSpec, asConnected,
            resPool, license);
      // success or failure of the task
      boolean taskSucceeded = VcServiceUtil.waitForTaskSuccess(task, hostSpec);

      if (taskSucceeded && asConnected) {
         // get added host
         final Task addHostTask = VcServiceUtil.getVcService(hostSpec)
               .getManagedObject(task);

         // wait for the task to finish
         final CountDownLatch finished = new CountDownLatch(1);
         new Thread(new Runnable() {
            @Override
            public void run() {
               while (addHostTask.getInfo().getResult() == null) {
               }
               finished.countDown();
            }
         }).start();

         boolean hostAddedSuccessfully = finished.await(
               BackendDelay.SMALL.getDuration(), TimeUnit.MILLISECONDS);
         if (!hostAddedSuccessfully) {
            throw new RuntimeException(
                  "Adding clustered host did not finnish on time");
         }

         ManagedObjectReference hostMor = (ManagedObjectReference) addHostTask
               .getInfo().getResult();

         HostSystem host = VcServiceUtil.getVcService(hostSpec)
               .getManagedObject(hostMor);
         // if host is added as connected wait with 1 min timeout to have it
         // really connected as task returns success earlier
         boolean isConnected = waitForHostState(HostState.CONNECTED, host,
               timeout);
         // exit host of Maintenance Mode if it is in
         if (isConnected && host.getRuntime().isInMaintenanceMode()) {
            return exitMaintenanceMode(hostSpec);
         }
         return isConnected;
      }
      return taskSucceeded;
   }

   /**
    * Enters a host into maintenance mode
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be
    *           manipulated, the parent should contain the ClusterSpec
    * @param evacuatePoweredOffVms
    *           This is a parameter only supported by VirtualCenter. If set to
    *           true, for a DRS disabled cluster, the task will not succeed
    *           unless all powered-off virtual machines have been manually
    *           reregistered; for a DRS enabled cluster, VirtualCenter will
    *           automatically reregister powered-off virtual machines and a
    *           powered-off virtual machine may remain at the host only for two
    *           reasons: (a) no compatible host found for reregistration, (b)
    *           DRS is disabled for the virtual machine. If set to false,
    *           powered-off virtual machines do not need to be moved.
    * @param stateTimeout
    *           Timeout in seconds to wait for the task to complete
    * @param vcTaskTimeout
    *           The task completes when the host successfully enters maintenance
    *           mode or the timeout expires, and in the latter case the task
    *           contains a Timeout fault. If the timeout is less than or equal
    *           to zero, there is no timeout. The timeout is specified in
    *           seconds.
    * @return True if the entered maintenance mode successfully, false
    *         otherwise.
    * @throws Exception
    *            if login to vc service fails
    */
   protected boolean enterMaintenanceMode(HostSpec hostSpec,
         boolean evacuatePoweredOffVms, int stateTimeout, int vcTaskTimeout)
         throws Exception {
      validateHostSpec(hostSpec);
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ManagedObjectReference task = host.enterMaintenanceMode(vcTaskTimeout,
            evacuatePoweredOffVms, null);
      return VcServiceUtil.waitForTaskSuccess(task, hostSpec)
            && waitForHostState(HostState.ENTER_MAINTENANCE_MODE, host,
                  stateTimeout);
   }

   /**
    * Exits a host from maintenance mode and waits for the host to exit
    * maintenance mode
    *
    * @param hostSpec
    *           The host name, port, and passwords for the host to be
    *           manipulated, the parent should contain the ClusterSpec
    * @param vcTaskTimeout
    *           The task completes when the host successfully exits maintenance
    *           mode or the timeout expires, and in the latter case the task
    *           contains a Timeout fault. If the timeout is less than or equal
    *           to zero, there is no timeout. The timeout is specified in
    *           seconds.
    * @param stateTimeout
    *           Timeout in seconds to wait for the task to complete
    *
    * @return True if the creation was successful, false otherwise.
    * @throws Exception
    */
   protected boolean exitMaintenanceMode(HostSpec hostSpec, int stateTimeout,
         int vcTaskTimeout) throws Exception {
      validateHostSpec(hostSpec);
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ManagedObjectReference task = host.exitMaintenanceMode(vcTaskTimeout);
      return VcServiceUtil.waitForTaskSuccess(task, hostSpec)
            && waitForHostState(HostState.EXIT_MAINTENANCE_MODE, host,
                  stateTimeout);
   }

   /**
    * Checks whether the host is in the specified state - connected,
    * disconnected, maintenance mode (in/out)
    *
    * @param host
    *           HostSystem object to check - using that type as this is private
    *           helper method and this is faster
    * @param state
    *           State that we want to check for
    * @return True if it is in the specified state, otherwise false
    * @throws IllegalArgumentException
    *            if no such state
    */
   protected boolean isHostInState(HostSystem host, HostState state)
         throws IllegalArgumentException {
      switch (state) {
      case CONNECTED:
         return host.getRuntime().getConnectionState()
               .equals(ConnectionState.connected);
      case DISCONNECTED:
         return host.getRuntime().getConnectionState()
               .equals(ConnectionState.disconnected);
      case ENTER_MAINTENANCE_MODE:
         return host.getRuntime().isInMaintenanceMode();
      case EXIT_MAINTENANCE_MODE:
         return !host.getRuntime().isInMaintenanceMode();
      case NOT_RESPONDING:
         return host.getRuntime().getConnectionState()
                .equals(ConnectionState.notResponding);
      default:
         throw new IllegalArgumentException("Incorrect host state.");
      }
   }

   /**
    * Get the config info for a host
    *
    * @param host
    *           HostSystem object to check
    * @return ConfigInfo the config info for selected host
    * @throws IllegalArgumentException
    *            if the host is null
    * @throws Exception
    *            if the host config info is not available
    */
   private ConfigInfo getHostConfigInfo(HostSystem host) throws Exception {
      if (host == null) {
         throw new IllegalArgumentException("Applied host is invalid!");
      }

      ConfigInfo configInfo = host.getConfig();
      if (configInfo == null) {
         throw new Exception(
               "Config Info is not available for host; host might be disocnnected!");
      }

      return configInfo;
   }

   /**
    * Method that enables or disable vmotion. It is used also in
    * HostCLientSrvApi
    *
    * @param enable
    *           - true to enable vmotion, false to disable it
    * @param host
    *           - host on which to perform action
    * @param service
    *           - service spec
    * @return true if successful, false otherwise
    * @throws Exception
    *            - if vc connection fails or there is problem with setting
    *            network parameters
    */
   // method is default as it is also used in HostClientSrv class
   protected boolean setVmotion(boolean enable, HostSystem host, VcService service)
         throws Exception {
      ConfigInfo configInfo = getHostConfigInfo(host);
      VirtualNicManagerInfo vnicInfo = configInfo.getVirtualNicManagerInfo();
      ManagedObjectReference vnicManagerMoRef = host.getConfigManager()
            .getVirtualNicManager();
      VirtualNicManager vnicManager = service
            .getManagedObject(vnicManagerMoRef);

      if (enable) {
         return enableVmotion(vnicManager, vnicInfo);
      } else {
         return disableVmotion(host, vnicManager, vnicInfo);
      }
   }

   /**
    * Method that enables or disable fault tolerance logging.
    *
    * @param enable
    *           - true to enable fault tolerance logging, false to disable it
    * @param host
    *           - host on which to perform action
    * @param service
    *           - service spec
    * @return true if successful, false otherwise
    * @throws Exception
    *            - if vc connection fails or there is problem with setting
    *            network parameters
    */
   protected boolean setFtLogging(boolean enable, HostSystem host,
         VcService service) throws Exception {
      ConfigInfo configInfo = getHostConfigInfo(host);
      VirtualNicManagerInfo vnicInfo = configInfo.getVirtualNicManagerInfo();
      ManagedObjectReference vnicManagerMoRef = host.getConfigManager()
            .getVirtualNicManager();
      VirtualNicManager vnicManager = service
            .getManagedObject(vnicManagerMoRef);

      if (enable) {
         return enableFtLogging(vnicManager, vnicInfo);
      } else {
         return disableFtLogging(host, vnicManager, vnicInfo);
      }
   }

   /**
    * Method that gets the fault tolerance logging state of a vnic on a host.
    *
    * @param host
    *           - host on which to perform action
    * @param service
    *           - service spec
    * @return true if fault tolerance logging is enabled, false otherwise
    * @throws Exception
    *            - if vc connection fails or there is problem with setting
    *            network parameters
    */
   protected boolean getFtLoggingState(HostSystem host, VcService service)
         throws Exception {
      ConfigInfo configInfo = getHostConfigInfo(host);
      VirtualNicManagerInfo vnicInfo = configInfo.getVirtualNicManagerInfo();

      NetConfig netConfig = getFtLoggingNic(vnicInfo);
      if (netConfig != null) {
         if (netConfig.getSelectedVnic() != null) {
            // fault tolerance logging is enabled
            return true;
         }
      }
      return false;
   }

   /**
    * Method that gets the vMotion state of a vnic on a host.
    *
    * @param host
    *           - host on which to perform action
    * @param service
    *           - service spec
    * @return true if vMotion is enabled, false otherwise
    * @throws Exception
    *            - if vc connection fails or there is problem with setting
    *            network parameters
    */
   protected boolean getVmotionState(HostSystem host, VcService service)
         throws Exception {
      ConfigInfo configInfo = getHostConfigInfo(host);
      VirtualNicManagerInfo vnicInfo = configInfo.getVirtualNicManagerInfo();

      NetConfig netConfig = getVmotionNic(vnicInfo);
      if (netConfig != null) {
         if (netConfig.getSelectedVnic() != null) {
            // vMotion is enabled
            return true;
         }
      }
      return false;
   }

   /**
    * Method that enables vmotion if it is not set; it is enabled for the first
    * available vnic
    *
    * @param vnicManager
    *           - host's virtual nic manager
    * @param vnicInfo
    *           - host's virtual nic manager info
    * @return true if successfully set or already set, false if not set and no
    *         available vnics found
    * @throws Exception
    */
   // Method is default as it is used in HostClientSrvApi also
   protected boolean enableVmotion(VirtualNicManager vnicManager,
         VirtualNicManagerInfo vnicInfo) throws Exception {
      NetConfig netConfig = getVmotionNic(vnicInfo);
      if (netConfig == null) {
         return false;
      }

      boolean hasSelectedNics = netConfig.getSelectedVnic() != null;
      if (hasSelectedNics) {
         // vmotion is already set
         return true;
      } else {
         // set vmotion
         VirtualNic[] vnicCandidates = netConfig.getCandidateVnic();
         if (vnicCandidates != null && vnicCandidates.length > 0) {
            // by default vmotion is enabled only for the first available
            // vnic
            setVnicForVmotion(true, vnicCandidates[0].getDevice(), vnicManager);
            return true;
         }
      }
      return false;
   }

   /**
    * Method that disables vmotion if it is not set; by default it gets disabled
    * for all vnics with enabled vmotion
    *
    * @param vnicManager
    *           - host's virtual nic manager
    * @param vnicInfo
    *           - host's virtual nic manager info
    * @return true if successfully disabled or already disabled, false otherwise
    * @throws Exception
    */
   // default as it is also used in HostClientSrvApi
   protected boolean disableVmotion(HostSystem host,
         VirtualNicManager vnicManager, VirtualNicManagerInfo vnicInfo)
         throws Exception {
      NetConfig netConfig = getVmotionNic(vnicInfo);
      if (netConfig == null) {
         return false;
      }
      boolean isVmotionNic = netConfig.getNicType().equals(
            VirtualNicManager.NicType.vmotion.name());
      if (isVmotionNic) {
         String[] vnicSelected = netConfig.getSelectedVnic();
         if (vnicSelected != null && vnicSelected.length > 0) {
            // vmotion is disabled for all vnics with enabled vmotion
            for (String vnic : vnicSelected) {
               setVnicForVmotion(false, getVnicDeviceByKey(host, vnic),
                     vnicManager);
            }
            return true;
         }
      }
      return true;
   }

   /**
    * Method that enables fault tolerance logging if it is not set; it is
    * enabled for the first available vnic
    *
    * @param vnicManager
    *           - host's virtual nic manager
    * @param vnicInfo
    *           - host's virtual nic manager info
    * @return true if successfully set or already set, false if not set and no
    *         available vnics found
    * @throws Exception
    */
   protected boolean enableFtLogging(VirtualNicManager vnicManager,
         VirtualNicManagerInfo vnicInfo) throws Exception {
      NetConfig netConfig = getFtLoggingNic(vnicInfo);
      if (netConfig == null) {
         return false;
      }

      boolean hasSelectedNics = netConfig.getSelectedVnic() != null;
      if (hasSelectedNics) {
         // fault tolerance logging is already set
         return true;
      } else {
         // set fault tolerance logging
         VirtualNic[] vnicCandidates = netConfig.getCandidateVnic();
         if (vnicCandidates != null && vnicCandidates.length > 0) {
            // by default ft logging is enabled only for the first available
            // vnic
            setVnicForFtLogging(true, vnicCandidates[0].getDevice(),
                  vnicManager);
            return true;
         }
      }
      return false;
   }

   /**
    * Method that disables fault tolerance logging if it is not set; by default
    * it gets disabled for all vnics with enabled fault tolerance logging
    *
    * @param vnicManager
    *           - host's virtual nic manager
    * @param vnicInfo
    *           - host's virtual nic manager info
    * @return true if successfully disabled or already disabled, false otherwise
    * @throws Exception
    */
   protected boolean disableFtLogging(HostSystem host,
         VirtualNicManager vnicManager, VirtualNicManagerInfo vnicInfo)
         throws Exception {
      NetConfig netConfig = getFtLoggingNic(vnicInfo);
      if (netConfig == null) {
         return false;
      }

      boolean isFtNic = netConfig.getNicType().equals(
            VirtualNicManager.NicType.faultToleranceLogging.name());
      if (isFtNic) {
         String[] vnicSelected = netConfig.getSelectedVnic();
         if (vnicSelected != null && vnicSelected.length > 0) {
            // ft logging is disabled for all vnics with enabled ft logging
            for (String vnic : vnicSelected) {
               setVnicForFtLogging(false, getVnicDeviceByKey(host, vnic),
                     vnicManager);
            }
            return true;
         }
      }
      return true;
   }

   /**
    * Builds the ConnectSpec object based on the HostSpec that is used for host
    * operations in the api
    *
    * @param hostSpec
    *           host specification of the host
    * @return ConnectSpec object based on info in HostSpec
    */
   protected ConnectSpec buildConnectSpec(HostSpec hostSpec) {
      return new ConnectSpec(hostSpec.name.get(), hostSpec.port.get(), null,
            hostSpec.userName.get(), hostSpec.password.get(), null, true, null,
            null, null, LockdownMode.lockdownDisabled, null);
   }

   /**
    * Retrieves the HostSystem object from inventory by HostSpec
    *
    * @param hostSpec
    *           - spec of host system object
    * @return - host system that corresponds to provided spec and is found in
    *         inventory
    */
   public HostSystem getHostSystem(HostSpec hostSpec) {
      validateHostSpec(hostSpec);
      ServiceSpec serviceSpec = hostSpec.service.get();

      try {
         return ManagedEntityUtil.getManagedObject(hostSpec, serviceSpec);
      } catch (Exception e) {
         String errorMessage = String.format(
               "Cannot retrieve host system for %s", hostSpec.name.get());
         _logger.error(errorMessage);
         throw new RuntimeException(errorMessage, e);
      }
   }

   /**
    * Method that gets the device of a virtual nic by the supplied virtual nic
    * key
    *
    * @param host
    *           - host of the virtual nic
    * @param vnicKey
    *           - key of the virtual nic
    * @return - the value of the device
    * @throws Exception
    *            - if invalid parameters are supplied, vnic doesn't belong to
    *            host, host is disconnected or host doesn't have a network
    */
   private String getVnicDeviceByKey(HostSystem host, String vnicKey)
         throws Exception {
      if (host == null || vnicKey == null) {
         throw new IllegalArgumentException(
               "Please, supply valid host and vnic key!");
      }

      ConfigInfo configInfo = host.getConfig();

      if (configInfo == null) {
         throw new Exception(
               "No configuration info for host; host migth be disconnected!");
      }

      NetworkInfo nwInfo = configInfo.getNetwork();

      if (nwInfo == null) {
         throw new Exception("No network information for host!");
      }

      VirtualNic[] vnics = nwInfo.getVnic();

      for (VirtualNic vnic : vnics) {
         if (vnicKey.contains(vnic.getKey())) {
            return vnic.getDevice();
         }
      }

      throw new IllegalArgumentException(
            "The supplied vnic key is not found in the supplied host!");
   }

   /**
    * Method that gets the network configuration for vmotion
    *
    * @param vnicInfo
    *           - host's VirtualNicManagerInfo that gets vmotion configuration
    * @return Network configuration for vmotion or null
    * @throws Exception
    */
   private NetConfig getVmotionNic(VirtualNicManagerInfo vnicInfo)
         throws Exception {
      if (vnicInfo == null) {
         throw new IllegalArgumentException("Invalid Virtual Nic Info!");
      }

      for (NetConfig netConfig : vnicInfo.getNetConfig()) {
         boolean isVmotionNic = netConfig.getNicType().equals(
               VirtualNicManager.NicType.vmotion.name());
         if (isVmotionNic) {
            return netConfig;
         }
      }

      return null;
   }

   /**
    * Method that either sets or unsets vnic for vmotion
    *
    * @param enable
    *           - true to set vnic for vmotion and false otherwise
    * @param vnicDevice
    *           - The device that uniquely identifies the VirtualNic
    * @param vnicManager
    *           - host's VirtualNicManager manager that will set nic
    * @throws Exception
    *            if cannot connect to vCenter server or host is not found in
    *            inventory
    */
   private void setVnicForVmotion(boolean enable, String vnicDevice,
         VirtualNicManager vnicManager) throws Exception {
      if (vnicManager == null) {
         throw new IllegalArgumentException("Invlaid Virtual Nic Manager!");
      }

      if (enable) {
         vnicManager.selectVnic(VirtualNicManager.NicType.vmotion.name(),
               vnicDevice);
      } else {
         vnicManager.deselectVnic(VirtualNicManager.NicType.vmotion.name(),
               vnicDevice);
      }
   }

   /**
    * Method that gets the network configuration for fault tolerance logging
    *
    * @param vnicInfo
    *           - host's VirtualNicManagerInfo that gets fault tolerance logging
    *           configuration
    * @return Network configuration for fault tolerance logging or null
    * @throws Exception
    */
   private NetConfig getFtLoggingNic(VirtualNicManagerInfo vnicInfo)
         throws Exception {
      if (vnicInfo == null) {
         throw new IllegalArgumentException("Invalid Virtual Nic Info!");
      }

      for (NetConfig netConfig : vnicInfo.getNetConfig()) {
         boolean isFtNic = netConfig.getNicType().equals(
               VirtualNicManager.NicType.faultToleranceLogging.name());
         if (isFtNic) {
            return netConfig;
         }
      }
      return null;
   }

   /**
    * Method that either sets or unsets vnic for fault tolerance logging
    *
    * @param enable
    *           - true to set vnic for fault tolerance logging and false
    *           otherwise
    * @param vnicDevice
    *           - The device that uniquely identifies the VirtualNic
    * @param vnicManager
    *           - host's VirtualNicManager manager that will set nic
    * @throws Exception
    *            if cannot connect to vCenter server or host is not found in
    *            inventory
    */
   private void setVnicForFtLogging(boolean enable, String vnicDevice,
         VirtualNicManager vnicManager) throws Exception {
      if (vnicManager == null) {
         throw new IllegalArgumentException("Invlaid Virtual Nic Manager!");
      }

      if (enable) {
         vnicManager.selectVnic(
               VirtualNicManager.NicType.faultToleranceLogging.name(),
               vnicDevice);
      } else {
         vnicManager.deselectVnic(
               VirtualNicManager.NicType.faultToleranceLogging.name(),
               vnicDevice);
      }
   }

   /**
    * Validates that the hostSpec is usable, i.e. that it is not null, that
    * there is a parent, and assigned host name
    *
    * @param hostSpec
    *           Host specification of the host
    * @throws IllegalArgumentException
    *            if host spec is null, doesn't have name or parent
    */
   private void validateHostSpec(HostSpec hostSpec)
         throws IllegalArgumentException {
      if (hostSpec == null || !hostSpec.parent.isAssigned()) {
         throw new IllegalArgumentException("Host/Parent spec is not set");
      }
      if (!hostSpec.name.isAssigned() || hostSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("Host name is not set");
      }
   }

   /**
    * Waits for the specified timeout in seconds to verify the host has entered
    * the specified state
    *
    * @param state
    *           Host state to be reached
    * @param timeout
    *           Timeout in seconds to wait for the host to get to the specified
    *           state, i.e. 60 means 60 seconds
    * @return true if host entered the state in the specified timeout and false
    *         if timeout is over or host cannot enter the state
    * @throws Exception
    *            if illegal state is supplied
    */
   private boolean waitForHostState(HostState state, HostSystem host,
         int timeout) throws Exception {

      // wait for timeout to get host in state
      int retries = timeout;
      while (!isHostInState(host, state) && retries > 0) {
         Thread.sleep(1000);
         retries--;
      }
      // if timeout hasn't passed return true
      if (retries > 0) {
         return true;
      }

      _logger.error("Timeout has passed without reaching the expected state: "
            + state.name());
      return false;
   }

   /**
    * Method that allows enabling / disabling of firewall rulesets
    *
    * @param hostSpec
    *           - host on which teh action will be performed
    * @param ruleSet
    *           - rule
    * @param enable
    *           -true to enable, false to disable
    * @throws Exception
    */
   private void setFirewallRuleSet(HostSpec hostSpec, String ruleSet, boolean enable)
         throws Exception {
      validateHostSpec(hostSpec);
      // Get the host
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      FirewallSystem firewallSystem =
            VcServiceUtil.getVcService(hostSpec.service.get()).getManagedObject(
                  host.getConfigManager().getFirewallSystem());
      if (enable) {
         firewallSystem.enableRuleset(ruleSet);
      } else {
         firewallSystem.disableRuleset(ruleSet);
      }
   }

   private HostState getHostState(HostUtil.HostStates hostState) {
      switch (hostState) {
         case CONNECTED:
            return HostState.CONNECTED;
         case DISCONNECTED:
            return HostState.DISCONNECTED;
         case ENTER_MAINTENANCE_MODE:
            return HostState.ENTER_MAINTENANCE_MODE;
         case EXIT_MAINTENANCE_MODE:
            return HostState.EXIT_MAINTENANCE_MODE;
         case NOT_RESPONDING:
            return HostState.NOT_RESPONDING;
         default:
            throw new RuntimeException("Host state not supported!");
      }
   }
}
