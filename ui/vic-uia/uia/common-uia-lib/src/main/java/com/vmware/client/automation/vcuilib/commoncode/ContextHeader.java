package com.vmware.client.automation.vcuilib.commoncode;

import java.util.Map;

/**
 * Header interface.
 *
 * NOTE: this interface is a copy of the one from VCUI-QE-LIB
 */
public interface ContextHeader {

   /**
    * Check and uncheck column check boxes in context header
    *
    * @param String[] selectColumns - columns which check boxes to be checked
    * @param String[] deselectColumns - columns which check boxes to be unchecked
    */
   public void selectDeselectColums(String[] selectColumns, String[] deselectColumns);

   /**
    * Click OK button of context header
    *
    */
   public void clickOK();

   /**
    * Click close button of context header
    *
    */
   public void clickClose();

   /**
    * Click link in context header
    *
    */
   public void clickLink();

   /**
    * Get name of the link in context header
    *
    */
   public String getLinkName();


   /**
    * Return true if OK button is enables, otherwise false
    *
    * @return true
    */
   public boolean isOKButtonEnabled();

   /**
    * Return column check boxes state in context header
    *
    * @param String[] columns - columns which check box state to be return
    * @return Map<String, Boolean> - key is column name, value is true
    *         if the check box is checked otherwise false
    */
   public Map<String, Boolean> getColumnsSelectedState(String[] columns);


   /**
    * Get all columns which are displayed in the ContextHeader
    */
   public String[] getAllColumns();
}
