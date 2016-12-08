/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.Datacenter;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.host.DatastoreSystem;
import com.vmware.vim.binding.vim.host.VmfsDatastoreCreateSpec;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * Class that is used to perform specific operations directly on the ESX host
 */
public class HostClientBasicSrvApi {

   private static final Logger _logger =
         LoggerFactory.getLogger(HostClientBasicSrvApi.class);
   private static HostClientBasicSrvApi instance = null;
   protected HostClientBasicSrvApi() {}

   /**
    * Get instance of HostClientBasicSrvApi.
    *
    * @return  created instance
    */
   public static HostClientBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized(HostClientBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing HostClientSrvApi.");
               instance = new HostClientBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Method that creates a VMFS datastore on a ESX host through connecting to the host directly
    * @param datastoreSpec - the datastore which to create, the service spec should be HostServiceSpec
    * of the host on which the datastore will be created
    * @return true if successful, and false otherwise
    * @throws Exception - if the specs are not correct, or the requested datastore is not of VMFS type,
    * or if there is a duplicate name or configuration fault with the host, or if cannot connect to the host,
    * or if there is no available space for creation of a vmfs storage
    */
   public boolean createVmfsDatastore(DatastoreSpec datastoreSpec) throws Exception{
      DatastoreBasicSrvApi.getInstance().validateDatastoreSpec(datastoreSpec);

      if (!datastoreSpec.type.get().equals(DatastoreType.VMFS)) {
         throw new Exception("Datastore type shoudl be VMFS");
      }

      DatastoreSystem datastoreSystem = getDatastoreSystem(datastoreSpec.parent.get());
      VmfsDatastoreCreateSpec vmfsCreateSpec = DatastoreBasicSrvApi.getInstance().getVmfsSpec(datastoreSystem, datastoreSpec);

      if (vmfsCreateSpec == null) {
         throw new Exception("There are no available disks for creation of a vmfs storage");
      }

      return datastoreSystem.createVmfsDatastore(vmfsCreateSpec) != null;
   }

   public boolean createNasDatastore(DatastoreSpec datastoreSpec) throws Exception{
      DatastoreBasicSrvApi.getInstance().validateDatastoreSpec(datastoreSpec);

      if (!datastoreSpec.type.get().equals(DatastoreType.VMFS)) {
         throw new Exception("Datastore type shoudl be VMFS");
      }

      DatastoreSystem datastoreSystem = getDatastoreSystem(datastoreSpec.parent.get());

      VmfsDatastoreCreateSpec vmfsCreateSpec = DatastoreBasicSrvApi.getInstance().getVmfsSpec(datastoreSystem, datastoreSpec);

      if (vmfsCreateSpec == null) {
         throw new Exception("There are no available disks for creation of a vmfs storage");
      }

      return datastoreSystem.createVmfsDatastore(vmfsCreateSpec) != null;
   }


   /**
    * Method that deletes all of the vmfs datastores of the host
    * @param hostSpec - the host whose datastores to delete, with service spec of type
    * HostServiceSpec of the host
    * @return - true if successful, and false otherwise
    * @throws Exception - if the host spec is incorrect, or if cannot connect to the host
    */
   //TODO cannot delete a vmfs datastore if it is shared; it seems that through host client
   // connection teh object that shows hosts using the datastore is empty and it is unclear
   // if through host client connection there is the matter of shared datastore
   public boolean deleteAllVmfsDatastores(HostSpec hostSpec) throws Exception {
      boolean destroyed = true;

      validateHostSpec(hostSpec);

      List<Datastore> datastores = getDatastoresByType(hostSpec, DatastoreType.VMFS);


      for (Datastore datastore : datastores) {
         destroyed = destroyed && VcServiceUtil.waitForTaskSuccess(datastore.destroy(), hostSpec);
      }

      return destroyed;
   }

   /**
    * Method that deletes a vmfs datastores of the host by datastroe's name
    * @param datastoreSpec - the datastore to delete, with service spec of type
    * HostServiceSpec of the host
    * @return - true if successful, and false otherwise
    * @throws Exception - if the datastore spec is incorrect, or if cannot connect to the host
    */
   public boolean deleteVmfsDatastore(DatastoreSpec datastoreSpec) throws Exception {

      DatastoreBasicSrvApi.getInstance().validateDatastoreSpec(datastoreSpec);

      Datastore datastore = getDatastoresByName(datastoreSpec, datastoreSpec.name.get());

      return VcServiceUtil.waitForTaskSuccess(datastore.destroy(), datastoreSpec);
   }

   /**
    * Return true if datastore with the specified name is found on the host.
    * @param datastoreSpec
    * @return
    * @throws Exception
    */
   public boolean findVmfsDatastoreByName(DatastoreSpec datastoreSpec) throws Exception{
      DatastoreBasicSrvApi.getInstance().validateDatastoreSpec(datastoreSpec);

      try {
         getDatastoresByName(datastoreSpec, datastoreSpec.name.get());
         return true;
      } catch (IllegalArgumentException iae) {
         return false;
      }
   }


   /**
    * Method that enables vmotion if it is not set
    * @param hostSpec - specification of the host on which to enable vmotion
    * @return true if successfully set or already set, false if not set and no available
    * vnics found
    * @throws Exception
    */
   public boolean enableVmotion(HostSpec hostSpec) throws Exception {
      return HostBasicSrvApi.getInstance().setVmotion(true, getHostSystem(hostSpec), VcServiceUtil.getVcService(hostSpec.service.get()));
   }

   /**
    * Method that disables vmotion if it is not set
    * @param hostSpec - specification of the host on which to enable vmotion
    * @return true if successfully disabled or already disabled, false otherwise
    * @throws Exception
    */
   public boolean disableVmotion(HostSpec hostSpec) throws Exception {
      return HostBasicSrvApi.getInstance().setVmotion(false, getHostSystem(hostSpec), VcServiceUtil.getVcService(hostSpec.service.get()));
   }
   //---------------------------------------------------------------------------
   // Private methods

   // Method that validates the host spec is assigned and has a name
   private void validateHostSpec(HostSpec hostSpec) {

      if (hostSpec == null) {
         throw new IllegalArgumentException("Host spec is not set");
      }
      if (!hostSpec.name.isAssigned() || hostSpec.name.get().isEmpty()) {
         throw new IllegalArgumentException("Host name is not set");
      }
   }

   // Method that returns the datastore system of the specified host
   private DatastoreSystem getDatastoreSystem(EntitySpec entitySpec) throws Exception {

      VcService service = VcServiceUtil.getVcService(entitySpec);
      HostSystem host = getHostSystem((HostSpec) entitySpec);
      return service.getManagedObject(host.getConfigManager().getDatastoreSystem());

   }

   // Method that gets all host's datastores of a certain type
   private List<Datastore> getDatastoresByType(EntitySpec entitySpec, DatastoreType dsType)
         throws Exception {

      VcService service = VcServiceUtil.getVcService(entitySpec);
      List<Datastore> datastores = new ArrayList<Datastore>();

      Folder folder = service.getManagedObject(getDatacenter(entitySpec).getDatastoreFolder());
      for (ManagedObjectReference childMor : folder.getChildEntity()) {
         Datastore datastore = service.getManagedObject(childMor);

         if (datastore.getSummary().getType().equals(dsType.name())) {
            datastores.add(datastore);
         }
      }

      return datastores;
   }

   // Method that returns the host's datastore specified by the name parameter
   private Datastore getDatastoresByName(EntitySpec entitySpec, String name)
         throws Exception {

      VcService service = VcServiceUtil.getVcService(entitySpec);

      Folder folder = service.getManagedObject(getDatacenter(entitySpec).getDatastoreFolder());
      for (ManagedObjectReference childMor : folder.getChildEntity()) {
         Datastore datastore = service.getManagedObject(childMor);

         if (datastore.getName().equals(name)) {
            return datastore;
         }
      }

      throw new IllegalArgumentException(String.format("Datastore with name %s is not found on host", name));
   }

   // Method that gets the host datacenter - in host inventory there is a logical
   // datacenter and root folder
   private Datacenter getDatacenter(EntitySpec entitySpec) throws Exception {

      VcService service = VcServiceUtil.getVcService(entitySpec);
      Folder rootFolder =
            service
                  .getManagedObject(service.getServiceInstanceContent().getRootFolder());
      for (ManagedObjectReference childMor : rootFolder.getChildEntity()) {
         if (childMor.getType().equals(
               ManagedEntityUtil.EntityType.DATACENTER.getValue())) {
            Datacenter datacenter = service.getManagedObject(childMor);
            return datacenter;
         }
      }
      return null;
   }

   // Method that gets the host by its spec
   private HostSystem getHostSystem(HostSpec hostSpec) throws Exception {
      validateHostSpec(hostSpec);

      VcService service = VcServiceUtil.getVcService(hostSpec);

      Datacenter datacenter = getDatacenter(hostSpec);
      if (datacenter == null) {
         throw new Exception("The host doesn't have a datacenter!");
      }
      Folder hostFolder = service.getManagedObject(datacenter.getHostFolder());

      // Actually it is an only child and is a logical compute resource for the host
      ComputeResource cr = service.getManagedObject(hostFolder.getChildEntity()[0]);
      // Same here - it is an only host
      HostSystem host = service.getManagedObject(cr.getHost()[0]);

      return host;
   }
}
