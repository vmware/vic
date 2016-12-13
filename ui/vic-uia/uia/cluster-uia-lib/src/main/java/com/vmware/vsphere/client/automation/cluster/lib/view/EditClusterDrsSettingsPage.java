/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.view;

import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.srv.common.spec.ClusterSpec.DrsAutoLevel;

/**
 * Represents the "vSphere DRS" page in the cluster edit dialog.
 */
public class EditClusterDrsSettingsPage extends EditClusterPageNavigator {

    private static final IDGroup DRS_ENABLE_CB =
            IDGroup.toIDGroup("enableDrsCheckBox");
    private static final IDGroup DRS_AUT_LEVEL_DD =
            IDGroup.toIDGroup("drsAutLevelDropDown");
    private static final IDGroup DRS_AUTOMATION_SECTION_IMG =
            IDGroup.toIDGroup("arrowImagedrsAutomationBlock");
    private static final IDGroup ENABLE_VM_DRS_AUT_LEVEL_CB =
            IDGroup.toIDGroup("_DrsConfigForm_CheckBox2");
    private static final IDGroup DRS_MAN_LEVEL_RB =
            IDGroup.toIDGroup("_DrsConfigForm_RadioButton1");
    private static final IDGroup DRS_PARTIALLY_AUTOMATED_LEVEL_RB =
            IDGroup.toIDGroup("_DrsConfigForm_RadioButton2");
    private static final IDGroup DRS_FULLY_AUTOMATED_LEVEL_RB =
            IDGroup.toIDGroup("_DrsConfigForm_RadioButton3");

    private static final IDGroup DRS_POLICIES_SECTION_IMG =
          IDGroup.toIDGroup("arrowImagedrsPoliciesBlock");
    private static final IDGroup DRS_POLICIES_ENFORCE_EVEN_DIST_CB =
          IDGroup.toIDGroup("_DrsConfigForm_CheckBox3");
    private static final IDGroup DRS_POLICIES_CONSUMED_MEM_CB =
          IDGroup.toIDGroup("_DrsConfigForm_CheckBox4");
    private static final IDGroup DRS_POLICIES_CPU_OVERCOM_CB =
          IDGroup.toIDGroup("_DrsConfigForm_CheckBox5");
    private static final IDGroup DRS_POLICIES_CPU_OVERCOM_TB =
          IDGroup.toIDGroup("tiwoDialog/drsPoliciesBlock/_DrsConfigForm_PropertyGridRow8/nscpuPercent");
    private static final IDGroup DRS_POLICIES_CPU_OVERCOM_INC =
          IDGroup.toIDGroup("nscpuPercent/incrementButton");
    private static final IDGroup DRS_POLICIES_CPU_OVERCOM_DEC =
          IDGroup.toIDGroup("nscpuPercent/decrementButton");
    /**
     * Check "Turn on vSphere DRS" check box.
     *
     * @param value   if the check box should be checked or not
     */
    public void setDrsEnabled(boolean value) {
        UI.component.value.set(value, DRS_ENABLE_CB);
    }

    /**
     * Gets the status(checked or not) of the "Turn on vSphere DRS" check box.
     */
    public boolean isDrsEnabled() {
        return Boolean.valueOf(UI.component.value.get(DRS_ENABLE_CB));
    }

    /**
     * Gets status of Turn on vSphere DRS check box: disabled, enabled for selection.
     */
    public boolean isDrsCbEnabled() {
        return Boolean.valueOf(UI.component.property.getBoolean(Property.ENABLED, DRS_ENABLE_CB));
    }

    /**
     * Gets status of Manual Automation Level radio button: disabled, enabled for selection.
     */
    public boolean isManualAutomatedRbEnabled() {
        return UI.component.property.getBoolean(Property.ENABLED, DRS_MAN_LEVEL_RB);
    }

    /**
     * Gets status of Partially Automated radio button: disabled, enabled for selection.
     */
    public boolean isPartiallyAutomatedRbEnabled() {
        return UI.component.property.getBoolean(Property.ENABLED, DRS_PARTIALLY_AUTOMATED_LEVEL_RB);
    }

    /**
     * Gets status of Fully Automated radio button: disabled, enabled for selection.
     */
    public boolean isFullyAutomatedRbEnabled() {
        return UI.component.property.getBoolean(Property.ENABLED, DRS_FULLY_AUTOMATED_LEVEL_RB);
    }

    /**
     * Click(expand/collapse) "DRS Automation" section of this view.
     */
    public void clickDrsAutomationSection() {
        UI.component.click(DRS_AUTOMATION_SECTION_IMG);
    }

    /**
     * Sets DRS automation level.
     *
     * @param value     automation level to be set
     */
    public void setDrsAutomationLevel(String value) {
        UI.component.value.set(value, DRS_AUT_LEVEL_DD);
    }

    /**
     * Enable/Disable individual virtual machines DRS automation levels.
     *
     * @param value   if the check box should be checked or not
     */
    public void setVmDrsAutomationLevelEnabled(boolean value) {
        if (!UI.component.exists(ENABLE_VM_DRS_AUT_LEVEL_CB)) {
            clickDrsAutomationSection();
        }

        UI.component.value.set(value, ENABLE_VM_DRS_AUT_LEVEL_CB);
    }


    /**
     * Select drs automation level - "Fully Automated" or "Partially Automated"
     * @param autoLevel - PARTIAL_AUTO or FULL_AUTO from the DrsAutoLevel enum
     */
    public void selectDrsAutoLevel(DrsAutoLevel autoLevel) throws Exception {
       switch (autoLevel){
          case PARTIAL_AUTO:
             UI.component.value.set(true, EditClusterDrsSettingsPage.DRS_PARTIALLY_AUTOMATED_LEVEL_RB);
             break;
          case FULL_AUTO:
             UI.component.value.set(true, EditClusterDrsSettingsPage.DRS_FULLY_AUTOMATED_LEVEL_RB);
             break;
          case MANUAL:
              UI.component.value.set(true, EditClusterDrsSettingsPage.DRS_MAN_LEVEL_RB);
              break;
          default:
             throw new IllegalArgumentException("No such automation level.");
       }
    }

    /**
     * Click(expand/collapse) "DRS Policies" section of this view.
     */
    public void clickDrsPoliciesSection() {
        UI.component.click(DRS_POLICIES_SECTION_IMG);
    }

    /**
     * Enable/Disable DRS Policies-> Enforce Even Distribution checkbox.
     *
     * @param value   if the check box should be checked or not
     */
    public void setEnforceEvenDistribution(boolean value) {
        if (!UI.component.exists(DRS_POLICIES_ENFORCE_EVEN_DIST_CB)) {
           clickDrsPoliciesSection();
        }

        UI.component.value.set(value, DRS_POLICIES_ENFORCE_EVEN_DIST_CB);
    }

    /**
     * Enable/Disable DRS Policies-> Enforce Even Distribution checkbox.
     *
     * @param value   if the check box should be checked or not
     */
    public void setConsumedMemory(boolean value) {
        if (!UI.component.exists(DRS_POLICIES_CONSUMED_MEM_CB)) {
           clickDrsPoliciesSection();
        }

        UI.component.value.set(value, DRS_POLICIES_CONSUMED_MEM_CB);
    }

    /**
     * Enable/Disable DRS Policies-> Enforce Even Distribution checkbox.
     *
     * @param value   if the check box should be checked or not
     */
    public void setCpuOvercommitment(boolean value) {
        if (!UI.component.exists(DRS_POLICIES_CPU_OVERCOM_CB)) {
           clickDrsPoliciesSection();
        }

        UI.component.value.set(value, DRS_POLICIES_CPU_OVERCOM_CB);
    }

    /**
     * Set value of CPU value textbox
     *
     * @param cpuValue
     * @throws InterruptedException
     */
   public void setCpuOvercommitmentValue(String cpuValue) throws InterruptedException {
      if (!UI.component.exists(DRS_POLICIES_CPU_OVERCOM_TB)) {
         clickDrsPoliciesSection();
      }
      waitForControlEnabled(IDGroup.toIDGroup(DRS_POLICIES_CPU_OVERCOM_TB));
      UI.component.value.set(cpuValue, DRS_POLICIES_CPU_OVERCOM_TB);
   }
}
