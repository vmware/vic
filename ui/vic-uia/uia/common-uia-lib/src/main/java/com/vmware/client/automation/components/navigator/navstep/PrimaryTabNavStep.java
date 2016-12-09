/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.navstep;

import java.text.MessageFormat;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.PropertiesUtil;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.delay.Delay;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.custom.ViClientPermanentTabBar;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * A <code>NavigationStep</code> for selecting a primary tab such as Summary,
 * Manage, and Monitor.
 */
public class PrimaryTabNavStep extends BaseNavStep {

   private static final String ID_TAB_NAVIGATOR = "tabNavigator";
   private static final String ID_PATTERN_MAIN_TAB = "tabNavigator/tabBar/label=%s";
   private int _tabPosition = -1;
   private String _tabName = "";

   // Queries for finding the Main View index
   private static final String CHILDREN_CLASSNAME_PROPERTY_QUERY = "rawChildren.[].className";
   private static final String CHILDREN_X_PROPERTY_QUERY = "rawChildren.[{0}].[].x";
   private static final String CHILDREN_LABEL_PROPERTY_QUERY = "rawChildren.[{0}].[].label";
   private static final String TABBAR_CLASSNAME = "TabBar";

   /**
    * The constructor defines a mapping between the navigation identifier and
    * the tab position. We deprecated the method because we would like to make
    * the tab navigation by name and not by index as navigation by index is not
    * always reliable. We do not plan to deprecate the class.
    *
    * @param nId
    * @see <code>NavigationStep</code>.
    *
    * @param tabPosition
    *           Zero based position of the tab.
    *
    * @deprecated Please use the following constructor:
    *             <code>PrimaryTabNavStep(String nid, String tabName)</code>
    */
   @Deprecated
   public PrimaryTabNavStep(String nid, int tabPosition) {
      super(nid);

      _tabPosition = tabPosition;
   }

   /**
    * The constructor defines a mapping between the navigation identifier and
    * the tab name.
    *
    * @param nId
    * @see <code>NavigationStep</code>.
    *
    * @param tabName
    *           Tab name.
    */
   public PrimaryTabNavStep(String nid, String tabName) {
      super(nid);

      _tabName = tabName;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) {
      SUITA.Factory.UI_AUTOMATION_TOOL.audit.snapshotAppScreen(
            SubToolAudit.getFPID(), "PRE_CLICK_PRIMARY_TAB_INDEX_"
                  + _tabPosition);

      if (!Strings.isNullOrEmpty(_tabName)) {
         selectView(_tabName);
      } else if (_tabPosition >= 0) {
         // TODO: this should be removed, once the deprecated constructor is
         // deleted
         new BaseView().waitForPageToRefresh();
         selectView(getOrderedMainViews().get(_tabPosition));
      } else {
         throw new IllegalArgumentException(
               "tabPosition and tabName cannot be both undefined at the same time");
      }
   }

   /**
    * Selects entity's main view by name
    *
    * @param viewName
    */
   private void selectView(String viewName) {
      ViClientPermanentTabBar tabsBar = new ViClientPermanentTabBar(
            ID_TAB_NAVIGATOR, BrowserUtil.flashSelenium);
      tabsBar.waitForElementEnable(Delay.timeout.forSeconds(30).getDuration());
      new BaseView().waitForPageToRefresh();

      String tabSelector = String.format(ID_PATTERN_MAIN_TAB, viewName);
      UIComponent tab = new UIComponent(tabSelector, BrowserUtil.flashSelenium);

      String tabSelectorConfig = String.format(ID_PATTERN_MAIN_TAB,
            CommonUtil.getLocalizedString("entity.tab.primary.configure"));
      long endTime = System.currentTimeMillis()
            + UiDelay.UI_OPERATION_TIMEOUT.getDuration();

      while (System.currentTimeMillis() < endTime) {
         if (tab.isVisibleOnPath()) {
            tab.leftMouseClick();
            return;
         }
         // Check if 'Manage' was renamed to 'Configure' and try to click it 
         if (viewName.equals(CommonUtil
               .getLocalizedString("entity.tab.primary.manage"))
               && UI.component.exists(tabSelectorConfig)) {
            UI.component.click(tabSelectorConfig);
            return;
         }
      }
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

      List<Pair<Integer, String>> elements = new ArrayList<>();

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
