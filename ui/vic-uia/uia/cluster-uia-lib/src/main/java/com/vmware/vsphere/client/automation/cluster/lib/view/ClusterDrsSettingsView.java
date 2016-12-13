/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.cluster.lib.view;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * This class represent the parent cluster DRS settings info view.
 * View Path: cluster -> manage -> settings -> vSphere DRS
 */
public class ClusterDrsSettingsView extends BaseView {
   protected static final Logger _logger = LoggerFactory
         .getLogger(ClusterDrsSettingsView.class);

    private static final String DRS_ENABLED_LABEL =
        CommonUtil.getLocalizedString("manageCluster.dialog.drsEnabled.label.enabled");

    private static final IDGroup DRS_ENABLED_LBL = IDGroup.toIDGroup("titleLabel");
    private static final IDGroup EDIT_DRS_SETTINGS_BUTTON =
            IDGroup.toIDGroup("editDrsButton");
    private static final IDGroup DRS_POLICIES_SECTION_IMG =
          IDGroup.toIDGroup("arrowImagedrsPoliciesBlock");

    private static final IDGroup DRS_POLICIES_EVEN_DISTRIBUTION_LBL =
          IDGroup.toIDGroup("_DrsConfigView_PropertyGridRow6/labelItem/className=UITextField");
    private static final IDGroup DRS_POLICIES_CONSUMED_MEMORY_LBL =
          IDGroup.toIDGroup("_DrsConfigView_PropertyGridRow7/labelItem/className=UITextField");
    private static final IDGroup DRS_POLICIES_CPU_OVERCOMMITMENT_LBL =
          IDGroup.toIDGroup("_DrsConfigView_PropertyGridRow8/labelItem/className=UITextField");
    private static final IDGroup DRS_POLICIES_EVEN_DISTRIBUTION_TXT =
          IDGroup.toIDGroup("_DrsConfigView_Text4");
    private static final IDGroup DRS_POLICIES_CONSUMED_MEMORY_TXT =
          IDGroup.toIDGroup("_DrsConfigView_Text5");
    private static final IDGroup DRS_POLICIES_CPU_OVERCOMMITMENT_TXT =
          IDGroup.toIDGroup("_DrsConfigView_Text6");
    private static final IDGroup ADVANCED_OPTION_SECTION_IMG =
          IDGroup.toIDGroup("arrowImagedrsAdvOptionsBlock");

    private static final IDGroup ADVANCED_OPTIONS_GRID =
          IDGroup.toIDGroup("advancedOptionsGrid");
    /**
     * Gets the text of DRS enabled label.
     */
    public String getDrsEnabled() {
        return UI.component.property.get(Property.TEXT, DRS_ENABLED_LBL);
    }

    /**
     * Checks if DRS is enabled based on the displayed label.
     */
    public boolean isDrsEnabled() {
        return getDrsEnabled().equals(DRS_ENABLED_LABEL) ? true : false;
    }

    /**
     * Click on edit DRS settings button.
     */
    public void clickEditDrsSettingsButton() {
        UI.component.click(EDIT_DRS_SETTINGS_BUTTON);
    }
    /**
     * Click(expand/collapse) "DRS Policies" section of this view.
     */
    public void clickDrsPoliciesSection() {
        UI.component.click(DRS_POLICIES_SECTION_IMG);
    }
    /**
     * Click(expand/collapse) "Advanced Options" section of this view.
     */
    public void clickAdvancedOptionsSection() {
        UI.component.click(ADVANCED_OPTION_SECTION_IMG);
    }

    /**
     * Expand "DRS Policies" section of this view.
     */
    public void expandDrsPoliciesSection() {
       if (!UI.component.exists(DRS_POLICIES_EVEN_DISTRIBUTION_TXT)) {
          clickDrsPoliciesSection();
      }
    }
     /**
      * Expand "Advanced Options" section of this view.
      */
     public void expandAdvancedOptionsSection() {
       if (!UI.component.exists(ADVANCED_OPTIONS_GRID)) {
          clickAdvancedOptionsSection();
       }
    }

     /**
      * Gets the text of DRS Policies->Even Dist label.
      */
     public String getDrsPoliciesEvenDistributionLabelText() {
         return UI.component.property.get(Property.TEXT, DRS_POLICIES_EVEN_DISTRIBUTION_LBL);
     }
     /**
      * Gets the text of DRS Policies->Consumed Mem label.
      */
     public String getDrsPoliciesConsumedMemoryLabelText() {
         return UI.component.property.get(Property.TEXT, DRS_POLICIES_CONSUMED_MEMORY_LBL);
     }
     /**
      * Gets the text of DRS Policies->CPU Overcom label.
      */
     public String getDrsPoliciesCpuOvercommitmentLabelText() {
         return UI.component.property.get(Property.TEXT, DRS_POLICIES_CPU_OVERCOMMITMENT_LBL);
     }

    /**
     * Gets the text to the right of of DRS Policies->Even Dist label.
     */
    public String getDrsPoliciesEvenDistributionText() {
        return UI.component.property.get(Property.TEXT, DRS_POLICIES_EVEN_DISTRIBUTION_TXT);
    }
    /**
     * Gets the text to the right of of DRS Policies->Consumed Mem label.
     */
    public String getDrsPoliciesConsumedMemoryText() {
        return UI.component.property.get(Property.TEXT, DRS_POLICIES_CONSUMED_MEMORY_TXT);
    }
    /**
     * Gets the text to the right of DRS Policies->CPU Overcom label.
     */
    public String getDrsPoliciesCpuOvercommitmentText() {
        return UI.component.property.get(Property.TEXT, DRS_POLICIES_CPU_OVERCOMMITMENT_TXT);
    }

    /**
     * Finds and returns the advanced data grid on 'Advanced Options' section in the view stack.
     */
    private AdvancedDataGrid getGrid() {
       return GridControl.findGrid(IDGroup.toIDGroup(ADVANCED_OPTIONS_GRID));
    }

    /**
     * Checks whether advanced options item is listed in the data grid
     *
     * @param option
     *           Name of the option item to be searched for
     *
     * @return true if the option item is found in the grid, false otherwise
     */
    public boolean isOptionFoundInGrid(String option) {
       int row = GridControl.getEntityIndex(getGrid(), option);
       _logger.info("Option: [" + option + "] found in row index:[" + row + "]");
       return row >= 0;
    }

}
