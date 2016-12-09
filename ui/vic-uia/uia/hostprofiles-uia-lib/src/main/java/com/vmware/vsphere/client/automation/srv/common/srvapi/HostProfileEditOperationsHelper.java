/*
 * Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.vim.binding.impl.vmodl.KeyAnyValueImpl;
import com.vmware.vim.binding.vim.PasswordField;
import com.vmware.vim.binding.vim.fault.DuplicateName;
import com.vmware.vim.binding.vim.fault.ProfileUpdateFailed;
import com.vmware.vim.binding.vim.profile.*;
import com.vmware.vim.binding.vim.profile.host.*;
import com.vmware.vim.binding.vmodl.KeyAnyValue;
import com.vmware.vsphere.client.automation.common.HostProfilesUtil;
import com.vmware.vsphere.client.automation.common.spec.HostPolicyUpdateSpec;
import com.vmware.vsphere.client.automation.common.spec.PolicyOptionType;
import com.vmware.vsphere.client.automation.srv.common.spec.HostProfileSpec;
import org.apache.commons.lang.ArrayUtils;

/**
 * Helper class that encapsulates all specific edit operations needed by HostProfileEditOperationsSrvApi class.
 */
class HostProfileEditOperationsHelper {
    private static final String KEY = "key";
    private static final String NTP_CLIENT = "key-vim-profile-host-RulesetProfile-ntpClient";
    private static final String CFG_ETC_MOTD = "key-vim-profile-host-OptionProfile-Config_Etc_motd";
    private static final String PARAM_NAME = "parameterName";
    private static final String PARAM_VALUE = "parameterValue";
    private static final String GENERIC_NETSTACK_INSTANCE_PROFILE = "GenericNetStackInstanceProfile";
    private static final String GENERIC_DNS_CFG_PROFILE = "GenericDnsConfigProfile";
    private static final String DNS_CFG_POLICY = "DnsConfigPolicy";
    private static final String DHCP = "dhcp";
    private static final String MGMT_NW_HOST_PORT_GROUP = "key-vim-profile-host-HostPortgroupProfile-ManagementNetwork";
    private static final String ADMIN_PWD_POLICY = "AdminPasswordPolicy";
    private static final String SECURITY_PROFILE =  "security_SecurityProfile_SecurityConfigProfile";
    private static final String SECURITY_USER_ACCOUNT = "security_UserAccountProfile_UserAccountProfile";
    private static final String SECURITY_PASSWORD_POLICY = "security.UserAccountProfile.PasswordPolicy";

    static final String KERNEL_MODULE_CFG_PROFILE = "kernelModule_moduleProfile_KernelModuleConfigProfile";
    static final String KERNEL_MODULE_PROFILE = "kernelModule_moduleProfile_KernelModuleProfile";
    static final String KERNEL_MODULE_HEALTHCHK_KEY = "KernelModuleProfile-healthchk-key";
    static final String KERNEL_MODULE_PARAM_PROFILE = "kernelModule_moduleProfile_KernelModuleParamProfile";
    static final String KERNEL_MODULE_HEAPSIZE_KEY = "KernelModuleParamProfile-heapsize-key";
    static final String HOST_CACHE_CFG_PROFILE = "hostCache_hostCacheConfig_HostCacheConfigProfile";

    /**
     * Method that is used only for the purpose of Favorite tests prerequisites. It
     * sets/unsets the first advanced option as favorite and the ip address settings in
     * Host port group -> Management Network as favorite.
     *
     * @param hostProfileSpec - specification of the profile whose policies are to be set
     * @param set - true to mark them as favorites and false to unmark them
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC
     */
    boolean setSomeFavorites(final HostProfileSpec hostProfileSpec, boolean set) throws Exception {
        boolean result = false;
        HostProfile hostProfile = validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();

        // set first adv option as favorite
        OptionProfile[] advSttngsOptions = hostApplyProfile.getOption();
        if (advSttngsOptions.length > 0) {
            boolean found = false;
            for (int i = 0; i < advSttngsOptions.length; i++) {
                if (CFG_ETC_MOTD.equals(advSttngsOptions[i].getKey())) {
                    advSttngsOptions[i].setFavorite(set);
                    found = true;
                    break;
                }
            }
            if (!found) {
                advSttngsOptions[0].setFavorite(set);
            }
            result = true;
        }

        // set ip address settings of Host prot group management nw as favorite
        NetworkProfile nwProfile = hostApplyProfile.getNetwork();
        if (nwProfile != null) {
            HostPortGroupProfile[] hpgProfile = nwProfile.getHostPortGroup();
            if (hpgProfile.length > 0) {
                hpgProfile[0].getIpConfig().setFavorite(set);
                result = true;
            }
        }

        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.setName(hostProfile.getName());
        completeConfigSpec.setApplyProfile(hostApplyProfile);

        hostProfile.update(completeConfigSpec);

        return result;
    }

   boolean getSomeFavorites(HostProfileSpec hostProfileSpec) throws Exception {
      HostApplyProfile hostApplyProfile = getHostApplyProfile(hostProfileSpec);

      boolean result = isEtcMotdFavorite(hostApplyProfile);
      result &= isHpgProfileFavorite(hostApplyProfile);

      return result;
   }

   private HostApplyProfile getHostApplyProfile(
         HostProfileSpec hostProfileSpec) throws Exception {
      HostProfile hostProfile = validateSpecAndGetProfile(hostProfileSpec);
      return getHostApplyProfile(hostProfile);
   }

   private boolean isHpgProfileFavorite(HostApplyProfile hostApplyProfile) {
      boolean result = false;
      NetworkProfile nwProfile = hostApplyProfile.getNetwork();
      if (nwProfile != null) {
         HostPortGroupProfile[] hpgProfile = nwProfile.getHostPortGroup();
         if (hpgProfile.length > 0) {
            result = hpgProfile[0].getIpConfig().getFavorite();
         }
      }
      return result;
   }

   private boolean isEtcMotdFavorite(HostApplyProfile hostApplyProfile) {
      boolean result = false;
      OptionProfile[] advSttngsOptions = hostApplyProfile.getOption();
      if (advSttngsOptions.length > 0) {
         boolean found = false;
         for (int i = 0; i < advSttngsOptions.length; i++) {
            if (CFG_ETC_MOTD.equals(advSttngsOptions[i].getKey())) {
               result = advSttngsOptions[i].getFavorite();
               found = true;
               break;
            }
         }
         if (!found) {
            result = advSttngsOptions[0].getFavorite();
         }
      }
      return result;
   }

   /**
     * For testing purposes this method sets the ntpClient to enable/disable ignore compliance check
     *
     * @param hostProfileSpec - host profile to be modified
     * @param set - true to enable icc, false to disable it
     * @return - true if successful, false otherwise
     * @throws Exception - if there is a problem with vc connection, the profile is not found or the firewall or ruleset
     *             profile are missing in the host profile
     */
    boolean setSomeIgnoreCc(final HostProfileSpec hostProfileSpec, boolean set) throws Exception {
        boolean result = false;
        HostProfile hostProfile = validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();

        FirewallProfile firewall = hostApplyProfile.getFirewall();
        HostProfilesUtil.ensureNotNull(firewall, "There is no firewall profile!");

        FirewallProfile.RulesetProfile[] ruleset = firewall.getRuleset();
        HostProfilesUtil.ensureNotNull(ruleset, "No ruleset profiles found!");

        for (int i = 0; i < ruleset.length; i++) {
            if (NTP_CLIENT.equals(ruleset[i].getKey())) {
//                ruleset[i].setIgnoreComplianceCheck(set);
                result = true;
                break;
            }
        }

        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.setName(hostProfile.getName());
        completeConfigSpec.setApplyProfile(hostApplyProfile);

        hostProfile.update(completeConfigSpec);

        return result;
    }

   boolean getSomeIgnoreCc(HostProfileSpec profileSpec) throws Exception {
      HostApplyProfile hostApplyProfile = getHostApplyProfile(profileSpec);
      FirewallProfile firewall = hostApplyProfile.getFirewall();
      HostProfilesUtil.ensureNotNull(firewall, "There is no firewall profile!");

      FirewallProfile.RulesetProfile[] ruleSet = firewall.getRuleset();
      HostProfilesUtil.ensureNotNull(ruleSet, "No rule-set profiles found!");

      boolean result = false;
      for (int i = 0; i < ruleSet.length; i++) {
         if (NTP_CLIENT.equals(ruleSet[i].getKey())) {
//            result = ruleSet[i].getIgnoreComplianceCheck();
            break;
         }
      }

      return result;
   }

   private HostApplyProfile getHostApplyProfile(HostProfile hostProfile) {
      HostProfile.ConfigInfo configInfo;
      configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
      return configInfo.getApplyProfile();
   }

   /**
     * Method that sets/unsets the property Host Config Cache ( if it exists)
     * either to Enabled Ignore Compliance Check or Disabled Ignore Compliance
     * Check
     *
     * @param hostProfileSpec - host profile to modify
     * @param set - true if enable icc should be enabled, false to disable it
     * @return - true if successful, false otherwise
     * @throws Exception - if there is an issue with VC connection, or with
     *             the profile
     */
    boolean setHostCacheCfgIgnoreCC(final HostProfileSpec hostProfileSpec, boolean set) throws Exception {
        boolean result = false;
        HostProfile hostProfile = validateSpecAndGetProfile(hostProfileSpec);
        HostApplyProfile hostApplyProfile = getApplyProfile(hostProfile);

        ApplyProfileProperty[] appProfileProps = hostApplyProfile.getProperty();
        ApplyProfileProperty hostCacheCfg = null;

        for (int i = 0; i < appProfileProps.length; i++) {
            if (HOST_CACHE_CFG_PROFILE.equals(appProfileProps[i].getPropertyName())) {
                hostCacheCfg = appProfileProps[i];
                result = true;
                break;
            }
        }

        HostProfilesUtil.ensureNotNull(hostCacheCfg, "There is no Host Cache Configuration in the Host Profile!");

        ApplyProfile[] hostCacheCfgProfile = hostCacheCfg.getProfile();
        HostProfilesUtil.ensureNotNull(hostCacheCfgProfile, "There is no profile for the Host Cache Configuration!");
        if (hostCacheCfgProfile.length == 0) {
            throw new RuntimeException("There is no profile for the Host Cache Configuration!");
        }
//        hostCacheCfgProfile[0].setIgnoreComplianceCheck(set);

        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.setName(hostProfile.getName());
        completeConfigSpec.setApplyProfile(hostApplyProfile);

        hostProfile.update(completeConfigSpec);

        return result;
    }

    /**
     * Method that sets the UI policy Security and Services -> Security Settings ->
     * Security configuration : Administrator password to Configure a fixed administrator
     * password, if a password is supplied, and sets the supplied password. If no password
     * is supplied, it is set to Leave unchanged the administrator password.
     *
     * @param hostProfileSpec - the host profile to modify
     * @param password - the string value of the password to set
     * @return true if successful, otherwise false
     * @throws Exception - if connection with VC fails, or if any of the needed values
     *             cannot be gotten
     */
    boolean setSecurityConfigurationAdminPassword(final HostProfileSpec hostProfileSpec, String password)
        throws Exception {
        boolean result = false;
        HostProfile hostProfile = validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();

        ApplyProfileProperty[] applyProfileProperties = hostApplyProfile.getProperty();

        HostProfilesUtil.ensureNotNull(applyProfileProperties, "Profile properties not found!");

        // The health check can be found in properties:
        // security_SecurityProfile_SecurityConfigProfile ->
        // security_SecurityProfile_SecurityConfigProfile ->
        // security_UserAccountProfile_UserAccountProfile ->
        // security_UserAccountProfile_UserAccountProfile
        ApplyProfile[] applyProfiles = getApplyProfileFromApplyProfileProps(SECURITY_PROFILE, applyProfileProperties);

        ApplyProfileProperty[] appProfileProps = getApplyProfilePropFromApplyProfile(
            SECURITY_PROFILE,
            applyProfiles,
            false);
        ApplyProfile[] appProfilesUser = getApplyProfileFromApplyProfileProps(SECURITY_USER_ACCOUNT, appProfileProps);

        ApplyProfile rootApplyProfile = getPwdApplyProfilePerUser(appProfilesUser, "root");

        HostProfilesUtil.ensureNotNull(rootApplyProfile, "No root apply profile found!");

        Policy[] rootPolicies = rootApplyProfile.getPolicy();

        HostProfilesUtil.ensureNotNull(rootPolicies, "No root policies found!");

        PolicyOption rootPwd = null;
        for (Policy rootPolicy : rootPolicies) {
            if (SECURITY_PASSWORD_POLICY.equals(rootPolicy.getId())) {
                rootPwd = rootPolicy.getPolicyOption();
                result = true;
                break;
            }
        }

        HostProfilesUtil.ensureNotNull(rootPwd, "No root password found!");

        if (password != null) {
            rootPwd.setId("security.UserAccountProfile.FixedPasswordConfigOption");
            KeyAnyValue pwd = new KeyAnyValueImpl();
            pwd.setKey("password");
            pwd.setValue(new PasswordField("TestPassword"));
            rootPwd.setParameter(new KeyAnyValue[] { pwd });
        } else {
            rootPwd.setId("security.UserAccountProfile.DefaultAccountPasswordUnchangedOption");
            rootPwd.setParameter(null);
        }

        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.setName(hostProfile.getName());
        completeConfigSpec.setApplyProfile(hostApplyProfile);

        hostProfile.update(completeConfigSpec);

        return result;
    }

    /**
     * Method that sets the heap size of the Kernel Module parameter Health Check to 512;
     * this is used as a prerequisite for tests that would need a reboot requirement for
     * remediate
     *
     * @param hostProfileSpec - the host profile to modify
     * @param set - whether to set it or to unset it
     * @return - true if successful
     * @throws Exception
     */
    boolean setKernelModuleHealthCheckHeapSize512(final HostProfileSpec hostProfileSpec, boolean set) throws Exception {
        boolean result = false;
        HostProfile hostProfile = validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();

        ApplyProfileProperty[] applyProfileProperties = hostApplyProfile.getProperty();

        HostProfilesUtil.ensureNotNull(applyProfileProperties, "Profile properties not found!");

        // The health check can be found in properties:
        // kernelModule_moduleProfile_KernelModuleConfigProfile ->
        // kernelModule_moduleProfile_KernelModuleConfigProfile ->
        // kernelModule_moduleProfile_KernelModuleProfile ->
        // KernelModuleProfile-healthchk-key ->
        // kernelModule_moduleProfile_KernelModuleParamProfile ->
        // KernelModuleParamProfile-heapsize-key
        ApplyProfile[] applyProfiles = getApplyProfileFromApplyProfileProps(
            KERNEL_MODULE_CFG_PROFILE,
            applyProfileProperties);

        ApplyProfileProperty[] appProfilePropsKernelModule = getApplyProfilePropFromApplyProfile(
            KERNEL_MODULE_CFG_PROFILE,
            applyProfiles,
            false);
        ApplyProfile[] appProfileKernelModule = getApplyProfileFromApplyProfileProps(
            KERNEL_MODULE_PROFILE,
            appProfilePropsKernelModule);

        ApplyProfileProperty[] kernelModuleHealtchChkProps = getApplyProfilePropFromApplyProfile(
            KERNEL_MODULE_HEALTHCHK_KEY,
            appProfileKernelModule,
            true);

        ApplyProfile[] healthChkApplyProfile = getApplyProfileFromApplyProfileProps(
            KERNEL_MODULE_PARAM_PROFILE,
            kernelModuleHealtchChkProps);


        Policy[] healthCheckHeapsizePolicy = null;
        for (int i = 0; i < healthChkApplyProfile.length; i++) {
            if ((KERNEL_MODULE_HEAPSIZE_KEY).equals(((ApplyProfileElement) healthChkApplyProfile[i]).getKey())) {
                healthCheckHeapsizePolicy = ((ApplyProfileElement) healthChkApplyProfile[i]).getPolicy();
                result = true;
                break;
            }
        }

        HostProfilesUtil.ensureNotNull(
            healthCheckHeapsizePolicy,
            "Kernel Module Profile Healcth Check Heap size policies not found!");

        if (healthCheckHeapsizePolicy.length == 0) {
            throw new RuntimeException("Kernel Module Profile Healcth Check Heap size policies empty!");
        }

        PolicyOption healthcchkHeapsizePolicyOption = healthCheckHeapsizePolicy[0].getPolicyOption();
        KeyAnyValue[] parameters = healthcchkHeapsizePolicyOption.getParameter();
        KeyAnyValue[] newValuesArray;
        if (set) {
            if (parameters.length == 1) {
                newValuesArray = new KeyAnyValue[2];
                newValuesArray[0] = parameters[0];
                KeyAnyValue value = new KeyAnyValueImpl();
                value.setKey(PARAM_VALUE);
                value.setValue("512");
                newValuesArray[1] = value;
            } else {
                parameters[1].setKey(PARAM_VALUE);
                parameters[1].setValue("512");
                newValuesArray = parameters;
            }
        } else {
            newValuesArray = new KeyAnyValue[] { new KeyAnyValueImpl() };
            newValuesArray[0].setKey(PARAM_NAME);
            newValuesArray[0].setValue("heapSize");
        }
        healthcchkHeapsizePolicyOption.setParameter(newValuesArray);
        healthCheckHeapsizePolicy[0].setPolicyOption(healthcchkHeapsizePolicyOption);

        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.setName(hostProfile.getName());
        completeConfigSpec.setApplyProfile(hostApplyProfile);

        hostProfile.update(completeConfigSpec);

        return result;
    }

    /**
     * Method that gets the profiles from properties array by the property name
     *
     * @param propName - name of the property whose profiles to get
     * @param appProfileProps - array of the properties
     * @return - the array of profiles associated with the specified property
     * @exception - if there is no profile with such a name
     */
    ApplyProfile[] getApplyProfileFromApplyProfileProps(String propName, ApplyProfileProperty[] appProfileProps) {

        if (propName == null || appProfileProps == null) {
            throw new IllegalArgumentException("Supplied params are null!");
        }

        for (int i = 0; i < appProfileProps.length; i++) {
            if (propName.equals(appProfileProps[i].getPropertyName())) {
                return appProfileProps[i].getProfile();
            }
        }

        throw new RuntimeException("No ApplyProfile found!");
    }

    /**
     * Method to get the profile properties from a profile through filtering by key or by
     * profile type name
     *
     * @param filter - the profile
     * @param applyProfile - the array of profiles in which to search
     * @param byKey - boolean parameter to show whether to search in the key or profile
     *            type name
     * @return - array of profile properties
     * @exception - if there is no profile with such key/type name
     */
    ApplyProfileProperty[] getApplyProfilePropFromApplyProfile(String filter, ApplyProfile[] applyProfile,
        boolean byKey) {
        if (filter == null || applyProfile == null) {
            throw new IllegalArgumentException("Supplied params are null!");
        }
        for (int i = 0; i < applyProfile.length; i++) {
            String value;
            if (byKey) {
                value = ((ApplyProfileElement) applyProfile[i]).getKey();
            } else {
                value = applyProfile[i].getProfileTypeName();
            }
            if (filter.equals(value)) {
                return applyProfile[i].getProperty();
            }
        }

        throw new RuntimeException("No ApplyProfileProperty found!");
    }

    /**
     * Method that sets as enabled/disabled the DHCP in netStack
     *
     * @param network - network profile to edit
     * @param set - true to set, false otherwise
     */
    void setNetStackDhcp(NetworkProfile network, Boolean set) {
        ApplyProfileProperty[] properties = network.getProperty();
        HostProfilesUtil.ensureNotNull(properties, "There are no network properties!");

        ApplyProfile[] applyProfiles = getApplyProfileFromApplyProfileProps(
            GENERIC_NETSTACK_INSTANCE_PROFILE,
            properties);
        ApplyProfile genericNetStack = null;
        for (ApplyProfile applyProfile : applyProfiles) {
            if (GENERIC_NETSTACK_INSTANCE_PROFILE.equals(applyProfile.getProfileTypeName())) {
                genericNetStack = applyProfile;
                break;
            }
        }
        HostProfilesUtil.ensureNotNull(genericNetStack, "There is no Generic NetStack instance!");

        ApplyProfileProperty genericDnsCfgProfile = null;
        for (ApplyProfileProperty property : genericNetStack.getProperty()) {
            if (GENERIC_DNS_CFG_PROFILE.equals(property.getPropertyName())) {
                genericDnsCfgProfile = property;
                break;
            }
        }
        HostProfilesUtil.ensureNotNull(genericDnsCfgProfile, "There is no Generic NetStack profile!");

        ApplyProfile[] dnsProfiles = genericDnsCfgProfile.getProfile();
        ApplyProfile dnsProfile = null;
        for (ApplyProfile profile : dnsProfiles) {
            if (GENERIC_DNS_CFG_PROFILE.equals(profile.getProfileTypeName())) {
                dnsProfile = profile;
                break;
            }
        }
        HostProfilesUtil.ensureNotNull(dnsProfile, "There is no DNS profile!");

        Policy[] policies = dnsProfile.getPolicy();
        Policy dnsCfgPolicy = null;
        for (Policy policy : policies) {
            if (DNS_CFG_POLICY.equals(policy.getId())) {
                dnsCfgPolicy = policy;
                break;
            }
        }
        HostProfilesUtil.ensureNotNull(dnsCfgPolicy, "There is no DNS configuration policy!");

        PolicyOption policyOption = dnsCfgPolicy.getPolicyOption();
        KeyAnyValue[] parameters = policyOption.getParameter();
        KeyAnyValue dhcp = null;
        for (KeyAnyValue parameter : parameters) {
            if (DHCP.equals(parameter.getKey())) {
                dhcp = parameter;
                break;
            }
        }

        if (set != dhcp.getValue()) {
            dhcp.setValue(set);
        }

    }

   /**
    * Method that gets the Host Port Group -> Management network
    * @param network - the network profile
    * @return - the management network
    */
    HostPortGroupProfile getMgmntNw(NetworkProfile network) {
        HostProfilesUtil.ensureNotNull(network, "No network profile found in host profile!");

        HostPortGroupProfile[] hostPortGroups = network.getHostPortGroup();
        HostProfilesUtil.ensureNotNull(hostPortGroups, "No host portgroup found in host profile!");

        HostPortGroupProfile hostPortGroup = null;
        for (HostPortGroupProfile hpGroup : hostPortGroups) {
            if (MGMT_NW_HOST_PORT_GROUP.equals(hpGroup.getKey())) {
                hostPortGroup = hpGroup;
                break;
            }
        }
        HostProfilesUtil.ensureNotNull(hostPortGroup, "No Management Network found!");

       return hostPortGroup;
    }

    /**
     * Method that removes a setting from profile, this is actually as delete
     *
     * @param parentSubprofile - subprofile to edit
     * @param sttngName - setting to remove
     * @return - modified subprofile
     */
    ApplyProfile[] removeSetting(ApplyProfile[] parentSubprofile, String sttngName) {

        for (ApplyProfile element : parentSubprofile) {
            if (sttngName.equals(((ApplyProfileElement) element).getKey())) {
                return (ApplyProfile[]) ArrayUtils.removeElement(parentSubprofile, element);
            }
        }

        throw new RuntimeException("Setting: " + sttngName + " not found!");
    }

    /**
     * Returns a CompleteConfigSpec from a host profile.
     *
     * @param hostProfile host profile
     * @param applyProfile
     * @return complete config spec
     */
    HostProfile.CompleteConfigSpec getCompleteConfigSpec(HostProfile hostProfile, HostApplyProfile applyProfile) {
        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.name = hostProfile.getName();

        completeConfigSpec.applyProfile = applyProfile;
        return completeConfigSpec;
    }

    /**
     * Executes the update action on a host profile, using the complete config spec as new
     * values.
     *
     * @param hostProfile the host profile
     * @param completeConfigSpec
     */
    void executeUpdateAction(HostProfile hostProfile, HostProfile.CompleteConfigSpec completeConfigSpec) {
        try {
            hostProfile.update(completeConfigSpec);
        } catch (DuplicateName duplicateName) {
            throw new RuntimeException("Duplicate host profile name.", duplicateName);
        } catch (ProfileUpdateFailed profileUpdateFailed) {
            throw new RuntimeException("Update of host profile failed", profileUpdateFailed);
        }
    }

    /**
     * Get the apply profile from a host profile.
     *
     * @param hostProfile host profile
     * @return apply profile
     */
    HostApplyProfile getApplyProfile(HostProfile hostProfile) {
        final HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        return configInfo.applyProfile;
    }

    /**
     * Finds and updates apply profile values
     *
     * @param applyProfile apply profile
     * @param hostPolicyUpdateSpec policy option spec
     */
    void findAndUpdateApplyProfileValues(HostApplyProfile applyProfile, HostPolicyUpdateSpec[] hostPolicyUpdateSpec) {

        final OptionProfile[] optionProfiles = applyProfile.getOption();

        for (HostPolicyUpdateSpec policyUpdateSpec : hostPolicyUpdateSpec) {

            final String propToChange = policyUpdateSpec.originalPropertyName.get();

            for (int i = 0; i < optionProfiles.length; i++) {
                OptionProfile optionProfile = optionProfiles[i];
                Policy policy = optionProfile.policy[0];
                PolicyOption policyOption = policy.policyOption;
                KeyAnyValue parameter = policyOption.parameter[0];
                String value = (String) parameter.getValue();

                if (propToChange.equals(value)) {
                    updatePolicyWithNewValues(policy, policyUpdateSpec);
                    break;
                }

            }
        }
    }

    /**
     * Updates the host profile policy with new values from the policy option spec
     *
     * @param policy the policy to update
     * @param hostPolicyUpdateSpec the policy option spec with the new values
     */
    void updatePolicyWithNewValues(Policy policy, HostPolicyUpdateSpec hostPolicyUpdateSpec) {
        PolicyOption policyOption = policy.getPolicyOption();

        setPolicyType(hostPolicyUpdateSpec, policyOption);
        setPolicyParameters(hostPolicyUpdateSpec, policyOption);

        policy.setPolicyOption(policyOption);
    }

    void setPolicyParameters(HostPolicyUpdateSpec hostPolicyUpdateSpec, PolicyOption policyOption) {
        KeyAnyValue[] parameters = getKeyAnyValueArray(hostPolicyUpdateSpec);
        policyOption.setParameter(parameters);
    }

    KeyAnyValue[] getKeyAnyValueArray(HostPolicyUpdateSpec hostPolicyUpdateSpec) {
        KeyAnyValue keyParameter = new KeyAnyValueImpl();
        keyParameter.setKey(KEY);
        String value = hostPolicyUpdateSpec.newPropertyName.get();
        keyParameter.setValue(value);

        KeyAnyValue[] policyOptionParameters = { keyParameter };
        return policyOptionParameters;
    }

    void setPolicyType(HostPolicyUpdateSpec hostPolicyUpdateSpec, PolicyOption policyOption) {
        PolicyOptionType policyType = hostPolicyUpdateSpec.newPolicyType.get();
        String policyTypeString = policyType.getValue();
        policyOption.setId(policyTypeString);
    }

    HostProfile validateSpecAndGetProfile(final HostProfileSpec hostProfileSpec) throws Exception {
        TestSpecValidator.ensureNotNull(hostProfileSpec, "Supply HostProfileSpec!");
        HostProfile hostProfile = HostProfileSrvApi.getInstance().getHostProfileByName(hostProfileSpec);
        HostProfilesUtil.ensureNotNull(hostProfile, "No such host profile: " + hostProfileSpec.name.get());

        return hostProfile;
    }

    private ApplyProfile getPwdApplyProfilePerUser(ApplyProfile[] appProfiles, String userName) {
        for (ApplyProfile appProfile : appProfiles) {
            if ("security_UserAccountProfile_UserAccountProfile".equals(appProfile.getProfileTypeName())) {
                Policy[] policies = appProfile.getPolicy();
                for (Policy policy : policies) {
                    if ("security.UserAccountProfile.UserPolicy".equals(policy.getId())) {
                        PolicyOption policyOpt = policy.getPolicyOption();
                        KeyAnyValue[] params = policyOpt.getParameter();
                        for (KeyAnyValue param : params) {
                            if (userName.equals(param.getValue())) {
                                return appProfile;
                            }
                        }
                    }
                }
            }
        }
        return null;
    }
}
