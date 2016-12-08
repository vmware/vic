/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutionException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.delay.Delay;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.pbm.ServiceInstance;
import com.vmware.vim.binding.pbm.ServiceInstanceContent;
import com.vmware.vim.binding.pbm.capability.CapabilityInstance;
import com.vmware.vim.binding.pbm.capability.CapabilityMetadata;
import com.vmware.vim.binding.pbm.capability.ConstraintInstance;
import com.vmware.vim.binding.pbm.capability.PropertyInstance;
import com.vmware.vim.binding.pbm.capability.types.BuiltinTypesEnum;
import com.vmware.vim.binding.pbm.capability.types.DiscreteSet;
import com.vmware.vim.binding.pbm.profile.CapabilityBasedProfile.ProfileCategoryEnum;
import com.vmware.vim.binding.pbm.profile.CapabilityBasedProfileCreateSpec;
import com.vmware.vim.binding.pbm.profile.CapabilityBasedProfileUpdateSpec;
import com.vmware.vim.binding.pbm.profile.Profile;
import com.vmware.vim.binding.pbm.profile.ProfileId;
import com.vmware.vim.binding.pbm.profile.ProfileManager;
import com.vmware.vim.binding.pbm.profile.ProfileOperationOutcome;
import com.vmware.vim.binding.pbm.profile.ResourceType;
import com.vmware.vim.binding.pbm.profile.SubProfileCapabilityConstraints;
import com.vmware.vim.binding.pbm.profile.SubProfileCapabilityConstraints.SubProfile;
import com.vmware.vim.binding.vim.ComputeResource;
import com.vmware.vim.binding.vim.ComputeResource.ConfigSpec;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.vmomi.client.Client;
import com.vmware.vim.vmomi.core.RequestContext;
import com.vmware.vim.vmomi.core.Stub;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vim.vmomi.core.impl.RequestContextImpl;
import com.vmware.vim.vmomi.core.types.VmodlType;
import com.vmware.vim.vmomi.core.types.VmodlTypeMap;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicyComponentSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicyRuleSetSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicyRuleSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.StoragePolicySpec;

/**
 * API commands for managing the storage policies.
 */
public class StoragePolicyBasicSrvApi {

   private static final Logger _logger = LoggerFactory
         .getLogger(StoragePolicyBasicSrvApi.class);
   private static final ResourceType STORAGE_TYPE = new ResourceType("STORAGE");
   private static final String COMPONENT_AS_CAPABILITY_ID = "com.vmware.storageprofile.dataservice";

   private static StoragePolicyBasicSrvApi instance = null;

   protected StoragePolicyBasicSrvApi() {
   }

   /**
    * Get instance of StoragePolicySrvApi.
    *
    * @return created instance
    */
   public static StoragePolicyBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (StoragePolicyBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing StoragePolicySrvApi.");
               instance = new StoragePolicyBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a storage policy that binds to specified tags. All tags has to be
    * assigned to datastores and they should be members of the passed tag
    * category.
    *
    * @param storagePolicySpec
    *           storage policy specification that will be used
    * @param categorySpec
    *           tags category, the tags should belong to this category
    * @param tagSpecs
    *           list of tag specs, all has to be applied to a datastores
    *
    * @return true if the creation was successful, false otherwise
    * @throws Exception
    */
   public boolean createStoragePolicy(StoragePolicySpec storagePolicySpec,
         BackingCategorySpec categorySpec, List<BackingTagSpec> tagSpecs)
         throws Exception {

      validateStoragePolicyCreateSpec(storagePolicySpec, categorySpec, tagSpecs);

      ProfileManager profileManager = getProfileManager(storagePolicySpec);
      CapabilityBasedProfileCreateSpec spec = createStorageProfileCreateSpec(
            storagePolicySpec, categorySpec, tagSpecs);

      BlockingFuture<ProfileId> result = new BlockingFuture<ProfileId>();
      profileManager.create(spec, result);

      return result.get() != null;
   }

   /**
    * Creates a storage policy component.
    *
    * @param storagePolicyComponentSpec
    *           the component to create
    * @return true if the component has been created successfully, false
    *         otherwise
    * @throws Exception
    */
   public boolean createStoragePolicyComponent(
         StoragePolicyComponentSpec storagePolicyComponentSpec)
         throws Exception {
      ProfileManager profileManager = getProfileManager(storagePolicyComponentSpec);
      CapabilityBasedProfileCreateSpec spec = toCapabilityBasedProfileCreateSpec(storagePolicyComponentSpec);
      BlockingFuture<ProfileId> result = new BlockingFuture<ProfileId>();
      profileManager.create(spec, result);
      return result.get() != null;
   }

   /**
    * Creates a storage policy.
    *
    * @param storagePolicySpec
    *           the policy to create
    * @return true if the policy has been created successfully, false otherwise
    * @throws Exception
    */
   public boolean createStoragePolicy(StoragePolicySpec storagePolicySpec)
         throws Exception {
      ProfileManager profileManager = getProfileManager(storagePolicySpec);
      CapabilityBasedProfileCreateSpec spec = toCapabilityBasedProfileCreateSpec(storagePolicySpec);
      BlockingFuture<ProfileId> result = new BlockingFuture<ProfileId>();
      profileManager.create(spec, result);
      return result.get() != null;
   }

   /**
    * Checks if the storage policy is present and if it is, then attempts to
    * delete it.
    *
    * @param storagePolicySpec
    * @return true if the policy has been deleted successfully or no such policy
    *         has been found, false if he storage policy could not be deleted
    * @throws VcException
    */
   public boolean deleteStoragePolicySafely(StoragePolicySpec storagePolicySpec)
         throws VcException {
      ProfileManager profileManager = getProfileManager(storagePolicySpec);
      ProfileId profileId = getProfileByName(storagePolicySpec.name.get(),
            profileManager);
      if (profileId != null) {
         return doDeleteProfile(profileId, profileManager);
      }
      return true;
   }

   /**
    * Checks if the input storage policy component is present and if it is, then
    * attempts to delete it.
    *
    * @param storagePolicyComponentSpec
    *           the component to delete
    * @return true if the component has been deleted successfully or no such
    *         component has been found, false if the component could not be
    *         deleted
    * @throws VcException
    */
   public boolean deleteStoragePolicyComponentSafely(
         StoragePolicyComponentSpec storagePolicyComponentSpec)
         throws VcException {
      ProfileManager profileManager = getProfileManager(storagePolicyComponentSpec);
      ProfileId profileId = getProfileByName(
            storagePolicyComponentSpec.name.get(), profileManager);
      if (profileId != null) {
         return doDeleteProfile(profileId, profileManager);
      }
      return true;
   }

   /**
    * Check if storage policy exists in the system.
    *
    * @param storagePolicySpec
    *           the spec of the policy that will be searched for
    * @return true if the policy is present, false otherwise
    */
   public boolean isStoragePolicyPresent(StoragePolicySpec storagePolicySpec) {
      validateStoragePolicySpec(storagePolicySpec);
      return isProfilePresent(storagePolicySpec);
   }

   /**
    * Checks if the input storage policy component exists
    *
    * @param storagePolicyComponentSpec
    *           the component to search for
    * @return true if the component exists, false otherwise
    */
   public boolean isStoragePolicyComponentPresent(
         StoragePolicyComponentSpec storagePolicyComponentSpec) {
      return isProfilePresent(storagePolicyComponentSpec);
   }

   /**
    * Method to enable/disable VM Storage Policies Profile on cluster or host
    * resource.
    *
    * @param entity
    *           cluster or host on which to enable/disable storage policy
    * @param enable
    *           enable or disable (true or false)
    * @return true if it was successful, false otherwise
    * @throws Exception
    *            if supplied resource is not host or cluster, or if problem with
    *            VC connection
    */
   public boolean setStoragePolicyEnabled(ManagedEntitySpec entity,
         boolean enable) throws Exception {
      ComputeResource parentComputeResource;

      // get the parent compute resource
      if (entity instanceof HostSpec) {
         HostSystem host = ManagedEntityUtil.getManagedObject(entity);
         parentComputeResource = VcServiceUtil.getVcService(entity)
               .getManagedObject(host.getParent());
      } else if (entity instanceof ClusterSpec) {
         parentComputeResource = ManagedEntityUtil.getManagedObject(entity);
      } else {
         throw new IllegalArgumentException(
               "Incompatible resource for storage policy assignment - it should be a cluster or host");
      }

      // if spbm status differs from applied, reconfigure
      if (parentComputeResource.getConfigurationEx().getSpbmEnabled() == null
            || parentComputeResource.getConfigurationEx().getSpbmEnabled()
                  .booleanValue() != enable) {
         ConfigSpec configSpec = new ConfigSpec();
         configSpec.setSpbmEnabled(enable);
         ManagedObjectReference taskMoRef = parentComputeResource
               .reconfigureEx(configSpec, true);
         return VcServiceUtil.waitForTaskSuccess(taskMoRef, entity);
      }

      return true;
   }

   /**
    * Deletes storage policy. The policy has to be present in the system.
    *
    * @param storagePolicySpec
    *           the spec of a storage policy which will be deleted
    * @return true if deletion was successful, false otherwise
    */
   public boolean deleteStoragePolicy(StoragePolicySpec storagePolicySpec)
         throws VcException {

      validateStoragePolicySpec(storagePolicySpec);

      ProfileManager profileManager = getProfileManager(storagePolicySpec);
      ProfileId profileId = getProfileByName(storagePolicySpec.name.get(),
            profileManager);

      BlockingFuture<ProfileOperationOutcome[]> result = new BlockingFuture<ProfileOperationOutcome[]>();
      profileManager.delete(new ProfileId[] { profileId }, result);

      try {
         ProfileOperationOutcome[] profileOperationOutcomes = result.get();
         if (profileOperationOutcomes != null) {
            for (ProfileOperationOutcome outcome : profileOperationOutcomes) {
               if (outcome.fault != null) {
                  _logger.error("Error deleting storage profile: "
                        + outcome.fault);
                  return false;
               }
            }
         }
      } catch (ExecutionException e) {
         _logger.error("Error deleting storage policy!");
         e.printStackTrace();
         return false;
      } catch (InterruptedException e) {
         _logger.error("Error deleting storage policy!");
         e.printStackTrace();
         return false;
      }

      return true;
   }

   /**
    * Update the name of the target storage policy.
    *
    * @param targetStoragePolicy
    *           storage policy which name will be changed
    * @param newStoragePolicy
    *           spec that holds the new storage policy name
    * @return if the update operation was successful
    * @throws VcException
    *            if there is VC connection problem
    */
   public boolean updateStoragePolicyName(
         StoragePolicySpec targetStoragePolicy,
         StoragePolicySpec newStoragePolicy) throws VcException {
      validateStoragePolicySpec(targetStoragePolicy);
      validateStoragePolicySpec(newStoragePolicy);

      ProfileManager profileManager = getProfileManager(targetStoragePolicy);
      ProfileId profileId = getProfileByName(targetStoragePolicy.name.get(),
            profileManager);

      CapabilityBasedProfileUpdateSpec updateSpec = new CapabilityBasedProfileUpdateSpec(
            newStoragePolicy.name.get(), null, null);
      BlockingFuture<Void> result = new BlockingFuture<Void>();
      profileManager.update(profileId, updateSpec, result);

      try {
         result.get();
      } catch (ExecutionException e) {
         _logger.error("Error updating storage policy!");
         e.printStackTrace();
         return false;
      } catch (InterruptedException e) {
         _logger.error("Error updating storage policy!");
         e.printStackTrace();
         return false;
      }

      return true;
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Method that validates the storage policy specification: - storage policy
    * name should be assigned - at least one tag assigned to a datastore should
    * be available - tag category name should be assigned
    *
    * @param storagePolicySpec
    *           storage policy spec that will be validated
    * @param categorySpec
    *           tag category, which contains all the tags
    * @param gatSpecs
    *           list of tag specs, that should be assigned to datastores
    * @throws IllegalArgumentException
    *            if storage policies spec requirements aren't met
    */
   private void validateStoragePolicyCreateSpec(
         StoragePolicySpec storagePolicySpec, BackingCategorySpec categorySpec,
         List<BackingTagSpec> tagSpecs) throws IllegalArgumentException {
      validateStoragePolicySpec(storagePolicySpec);

      if (!categorySpec.name.isAssigned() || categorySpec.name.get().isEmpty()) {
         _logger.error("Tag category name is not set.");
         throw new IllegalArgumentException("Tag category name is not set.");
      }

      if (tagSpecs == null || tagSpecs.isEmpty()) {
         _logger.error("No Storage policies tags are avaibale.");
         throw new IllegalArgumentException(
               "No Storage policies tags are avaibale.");
      }
   }

   /**
    * Method that validates the storage policy specification: - storage policy
    * name should be assigned
    *
    * @param storagePolicySpec
    *           storage policy spec that will be validated
    * @throws IllegalArgumentException
    *            if storage policies spec requirements aren't met
    */
   private void validateStoragePolicySpec(StoragePolicySpec storagePolicySpec)
         throws IllegalArgumentException {
      if (!storagePolicySpec.name.isAssigned()
            || storagePolicySpec.name.get().isEmpty()) {
         _logger.error("Storage policy name is not set.");
         throw new IllegalArgumentException("Storage policy name is not set.");
      }
   }

   /**
    * Obtains PBM service instance content.
    *
    * @param pbmService
    *           PBM service, whose content will be retrieved
    * @return ` the obtained service content
    */
   private ServiceInstanceContent getServiceInstanceContent(
         ServiceInstance pbmService) {

      ServiceInstanceContent content = null;
      BlockingFuture<ServiceInstanceContent> future = new BlockingFuture<ServiceInstanceContent>();
      pbmService.getContent(future);
      try {
         content = future.get();
      } catch (Exception e) {
         _logger.error("Failed to obtain PBM service content!");
         throw new RuntimeException(e);
      }
      return content;
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
   private ServiceInstance getPbmService(Client client, ServiceSpec serviceSpec)
         throws VcException {
      return createStub("PbmServiceInstance", "ServiceInstance", client,
            serviceSpec);
   }

   /**
    * Obtains profile manager.
    *
    * @return initialized profile manager
    * @throws VcException
    *            if connection to VC does not succeeds
    */
   ProfileManager getProfileManager(EntitySpec entitySpec) throws VcException {
      Client client = VcServiceUtil.getPbmVsliClient(entitySpec);
      ServiceInstanceContent content = getServiceInstanceContent(getPbmService(
            client, entitySpec.service.get()));
      return getManagedObject(content.getProfileManager(), client,
            entitySpec.service.get());
   }

   @SuppressWarnings("unchecked")
   private <T extends ManagedObject> T createStub(String moRefType,
         String moRefId, Client client, ServiceSpec serviceSpec)
         throws VcException {
      ManagedObjectReference moRef = new ManagedObjectReference(moRefType,
            moRefId);
      return (T) getManagedObject(moRef, client, serviceSpec);
   }

   private <T extends ManagedObject> T getManagedObject(
         ManagedObjectReference moRef, Client client, ServiceSpec serviceSpec)
         throws VcException {

      RequestContext sessionContext = new RequestContextImpl();
      sessionContext.put("vcSessionCookie",
            VcServiceUtil.getSessionCookie(serviceSpec));

      VmodlTypeMap typeMap = VmodlTypeMap.Factory.getTypeMap();
      VmodlType vmodlType = typeMap.getVmodlType(moRef.getType());
      Class<T> typeClass = (Class<T>) vmodlType.getTypeClass();

      T result = client.createStub(typeClass, moRef);
      // It is required that we pass the VC session cookie , otherwise
      // the service will complain about an invalid session.
      ((Stub) result)._setRequestContext(sessionContext);
      return result;
   }

   private CapabilityBasedProfileCreateSpec createStorageProfileCreateSpec(
         StoragePolicySpec storagePolicySpec, BackingCategorySpec categorySpec,
         List<BackingTagSpec> tagSpecs) {

      List<String> tagNames = new ArrayList<String>();
      for (BackingTagSpec tag : tagSpecs) {
         tagNames.add(tag.name.get());
      }

      DiscreteSet discretSet = new DiscreteSet();
      discretSet.setValues(tagNames.toArray(new Object[0]));

      List<PropertyInstance> propertyInstances = new ArrayList<PropertyInstance>();

      PropertyInstance propInstance = new PropertyInstance();
      propInstance.id = String.format("com.vmware.storage.tag.%s.property",
            categorySpec.name.get());
      propInstance.value = discretSet;

      propertyInstances.add(propInstance);

      List<ConstraintInstance> constraintInstances = new ArrayList<ConstraintInstance>();
      constraintInstances.add(new ConstraintInstance(propertyInstances
            .toArray(new PropertyInstance[0])));

      CapabilityInstance capabilityInstance = new CapabilityInstance();
      capabilityInstance.setId(new CapabilityMetadata.UniqueId(
            "http://www.vmware.com/storage/tag", categorySpec.name.get()));
      capabilityInstance.setConstraint(constraintInstances
            .toArray(new ConstraintInstance[0]));

      SubProfile subProfile = new SubProfileCapabilityConstraints.SubProfile(
            "subProfile", new CapabilityInstance[] { capabilityInstance },
            Boolean.TRUE);

      SubProfileCapabilityConstraints constraints = new SubProfileCapabilityConstraints(
            new SubProfileCapabilityConstraints.SubProfile[] { subProfile });

      CapabilityBasedProfileCreateSpec spec = new CapabilityBasedProfileCreateSpec();
      spec.setName(storagePolicySpec.name.get());
      if (storagePolicySpec.description.isAssigned()) {
         spec.setDescription(storagePolicySpec.description.get());
      } else {
         // workaround for required description of the policy
         spec.setDescription(storagePolicySpec.name.get());
      }

      spec.setResourceType(STORAGE_TYPE);
      spec.setConstraints(constraints);

      return spec;
   }

   /**
    * Transforms {@link StoragePolicyComponentSpec} into {@link SubProfile}.
    *
    * @param storagePolicyComponentSpec
    *           the component spec to transform
    * @return sub-profile with the name and capabilities of the input component
    *         spec
    * @throws IllegalArgumentException
    */
   public SubProfile toSubProfile(
         StoragePolicyComponentSpec storagePolicyComponentSpec)
         throws IllegalArgumentException {
      if (!storagePolicyComponentSpec.name.isAssigned()) {
         throw new IllegalArgumentException(
               "Invalid spec - component has no name.");
      }
      CapabilityInstance[] capabilityInstances = new CapabilityInstance[storagePolicyComponentSpec.rules
            .getAll().size()];
      for (int i = 0; i < capabilityInstances.length; i++) {
         capabilityInstances[i] = storagePolicyComponentSpec.rules.getAll()
               .get(i).asCapabilityInstance();
      }
      SubProfile subProfile = new SubProfile();
      subProfile.setName(storagePolicyComponentSpec.name.get());
      subProfile.setCapability(capabilityInstances);
      return subProfile;
   }

   /**
    * Transforms {@link StoragePolicyComponentSpec} into list of
    * {@link CapabilityInstance}s. Each {@link StoragePolicyRuleSpec} is
    * transformed into {@link CapabilityInstance}.
    *
    * @param storagePolicyComponentSpec
    *           the component spec to transform
    * @return list of capability instances
    */
   public List<CapabilityInstance> extractCapabilityInstances(
         StoragePolicyComponentSpec storagePolicyComponentSpec) {
      if (storagePolicyComponentSpec == null) {
         throw new IllegalArgumentException(
               "Storage policy component cannot be null.");
      }
      List<CapabilityInstance> capabilityInstances = new ArrayList<>();
      for (StoragePolicyRuleSpec ruleSpec : storagePolicyComponentSpec.rules
            .getAll()) {
         capabilityInstances.add(ruleSpec.asCapabilityInstance());
      }
      return capabilityInstances;
   }

   /**
    * Transforms {@link ProfileId} into {@link CapabilityInstance} for the cases
    * when an existing component (represented as a storage policy) must be
    * referenced in a storage policy.
    *
    * @param componentPolicyId
    *           id of the existing component
    * @return capability instance referencing the existing component
    */
   public CapabilityInstance createComponentPolicyReference(
         ProfileId componentPolicyId) {
      CapabilityInstance capabilityInstance = new CapabilityInstance();
      capabilityInstance.setId(new CapabilityMetadata.UniqueId(
            COMPONENT_AS_CAPABILITY_ID, componentPolicyId.getUniqueId()));
      PropertyInstance[] properties = { new PropertyInstance(
            componentPolicyId.getUniqueId(), null,
            componentPolicyId.getUniqueId()) };
      ConstraintInstance constraint = new ConstraintInstance(properties);
      capabilityInstance.setConstraint(new ConstraintInstance[] { constraint });
      return capabilityInstance;
   }

   /**
    * Transforms {@link StoragePolicyRuleSetSpec} into {@link SubProfile}.<br />
    * If the component does not exist, rule specs from the rule-set spec are
    * merged with the rule specs from the component specs, contained in the
    * rule-set spec. All the rule specs are then transformed into capability
    * instances and added to the sub-profile.<br />
    * If the component exists, then its id is packaged in a property instance
    * with a data type of {@link BuiltinTypesEnum#VMW_POLICY}. The property
    * instance is in turn assigned to a capability instance from the
    * "com.vmware.storageprofile.dataservice" namespace, which is added to the
    * sub-profile.
    *
    * @param storagePolicyRuleSetSpec
    *           the rule-set spec to transform
    * @return sub-profile with the name and capabilities of the input rule-set
    *         spec
    * @throws IllegalArgumentException
    * @throws VcException
    */
   public SubProfile toSubProfile(
         StoragePolicyRuleSetSpec storagePolicyRuleSetSpec)
         throws IllegalArgumentException, VcException {
      List<CapabilityInstance> capabilityInstances = new ArrayList<>();
      for (StoragePolicyRuleSpec ruleSpec : storagePolicyRuleSetSpec.rules
            .getAll()) {
         capabilityInstances.add(ruleSpec.asCapabilityInstance());
      }
      for (StoragePolicyComponentSpec componentSpec : storagePolicyRuleSetSpec.components
            .getAll()) {
         final ProfileManager profileManager = getProfileManager(componentSpec);
         if (!componentSpec.name.isAssigned()) {
            throw new IllegalArgumentException(
                  "Invalid spec - component has no name.");
         }
         final ProfileId profileId = getProfileByName(componentSpec.name.get(),
               profileManager);
         if (profileId == null) {
            // No storage policy representing this component, so flatten all the
            // component's capabilities into the current SubProfile.
            capabilityInstances
                  .addAll(extractCapabilityInstances(componentSpec));
         } else {
            // Storage policy for this component is existing, so reference the
            // policy id in a capability instance.
            capabilityInstances.add(createComponentPolicyReference(profileId));
         }
      }
      CapabilityInstance[] capabilityInstancesArray = capabilityInstances
            .toArray(new CapabilityInstance[capabilityInstances.size()]);
      SubProfile subProfile = new SubProfile();
      subProfile.setName(storagePolicyRuleSetSpec.name.get());
      subProfile.setCapability(capabilityInstancesArray);
      return subProfile;
   }

   /**
    * Transforms {@link StoragePolicySpec} into
    * {@link CapabilityBasedProfileCreateSpec} with a <i>requirement</i>
    * category.
    *
    * @param storagePolicySpec
    *           the policy spec to transform
    * @return the corresponding create spec
    * @throws Exception
    */
   @SuppressWarnings("deprecation")
   public CapabilityBasedProfileCreateSpec toCapabilityBasedProfileCreateSpec(
         StoragePolicySpec storagePolicySpec) throws Exception {
      SubProfile[] subProfilesArray = new SubProfile[storagePolicySpec.rulesets
            .getAll().size()];
      for (int i = 0; i < subProfilesArray.length; i++) {
         subProfilesArray[i] = toSubProfile(storagePolicySpec.rulesets.getAll()
               .get(i));
      }
      SubProfileCapabilityConstraints constraints = new SubProfileCapabilityConstraints(
            subProfilesArray);
      CapabilityBasedProfileCreateSpec capabilityBasedProfileCreateSpec = new CapabilityBasedProfileCreateSpec();
      capabilityBasedProfileCreateSpec.setResourceType(STORAGE_TYPE);
      capabilityBasedProfileCreateSpec.setName(storagePolicySpec.name.get());
      capabilityBasedProfileCreateSpec
            .setCategory(ProfileCategoryEnum.REQUIREMENT.toString());
      capabilityBasedProfileCreateSpec.setConstraints(constraints);
      if (storagePolicySpec.description.isAssigned()) {
         capabilityBasedProfileCreateSpec
               .setDescription(storagePolicySpec.description.get());
      }
      return capabilityBasedProfileCreateSpec;
   }

   /**
    * Transforms {@link StoragePolicyComponentSpec} into
    * {@link CapabilityBasedProfileCreateSpec} with a <i>data service</i>
    * category.
    *
    * @param storagePolicyComponentSpec
    *           the component spec to transform
    * @return the corresponding create spec
    * @throws Exception
    */
   @SuppressWarnings("deprecation")
   public CapabilityBasedProfileCreateSpec toCapabilityBasedProfileCreateSpec(
         StoragePolicyComponentSpec storagePolicyComponentSpec)
         throws Exception {
      if (!storagePolicyComponentSpec.rules.isAssigned()
            || storagePolicyComponentSpec.rules.getAll().size() == 0) {
         throw new IllegalArgumentException(
               "There are no rules in this component.");
      }
      SubProfile subProfile = toSubProfile(storagePolicyComponentSpec);
      SubProfileCapabilityConstraints constraints = new SubProfileCapabilityConstraints(
            new SubProfile[] { subProfile });
      CapabilityBasedProfileCreateSpec capabilityBasedProfileCreateSpec = new CapabilityBasedProfileCreateSpec();
      capabilityBasedProfileCreateSpec.setResourceType(STORAGE_TYPE);
      capabilityBasedProfileCreateSpec.setName(storagePolicyComponentSpec.name
            .get());
      if (storagePolicyComponentSpec.description.isAssigned()) {
         capabilityBasedProfileCreateSpec
               .setDescription(storagePolicyComponentSpec.description.get());
      } else {
         // When creating storage policy component as a policy, description is mandatory,
         // as PbmExtendedElementDescription throws a serialization error when policies are
         // requested.
         capabilityBasedProfileCreateSpec.setDescription("");
      }
      capabilityBasedProfileCreateSpec
            .setCategory(ProfileCategoryEnum.DATA_SERVICE_POLICY.toString());
      capabilityBasedProfileCreateSpec.setConstraints(constraints);
      return capabilityBasedProfileCreateSpec;
   }

   /**
    * Gets storage policy profile by name.
    *
    * @param profileName
    *           the name of the storage policy
    * @param profileManager
    *           the profile manager
    * @return profile ID
    */
   ProfileId getProfileByName(String profileName, ProfileManager profileManager) {
      ProfileId profileId = null;

      boolean finished = false;
      long timeout = Delay.timeout.forSeconds(10).getDuration();
      long elapsedTime = 0;
      long startTime = System.currentTimeMillis();

      // In profiles array sometimes has more data because delete action is not
      // finished. Check if profileName is found in while loop for 10 seconds
      while (!finished && elapsedTime < timeout) {
         try {

            BlockingFuture<ProfileId[]> blockingResult = new BlockingFuture<ProfileId[]>();
            profileManager.queryProfile(STORAGE_TYPE, null, blockingResult);

            try {
               ProfileId[] profileIds = blockingResult.get();
               BlockingFuture<Profile[]> profileResult = new BlockingFuture<Profile[]>();
               profileManager.retrieveContent(profileIds, profileResult);

               Profile[] profiles = profileResult.get();
               for (Profile profile : profiles) {
                  if (profile.getName().equals(profileName)) {
                     profileId = profile.getProfileId();
                  }
               }
            } catch (ExecutionException | InterruptedException e) {
               throw new RuntimeException(
                     "Could not retrieve storage profile with name: "
                           + profileName, e);
            }
            finished = true;
         } catch (RuntimeException e) {
            _logger
                  .error("Error getting profile by name: in profiles array has more data because of performing action");
            if (elapsedTime >= timeout) {
               // It wasn't successful after 10 seconds
               throw new RuntimeException(
                     "Could not get storage profile by name ", e);
            }
         }
         elapsedTime = System.currentTimeMillis() - startTime;
      }
      return profileId;
   }

   /**
    * Checks for <i>requirement</i> or <i>data service</i> profiles existence,
    * i.e. storage policy or storage policy component.
    *
    * @param spec
    *           <i>requirement</i> or <i>data service</i> profile
    * @return true if the profile is present, false if it is not or an error has
    *         occurred
    */
   private boolean isProfilePresent(ManagedEntitySpec spec) {
      if (!(spec instanceof StoragePolicySpec)
            && !(spec instanceof StoragePolicyComponentSpec)) {
         return false;
      }
      ProfileManager profileManager;
      try {
         profileManager = getProfileManager(spec);
      } catch (VcException e) {
         return false;
      }
      return getProfileByName(spec.name.get(), profileManager) != null;
   }

   /**
    * Deletes <i>requirement</i> or <i>data service</i> profiles.
    *
    * @param profileId
    *           the id of the profile to delete
    * @param profileManager
    *           the profile manager to use for deletion
    * @return true if the profile has been deleted successfully, false otherwise
    */
   private boolean doDeleteProfile(ProfileId profileId,
         ProfileManager profileManager) {
      BlockingFuture<ProfileOperationOutcome[]> result = new BlockingFuture<ProfileOperationOutcome[]>();
      profileManager.delete(new ProfileId[] { profileId }, result);
      try {
         ProfileOperationOutcome[] profileOperationOutcomes = result.get();
         if (profileOperationOutcomes != null) {
            for (ProfileOperationOutcome outcome : profileOperationOutcomes) {
               if (outcome.fault != null) {
                  _logger.error("Error deleting storage profile: "
                        + outcome.fault);
                  return false;
               }
            }
         }
      } catch (ExecutionException | InterruptedException e) {
         _logger.error("Error deleting storage profile!");
         e.printStackTrace();
         return false;
      }
      return true;
   }
}
