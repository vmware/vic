/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.view;


/**
 * Implements the "simple search" view
 */
public class SimpleSearchView extends BaseSearchView {
   private static String ID_TI_SEARCH = "simpleSearch/searchInput";
   private static String ID_BTN_SEARCH = "simpleSearch/searchButton";

   /**
    * Sets the search text
    * @param searchText - the search text
    */
   public void setSearchText(String searchText) {
      UI.component.value.set(searchText, ID_TI_SEARCH);
   }

   /**
    * Clicks the "Search" button
    */
   public void clickSearchButton() {
      UI.component.click(ID_BTN_SEARCH);
   }
}
