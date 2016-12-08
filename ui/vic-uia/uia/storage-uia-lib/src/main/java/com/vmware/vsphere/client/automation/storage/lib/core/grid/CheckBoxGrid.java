package com.vmware.vsphere.client.automation.storage.lib.core.grid;

import java.text.MessageFormat;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.Property;

public class CheckBoxGrid {

   private static final String GRID_CHECK_BOX_SELECTOR_PATTERN = "{0}/className=AdvancedListBaseContentHolder"
         + "/className=AdvancedListBaseContentHolder/{1}={2}/className=CheckBox";

   private final String gridComponentSelector;
   private final String checkBoxProperty;

   /**
    * Initializes new instance of {@link CheckBoxGrid}
    *
    * @param gridComponentSelector
    *           the id of the grid
    * @param checkBoxProperty
    *           the property of the checkbox element which value corresponds to
    *           key column value
    */
   public CheckBoxGrid(String gridComponentSelector, String checkBoxProperty) {
      this.gridComponentSelector = gridComponentSelector;
      this.checkBoxProperty = checkBoxProperty;
   }

   /**
    * Select row by value of the key column
    *
    * @param keyColumnValue
    */
   public void select(String keyColumnValue) {
      String checkBoxSelector = buildCheckBoxSelector(keyColumnValue);

      SUITA.Factory.UI_AUTOMATION_TOOL.component.click(checkBoxSelector);

      if (!isSelected(keyColumnValue)) {
         throw new RuntimeException(
               String.format(
                     "Failed to interact with checkbox in grid %s for key column value of %s.",
                     this.gridComponentSelector, keyColumnValue));
      }
   }

   /**
    * Returns whether given element is selected
    *
    * @param keyColumnValue
    * @return
    */
   public boolean isSelected(String keyColumnValue) {
      String checkBoxSelector = buildCheckBoxSelector(keyColumnValue);

      return SUITA.Factory.UI_AUTOMATION_TOOL.component.value
            .getBoolean(checkBoxSelector);
   }

   /**
    * Returns whether given element is enabled
    *
    * @param keyColumnValue
    * @return
    */
   public boolean isEnabled(String keyColumnValue) {
      String checkBoxSelector = buildCheckBoxSelector(keyColumnValue);

      return SUITA.Factory.UI_AUTOMATION_TOOL.component.property.getBoolean(
            Property.ENABLED, checkBoxSelector);
   }

   private String buildCheckBoxSelector(String keyColumnValue) {
      String checkBoxSelector = MessageFormat.format(
            GRID_CHECK_BOX_SELECTOR_PATTERN, this.gridComponentSelector,
            this.checkBoxProperty, keyColumnValue);
      return checkBoxSelector;
   }
}
