/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.view;

import com.vmware.client.automation.common.view.BaseView;

/**
 * Implements the "home" view
 */
public class HomeView extends BaseSearchView {
   private static String ID_TI_SEARCH = "searchInput";
   private static String ID_BTN_SEARCH = "searchButton";
   private static String ID_LINK_FIRST_SEARCH_RESULT = "category_0_link_0";

   /**
    * Sets the search text
    *
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

   /**
    * Clicks on the link of first result.
    *
    * @result - true if link exists and is clicked, false if link doesn't exist
    */
   public boolean clickFirstResult() {
      new BaseView().waitForPageToRefresh();
      if (UI.component.exists(ID_LINK_FIRST_SEARCH_RESULT)) {
         UI.component.click(ID_LINK_FIRST_SEARCH_RESULT);
         return true;
      } else {
         return false;
      }
   }
}
