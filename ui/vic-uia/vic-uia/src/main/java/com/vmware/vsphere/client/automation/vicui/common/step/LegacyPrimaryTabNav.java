package com.vmware.vsphere.client.automation.vicui.common.step;

/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
//package com.vmware.client.automation.components.navigator.navstep;

import java.text.MessageFormat;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.PropertiesUtil;
import com.vmware.client.automation.delay.Delay;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.custom.ViClientPermanentTabBar;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * Selecting a primary tab such as Summary, Manage, and Monitor.
 */
public class LegacyPrimaryTabNav {

   private static final String ID_TAB_NAVIGATOR = "tabNavigator";
   private static final String ID_PATTERN_MAIN_TAB = "tabNavigator/tabBar/label=%s";

   // Queries for finding the Main View index
   private static final String CHILDREN_CLASSNAME_PROPERTY_QUERY = "rawChildren.[].className";
   private static final String CHILDREN_X_PROPERTY_QUERY = "rawChildren.[{0}].[].x";
   private static final String CHILDREN_LABEL_PROPERTY_QUERY = "rawChildren.[{0}].[].label";
   private static final String TABBAR_CLASSNAME = "TabBar";

   /**
    * Selects primary tab by name
    *
    * @param tabName
    */
   public void selectPrimaryTab(String tabName) {
      ViClientPermanentTabBar tabsBar = new ViClientPermanentTabBar(
            ID_TAB_NAVIGATOR, BrowserUtil.flashSelenium);
      tabsBar.waitForElementEnable(Delay.timeout.forSeconds(30).getDuration());
      new BaseView().waitForPageToRefresh();

      String tabSelector = String.format(ID_PATTERN_MAIN_TAB, tabName);
      UIComponent tab = new UIComponent(tabSelector, BrowserUtil.flashSelenium);

      long endTime = System.currentTimeMillis()
            + UiDelay.UI_OPERATION_TIMEOUT.getDuration();

      while (System.currentTimeMillis() < endTime) {
         if (!tab.isVisibleOnPath()) {
            int mainViewIndex = getMainViewIndex(tabName);
            if (mainViewIndex >= 0) {
               tabsBar.selectTabAt(Integer.toString(mainViewIndex),
                     Delay.timeout.forSeconds(10).getDuration());
            }
         }
      }
   }

   /**
    * This method returns by name the index of a Primary tab
    *
    * @param mainViewLabel
    *           : The name of the Primary tab
    * @return
    */
   private int getMainViewIndex(String mainViewLabel) {
      List<String> mainViews = getOrderedMainViews();
      int mainViewIndex = mainViews.indexOf(mainViewLabel);
      return mainViewIndex;
   }

   /**
    * This method extracts and returns an ordered by index List of the Names of
    * Primary Tabs on the current view
    *
    * @return
    */
   private List<String> getOrderedMainViews() {
      int tabbarChildIndex = getTabBarChildIndex();

      final String X_QUERY = MessageFormat.format(CHILDREN_X_PROPERTY_QUERY,
            tabbarChildIndex);
      final String LABEL_QUERY = MessageFormat.format(
            CHILDREN_LABEL_PROPERTY_QUERY, tabbarChildIndex);
      Map<String, String[]> properties = PropertiesUtil.getProperties(
            ID_TAB_NAVIGATOR, X_QUERY, LABEL_QUERY);
      String[] xProperties = properties.get(X_QUERY);
      String[] labelProperties = properties.get(LABEL_QUERY);

      if (xProperties.length != labelProperties.length) {
         throw new RuntimeException("'xProperties' array has length "
               + xProperties.length
               + " different from the length of 'labelProperties' array "
               + labelProperties.length);
      }

      List<Pair<Integer, String>> elements = new ArrayList<Pair<Integer, String>>();

      for (int i = 0; i < labelProperties.length; i++) {
         elements.add(new ImmutablePair<Integer, String>(Integer
               .parseInt(xProperties[i]), labelProperties[i]));
      }

      // Tabs should be sorted by their x properties in order to be able to find
      // the correct index in the UI
      Collections.sort(elements);

      // Return the list of tab names oder by their index
      List<String> result = new ArrayList<String>();
      for (Pair<Integer, String> pair : elements) {
         result.add(pair.getRight());
      }

      return result;
   }

   /**
    * Gets the child index of the TabBar in the rawChildred list of the
    * TabNavigator <br>
    * <b>NOTE:</b> TabBar component isn't in a child-parent hierarchy which the
    * Flex part can access, that's why we should access the component by
    * manually finding it from the closest accessible component (TabNavigator)
    *
    * @return the index of the TabBar
    */
   private int getTabBarChildIndex() {
      String[] classNames = PropertiesUtil.getProperties(ID_TAB_NAVIGATOR,
            CHILDREN_CLASSNAME_PROPERTY_QUERY).get(
            CHILDREN_CLASSNAME_PROPERTY_QUERY);
      for (int i = 0; i < classNames.length; i++) {
         if (classNames[i].equals(TABBAR_CLASSNAME)) {
            return i;
         }
      }

      throw new RuntimeException(TABBAR_CLASSNAME
            + " child isn't available under element with id: "
            + ID_TAB_NAVIGATOR);
   }
}
