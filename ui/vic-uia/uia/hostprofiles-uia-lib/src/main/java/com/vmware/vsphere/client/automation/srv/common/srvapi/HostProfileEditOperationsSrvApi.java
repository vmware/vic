/*
 * Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import com.vmware.vim.binding.impl.vmodl.KeyAnyValueImpl;
import com.vmware.vim.binding.vim.profile.ApplyProfile;
import com.vmware.vim.binding.vim.profile.ApplyProfileProperty;
import com.vmware.vim.binding.vim.profile.Policy;
import com.vmware.vim.binding.vim.profile.PolicyOption;
import com.vmware.vim.binding.vim.profile.host.*;
import com.vmware.vim.binding.vmodl.KeyAnyValue;
import com.vmware.vsphere.client.automation.common.HostProfilesUtil;
import com.vmware.vsphere.client.automation.common.spec.HostPolicyUpdateSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostProfileSpec;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Class that represents different editing operation on a host profile through the VC API
 */
public class HostProfileEditOperationsSrvApi {

    private HostProfileEditOperationsHelper helper;

    private static final Logger logger = LoggerFactory.getLogger(HostProfileEditOperationsSrvApi.class);
    private static final String MODEL = "Model";
    private static final String SIZE = "SizeMB";
    private static final String HOST_CACHE_CFG_POLICY_OPT = "hostCache.hostCacheConfig.HostCacheConfigProfilePolicyOption";
    private static final String HOST_CACHE_CFG_PROFILE_POLICY = "hostCache.hostCacheConfig.HostCacheConfigProfilePolicy";
    private static final String IP_ADDRESS_POLICY = "IpAddressPolicy";
    private static final String DHCP_FIXED_OPT = "FixedDhcpOption";
    private static final String USER_REQ_IP_ADDRESS = "UserInputIPAddress_UseDefault";
    private static final String MTU_POLICY ="MtuPolicy";

    private static HostProfileEditOperationsSrvApi INSTANCE = null;

    protected HostProfileEditOperationsSrvApi() {
        // Disallow the creation of this object outside of the class.
        helper = new HostProfileEditOperationsHelper();
    }

    /**
     * Get instance of HostProfileEditOperationsSrvApi.
     *
     * @return created instance
     */
    public static HostProfileEditOperationsSrvApi getInstance() {
        if (INSTANCE == null) {
            synchronized (HostProfileEditOperationsSrvApi.class) {
                if (INSTANCE == null) {
                    logger.info("Initializing HostProfileEditOperationsSrvApi.");
                    INSTANCE = new HostProfileEditOperationsSrvApi();
                }
            }
        }

        return INSTANCE;
    }

    /**
     * Method that is used only for the purpose of Favorite tests prerequisites. It sets
     * the first advanced option as favorite and the ip address settings in Host port group
     * -> Management Network as favorite.
     *
     * @param hostProfileSpec - specification of the profile whose policies are to be set
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC
     */
    public boolean markSomeOptionsAsFavorites(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setSomeFavorites(hostProfileSpec, true);
    }

    /**
     * Method that is used only for the purpose of Favorite tests prerequisites. It unsets
     * the first advanced option as favorite and the ip address settings in Host port group
     * -> Management Network as favorite.
     *
     * @param hostProfileSpec - specification of the profile whose policies are to be set
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC
     */
    public boolean unmarkSomeOptionsAsFavorites(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setSomeFavorites(hostProfileSpec, false);
    }

   /**
    * Method which checks if the first advanced option and ip address
    * settings are marked as favorite
    *
    * @param profileSpec Host profile spec
    * @return true if both options are set as favorite
    * @throws Exception
    */
   public boolean getSomeFavorites(
         HostProfileSpec profileSpec) throws Exception {
      return helper.getSomeFavorites(profileSpec);
   }

   /**
     * Method that is used only for the purpose of having a setting with enabled ignore compliance check.
     * It sets the Security and Services -> Firewall Configuration -> Firewall configuration -> Ruleset
     * Configuration-> ntpClient to enabled ignore compliance check.
     *
     * @param hostProfileSpec - host profile to edit
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC or execute task
     */
    public boolean markNtpClientEnableIcc(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setSomeIgnoreCc(hostProfileSpec, true);
    }

    /**
     * Method that is used only for the purpose of having a setting with disabled ignore compliance check.
     * It sets the Security and Services -> Firewall Configuration -> Firewall configuration -> Ruleset
     * Configuration-> ntpClient to disabled ignore compliance check.
     *
     * @param hostProfileSpec - host profile to edit
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC or execute task
     */
    public boolean unmarkNtpClientEnableIcc(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setSomeIgnoreCc(hostProfileSpec, false);
    }

   /**
    * Method which checks if the ntp client is marked with Ignore Compliance
    * Check.
    *
    * @param profileSpec the host profile spec
    * @return true if the the ntp client is marked with ignore compliance check.
    * @throws Exception
    */
   public boolean getSomeIgnoreCc(
         HostProfileSpec profileSpec) throws Exception {
      return helper.getSomeIgnoreCc(profileSpec);
   }

    /**
     * Method that enables the Ignore Compliance Check for Host Cache Config
     *
     * @param hostProfileSpec - host profile to modify
     * @return true if successful, false otherwise
     * @throws Exception - if there is an issue with VC connection or with the
     *             profile
     */
    public boolean markHostConfigCacheEnableIcc(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setHostCacheCfgIgnoreCC(hostProfileSpec, true);
    }

    /**
     * Method that disables the Ignore Compliance Check for Host Cache Config
     *
     * @param hostProfileSpec - host profile to modify
     * @return true if successful, false otherwise
     * @throws Exception - if there is an issue with VC connection or with the
     *             profile
     */
    public boolean unmarkHostConfigCacheEnableIcc(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setHostCacheCfgIgnoreCC(hostProfileSpec, false);
    }

    /**
     * Method that is used only for the purpose of Remediate tests prerequisites. It sets
     * the kernel module health check parameter's heap size to 512
     *
     * @param hostProfileSpec - specification of the profile whose policies are to be set
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC
     */
    public boolean setRebootParameter(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setKernelModuleHealthCheckHeapSize512(hostProfileSpec, true);
    }

    /**
     * Method that is used only for the purpose of Remediate tests prerequisites. It sets
     * the kernel module health check parameter's heap size to none
     *
     * @param hostProfileSpec - specification of the profile whose policies are to be set
     * @return true if successful, false otherwise
     * @throws Exception if cannot connect to VC
     */
    public boolean unsetRebootParameter(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setKernelModuleHealthCheckHeapSize512(hostProfileSpec, false);
    }

    /**
     * Method that sets the UI policy Security and Services -> Security Settings ->
     * Security configuration : Administrator password to Configure a fixed administrator
     * password and sets the supplied password.
     *
     * @param hostProfileSpec - profile to modify
     * @param password - password string value to set
     * @return true if successful, false otherwise
     * @throws Exception - if problems with VC connection, profile not found, or policy
     *             option not found
     */
    public boolean setAdminPasswordToValue(final HostProfileSpec hostProfileSpec, String password) throws Exception {
        return helper.setSecurityConfigurationAdminPassword(hostProfileSpec, password);
    }

    /**
     * Method that sets the UI policy Security and Services -> Security Settings ->
     * Security configuration : Administrator password to Leave administrator password
     * unchanged, which is teh default value.
     *
     * @param hostProfileSpec - profile to modify
     * @return true if successful, false otherwise
     * @throws Exception - if problems with VC connection, profile not found, or policy
     *             option not found
     */
    public boolean setAdminPasswordUnchanged(final HostProfileSpec hostProfileSpec) throws Exception {
        return helper.setSecurityConfigurationAdminPassword(hostProfileSpec, null);
    }

    /**
     * Mathod that adds the subprofile Host Cache configuration under General
     * System Settings -> Host Cache Configuration and makes it enabled
     *
     * @param hostProfileSpec
     * @param name
     * @return
     * @throws Exception
     */
    public boolean addHostCacheCfg(final HostProfileSpec hostProfileSpec, String name) throws Exception {
        boolean result = false;
        HostProfile hostProfile = helper.validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();

        ApplyProfileProperty[] appProfileProps = hostApplyProfile.getProperty();
        ApplyProfileProperty hostCacheCfg = null;

        for (int i = 0; i < appProfileProps.length; i++) {
            if (helper.HOST_CACHE_CFG_PROFILE.equals(appProfileProps[i].getPropertyName())) {
                hostCacheCfg = appProfileProps[i];
                result = true;
                break;
            }
        }

        HostProfilesUtil.ensureNotNull(hostCacheCfg, "No Host cache configuration present!");

        if (name != null) {
            ApplyProfile[] hostCacheCfgApplyProfile = new ApplyProfile[1];
            ApplyProfile hostCacheProfile = new ApplyProfile();

            KeyAnyValue parameterModel = new KeyAnyValueImpl();
            parameterModel.setKey(MODEL);
            parameterModel.setValue(name);
            KeyAnyValue parameterSize = new KeyAnyValueImpl();
            parameterSize.setKey(SIZE);
            parameterSize.setValue(0);
            KeyAnyValue[] parameters = new KeyAnyValue[] { parameterModel, parameterSize };

            PolicyOption policyOption = new PolicyOption();
            policyOption.setId(HOST_CACHE_CFG_POLICY_OPT);
            policyOption.setParameter(parameters);

            Policy policy = new Policy();
            policy.setId(HOST_CACHE_CFG_PROFILE_POLICY);
            policy.setPolicyOption(policyOption);

            hostCacheProfile.setProfileTypeName(helper.HOST_CACHE_CFG_PROFILE);
            hostCacheProfile.setPolicy(new Policy[] { policy });
            hostCacheProfile.setEnabled(true);
            hostCacheCfgApplyProfile[0] = hostCacheProfile;
            hostCacheCfg.setProfile(hostCacheCfgApplyProfile);
        } else {
            hostCacheCfg.setProfile(null);
        }

        HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
        completeConfigSpec.setName(hostProfile.getName());
        completeConfigSpec.setApplyProfile(hostApplyProfile);

        hostProfile.update(completeConfigSpec);

        return result;
    }

    /**
     * Method that sets the IP Address setting in Network Configuration -> Management network -> Host port group to
     * Prompt the user ..., i.e. user required setting, or to its default value - use dhcp
     *
     * @param hostProfileSpec - teh spec to modify
     * @param set - if true - sets it to prompt the user.., else sets it to default value
     * @return true if succesful
     * @throws Exception if VC connection fails, if setting not found or update task fails
     */
    public boolean setIpv4AddrInMgmntNwUserInput(final HostProfileSpec hostProfileSpec, boolean set) throws Exception {
        boolean result;
        HostProfile hostProfile = helper.validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();
        NetworkProfile network = hostApplyProfile.getNetwork();
        HostPortGroupProfile hostPortGroup = helper.getMgmntNw(network);
        IpAddressProfile ipConfig = hostPortGroup.getIpConfig();
        HostProfilesUtil.ensureNotNull(ipConfig, "No ip configuration found!");
        Policy[] policies = ipConfig.getPolicy();
        HostProfilesUtil.ensureNotNull(policies, "No policy found!");

        Policy ipAddressPolicy = null;
        for (Policy policy : policies) {
            if (IP_ADDRESS_POLICY.equals(policy.getId())) {
                ipAddressPolicy = policy;
                break;
            }
        }
        HostProfilesUtil.ensureNotNull(ipAddressPolicy, "No policy with id IPAddressPolicy found!");

        result = true;
        String newValue;
        boolean isDhcpNetStack = false;
        if (set) {
            newValue = USER_REQ_IP_ADDRESS;
        } else {
            newValue = DHCP_FIXED_OPT;
            isDhcpNetStack = true;
        }

        // in order to set ipv4 address to user prompt dhcp in NetStack has to be disabled
        helper.setNetStackDhcp(network, isDhcpNetStack);

        String currentValue = ipAddressPolicy.getPolicyOption().getId();
        if (!newValue.equals(currentValue)) {
            ipAddressPolicy.getPolicyOption().setId(newValue);

            HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
            completeConfigSpec.setName(hostProfile.getName());
            completeConfigSpec.setApplyProfile(hostApplyProfile);

            hostProfile.update(completeConfigSpec);
        }

        return result;
    }

    /**
     * Methdo that sets a new value to the MTU parameter of Host Port Group ->
     * Management network
     *
     * @param hostProfileSpec - host profile to modify
     * @param value - the new value to set
     * @return - true if successful
     * @throws Exception -  if mtu is not found, vc connection error, or profile does not
     *             exist
     */
    public boolean setMgmntNwMtuValue(final HostProfileSpec hostProfileSpec, int value) throws Exception {
        boolean result;
        HostProfile hostProfile = helper.validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();
        NetworkProfile network = hostApplyProfile.getNetwork();
        HostPortGroupProfile hostPortGroup = helper.getMgmntNw(network);
        Policy[] policies = hostPortGroup.getPolicy();
        Policy mtu = null;

        for (Policy policy : policies) {
            if (MTU_POLICY.equals(policy.getId())) {
                mtu = policy;
                break;
            }
        }

        HostProfilesUtil.ensureNotNull(mtu, "There is no MTU policy in Management Network!");

        result = true;
        int currentValue = (Integer) mtu.getPolicyOption().getParameter()[0].getValue();
        if (value != currentValue) {
            mtu.getPolicyOption().getParameter()[0].setValue(value);

            HostProfile.CompleteConfigSpec completeConfigSpec = new HostProfile.CompleteConfigSpec();
            completeConfigSpec.setName(hostProfile.getName());
            completeConfigSpec.setApplyProfile(hostApplyProfile);

            hostProfile.update(completeConfigSpec);
        }

        return result;
    }

   /**
    * Method that gets the current value of MTU parameter of Host Port Group ->
    * Management Network
    * @param hostProfileSpec - hots profile whose value to get
    * @return - value of MTU
    * @throws Exception -  if mtu is not found, vc connection error, or profile does not
    *             exist
    */
    public Integer getMgmntNwMtuValue(final HostProfileSpec hostProfileSpec) throws Exception {
        HostProfile hostProfile = helper.validateSpecAndGetProfile(hostProfileSpec);
        HostProfile.ConfigInfo configInfo = (HostProfile.ConfigInfo) hostProfile.getConfig();
        HostApplyProfile hostApplyProfile = configInfo.getApplyProfile();
        NetworkProfile network = hostApplyProfile.getNetwork();
        HostPortGroupProfile hostPortGroup = helper.getMgmntNw(network);
        Policy[] policies = hostPortGroup.getPolicy();
        Policy mtu = null;

        for (Policy policy : policies) {
            if (MTU_POLICY.equals(policy.getId())) {
                mtu = policy;
                break;
            }
        }

        HostProfilesUtil.ensureNotNull(mtu, "There is no MTU policy in Management Network!");

        return (Integer) mtu.getPolicyOption().getParameter()[0].getValue();
    }

    /**
     * Method that removes the kernel module health check profile
     *
     * @param hostProfileSpec - host profile to modify
     * @return true if successful, false otherwise
     * @throws Exception - if the VC connection fails, or any of the path elements is not found, or
     *             update doesn't succeed
     */
    public boolean removeKernelModuleHealthchk(HostProfileSpec hostProfileSpec) throws Exception {
        boolean result = false;
        HostProfile hostProfile = helper.validateSpecAndGetProfile(hostProfileSpec);
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
        ApplyProfile[] applyProfiles = helper
            .getApplyProfileFromApplyProfileProps(helper.KERNEL_MODULE_CFG_PROFILE, applyProfileProperties);

        ApplyProfileProperty[] appProfilePropsKernelModule = helper
            .getApplyProfilePropFromApplyProfile(helper.KERNEL_MODULE_CFG_PROFILE, applyProfiles, false);
        ApplyProfile[] appProfileKernelModule = helper
            .getApplyProfileFromApplyProfileProps(helper.KERNEL_MODULE_PROFILE, appProfilePropsKernelModule);

        ApplyProfile[] modifiedAppProfileKernelModule = helper
            .removeSetting(appProfileKernelModule, helper.KERNEL_MODULE_HEALTHCHK_KEY);

        for (ApplyProfileProperty property : appProfilePropsKernelModule) {
            if (helper.KERNEL_MODULE_PROFILE.equals(property.getPropertyName())) {
                property.setProfile(modifiedAppProfileKernelModule);
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

    /**
     * Method used to update Host Profile properties (Advanced Configuration Settings >
     * Advanced Options > *)
     * <p/>
     * NOTE: If a profile policy is set to USER_EXPLICIT_POLICY_CHOICE, then the first such
     * property will be updated.
     *
     * @param hostProfileSpec the host profile spec
     * @param hostPolicyUpdateSpec the values to set
     */
    public void updateAdvancedOptions(HostProfileSpec hostProfileSpec, HostPolicyUpdateSpec[] hostPolicyUpdateSpec)
        throws Exception {
        final HostProfile hostProfile = helper.validateSpecAndGetProfile(hostProfileSpec);
        final HostApplyProfile applyProfile = helper.getApplyProfile(hostProfile);

        helper.findAndUpdateApplyProfileValues(applyProfile, hostPolicyUpdateSpec);

        // get complete config spec
        HostProfile.CompleteConfigSpec completeConfigSpec;
        completeConfigSpec = helper.getCompleteConfigSpec(hostProfile, applyProfile);

        helper.executeUpdateAction(hostProfile, completeConfigSpec);
    }
}
