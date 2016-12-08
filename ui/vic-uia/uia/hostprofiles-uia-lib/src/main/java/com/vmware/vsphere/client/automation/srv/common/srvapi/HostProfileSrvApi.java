/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.ServiceInstanceContent;
import com.vmware.vim.binding.vim.profile.ComplianceManager;
import com.vmware.vim.binding.vim.profile.ProfileManager;
import com.vmware.vim.binding.vim.profile.host.HostProfile;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.FileUtils;
import com.vmware.vsphere.client.automation.common.HostProfilesUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostProfileSpec;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * API server wrapper for host profile related calls.
 */
public class HostProfileSrvApi {
    private static final Logger _logger = LoggerFactory.getLogger(HostProfileSrvApi.class);
   private static final String HOST_SYSTEM = "HostSystem";

    private static HostProfileSrvApi INSTANCE = null;

    protected HostProfileSrvApi() {
        // Disallow the creation of this object outside of the class.
    }

    /**
     * Get instance of HostProfileSrvApi.
     *
     * @return created instance
     */
    public static HostProfileSrvApi getInstance() {
        if (INSTANCE == null) {
            synchronized (HostProfileSrvApi.class) {
                if (INSTANCE == null) {
                    _logger.info("Initializing HostProfileSrvApi.");
                    INSTANCE = new HostProfileSrvApi();
                }
            }
        }

        return INSTANCE;
    }

    /**
     * Returns an array of the hosts which are attached to the host profile
     *
     * @param hpSpec host profile spec
     * @return an array of host MORs
     */
    public ManagedObjectReference[] getAttachedHosts(HostProfileSpec hpSpec) {
        ManagedObjectReference[] attachedEntities = getAttachedEntities(hpSpec);
        ManagedObjectReference[] hostList = filterOnlyHosts(attachedEntities);

        return hostList;
    }

    /**
     * Method that deletes a specified Host Profile.
     *
     * @param hostProfileSpec - the specification of the host profile to be deleted
     * @throws Exception if the profile has not been found
     */
    public void deleteHostProfile(final HostProfileSpec hostProfileSpec) throws Exception {
        HostProfile target = getHostProfileByName(hostProfileSpec);

       HostProfilesUtil.ensureNotNull(target, String.format("Host profile " +
          "'%s' has not been found.", hostProfileSpec.name.get()));
        _logger.info(String.format("Deleting host profile '%s'", hostProfileSpec.name.get()));

        // Delete the profile.
        target.destroy();
    }

    /**
     * Method that checks the existence of a Host Profile
     *
     * @param hostProfileSpec - the specification of the host profile
     * @returns true if host profile exists, otherwise false
     * @throws Exception
     */
    public boolean checkHostProfileExists(final HostProfileSpec hostProfileSpec) throws Exception {
        if (getHostProfileByName(hostProfileSpec) == null) {
            return false;
        }

        return true;
    }

    /**
     * Method that creates a host profile
     *
     * @param hostProfileSpec - the specification of the profile to be created, it contains
     *            the reference host from which to extract the profile
     * @throws Exception - throws exception if the profile cannot be created
     */
    public void createHostProfileFromHost(final HostProfileSpec hostProfileSpec) throws Exception {
        ManagedObjectReference hostMor = ManagedEntityUtil.getManagedObject(hostProfileSpec.referenceHost.get())
            ._getRef();

        HostProfile.HostBasedConfigSpec createSpec = new HostProfile.HostBasedConfigSpec();
        if (hostProfileSpec.description.isAssigned()) {
            createSpec.setAnnotation(hostProfileSpec.description.get());
        } else {
            createSpec.setAnnotation(null);
        }
        createSpec.setEnabled(true);
        createSpec.setName(hostProfileSpec.name.get());
        createSpec.setProfilesToExtract(null);
        createSpec.setHost(hostMor);
        createSpec.setUseHostProfileEngine(true);

        ProfileManager profileManager = getHostProfileManager(hostProfileSpec.service.get());
        ManagedObjectReference newProfileRef = profileManager.createProfile(createSpec);

        final HostProfile newProfile = ManagedEntityUtil
            .getManagedObjectFromMoRef(newProfileRef, hostProfileSpec.service.get());
        newProfile.updateReferenceHost(hostMor);
    }

    /**
     * Method that associates a host or a cluster to a specified host profile
     *
     * @param hostProfileSpec - specification of the host profile
     * @param entities - specification of the entities to attach
     * @return true if successful, false if host profile or host do not exist
     * @throws Exception - if there is a vc exception
     */
    public boolean attachHostProfile(final HostProfileSpec hostProfileSpec, final ManagedEntitySpec... entities)
        throws Exception {
        TestSpecValidator.ensureNotNull(hostProfileSpec, "No HostProfileSpec!");
        TestSpecValidator.ensureNotEmpty(Arrays.asList(entities), "Supply entities to attach!");

        HostProfile hostProfile = getHostProfileByName(hostProfileSpec);
        if (hostProfile == null) {
            return false;
        }

        try {
            ManagedObjectReference[] mors = getMors(entities);
            hostProfile.associateEntities(mors);
        } catch (ObjectNotFoundException e) {
            return false;
        }

        return true;
    }

    /**
     * Method that dissociates a host to a specified host profile
     *
     * @param hostProfileSpec - specification of the host profile
     * @param entities - specification of the entities to detach
     * @return true if successful, false if host profile or host do not exist
     * @throws Exception - if there is a vc exception
     */
    public boolean dettachHostProfile(final HostProfileSpec hostProfileSpec, final ManagedEntitySpec... entities)
        throws Exception {
        TestSpecValidator.ensureNotNull(hostProfileSpec, "Supply HostProfielSpec!");
        TestSpecValidator.ensureNotEmpty(Arrays.asList(entities), "Supply entities to detach!");

        HostProfile hostProfile = getHostProfileByName(hostProfileSpec);
        if (hostProfile == null) {
            return false;
        }

        try {
            ManagedObjectReference[] mors = getMors(entities);
            hostProfile.dissociateEntities(mors);
        } catch (ObjectNotFoundException e) {
            return false;
        }
        return true;
    }

   /**
    * Exports a Host Profile to a .vpf file.
    *
    * @param profileSpec the host profile spec to export from
    * @param fileName    the file name to export to
    */
   public void exportHostProfile(HostProfileSpec profileSpec, String fileName) {
      TestSpecValidator.ensureNotNull(profileSpec, "HostProfileSpec is null");
      verifyNotNullOrEmpty(fileName, "Filename is null or empty");

      HostProfile hostProfile = getHostProfile(profileSpec);
      String exportedProfile = hostProfile.exportProfile();
      FileUtils.getInstance().writeStringToTextFile(exportedProfile, fileName);
   }

   private void verifyNotNullOrEmpty(String string, String errorMessage) {
      if (Strings.isNullOrEmpty(string)) {
         throw new IllegalArgumentException(errorMessage);
      }
   }

   private HostProfile getHostProfile(HostProfileSpec profileSpec) {
      HostProfile hostProfile;
      try {
         hostProfile = getHostProfileByName(profileSpec);
      } catch (Exception e) {
         throw new RuntimeException("Error getting host profile", e);
      }
      return hostProfile;
   }

   /**
     * Method that checks compliance for specified entity
     *
     * @param entity
     * @return true if successful, false if no such profile or task fails
     * @throws Exception - if no such entity
     */
    public boolean checkHostProfileCompliance(final ManagedEntitySpec entity) throws Exception {
        TestSpecValidator.ensureNotNull(entity, "Supply entity to check compliance!");

        VcService vcService = VcServiceUtil.getVcService(entity.service.get());
        ServiceInstanceContent serviceContent = vcService.getServiceInstanceContent();

       HostProfilesUtil.ensureNotNull(serviceContent, "ServiceInstanceContent was found to be null");

        ManagedObjectReference entityMor = getMors(entity)[0];

        // this manager is used to invoke check compliance
        ManagedObjectReference complManagerMor = serviceContent.getComplianceManager();
        ComplianceManager complMgr = vcService.getManagedObject(complManagerMor);

        // get the profile associated to the entity which has to be gotten
        ProfileManager profileManager = getHostProfileManager(entity.service.get());
        ManagedObjectReference[] assocToEntityProfiles = profileManager.findAssociatedProfile(entityMor);
       HostProfilesUtil.ensureNotNull(assocToEntityProfiles, "The entity: " +
          entity.name.get() + " has no associated host profile!");

        ManagedObjectReference task = complMgr
            .checkCompliance(assocToEntityProfiles, new ManagedObjectReference[] { entityMor });
        // success or failure of the task
        return VcServiceUtil.waitForTaskSuccess(task, entity.service.get());
    }

    // Private methods

    /**
     * Method that finds a specified by its HostProfileSpec host profile.
     *
     * @param hostProfileSpec - specification of the host profile to look for
     * @return The HostProfile object that represents the host profile from the spec, or
     *         null if it doesn't exist
     * @throws Exception if the HostProfileManager cannot be gotten or the ManagedObject of
     *             the host profile
     */
    HostProfile getHostProfileByName(final HostProfileSpec hostProfileSpec) throws Exception {
        HostProfile targetProfile = null;
        ProfileManager profileManager = getHostProfileManager(hostProfileSpec.service.get());
        ManagedObjectReference[] allProfiles = profileManager.getProfile();

        for (ManagedObjectReference hostProfileMor : allProfiles) {
            HostProfile hostProfile = ManagedEntityUtil
                .getManagedObjectFromMoRef(hostProfileMor, hostProfileSpec.service.get());
            if (hostProfile.getName().equals(hostProfileSpec.name.get())) {
                targetProfile = hostProfile;
                break;
            }
        }

        return targetProfile;
    }

    /**
     * Retrieves the ProfileManager.
     *
     * @param serviceSpec
     * @return
     * @throws Exception
     */
    private ProfileManager getHostProfileManager(final ServiceSpec serviceSpec) throws Exception {
        ProfileManager profileManager = null;

        VcService vcService = VcServiceUtil.getVcService(serviceSpec);
        ServiceInstanceContent serviceContent = vcService.getServiceInstanceContent();

        HostProfilesUtil.ensureNotNull(serviceContent, "ServiceInstanceContent was found to be null");

        ManagedObjectReference hostProfileManagerMor = serviceContent.getHostProfileManager();
        profileManager = ManagedEntityUtil.getManagedObjectFromMoRef(hostProfileManagerMor, serviceSpec);

        return profileManager;
    }

    /**
     * Returns a list of the entities attached to the host profile
     *
     * @param hostProfileSpec
     * @return list of attached MORs
     */
    private ManagedObjectReference[] getAttachedEntities(HostProfileSpec hostProfileSpec) {
        HostProfile hostProfile;
        try {
            hostProfile = getHostProfileByName(hostProfileSpec);
        } catch (Exception e) {
            throw new RuntimeException("Error getting host profile");
        }

        ManagedObjectReference[] entities = hostProfile.getEntity();

        return entities;
    }

    /**
     * Returns a list containing only the MORs which are of type HostSystem
     *
     * @param entities
     * @return list of host MORs
     */
    private ManagedObjectReference[] filterOnlyHosts(ManagedObjectReference[] entities) {
        List<ManagedObjectReference> hosts = new ArrayList<>();

        for (ManagedObjectReference attachedEntity : entities) {
            if (attachedEntity.getType().equals(HOST_SYSTEM)) {
                hosts.add(attachedEntity);
            }
        }

        // transform the List to an array
        ManagedObjectReference[] hostList;
        hostList = new ManagedObjectReference[hosts.size()];
        hosts.toArray(hostList);

        return hostList;
    }

   /**
    * Method that creates an array of references to managed entities, i.e. clusters and
    * hosts
    *
    * @param entities - host profile spec should contain a list of hosts
    * @return array of host references
    * @throws Exception ObjectNotFound exception if a host doesn't exist in Inventory
    */
   public ManagedObjectReference[] getMors(final ManagedEntitySpec... entities) throws Exception {
      HostProfilesUtil.ensureNotNull(entities, "No entities supplied");

      int size = entities.length;
      ManagedObjectReference[] hostMors = new ManagedObjectReference[size];
      for (int i = 0; i < size; i++) {
         ManagedObject host = ManagedEntityUtil.getManagedObject(entities[i]);
         hostMors[i] = host._getRef();
      }

      return hostMors;
   }
}
