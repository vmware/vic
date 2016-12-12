/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.net.URI;
import java.util.ArrayList;
import java.util.List;
import java.util.Stack;

import javax.activation.UnsupportedDataTypeException;

import org.apache.commons.lang3.ArrayUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.sso.SsoClient;
import com.vmware.client.automation.util.SsoUtil;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.pbm.ServerObjectRef;
import com.vmware.vim.binding.vim.ClusterComputeResource;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.Datacenter;
import com.vmware.vim.binding.vim.Datastore;
import com.vmware.vim.binding.vim.DistributedVirtualSwitch;
import com.vmware.vim.binding.vim.Folder;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.ManagedEntity;
import com.vmware.vim.binding.vim.Network;
import com.vmware.vim.binding.vim.ResourcePool;
import com.vmware.vim.binding.vim.ServiceInstanceContent;
import com.vmware.vim.binding.vim.StoragePod;
import com.vmware.vim.binding.vim.VirtualApp;
import com.vmware.vim.binding.vim.VirtualMachine;
import com.vmware.vim.binding.vim.dvs.DistributedVirtualPortgroup;
import com.vmware.vim.binding.vim.vApp.EntityConfigInfo;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.binding.vmodl.fault.ManagedObjectNotFound;
import com.vmware.vim.query.client.Client;
import com.vmware.vim.vmomi.cis.CisIdConverter;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatacenterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DvPortgroupSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DvsSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.NetworkSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ResourcePoolSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VappSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec.FaultToleranceVmRoles;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * This is a class used for traversal of the Vc Inventory and finding the
 * ManagedObject that corresponds to the supplied ManageEntitySpecification
 * (folder, datacenter, cluster, etc.)
 */
public class ManagedEntityUtil {

   private static final Logger _logger = LoggerFactory
         .getLogger(ManagedEntityUtil.class);

   private static final String DATACENTER_TYPE = "Datacenter";
   private static final String FOLDER_TYPE = "Folder";
   private static final String HOST_TYPE = "HostSystem";
   private static final String COMPUTE_RESOURCE_TYPE = "ComputeResource";
   private static final String CLUSTER_TYPE = "ClusterComputeResource";
   private static final String RESOURCE_POOL_TYPE = "ResourcePool";
   private static final String VAPP_TYPE = "VirtualApp";
   private static final String VM_TYPE = "VirtualMachine";
   private static final String DATASTORE_TYPE = "Datastore";
   private static final String DVS_TYPE = "VmwareDistributedVirtualSwitch";
   private static final String VC_TYPE = "TODO: add correct type";
   private static final String DVPORTGROUP_TYPE = "DistributedVirtualPortgroup";
   private static final String DS_CLUSTER_TYPE = "StoragePod";
   private static final String NETWORK_TYPE = "Network";

   // VAPI moref template. To address vmodl entities in the VAPI
   // requests(vmodl2)
   // their moref have to be converted to the given format.
   // The three parameters are the ManagedObjectReference type, value and
   // server guid.
   private static final String VAPI_MOREF_TEMPLATE = "urn:vmomi:%s:%s:%s";// "moref:%s:%s:%s";

   /**
    * Enum that holds the Vim Types of the Inventory objects; enum is with default access as
    * it is used in HostClientSrvApi
    */
   static enum EntityType {
      DATACENTER(DATACENTER_TYPE), FOLDER(FOLDER_TYPE), HOST(HOST_TYPE),
      CLUSTER(CLUSTER_TYPE), RESOURCE_POOL(RESOURCE_POOL_TYPE), VAPP(VAPP_TYPE), VM(VM_TYPE),
      DATASTORE(DATASTORE_TYPE), DVS(DVS_TYPE), DVPORTGROUP(DVPORTGROUP_TYPE),
      VC(VC_TYPE), DATASTORE_CLUSTER(DS_CLUSTER_TYPE), NETWORK(NETWORK_TYPE);

      private final String value;

      private EntityType(String value) {
         this.value = value;
      }

      public String getValue() {
         return value;
      }
   }

   /**
    * Class that holds the representation of a node in the inventory - name,
    * type and if it is a folder, it holds the folder type
    */
   private final class HierarchyItem {
      private String _childName;
      private EntityType _childType;
      private FolderType _folderType;

      /**
       * The constructor gets a ManagedEntitytSpec and extracts it in a suitable
       * form to create a hierarchy item
       */
      public HierarchyItem(ManagedEntitySpec spec) {
         // Get the child type - folder, host, etc.
         setChildType(getType(spec));
         // If the child is folder, then the type of the folder will be
         // needed
         if (getChildType().equals(EntityType.FOLDER)) {
            FolderSpec fSpec = (FolderSpec) spec;
            setFolderType(fSpec.type.get());
         }
         // Set the name as this and its type will actually uniquely identify
         // the entity in the tree
         setChildName(spec.name.get());
      }

      public String getChildName() {
         return _childName;
      }

      private void setChildName(String childName) {
         this._childName = childName;
      }

      public EntityType getChildType() {
         return _childType;
      }

      private void setChildType(EntityType type) {
         this._childType = type;
      }

      public FolderType getFolderType() {
         return _folderType;
      }

      private void setFolderType(FolderType folderType) {
         this._folderType = folderType;
      }

      @Override
      public String toString() {
         return String.format("\"%s\" (%s)", this._childName, this._childType);
      }
   }

   /**
    * Method that gets the Managed Object that corresponds to the supplied
    * specification
    *
    * @param spec
    *            specification of an inventory entity
    * @return Managed object that corresponds to the specified entity or {@link ObjectNotFoundException} if the object
    *         was not found in the
    *         inventory
    * @throws Exception
    * @throws ObjectNotFoundException
    */
   public static <T extends ManagedObject> T getManagedObject(ManagedEntitySpec spec) throws Exception {
      return getManagedObject(spec, spec.service.get());
   }

   /**
    * Getting ManagedObjectReference for object with given spec.
    *
    * @param spec
    * @return
    * @throws Exception
    */
   public static ManagedObjectReference getManagedObjectReference(
         ManagedEntitySpec spec) throws Exception {
      ManagedObjectReference mor = null;
      // As categories and tags are not part of the hierarchy it can't be used
      // the same logic for getting the moRef of VC object.
      if (spec instanceof BackingCategorySpec) {
         mor = getTagCategoryMoRef(spec);
      } else if (spec instanceof BackingTagSpec) {
         mor = getTagMoRef((BackingTagSpec) spec);
      } else if (spec instanceof VmSpec) {
         FaultToleranceVmRoles role = ((VmSpec) spec).ftRole.isAssigned() ? ((VmSpec) spec).ftRole
               .get() : null;
         mor = getHierarchyItemMoRef(spec, role);
      } else {
         mor = getHierarchyItemMoRef(spec, null);
      }
      return mor;
   }

   /**
    * Method creating the hierarchy and searching through it to get
    * ManagedObjectReference for the object, which spec is passed
    *
    * @param spec
    * @param ftRole
    *           - the fault tolerance role we are looking for, valid only for
    *           VMs
    * @return
    * @throws Exception
    * @throws ObjectNotFoundException
    */
   private static ManagedObjectReference getHierarchyItemMoRef(
         ManagedEntitySpec spec, FaultToleranceVmRoles ftRole) throws Exception {
      ManagedObjectReference mor;
      Stack<HierarchyItem> hierarchy = getHierarchy(spec);
      VcService service = VcServiceUtil.getVcService(spec.service.get());
      ServiceInstanceContent serviceContent = service
            .getServiceInstanceContent();

      // Starts search from root folder always and finish when it gets the
      // mor of the last element of the hierarchy - this is the one that is
      // sought
      mor = serviceContent.getRootFolder();
      List<HierarchyItem> searchedPath = new ArrayList<>();
      while (!hierarchy.empty()) {
         HierarchyItem currentElement = hierarchy.pop();
         searchedPath.add(currentElement);
         mor = getVcObjectManagedObjectRerefence(mor, currentElement, service,
               ftRole);
      }

      // in case an empty hierarchy is passed, the root folder shouldn't be
      // returned
      if (mor == null) {
         throw new ObjectNotFoundException(String.format(
               "Object on path: \"%s\" was not found",
               ArrayUtils.toString(searchedPath.toArray())));
      }
      // Manually creating the mor and setting the properties. Without doing
      // this serverGuid property will always be null.
      mor = new ManagedObjectReference(mor.getType(), mor.getValue(),
            service.getServiceGuid());
      return mor;
   }

   /**
    * Method that gets the Managed Object that corresponds to the supplied
    * specification
    *
    * @param spec
    *           specification of an inventory entity
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return Managed object that corresponds to the specified entity or {@link ObjectNotFoundException} if the object
    *         was not found in the
    *         inventory
    * @throws Exception
    * @throws ObjectNotFoundException
    */
   public static <T extends ManagedObject> T getManagedObject(
         ManagedEntitySpec spec, ServiceSpec serviceSpec) throws Exception {
      ManagedObjectReference mor = getManagedObjectReference(spec);
      VcService vcService = VcServiceUtil.getVcService(spec.service.get());
      return vcService.getManagedObject(mor);
   }

   /**
    * Waits for specified timeout to verify that a managed object is removed
    * from the VC Inventory
    *
    * @param spec
    *            Specification of the managed entity
    * @param timeoutIn
    *            Max number of seconds to wait
    * @return True if object is deleted successfulluy in the specified timeout,
    *         otherwise false
    * @throws Exception
    *             if failed to login to vc service
    */
   public static boolean waitForEntityDeletion(ManagedEntitySpec spec, int timeout) throws Exception {
      int retries = timeout;
      int logRetries = 0;
      String name = "";
      // timeout of specified timeout to let the operation really finish
      while (retries > 0) {
         _logger.info("Entering retry " + logRetries);
         Thread.sleep(1000);
         try {
            ManagedEntityUtil.getManagedObject(spec);
         } catch (ObjectNotFoundException e) {
            if (spec.name.isAssigned()) {
               name = spec.name.get();
            }
            _logger.info("Object " + name + " successfully deleted");
            return true;
         }
         retries--;
         logRetries++;
      }

      _logger.error("Timeout is over and object " + name + " is still present in Inventory");
      return false;
   }

   /**
    * Utility method to convert the <code>ManagedObjectReference</code> object
    * provided by the <code>managedEntitySpec</code> to a moref format used in
    * the VAPI calls. (For example see the <code>VdcClusterService.attach</code> )
    *
    * @param managedEntitySpec
    *            provides the ManagedEntitySpec
    * @return URI string format used in the VAPI.
    * @throws Exception
    *             if no entity for the provided <code>managedEntitySpec</code> is
    *             found.
    * @deprecated Use getVapiMoRefUrl(ManagedEntitySpec managedEntitySpec,
    *             ServiceSpec serviceSpec)
    */
   @Deprecated
   public static String getVapiMoRefUrl(ManagedEntitySpec managedEntitySpec) throws Exception {
      return getVapiMoRefUrl(managedEntitySpec, managedEntitySpec.service.get());
   }

   /**
    * Utility method to convert the <code>ManagedObjectReference</code> object
    * provided by the <code>managedEntitySpec</code> to a moref format used in
    * the VAPI calls. (For example see the <code>VdcClusterService.attach</code> )
    *
    * @param managedEntitySpec
    *           provides the ManagedEntitySpec
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return URI string format used in the VAPI.
    * @throws Exception
    *             if no entity for the provided <code>managedEntitySpec</code> is
    *             found.
    */
   public static String getVapiMoRefUrl(ManagedEntitySpec managedEntitySpec, ServiceSpec serviceSpec) throws Exception {
      ManagedObjectReference managedObjectReference = getManagedObject(managedEntitySpec, serviceSpec)._getRef();
      return CisIdConverter.toGlobalCisId(
            managedObjectReference,
            VcServiceUtil.getVcService(serviceSpec).getServiceGuid()
            );
   }

   public static String getOldVapiMoRefUrl(ManagedEntitySpec managedEntitySpec,
         ServiceSpec serviceSpec) throws Exception {
      ManagedObjectReference managedObjectReference =
            getManagedObject(managedEntitySpec, serviceSpec)._getRef();
      return String.format(
            VAPI_MOREF_TEMPLATE,
            managedObjectReference.getType(),
            managedObjectReference.getValue(),
            VcServiceUtil.getVcService(serviceSpec).getServiceGuid());
   }

   /**
    * Utility method to convert a <code>ManagedObjectReference</code> object
    * provided by the <code>resourceSpec</code> to a URI used in the VAPI calls.
    * )
    *
    * @param resourceSpec
    *            provides the ManagedObjectReference object
    * @return URI format used in the vAPI.
    * @throws Exception
    *             if no entity for the provided <code>resourceSpec</code> is
    *             found.
    */
   public static URI getResourceUri(ManagedEntitySpec resourceSpec) throws Exception {
      String uri = ManagedEntityUtil.getVapiMoRefUrl(resourceSpec);
      return URI.create(uri);
   }

   /**
    * Utility method to convert a list of <code>ManagedObjectReference</code> objects provided by the
    * <code>resourceSpecs</code> to a list of URIs used
    * in the VAPI calls. )
    *
    * @param resourceSpecs
    *            provides the ManagedObjectReference list
    * @return List<URI> format used in the VAPI.
    * @throws Exception
    *             if no entity for an item in the provided <code>resourceSpecs</code> is found.
    */
   public static List<URI> getResourceUris(List<? extends ManagedEntitySpec> resourceSpecs) throws Exception {
      List<URI> resourceUris = new ArrayList<URI>();
      for (ManagedEntitySpec resSpec : resourceSpecs) {
         String uri = ManagedEntityUtil.getOldVapiMoRefUrl(
               resSpec,
               resSpec.service.isAssigned() ? resSpec.service.get() : null
            );
         resourceUris.add(URI.create(uri));
      }
      return resourceUris;
   }

   /**
    * Gets the names of the ManagedEntitySpec objects in the specs list. Useful
    * for passing the names to grid selection methods that use String arguments.
    *
    * @param specs
    *            - the list of ManagedEntitySpec objects
    * @return a list with the names of the ManagedEntitySpec objects
    */
   public static List<String> getEntityNames(List<? extends ManagedEntitySpec> specs) {
      List<String> names = new ArrayList<String>();
      for (ManagedEntitySpec spec : specs) {
         names.add(spec.name.get());
      }
      return names;
   }

   /**
    * Method that can return the name of the Managed Object from its Managed
    * Object Reference String representation in the form of:
    * "serverGUID:type:value"
    *
    * @param stringMor
    *            - string representation of the managed object reference -
    *            "serverGUID:type:value"
    * @return String with the name of the object
    * @throws Exception
    *             if the connection to vc cannot be established or if type of
    *             Managed Object Reference is not supported
    * @deprecated getNameFromStringMor(String stringMor, ServiceSpec
    *             serviceSpec)
    */
   @Deprecated
   public static String getNameFromStringMor(String stringMor) throws Exception {
      return getNameFromStringMor(stringMor, null);
   }

   /**
    * Method that can return the name of the Managed Object from its Managed
    * Object Reference String representation in the form of:
    * "serverGUID:type:value"
    *
    * @param stringMor
    *            - string representation of the managed object reference -
    *            "serverGUID:type:value"
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return String with the name of the object
    * @throws Exception
    *             if the connection to vc cannot be established or if type of
    *             Managed Object Reference is not supported
    */
   public static String getNameFromStringMor(String stringMor, ServiceSpec serviceSpec) throws Exception {
      String stringMorValues[] = stringMor.split(":");
      ManagedObjectReference moRef = new ManagedObjectReference(stringMorValues[1], stringMorValues[2],
            stringMorValues[0]);

      return getNameFromMoRef(moRef, serviceSpec);
   }

   /**
    * Method that can return the name of the Managed Object from its Managed
    * Object Reference
    *
    * @param moRef
    *           - the Managed Object Reference
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return String with the name of the object
    * @throws Exception
    *             if the connection to vc cannot be established or if type of
    *             Managed Object Reference is not supported
    */
   public static String getNameFromMoRef(ManagedObjectReference moRef, ServiceSpec serviceSpec) throws VcException {
      VcServiceUtil.getVcService(serviceSpec);

      return ((ManagedEntity) getManagedObjectFromMoRef(moRef, serviceSpec)).getName();
   }

   /**
    * Method that can return the Managed Object from its Managed Object
    * Reference
    *
    * @param moRef
    *           - the Managed Object Reference
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return the Managed Object itself
    * @throws Exception
    *             if the connection to vc cannot be established or if type of
    *             Managed Object Reference is not supported
    */
   public static <T extends ManagedObject> T getManagedObjectFromMoRef(ManagedObjectReference moRef,
         ServiceSpec serviceSpec) throws VcException {
      VcService service = VcServiceUtil.getVcService(serviceSpec);

      return service.getManagedObject(moRef);
   }

   /**
    * Converts ManagedObjectReference to ServerObjectRef.
    *
    * @param moRef
    *           the ManagedObjectReference which will be converted to
    *           ServerObjectRef.
    * @return newly created ServerObjectRef.
    */
   public static ServerObjectRef managedObjectRefToServerObjectRef(
         ManagedObjectReference moRef) {
      ServerObjectRef soRef = new ServerObjectRef();
      soRef.key = moRef.getValue();
      soRef.objectType = moRefTypeToSoRefType(moRef.getType());
      return soRef;
   }

   /**
    * Converts ManagedObjectReference type to ServerObjectRef type.
    *
    * @param moRefType ManagedObjectReference type to be converted.
    * @return Value from the ServerObjectRef.ObjectType enum.
    */
   private static String moRefTypeToSoRefType(String moRefType) {
      if (VM_TYPE.equals(moRefType)) {
         return ServerObjectRef.ObjectType.virtualMachine.toString();
      }
      if (DATASTORE_TYPE.equals(moRefType)) {
         return ServerObjectRef.ObjectType.datastore.toString();
      }
      return moRefType;
   }

   /**
    * The method takes a spec and traverses all it parents creating a list of
    * HierarchyItems The first element of the list is the one that is located in
    * the root folder, each next element is the child of the previous
    */
   private static Stack<HierarchyItem> getHierarchy(ManagedEntitySpec spec) {
      Stack<HierarchyItem> result = new Stack<HierarchyItem>();

      if (spec == null) {
         return result;
      }
      ManagedEntitySpec tempSpec = spec;
      // ManagedEntitySpec the hierarchy of children as the one that is actually being
      // sought remains as the last element of the list
      while (tempSpec != null) {
         HierarchyItem child = new ManagedEntityUtil().new HierarchyItem(tempSpec);
         // Add ea element at the beginning of the list so the structure
         // corresponds to the structure in vc starting from root folder
         result.push(child);
         // Get the next element
         if (tempSpec.parent.isAssigned() && !(tempSpec.parent.get() instanceof VcSpec)) {
            tempSpec = tempSpec.parent.get();
         } else {
            tempSpec = null;
         }
      }

      return result;
   }

   /**
    * Method that gets the entity's type from its specification
    *
    * @param managedEntitySpec
    * @return
    */
   private static EntityType getType(EntitySpec managedEntitySpec) {

      if (managedEntitySpec instanceof DatacenterSpec) {
         return EntityType.DATACENTER;
      } else if (managedEntitySpec instanceof ClusterSpec) {
         return EntityType.CLUSTER;
      } else if (managedEntitySpec instanceof FolderSpec) {
         return EntityType.FOLDER;
      } else if (managedEntitySpec instanceof HostSpec) {
         return EntityType.HOST;
      } else if (managedEntitySpec instanceof VappSpec) {
         return EntityType.VAPP;
      } else if (managedEntitySpec instanceof ResourcePoolSpec) {
         return EntityType.RESOURCE_POOL;
      } else if (managedEntitySpec instanceof VmSpec) {
         return EntityType.VM;
      } else if (managedEntitySpec instanceof DatastoreSpec) {
         return EntityType.DATASTORE;
      } else if (managedEntitySpec instanceof DvsSpec) {
         return EntityType.DVS;
      } else if (managedEntitySpec instanceof DvPortgroupSpec) {
         return EntityType.DVPORTGROUP;
      } else if (managedEntitySpec instanceof DatastoreClusterSpec) {
         return EntityType.DATASTORE_CLUSTER;
      } else if (managedEntitySpec instanceof NetworkSpec) {
         return EntityType.NETWORK;
      } else if (managedEntitySpec instanceof VcSpec) {
         return EntityType.VC;
      } else {
         throw new UnsupportedOperationException(
               "No Entity type corresponds to the supplied specification: "
                     + managedEntitySpec.getClass());
      }
   }

   /**
    * Method that gets the MOR of a specified HierarchyItem and parent - mor
    *
    * @param parentMor
    *           - the mor of the parent that for whose child is sought
    * @param child
    *           - the child whose mor is sought
    * @param ftRole
    *           - the fault tolerance role we are looking for (1 is primary, 2
    *           is secondary), valid only for VMs
    *
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getVcObjectManagedObjectRerefence(
         ManagedObjectReference parentMor, HierarchyItem child,
         VcService service, FaultToleranceVmRoles ftRole) throws Exception {
      // null will be retuned if no object is found
      ManagedObjectReference result = null;
      // Get the managed object of the parent in whose children the child will
      // be sought
      ManagedObject parent = service.getManagedObject(parentMor);

      // Depending on the child type, the search differs
      switch (child.getChildType()) {
      case VC:
         // As the use case for retrieving the "vc managed object" is to be used for navigation the current
         // implementation return the root folder of the inventory.
         result = service.getServiceInstanceContent().getRootFolder();
         break;
      case DATACENTER:
         // In the case for search for datacenter it could only be in a
         // folder - root or dc folder
         // So parent is cast to a folder
         result = getDatacenterMoRef(parent, child, service);
         break;
      case FOLDER:
         // Folder can be in a folder or in a datacenter's folders, in the
         // latter
         // case parent depends on the type of the folder
         // get parent
         result = getFolderMoRef(parent, child, service);
         break;
      case CLUSTER:
         // Clusters are located in the host folder of a datacenter
         result = getClusterMoRef(parent, child, service);
         break;
      case HOST:
         // Hosts are located in a compute resource
         result = getHostMoRef(parent, child, service);
         break;
      case RESOURCE_POOL:
         // Resource pools are located in a compute resource
         result = getResourcePoolMoRef(parent, child, service);
         break;
      case VAPP:
         // Vapps are located in the VM folder of a datacenter
         result = getVappMoRef(parent, child, service);
         break;
      case VM:
         // VMs are located in the VM folder of a datacenter
         result = getVmMoRef(parent, child, service, ftRole);
         break;
      case DATASTORE:
         // Datastores are located in the Storage folder of a datacenter
         result = getDatastoreMoRef(parent, child, service);
         break;
      case DVS:
         // DVSes are located in the Networking folder of a datacenter
         result = getDvsMoRef(parent, child, service);
         break;
      case DVPORTGROUP:
         // dvPortgroups are located in dvs
         result = getDvPortgroupMoRef(parent, child, service);
         break;
      case DATASTORE_CLUSTER:
         // DS clusters are located under Datastore folder of a datacenter
         result = getDsClusterMoRef(parent, child, service);
         break;
      case NETWORK:
         // Networks are located in the Networking folder of a datacenter
         result = getNetworkMoRef(parent, child, service);
         break;
      default:
         throw new UnsupportedDataTypeException("No such Entity Type: " + child.getChildType());
      }

      return result;
   }

   /**
    * Gets the Managed Object Reference of a Tag Category by creating connection
    * to the TaggingManager, which retrieves the information about the tags and
    * categories objects.
    *
    * @param spec
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getTagCategoryMoRef(
         ManagedEntitySpec spec) throws Exception {
      SsoClient ssoConnection = SsoUtil.getVcConnector(spec).getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      return BackingTagsBasicSrvApi.getInstance().getBackingCategory(
            queryClient, spec);
   }

   /**
    * Gets the Managed Object Reference of a Tag by creating connection to the
    * TaggingManager, which retrieves the information about the tags and
    * categories objects.
    *
    * @param spec
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getTagMoRef(BackingTagSpec spec)
         throws Exception {
      SsoClient ssoConnection = SsoUtil.getVcConnector(spec).getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      BackingCategorySpec backingCategorySpec = spec.category.get();
      ManagedObjectReference backingCategory = BackingTagsBasicSrvApi
            .getInstance().getBackingCategory(queryClient, backingCategorySpec);
      return BackingTagsBasicSrvApi.getInstance().getBackingTag(queryClient,
            backingCategory, spec);
   }

   /**
    * Gets the Managed Object Reference of a Datastore by searching for it in
    * the Child Entities of the parent Datacenter's Storage folder depending on
    * its type.
    *
    * @param parent
    *            the managed object that is parent of the Datastore it could be
    *            Datacenter or HostSystem
    * @param child
    *            the hierarchy item that corresponds to the Datastore sought
    * @return
    * @throws Exception
    *             if the VcService cannot be obtained
    */
   private static ManagedObjectReference getDatastoreMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      ManagedObjectReference[] datastores = null;
      if (parent instanceof Datacenter) {
         Folder parentFolder = service.getManagedObject(((Datacenter) parent).getDatastoreFolder());
         datastores = parentFolder.getChildEntity();
      } else if (parent instanceof HostSystem) {
         datastores = ((HostSystem) parent).getDatastore();
      }

      for (ManagedObjectReference childMor : datastores) {
         if (childMor.getType().equals(EntityType.DATASTORE.getValue())) {
            Datastore datastore = service.getManagedObject(childMor);
            if (datastore.getName().equals(child.getChildName())) {
               return datastore._getRef();
            }
         }
      }

      return null;
   }

   /**
    * Gets the Managed Object Reference of a VM by searching for it in the Child
    * Entities of the parent Datacenter's VM folder depending on its type
    * (VirtualMachine);
    *
    * @Note Currently implementation works only if VM's parent is set as host
    *
    * @param parent
    *           - the managed object that is parent (host, datacenter, folder)
    *           of the VM
    * @param child
    *           - the hierarchy item that corresponds to the VM sought
    * @param ftRole
    *           - the fault tolerance role we are looking for (1 is primary, 2
    *           is secondary)
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getVmMoRef(ManagedObject parent,
         HierarchyItem child, VcService service, FaultToleranceVmRoles ftRole)
         throws Exception {
      Folder parentFolder = null;
      if (parent instanceof Folder) {
         parentFolder = (Folder) parent;
      } else if (parent instanceof Datacenter) {
         parentFolder = service.getManagedObject(((Datacenter) parent).getVmFolder());
      } else if ((parent instanceof HostSystem) || (parent instanceof ClusterComputeResource)
            || (parent instanceof VirtualApp) || (parent instanceof ResourcePool)) {
         Datacenter parentDC = getParentDatacenter((ManagedEntity) parent, service);
         parentFolder = service.getManagedObject(parentDC.getVmFolder());
      } else {
         throw new RuntimeException(String.format(
               "Get VM %s failed, because it's parent is not supported: %s",
               child.getChildName(),
               parent._getRef().getType()));
      }
      return getVmMorefFromNestedVmFolders(parentFolder, child.getChildName(),
            service, ftRole);
   }

   /**
    * Helper method that searches for a vApp recursively in the specified VM
    * folder location and in all nested VM folders. Searching for nested vApp
    * inside vApp is not yet supported.
    *
    * @param parentFolder
    *            - the folder from which to start the search
    * @param vappName
    *            - the name of the vApp whose ManagedObjectReference to get
    * @return - the ManagedObjectReference for the found vApp; null if not found
    * @throws Exception
    */
   private static ManagedObjectReference getVappMorefFromNestedVmFolders(Folder parentFolder, String vappName, VcService service)
         throws Exception {
      ManagedObjectReference result = null;

      for (ManagedObjectReference childMor : parentFolder.getChildEntity()) {

         try {
            // Search the first level vApps
            // TODO Add support for nested vApps inside vApps
            if (childMor.getType().equals(EntityType.VAPP.getValue())) {
               VirtualApp vapp = service.getManagedObject(childMor);
               if (vapp.getName().equals(vappName)) {
                  result = vapp._getRef();
                  break;
               }
            }

            // Recursively search the folders
            if (childMor.getType().equals(EntityType.FOLDER.getValue())) {
               result = getVappMorefFromNestedVmFolders(
                           (Folder) service.getManagedObject(childMor),
                           vappName,
                           service);
               if (result != null) {
                  break;
               }
            }
         } catch (ManagedObjectNotFound monf) {
            _logger.warn("Managed object not available at the system:" + monf.getMessage());
         }
      }

      return result;
   }

   /**
    * Helper method that iterates in the default VM folder and all custom VM
    * folders (that are actually subfolders of the default vM folder) for the VM
    * we are looking for.
    *
    * @param parentFolder
    *           - the folder from which to start the search
    * @param vmName
    *           - the name of teh VM whoe ManagedObjectReference to get
    * @param ftRole
    *           - the fault tolerance role we are looking for
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getVmMorefFromNestedVmFolders(Folder parentFolder, String vmName,
         VcService service, FaultToleranceVmRoles ftRole) throws Exception {
      ManagedObjectReference result = null;

      for (ManagedObjectReference childMor : parentFolder.getChildEntity()) {
         // search the first level VMs
         try {
            if (childMor.getType().equals(EntityType.VM.getValue())) {
               VirtualMachine vm = service.getManagedObject(childMor);
               if (vm.getName().equals(vmName)) {
                  if (ftRole != null) {
                     try {
                        if (ftRole.getRole() == vm.getConfig().getFtInfo()
                              .getRole()) {
                           result = vm._getRef();
                           break;
                        }
                     } catch (NullPointerException e) {
                        // Fault tolerance has not been turned on yet
                        result = vm._getRef();
                        break;
                     }
                  } else {
                     result = vm._getRef();
                     break;
                  }
               }
            }
         } catch (ManagedObjectNotFound monf) {
            _logger.warn("Vm object not available on the system." + monf.toString());
            continue;
         }

         // search the VMs inside vApps
         if (childMor.getType().equals(EntityType.VAPP.getValue())) {
            try {
               VirtualApp vapp = service.getManagedObject(childMor);
               for (ManagedObjectReference vmMor : vapp.getVm()) {
                  VirtualMachine vm = service.getManagedObject(vmMor);
                  if (vm.getName().equals(vmName)) {
                     if (ftRole != null) {
                        try {
                           if (ftRole.getRole() == vm.getConfig().getFtInfo()
                                 .getRole()) {
                              result = vm._getRef();
                              break;
                           }
                        } catch (NullPointerException e) {
                           // Fault tolerance has not been turned on yet
                           result = vm._getRef();
                           break;
                        }
                     } else {
                        result = vm._getRef();
                        break;
                     }
                  }
               }
            } catch (ManagedObjectNotFound monf) {
               _logger.warn("vApp object not available on the system." + monf.toString());
               continue;
            }
            if (result != null) {
               break;
            }
         }

         // recursively search the folders
         if (childMor.getType().equals(EntityType.FOLDER.getValue())) {
            result = getVmMorefFromNestedVmFolders(
                  (Folder) service.getManagedObject(childMor), vmName, service,
                  ftRole);
            if (result != null) {
               break;
            }
         }
      }

      return result;

   }

   /**
    * Gets the Managed Object Reference of a datacenter as it searches it in the
    * Child Entities of the parent folder of the datacenter
    *
    * @param parent
    *            - the managed object that is parent of the specified datacenter
    * @param child
    *            - the hierarchy item that corresponds to the datacenter sought
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getDatacenterMoRef(ManagedObject parent, HierarchyItem child,
         VcService service) throws Exception {
      _logger.info("Looking for DC with name:" + child.getChildName());
      Folder parentFolder = (Folder) parent;
      for (ManagedObjectReference childMor : parentFolder.getChildEntity()) {
         if (childMor.getType().equals(EntityType.DATACENTER.getValue())) {
            try {
               Datacenter datacenter = service.getManagedObject(childMor);
               if (datacenter.getName().equals(child.getChildName())) {
                  return datacenter._getRef();
               }
            } catch (ManagedObjectNotFound monf) {
               _logger.error(monf.getMessage());
               continue;
            }
         }
      }
      return null;
   }

   /**
    * Gets the Managed Object Reference of a Folder by searching for it in the
    * Child Entities of the parent folder or the specific Datacenter folder
    * depending on its type
    *
    * @param parent
    *            - the managed object that is parent of the folder
    * @param child
    *            - the hierarchy item that corresponds to the folder sought
    * @param service
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getFolderMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      Folder parentFolder = null;
      if (parent._getRef().getType().equals(FOLDER_TYPE)) {
         parentFolder = (Folder) parent;
      } else {
         // get parent depending on folder type
         Datacenter parentDC = (Datacenter) parent;
         switch (child.getFolderType()) {
         case HOST:
            parentFolder = (Folder) service.getManagedObject(parentDC.getHostFolder());
            break;
         case NETWORK:
            parentFolder = (Folder) service.getManagedObject(parentDC.getNetworkFolder());
            break;
         case STORAGE:
            parentFolder = (Folder) service.getManagedObject(parentDC.getDatastoreFolder());
            break;
         case VM:
            parentFolder = (Folder) service.getManagedObject(parentDC.getVmFolder());
            break;
         default:
            parentFolder = (Folder) service.getManagedObject(service.getServiceInstanceContent().getRootFolder());
         }
      }
      for (ManagedObjectReference childMor : parentFolder.getChildEntity()) {
         if (childMor.getType().equals(EntityType.FOLDER.getValue())) {
            Folder folder = service.getManagedObject(childMor);
            if (folder.getName().equals(child.getChildName())) {
               return folder._getRef();
            }
         }
      }
      return null;
   }

   /**
    * Gets the Managed Object Reference of a Cluster by searching for it in the
    * Child Entities of the parent Datacenter's host folder depending on its
    * type
    *
    * @param parent
    *            - the managed object that is parent of the cluster
    * @param child
    *            - the hierarchy item that corresponds to the cluster sought
    * @param service
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getClusterMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      Datacenter parentDC = (Datacenter) parent;
      Folder parentFolder = service.getManagedObject(parentDC.getHostFolder());
      for (ManagedObjectReference childMor : parentFolder.getChildEntity()) {
         try {
            if (childMor.getType().equals(EntityType.CLUSTER.getValue())) {
               ClusterComputeResource cluster = service.getManagedObject(childMor);
               if (cluster.getName().equals(child.getChildName())) {
                  return cluster._getRef();
               }
            }
         } catch (ManagedObjectNotFound monf) {
            _logger.error(monf.getMessage());
            continue;
         }
      }
      return null;
   }

   /**
    * Gets the Managed Object Reference of a Vapp by searching for it in the
    * Child Entities of the parent Datacenter's VM folder depending on its type
    *
    * @param parent
    *            - the managed object that is parent of the Vapp
    * @param child
    *            - the hierarchy item that corresponds to the vApp sought
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getVappMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      Folder parentFolder = null;
      if (parent instanceof Folder) {
         parentFolder = (Folder) parent;
      } else if (parent instanceof Datacenter) {
         parentFolder = service.getManagedObject(((Datacenter) parent).getVmFolder());
      } else if (parent instanceof VirtualApp) {
         for (EntityConfigInfo entityCfgInfo : ((VirtualApp) parent).getVAppConfig()
               .getEntityConfig()) {
            ManagedObject mo = service.getManagedObject(entityCfgInfo.getKey());
            if (mo instanceof VirtualApp) {
               if (((VirtualApp) mo).getName().equals(child.getChildName())) {
                  return mo._getRef();
               }
            }
         }
         return null;
      } else if ((parent instanceof HostSystem) || (parent instanceof ClusterComputeResource)
            || (parent instanceof ResourcePool)) {
         Datacenter parentDC = getParentDatacenter((ManagedEntity) parent, service);
         parentFolder = service.getManagedObject(parentDC.getVmFolder());
      } else {
         throw new RuntimeException(String.format(
               "Get vApp %s failed, because it's parent is not supported: %s",
               child.getChildName(),
               parent._getRef().getType()));
      }

      return getVappMorefFromNestedVmFolders(parentFolder, child.getChildName(), service);
   }

   /**
    * Gets the Managed Object Reference of a Host by searching for it in the
    * Child Entities of the Compute Resource that is its parent
    *
    * @param parent
    *            - the managed object that is parent of the host
    * @param child
    *            - the hierarchy item that corresponds to the host sought
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getHostMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      // The case of standalone host
      if (parent._getRef().getType().equals(DATACENTER_TYPE)) {
         Datacenter dc = (Datacenter) parent;
         Folder hostFolder = (Folder) service.getManagedObject(dc.getHostFolder());
         for (ManagedObjectReference moRefCr : hostFolder.getChildEntity()) {
            if (moRefCr.getType().equals(COMPUTE_RESOURCE_TYPE) || moRefCr.getType().equals(CLUSTER_TYPE)) {
               ComputeResource cr = service.getManagedObject(moRefCr);
               for (ManagedObjectReference moRefClHost : cr.getHost()) {
                  HostSystem host = service.getManagedObject(moRefClHost);
                  if (host.getName().equals(child.getChildName())) {
                     return host._getRef();

                  }
               }
            }
         }
      } else {
         // The case of clustered host
         ComputeResource parentCR = (ComputeResource) parent;
         for (ManagedObjectReference childMor : parentCR.getHost()) {
            HostSystem host = service.getManagedObject(childMor);
            if (host.getName().equals(child.getChildName())) {
               return host._getRef();
            }
         }
      }
      return null;
   }

   /**
    * Gets the Managed Object Reference of a resource pool by searching for it
    * in the child entities of its parent Compute Resource.
    *
    * Note: Parent of type cluster is only supported currently.
    *
    * @param parent
    *            - the managed object that is parent of the resource pool
    * @param child
    *            - the hierarchy item that corresponds to the resource pool
    *            sought
    * @return the Managed Object Reference for the found resource pool
    * @throws Exception
    */
   private static ManagedObjectReference getResourcePoolMoRef(
         ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      if (parent._getRef().getType().equals(CLUSTER_TYPE)) {
         // Resource pool with cluster parent
         ComputeResource clusterParentCR = (ComputeResource) parent;
         ResourcePool defaultClusterResPool = service
               .getManagedObject(clusterParentCR.getResourcePool());
         for (ManagedObjectReference childMor : defaultClusterResPool
               .getResourcePool()) {
            ResourcePool resPool = service.getManagedObject(childMor);
            if (resPool.getName().equals(child.getChildName())) {
               return resPool._getRef();
            }
         }
      } else if (parent._getRef().getType().equals(HOST_TYPE)) {
         // Resource pool with non-clustered host as parent
         Datacenter parentDC = getParentDatacenter((ManagedEntity) parent,
               service);
         Folder hostFolder = (Folder) service.getManagedObject(parentDC
               .getHostFolder());
         for (ManagedObjectReference childMor : hostFolder.getChildEntity()) {
            if (childMor.getType().equals(COMPUTE_RESOURCE_TYPE)) {
               ComputeResource hostCR = service.getManagedObject(childMor);
               ResourcePool defaultHostResPool = service
                     .getManagedObject(hostCR.getResourcePool());
               for (ManagedObjectReference resPoolMor : defaultHostResPool
                     .getResourcePool()) {
                  ResourcePool resPool = service.getManagedObject(resPoolMor);
                  if (resPool.getName().equals(child.getChildName())) {
                     return resPool._getRef();
                  }
               }
            }
         }
      } else {
         throw new RuntimeException(
               String.format(
                     "Get resource pool %s failed, because it's parent is not supported: %s",
                     child.getChildName(), parent._getRef().getType()));
      }
      return null;
   }

   /**
    * Retrieves the datacenter that holds the specified managed entity, e.g. a
    * HostSystem, VirtualApp, ResourcePool object.
    *
    * @param entity
    *            - the managed entity whose datacenter is looked for
    * @return the parent datacenter object
    * @throws Exception
    */
   private static Datacenter getParentDatacenter(ManagedEntity entity, VcService service) throws Exception {
      ManagedEntity entityParent = (ManagedEntity) service.getManagedObject(entity.getParent());
      ManagedObjectReference parent = entityParent.getParent();

      while (!parent.getType().equals(DATACENTER_TYPE)) {
         parent = ((ManagedEntity) service.getManagedObject(parent)).getParent();
      }
      return (Datacenter) service.getManagedObject(parent);
   }

   /**
    * Gets the Managed Object Reference of a distributed virtual switch as it
    * searches it in the child entities of the networking folder of the
    * datacenter.
    *
    * @param parent the datacenter managed object that is parent of the specified dvs
    * @param child the hierarchy item that corresponds to the DVS
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getDvsMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      Datacenter dc = (Datacenter) parent;
      Folder networkFolder = (Folder) service.getManagedObject(dc.getNetworkFolder());
      for (ManagedObjectReference moRefCr : networkFolder.getChildEntity()) {
         if (moRefCr.getType().equals(DVS_TYPE)) {
            DistributedVirtualSwitch dvs = service.getManagedObject(moRefCr);
            if (dvs.getName().equals(child.getChildName())) {
               return dvs._getRef();
            }
         }
      }

      return null;
   }

   /**
    * Gets the Managed Object Reference of a datastore cluster as it
    * searches it in the child entities of the datastore folder of the
    * datacenter.
    *
    * @param parent
    * @param child
    * @param service
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getDsClusterMoRef(ManagedObject parent, HierarchyItem child, VcService service)
         throws Exception {
      Datacenter dc = (Datacenter) parent;
      Folder dsFolder = (Folder) service.getManagedObject(dc.getDatastoreFolder());
      for (ManagedObjectReference moRefCr : dsFolder.getChildEntity()) {
         if (moRefCr.getType().equals(DS_CLUSTER_TYPE)) {
            StoragePod dsCluster = service.getManagedObject(moRefCr);
            if (dsCluster.getName().equals(child.getChildName())) {
               return dsCluster._getRef();
            }
         }
      }
      return null;
   }

   /**
    * Gets the Managed Object Reference of a network
    *
    * @param parent  the host managed object that is parent of the specified network
    * @param child   the hierarchy item that corresponds to the network
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getNetworkMoRef(
         ManagedObject parent, HierarchyItem child, VcService service) throws Exception {
      Datacenter dc = (Datacenter) parent;
      Folder networkFolder = (Folder) service.getManagedObject(dc.getNetworkFolder());
      for (ManagedObjectReference moRefCr : networkFolder.getChildEntity()) {
         if (moRefCr.getType().equals(NETWORK_TYPE)) {
            Network network = service.getManagedObject(moRefCr);
            if (network.getName().equals(child.getChildName())) {
               return network._getRef();
            }
         }
      }

      return null;
   }

   /**
    * Gets the Managed Object Reference of a distributed virtual portgroup
    * as it searches it in the child entities of the networking folder of the datacenter.
    *
    * @param parent the datacenter managed object that is parent of the specified dvs
    * @param child the hierarchy item that corresponds to the DVS
    * @return
    * @throws Exception
    */
   private static ManagedObjectReference getDvPortgroupMoRef(ManagedObject parent, HierarchyItem child,
         VcService service) throws Exception {
      DistributedVirtualSwitch dvs = (DistributedVirtualSwitch) parent;
      for (ManagedObjectReference moRefCr : dvs.getPortgroup()) {
         if (moRefCr.getType().equals(DVPORTGROUP_TYPE)) {
            DistributedVirtualPortgroup dvPg = service.getManagedObject(moRefCr);
            if (dvPg.getName().equals(child.getChildName())) {
               return dvPg._getRef();
            }
         }
      }

      return null;
   }
}